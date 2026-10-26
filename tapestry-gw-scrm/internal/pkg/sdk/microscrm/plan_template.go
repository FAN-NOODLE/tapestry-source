package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type PlanTemplateApi struct{}

func NewPlanTemplateApi() PlanTemplateApi {
	return PlanTemplateApi{}
}

func (a PlanTemplateApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/plan/template/lists"), request)
}

func (a PlanTemplateApi) Create(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/plan/template/create"), request)
}

func (a PlanTemplateApi) Detail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/plan/template/details"), request)
}

func (a PlanTemplateApi) Update(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/plan/template/update"), request)
}

func (a PlanTemplateApi) Status(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/plan/template/status"), request)
}
