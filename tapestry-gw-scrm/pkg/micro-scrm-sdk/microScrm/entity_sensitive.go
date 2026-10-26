package microScrm

// SensitiveListRequest 获取敏感词列表请求参数
type SensitiveListRequest struct {
	SensitiveText string `json:"sensitive_text" db:"sensitive_text" desc:"敏感词搜索"`
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id" binding:"required" desc:"企业ID"`
	Page          int    `json:"page" desc:"页码"`
	Size          int    `json:"size" desc:"分页"`
}

// SensitiveList 敏感词字段
type SensitiveListResponse []struct {
	ID              int64  `json:"id" db:"id" desc:"主键id"`
	SensitiveText   string `json:"sensitive_text" db:"sensitive_text" desc:"敏感词"`
	SensitiveType   int8   `json:"sensitive_type" db:"sensitive_type" desc:"敏感词类型 1 名片 2 文本"`
	EnterpriseID    int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	TriggerNum      int32  `json:"trigger_num" db:"trigger_num" desc:"总触发次数"`
	TodayTriggerNum int32  `json:"today_trigger_num" desc:"今日触发次数"`
	CreatedAt       int64  `json:"created_at" db:"created_at" desc:"创建时间"`
}

// SensitiveSaveRequest 敏感词保存请求参数
type SensitiveSaveRequest struct {
	ID            int64  `json:"id" db:"id" desc:"ID"`
	SensitiveType int8   `json:"sensitive_type" db:"sensitive_type" validate:"required" required:"敏感词类型不能为空" desc:"敏感词类型 1 名片 2 文本"`
	SensitiveText string `json:"sensitive_text" db:"sensitive_text" desc:"敏感词"`
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空" desc:"企业ID"`
	TriggerNum    int32  `json:"trigger_num" db:"trigger_num" desc:"总触发次数"`
}

// SensitiveDeleteRequest 敏感词删除请求参数
type SensitiveDeleteRequest struct {
	ID int64 `json:"id" db:"id"  validate:"required" required:"敏感词ID不能为空" desc:"ID"`
}
