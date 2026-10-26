package marketinghttp

import (
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/customer"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	customerMysql "micro-scrm/internal/modules/customer/respository/mysql"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/marketing"
	"micro-scrm/internal/modules/marketing/respository/mysql"
	"micro-scrm/internal/modules/sms"
	smsMysql "micro-scrm/internal/modules/sms/respository/mysql"
	"micro-scrm/internal/modules/user"
	userMysql "micro-scrm/internal/modules/user/respository/mysql"
)

//Endpoints 路由管理
func Endpoints() []*swagger.Endpoint {
	service = marketing.New(mysql.New(), userEs.NewEsRepo())
	customerService = customer.New(customerMysql.NewCustomer(), customerMysql.NewWorkUser(), userEs.NewEsRepo())
	userService = user.NewUserService(userMysql.NewUser(), userMysql.NewEvent(), userEs.NewEsRepo(), enterpriseMysql.New())
	smsService = sms.New(smsMysql.New())
	service.SmsService = smsService
	enterpriseService = enterprise.New(enterpriseMysql.New())
	return []*swagger.Endpoint{
		endpoint.Post("/marketing-strategies/store", "添加营销策略",
			endpoint.Handler(marketingStrategy),
			endpoint.Tags("营销策略"),
			endpoint.Body(marketing.MarketingStrategy{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: marketing.MarketingAddResponse{},
			}, "ActivityPages", "接口调用成功")),
		endpoint.Post("/marketing-strategies/inner-store", "添加内置营销策略",
			endpoint.Handler(marketingStrategyInner),
			endpoint.Tags("营销策略"),
			endpoint.Body(marketing.MarketingStrategy{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: marketing.MarketingAddResponse{},
			}, "ActivityPages", "接口调用成功")),
		endpoint.Post("/marketing-strategies/list", "营销策略列表",
			endpoint.Handler(marketingStrategies),
			endpoint.Tags("营销策略"),
			endpoint.Body(marketing.ListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: marketing.MarketingStrategyListResponse{},
			}, "ActivityPages", "接口调用成功")),
		endpoint.Post("/marketing-strategies/update", "更新营销策略",
			endpoint.Handler(marketingStrategyUpdate),
			endpoint.Tags("营销策略"),
			endpoint.Body(marketing.MarketingStrategy{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: marketing.MarketingStrategy{},
			}, "ActivityPages", "接口调用成功")),
		endpoint.Post("/marketing-strategies/detail", "营销策略详情",
			endpoint.Handler(marketingStrategyDetail),
			endpoint.Tags("营销策略"),
			endpoint.Body(marketing.MarketingDetailRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: marketing.MarketingStrategy{},
			}, "ActivityPages", "接口调用成功")),

		endpoint.Post("/marketing-strategies/data/detail", "营销策略数据详情",
			endpoint.Handler(marketingStrategyDataDetail),
			endpoint.Tags("营销策略"),
			endpoint.Body(marketing.MarketingDataDetailRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: marketing.MarketingStrategy{},
			}, "ActivityPages", "接口调用成功")),
		endpoint.Post("/marketing-strategies/delete", "删除营销策略",
			endpoint.Handler(marketingStrategyDelete),
			endpoint.Tags("营销策略"),
			endpoint.Body(marketing.DeleteRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "ActivityPages", "接口调用成功")),
		endpoint.Post("/marketing-strategies/update-status", "营销策略状态修改",
			endpoint.Handler(marketingStrategyUpdateStatus),
			endpoint.Tags("营销策略"),
			endpoint.Body(marketing.DeleteRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "ActivityPages", "接口调用成功")),

		endpoint.Post("/marketing-strategies/selector-text", "选择项",
			endpoint.Handler(SelectorText),
			endpoint.Tags("营销策略"),
			endpoint.Body(marketing.DeleteRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "ActivityPages", "接口调用成功")),

		endpoint.Post("/marketing-strategies/market-trig-exec", "选择项",
			endpoint.Handler(MarketTrigExec),
			endpoint.Tags("营销触达"),
			endpoint.Body(marketing.DeleteRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "ActivityPages", "接口调用成功")),

		endpoint.Post("/marketing-strategies/update-strategy-time", "设置策略时间",
			endpoint.Handler(UpdateStrategyTime),
			endpoint.Tags("营销触达"),
			endpoint.Body(marketing.DeleteRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "ActivityPages", "接口调用成功")),

		endpoint.Post("/marketing-strategies/tag-group-list", "客户标签列表",
			endpoint.Handler(tagGroupList),
			endpoint.Tags("客户标签列表"),
			endpoint.Response(http.StatusOK, "", "tagGroupList", "接口调用成功")),

		endpoint.Post("/marketing-strategies/tag-group-user-count", "客户标签列表",
			endpoint.Handler(tagGroupUserCount),
			endpoint.Tags("客户标签列表"),
			endpoint.Response(http.StatusOK, "", "tagGroupList", "接口调用成功")),

		endpoint.Post("/marketing-strategies/test-group-send", "群发结果获取测试接口",
			endpoint.Handler(getGroupMsgResult),
			endpoint.Tags("客户标签列表"),
			endpoint.Response(http.StatusOK, "", "getGroupMsgResult", "接口调用成功")),

		endpoint.Post("/marketing-strategies/trigger-update", "批量更新营销策略的触达人数",
			endpoint.Handler(MarketingTriggerUpdate),
			endpoint.Tags("批量更新营销策略的触达人数"),
			endpoint.Response(http.StatusOK, "", "MarketingTriggerUpdate", "批量更新营销策略的触达人数成功")),

		endpoint.Post("/marketing-strategies/get-wx-channel", "获取场景值列表",
			endpoint.Handler(wxChannelList),
			endpoint.Tags("场景值列表"),
			endpoint.Response(http.StatusOK, "", "wxChannelList", "接口调用成功")),
		endpoint.Post("/marketing-strategies/get-spu-detail", "获取商品详情",
			endpoint.Handler(GetSpuDetail),
			endpoint.Tags("获取商品详情"),
			endpoint.Response(http.StatusOK, "", "GetSpuDetail", "获取商品详情")),
		endpoint.Post("/marketing-strategies/short-url-add", "添加短链",
			endpoint.Handler(ShortUrlAdd),
			endpoint.Tags("添加短链"),
			endpoint.Response(http.StatusOK, "", "ShortUrlAdd", "添加短链")),
		endpoint.Post("/marketing-strategies/short-url-data", "短链数据",
			endpoint.Handler(ShortUrlData),
			endpoint.Tags("短链数据"),
			endpoint.Response(http.StatusOK, "", "ShortUrlData", "短链数据")),
		endpoint.Post("/marketing-strategies/short-url-list", "获取短链列表",
			endpoint.Handler(ShortUrlList),
			endpoint.Tags("获取短链列表"),
			endpoint.Response(http.StatusOK, "", "ShortUrlList", "获取短链列表")),
		endpoint.Post("/marketing-strategies/short-url-del", "删除短链",
			endpoint.Handler(ShortUrlDel),
			endpoint.Tags("删除短链"),
			endpoint.Response(http.StatusOK, "", "ShortUrlDel", "删除短链")),
		endpoint.Post("/marketing-strategies/two-level-ids", "获取二级渠道ids",
			endpoint.Handler(TwoLevelIds),
			endpoint.Tags("获取二级渠道ids"),
			endpoint.Response(http.StatusOK, "", "TwoLevelIds", "获取二级渠道ids")),
		endpoint.Post("/marketing-strategies/get-origin-id", "获取小程序原始id",
			endpoint.Handler(GetOriginId),
			endpoint.Tags("获取小程序原始id"),
			endpoint.Response(http.StatusOK, "", "GetOriginId", "获取小程序原始id")),
		endpoint.Post("/marketing-strategies/all-short-url-list", "获取全部短链列表",
			endpoint.Handler(GetAllShortUrlList),
			endpoint.Tags("获取全部短链列表"),
			endpoint.Response(http.StatusOK, "", "GetAllShortUrlList", "获取全部短链列表")),
	}
}
