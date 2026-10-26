package couponhttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gw-work-scrm/internal/config"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"net/http"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {
	couponSdk = microScrm.NewCoupon(config.CfgData.Service.MicroScrm.Domain)
	return []*swagger.Endpoint{
		endpoint.Post(
			"/wx-coupon/list",
			"微信优惠优惠券列表",
			endpoint.Handler(wxCouponList),
			endpoint.Tags("侧边栏社群详情"),
			endpoint.Body(microScrm.GroupChatInfoRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"groupChatInfo",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/wx-coupon/goods-list",
			"查询优惠券关联的商品",
			endpoint.Handler(couponGoods),
			endpoint.Tags("优惠券"),
			endpoint.Body(microScrm.GroupChatInfoRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"groupChatInfo",
				"接口调用成功",
			),
		),
	}
}
