package customerhttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gw-work-scrm/internal/config"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"net/http"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {
	customerSdk = microScrm.NewCustomer(config.CfgData.Service.MicroScrm.Domain)
	siderbarSdk = microScrm.NewSidebarGroupChat(config.CfgData.Service.MicroScrm.Domain)
	return []*swagger.Endpoint{
		//endpoint.Post(
		//	"/customer/footprint-list",
		//	"用户足迹列表-黄晓敏",
		//	endpoint.Handler(footPrint),
		//	endpoint.Tags("侧边栏-客户画像"),
		//	endpoint.Body(microScrm.CustomerRequest{}, "请求参数", true),
		//	endpoint.Response(
		//		http.StatusOK,
		//		"",
		//		"footprint-list",
		//		"接口调用成功",
		//	),
		//),

		//endpoint.Post(
		//	"/customer/detail",
		//	"客户详情-黄晓敏",
		//	endpoint.Handler(detail),
		//	endpoint.Tags("侧边栏-客户画像"),
		//	endpoint.Body(microScrm.CustomerRequest{}, "请求参数", true),
		//	endpoint.Response(
		//		http.StatusOK,
		//		"",
		//		"detail",
		//		"接口调用成功",
		//	),
		//),

		endpoint.Post(
			"/customer/dl-order-list",
			"客户详情-订单列表-段龙",
			endpoint.Handler(dlOrderList),
			endpoint.Tags("侧边栏-客户画像-订单列表"),
			endpoint.Body(microScrm.CustomerRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"dl-order-list",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/customer/dl-points-list",
			"客户详情-积分明细列表",
			endpoint.Handler(dlPointsDetailList),
			endpoint.Tags("侧边栏-客户画像-订单列表"),
			endpoint.Body(microScrm.CustomerRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"dl-order-list",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/customer/dl-coupons-list",
			"客户详情-卡券列表",
			endpoint.Handler(dlCouponsList),
			endpoint.Tags("侧边栏-客户画像-订单列表"),
			endpoint.Body(microScrm.CustomerRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"dl-order-list",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/customer/dl-order-goods",
			"客户详情-订单商品-段龙",
			endpoint.Handler(dlOrderGoods),
			endpoint.Tags("侧边栏-客户画像-订单列表-订单商品"),
			endpoint.Body(microScrm.CustomerRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"dl-order-list",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/customer/belong-to-members",
			"添加的员工列表-黄晓敏",
			endpoint.Handler(belongToMembers),
			endpoint.Tags("侧边栏-客户画像"),
			endpoint.Body(microScrm.CustomerRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"belong-to-members",
				"接口调用成功",
			),
		),

		//endpoint.Post(
		//	"/customer/group-list",
		//	"加入社群的列表-黄晓敏",
		//	endpoint.Handler(groupChatList),
		//	endpoint.Tags("侧边栏-客户画像"),
		//	endpoint.Body(microScrm.CustomerRequest{}, "请求参数", true),
		//	endpoint.Response(
		//		http.StatusOK,
		//		"",
		//		"group-list",
		//		"接口调用成功",
		//	),
		//),

		endpoint.Post(
			"/customer/remark-save",
			"客户备注信息保存-刘鹏伟",
			endpoint.Handler(RemarkSave),
			endpoint.Tags("侧边栏-客户画像"),
			endpoint.Body(microScrm.RemarkInfo{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"remark-save",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/customer/order-list",
			"客户订单列表-刘鹏伟",
			endpoint.Handler(orderList),
			endpoint.Tags("侧边栏-客户画像"),
			endpoint.Body(microScrm.RequestOrderList{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"order-list",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/customer/order-detail",
			"客户订单详情-刘鹏伟",
			endpoint.Handler(orderDetail),
			endpoint.Tags("侧边栏-客户画像"),
			endpoint.Body(microScrm.RequestOrderDetail{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"order-detail",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/customer/set-external-tags",
			"设置外部联系人企业标签-刘鹏伟",
			endpoint.Handler(SetExternalTags),
			endpoint.Tags("侧边栏-客户画像"),
			endpoint.Body(microScrm.SetExternalTags{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"remark-save",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/customer/coupon-list",
			"用户优惠券列表",
			endpoint.Handler(couponList),
			endpoint.Tags("侧边栏-客户画像"),
			endpoint.Body(microScrm.CustomerRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"coupon-list",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/customer/address-list",
			"用户地址列表",
			endpoint.Handler(addressList),
			endpoint.Tags("侧边栏-客户画像"),
			endpoint.Body(microScrm.CustomerRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"address-list",
				"接口调用成功",
			),
		),
	}
}
