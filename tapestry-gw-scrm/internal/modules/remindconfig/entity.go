package remindconfig

// ConfigInfoRequest 提醒设置信息获取（结构体）
type ConfigInfoRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	NoticeType   int32 `json:"notice_type" db:"notice_type" validate:"required" required:"查询类型不能为空" desc:"提醒开关类型 1-超时 2-敏感词"`
}

// ConfigInfoResponse 提醒设置信息获取（结构体）
type ConfigInfoResponse struct {
	ID           int64     `json:"id" db:"id" desc:"ID"`
	OverTime     int32     `json:"overtime" db:"overtime" desc:"聊天记录超时分钟"`
	ReceiveStart string    `json:"receive_start" db:"receive_start" desc:"接待开始时间"`
	ReceiveEnd   string    `json:"receive_end" db:"receive_end" desc:"接待结束时间"`
	NoticeStatus int32     `json:"notice_status" db:"notice_status" desc:"提醒开关"`
	StatTxt      string    `db:"stat_txt" desc:"设置员工字段json,企微userid"`
	StatList     []StatTxt `json:"stat_list"`
}

type StatTxt struct {
	Id   int64  `json:"id" desc:"id"`
	Name string `json:"name" desc:"name"`
}
