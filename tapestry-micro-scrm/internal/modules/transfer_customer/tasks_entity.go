package transfer_customer

import (
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/transfer_customer"
	"micro-scrm/internal/orm/service"
)

// TasksGetRequest 获取单个任务信息请求
type TasksGetRequest struct {
	EnterpriseID int64 `json:"enterprise_id" binding:"required" desc:"企业ID"`
	ID           int64 `json:"id" binding:"required" desc:"任务ID"`
}

// TasksGetListRequest 获取单任务信息集合请求
type TasksGetListRequest struct {
	EnterpriseID int64 `json:"enterprise_id" binding:"required" desc:"企业ID"`
	StartTime    int64 `json:"start_time" desc:"任务开始的时间戳"`
	EndTime      int64 `json:"end_time" desc:"任务结束的时间戳"`
	Status       uint8 `json:"status" desc:"任务状态"`
	service.PaginateParam
}

// TasksGetListResponse 获取单任务信息集合
type TasksGetListResponse struct {
	Request TasksGetListRequest                            `json:"request"`
	Total   int64                                          `json:"total" desc:"总数据条数"`
	Data    []transfer_customer.TransferCustomerTasksModel `json:"data"`
}

// TasksAddRequest 新增任务请求
type TasksAddRequest struct {
	EnterpriseID int64  `json:"enterprise_id" binding:"required" desc:"企业ID"`
	Name         string `json:"name" binding:"required,max=40" desc:"任务名称"`
	TransferType uint8  `json:"transfer_type" binding:"required" desc:"继承类型"`
	IsRun        uint8  `json:"is_run" desc:"是否立即启动"`
}

// TasksUpdateRequest 更新任务请求
type TasksUpdateRequest struct {
	EnterpriseID int64  `json:"enterprise_id" binding:"required" desc:"企业ID"`
	ID           int64  `json:"id" binding:"required" desc:"任务ID"`
	Name         string `json:"name" desc:"任务名称"`
	TransferType uint8  `json:"transfer_type" desc:"继承类型"`
	Status       uint8  `json:"status" desc:"任务状态"`
}

// TasksDeleteRequest 删除任务请求
type TasksDeleteRequest struct {
	EnterpriseID int64 `json:"enterprise_id" binding:"required" desc:"企业ID"`
	ID           int64 `json:"id" binding:"required" desc:"任务ID"`
}

// TasksReportRequest 任务报表请求
type TasksReportRequest struct {
	EnterpriseID int64 `json:"enterprise_id" binding:"required" desc:"企业ID"`
	ID           int64 `json:"id" binding:"required" desc:"任务ID"`
	IsOnlyList   bool  `json:"is_only_list" desc:"是否只请求客户列表"`
	service.PaginateParam
}

// TasksReportResponse 任务报表响应
type TasksReportResponse struct {
	Request     TasksReportRequest                           `json:"request"`
	TaskInfo    transfer_customer.TransferCustomerTasksModel `json:"task_info"`
	ReportTitle TasksReportTitle                             `json:"report_title"`
	ReportBody  TasksReportBody                              `json:"report_body"`
}

type TasksReportTitle struct {
	CustomerTotalNum   int `json:"customer_total_num" desc:"客户总人数"`
	CustomerSuccessNum int `json:"customer_success_num" desc:"成功继承客户数"`
	CustomerFailNum    int `json:"customer_fail_num" desc:"继承失败客户数"`
	FinishBANum        int `json:"finish_ba_num" desc:"完成继承BA数"`
	FinishRate         int `json:"finish_rate" desc:"继承完成率"`
}

type TasksReportBody struct {
	service.PaginateParam
	DataList []TasksReportListData `json:"data_list"`
}

type TasksReportListData struct {
	Name               string `json:"name" desc:"客户名称"`
	FriendNum          int    `json:"friend_num" desc:"好友数"`
	CurrentAdminBAName string `json:"current_admin_ba_name" desc:"当前管理BA"`
	BADepartmentName   string `json:"ba_department_name" desc:"所属部门"`
	TransferStatus     string `json:"transfer_status" desc:"继承状态"`
	TransferStatusMsg  string `json:"transfer_status_msg" desc:"继承状态文本"`
	FailMsg            string `json:"fail_msg" desc:"失败明文"`
}
