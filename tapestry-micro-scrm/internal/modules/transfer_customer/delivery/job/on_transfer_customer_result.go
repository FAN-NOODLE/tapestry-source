package job

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/transfer_customer"
	"micro-scrm/internal/orm/service"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/workwechat"
	"micro-scrm/internal/pkg/xjob"
)

// 在职继承任务运行Job

type onTransferResultJob struct {
	name string
}

func NewOnTransferResultJob() *onTransferResultJob {
	return &onTransferResultJob{
		name: fmt.Sprintf("%s::获取在职继承任务运行结果Job", config.CfgData.App),
	}
}
func (job onTransferResultJob) Run(ctx ctx.Context, jobHandler *xjob.JobHandler) error {
	fmt.Printf("%s,任务执行Running...\n", job.name)

	// 获取所有品牌的企业ID
	enterpriseInfoList, enterpriseInfoErr := orm.BossService.Enterprise.GetAllInfo()
	if enterpriseInfoErr != nil {
		log.Info(ctx).Msgf("%s::获取所有品牌的信息时报错，错误为：%v", job.name, enterpriseInfoErr)
		return enterpriseInfoErr
	}

	var wg sync.WaitGroup
	for _, enterpriseInfo := range enterpriseInfoList {
		wg.Add(1)
		go func(enterpriseInfo boss.EnterpriseModel) {
			defer wg.Done()
			done := make(chan struct{})
			defer close(done)
			go newOnTransferCustomer(job.name, enterpriseInfo).execTransferResult(ctx, done)

		Label:
			for {
				select {
				case <-done:
					break Label
				default:
				}
			}
		}(enterpriseInfo)
	}
	wg.Wait()

	fmt.Printf("%s,任务执行结束\n", job.name)

	return nil

}
func (job onTransferCustomer) execTransferResult(ctx ctx.Context, done chan struct{}) error {
	// 获取启动的任务
	taskList, taskListErr := orm.Service.TransferCustomer.TasksPullDataStatus(job.enterpriseInfo.ID,
		service.PaginateParam{PageSize: 99999999, Page: 1})
	if taskListErr != nil {
		log.Info(ctx).Msgf("%s,获取该企业在职继承任务失败:%v", job.name, taskListErr)
		return taskListErr
	}

	// 开始品牌的在职继承
	var wg sync.WaitGroup
	for _, item := range taskList {
		// 如果没有把数据拉取完成，则不进行客户继承
		if item.PullDataStatus != transfer_customer.TCT_PullDataStatus_Yes {
			continue
		}

		wg.Add(1)
		go func(item transfer_customer.TransferCustomerTasksModel) {
			defer wg.Done()
			// 实时监测，任务是否在拉取继承数据时取消
			taskDone := make(chan struct{})
			defer close(taskDone)
			taskCtx, taskCancel := context.WithCancel(context.Background())
			go job.listenTaskStatus(ctx, taskDone, item.ID)                // 监测任务是否取消
			go job.onTransferResultStatus(ctx, taskCtx, taskDone, item.ID) // 查询在职继承的状态

		Label:
			for {
				select {
				case <-taskDone:
					taskCancel()
					time.Sleep(1 * time.Second)
					break Label
				default:
				}
			}
		}(item)
	}
	wg.Wait()
	done <- struct{}{}
	return nil
}

