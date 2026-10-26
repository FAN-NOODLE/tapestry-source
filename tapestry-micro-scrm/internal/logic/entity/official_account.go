package entity

// js-config
type JsConfig struct {
	AppID     string `json:"appId"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonce_str"`
	Signature string `json:"signature"`
}

type ShuJsConfig struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    ShuData `json:"data"`
}
type ShuData struct {
	AppID     string `json:"AppID"`
	Timestamp string `json:"Timestamp"`
	NonceStr  string `json:"NonceStr"`
	Signature string `json:"Signature"`
}

type ReqShuJsConfig struct {
	Id        string `json:"id"`
	Slat      string `json:"slat"`
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
	Url       string `json:"url"`
}

type ErrorResponse struct {
	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

// 欧莱雅第三方平台接口的access-token
type ThirdTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type ThirdTokenResponse struct {
	ErrorResponse
	AccessToken string `json:"access_token"`
	ExpireIn    int64  `json:"expires_in"`
}

// 欧莱雅第三方平台接口的js-ticket
type ThirdJsTicketRequest struct {
	AccessToken string `json:"access_token"`
}

type ThirdJsTicketResponse struct {
	ErrorResponse
	Data ThirdTicketData `json:"data"`
}

type ThirdTicketData struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}
