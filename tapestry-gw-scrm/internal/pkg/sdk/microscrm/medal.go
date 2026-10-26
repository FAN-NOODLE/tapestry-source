package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type MedalApi struct{}

func NewMedalApi() MedalApi {
	return MedalApi{}
}

func (a MedalApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/medal/list"), request)
}

func (a MedalApi) Create(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/medal/create"), request)
}
func (a MedalApi) Detail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/medal/detail"), request)
}
func (a MedalApi) Update(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/medal/update"), request)
}
func (a MedalApi) Remove(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/medal/remove"), request)
}
func (a MedalApi) Status(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/medal/status"), request)
}

type MedalRecordApi struct{}

func NewMedalRecordApi() MedalRecordApi {
	return MedalRecordApi{}
}

func (a MedalRecordApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/medal/record-list"), request)
}

func (a MedalRecordApi) Reissue(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/medal/record-reissue"), request)
}

func (a MedalRecordApi) Remove(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/medal/record-remove"), request)
}

type MedalRuleApi struct{}

func NewMedalRuleApi() MedalRuleApi {
	return MedalRuleApi{}
}

func (a MedalRuleApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/medal/rule-list"), request)
}
