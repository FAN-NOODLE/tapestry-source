package index_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/index"
	"gw-scrm/internal/pkg/sdk/microboss"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/index/qw/friend", "首页企微好友群成员",
			endpoint.Handler(indexData),
			endpoint.Tags("首页企微好友群成员"),
			endpoint.Body(index.PrivateDomain{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: index.PrivateDomain{},
			}, "indexData", "接口调用成功")),

		endpoint.Post("/index/shoppers/list", "导购业绩排行榜",
			endpoint.Handler(ShopperList),
			endpoint.Tags("导购业绩排行榜"),
			endpoint.Body(index.ShopperListRes{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: index.ShopperListRes{},
			}, "indexData", "接口调用成功")),

		endpoint.Post("/index/shop/fee", "首页商城变现",
			endpoint.Handler(ShopFee),
			endpoint.Tags("首页商城变现"),
			endpoint.Body(index.ShopFeeRes{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: index.ShopFeeRes{},
			}, "indexData", "接口调用成功")),

		endpoint.Post("/index/shop/conversion", "首页客户总转化变现",
			endpoint.Handler(TotalUv),
			endpoint.Tags("首页客户总转化变现"),
			endpoint.Body(index.ShopConversion{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: index.ShopConversion{},
			}, "indexData", "接口调用成功")),

		endpoint.Post("/index/help", "首页公告帮助中心",
			endpoint.Handler(NoticeHelp),
			endpoint.Tags("首页公告"),
			endpoint.Body(index.NoticeHelpRes{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: index.NoticeHelpRes{},
			}, "indexData", "接口调用成功")),

		endpoint.Post("/index/store/list", "首页门店列表",
			endpoint.Handler(StoreList),
			endpoint.Tags("首页门店列表"),
			endpoint.Body(index.StoreListRes{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: index.StoreListRes{},
			}, "indexData", "接口调用成功")),

		endpoint.Post("/index/base/config", "首页基本配置",
			endpoint.Handler(BaseConfigCheck),
			endpoint.Tags("首页基本配置"),
			endpoint.Body(index.BaseInfoCheck{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: index.BaseInfoCheck{},
			}, "indexData", "接口调用成功")),

		endpoint.Post("/index/buy/notice", "帮助中心列表",
			endpoint.Handler(BuyNotice),
			endpoint.Tags("付费通知"),
			endpoint.Body(microboss.BuyNoticeRes{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: microboss.BuyNoticeRes{},
			}, "indexData", "接口调用成功")),
	}
}
