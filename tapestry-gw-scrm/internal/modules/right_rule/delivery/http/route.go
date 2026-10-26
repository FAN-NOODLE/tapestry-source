package right_rule_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/right-rule/get-score-name", "获取企业积分名称",
			endpoint.Handler(GetScoreName),
			endpoint.Tags("获取企业积分名称")),

		endpoint.Post("/right-rule/save-score-name", "设置企业积分名称",
			endpoint.Handler(SaveScoreName),
			endpoint.Tags("设置企业积分名称")),

		endpoint.Post("/right-rule/score-rule-save", "积分权益规则新增保存",
			endpoint.Handler(scoreRuleSave),
			endpoint.Tags("积分权益规则新增保存")),

		endpoint.Post("/right-rule/score-detail", "积分权益规则详情",
			endpoint.Handler(scoreDetail),
			endpoint.Tags("积分权益规则详情")),

		endpoint.Post("/right-rule/set-use", "权益规则-启用/禁用",
			endpoint.Handler(setUse),
			endpoint.Tags("权益规则-启用/禁用")),

		endpoint.Post("/right-rule/delete", "权益规则-删除",
			endpoint.Handler(ruleDelete),
			endpoint.Tags("权益规则-删除")),

		endpoint.Post("/right-rule/save-coupon-rule", "添加优惠券权限规则",
			endpoint.Handler(saveCouponRule),
			endpoint.Tags("添加优惠券权限规则")),

		endpoint.Post("/right-rule/wx-coupon-list", "优惠券列表",
			endpoint.Handler(wxCouponList),
			endpoint.Tags("优惠券列表")),

		endpoint.Post("/right-rule/info-coupon-rule", "优惠券权益规则详情",
			endpoint.Handler(infoCouponRule),
			endpoint.Tags("优惠券权益规则详情")),

		endpoint.Post("/right-rule/coupon-package-list", "券包列表",
			endpoint.Handler(couponPackageList),
			endpoint.Tags("优惠券权益规则详情")),

		endpoint.Post("/right-rule/list", "权益规则列表",
			endpoint.Handler(list),
			endpoint.Tags("权益规则列表")),
	}
}
