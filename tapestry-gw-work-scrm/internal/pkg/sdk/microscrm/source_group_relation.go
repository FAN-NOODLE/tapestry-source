package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/pkg/sdk"
)

type SourceGroupRelationApi struct{}

func NewSourceGroupRelationApi() SourceGroupRelationApi {
	return SourceGroupRelationApi{}
}

func (a SourceGroupRelationApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/group-relation/list"), request)
}

func (a SourceGroupRelationApi) Save(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/group-relation/save"), request)
}

func (a SourceGroupRelationApi) GetListForStore(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/group-relation/get-list-for-store"), request)
}

func (a SourceGroupRelationApi) SelectList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/select-list"), request)
}

func (a SourceGroupRelationApi) GetNewSourceCount(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/get-new-num"), request)
}

func (a SourceGroupRelationApi) CollectSelectList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/collect-select-list"), request)
}
