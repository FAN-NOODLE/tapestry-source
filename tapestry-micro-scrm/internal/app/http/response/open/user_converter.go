package openResp

// ConvertOpenUserIdResponse 将open_userid转为userid响应
type ConvertOpenUserIdResponse struct {
	UserID string `json:"user_id"` // 自建应用的userid
}
