package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type SourceImageApi struct{}

func NewSourceImageApi() SourceImageApi {
	return SourceImageApi{}
}

func (a SourceImageApi) Update(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/image/update"), request)
}

func (a SourceImageApi) Detail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/image/detail"), request)
}
