package excitationModels

import (
	"micro-scrm/internal/app/models"
)

type ExcitationIntegralRuleModel struct {
	models.MicroBaseModel

	ID           int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20);autoIncrement"`
	EnterpriseID int64  `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20);default:0"`
	Type         int64  `json:"type" gorm:"column:type;type:bigint(20);not null"`
	RuleName     string `json:"rule_name" gorm:"column:rule_name;type:varchar(255);comment:规则名称"`
	UpdatedAt    int64  `gorm:"column:updated_at;type:int(10);not null;default:0" json:"updated_at"`

	//关联维度配置
	IncentiveConfig ExcitationIncentiveConfigModel `gorm:"foreignKey:Type;references:Type"`
	//关联能力图谱设置表
	IncentiveRelation ExcitationIncentiveRelationModel `gorm:"foreignKey:Type;references:IncentiveConfigType"`
	//关联规则信息
	Details []ExcitationIntegralRuleDetailModel `gorm:"foreignKey:RuleID;references:ID"`
}

func (ExcitationIntegralRuleModel) TableName() string {
	return "excitation_integral_rule"
}
