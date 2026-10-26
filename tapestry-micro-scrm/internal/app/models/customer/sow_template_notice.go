package customerModels

import "micro-scrm/internal/app/models"

type SowTemplateNoticeModel struct {
	TemplateID int64  `json:"template_id" gorm:"column:template_id;type:bigint(20);not null;default:0;comment:模版id"`
	DelayType  int    `json:"delay_type" gorm:"column:delay_type;type:tinyint(1);not null;default:0;comment:延迟提醒类型：1天 2周 3月"`
	DelayValue int    `json:"delay_value" gorm:"column:delay_value;type:int(10);not null;default:0;comment:延迟提醒的时间数值"`
	DelayTime  string `json:"delay_time" gorm:"column:delay_time;type:varchar(50);not null;default:'';comment:延迟提醒当天的提醒时间，时:分"`
	Content    string `json:"content" gorm:"column:content;type:varchar(255);not null;default:'';comment:提醒文案"`
	models.MicroBaseModel
}

// TableName 播种任务
func (SowTemplateNoticeModel) TableName() string {
	return "sow_template_notice"
}
