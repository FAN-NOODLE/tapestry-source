package excitationModels

import "micro-scrm/internal/app/models"

const (
	TimeRule1 = iota + 1 // 时间规则: 1每次(人) 2每天 3每周 4每月 5每季度 6每年
	TimeRule2
	TimeRule3
	TimeRule4
	TimeRule5
	TimeRule6

	StatusTypeOn  = 1 // 开启状态: 1开启 2关闭
	StatusTypeOff = 2

	IsDiffStore_No  = 0
	IsDiffStore_Yes = 1
)

type ExcitationIncentiveRelationModel struct {
	models.MicroBaseModel

	IncentiveConfigType int `json:"incentive_config_type" gorm:"column:incentive_config_type;type:tinyint(4);not null;default:0;comment:incentive_config表type"`
	Status              int `json:"status" gorm:"column:status;type:tinyint(4);not null;default:0;comment:开启状态: 1开启 2关闭"`
	AbilityIndex        int `json:"ability_index" gorm:"column:ability_index;type:tinyint(4);not null;default:0;comment:能力指标条件: 1邀请目标达成率(每月) 2任务完成率(每月) 3总客户互动积分排行榜(每周) 4总培训积分排行榜(每周) 5目标达成率(每月)"`
	TimeRule            int `json:"time_rule" gorm:"column:time_rule;type:tinyint(4);not null;default:0;comment:时间规则: 1每次(人) 2每天 3每周 4每月 5每季度 6每年"`
	IsDiffStore         int `json:"is_diff_store" gorm:"column:is_diff_store;type:tinyint(4);not null;default:0;comment:是否区分奥莱零售 0 不区分 1区分奥莱零售"`
	// 关联维度配置规则
	IncentiveRelationRules []ExcitationIncentiveRelationRuleModel `json:"incentive_relation_rules" gorm:"foreignKey:IncentiveRelationID;references:ID"`
}

// TableName returns the table name of the ExcitationIncentiveRelation model
func (ExcitationIncentiveRelationModel) TableName() string {
	return "excitation_incentive_relation"
}
