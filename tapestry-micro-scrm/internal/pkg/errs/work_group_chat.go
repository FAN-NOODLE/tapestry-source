package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	GroupChatImportAgain    = errorx.NewErrorMsg(2100001, "请重新上传数据")
	GroupChatTypeNotDeleted = errorx.NewErrorMsg(2100002, "该群类型下存在分组，不可删除")
	GroupChatTypeExisted    = errorx.NewErrorMsg(2100003, "该名称已存在，请重新输入")
)
