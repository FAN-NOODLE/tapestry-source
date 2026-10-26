package qyWork

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

const (
	// 发送应用消息
	messageSend = "cgi-bin/message/send"
	// 撤回应用消息
	messageRecall = "cgi-bin/message/recall"
)

// SendMiniMessage 发送小程序应用消息
func (w *QyWork) SendMiniMessage(c ctx.Context, req SendMiniMessageReq) (SendMessageResp, error) {
	var resp SendMessageResp
	// 获取Token
	token, err := w.getToken(c, StoreAssistantSelf)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(messageSend, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SendTemplateCardMessage 发送模板卡片消息
func (w *QyWork) SendTemplateCardMessage(c ctx.Context, identifier string, req SendTemplateCardMessageReq) (SendMessageResp, error) {
	var resp SendMessageResp
	// 获取Token
	token, err := w.getToken(c, identifier)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(messageSend, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RecallMessage 撤回应用消息
func (w *QyWork) RecallMessage(c ctx.Context, msgId string) (QyResponse, error) {
	var resp QyResponse
	// 获取Token
	token, err := w.getToken(c, StoreAssistantSelf)
	if err != nil {

	}

	err = w.requestPostJson(c, w.getTokenUrl(messageRecall, token), map[string]string{"msgid": msgId}, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
