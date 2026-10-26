package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type LevelRuleApi struct{}

func NewLevelRuleApi() LevelRuleApi {
	return LevelRuleApi{}
}

func (a LevelRuleApi) Save(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/level-rule/save"), request)
}

func (a LevelRuleApi) GetDetailList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/level-rule/detail-list"), request)
}

func (a LevelRuleApi) GetList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/level-rule/list"), request)
}
