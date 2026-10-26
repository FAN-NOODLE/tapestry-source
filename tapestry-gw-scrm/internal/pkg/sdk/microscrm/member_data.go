package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type MemberDataApi struct{}

func NewMemberDataApi() MemberDataApi {
	return MemberDataApi{}
}

func (a MemberDataApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/member-data/list"), request)
}
