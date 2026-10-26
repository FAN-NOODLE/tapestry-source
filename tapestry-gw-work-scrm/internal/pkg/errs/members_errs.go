package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// MemberIDIsEmpty 成员id不能为空错误
	MemberIDIsEmpty = errorx.NewErrorMsg(120001, "成员ID不能为空")
	// EnterpriseIDIsEmpty 企业id不能为空错误
	EnterpriseIDIsEmpty = errorx.NewErrorMsg(120002, "企业id不能为空")
)
