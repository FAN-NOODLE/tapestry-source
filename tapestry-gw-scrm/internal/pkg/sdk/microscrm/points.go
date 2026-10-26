package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type PointsApi struct{}

func NewPointsApi() PointsApi {
	return PointsApi{}
}

func (a PointsApi) RecordList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/excitation/points/record-list"), request)
}
