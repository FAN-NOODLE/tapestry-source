package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// 外部联系人错误
	SpuBaseDataErr = errorx.NewErrorMsg(1800001, "spu基本信息错误")
)
