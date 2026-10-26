package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// TokenNotFind 请求未携带token
	TokenNotFind = errorx.NewErrorMsg(-1, "请求未携带token，无权限访问")
	// TokenIsExp 授权已过期
	TokenIsExp = errorx.NewErrorMsg(-2, "授权已过期")
	// UserIdNotFind 未查询到人员有效信息
	UserIdNotFind = errorx.NewErrorMsg(-3, "未查询到人员有效信息")
	// AuthTokenIllegal 非法的token
	AuthTokenIllegal = errorx.NewErrorMsg(-4, "非法的token")
	// PermissionDenied 没有该接口权限
	PermissionDenied = errorx.NewErrorMsg(-5, "没有该接口权限")
)
