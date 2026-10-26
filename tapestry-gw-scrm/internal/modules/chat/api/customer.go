package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type ChatApi struct{}

func (api *ChatApi) ChatMsgTypeList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ChatMsgTypeList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *ChatApi) ChatMsgList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ChatMsgList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *ChatApi) ChatGroupMsgList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ChatGroupMsgList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *ChatApi) ChatMsgListByID(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ChatMsgListByID(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *ChatApi) ChatAgreeInfo(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ChatAgreeInfo(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *ChatApi) ChatMemberList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ChatMemberList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *ChatApi) ChatObjectToExternalList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ChatObjectToExternalList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *ChatApi) ChatObjectToMemberList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ChatObjectToMemberList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *ChatApi) ChatObjectToGroupList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ChatObjectToGroupList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *ChatApi) GetUserInfo(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.GetUserInfo(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
