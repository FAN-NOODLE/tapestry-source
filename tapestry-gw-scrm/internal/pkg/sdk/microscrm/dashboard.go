package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type DashboardApi struct{}

func NewDashboardApi() DashboardApi {
	return DashboardApi{}
}

func (a DashboardApi) AreaTree(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/area/tree"), request)
}

func (a DashboardApi) NotRelationStoreList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/area/not-relation-store-list"), request)
}

//func (a DashboardApi) StoreList(c ctx.Context, request interface{}) (interface{}, error) {
//	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/area/store-list"), request)
//}

func (a DashboardApi) AreaSave(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/area/save"), request)
}

func (a DashboardApi) AreaRemove(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/area/remove"), request)
}

func (a DashboardApi) ManagerAuthSave(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/manager/auth-save"), request)
}

func (a DashboardApi) ManagerAuthDelete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/manager/auth-delete"), request)
}

func (a DashboardApi) ManagerAuthDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/manager/auth-detail"), request)
}

func (a DashboardApi) ManagerAuthList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/manager/auth-list"), request)
}

func (a DashboardApi) ManagerAllMemberIds(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/manager/all-member-ids"), request)
}

func (a DashboardApi) AreaList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/area/list"), request)
}

func (a DashboardApi) HkStoreSync(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/hkmo-store/sync"), request)
}

func (a DashboardApi) HkStoreLastSyncTime(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/hkmo-store/last-sync-time"), request)
}

func (a DashboardApi) HkStoreList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/hkmo-store/list"), request)
}

func (a DashboardApi) HkStoreDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/hkmo-store/detail"), request)
}

func (a DashboardApi) HkStoreSave(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/hkmo-store/save"), request)
}

func (a DashboardApi) HkStoreClose(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/hkmo-store/close"), request)
}

func (a DashboardApi) HkManagerStoreList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/hkmo-manager/store-list"), request)
}

func (a DashboardApi) HkManagerAuthList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/hkmo-manager/auth-list"), request)
}

func (a DashboardApi) HkManagerAuthSave(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/hkmo-manager/auth-save"), request)
}

func (a DashboardApi) HkManagerAuthDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/hkmo-manager/auth-detail"), request)
}

func (a DashboardApi) HkManagerAuthDelete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/hkmo-manager/auth-delete"), request)
}

func (a DashboardApi) HkManagerAllMemberIds(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/hkmo-manager/all-member-ids"), request)
}

func (a DashboardApi) HkAreaList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/dashboard/hkmo-area/list"), request)
}
