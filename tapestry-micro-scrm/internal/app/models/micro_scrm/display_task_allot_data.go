package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

const (
	VmImageSuccess = iota
	VmImageFail
)

type DisplayTaskAllotDataModel struct {
	models.BaseModel
	DisplayTaskID int64                         `json:"display_task_id" gorm:"column:display_task_id;type:bigint(20);not null;default:0;comment:任务表id"`
	AllotID       int64                         `json:"allot_id" gorm:"column:allot_id;type:bigint(20);not null;index:allot_index;default:0;comment:分配表id"`
	Identifier    string                        `json:"identifier" gorm:"column:identifier;type:varchar(255);not null;comment:模板字段标识"`
	ImageURL      string                        `json:"image_url" gorm:"column:image_url;type:varchar(255);not null;comment:上传的图片地址"`
	Status        int32                         `json:"status" gorm:"column:status;type:tinyint(1);not null;default:0;comment:图片的审批状态，0为默认通过，1为未通过"`
	DraftStatus   int32                         `json:"draft_status" gorm:"column:draft_status;type:tinyint(1);not null;default:0;comment:图片的草稿审核状态，0为默认通过，1为未通过"`
	IsDraft       int32                         `json:"is_draft" gorm:"column:is_draft;type:tinyint(1);not null;default:0;comment:是否存草稿，0否，1是"`
	Remark        []DisplayTaskAllotRemarkModel `json:"remark" gorm:"foreignKey:AllotDataID;references:ID"`
}

// TableName returns the table name of the DisplayTaskAllotData model
func (DisplayTaskAllotDataModel) TableName() string {
	return "display_task_allot_data"
}
