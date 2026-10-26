package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk/gomicrouser"
)

type MemberLevelApi struct {
}

// 会员等级列表
func (s *MemberLevelApi) LevelList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.MemberLevelList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 会员等级详情
func (s *MemberLevelApi) LevelDetail(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.MemberLevelDetail(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 会员等级新增
func (s *MemberLevelApi) LevelSet(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.MemberLevelSet(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 会员等级状态修改
func (s *MemberLevelApi) UpdateStatus(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.UpdateStatus(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 等级内容排序
func (s *MemberLevelApi) ContentSort(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.ContentSort(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 下拉框会员等级列表
func (s *MemberLevelApi) SelectLevelList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.SelectLevelList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
