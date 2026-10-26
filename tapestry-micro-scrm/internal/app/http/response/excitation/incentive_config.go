package excitationResp

import (
	excitationModels "micro-scrm/internal/app/models/excitation"
)

type IncentiveConfigListResp struct {
	EnterpriseId          int64                 `json:"enterprise_id"`
	Id                    int64                 `json:"id"`
	RuleName              string                `json:"rule_name"`
	Status                int32                 `json:"status"`
	Type                  int32                 `json:"type"`
	IncentiveRelationResp IncentiveRelationResp `json:"incentive_relation" `
	IntegralRuleResp      IntegralRuleResp      `json:"integral_rule" `
}

type IncentiveRelationResp struct {
	Id                  int64 `json:"id"`
	IncentiveConfigType int   `json:"incentive_config_type"`
	Status              int   `json:"status" gorm:"comment:开启状态: 1开启 2关闭"`
	AbilityIndex        int   `json:"ability_index" `
	TimeRule            int   `json:"time_rule"`
}

type IntegralRuleResp struct {
	ID           int64  `json:"id"`
	EnterpriseID int64  `json:"enterprise_id"`
	Type         int64  `json:"type"`
	RuleName     string `json:"rule_name"`
}

type RelationDetailResp struct {
	RelationDetail      excitationModels.ExcitationIncentiveRelationModel `json:"relation_detail"`
	RelationRuleDetails []RelationRuleDetailResp                          `json:"relation_rule_details"`
}

type RelationRuleDetailResp struct {
	Id                  int64  `json:"id"`
	IncentiveRelationId int64  `json:"incentive_relation_id" gorm:"comment:incentive_relation_id"`
	IncentiveConfigType int    `json:"incentive_config_type" gorm:"comment:incentive_config表type"`
	RuleCondition       int    `json:"rule_condition" gorm:"comment:指标条件: 1top榜  2高于 3达到 4低于"`
	RuleValue           int    `json:"rule_value"`
	RuleGive            string `json:"rule_give" gorm:"comment:能力指标获得值"`
	StoreType           int    `json:"store_type"`
}

type RelationRuleDataResp struct {
	Id                  int64  `json:"id"`
	IncentiveRelationId int64  `json:"incentive_relation_id" gorm:"comment:incentive_relation_id"`
	IncentiveConfigType int    `json:"incentive_config_type" gorm:"comment:incentive_config表type"`
	RuleCondition       int    `json:"rule_condition" gorm:"comment:指标条件: 1top榜  2高于 3达到 4低于"`
	RuleValue           int    `json:"rule_value"`
	RuleGive            string `json:"rule_give" gorm:"comment:能力指标获得值"`
	SortKey             int    `json:"sort_key"`
	StoreType           int    `json:"store_type"`
}
