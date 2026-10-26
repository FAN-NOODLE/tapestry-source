package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type StoreVisitApi struct {
}

func NewStoreVisitApi() StoreVisitApi {
	return StoreVisitApi{}
}

// List 巡店记录列表
func (a StoreVisitApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/store-visit/admin/list"), request)
}

func (a StoreVisitApi) Close(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/store-visit/admin/close"), request)
}

func (a StoreVisitApi) DeptNames(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/store-visit/dept-name-list"), request)
}

func (a StoreVisitApi) Dashboard(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/store-visit/admin/dashboard"), request)
}

func (a StoreVisitApi) Export(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/store-visit/admin/export"), request)
}

func (a StoreVisitApi) DashboardExport(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/store-visit/admin/dashboard/export"), request)
}
