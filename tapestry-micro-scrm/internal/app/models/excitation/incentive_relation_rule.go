package excitationModels

import "micro-scrm/internal/app/models"

const (
	// RuleCondition1 指标条件: 1top榜  2高于 3达到 4低于 5top百分比'
	RuleCondition1 = iota + 1
	RuleCondition2
	RuleCondition3
	RuleCondition4
	RuleCondition5
)

type ExcitationIncentiveRelationRuleModel struct {
	models.MicroBaseModel

	IncentiveRelationID int64  `json:"incentive_relation_id" gorm:"column:incentive_relation_id;type:bigint(20) unsigned;not null;default:0;comment:incentive_relation表Id"`
	IncentiveConfigType int    `json:"incentive_config_type" gorm:"column:incentive_config_type;type:tinyint(4) unsigned;not null;default:0;comment:incentive_config表type"`
	RuleCondition       int    `json:"rule_condition" gorm:"column:rule_condition;type:tinyint(4);not null;default:0;comment:指标条件: 1top榜  1高于 2达到 3低于"`
	RuleValue           int    `json:"rule_value" gorm:"column:rule_value;type:int(11);not null;default:0;comment:指标条件值(整数)"`
	RuleGive            string `json:"rule_give" gorm:"column:rule_give;type:varchar(50);not null;comment:能力指标获得值"`
	StoreType           int    `json:"store_type" gorm:"column:store_type;type:tinyint(4) unsigned;not null;default:0;comment:0 通用规则 1奥莱 2零售"`
}

// TableName returns the table name of the ExcitationIncentiveRelationRule model
func (ExcitationIncentiveRelationRuleModel) TableName() string {
	return "excitation_incentive_relation_rule"
}
