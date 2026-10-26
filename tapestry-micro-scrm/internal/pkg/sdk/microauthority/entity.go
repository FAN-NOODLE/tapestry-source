package microauthority

// WxQrcodeResponse 生成小程序太阳码返回格式
type WxQrcodeResponse struct {
	ResponseBase
	Data struct {
		Url      string `json:"url"`
		ShareImg string `json:"share_img"`
	}
}

//短信发送成功返回值
type SendMsgResponse struct {
	ResponseBase
}
