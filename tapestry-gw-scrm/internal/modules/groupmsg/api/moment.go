package api

import (
	"encoding/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/modules/groupmsg"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type GroupMsgApi struct {
}

// List 群发任务列表
func (m *GroupMsgApi) List(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.List(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GroupSendAdd 新增群发消息
func (m *GroupMsgApi) GroupSendAdd(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.GroupSendAdd(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GroupSendEdit 编辑群发消息
func (m *GroupMsgApi) GroupSendEdit(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.GroupSendEdit(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// MemberList 任务明细-执行群主下拉列表
func (m *GroupMsgApi) MemberList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.TaskMemberList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TaskInfo 群发任务详情
func (m *GroupMsgApi) TaskInfo(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.TaskInfo(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TaskDetail 群发任务明细
func (m *GroupMsgApi) TaskDetail(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.TaskDetail(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *GroupMsgApi) TaskMemberDetail(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.TaskMemberDetail(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *GroupMsgApi) TaskMemberSendList(c ctx.Context, req interface{}) (resp groupmsg.TaskMemberSendListResp, err error) {
	res, err := microscrm.TaskMemberSendList(c, req)
	if err != nil {
		return resp, err
	}
	resByte, err := json.Marshal(res)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(resByte, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
