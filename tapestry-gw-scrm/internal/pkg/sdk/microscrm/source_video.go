package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type SourceVideoApi struct{}

func NewSourceVideoApi() SourceVideoApi {
	return SourceVideoApi{}
}

func (a SourceVideoApi) Update(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/video/update"), request)
}

func (a SourceVideoApi) Detail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/video/detail"), request)
}
