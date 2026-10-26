package job

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/transfer_customer"
	"micro-scrm/internal/orm/service"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/workwechat"
	"micro-scrm/internal/pkg/xjob"
)

// 在职继承任务运行Job

type onTransferCustomerTasksRunningJob struct {
	name string
}

func NewOnTransferCustomerTasksRunningJob() *onTransferCustomerTasksRunningJob {
	return &onTransferCustomerTasksRunningJob{
		name: fmt.Sprintf("%s::在职继承任务运行Job", config.CfgData.App),
	}
}

func (job onTransferCustomerTasksRunningJob) Run(ctx ctx.Context, jobHandler *xjob.JobHandler) error {
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
			go newOnTransferCustomer(job.name, enterpriseInfo).exec(ctx, done)

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

// ===================================================================================================

type onTransferCustomer struct {
	name           string
	enterpriseInfo boss.EnterpriseModel
	message        micro_scrm.TransferCustomerSuccessMessageModel
}

func newOnTransferCustomer(name string, enterpriseInfo boss.EnterpriseModel) *onTransferCustomer {
	obj := &onTransferCustomer{
		name:           fmt.Sprintf("%s,企业:%s,企业ID:%d", name, enterpriseInfo.Name, enterpriseInfo.ID),
		enterpriseInfo: enterpriseInfo,
	}

	// 获取继承消息
	orm.Service.TransferCustomer.Where("enterprise_id = ?", enterpriseInfo.ID).Take(&obj.message)
	return obj
}

func (task onTransferCustomer) exec(ctx ctx.Context, done chan struct{}) {
	// 获取启动的任务
	taskList, taskListErr := orm.Service.TransferCustomer.TasksGetList(task.enterpriseInfo.ID, 0, 0,
		transfer_customer.TCT_STATUS_RUNNING,
		service.PaginateParam{PageSize: 99999999, Page: 1})
	if taskListErr != nil {
		log.Info(ctx).Msgf("%s,获取该企业在职继承任务失败:%v", task.name, taskListErr)
		return
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
			go task.listenTaskCancel(ctx, taskDone, item.ID) // 监测任务是否取消
			go task.operation(ctx, taskCtx, item)            // 拉取任务的数据

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

// 监听任务取消
// ctx 系统默认上下文； done 退出标志；taskID 执行的任务ID
func (task onTransferCustomer) listenTaskCancel(ctx ctx.Context, done chan struct{}, taskID int64) {
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

		var count int64
		orm.Service.TransferCustomer.
			Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, task.enterpriseInfo.ID)).
			Model(&transfer_customer.TransferCustomerDataModel{}).
			Where("task_id = ?", taskID).
			Where("transfer_customer_status = ? AND transfer_customer_type!=?", 0, 2).Count(&count)

		// 拉取停止结束
		// 1、任务取消；2、拉取数据结束
		if taskInfo.Status == transfer_customer.TCT_STATUS_CANCEL || count == 0 {
			done <- struct{}{}
			break
		}
	}
}

// 执行客户继承
// ctx 系统默认上下文； taskCtx 该任务操作的上下文；taskInfo 执行的任务信息集合
func (task onTransferCustomer) operation(ctx ctx.Context, taskCtx context.Context, taskInfo transfer_customer.TransferCustomerTasksModel) {
	// 获取需要继承的客户
	customerListChan := make(chan string, EXTRACT_CUSTOMER_NUM)
	defer close(customerListChan)
	go func(dataListChan chan string) {
		page := 1
		for {
			unionIDList := []string{}
			if err := orm.Service.TransferCustomer.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, task.enterpriseInfo.ID)).
				Model(&transfer_customer.TransferCustomerDataModel{}).
				Distinct().Offset((page-1)*EXTRACT_CUSTOMER_NUM).Where(map[string]interface{}{
				"task_id":                  taskInfo.ID,
				"transfer_customer_status": transfer_customer.TCD_STATUS_NO,
			}).
				Limit(EXTRACT_CUSTOMER_NUM).Pluck("external_unionid", &unionIDList).Error; err != nil {
				fmt.Println(err)
				continue
			}

			if len(unionIDList) == 0 {
				break
			}

			for _, unionID := range unionIDList {
				dataListChan <- unionID
			}

			page++
		}
	}(customerListChan)

	wwObj, wwObjErr := workwechat.NewWorkWeChat(ctx, task.enterpriseInfo.ID)
	if wwObjErr != nil {
		log.Info(ctx).Msgf("%s::客户在职继承获取企微对象失败:%v", task.name, wwObjErr)
		return
	}
	for {
		select {
		case <-taskCtx.Done():
			return
		case unionID := <-customerListChan:
			dataList := []transfer_customer.TransferCustomerDataModel{}
			if err := orm.Service.TransferCustomer.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, task.enterpriseInfo.ID)).
				Where("external_unionid = ?", unionID).
				Where("task_id = ?", taskInfo.ID).
				Find(&dataList).Error; err != nil {
				log.Info(ctx).Msgf("%s::客户在职继承已成:%v", task.name, err)
				continue
			}

			switch taskInfo.Type {
			case transfer_customer.TCT_Type_Only:
				task.transferOnly(ctx, wwObj, taskInfo, dataList)
			case transfer_customer.TCT_Type_Random:
				task.transferRandom(ctx, wwObj, taskInfo, dataList)
			}
		}
	}
}

