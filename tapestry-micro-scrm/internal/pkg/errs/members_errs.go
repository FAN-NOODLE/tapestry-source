package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// MemberIDIsEmpty 成员id不能为空错误
	MemberIDIsEmpty = errorx.NewErrorMsg(120001, "成员ID不能为空")
	// EnterpriseIDIsEmpty 企业id不能为空错误
	EnterpriseIDIsEmpty = errorx.NewErrorMsg(120002, "企业id不能为空")
	// MemberNotFund 成员不存在
	MemberNotFund = errorx.NewErrorMsg(120003, "成员不存在")
	// SyncContactRunning 人员同步正在进行
	SyncContactRunning = errorx.NewErrorMsg(12004, "同步通讯录任务正在进行")
)
