package excitationReq

import "micro-scrm/internal/app/http/request"

type IncentiveConfigListReq struct {
	request.EnterpriseBaseReq
}

type IncentiveRelationSaveReq struct {
	request.EnterpriseBaseReq
	ID                  int64                          `json:"id"`
	IncentiveConfigType int                            `json:"incentive_config_type"`
	Status              int                            `json:"status" gorm:"comment:开启状态: 1开启 2关闭"`
	AbilityIndex        int                            `json:"ability_index" gorm:"comment:能力指标条件: 1邀请目标达成率(每月) 2任务完成率(每月) 3总客户互动积分排行榜(每周) 4总培训积分排行榜(每周) 5目标达成率(每月)"`
	TimeRule            int                            `json:"time_rule" gorm:"comment:时间规则: 1每次(人) 2每天 3每周 4每月 5每季度 6每年"`
	IsDiffStore         int                            `json:"is_diff_store"`
	RelationRule        []IncentiveRelationRuleSaveReq `json:"relation_rule"`
}

type IncentiveRelationRuleSaveReq struct {
	request.EnterpriseBaseReq
	IncentiveRelationID int64  `json:"incentive_relation_id" gorm:"comment:incentive_relation表Id"`
	IncentiveConfigType int    `json:"incentive_config_type" gorm:"comment:incentive_config表type"`
	RuleCondition       int    `json:"rule_condition" gorm:"comment:指标条件: 1top榜  2高于 3达到 4低于"`
	RuleValue           int    `json:"rule_value" gorm:"comment:指标条件值(整数)"`
	RuleGive            string `json:"rule_give" gorm:"comment:能力指标获得值"`
	StoreType           int    `json:"store_type"`
}

type IncentiveRelationDetailReq struct {
	request.EnterpriseBaseReq
	ID int64 `json:"id"`
}

type IncentiveRelationDetailListReq struct {
	request.EnterpriseBaseReq
}
