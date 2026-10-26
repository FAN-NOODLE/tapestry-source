package excitationModels

import (
	"micro-scrm/internal/app/models"
)

type ExcitationMemberDataModel struct {
	models.MicroBaseModel
	MemberID          int64   `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:member_id_index;default:0;comment:员工id"`
	PointValue        string  `json:"point_value" gorm:"column:point_value;type:varchar(64);not null;comment:积分(整数)"`
	PointUpdatedAt    int     `json:"point_updated_at" gorm:"column:point_updated_at;type:int(11) unsigned;not null;default:0;comment:积分更新时间"`
	Coins             float64 `json:"coins" gorm:"column:coins;type:float;not null;comment:金币"`
	CoinsUpdatedAt    int     `json:"coins_updated_at" gorm:"column:coins_updated_at;type:int(11) unsigned;not null;default:0;comment:金币更新时间"`
	SpendCoins        float64 `json:"spend_coins" gorm:"column:spend_coins;type:float;not null;comment:总花费金币"`
	TotalCoins        float64 `json:"total_coins" gorm:"column:total_coins;type:float;not null;comment:总获得金币"`
	Level             string  `json:"level" gorm:"column:level;type:varchar(50);not null;comment:等级名称"`
	LevelType         int     `json:"level_type" gorm:"column:level_type;type:tinyint(4);not null;default:0;comment:等级类型"`
	LevelID           int64   `json:"level_id" gorm:"column:level_id;type:bigint(20) unsigned;not null;default:0;comment:level_rule表Id"`
	LevelUpdatedAt    int     `json:"level_updated_at" gorm:"column:level_updated_at;type:int(11) unsigned;not null;default:0;comment:等级更新时间"`
	MedalNum          int     `json:"medal_num" gorm:"column:medal_num;type:int(11);not null;default:0;comment:奖章数量(整数)"`
	MedalUpdatedAt    int     `json:"medal_updated_at" gorm:"column:medal_updated_at;type:int(11) unsigned;not null;default:0;comment:奖章更新时间"`
	InviteValue       string  `json:"invite_value" gorm:"column:invite_value;type:varchar(50);not null;comment:邀约能力指标"`
	InviteRuleValue   int     `json:"invite_rule_value" gorm:"column:invite_rule_value;type:int(11);not null;default:0;comment:邀约目标达成率值(整数)"`
	InviteRuleStr     string  `json:"invite_rule_str" gorm:"column:invite_rule_str;type:varchar(100);not null;comment:邀约能力指标说明"`
	InviteUpdatedAt   int     `json:"invite_updated_at" gorm:"column:invite_updated_at;type:int(11) unsigned;not null;default:0;comment:邀约能力指标更新时间"`
	TaskValue         string  `json:"task_value" gorm:"column:task_value;type:varchar(50);not null;comment:任务能力指标"`
	TaskRuleValue     int     `json:"task_rule_value" gorm:"column:task_rule_value;type:int(11);not null;default:0;comment:任务目标达成率值(整数)"`
	TaskRuleStr       string  `json:"task_rule_str" gorm:"column:task_rule_str;type:varchar(100);not null;comment:任务能力指标说明"`
	TaskUpdatedAt     int     `json:"task_updated_at" gorm:"column:task_updated_at;type:int(11) unsigned;not null;default:0;comment:任务能力指标更新时间"`
	InteractValue     string  `json:"interact_value" gorm:"column:interact_value;type:varchar(50);not null;comment:互动能力指标"`
	InteractRuleValue int     `json:"interact_rule_value" gorm:"column:interact_rule_value;type:int(11);not null;default:0;comment:互动目标达成率值(整数)"`
	InteractRuleStr   string  `json:"interact_rule_str" gorm:"column:interact_rule_str;type:varchar(100);not null;comment:互动能力指标说明"`
	InteractUpdatedAt int     `json:"interact_updated_at" gorm:"column:interact_updated_at;type:int(11) unsigned;not null;default:0;comment:互动能力指标更新时间"`
	TrainValue        string  `json:"train_value" gorm:"column:train_value;type:varchar(50);not null;comment:培训能力指标"`
	TrainRuleValue    int     `json:"train_rule_value" gorm:"column:train_rule_value;type:int(11);not null;default:0;comment:培训目标达成率值(整数)"`
	TrainRuleStr      string  `json:"train_rule_str" gorm:"column:train_rule_str;type:varchar(100);not null;comment:培训能力指标说明"`
	TrainUpdatedAt    int     `json:"train_updated_at" gorm:"column:train_updated_at;type:int(11) unsigned;not null;default:0;comment:培训能力指标更新时间"`
	OnlineValue       string  `json:"online_value" gorm:"column:online_value;type:varchar(50);not null;comment:线上能力指标"`
	OnlineRuleValue   int     `json:"online_rule_value" gorm:"column:online_rule_value;type:int(11);not null;default:0;comment:线上目标达成率值(整数)"`
	OnlineRuleStr     string  `json:"online_rule_str" gorm:"column:online_rule_str;type:varchar(100);not null;comment:线上能力指标说明"`
	OnlineUpdatedAt   int     `json:"online_updated_at" gorm:"column:online_updated_at;type:int(11) unsigned;not null;default:0;comment:线上能力指标更新时间"`
}

// TableName returns the table name of the ExcitationMemberData model
func (ExcitationMemberDataModel) TableName() string {
	return "excitation_member_data"
}
