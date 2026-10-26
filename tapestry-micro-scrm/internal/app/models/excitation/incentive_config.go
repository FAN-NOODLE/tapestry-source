package excitationModels

import "micro-scrm/internal/app/models"

const (
	IncentiveConfigTypeInvite   = iota + 1 // 邀约维度
	IncentiveConfigTypeTask                // 任务维度
	IncentiveConfigTypeInteract            // 客户互动维度
	IncentiveConfigTypeTrain               // 培训维度
	IncentiveConfigTypeOnline              // 线上能力

	IncentiveConfigStatusOn  = 1 // 开启状态: 1开启 2关闭
	IncentiveConfigStatusOff = 2
)

type ExcitationIncentiveConfigModel struct {
	models.MicroBaseModel

	RuleName string `json:"rule_name" gorm:"column:rule_name;type:varchar(100);not null;comment:规则名称"`
	Status   int32  `json:"status" gorm:"column:status;type:tinyint(4);not null;default:0;comment:开启状态: 1开启 2关闭"`
	Type     int32  `json:"type" gorm:"column:type;type:tinyint(4);not null;default:0;comment:激励维度: 1邀约维度 2任务维度 3客户互动维度 4培训维度 5线上能力"`

	// 关联能力图谱关系
	IncentiveRelation ExcitationIncentiveRelationModel `json:"incentive_relation"  gorm:"foreignKey:Type;references:IncentiveConfigType"`
}

// TableName returns the table name of the ExcitationIncentiveConfig model
func (ExcitationIncentiveConfigModel) TableName() string {
	return "excitation_incentive_config"
}
