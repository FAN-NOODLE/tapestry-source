package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type InviteApi struct{}

func NewInviteApi() InviteApi {
	return InviteApi{}
}

func (a InviteApi) TemplateSave(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/invite/template/save"), request)
}

func (a InviteApi) TemplateList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/invite/template/list"), request)
}

func (a InviteApi) TemplateInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/invite/template/info"), request)
}

func (a InviteApi) TemplateDelete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/invite/template/delete"), request)
}

func (a InviteApi) RecordOrderList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/invite/record-order-list"), request)
}

func (a InviteApi) UpdateStatus(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/invite/template/update-status"), request)
}
