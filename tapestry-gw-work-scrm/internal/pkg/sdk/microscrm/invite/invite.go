package invite

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gw-work-scrm/internal/app/http/request/invite"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

type Api struct {
	microScrm.BaseSDK
}

func NewInviteApi() *Api {
	return &Api{}
}

// ConfigList 邀请函配置页列表
func (a Api) ConfigList(c ctx.Context, request invite.TemplateConfigListReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/invite/template/config-list"), request)

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

// ConfigDetail 邀请函模板配置详情
func (a Api) ConfigDetail(c ctx.Context, request invite.TemplateConfigDetailReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/invite/template/config-detail"), request)

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

// CreateRecord 发送邀约单记录
func (a Api) CreateRecord(c ctx.Context, request invite.RecordInfoReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/invite/create-record"), request)

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

// RecordInfo 用户查看邀请函信息
func (a Api) RecordInfo(c ctx.Context, request invite.GetRecordInfoReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/invite/get-record-info"), request)

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

func (a Api) RecordDelete(c ctx.Context, request invite.GetRecordInfoReq) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/invite/record-delete"), request)

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
