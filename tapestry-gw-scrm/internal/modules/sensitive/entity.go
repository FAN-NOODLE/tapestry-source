package sensitive

// ListRequest 获取敏感词列表请求参数
type ListRequest struct {
	SensitiveText string `json:"sensitive_text" desc:"敏感词搜索"`
	EnterpriseID  int64  `json:"enterprise_id" desc:"企业ID"`
	Page          int    `json:"page" desc:"页码"`
	Size          int    `json:"size" desc:"分页"`
}

// ResponseList 返回敏感词列表
type ResponseList struct {
	List  []Detail `json:"list"`
	Total int64    `json:"total"`
	Page  int      `json:"page"`
	Size  int      `json:"size"`
}

// Detail 敏感词字段
type Detail struct {
	ID              int64  `json:"id" desc:"主键id"`
	SensitiveText   string `json:"sensitive_text" desc:"敏感词"`
	SensitiveType   int8   `json:"sensitive_type" desc:"敏感词类型 1 名片 2 文本"`
	EnterpriseID    int64  `json:"enterprise_id" desc:"企业id"`
	TriggerNum      int32  `json:"trigger_num" desc:"总触发次数"`
	TodayTriggerNum int32  `json:"today_trigger_num" desc:"今日触发次数"`
	CreatedAt       int64  `json:"created_at" desc:"创建时间"`
}

// CountBySensitive 获取敏感词每天对应的聊天数
type CountBySensitive struct {
	SensitiveID     int64 `json:"sensitive_id" desc:"敏感词id"`
	TodayTriggerNum int32 `json:"today_trigger_num" desc:"今日触发次数"`
}
