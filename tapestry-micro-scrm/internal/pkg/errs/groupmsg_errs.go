package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	GroupMsgContentParseErr  = errorx.NewErrorMsg(2100001, "群发消息内容解析失败")
	GroupMsgContentEmptyErr  = errorx.NewErrorMsg(2100002, "群发消息内容为空")
	GroupMsgTaskNotExist     = errorx.NewErrorMsg(2100003, "群发任务为空")
	GroupMsgTaskQueueDataErr = errorx.NewErrorMsg(2100004, "任务队列数据错误")
)
