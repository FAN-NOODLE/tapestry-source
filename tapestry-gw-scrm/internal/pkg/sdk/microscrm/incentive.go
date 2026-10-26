package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type IncentiveApi struct{}

func NewIncentiveApi() IncentiveApi {
	return IncentiveApi{}
}

func (a IncentiveApi) GetConfig(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/incentive/config"), request)
}

func (a IncentiveApi) RelationSave(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/incentive/relation/save"), request)
}

func (a IncentiveApi) RelationDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/incentive/relation/detail"), request)
}

func (a IncentiveApi) RelationDetailList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/incentive/relation/detail-list"), request)
}
