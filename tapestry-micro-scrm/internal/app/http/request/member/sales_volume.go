package memberReq

type MemberSalesVolumeReq struct {
	EnterpriseId int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
	StartTime    int64 `json:"start_time"`
	EndTime      int64 `json:"end_time"`
	SearchType   int64 `json:"search_type"`
}

type MemberSalesVolume struct {
	MemberId int64   `json:"member_id"`
	Total    float64 `json:"total"`
}

type AchievementParams struct {
	EnterpriseId int64   `json:"enterprise_id"`
	MemberIds    []int64 `json:"member_ids"`
	SearchType   int64   `json:"search_type"`
	StartTime    int64   `json:"start_time"`
	EndTime      int64   `json:"end_time"`
}
