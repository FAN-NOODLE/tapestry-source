package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// UserGroupIDIsEmpty 客户分组ID不能为空
	UserGroupIDIsEmpty = errorx.NewErrorMsg(1400001, "客户分组ID不能为空")
)
