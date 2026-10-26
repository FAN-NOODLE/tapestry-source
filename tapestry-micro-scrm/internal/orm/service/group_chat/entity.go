package group_chat

type GropuNameQueryconditions struct {
	EnterpriseID  int64   `json:"enterprise_id"`
	ShoppingGuide []int64 `json:"shopping_guide"`
	StartTime     int64   `json:"start_time"`
	EndTime       int64   `json:"end_time"`
	Page          Page    `json:"page"`
}

// Page 分页参数
type Page struct {
	PageSize int `json:"page_size"`
	PageNum  int `json:"page_num"`
}
