package excitationModels

import "micro-scrm/internal/app/models"

type ExcitationIntegralRuleDescriptionModel struct {
	models.MicroBaseModel
	ID           int64  `json:"id" gorm:"primaryKey;column:id;type:int(11);autoIncrement"`
	EnterpriseID int64  `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20);default:0"`
	Content      string `json:"content" gorm:"column:content;type:text;not null"`
}

func (ExcitationIntegralRuleDescriptionModel) TableName() string {
	return "excitation_integral_rule_description"
}
