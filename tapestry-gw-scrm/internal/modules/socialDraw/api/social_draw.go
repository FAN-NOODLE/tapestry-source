package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type SocialDrawApi struct {
}

func (api *SocialDrawApi) GetWorkGroupChatList(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.WorkGroupChatList(c, req)
}

func (api *SocialDrawApi) GetWorkGroupChatInfo(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.WorkGroupChatInfo(c, req)
}

func (api *SocialDrawApi) GetWorkGroupChatTypeList(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.WorkGroupChatTypeList(c, req)
}

func (api *SocialDrawApi) ManageWorkGroupChatType(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.WorkGroupChatTypeManage(c, req)
}

func (api *SocialDrawApi) SetWorkGroupChat(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.WorkGroupChatSetting(c, req)
}

func (api *SocialDrawApi) GetWorkGroupChatFilters(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.WorkGroupChatFilters(c, req)
}

func (api *SocialDrawApi) GetWorkGroupChatMembers(c ctx.Context, req interface{}) (interface{}, error) {
	return microscrm.WorkGroupChatMembers(c, req)
}
