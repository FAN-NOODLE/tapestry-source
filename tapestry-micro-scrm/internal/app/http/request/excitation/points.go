package excitationReq

import "micro-scrm/internal/app/http/request"

type ManualPointsReq struct {
	MedalPointsData []ManualPointsRow `json:"medal_points_data"`
	EnterpriseId    int64             `json:"enterprise_id"`
}

type RecordListReq struct {
	request.MemberSearchReq
	IntegralRuleName string `json:"integral_rule_name"`
	IncentiveType    int64  `json:"incentive_type"`
	SourceType       int64  `json:"source_type"`
	StartTime        int64  `json:"start_time"`
	EndTime          int64  `json:"end_time"`
	Page             int    `json:"page"`
	Size             int    `json:"size"`
	HasCoins         bool   `json:"has_coins"`
}

type GetPointsListReq struct {
	request.EnterpriseBaseReq
	MemberId int64 `json:"member_id"`
	Page     int   `json:"page"`
	Size     int   `json:"size"`
}

type GetCoinsListReq struct {
	GetPointsListReq
}
type CoinsInfoReq struct {
	request.EnterpriseBaseReq
	MemberId int64 `json:"member_id"`
}

type SumPointsReq struct {
	MemberId     int64 `json:"member_id"`
	EnterpriseId int64 `json:"enterprise_id"`
	Type         int64 `json:"type"`
	StartTime    int64 `json:"start_time"`
	EndTime      int64 `json:"end_time"`
}

type SumPointsCategoryReq struct {
	MemberId     int64   `json:"member_id"`
	EnterpriseId int64   `json:"enterprise_id"`
	CategoryIds  []int64 `json:"category_ids"`
	StartTime    int64   `json:"start_time"`
	EndTime      int64   `json:"end_time"`
}

type TaskFinishedReq struct {
	TaskAllotLogId int64  `json:"task_allot_log_id"`
	MemberId       int64  `json:"member_id"`
	EnterpriseId   int64  `json:"enterprise_id"`
	Mobile         string `json:"mobile"`
}
