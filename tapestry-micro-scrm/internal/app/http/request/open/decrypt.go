package openReq

// DecryptReq PI 信息批量解密请求。
type DecryptReq struct {
	Items []DecryptItem `json:"items"`
}

// DecryptItem 单个待解密字段。
type DecryptItem struct {
	RequestKey       string `json:"request_key"`
	Ciphertext       string `json:"ciphertext"`
	EncryptedDataKey string `json:"encrypted_data_key"`
	IV               string `json:"iv"`
	AuthTag          string `json:"auth_tag"`
	Algorithm        string `json:"algorithm"`
}
