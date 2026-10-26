package transfer_customer

import "micro-scrm/internal/orm/model"

const (
	TCD_STATUS_NO = iota
	TCD_STATUS_Finish
	TCD_STATUS_Wait
	TCD_STATUS_Refuse
	TCD_STATUS_Ceiling
	TCD_STATUS_NoRecord
	TCD_STATUS_FAIL
)

const (
	TCD_TYPE_Only = iota + 1
	TCD_TYPE_Random
)

type TransferCustomerDataModel struct {
	model.BranchModel

	HandoverUserid         string `gorm:"column:handover_userid" json:"handover_userid"` // 原跟进成员的userid
	HandoverName           string `gorm:"handover_name" json:"handover_name"`
	TakeoverUserid         string `gorm:"column:takeover_userid" json:"takeover_userid"` // 接替成员的userid
	TakeoverName           string `gorm:"takeover_name" json:"takeover_name"`
	ExternalUserid         string `gorm:"column:external_userid" json:"external_userid"`                             // 客户的external_userid
	ExternalUnionID        string `gorm:"column:external_unionid" json:"external_unionid"`                           // 客户的unionid
	TransferSuccessMsgId   int64  `gorm:"column:transfer_success_msg_id" json:"transfer_success_msg_id"`             // 在职继承欢迎语编号,为空的
	TransferCustomerStatus int    `gorm:"column:transfer_customer_status;default:0" json:"transfer_customer_status"` // 接替状态，0 未开始、 1-接替完毕 2-等待接替 3-客户拒绝 4-接替成员客户达到上限 5-无接替记录 6 已执行在职继承任务
	TransferCustomerType   int    `gorm:"column:transfer_customer_type;NOT NULL" json:"transfer_customer_type"`      // 在职继承类型 1完全继承、2 被选择但不继承
	TaskId                 int64  `gorm:"column:task_id" json:"task_id"`                                             // 执行的任务编号
	// 删除时间
}

func (TransferCustomerDataModel) TableName() string {
	return "transfer_customer_data"
}

type TransferCustomerDataModel_TasksReportTitle struct {
	CustomerTotalNum   int `json:"customer_total_num" desc:"客户总人数"`
	CustomerSuccessNum int `json:"customer_success_num" desc:"成功继承客户数"`
	CustomerFailNum    int `json:"customer_fail_num" desc:"继承失败客户数"`
	FinishBANum        int `json:"finish_ba_num" desc:"完成继承BA数"`
	FinishRate         int `json:"finish_rate" desc:"继承完成率"`
}

type TransferCustomerDataModel_TasksReportListData struct {
	Name               string `json:"name" desc:"客户名称"`
	FriendNum          int    `json:"friend_num" desc:"好友数"`
	CurrentAdminBAName string `json:"current_admin_ba_name" desc:"当前管理BA"`
	BADepartmentName   string `json:"ba_department_name" desc:"所属部门"`
	TransferStatus     string `json:"transfer_status" desc:"继承状态"`
	TransferStatusMsg  string `json:"transfer_status_msg" desc:"继承状态文本"`
	FailMsg            string `json:"fail_msg" desc:"失败明文"`
}

// TransferCustomerDataModel_GetTransferDataList 获取继承数据的实体
type TransferCustomerDataModel_GetTransferDataList struct {
	HandoverUserID     string `json:"handover_user_id" desc:"原跟进员工user_id"`
	TakeoverUserID     string `json:"takeover_user_id" desc:"现接替员工user_id"`
	WorkExternalUserID string `json:"work_external_user_id" desc:"被接替的客户的企微ID"`
	UnionID            string `json:"union_id" desc:"被接替的客户的union_id"`
}
