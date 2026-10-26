package enterprise

// Corp 企业信息
type Corp struct {
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	CorpID        string `json:"corp_id" db:"corp_id"`
	Domain        string `json:"domain" db:"domain"`
	Hash          string `json:"hash" db:"hash"`
	TagInitStatus int8   `json:"init_status" db:"tag_init_status"`
	Reason        string `json:"fail_reason" db:"reason"`
	IsNotCallback int    `json:"is_callback" desc:"需要配置回调1=不需要 0=需要配置"`
}

// CorpListRequest 企业列表
type CorpListRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// AgentConfigRequest 签名请求结构体
type JsConfigRequest struct {
	Url          string `json:"url" desc:"url"`
	Timestamp    int64  `json:"timestamp" desc:"时间戳"`
	Nonce        string `json:"nonce" desc:"随机字符串"`
	EnterpriseID int64  `json:"enterprise_id"`
}

// jsSign 公共签名
type jsSign struct {
	Url       string `json:"url" desc:"url"`
	Timestamp int64  `json:"timestamp" desc:"时间戳"`
	Nonce     string `json:"noncestr" desc:"随机字符串"`
	Sign      string `json:"signature" desc:"签名"`
}

// JsConfigResponse 企业签名响应
type JsConfigResponse struct {
	jsSign
	AppId string `json:"appid" desc:"企业ID"`
}

// AgentConfigResponse 应用签名响应
type AgentConfigResponse struct {
	jsSign
	CorpId  string `json:"corp_id" desc:"企业ID"`
	AgentId int64  `json:"agent_id" desc:"应用ID"`
}

// OfficialRedirectURLReq 授权地址请求参数
type OfficialRedirectURLReq struct {
	EnterpriseID int64  `json:"enterprise_id"  binding:"required"`
	URL          string `json:"url"  binding:"required"`
	Scope        string `json:"scope"  binding:"required"`
	State        string `json:"state"`
}

// OfficialRedirectURLRes 授权地址请求响应参数
type OfficialRedirectURLRes struct {
	URL string `json:"url"`
}

// OfficialInfo 公众号信息
type OfficialInfo struct {
	AppID  string `json:"app_id"`
	Secret string `json:"secret"`
}

// EpConfig 企业配置结构体
type EpConfig struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	Hash         string `json:"hash" db:"hash"`
	Key          string `json:"key" db:"key"`
	Value        string `json:"value" db:"value"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// OfficialAuthReq 公众号授权请求参数
type OfficialAuthReq struct {
	EnterpriseID int64  `json:"enterprise_id"  binding:"required"`
	Code         string `json:"code" binding:"required"`
}

// OfficialUserReq 单挑请求实体
type OfficialUserReq struct {
	ID           int64  `json:"id" db:"id"`
	UnionID      string `json:"union_id" db:"union_id"`
	OpenID       string `json:"open_id" db:"open_id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
}

// OfficialUser 公众号用户信息
type OfficialUser struct {
	ID           int64  `json:"id" db:"id"`
	OpenID       string `json:"open_id" db:"open_id"`
	Nickname     string `json:"nickname" db:"nick_name"`
	Sex          int32  `json:"sex" db:"sex"`
	Province     string `json:"province" db:"province"`
	City         string `json:"city" db:"city"`
	Country      string `json:"country" db:"country"`
	HeadImgURL   string `json:"img_url" db:"img_url"`
	Privilege    string `json:"privilege" db:"privilege"`
	UnionID      string `json:"union_id" db:"union_id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

type ConfigReq struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	Key          string `json:"key" db:"key"`
}

type ConfigValue struct {
	Value string `json:"value"  db:"value"`
}

// WxAppConfigEntity 公众号信息结构体
type WxAppConfigEntity struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	AppID        string `json:"appid" db:"appid"`
	AppType      int    `db:"app_type"`
	Name         string `json:"name" db:"name"`
	BindStatus   int    `json:"bind_status" db:"bind_status"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `db:"updated_at"`
}

// 欧莱雅第三方平台，回传用户信息结构
type ThirdUserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int32    `json:"sex"`
	Language   string   `json:"language"`
	City       string   `json:"city"`
	Province   string   `json:"province"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid"`
}
