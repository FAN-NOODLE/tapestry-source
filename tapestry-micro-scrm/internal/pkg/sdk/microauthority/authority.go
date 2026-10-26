package microauthority

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// WxQrcode 生成太阳码
func WxQrcode(c ctx.Context, request interface{}) (*WxQrcodeResponse, error) {
	var response WxQrcodeResponse
	err := httpClient.Post(c, GetURL(urlEnum.WxQrcode), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SendMsg 发送短信
func SendMsg(c ctx.Context, request interface{}) (*SendMsgResponse, error) {
	var response SendMsgResponse
	err := httpClient.Post(c, GetURL(urlEnum.SendMsg), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
