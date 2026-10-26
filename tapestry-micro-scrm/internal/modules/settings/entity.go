package settings

// 会话存档配置实体
type ChatSettings struct {
	Id                int64    `json:"id" db:"id" desc:"ID"`
	Secret            string   `json:"secret" db:"secret" desc:"会话存档secret"`
	MessagePublicKey  string   `json:"message_public_key" db:"message_public_key" desc:"消息加密公钥"`
	MessagePrivateKey string   `json:"message_private_key" db:"message_private_key" desc:"消息加密私钥"`
	Limit             int64    `json:"limit" db:"limit" desc:"一次拉多少条消息"`
	EnterpriseId      int64    `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt         int64    `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt         int64    `json:"updated_at" db:"updated_at" desc:"更新时间"`
	TrustedIps        []string `json:"trusted_ips" desc:"可信ip"`
}

// ConfigInfoRequest 超时提醒设置信息获取（结构体）
type SettingsInfoRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
}

type GenRsaKeyInfo struct {
	PrivateKey string `json:"private_key" desc:"私钥"`
	PublicKey  string `json:"public_key" desc:"公钥"`
}