// 执行客户继承
// ctx 系统默认上下文； taskCtx 该任务操作的上下文；taskInfo 执行的任务信息集合
func (task onTransferCustomer) onTransferResultStatus(ctx ctx.Context, taskCtx context.Context, done chan struct{}, taskId int64) {
	receiveChan := make(chan *transfer_customer.TransferCustomerDataModel, 5000)
	go GetTransferResultData(receiveChan, task.enterpriseInfo.ID, taskId)
	var wg sync.WaitGroup
	for v := range receiveChan {
		wg.Add(1)
		go func(p *transfer_customer.TransferCustomerDataModel) {
			defer wg.Done()
			params := work.TransferResultRequst{
				HandoverUserid: p.HandoverUserid,
				TakeoverUserid: p.TakeoverUserid,
			}
			if requst, err := workwechat.NewWorkWeChat(ctx, task.enterpriseInfo.ID); err == nil {
				if repone, err := requst.OnTransfer_QueryCustomerStatus(params); err == nil {
					for _, value := range repone.Customer {
						if value.ExternalUserid == p.ExternalUserid {
							if value.Status == 1 || value.Status == 3 || value.Status == 4 {
								p.TransferCustomerStatus = value.Status
								transfercustomerlog := transfer_customer.TransferCustomerLogModel{
									HandoverUserid:         p.HandoverUserid,
									TakeoverUserid:         p.TakeoverUserid,
									ExternalUserid:         p.ExternalUserid,
									TransferSuccessMsgId:   p.TransferSuccessMsgId,
									TransferCustomerStatus: p.TransferCustomerStatus,
									TakeoverTime:           value.TakeoverTime,
									TaskId:                 taskId,
									Errcode:                repone.ErrCode,
									Errmsg:                 repone.ErrMsg,
									TransferCustomerDataId: p.ID,
								}
								if err := orm.Service.TransferCustomer.UpdatetransfercustomerStatus(task.enterpriseInfo.ID, p, transfercustomerlog); err != nil {
									log.Info(ctx).Msgf("%s::监测任务修改在职继承人员根据数据编号ID:%d,查询任务信息失败:%v", task.name, p.ID, err)
								}
							}
						}
					}
				} else {
					if err := orm.Service.TransferCustomer.Deletetransfercustomer(task.enterpriseInfo.ID, p); err != nil {
						log.Info(ctx).Msgf("%s::修改任务状态:%d,查询任务信息失败:%v", task.name, taskId, err)
					}
				}
			}

		}(v)
	}
	wg.Wait()
	// 查询任务下是否还有待接替的用户如果没有表示任务已经结束。
	if count, err := orm.Service.TransferCustomer.FindNeedTransferrResultDateCount(task.enterpriseInfo.ID, taskId); err == nil && count == 0 {
		// 表示数据已经查询结束 更新计划任务状态
		orm.Service.TransferCustomer.UpdateTaskStatus(task.enterpriseInfo.ID, taskId)
	}
	done <- struct{}{}
}

func GetTransferResultData(receiveChan chan *transfer_customer.TransferCustomerDataModel, enterpriseID, taskId int64) {
	paginate := service.PaginateParam{
		PageSize: 5000,
	}
	for {
		paginate.Page++
		datas, err := orm.Service.TransferCustomer.FindNeedTransferrResultDate(enterpriseID, taskId, paginate)
		if err != nil {
			close(receiveChan)
			break
		}
		for _, v := range datas {
			receiveChan <- v
		}
		if len(datas) != paginate.PageSize {
			close(receiveChan)
			break
		}
	}
}
func (task onTransferCustomer) listenTaskStatus(ctx ctx.Context, done chan struct{}, taskID int64) {
	for {
		time.Sleep(TASK_LISTEN_INTERVAL)
		taskInfo, taskInfoErr := orm.Service.TransferCustomer.TasksGet(task.enterpriseInfo.ID, transfer_customer.TransferCustomerTasksModel{
			BranchModel: model.BranchModel{
				BaseModel: model.BaseModel{
					ID: taskID,
				},
			},
		})
		if taskInfoErr != nil {
			log.Info(ctx).Msgf("%s::监测任务退出功能根据ID:%d,查询任务信息失败:%v", task.name, taskInfo.ID, taskInfoErr)
			continue
		}
		// 拉取停止结束
		// 任务已经在其他状态
		if taskInfo.Status != transfer_customer.TCT_STATUS_RUNNING {
			done <- struct{}{}
			break
		}
	}
}
