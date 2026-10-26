package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type CoinsApi struct{}

func NewCoinsApi() CoinsApi {
	return CoinsApi{}
}

func (a CoinsApi) RecordList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/coins/record-list"), request)
}

func (a CoinsApi) CoinsSettingsInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/coins/coins-settings-info"), request)
}

func (a CoinsApi) SaveCoinsSettings(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/coins/save-coins-settings"), request)
}
