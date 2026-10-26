package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type SourceCategoryApi struct{}

func NewSourceCategoryApi() SourceCategoryApi {
	return SourceCategoryApi{}
}

func (a SourceCategoryApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/category/list"), request)
}

func (a SourceCategoryApi) Create(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/category/create"), request)
}

func (a SourceCategoryApi) Update(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/category/update"), request)
}

func (a SourceCategoryApi) SourceRelation(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/category/source_relation"), request)
}

func (a SourceCategoryApi) Delete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/category/delete"), request)
}
