package couponHttp

import (
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-scrm/internal/modules/coupon"
	"micro-scrm/internal/modules/coupon/respository/mysql"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/goods"
	goodsMysql "micro-scrm/internal/modules/goods/respository/mysql"
	"micro-scrm/internal/modules/user"
	userMysql "micro-scrm/internal/modules/user/respository/mysql"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {
	userService = user.NewUserService(userMysql.NewUser(), userMysql.NewEvent(), userEs.NewEsRepo(), enterpriseMysql.New())
	couponService = coupon.New(mysql.NewCoupon())
	goodService = goods.NewGoods(goodsMysql.New())
	return []*swagger.Endpoint{
		endpoint.Post("/coupon/customer", "用户优惠券列表",
			endpoint.Handler(customerCoupon),
			endpoint.Tags("优惠券"),
			endpoint.Response(http.StatusOK, "", "coupon", "接口调用成功")),

		endpoint.Post("/coupon/wx-coupon-list", "微信优惠券列表",
			endpoint.Handler(wxCouponList),
			endpoint.Tags("微信优惠券列表"),
			endpoint.Response(http.StatusOK, "", "coupon", "接口调用成功")),

		endpoint.Post(
			"/coupon/goods-list",
			"部分可用商品列表",
			endpoint.Handler(couponGoods),
			endpoint.Response(http.StatusOK, "", "coupon", "接口调用成功"),
		),
		endpoint.Post("/coupon/wx-coupon-available-list", "可用的微信优惠券列表",
			endpoint.Handler(wxCouponAvailableList),
			endpoint.Tags("可用的微信优惠券列表"),
			endpoint.Response(http.StatusOK, "", "coupon", "接口调用成功")),
	}
}