// 把客户的关系继承到一个管理BA上
func (task onTransferCustomer) transferOnly(ctx ctx.Context, ww *workwechat.WorkWeChat, taskInfo transfer_customer.TransferCustomerTasksModel, dataList []transfer_customer.TransferCustomerDataModel) {
	for _, data := range dataList {
		response, responseErr := ww.OnTransfer_AllotCustomer(work.TransferCustomerAllotRequest{
			HandoverUserID:     data.HandoverUserid,
			TakeoverUserID:     data.TakeoverUserid,
			ExternalUserID:     []string{data.ExternalUserid},
			TransferSuccessMsg: task.message.Content,
		})
		if responseErr != nil {
			log.Info(ctx).Msgf("%s::客户在职继承失败:%v", task.name, responseErr)
			orm.Service.TransferCustomer.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, task.enterpriseInfo.ID)).
				Where("id = ?", data.ID).
				Updates(map[string]interface{}{
					"transfer_customer_status": transfer_customer.TCD_STATUS_FAIL,
					"transfer_customer_type":   transfer_customer.TCD_TYPE_Only,
				})
		} else {
			if response.ErrCode == 0 {
				orm.Service.TransferCustomer.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, task.enterpriseInfo.ID)).
					Where("id = ?", data.ID).
					Updates(map[string]interface{}{
						"transfer_customer_status": transfer_customer.TCD_STATUS_Wait,
						"transfer_customer_type":   transfer_customer.TCD_TYPE_Only,
					})
			} else {
				orm.Service.TransferCustomer.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, task.enterpriseInfo.ID)).
					Where("id = ?", data.ID).
					Updates(map[string]interface{}{
						"transfer_customer_status": transfer_customer.TCD_STATUS_FAIL,
						"transfer_customer_type":   transfer_customer.TCD_TYPE_Only,
					})
			}
		}
	}
}

// 随机把客户的一个关系继承到一个管理BA上
func (task onTransferCustomer) transferRandom(ctx ctx.Context, ww *workwechat.WorkWeChat, taskInfo transfer_customer.TransferCustomerTasksModel, dataList []transfer_customer.TransferCustomerDataModel) {
	if len(dataList) == 0 {
		return
	}

	randNum := rand.Intn(len(dataList))
	response, responseErr := ww.OnTransfer_AllotCustomer(work.TransferCustomerAllotRequest{
		HandoverUserID:     dataList[randNum].HandoverUserid,
		TakeoverUserID:     dataList[randNum].TakeoverUserid,
		ExternalUserID:     []string{dataList[randNum].ExternalUserid},
		TransferSuccessMsg: task.message.Content,
	})
	if responseErr != nil {
		log.Info(ctx).Msgf("%s::客户在职继承失败:%v", task.name, responseErr)
		orm.Service.TransferCustomer.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, task.enterpriseInfo.ID)).
			Where("id = ?", dataList[randNum].ID).
			Updates(map[string]interface{}{
				"transfer_customer_status": transfer_customer.TCD_STATUS_FAIL,
				"transfer_customer_type":   transfer_customer.TCD_TYPE_Only,
			})
	} else {
		if response.ErrCode == 0 {
			orm.Service.TransferCustomer.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, task.enterpriseInfo.ID)).
				Where("id = ?", dataList[randNum].ID).
				Updates(map[string]interface{}{
					"transfer_customer_status": transfer_customer.TCD_STATUS_Wait,
					"transfer_customer_type":   transfer_customer.TCD_TYPE_Only,
				})
		} else {
			orm.Service.TransferCustomer.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, task.enterpriseInfo.ID)).
				Where("id = ?", dataList[randNum].ID).
				Updates(map[string]interface{}{
					"transfer_customer_status": transfer_customer.TCD_STATUS_FAIL,
					"transfer_customer_type":   transfer_customer.TCD_TYPE_Only,
				})
		}
	}

	for index, data := range dataList {
		if index == randNum {
			continue
		}

		orm.Service.TransferCustomer.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, task.enterpriseInfo.ID)).
			Where("id = ?", data.ID).
			Updates(map[string]interface{}{
				"transfer_customer_type": transfer_customer.TCD_TYPE_Random,
			})
	}
}
