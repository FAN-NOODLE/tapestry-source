package customerModels

import "micro-scrm/internal/app/models"

const (
	InheritChatStatusError   = iota // 0 分配失败
	InheritChatStatusSuccess        // 1 分配成功
)

type CustomerChatInheritModel struct {
	models.MicroBaseModel
	HandoverUserid         string `json:"handover_userid" gorm:"column:handover_userid;type:varchar(255);not null;index:handover_userid_index;comment:离职成员的userid"`
	HandoverName           string `json:"handover_name" gorm:"column:handover_name;type:varchar(255);not null;comment:离职成员姓名"`
	HandoverDepartmentID   int64  `json:"handover_department_id" gorm:"column:handover_department_id;type:bigint(10);not null;default:0;comment:离职成员部门id"`
	HandoverDepartmentName string `json:"handover_department_name" gorm:"column:handover_department_name;type:varchar(255);not null;comment:离职成员部门名"`
	ChatID                 string `json:"chat_id" gorm:"column:chat_id;type:varchar(255);not null;index:chat_id_index;comment:企微群Id"`
	ChatName               string `json:"chat_name" gorm:"column:chat_name;type:varchar(255);not null;comment:企微群名称"`
	TakeoverUserid         string `json:"takeover_userid" gorm:"column:takeover_userid;type:varchar(255);not null;comment:接替人userid"`
	TakeoverName           string `json:"takeover_name" gorm:"column:takeover_name;type:varchar(255);not null;comment:接替人姓名"`
	TakeoverDepartmentName string `json:"takeover_department_name" gorm:"column:takeover_department_name;type:varchar(255);not null;comment:接替人部门名"`
	Status                 int32  `json:"status" gorm:"column:status;type:tinyint(1);not null;default:0;comment:状态: 0 分配失败 1 分配成功"`
}

// TableName returns the table name of the CustomerChatInherit model
func (CustomerChatInheritModel) TableName() string {
	return "customer_chat_inherit"
}
