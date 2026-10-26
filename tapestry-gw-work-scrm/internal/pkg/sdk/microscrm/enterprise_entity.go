package microScrm

// AgentConfigRequest 签名请求结构体
type JsConfigRequest struct {
	Url          string `json:"url" desc:"url"`
	Timestamp    int64  `json:"timestamp" desc:"时间戳"`
	Nonce        string `json:"nonce" desc:"随机字符串"`
	EnterpriseID int64  `json:"enterprise_id"`
	Hash         string `json:"hash"`
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
