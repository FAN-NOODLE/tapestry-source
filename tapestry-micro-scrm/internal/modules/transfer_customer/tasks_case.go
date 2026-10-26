package transfer_customer

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/transfer_customer"
	"micro-scrm/internal/pkg/helper"
	"sync"
	"time"
)

var casing *tasksCase

type tasksCase struct {
	Name string
}

var once sync.Once

func NewTasksCase() *tasksCase {
	once.Do(func() {
		casing = &tasksCase{
			Name: "在职继承任务逻辑层",
		}
	})

	return casing
}

// Report 查看任务的报表
func (casing tasksCase) Report(ctx ctx.Context, param TasksReportRequest) (data TasksReportResponse, err error) {
	// 响应填充请求参数
	data.Request = param

	// 请求数据头
	if !param.IsOnlyList {
		if taskInfo, taskInfoErr := orm.Service.TransferCustomer.TasksGet(param.EnterpriseID, transfer_customer.TransferCustomerTasksModel{
			BranchModel: model.BranchModel{
				BaseModel: model.BaseModel{
					ID: param.ID,
				},
			},
		}); taskInfoErr == nil {
			data.TaskInfo = taskInfo
		} else {
			err = errors.New("该任务报表信息不存在")
			return
		}

		if titleData, titleErr := orm.Service.TransferCustomer.TasksSummaryData(param.EnterpriseID, param.ID); titleErr == nil {
			data.ReportTitle.CustomerTotalNum = titleData.CustomerTotalNum
			data.ReportTitle.CustomerSuccessNum = titleData.CustomerSuccessNum
			data.ReportTitle.CustomerFailNum = titleData.CustomerFailNum
			data.ReportTitle.FinishBANum = titleData.FinishBANum
			fmt.Println(float64(titleData.CustomerSuccessNum), float64(titleData.CustomerTotalNum))

			if titleData.CustomerSuccessNum == 0 || titleData.CustomerTotalNum == 0 {
				data.ReportTitle.FinishRate = 0
			} else {
				data.ReportTitle.FinishRate = int(float64(titleData.CustomerSuccessNum) / float64(titleData.CustomerTotalNum) * 100)
			}

		} else {
			err = titleErr
			return
		}
	}

	// 响应填充客户列表
	if dataList, listErr := orm.Service.TransferCustomer.TasksTransferCustomerList(param.EnterpriseID, param.ID, param.PaginateParam); listErr == nil {
		data.ReportBody.DataList = []TasksReportListData{}
		for _, item := range dataList {
			data.ReportBody.DataList = append(data.ReportBody.DataList, TasksReportListData{
				Name:               item.Name,
				FriendNum:          item.FriendNum,
				CurrentAdminBAName: item.CurrentAdminBAName,
				BADepartmentName:   item.BADepartmentName,
				TransferStatus:     item.TransferStatus,
				TransferStatusMsg:  item.TransferStatusMsg,
				FailMsg:            item.FailMsg,
			})
		}

		data.ReportBody.PageSize = param.PageSize
		data.ReportBody.Page = param.Page
	} else {
		err = listErr
		return
	}
	return
}

// Delete 删除任务
func (casing tasksCase) Delete(ctx ctx.Context, param TasksDeleteRequest) error {
	data, err := orm.Service.TransferCustomer.TasksGet(param.EnterpriseID, transfer_customer.TransferCustomerTasksModel{
		BranchModel: model.BranchModel{
			BaseModel: model.BaseModel{
				ID: param.ID,
			},
		},
	})
	if err != nil {
		return err
	}

	if data.Status != transfer_customer.TCT_STATUS_NOT {
		return errors.New("任务已启动，不能删除")
	}

	data.DeletedAt = time.Now().Unix()
	return orm.Service.TransferCustomer.TasksUpdate(param.EnterpriseID, data)
}

// Update 更新任务
func (casing tasksCase) Update(ctx ctx.Context, param TasksUpdateRequest) error {
	data, err := orm.Service.TransferCustomer.TasksGet(param.EnterpriseID, transfer_customer.TransferCustomerTasksModel{
		BranchModel: model.BranchModel{
			BaseModel: model.BaseModel{
				ID: param.ID,
			},
		},
	})
	if err != nil {
		return err
	}

	if param.Name != "" {
		data.Name = param.Name
	}

	if param.TransferType != 0 {
		data.Type = param.TransferType
	}

	if param.Status != 0 {
		data.Status = param.Status
		if param.Status == transfer_customer.TCT_STATUS_RUNNING {
			data.StartTime = time.Now().Unix()
		}
	}

	return orm.Service.TransferCustomer.TasksUpdate(param.EnterpriseID, data)
}

// Add 新增任务
func (casing tasksCase) Add(ctx ctx.Context, param TasksAddRequest) error {
	// 如果isRun = 1 则为保存并运行
	var status uint8 = 0
	switch param.IsRun {
	case 1:
		status = transfer_customer.TCT_STATUS_RUNNING
	default:
		status = transfer_customer.TCT_STATUS_NOT
	}

	if count, countErr := orm.Service.TransferCustomer.TasksCheck(param.EnterpriseID, param.Name); countErr != nil {
		return countErr
	} else {
		if count > 0 {
			return errors.New("任务已存在，请勿重复添加")
		}
	}

	return orm.Service.TransferCustomer.TasksCreate(param.EnterpriseID, transfer_customer.TransferCustomerTasksModel{
		Name: param.Name,
		Type: func(transferType uint8) uint8 {
			if transferType == 0 {
				return uint8(transfer_customer.TCT_Type_Random)
			}
			return transferType
		}(param.TransferType),
		Status:     status,
		StartTime:  time.Now().Unix(),
		TransferID: fmt.Sprintf("jc%s", time.Now().Format(helper.YYYYMMDDHHMMSS)),
	})
}

// Get 根据任务的ID，获取任务的详情
func (casing tasksCase) Get(ctx ctx.Context, param TasksGetRequest) (data transfer_customer.TransferCustomerTasksModel, err error) {
	data, err = orm.Service.TransferCustomer.TasksGet(param.EnterpriseID, transfer_customer.TransferCustomerTasksModel{
		BranchModel: model.BranchModel{
			BaseModel: model.BaseModel{
				ID: param.ID,
			},
		},
	})

	return
}

// GetList 根据条件查询任务的详情
func (casing tasksCase) GetList(ctx ctx.Context, param TasksGetListRequest) (data TasksGetListResponse, err error) {
	data.Request = param

	var wg sync.WaitGroup
	wg.Add(1)
	go func(err *error) {
		defer wg.Done()
		if count, resErr := orm.Service.TransferCustomer.TasksGetListCount(param.EnterpriseID, param.StartTime, param.EndTime, param.Status); resErr != nil {
			*err = resErr
			return
		} else {
			data.Total = count
		}
	}(&err)

	wg.Add(1)
	go func(err *error) {
		defer wg.Done()
		if dataList, resErr := orm.Service.TransferCustomer.TasksGetList(param.EnterpriseID, param.StartTime, param.EndTime, param.Status, param.PaginateParam); resErr != nil {
			*err = resErr
			return
		} else {
			data.Data = append([]transfer_customer.TransferCustomerTasksModel{}, dataList...)
		}
	}(&err)
	wg.Wait()

	return
}
