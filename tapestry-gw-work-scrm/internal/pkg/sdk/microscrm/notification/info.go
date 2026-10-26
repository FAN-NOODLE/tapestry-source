package notification

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

type Info struct {
	microScrm.BaseSDK
}

func NewInfo() *Info {
	return &Info{}
}

func (a *Info) MergeList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/notification/merge/list-show"), request)

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

func (a *Info) InfoList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/notification/list"), request)

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

func (a *Info) Info(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/notification/info"), request)

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

func (a *Info) MergeDetail(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/notification/merge/detail"), request)

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

func (a *Info) Recall(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/notification/recall"), request)

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

func (a *Info) Cancel(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/notification/cancel"), request)

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

func (a *Info) Check(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/notification/check"), request)

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

func (a *Info) Manager(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/notification/manager"), request)

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

func (a *Info) SendFileByEmail(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/notification/send-file-by-email"), request)

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

func (a *Info) ListReadNum(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/notification/list-read-num"), request)

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

func (a *Info) SavePdf(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/notification/merge/save-pdf"), request)

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
