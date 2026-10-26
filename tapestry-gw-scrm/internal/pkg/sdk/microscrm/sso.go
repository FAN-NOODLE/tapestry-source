package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type UserApi struct {
}

func NewUserApi() UserApi {
	return UserApi{}
}

// SSOLogin SSO登录
func (u UserApi) SSOLogin(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/sso/login"), request)
}
