package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type RemindconfigApi struct {
}

//  ConfigInfo 提醒设置获取
func (s *RemindconfigApi) ConfigInfo(c ctx.Context, req interface{}) interface{} {
	res, err := microscrm.ConfigInfo(c, req)
	if err != nil {
		return nil
	}
	return res
}

// ConfigSave 提醒设置保存
func (s *RemindconfigApi) ConfigSave(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.ConfigSave(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
