package excitationResp

type PointsRecordList struct {
	Id               int64   `json:"id"`
	IntegralRuleName string  `json:"integral_rule_name"`
	RuleCategoryID   int64   `json:"rule_category_id"`
	IncentiveType    int64   `json:"incentive_type"`
	RuleCategoryName string  `json:"rule_category_name"`
	PointsRemark     string  `json:"points_remark"`
	FyTimeStr        string  `json:"fy_time_str"`
	Points           float64 `json:"points"`
	SourceType       int64   `json:"source_type"`
	CreatedAt        int64   `json:"created_at"`
	CreatedTimeStr   string  `json:"created_time_str"`
	IncentiveTypeStr string  `json:"incentive_type_str"`
	SourceTypeStr    string  `json:"source_type_str"`
	Coins            float64 `json:"coins"`
	MemberInfo
}

type MemberInfo struct {
	MemberUserId   string `json:"member_user_id"`
	MemberName     string `json:"member_name"`
	DepartmentName string `json:"department_name"`
	ThirdId        int    `json:"third_id"`
}

type GetPointsListResp struct {
	Id            int64  `json:"id"`
	PointsRemark  string `json:"points_remark"`
	Points        string `json:"points"`
	CreatedAt     int64  `json:"created_at"`
	IncentiveType int64  `json:"incentive_type"`
}

type MemberSumPointsList struct {
	MemberId  int64   `json:"member_id"`
	SumPoints float64 `json:"sum_points"`
}

type SumPointsResp struct {
	CategoryId int64   `json:"category_id"`
	SumPoints  float64 `json:"sum_points"`
}
