package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk/gomicrouser"
)

type RightRuleApi struct {
}

// ScoreRuleSave 积分权益新增保存
func (s *RightRuleApi) ScoreRuleSave(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.ScoreRuleSave(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ScoreDetail 积分权益规则详情
func (s *RightRuleApi) ScoreDetail(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.ScoreDetail(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SetUse 权益规则-启用/禁用
func (s *RightRuleApi) SetUse(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.SetUse(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// RuleDelete 权益规则-删除
func (s *RightRuleApi) RuleDelete(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.RuleDelete(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SaveCouponRule 新增编辑优惠券权益规则
func (s *RightRuleApi) SaveCouponRule(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.SaveCouponRule(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// wxCouponList 优惠券列表
func (s *RightRuleApi) WxCouponList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.WxCouponList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// InfoCouponRule 优惠券权益规则详情
func (s *RightRuleApi) InfoCouponRule(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.InfoCouponRule(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// InfoCouponRule 优惠券权益规则详情
func (s *RightRuleApi) CouponPackageList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.CouponPackageList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// RuleList 权益规则列表
func (s *RightRuleApi) RuleList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := gomicrouser.RuleList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
