package microScrmReserve

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gw-work-scrm/internal/app/http/request/reserve"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

type OrderApi struct {
	microScrm.BaseSDK
}

func NewOrderApi() *OrderApi {
	return &OrderApi{}
}

// CreateOrder 创建邀约单
func (a OrderApi) CreateOrder(c ctx.Context, request reserve.OrderInfoReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/reserve/order/create"), request)

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

// OrderLists 邀约单列表
func (a OrderApi) OrderLists(c ctx.Context, request reserve.OrderListReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/reserve/order/lists"), request)

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

// TimeList 时间段的邀约单数
func (a OrderApi) TimeList(c ctx.Context, request reserve.OrderTimeListReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/reserve/order/time-list"), request)

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

// MemberServiceLists 导购可创建邀约的服务列表
func (a OrderApi) MemberServiceLists(c ctx.Context, request reserve.ReserveServiceListByMemberReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/reserve/service/member-list"), request)

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

func (a OrderApi) OrderShare(c ctx.Context, request reserve.OrderBaseReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/reserve/order/share"), request)

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

func (a OrderApi) MemberDetail(c ctx.Context, request reserve.MemberDetailReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/member/get-detail"), request)

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

func (a OrderApi) OrderDelete(c ctx.Context, request reserve.OrderBaseReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/reserve/order/delete"), request)

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

func (a OrderApi) CheckReserveTime(c ctx.Context, request reserve.CheckReserveTime) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/reserve/order/check-reserve-time"), request)

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

func (a OrderApi) VerifyShare(c ctx.Context, request reserve.OrderAssignReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/reserve/order/verify-share"), request)

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

func (a OrderApi) OrderChangStatus(c ctx.Context, request reserve.ChangeStatus) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/reserve/order/chang-status"), request)

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

func (a OrderApi) OrderShow(c ctx.Context, request reserve.OrderShowReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/reserve/order/show"), request)

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

// CustomerReserveLists 获取用户的所有邀约单数据
func (a OrderApi) CustomerReserveLists(c ctx.Context, request reserve.CustomerReserveListsReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/reserve/order/customer-reserve-lists"), request)

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
