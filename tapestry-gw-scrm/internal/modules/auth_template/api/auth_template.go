package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk/microshop"
)

type AuthTemplateApi struct {
}

// 授权页配置详情
func (s *AuthTemplateApi) Detail(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microshop.AuthorizationDetail(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 授权页配置保存
func (s *AuthTemplateApi) Save(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microshop.AuthorizationSave(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
