package activityCampaign

// 接收参数实体

type SignReceiveInfo struct {
	ID        int64 `json:"id"`         //签到活动id
	StartTime int64 `json:"start_time"` //开始时间
	Type      int   `json:"type"`
	EndTime   int64 `json:"end_time"` // 结束时间
	Page      Page  `json:"page"`     //分页
}

type SignDetailInfo struct {
	ID        int64  `json:"id"` //签到活动的id
	UnionID   string `json:"union_id"`
	StartTime int64  `json:"start_time"` //开始时间
	EndTime   int64  `json:"end_time"`   // 结束时间
	Page      Page   `json:"page"`       //分页
}

type Page struct {
	PageSize int `json:"page_size"`
	PageNum  int `json:"page_num"`
}

type UserSignInLogs struct {
	Date string `json:"date"`
}
