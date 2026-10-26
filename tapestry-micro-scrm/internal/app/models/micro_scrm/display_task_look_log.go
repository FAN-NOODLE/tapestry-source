package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type DisplayTaskLookLogModel struct {
	models.BaseModel
	AllotDataID int64 `json:"allot_data_id" gorm:"column:allot_data_id;type:bigint(20);not null;uniqueIndex:unique_allot_member_index;default:0"`
	MemberID    int64 `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;uniqueIndex:unique_allot_member_index;default:0"`
	AllotID     int64 `json:"allot_id" gorm:"column:allot_id;type:bigint(20);not null;default:0;comment:分配表id"`
	CommentID   int64 `json:"comment_id" gorm:"column:comment_id;type:bigint(20);not null;uniqueIndex:unique_allot_member_index;default:0;comment:评论表id"`
}

// TableName returns the table name of the DisplayTaskLookLog model
func (DisplayTaskLookLogModel) TableName() string {
	return "display_task_look_log"
}
