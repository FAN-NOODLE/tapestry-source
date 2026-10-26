package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// OpenAPIParamsErr OpenAPI参数错误。
	OpenAPIParamsErr = errorx.NewErrorMsg(400, "参数错误")
)
