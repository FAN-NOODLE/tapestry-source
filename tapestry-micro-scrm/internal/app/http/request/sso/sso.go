package ssoReq

// SSOLoginReq SSO登录请求
type SSOLoginReq struct {
	EnterpriseId int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
	Code         string `json:"code" binding:"required" comment:"Microsoft OAuth 授权码"`
}
