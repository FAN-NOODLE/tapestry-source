package templateModels

import "micro-scrm/internal/app/models"

type TemplatePlanRelationModel struct {
	models.MicroBaseModel
	TemplateID int   `json:"template_id" gorm:"column:template_id;type:int(10);not null;default:0;comment:template表id"`
	PlanID     int64 `json:"plan_id" gorm:"column:plan_id;type:bigint(20);not null;default:0;comment:计划任务的plan_id"`
	PlanStatus int32 `json:"plan_status" gorm:"column:plan_status;type:tinyint(1);not null;default:0;comment:任务的状态，0为使用中，1为已过期"`
}

func (TemplatePlanRelationModel) TableName() string {
	return "template_plan_relation"
}
