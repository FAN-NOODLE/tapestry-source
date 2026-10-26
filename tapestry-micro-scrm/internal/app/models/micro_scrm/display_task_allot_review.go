package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

const (
	VMReviewFail = iota + 1 // VM审核的状态，1为不通过，2为通过
	VMReviewSuccess
)

const (
	VMSpecialDraft = iota + 1
	VMSpecialSave
)

type DisplayTaskAllotReviewModel struct {
	models.BaseModel
	DisplayTaskID int64  `json:"display_task_id" gorm:"column:display_task_id;type:bigint(20);not null;default:0;comment:任务表id"`
	AllotID       int64  `json:"allot_id" gorm:"column:allot_id;type:bigint(20);not null;default:0;comment:分配表id"`
	ReviewID      int64  `json:"review_id" gorm:"column:review_id;type:bigint(20);not null;default:0;comment:审批者id"`
	ReviewType    int    `json:"review_type" gorm:"column:review_type;type:tinyint(1);not null;default:0;comment:审核状态，1为不通过，2为通过"`
	ReviewUserID  string `json:"review_user_id" gorm:"column:review_user_id;type:varchar(255);not null;comment:审批者工号"`
	ReviewName    string `json:"review_name" gorm:"column:review_name;type:varchar(255);not null;comment:审批者名称"`
}

// TableName returns the table name of the DisplayTaskAllotReview model
func (DisplayTaskAllotReviewModel) TableName() string {
	return "display_task_allot_review"
}
