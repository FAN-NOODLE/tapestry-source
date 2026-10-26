package microScrmCustomer

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

type CustomerApi struct {
	microScrm.BaseSDK
}

func NewCustomerApi() CustomerApi {
	return CustomerApi{}
}

func (a CustomerApi) CrmUserDetail(c ctx.Context, request interface{}) (interface{}, error) {

	var resp microScrm.CrmUserDetail
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/customer/external-user/detail"), request)
	if err != nil {
		return resp, err
	}
	result := &microScrm.CrmUserDetailRes{}
	_, err = res.Read(clients.JSON(result))

	if result.Code != 200 {
		return resp, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a CustomerApi) GetUserDetailTipsNum(c ctx.Context, request interface{}) (resp microScrm.UserDetailTipsNum, err error) {

	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/customer/user-detail-tips-num"), request)
	if err != nil {
		return resp, err
	}
	result := &microScrm.UserDetailTipsNumRes{}
	_, err = res.Read(clients.JSON(result))

	if result.Code != 200 {
		return resp, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a CustomerApi) GetUserBindRelation(c ctx.Context, request interface{}) (resp []microScrm.UserSaListData, err error) {

	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/customer/user-bind-relation"), request)
	if err != nil {
		return resp, err
	}
	result := &microScrm.UserSaListDataRes{}
	_, err = res.Read(clients.JSON(result))

	if result.Code != 200 {
		return resp, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a CustomerApi) FirstOpen(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/member/first-open"), request)
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
