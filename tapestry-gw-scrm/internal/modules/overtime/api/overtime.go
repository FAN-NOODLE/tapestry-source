package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type OvertimeApi struct {
}

// OvertimeList 超时记录列表
func (s *OvertimeApi) OvertimeList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.OvertimeList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
