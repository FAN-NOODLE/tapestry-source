package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// 海报生成错误
	GetPosterErr = errorx.NewErrorMsg(1900001, "海报生成错误")
)
