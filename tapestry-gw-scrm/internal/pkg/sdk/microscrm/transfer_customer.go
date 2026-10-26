package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
	"sync"
)

var trans *TransferCustomerApi

type TransferCustomerApi struct {
	Name string
}

var once sync.Once

func NewTransferCustomerApi() *TransferCustomerApi {
	once.Do(func() {
		trans = &TransferCustomerApi{
			Name: "客户继承",
		}
	})

	return trans
}

const (
	URL_TransferCustomer_TasksGet    = "%s/micro-scrm/v1/transfer_customer/task"
	URL_TransferCustomer_TasksList   = "%s/micro-scrm/v1/transfer_customer/tasks_list"
	URL_TransferCustomer_TasksAdd    = "%s/micro-scrm/v1/transfer_customer/tasks_add"
	URL_TransferCustomer_TasksUpdate = "%s/micro-scrm/v1/transfer_customer/tasks_update"
	URL_TransferCustomer_TasksDelete = "%s/micro-scrm/v1/transfer_customer/tasks_delete"
	URL_TransferCustomer_TasksReport = "%s/micro-scrm/v1/transfer_customer/tasks_report"
)

// TasksGet 获取客户继承任务
func (api TransferCustomerApi) TasksGet(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(URL_TransferCustomer_TasksGet), request)
}

// TasksList 获取客户继承任务列表
func (api TransferCustomerApi) TasksList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(URL_TransferCustomer_TasksList), request)
}

// TasksAdd 添加在职继承任务
func (api TransferCustomerApi) TasksAdd(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(URL_TransferCustomer_TasksAdd), request)
}

// TasksUpdate 更新在职继承任务
func (api TransferCustomerApi) TasksUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(URL_TransferCustomer_TasksUpdate), request)
}

// TasksDelete 删除在职继承任务
func (api TransferCustomerApi) TasksDelete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(URL_TransferCustomer_TasksDelete), request)
}

// TasksReport 获取任务报表
func (api TransferCustomerApi) TasksReport(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(URL_TransferCustomer_TasksReport), request)
}
