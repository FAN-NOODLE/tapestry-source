package microScrm

// OvertimeInfoRequest 超时提醒设置信息获取（结构体）
type ConfigInfoRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	NoticeType   int32 `json:"notice_type" db:"notice_type" validate:"required" required:"查询类型不能为空" desc:"提醒开关类型 1-超时 2-敏感词"`
}

// OvertimeInfoResponse 超时提醒设置信息获取（结构体）
type ConfigInfoResponse []struct {
	ID           int64     `json:"id" db:"id" desc:"ID"`
	OverTime     int32     `json:"overtime" db:"overtime" desc:"聊天记录超时分钟"`
	ReceiveStart string    `json:"receive_start" db:"receive_start" desc:"接待开始时间"`
	ReceiveEnd   string    `json:"receive_end" db:"receive_end" desc:"接待结束时间"`
	NoticeStatus int32     `json:"notice_status" db:"notice_status" desc:"提醒开关"`
	StatTxt      string    `db:"stat_txt" desc:"设置员工字段json,企微userid"`
	StatList     []StatTxt `json:"stat_list"`
}

type StatTxt []struct {
	Id     int64  `json:"member_id" desc:"member_id"`
	Name   string `json:"name" desc:"name"`
	UserId string `json:"user_id" desc:"user_id"`
}

// OvertimeInfoRequest 超时提醒设置信息设置（结构体）
type ConfigSaveRequest struct {
	ID           int64  `json:"id" db:"id" desc:"ID"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业id" validate:"required" required:"企业ID不能为空"`
	OverTime     int32  `json:"overtime" db:"overtime" desc:"聊天记录超时分钟"`
	ReceiveStart string `json:"receive_start" db:"receive_start" desc:"接待开始时间"`
	ReceiveEnd   string `json:"receive_end" db:"receive_end" desc:"接待结束时间"`
	NoticeType   int32  `json:"notice_type"  db:"notice_type" validate:"required" required:"设置类型不能为空" desc:"提醒开关 提醒开关类型 1-超时 2-敏感词"`
	NoticeStatus int32  `json:"notice_status" db:"notice_status" validate:"required" required:"提醒开关必须选择" desc:"提醒开关 1-开启 2-关闭"`
	StatTxt      string `json:"stat_txt" db:"stat_txt" desc:"设置详情"`
}

// ChatSettings 会话存档配置实体
type ChatSettings struct {
	Id                int64    `json:"id" db:"id" desc:"ID"`
	Secret            string   `json:"secret" db:"secret" desc:"会话存档secret"`
	MessagePublicKey  string   `json:"message_public_key" db:"message_public_key" desc:"消息加密公钥"`
	MessagePrivateKey string   `json:"message_private_key" db:"message_private_key" desc:"消息加密私钥"`
	EnterpriseID      int64    `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	TrustedIps        []string `json:"trusted_ips" desc:"可信ip"`
	CreatedAt         int64    `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt         int64    `json:"updated_at" db:"updated_at" desc:"更新时间"`
}

// ConfigInfoRequest 超时提醒设置信息获取（结构体）
type SettingsInfoRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
}
