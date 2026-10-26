package reserveReq

type ServiceScoreListReq struct {
	EnterpriseBase
	PageReq
	MemberIds        []int64 `json:"member_ids" comment:"服务导购的id"`
	StoreIds         []int64 `json:"store_ids" comment:"服务的门店id"`
	ReserveStartTime int     `json:"reserve_start_time" comment:"预约开始时间"`
	ReserveEndTime   int     `json:"reserve_end_time" comment:"预约结束时间"`
	ScoreStartTime   int     `json:"score_start_time" comment:"评分开始时间"`
	ScoreEndTime     int     `json:"score_end_time" comment:"评分结束时间"`
	Ids              []int64 `json:"ids" comment:"ids"`
}
