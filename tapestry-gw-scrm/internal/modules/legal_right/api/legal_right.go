package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk/gomicrouser"
)

type LegalRightApi struct {
}

// Save 权益保存
func (s *LegalRightApi) Save(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.Save(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// List 权益列表
func (s *LegalRightApi) List(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.List(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateRightStatus 修改权益状态
func (s *LegalRightApi) UpdateRightStatus(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.UpdateRightStatus(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
