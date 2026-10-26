package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

// ChatMsgTypeList 消息类型列表
func ChatMsgTypeList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.ChatMsgTypeList), request)
}

// ChatMsgList 消息列表
func ChatMsgList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.ChatMsgList), request)
}

// ChatGroupMsgList 消息列表
func ChatGroupMsgList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.ChatGroupMsgList), request)
}

// ChatMsgListByID 根据消息ID找附近的消息
func ChatMsgListByID(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.ChatMsgListByID), request)
}

// ChatAgreeInfo 会话同意信息
func ChatAgreeInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.ChatAgreeInfo), request)
}

// ChatMemberList 员工列表
func ChatMemberList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.ChatMemberList), request)
}

// ChatObjectToMemberList 员工的员工列表
func ChatObjectToMemberList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.ChatObjectToMemberList), request)
}

// ChatObjectToExternalList 员工的外部联系人列表
func ChatObjectToExternalList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.ChatObjectToExternalList), request)
}

// ChatObjectToGroupList 员工的社群列表
func ChatObjectToGroupList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.ChatObjectToGroupList), request)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.GetUserInfoUrl), request)
}
