package wxcouponhttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	wxcoupon "gw-scrm/internal/modules/wxCoupon"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/coupon/wx-coupon-available-list", "可用的微信优惠券列表",
			endpoint.Handler(wxCouponAvailableList),
			endpoint.Tags("微信优惠券"),
			endpoint.Body(wxcoupon.AvailableWxCouponRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "WxCoupon", "接口调用成功")),
	}
}
