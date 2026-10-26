package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type OnlineWarehouseApi struct {
}

func NewOnlineWarehouseApi() OnlineWarehouseApi {
	return OnlineWarehouseApi{}
}

// Submit 网仓订单提交发货
func (a OnlineWarehouseApi) Submit(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/online_warehouse/submit"), request)
}

// List 网仓订单列表
func (a OnlineWarehouseApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/online_warehouse/list"), request)
}

// Info 网仓订单详情
func (a OnlineWarehouseApi) Info(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/online_warehouse/info"), request)
}

// Edit 订单编辑
func (a OnlineWarehouseApi) Edit(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/online_warehouse/edit"), request)
}

// StoreList 门店列表
func (a OnlineWarehouseApi) StoreList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/online_warehouse/store_list"), request)
}

func (a OnlineWarehouseApi) StoreOrderList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/online_warehouse/store_order_list"), request)
}

func (a OnlineWarehouseApi) ImportOrder(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/online_warehouse/import_order"), request)
}

func (a OnlineWarehouseApi) BatchTransfer(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/online_warehouse/batch_transfer"), request)
}

func (a OnlineWarehouseApi) Refund(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/online_warehouse/refund"), request)
}

func (a OnlineWarehouseApi) DepList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/online_warehouse/dep-list"), request)
}

func (a OnlineWarehouseApi) TransferExport(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/online_warehouse/transfer-export"), request)
}

func (a OnlineWarehouseApi) TransferImport(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/online_warehouse/transfer-import"), request)
}
