package job

import (
	"context"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/transfer_customer"
	"micro-scrm/internal/orm/service"
	"micro-scrm/internal/pkg/xjob"
	"sync"
	"time"
)

// 在职继承任务启动Job

type onTransferCustomerTasksStartJob struct {
	name string
}

func NewOnTransferCustomerTasksStartJob() *onTransferCustomerTasksStartJob {
	return &onTransferCustomerTasksStartJob{
		name: fmt.Sprintf("%s::在职继承任务启动Job", config.CfgData.App),
	}
}

func (job onTransferCustomerTasksStartJob) Run(ctx ctx.Context, jobHandler *xjob.JobHandler) error {
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
			go newOnTransferCustomerTasksStart(enterpriseInfo, job.name).exec(ctx, done)

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

// =====================================================================================================================

type onTransferCustomerTasksStart struct {
	enterpriseInfo boss.EnterpriseModel
	name           string
}

func newOnTransferCustomerTasksStart(enterpriseInfo boss.EnterpriseModel, name string) *onTransferCustomerTasksStart {
	return &onTransferCustomerTasksStart{
		enterpriseInfo: enterpriseInfo,
		name:           fmt.Sprintf("%s::企业名称:%s,企业ID:%d", name, enterpriseInfo.Name, enterpriseInfo.ID),
	}
}

func (task onTransferCustomerTasksStart) exec(ctx ctx.Context, done chan struct{}) {
	// 获取启动的任务
	taskList, taskListErr := orm.Service.TransferCustomer.TasksGetList(task.enterpriseInfo.ID, 0, 0,
		transfer_customer.TCT_STATUS_RUNNING,
		service.PaginateParam{PageSize: 99999999, Page: 1})
	if taskListErr != nil {
		log.Info(ctx).Msgf("%s,获取该企业在职继承任务失败:%v", task.name, taskListErr)
		return
	}

	// 获取在职继承的数据
	var wg sync.WaitGroup
	for _, item := range taskList {
		// 如果任务已经把数据拉取完成，则不需要进行拉取数据
		if item.PullDataStatus == transfer_customer.TCT_PullDataStatus_Yes {
			continue
		}

		wg.Add(1)
		go func(item transfer_customer.TransferCustomerTasksModel) {
			defer wg.Done()
			// 实时监测，任务是否在拉取继承数据时取消
			taskDone := make(chan struct{})
			defer close(taskDone)
			taskCtx, taskCancel := context.WithCancel(context.Background())
			go task.listenTaskCancel(ctx, taskDone, item.ID) // 监测任务是否取消
			go task.pullTaskData(ctx, taskCtx, item.ID)      // 拉取任务的数据

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
}

// 监测任务退出
// 每隔一段时间检测一次任务是否取消/完成
func (task onTransferCustomerTasksStart) listenTaskCancel(ctx ctx.Context, done chan struct{}, taskID int64) {
	for {
		taskInfo, taskInfoErr := orm.Service.TransferCustomer.TasksGet(task.enterpriseInfo.ID, transfer_customer.TransferCustomerTasksModel{
			BranchModel: model.BranchModel{
				BaseModel: model.BaseModel{
					ID: taskID,
				},
			},
		})
		if taskInfoErr != nil {
			log.Info(ctx).Msgf("%s::监测任务退出功能根据ID:%d,查询任务信息失败:%v", task.name, taskID, taskInfoErr)
			continue
		}

		// 拉取停止结束
		// 1、任务取消；2、拉取数据结束
		if taskInfo.Status == transfer_customer.TCT_STATUS_CANCEL ||
			taskInfo.PullDataStatus == transfer_customer.TCT_PullDataStatus_Yes {
			done <- struct{}{}
			break
		}
		time.Sleep(TASK_LISTEN_INTERVAL)
	}
}

// 拉取任务数据
func (task onTransferCustomerTasksStart) pullTaskData(ctx ctx.Context, taskContext context.Context, taskID int64) {
	// 初始化页码
	page := 1
	for {
		select {
		case <-taskContext.Done():
			// 主动进入的为，取消任务
			if updateErr := orm.Service.TransferCustomer.TasksUpdate(task.enterpriseInfo.ID, transfer_customer.TransferCustomerTasksModel{
				BranchModel: model.BranchModel{
					BaseModel: model.BaseModel{
						ID: taskID,
					},
				},
				Status: transfer_customer.TCT_STATUS_CANCEL,
			}); updateErr != nil {
				log.Info(ctx).Msgf("%s::任务取消失败,任务ID：%d", task.name, taskID)
			}
			log.Info(ctx).Msgf("%s::任务ID:%d,被取消", task.name, taskID)
			return
		default:
			// 插入完成后，更新状态
			dataList, dataListErr := orm.Service.TransferCustomer.GetTransferDataList(task.enterpriseInfo.ID,
				service.PaginateParam{
					Page: page, PageSize: PAGE_SIZE,
				})
			if dataListErr != nil {
				log.Info(ctx).Msgf("%s,获取该企业在职继承任务数据失败:%v", task.name, dataListErr)
				break
			}

			if len(dataList) == 0 {
				var count int64
				if countErr := orm.Service.TransferCustomer.
					Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, task.enterpriseInfo.ID)).
					Model(&transfer_customer.TransferCustomerDataModel{}).Distinct("external_unionid").
					Where("task_id = ?", taskID).
					Count(&count).Error; countErr != nil {
					log.Info(ctx).Msgf("%s,获取该企业在职继承任务数据总数失败:%v", task.name, countErr)
				} else {
					// 更新任务拉取数据完成状态
					if updateErr := orm.Service.TransferCustomer.TasksUpdate(task.enterpriseInfo.ID, transfer_customer.TransferCustomerTasksModel{
						BranchModel: model.BranchModel{
							BaseModel: model.BaseModel{
								ID: taskID,
							},
						},
						CustomerNum:    int(count),
						PullDataStatus: transfer_customer.TCT_PullDataStatus_Yes,
					}); updateErr == nil {
						return
					}
				}
			}

			insertData := []*transfer_customer.TransferCustomerDataModel{}
			for _, data := range dataList {
				insertData = append(insertData, &transfer_customer.TransferCustomerDataModel{
					HandoverUserid:  data.HandoverUserID,
					TakeoverUserid:  data.TakeoverUserID,
					ExternalUserid:  data.WorkExternalUserID,
					ExternalUnionID: data.UnionID,
					TaskId:          taskID,
				})
			}
			if insertErr := orm.Service.TransferCustomer.InsertTransferCustomerData(task.enterpriseInfo.ID, insertData); insertErr != nil {
				log.Info(ctx).Msgf("%s,获取该企业在职继承数据数据插入时失败:%v", task.name, insertErr)
				break
			}

			// 页码递增
			page++
		}
	}
}
