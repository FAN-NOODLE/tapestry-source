package excitationReq

type PointsProduceSave struct {
	MemberId     int64  `json:"member_id"`
	Mobile       string `json:"mobile"`
	Type         int64  `json:"type"`
	Value        int64  `json:"value"`
	UniqueHash   string `json:"unique_hash"`
	EnterpriseId int64  `json:"enterprise_id"`
	SourceRemark string `json:"source_remark"` //需要记录积分的特殊来源，比如培训相关，
	AddTime      int64  `json:"add_time"`      //需要记录积分的添加时间，培训相关
}

type GetAverageReq struct {
	EnterpriseId   int64 `json:"enterprise_id"`
	RuleCategoryId int64 `json:"type"`
	MemberId       int64 `json:"member_id"`
}

type RuleCategory struct {
	Name       string `json:"name"`
	Rate       string `json:"rate"`
	Cumulative bool   `json:"cumulative"` //是否是累计型的
}

type ManualPointsSave struct {
	MemberId       int64  `json:"member_id"`
	RuleName       string `json:"rule_name"`
	RuleDetailName string `json:"rule_detail_name"`
	RuleDetailId   int64  `json:"rule_detail_id"`
	RuleId         int64  `json:"rule_id"`
	PointsRemark   string `json:"points_remark"`
	Points         int64  `json:"points"`
	NowTime        int64  `json:"now_time"`
	PointsRecordId int64  `json:"points_record_id"`
	MedalRecordId  int64  `json:"medal_record_id"`
	OrderNum       string `json:"order_num"`
	HasCoins       bool   `json:"has_coins"`
}

type ManualPointsRow struct {
	UserId         string `json:"user_id"`
	RuleId         int64  `json:"rule_id"`
	RuleDetailName string `json:"rule_detail_name"`
	PointsRemark   string `json:"points_remark"`
	Points         int64  `json:"points"`
	HasCoins       bool   `json:"has_coins"`
}

type MemberSumPointsReq struct {
	EnterpriseId int64 `json:"enterprise_id"`
	Type         int   `json:"type"`
	MemberId     int64 `json:"member_id"`
	StartTime    int   `json:"start_time"`
	EndTime      int   `json:"end_time"`
}

type CheckAstrictReq struct {
	PointsProduceSave
	NowTime   int64   `json:"now_time"`
	AddPoints float64 `json:"add_points"`
}
