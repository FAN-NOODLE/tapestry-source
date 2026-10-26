package excitationReq

type MedalRecordList struct {
	MemberUserId    string `json:"member_user_id"`
	MemberName      string `json:"member_name"`
	MedalRuleName   string `json:"medal_rule_name"`
	Source          int    `json:"source"`
	DepartmentId    int64  `json:"department_id"`
	RuleDimensionId int    `json:"rule_dimension_id"`
	MedalName       string `json:"medal_name"`
	StartTime       int    `json:"start_time"`
	EndTime         int    `json:"end_time"`
	MedalId         int64  `json:"medal_id"`
	EnterpriseId    int64  `json:"enterprise_id"`
	Page            int    `json:"page"`
	Size            int    `json:"size"`
}

type MedalRecordReissue struct {
	EnterpriseId int64                    `json:"enterprise_id"`
	ReissueData  []MedalRecordReissueData `json:"reissue_data" binding:"required,dive"`
}
type MedalRecordReissueData struct {
	MemberUserId string `json:"member_user_id" binding:"required"`
	MedalRuleId  int64  `json:"medal_rule_id" binding:"required"`
	MedalId      int64  `json:"medal_id" binding:"required"`
}
