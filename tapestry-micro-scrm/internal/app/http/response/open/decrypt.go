package openResp

// DecryptResp PI 信息批量解密响应。
type DecryptResp struct {
	Items []DecryptResult `json:"items"`
}

// DecryptResult 单个字段的解密结果。
type DecryptResult struct {
	RequestKey   string `json:"request_key"`
	Plaintext    string `json:"plaintext"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message"`
}
