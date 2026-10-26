package microScrmCustomer

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	customerReq "gw-work-scrm/internal/app/http/request/customer"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

type CustomerInfoApi struct {
	microScrm.BaseSDK
}

func NewCustomerInfoApi() CustomerInfoApi {
	return CustomerInfoApi{}
}

func (a CustomerInfoApi) SaveInfoAdd(c ctx.Context, request interface{}) (interface{}, error) {
	return a.Post(c, a.GetUrl("%s/micro-scrm/v1/customer/save-info-add"), request)
}

func (a CustomerInfoApi) GetInfoAdd(c ctx.Context, request interface{}) (microScrm.UserInfoAdd, error) {
	var resp microScrm.UserInfoAdd
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/customer/get-info-add"), request)
	if err != nil {
		return resp, err
	}
	result := &microScrm.UserInfoAddRes{}
	_, err = res.Read(clients.JSON(result))

	if result.Code != 200 {
		return resp, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a CustomerInfoApi) SimpleInfo(c ctx.Context, request customerReq.SimpleInfoRequest) (interface{}, error) {
	var resp microScrm.UserInfoAdd
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/customer/user-simple-info"), request)
	if err != nil {
		return resp, err
	}
	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))

	if result.Code != 200 {
		return resp, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}
