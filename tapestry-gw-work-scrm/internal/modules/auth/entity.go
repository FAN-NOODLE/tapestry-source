package auth

import microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"

// LoginAuth 授权登录请求
type LoginAuthRequest struct {
	EnterpriseID int64  `json:"enterprise_id" desc:"企业ID"`
	Hash         string `json:"hash"`
	Code         string `json:"code" desc:"授权码"`
	LoginType    string `json:"login_type" desc:"登录方式"`
}

// LoginAuth 授权登录响应
type LoginAuthResponse struct {
	*microScrm.AuthInfoResponse
	Token     string `json:"token" desc:"身份凭据"`
	BrandCode string `json:"brand_code"`
}

// RedirectUrlRequest 授权地址
type RedirectUrlRequest struct {
	EnterpriseID int64  `json:"enterprise_id" desc:"企业ID"`
	Hash         string `json:"hash"`
	RedirectUri  string `json:"redirect_uri" desc:"授权地址" validate:"required" required:"授权地址必填"`
}

// RedirectUrlResponse 授权登录响应
type RedirectUrlResponse struct {
	Url string `json:"url" desc:"授权地址"`
}
