package common

import "encoding/json"

// ContextKeyEnum ctx.Context中存储的数据key枚举
var ContextKeyEnum = struct {
	TokenInfo      string
	EnterpriseHash string
	EnterpriseID   string
}{
	TokenInfo:      "token_info",
	EnterpriseHash: "enterprise_hash",
	EnterpriseID:   "enterprise_id",
}

// HeaderKeyEnum 请求Header中存储的key枚举
var HeaderKeyEnum = struct {
	AuthToken      string
	EnterpriseHash string
	EnterpriseID   string
	XRequestID     string
}{
	AuthToken:      "Auth-Token",
	EnterpriseHash: "Enterprise-Hash",
	EnterpriseID:   "Enterprise-ID",
	XRequestID:     "X-Request-ID",
}

// EnterpriseInfo 企业信息结构体
type EnterpriseInfo struct {
	ID             int64  `json:"id"`
	EnterpriseID   int64  `json:"enterprise_id"`
	Hash           string `json:"hash"`
	Name           string `json:"name"`
	Logo           string `json:"logo"`
	IndustryID     int64  `json:"industry_id"`
	CorpID         string `json:"corp_id"`
	Domain         string `json:"domain"`
	Status         int    `json:"status"`
	BusinessSource int    `json:"business_source"`
	Version        int    `json:"version"`
	VersionTag     string `json:"version_tag"`
	IsSplitFlow    int    `json:"is_split_flow"`
	Identifier     string `json:"identifier"`
	BrandCode      string `json:"brand_code"`
}

// MarshalBinary 将结构体转码为字符串
func (e *EnterpriseInfo) MarshalBinary() (data []byte, err error) {
	return json.Marshal(e)
}

// UnmarshalBinary 解码将字符串转为结构体
func (e *EnterpriseInfo) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, e)
	if err != nil {
		return err
	}
	return nil
}
