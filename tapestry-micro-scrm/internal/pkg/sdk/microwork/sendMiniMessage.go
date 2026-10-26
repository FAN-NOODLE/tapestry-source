package microwork

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// Deprecated: 请使用 messageSvc.NewMessageService().SendMiniMessage()
// SendMiniMessage 发送智慧导购应用消息
func SendMiniMessage(c ctx.Context, request *SendMiniMessageRequest) (*SendMiniMessageResponse, error) {
	var response SendMiniMessageResponse
	err := httpClient.Post(c, GetURL("%s/service/v1/mini-program/send-message"), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
