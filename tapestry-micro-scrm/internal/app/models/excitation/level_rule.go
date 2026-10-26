package excitationModels

import "micro-scrm/internal/app/models"

const (
	// LV1 指标条件: 等级'
	LV1 = iota + 1
	LV2
	LV3
	LV4
	LV5
)

type ExcitationLevelRuleModel struct {
	models.MicroBaseModel
	LevelType             int    `json:"level_type" gorm:"column:level_type;type:tinyint(4);not null;default:0;comment:等级类型"`
	LevelName             string `json:"level_name" gorm:"column:level_name;type:varchar(50);not null;comment:等级名称"`
	LevelIcon             string `json:"level_icon" gorm:"column:level_icon;type:varchar(255);not null;comment:等级icon"`
	MinValue              int    `json:"min_value" gorm:"column:min_value;type:int(11);not null;default:0;comment:最小分值(整数)"`
	MaxValue              int    `json:"max_value" gorm:"column:max_value;type:int(11);not null;default:0;comment:最大分值(整数)"`
	RuleInviteID          int64  `json:"rule_invite_id" gorm:"column:rule_invite_id;type:bigint(20);not null;default:0;comment:incentive_relation_rule表Id"`
	RuleTaskID            int64  `json:"rule_task_id" gorm:"column:rule_task_id;type:bigint(20);not null;default:0;comment:incentive_relation_rule表Id"`
	RuleInteractID        int64  `json:"rule_interact_id" gorm:"column:rule_interact_id;type:bigint(20);not null;default:0;comment:incentive_relation_rule表Id"`
	RuleTrainID           int64  `json:"rule_train_id" gorm:"column:rule_train_id;type:bigint(20);not null;default:0;comment:incentive_relation_rule表Id"`
	RuleOnlineID          int64  `json:"rule_online_id" gorm:"column:rule_online_id;type:bigint(20);not null;default:0;comment:incentive_relation_rule表Id"`
	RuleConditionInvite   int    `json:"rule_condition_invite" gorm:"column:rule_condition_invite;type:tinyint(4);not null;default:0;comment:指标条件: 1top榜  2高于 3达到 4低于"`
	RuleConditionTask     int    `json:"rule_condition_task" gorm:"column:rule_condition_task;type:tinyint(4);not null;default:0;comment:指标条件: 1top榜  2高于 3达到 4低于"`
	RuleConditionInteract int    `json:"rule_condition_interact" gorm:"column:rule_condition_interact;type:tinyint(4);not null;default:0;comment:指标条件: 1top榜  2高于 3达到 4低于"`
	RuleConditionTrain    int    `json:"rule_condition_train" gorm:"column:rule_condition_train;type:tinyint(4);not null;default:0;comment:指标条件: 1top榜  2高于 3达到 4低于"`
	RuleConditionOnline   int    `json:"rule_condition_online" gorm:"column:rule_condition_online;type:tinyint(4);not null;default:0;comment:指标条件: 1top榜  2高于 3达到 4低于"`
}

// TableName returns the table name of the WjIncentiveConfig model
func (ExcitationLevelRuleModel) TableName() string {
	return "excitation_level_rule"
}
