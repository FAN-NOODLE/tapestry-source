package display

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

type Dashboard struct {
	microScrm.BaseSDK
}

func NewDashboard() *Dashboard {
	return &Dashboard{}
}

func (a *Dashboard) Tree(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/dashboard/area/tree"), request)

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

func (a *Dashboard) SearchStore(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/dashboard/area/search-store"), request)

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

func (a *Dashboard) MemberAuth(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/dashboard/area/manager-auth"), request)

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
