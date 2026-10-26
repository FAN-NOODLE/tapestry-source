package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type ReserveApi struct{}

func NewReserveApi() ReserveApi {
	return ReserveApi{}
}

func (a ReserveApi) ServiceCreate(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/reserve/service/create"), request)
}

func (a ReserveApi) ServiceUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/reserve/service/update"), request)
}

func (a ReserveApi) ServiceDelete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/reserve/service/delete"), request)
}

func (a ReserveApi) ServiceList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/reserve/service/list"), request)
}

func (a ReserveApi) ScoreList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/reserve/service-score/list"), request)
}

func (a ReserveApi) OrderList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPostNoLog(c, GetUrl("%s/micro-scrm/v1/reserve/order/lists-ck"), request)
}

func (a ReserveApi) OrderDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/reserve/order/detail"), request)
}

func (a ReserveApi) ReserveOrderList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPostNoLog(c, GetUrl("%s/micro-scrm/v1/reserve/order/reserve-order-list"), request)
}

func (a ReserveApi) ReserveServiceCheckDelete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/reserve/service/check-delete"), request)
}

func (a ReserveApi) ReserveServiceServiceInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/reserve/service/info"), request)
}

func (a ReserveApi) OrderExportList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPostNoLog(c, GetUrl("%s/micro-scrm/v1/reserve/order/export-lists"), request)
}
