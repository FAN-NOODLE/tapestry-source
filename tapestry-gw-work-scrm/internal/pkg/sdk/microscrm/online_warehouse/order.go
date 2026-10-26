package online_warehouse

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

type Order struct {
	microScrm.BaseSDK
}

func NewOrder() *Order {
	return &Order{}
}

func (a *Order) TipsNum(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/tips_num"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) Create(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/create"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) List(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/list"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) Info(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/info"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) DcList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/inventory/dc-list"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) InventoryList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/inventory/list"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) StoreRank(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/store-rank"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) SaRank(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/sa-rank"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) ExportMail(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/send-file-by-email"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err

}

func (a *Order) CustomerInfo(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/customer-order-info"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) CustomerEdit(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/customer-order-edit"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) LeaveSa(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/leave-sa"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err

}

func (a *Order) DepList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/dep-list"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) SaList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/sa-list"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) Cancel(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/cancel"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) PosterMedia(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/poster-media"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) CheckOrderSn(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/check-order"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) Dashboard(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/dashboard"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *Order) GoodsRank(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/online_warehouse/goods-rank"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}
