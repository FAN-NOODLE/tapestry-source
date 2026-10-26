package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// TokenNotFind 请求未携带token
	TokenNotFind = errorx.NewErrorMsg(-1, "[jwt]请求未携带token，无权限访问")
	// TokenIsExp 授权已过期
	TokenIsExp = errorx.NewErrorMsg(-2, "[jwt]授权已过期")
	// UserIdNotFind 未查询到人员有效信息
	UserIdNotFind = errorx.NewErrorMsg(-3, "[jwt]未查询到人员有效信息")
	// AuthInfoNotFind 未获取到企业微信授权信息
	AuthInfoNotFind = errorx.NewErrorMsg(-4, "[jwt]未获取到企业微信授权信息")
	// AuthTokenIllegal 非法的token
	AuthTokenIllegal = errorx.NewErrorMsg(-5, "[jwt]非法的token")
)
