package common

import (
	"encoding/json"
)

type TokenInfo struct {
	ExpireTime   int64  `json:"expire_time"`
	IsTmp        int    `json:"is_tmp"`
	EnterpriseId int64  `json:"enterprise_id"`
	AdminId      int64  `json:"admin_id"`
	ServeType    int    `json:"serve_type"`
	Mobile       string `json:"mobile"`
	Username     string `json:"username"`
}

func (t *TokenInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

// ctx.Context中存储的数据key枚举
var ContextKeyEnum = struct {
	TokenInfo      string
	EnterpriseHash string
}{
	TokenInfo:      "token_info",
	EnterpriseHash: "enterprise_hash",
}

// 请求Header中存储的key枚举
var HeaderKeyEnum = struct {
	AuthToken      string
	EnterpriseHash string
	EnterpriseID   string
	XRequestID     string
}{
	AuthToken:      "Auth-Token",
	EnterpriseHash: "Enterprise-Hash",
	EnterpriseID:   "Enterprise-Id",
	XRequestID:     "X-Request-ID",
}
