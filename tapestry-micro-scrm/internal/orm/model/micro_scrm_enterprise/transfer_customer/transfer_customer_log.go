package transfer_customer

import "micro-scrm/internal/orm/model"

type TransferCustomerLogModel struct {
	model.BranchModel

	HandoverUserid         string `gorm:"column:handover_userid" json:"handover_userid"`                             // 原跟进成员的userid
	TakeoverUserid         string `gorm:"column:takeover_userid" json:"takeover_userid"`                             // 接替成员的userid
	ExternalUserid         string `gorm:"column:external_userid" json:"external_userid"`                             // 客户的external_userid
	TransferSuccessMsgId   int64  `gorm:"column:transfer_success_msg_id" json:"transfer_success_msg_id"`             // 在职继承欢迎语编号,为空的
	TransferCustomerStatus int    `gorm:"column:transfer_customer_status;default:0" json:"transfer_customer_status"` // 接替状态， 1-接替完毕 2-等待接替 3-客户拒绝 4-接替成员客户达到上限 5-无接替记录
	TakeoverTime           int64  `gorm:"column:takeover_time" json:"takeover_time"`                                 // 接替客户的时间，如果是等待接替状态，则为未来的自动接替时间
	TaskId                 int64  `gorm:"column:task_id" json:"task_id"`                                             // 执行的任务编号
	Errcode                int    `gorm:"column:errcode" json:"errcode"`                                             // 请求企微接口返回状态码
	Errmsg                 string `gorm:"column:errmsg" json:"errmsg"`                                               // 请求企微接口返回状态码说明
	TransferCustomerDataId int64  `gorm:"column:transfer_customer_data_id" json:"transfer_customer_data_id"`         // 客户转移信息编号                                   // 删除时间                                              // 请求企微接口返回状态码说明

}

func (TransferCustomerLogModel) TableName() string {
	return "transfer_customer_log"
}
