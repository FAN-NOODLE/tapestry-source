package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type IntegralRuleApi struct{}

func NewIntegralRuleApi() IntegralRuleApi {
	return IntegralRuleApi{}
}

func (a IntegralRuleApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/integral-rule/list"), request)
}

func (a IntegralRuleApi) ParameterConfiguration(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/integral-rule/configuration"), request)
}

func (a IntegralRuleApi) Create(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/integral-rule/create"), request)
}

func (a IntegralRuleApi) Detail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/integral-rule/detail"), request)
}

func (a IntegralRuleApi) Destroy(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/integral-rule/destroy"), request)
}

func (a IntegralRuleApi) Update(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/integral-rule/update"), request)
}

func (a IntegralRuleApi) DescriptionSave(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/integral-rule/description-save"), request)
}

func (a IntegralRuleApi) DescriptionDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/integral-rule/description-detail"), request)
}
