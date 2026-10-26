package gomicrouser

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// ScoreRuleSave 积分权益新增保存
func ScoreRuleSave(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ScoreRuleSave), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// ScoreDetail 积分权益规则详情
func ScoreDetail(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ScoreDetail), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// SetUse 权益规则-启用/禁用
func SetUse(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.SetUse), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// RuleDelete 权益规则-删除
func RuleDelete(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.RuleDelete), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// SaveCouponRule 新增编辑优惠券权益规则
func SaveCouponRule(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.SaveCouponRuleUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// wxCouponList 优惠券列表
func WxCouponList(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.WxCouponListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// InfoCouponRule 优惠券权益规则详情
func InfoCouponRule(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.InfoCouponRuleUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// CouponPackageList 券包列表
func CouponPackageList(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.CouponPackageListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// RuleList 权益规则列表
func RuleList(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.RuleListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
