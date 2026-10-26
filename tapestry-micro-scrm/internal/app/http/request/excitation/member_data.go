package excitationReq

import "micro-scrm/internal/app/http/request"

type MemberDataListReq struct {
	request.EnterpriseBaseReq
	request.PageReq
	request.MemberInfoReq
	MemberThirdID         int    `json:"member_third_id"  comment:"third_id"`
	Level                 string `json:"level"  comment:"等级名称"`
	PointUpdatedStartTime int    `json:"point_updated_start_time"  comment:"积分更新时间"`
	PointUpdatedEndTime   int    `json:"point_updated_end_time" comment:"积分更新时间"`
	LevelUpdatedStartTime int    `json:"level_updated_start_time"  comment:"等级更新时间"`
	LevelUpdatedEndTime   int    `json:"level_updated_end_time" comment:"等级更新时间"`
	MedalUpdatedStartTime int    `json:"medal_updated_start_time"  comment:"等级更新时间"`
	MedalUpdatedEndTime   int    `json:"medal_updated_end_time" comment:"等级更新时间"`
	PointValue            string `json:"point_value"  comment:"积分数量排序"`
	Coins                 string `json:"coins"  comment:"金币排序"`
	SpendCoins            string `json:"spend_coins"  comment:"累计消耗金币排序"`
	TotalCoins            string `json:"total_coins"  comment:"累计总金币排序"`
	MedalNum              string `json:"medal_num"  comment:"奖章数量排序"`
}

//MemberDataIndex 小程序导购积分&等级
type MemberDataIndex struct {
	request.EnterpriseBaseReq
	Mobile string `json:"mobile"`
}

//GradePoints 导购等级积分
type GradePoints struct {
	request.EnterpriseBaseReq
	Mobile string `json:"mobile"`
}

// UpdatePointsReq  导购积分信息更新
type UpdatePointsReq struct {
	request.EnterpriseBaseReq
	MemberId int64   `json:"member_id"`
	Points   float64 `json:"points"`
}

// UpdateCoinsReq  导购金币信息更新
type UpdateCoinsReq struct {
	request.EnterpriseBaseReq
	MemberId int64   `json:"member_id"`
	Coins    float64 `json:"coins"`
}

// UpdateLevelReq  导购等级更新
type UpdateLevelReq struct {
	request.EnterpriseBaseReq
	MemberId int64  `json:"member_id"`
	Level    string `json:"level"`
	LevelId  int64  `json:"level_id"`
}

// UpdateMedalReq  导购奖章更新
type UpdateMedalReq struct {
	request.EnterpriseBaseReq
	MemberId int64 `json:"member_id"`
	MedalNum int   `json:"medal_num"`
}

// UpdateDimensionFieldReq  导购各维度字段更新
type UpdateDimensionFieldReq struct {
	request.EnterpriseBaseReq
	Id                  int64  `json:"id"`
	MemberId            int64  `json:"member_id"`
	IncentiveConfigType int    `json:"incentive_config_type"`
	DimensionValue      string `json:"dimension_value"`
	DimensionRuleValue  int    `json:"dimension_rule_value"`
	DimensionRuleStr    string `json:"dimension_rule_str"`
}

// RelationRuleArray 规则能力获得值数组
type RelationRuleArray struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

// GetPointsRulesReq 小程序获取积分规则信息
type GetPointsRulesReq struct {
	request.EnterpriseBaseReq
	RuleCategoryId []int64 `json:"rule_category_id"`
	MemberId       int64   `json:"member_id"`
}

// LevelRecordReq 小程序获取导购等级记录及最高等级时间
type LevelRecordReq struct {
	request.EnterpriseBaseReq
	MemberId int64 `json:"member_id"`
}
