package excitationResp

type LevelRuleDetailResp struct {
	EnterpriseId int64  `json:"enterprise_id"`
	ID           int64  `json:"id"`
	LevelType    int    `json:"level_type" gorm:"comment:等级类型"`
	LevelName    string `json:"level_name" gorm:"comment:等级名称"`
	LevelIcon    string `json:"level_icon" gorm:"comment:等级icon"`
	MinValue     int    `json:"min_value" gorm:"comment:最小分值(整数)"`
	MaxValue     int    `json:"max_value" gorm:"comment:最大分值(整数)"`

	RuleConditionInvite   int `json:"rule_condition_invite" gorm:"comment:指标条件: 1top榜  2高于 3达到 4低于"`
	RuleConditionTask     int `json:"rule_condition_task"`
	RuleConditionInteract int `json:"rule_condition_interact"`
	RuleConditionTrain    int `json:"rule_condition_train"`
	RuleConditionOnline   int `json:"rule_condition_online"`
	LevelRuleDetailListIDsResp
}

type LevelRuleDetailInfoResp struct {
	LevelRuleDetailResp
	LevelRuleDetailListIDsResp
}

type LevelRuleDetailListResp struct {
	LevelRuleDetailResp
	LevelRuleDetailListIDsResp
	LevelRuleAttachArr []LevelRuleAttachResp `json:"level_rule_attach_arr"`
}

type LevelRuleDetailListIDsResp struct {
	RuleInviteID   int64 `json:"rule_invite_id" gorm:"comment:incentive_relation_rule表Id"`
	RuleTaskID     int64 `json:"rule_task_id" gorm:"comment:incentive_relation_rule表Id"`
	RuleInteractID int64 `json:"rule_interact_id" gorm:"comment:incentive_relation_rule表Id"`
	RuleTrainID    int64 `json:"rule_train_id" gorm:"comment:incentive_relation_rule表Id"`
	RuleOnlineID   int64 `json:"rule_online_id" gorm:"comment:incentive_relation_rule表Id"`
}

type LevelRuleAttachResp struct {
	IncentiveRelationRuleId int64  `json:"incentive_relation_rule_id"`
	IncentiveConfigType     int    `json:"incentive_config_type" gorm:"comment:incentive_config表type"`
	RuleCondition           int    `json:"rule_condition" gorm:"comment:指标条件: 1top榜  2高于 3达到 4低于"`
	RuleGive                string `json:"rule_give" gorm:"comment:能力指标获得值"`
}
