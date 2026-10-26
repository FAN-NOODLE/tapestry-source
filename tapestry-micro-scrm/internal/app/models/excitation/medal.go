package excitationModels

import "micro-scrm/internal/app/models"

type ExcitationMedalModel struct {
	models.MicroBaseModel

	Name            string `json:"name" gorm:"column:name;type:varchar(100);not null;comment:奖章规则名称"`
	Source          int32  `json:"source" gorm:"column:source;type:tinyint(1);not null;default:1;comment:奖章来源 1手动上传,2规则统计"`
	RuleType        int32  `json:"rule_type" gorm:"column:rule_type;type:tinyint(1);not null;default:1;comment:规则有效期 1长期有效,2指定时间,3仅生效一次"`
	RuleStartTime   int    `json:"rule_start_time" gorm:"column:rule_start_time;type:int(10);not null;default:0;comment:规则有效期-开始"`
	RuleEndTime     int    `json:"rule_end_time" gorm:"column:rule_end_time;type:int(10);not null;default:0;comment:规则有效期-结束"`
	MedaType        int32  `json:"meda_type" gorm:"column:meda_type;type:tinyint(1);not null;default:1;comment:奖章类型 1普通奖章,2隐藏奖章"`
	MedaDescription string `json:"meda_description" gorm:"column:meda_description;type:varchar(255);not null;comment:奖章描述"`
	OtherRewards    int32  `json:"other_rewards" gorm:"column:other_rewards;type:tinyint(1);not null;default:1;comment:其它奖励 1无其它奖励,2仅奖励积分,3仅奖励实物奖品,4奖励积分+实物奖品"`
	Integral        int32  `json:"integral" gorm:"column:integral;type:smallint(5);not null;default:0;comment:积分"`
	ItemName        string `json:"item_name" gorm:"column:item_name;type:varchar(50);not null;comment:物品名称"`
	Status          int32  `json:"status" gorm:"column:status;type:tinyint(1);not null;default:2;comment:状态 1启用,2禁用"`
	RuleDimensionID int    `json:"rule_dimension_id" gorm:"column:rule_dimension_id;type:int(11);not null;default:0;comment:规则维度-规则表主键id"`
	RuleCategoryID  int    `json:"rule_category_id" gorm:"column:rule_category_id;type:int(11);not null;default:0;comment:规则类目-规则表主键id"`
	InviteTop       int32  `json:"invite_top" gorm:"column:invite_top;type:smallint(5);not null;default:0;comment:邀约业绩TOP"`

	MedalRule  ExcitationMedalRuleModel    `json:"medal_rule" gorm:"foreignKey:RuleCategoryID;references:ID"`
	MedalStyle []ExcitationMedalStyleModel `json:"medal_style" gorm:"foreignKey:MedalID;references:ID"`
}

// TableName returns the table name of the WjMedal model
func (ExcitationMedalModel) TableName() string {
	return "excitation_medal"
}

// 奖章类型 1普通奖章,2隐藏奖章
const (
	MedalMedaTypeOrdinary = 1
	MedalMedaTypeRare     = 2
)

// 奖章来源 1手动上传,2规则统计
const (
	MedalSourceManualUpload   = 1
	MedalSourceRuleStatistics = 2
)

// 规则有效期 1长期有效,2指定时间,3仅生效一次"
const (
	MedalRuleTypeLong    = 1
	MedalRuleTypeSpecify = 2
	MedalRuleTypeOnlyOne = 3
)

// 状态 1启用,2禁用
const (
	MedalStatusEnable  = 1
	MedalStatusDisable = 2
)
