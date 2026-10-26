package microScrmCustomer

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	customerReq "gw-work-scrm/internal/app/http/request/customer"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

type SowTaskApi struct {
	microScrm.BaseSDK
}

func NewSowTaskApi() *SowTaskApi {
	return &SowTaskApi{}
}

func (a SowTaskApi) Lists(c ctx.Context, request customerReq.ListReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/customer/sow-task/list"), request)

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

func (a SowTaskApi) Create(c ctx.Context, request customerReq.SowTaskCreate) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/customer/sow-task/create"), request)

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

func (a SowTaskApi) Update(c ctx.Context, request customerReq.SowTaskUpdate) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/customer/sow-task/update"), request)

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

func (a SowTaskApi) Finished(c ctx.Context, request customerReq.SowTaskFinished) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/customer/sow-task/finished"), request)

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

func (a SowTaskApi) Delete(c ctx.Context, request customerReq.SowTaskFinished) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/customer/sow-task/delete"), request)

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
