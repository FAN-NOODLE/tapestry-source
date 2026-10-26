package microwork

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

type SubscribeMessage struct {
}

func NewSubscribeMessage() *SubscribeMessage {
	return &SubscribeMessage{}
}

// Send 发送模板消息
func (s SubscribeMessage) Send(c ctx.Context, request SubscribeMessageSendReq) (SubscribeMessageSendResp, error) {
	var response SubscribeMessageSendResp
	err := httpClient.Post(c, GetURL("%s/service/v1/subscribe-message/send"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
