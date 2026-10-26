package settings

type ContactRequest struct {
	ExternalSecret string `json:"external_secret"`
	EnterpriseID   int64  `json:"enterprise_id"`
}

type WxqyTokenRequest struct {
	WxqyToken    string `json:"wxqy_token"`
	EnterpriseID int64  `json:"enterprise_id"`
}

type PluginRequest struct {
	AgentID      int64  `json:"agent_id"`
	CheckFile    int8   `json:"check_file"`
	Secret       string `json:"secret"`
	EnterpriseID int64  `json:"enterprise_id"`
	Type         int8   `json:"type"`
	FileUrl      string `json:"file_url"`
}

type CallbackRequest struct {
	AgentID      int64  `json:"agent_id"`
	CheckFile    int8   `json:"check_file"`
	Secret       string `json:"secret"`
	EnterpriseID int64  `json:"enterprise_id"`
}

type AuthRequest struct {
	EnterpriseID int64 `json:"enterprise_id"`
	AdminID      int64 `json:"admin_id"`
}

type ContactResponse struct {
	AuthUrl        string `json:"auth_url"`
	EnterpriseID   int64  `json:"enterprise_id"`
	IsOpen         int8   `json:"is_open"`
	Installed      bool   `json:"installed"`
	CorpID         string `json:"corp_id"`
	ExternalSecret string `json:"external_secret"`
}

type SettingResponse struct {
	ExternalSecret         string `json:"external_secret"`
	AgentID                int64  `json:"agent_id"`
	Secret                 string `json:"secret"`
	Domain                 string `json:"domain"`
	RedPacketsUrl          string `json:"red_packets_url"`
	CustomerInformationUrl string `json:"customer_information_url"`
	MaterialLibraryUrl     string `json:"material_library_url"`
	MarketingPromotionUrl  string `json:"marketing_promotion_url"`
	FissionActivityUrl     string `json:"fission_activity_url"`
	Key                    string `json:"key"`
	Token                  string `json:"token"`
	UcCallbackUrl          string `json:"uc_callback_url"`
	IsSuper                int    `json:"is_super"`
	CheckFile              int8   `json:"check_file"`
	FileUrl                string `json:"file_url"`
}

// 会话存档配置实体
type ChatSettingsResponse struct {
	Id                int64    `json:"id" db:"id" desc:"ID"`
	Secret            string   `json:"secret" db:"secret" desc:"会话存档secret"`
	MessagePublicKey  string   `json:"message_public_key" db:"message_public_key" desc:"消息加密公钥"`
	MessagePrivateKey string   `json:"message_private_key" db:"message_private_key" desc:"消息加密私钥"`
	EnterpriseId      int64    `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt         int64    `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt         int64    `json:"updated_at" db:"updated_at" desc:"更新时间"`
	TrustedIps        []string `json:"trusted_ips"  desc:"TrustedIps"`
}
