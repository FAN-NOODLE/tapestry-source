package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

const (
	VMSubmitDraft = iota
	VMSubmitCheck
)

type DisplayTaskAllotSubmitModel struct {
	models.BaseModel
	DisplayTaskID int64  `json:"display_task_id" gorm:"column:display_task_id;type:bigint(20);not null;default:0;comment:任务表id"`
	AllotID       int64  `json:"allot_id" gorm:"column:allot_id;type:bigint(20);not null;default:0;comment:分配表id"`
	SubmitType    int32  `json:"submit_type" gorm:"column:submit_type;type:tinyint(1);not null;default:0;comment:提交状态，0为草稿，1为提交审批"`
	MemberID      int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;comment:导购id"`
	MemberName    string `json:"member_name" gorm:"column:member_name;type:varchar(255);not null;comment:导购昵称"`
	MemberUserID  string `json:"member_user_id" gorm:"column:member_user_id;type:varchar(255);not null;comment:导购工号"`
}

// TableName returns the table name of the DisplayTaskAllotSubmit model
func (DisplayTaskAllotSubmitModel) TableName() string {
	return "display_task_allot_submit"
}
