package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk/gomicrouser"
)

type LevelDirectApi struct {
}

// 等级直达列表
func (s *LevelDirectApi) DirectList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.LevelDirectList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 等级直达添加
func (s *LevelDirectApi) DirectAdd(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.LevelDirectAdd(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 等级直达更新
func (s *LevelDirectApi) DirectUpdate(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.LevelDirectUpdate(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 等级直达删除
func (s *LevelDirectApi) DirectDelete(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.LevelDirectDelete(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 等级直达详情
func (s *LevelDirectApi) DirectDetail(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.LevelDirectDetail(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
