package openReq

type LoginRequest struct {
	SecretId  int64  `json:"secret_id" binding:"required" required:"秘钥Id必填"`
	SecretKey string `json:"secret_key" binding:"required" binding:"秘钥必填"`
}
