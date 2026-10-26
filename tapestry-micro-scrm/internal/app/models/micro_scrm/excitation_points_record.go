package microScrmModel

import "micro-scrm/internal/app/models"

const (
	ExcitationPointsRecordAutoType_Auto = iota + 1
	ExcitationPointsRecordAutoType_Manual
	ExcitationPointsRecordAutoType_Medal
)

type ExcitationPointsRecordModel struct {
	models.BaseModel
	IncentiveType          int64                      `json:"incentive_type" gorm:"column:incentive_type;type:tinyint(1);not null;index:incentive_type_index;default:0;comment:维度类型"`
	MemberID               int64                      `json:"member_id" gorm:"column:member_id;type:bigint(10) unsigned;not null;index:member_id_type_deleted_at_index;default:0;comment:导购member表id"`
	Type                   int64                      `json:"type" gorm:"column:type;type:tinyint(1) unsigned;not null;index:member_id_type_deleted_at_index;default:0;comment:触发类型:1:邀约数(每次) 2:邀约到店率(每周) 3:邀约成功率(每月) 4:邀约成功(每次) 5:邀约占比(每周) 6:邀约到店确认(每次 )7:任务完成(每次) 8:任务完成率(每月) 9:任务触达(每人) 10: 新增好友数(每人) 11:客户流失数(每人) 12:聊天总数(每月) 13:已回复聊天占比(每月) 14:完成播种按(次) 15:积分等比兑换（每分） 16:线上销售业绩率（每周）"`
	Value                  int64                      `json:"value" gorm:"column:value;type:int(10) unsigned;not null;default:0;comment:增加积分时的用户值(百分比值存百分比整数）"`
	ConditionValue         int64                      `json:"condition_value" gorm:"column:condition_value;type:int(10) unsigned;not null;default:0;comment:增加积分的条件值(百分比值存百分比整数）"`
	Points                 string                     `json:"points" gorm:"column:points;type:varchar(64);not null;comment:通过该记录获得的积分"`
	SourceType             int64                      `json:"source_type" gorm:"column:source_type;type:tinyint(1);not null;index:source_type_index;default:0;comment:0 自动（规则）1手动 2勋章"`
	IntegralRuleID         int64                      `json:"integral_rule_id" gorm:"column:integral_rule_id;type:bigint(20);not null;default:0;comment:消费时关联的规则id"`
	IntegralRuleName       string                     `json:"integral_rule_name" gorm:"column:integral_rule_name;type:varchar(255);not null;comment:消费时关联的规则名称"`
	IntegralRuleDetailName string                     `json:"integral_rule_detail_name" gorm:"column:integral_rule_detail_name;type:varchar(255);not null;comment:消费时关联的规则类目名称"`
	PointsRemark           string                     `json:"points_remark" gorm:"column:points_remark;type:varchar(255);not null;comment:添加积分的条件"`
	FyTimeStr              string                     `json:"fy_time_str" gorm:"column:fy_time_str;type:varchar(255);not null;comment:积分增加的财年财月相关时间"`
	CoinsRecord            ExcitationCoinsRecordModel `json:"coins_record" gorm:"foreignKey:PointsRecordID;references:ID"`
}

// TableName returns the table name of the ExcitationPointsRecord model
func (ExcitationPointsRecordModel) TableName() string {
	return "excitation_points_record"
}
