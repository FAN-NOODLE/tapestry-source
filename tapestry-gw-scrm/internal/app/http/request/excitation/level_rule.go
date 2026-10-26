package excitationReq

type LevelRuleListReq struct {
	EnterpriseId int64 `json:"enterprise_id"`
}

type LevelRuleSaveArrReq struct {
	EnterpriseId int64              `json:"enterprise_id"`
	LevelRules   []LevelRuleInfoReq `json:"level_rules" binding:"required"`
}

type LevelRuleInfoReq struct {
	ID                 int64                `json:"id"`
	LevelName          string               `json:"level_name" gorm:"comment:等级名称"`
	LevelType          int                  `json:"level_type" gorm:"comment:等级类型"`
	LevelIcon          string               `json:"level_icon" gorm:"comment:等级icon"`
	MinValue           int                  `json:"min_value" gorm:"comment:最小分值(整数)"`
	MaxValue           int                  `json:"max_value" gorm:"comment:最大分值(整数)"`
	IsDelete           int                  `json:"is_delete" gorm:"是否删除"`
	RuleInviteID       int64                `json:"rule_invite_id"`
	RuleTaskID         int64                `json:"rule_task_id"`
	RuleInteractID     int64                `json:"rule_interact_id"`
	RuleTrainID        int64                `json:"rule_train_id"`
	RuleOnlineID       int64                `json:"rule_online_id"`
	LevelRuleAttachArr []LevelRuleAttachReq `json:"level_rule_attach_arr"`
}

type LevelRuleAttachReq struct {
	IncentiveRelationRuleId int64  `json:"incentive_relation_rule_id"`
	IncentiveConfigType     int    `json:"incentive_config_type" gorm:"comment:incentive_config表type"`
	RuleCondition           int    `json:"rule_condition" gorm:"comment:指标条件: 1top榜  2高于 3达到 4低于"`
	RuleGive                string `json:"rule_give" gorm:"comment:能力指标获得值"`
}

type LevelRuleDetailReq struct {
	EnterpriseId int64 `json:"enterprise_id"`
}

type LevelRuleDetailListReq struct {
	LevelRuleDetailReq
}
