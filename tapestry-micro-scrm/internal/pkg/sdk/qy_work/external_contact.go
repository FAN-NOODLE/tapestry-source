package qyWork

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

const (
	// 发送欢迎语
	sendWelcome = "cgi-bin/externalcontact/send_welcome_msg"
	// 添加联系我二维码
	addExternalContact = "cgi-bin/externalcontact/add_contact_way"
	// 更新联系我二维码
	updateExternalContact = "cgi-bin/externalcontact/update_contact_way"
	//删除联系我二维码
	deleteExternalContact = "cgi-bin/externalcontact/del_contact_way"
	//查询联系我二维码
	getExternalContact = "cgi-bin/externalcontact/get_contact_way"
	//查询联系我二维码列表
	listExternalContact = "cgi-bin/externalcontact/list_contact_way"
	// 获取企业群发消息发送结果
	getGroupMsgResult = "cgi-bin/externalcontact/get_group_msg_result"
	// 添加企业群发消息任务
	addMsgTemplate = "cgi-bin/externalcontact/add_msg_template"
	// 编辑客户企业标签
	markTag = "cgi-bin/externalcontact/mark_tag"
)

// SendWelcome 发送欢迎语
func (w *QyWork) SendWelcome(c ctx.Context, req SendWelcomeReq) (QyResponse, error) {
	var resp QyResponse

	// 获取Token
	token, err := w.getToken(c, ManuCustomer)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(sendWelcome, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// AddContactWay 添加联系我二维码
func (w *QyWork) AddContactWay(c ctx.Context, req ContactCreateRequest) (ContactCreateResp, error) {
	var resp ContactCreateResp

	// 获取Token
	token, err := w.getToken(c, ManuCustomer)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(addExternalContact, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateContactWay  更新联系我二维码
func (w *QyWork) UpdateContactWay(c ctx.Context, req ContactUpdateRequest) (QyResponse, error) {
	var resp QyResponse

	// 获取Token
	token, err := w.getToken(c, ManuCustomer)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(updateExternalContact, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (w *QyWork) GetContactWay(c ctx.Context, req ContactDeleteRequest) (ContactGetResp, error) {
	var resp ContactGetResp

	// 获取Token
	token, err := w.getToken(c, ManuCustomer)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(getExternalContact, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (w *QyWork) GetContactWayList(c ctx.Context, req ContactListRequest) (ContactListResp, error) {
	var resp ContactListResp

	// 获取Token
	token, err := w.getToken(c, ManuCustomer)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(listExternalContact, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DeleteContactWay   删除联系我二维码
func (w *QyWork) DeleteContactWay(c ctx.Context, req ContactDeleteRequest) (QyResponse, error) {
	var resp QyResponse

	// 获取Token
	token, err := w.getToken(c, ManuCustomer)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(deleteExternalContact, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetGroupMsgResult 群发发送结果
func (w *QyWork) GetGroupMsgResult(c ctx.Context, msgId string) ([]GetGroupMsgResultResp, error) {
	var resp []GetGroupMsgResultResp

	// 获取Token
	token, err := w.getToken(c, ManuCustomer)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(getGroupMsgResult, token), map[string]any{"msgid": msgId}, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// AddMsgTemplate 添加企业群发消息
func (w *QyWork) AddMsgTemplate(c ctx.Context, req AddMsgTemplateReq) (AddMsgTemplateResp, error) {
	// 获取Token
	var resp AddMsgTemplateResp
	// 获取Token
	token, err := w.getToken(c, ManuCustomer)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(addMsgTemplate, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// MarkTag 编辑客户企业标签
func (w *QyWork) MarkTag(c ctx.Context, req MarkTagReq) (QyResponse, error) {
	var resp QyResponse
	// 获取Token
	token, err := w.getToken(c, ManuCustomer)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(markTag, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
