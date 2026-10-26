package microScrmModel

import "micro-scrm/internal/app/models"

type PlanTaskTemplateBackupModel struct {
	models.BaseModel
	TemplateID          int64  `json:"template_id" gorm:"column:template_id;type:bigint(20);not null;default:0;comment:主模版id"`
	TemplateContent     string `json:"template_content" gorm:"column:template_content;type:text;comment:任务模板的内容"`
	TemplateContentType int32  `json:"template_content_type" gorm:"column:template_content_type;type:tinyint(2);not null;default:0;comment:消息类型；1文本；2图片；3语音；4视频；5位置；6链接,7 小程序"`
	StartTime           int64  `json:"start_time" gorm:"column:start_time;type:int(10);not null;default:0;comment:生效开始时间"`
	EndTime             int64  `json:"end_time" gorm:"column:end_time;type:int(10);not null;default:0;comment:生效结束时间 0 为不结束"`
	GroupSendType       int    `json:"group_send_type" gorm:"column:group_send_type;type:tinyint(10);not null;default:0;comment:子任务群发类型，1企业群发助手，2，灵活沟通"`
	ALLowSelect         int    `json:"allow_select" gorm:"column:allow_select;type:tinyint(1);not null;default:0;comment:是否允许重新选择发送客户，0否，1是"`
	FinishRate          int    `json:"finish_rate" gorm:"column:finish_rate;type:tinyint(1);not null;default:100;comment:任务完成率"`
}

// TableName returns the table name of the PlanTaskTemplate model
func (PlanTaskTemplateBackupModel) TableName() string {
	return "plan_task_template_backup"
}
