package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/groupmsg"
)

var service *groupmsg.Service

// groupSendAdd 新增群发消息
func groupSendAdd(c *gin.Context) {
	var (
		request groupmsg.GroupMsgAddRequest
		err     error
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	err = service.GroupSendAdd(ctx.Wrap(c), &request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
	return
}

// groupSendEdit编辑群发消息
func groupSendEdit(c *gin.Context) {
	var (
		request groupmsg.GroupMsgEditRequest
		err     error
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	err = service.GroupSendEdit(ctx.Wrap(c), &request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
	return
}

// memberList任务明细-执行群主下拉列表
func memberList(c *gin.Context) {
	var (
		request groupmsg.MemberListRequest
		err     error
		members []*groupmsg.Member
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	members, err = service.MemberList(ctx.Wrap(c), &request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, members)
	return
}

// TaskInfo群发任务详情
func TaskInfo(c *gin.Context) {
	var (
		request groupmsg.TaskIfoRequest
		task    *groupmsg.TaskInfoResponse
		err     error
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	task, err = service.TaskInfo(ctx.Wrap(c), &request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, task)
	return
}

// msgTaskList 群发任务列表
func msgTaskList(c *gin.Context) {
	var request groupmsg.MsgTaskRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	res, err := service.MsgTaskList(ctx.Wrap(c), &request)

	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

// TaskDetail 群发任务明细
func TaskDetail(c *gin.Context) {
	var (
		request  groupmsg.TaskDetailRequest
		err      error
		response *groupmsg.TaskDetailResponse
	)
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	response, err = service.TaskDetail(ctx.Wrap(c), request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
	return
}

// TaskMemberDetail 群发员工任务明细
func TaskMemberDetail(c *gin.Context) {
	var (
		request  groupmsg.TaskDetailMemberRequest
		err      error
		response *groupmsg.TaskDetailMemberResponse
	)
	if err = c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	response, err = service.TaskDetailMember(ctx.Wrap(c), request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, response)
	return
}

// TaskMemberSendList 员工群发记录
func TaskMemberSendList(c *gin.Context) {
	var request groupmsg.TaskMemberSendListReq

	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	response, err := service.TaskMemberSendList(ctx.Wrap(c), request)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, &response)
	return
}
