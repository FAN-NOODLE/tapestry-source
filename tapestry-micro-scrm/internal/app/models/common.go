package models

type TopList struct {
	MemberId     int64  `json:"member_id"`
	Total        int64  `json:"total"`
	Ranking      string `json:"ranking"`
	RankingValue string `json:"ranking_value"`
}

type BehaviorList struct {
	MemberId int64 `json:"member_id"`
	Total    int64 `json:"total"`
}

type MemberTopList struct {
	MemberId int64 `json:"member_id"`
	Total    int64 `json:"total"`
	Top      int   `json:"top"`
}

type MemberTopSortList struct {
	MemberId     int64   `json:"member_id"`
	MemberUserid string  `json:"member_userid"`
	Value        float64 `json:"value" comment:"金额或积分"`
	Ratio        float64 `json:"ratio" comment:"占比"`
	Total        int64   `json:"total" comment:"数量"`
	SortTop      int     `json:"sort_top" comment:"排序"`
}
