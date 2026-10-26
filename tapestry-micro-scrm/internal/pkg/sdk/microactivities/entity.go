package microactivities

type ActivitiesByIdsRequest struct {
	Type         int8    `json:"type"`
	RuleIds      []int64 `json:"rule_ids"`
	EnterpriseID int64   `json:"enterprise_id"`
}

type ActivitiesByIdsResponse struct {
	ResponseBase
	Data []Activity `json:"data"`
}

type Activity struct {
	ActivityRuleID int64  `json:"activity_rule_id"`
	ActivityType   int8   `json:"activity_type"`
	Name           string `json:"name"`
	StartAt        int64  `json:"start_at"`
	EndAt          int64  `json:"end_at"`
}

type ActivityPagesRequest struct {
}

type ActivityPagesResponse struct {
	ResponseBase
	Data []ActivityPage `json:"data"`
}

type ActivityPage struct {
}
