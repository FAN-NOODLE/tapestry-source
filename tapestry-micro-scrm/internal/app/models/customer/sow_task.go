package customerModels

import "micro-scrm/internal/app/models"

type SowTaskModel struct {
	WorkExternalUserID int64  `json:"work_external_user_id" gorm:"column:work_external_user_id;type:bigint(20);not null;index:work_external_user_id_index;default:0;comment:客户id"`
	MemberID           int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:member_id_index;default:0;comment:员工id"`
	MemberUserID       string `json:"member_user_id" gorm:"column:member_user_id;type:varchar(100);not null;index:member_user_id_index;default:0;comment:员工企微Id"`
	Content            string `json:"content" gorm:"column:content;type:varchar(255);not null;comment:播种内容"`
	Status             int32  `json:"status" gorm:"column:status;type:tinyint(1);not null;default:1;comment:完成状态 1 未完成 2 已完成"`
	RemindTime         int64  `json:"remind_time" gorm:"column:remind_time;type:int(10);not null;default:0;comment:提醒时间"`
	IsPush             int32  `json:"is_push" gorm:"column:is_push;type:tinyint(1);not null;default:0;comment:是否已推送 1 已推送 0 未推送"`
	TemplateId         int64  `json:"template_id" gorm:"column:template_id;type:bigint(20);not null;index:template_id_index;default:0;comment:模版id，0为自定义"`
	models.MicroBaseModel
}

// TableName 播种任务
func (SowTaskModel) TableName() string {
	return "sow_task"
}

const (
	SowTaskStatusUnfinished = 1 // 1 未完成
	SowTaskStatusFinished   = 2 // 2 已完成
)

const (
	SowTaskIsPushTrue  = 1 // 1 已推送
	SowTaskIsPushFalse = 0 // 0 未推送
)
