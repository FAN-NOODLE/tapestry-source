package microScrmModel

import "micro-scrm/internal/app/models"

type PlanTaskConditionsModel struct {
	models.BaseModel
	PlanID      int64  `json:"plan_id" gorm:"column:plan_id;type:bigint(20);not null;default:0;comment:条件对应的子任务对应的计划的ID"`
	TaskID      int64  `json:"task_id" gorm:"column:task_id;type:bigint(20);not null;default:0;comment:条件对应的子任务的ID"`
	Conditions  string `json:"conditions" gorm:"column:conditions;type:text;not null;comment:条件内容，json结构;包含：多个维度的结构"`
	Content     string `json:"content" gorm:"column:content;type:text;not null;comment:条件内容，json结构；包含可以推送给顾客的所有信息：文本/图片/视频/小程序/链接"`
	ContentType int32  `json:"content_type" gorm:"column:content_type;type:tinyint(2);not null;default:0;comment:消息类型；1文本；2图片；3语音；4视频；5位置；6链接,7 小程序"`
}

// TableName returns the table name of the PlanTaskConditions model
func (PlanTaskConditionsModel) TableName() string {
	return "plan_task_conditions"
}
