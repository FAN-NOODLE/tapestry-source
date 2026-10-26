package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type SensitiveApi struct {
}

// SensitiveList 敏感词列表
func (s SensitiveApi) SensitiveList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.SensitiveList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SensitiveSave 敏感词保存
func (s *SensitiveApi) SensitiveSave(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.SensitiveSave(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SensitiveDelete 敏感词删除
func (s *SensitiveApi) SensitiveDelete(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.SensitiveDelete(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SensitiveDetail 敏感词详情
func (s *SensitiveApi) SensitiveDetail(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.SensitiveDetail(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SensitiveDetailList 敏感词明细列表
func (s *SensitiveApi) SensitiveDetailList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.SensitiveDetailList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
