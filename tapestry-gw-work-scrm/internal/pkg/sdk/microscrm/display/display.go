package display

import (
	"gw-work-scrm/internal/app/http/request/display"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

type Api struct {
	microScrm.BaseSDK
}

func NewApi() *Api {
	return &Api{}
}

func (a *Api) TaskInfo(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/task-info"), request)

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

func (a *Api) TaskList(c ctx.Context, request display.MemberTaskListReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/task-list"), request)

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

func (a *Api) SaveTask(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/save"), request)

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

func (a *Api) Look(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/look"), request)

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

func (a *Api) Auth(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/auth"), request)

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

func (a *Api) VmList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/task/vm-list"), request)

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

func (a *Api) StoreList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/store-list"), request)

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

func (a *Api) ListTips(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/list-tips"), request)

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

func (a *Api) CheckList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/check-list"), request)

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

func (a *Api) CheckStoreList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/check-store-list"), request)

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

func (a *Api) StoreTaskInfo(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/store-task-info"), request)

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

func (a *Api) VmCheck(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/vm-check"), request)

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

func (a *Api) VmSendNotice(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/send-notice"), request)

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

func (a *Api) SaveImage(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/save-image"), request)

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

func (a *Api) DelImage(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/del-image"), request)

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

func (a *Api) SaveComment(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/save-comment"), request)

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

func (a *Api) DelComment(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/del-comment"), request)

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

func (a *Api) AllList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/all-list"), request)

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

func (a *Api) Dashboard(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/vm-dashboard"), request)

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

func (a *Api) FailList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/fail-list"), request)

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

func (a *Api) DashboardList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/dashboard-list"), request)

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

func (a *Api) StoreTaskList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/store-task-list"), request)

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

func (a *Api) SpecialSave(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/display/api/special-save"), request)

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
