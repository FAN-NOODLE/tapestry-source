package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/pkg/sdk"
)

type SourceCategoryApi struct{}

func NewSourceCategoryApi() SourceCategoryApi {
	return SourceCategoryApi{}
}

func (a SourceCategoryApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/category/list"), request)
}
