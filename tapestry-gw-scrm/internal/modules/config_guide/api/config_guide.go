package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk/gomicrouser"
)

type ConfigGuideApi struct {
}

// 配置引导
func (s *ConfigGuideApi) ConfigList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.ConfigGuideList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
