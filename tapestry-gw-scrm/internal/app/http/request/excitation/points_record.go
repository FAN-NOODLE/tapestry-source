package excitationReq

type PointsRecordCreate struct {
	Name         string `json:"name"`
	MemberId     int64  `json:"member_id"`
	Type         int64  `json:"type"`
	Value        int64  `json:"value"`
	AverageValue int64  `json:"average_value"`
	UniqueHash   string `json:"unique_hash"`
	EnterpriseId int64  `json:"enterprise_id"`
}

type PointsRecordReq struct {
	UserId           string `json:"user_id"`
	UserName         string `json:"user_name"`
	ThirdId          int64  `json:"third_id"`
	EnterpriseId     int64  `json:"enterprise_id"`
	IntegralRuleName string `json:"integral_rule_name"`
	IncentiveType    int64  `json:"incentive_type"`
	SourceType       int64  `json:"source_type"`
	StartTime        int64  `json:"start_time"`
	EndTime          int64  `json:"end_time"`
	Page             int    `json:"page"`
	Size             int    `json:"size"`
}
