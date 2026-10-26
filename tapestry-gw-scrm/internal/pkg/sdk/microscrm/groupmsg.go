package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

// List 一个员工的朋友圈列表
func List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.GroupMsgListUrl), request)
}

// GroupSendAdd 新增群发任务
func GroupSendAdd(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.GroupSendAddUrl), request)
}

// GroupSendEdit 编辑群发任务
func GroupSendEdit(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.GroupSendEditUrl), request)
}

// TaskInfo 群发任务详情
func TaskInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TaskInfoUrl), request)
}

// TaskDetail 群发任务明细
func TaskDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TaskDetailUrl), request)
}

// TaskMemberList 任务明细-执行群主下拉列表
func TaskMemberList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TaskMemberListUrl), request)
}

func TaskMemberDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TaskMemberDetailUrl), request)
}

func TaskMemberSendList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TaskMemberSendListUrl), request)
}
