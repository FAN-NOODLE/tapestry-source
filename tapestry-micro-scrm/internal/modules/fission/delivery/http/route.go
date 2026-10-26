package http

import (
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/customer/respository/elastic"
	"micro-scrm/internal/modules/drainage"
	drainageMysql "micro-scrm/internal/modules/drainage/respository/mysql"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/fission"
	fissionBossApi "micro-scrm/internal/modules/fission/api/bossapi"
	fissionMysql "micro-scrm/internal/modules/fission/respository/mysql"
	"micro-scrm/internal/modules/user"
	userMysql "micro-scrm/internal/modules/user/respository/mysql"
	"micro-scrm/internal/pkg/common"
)

//Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {
	// 初始化服务
	userService := user.NewUserService(userMysql.NewUser(), userMysql.NewEvent(), elastic.NewEsRepo(), enterpriseMysql.New())
	fissionService = fission.NewService(
		fissionMysql.New(),
		drainage.NewContact(drainageMysql.NewDrainage()),
		userService,
		fissionBossApi.NewApi(),
		enterprise.New(enterpriseMysql.New()))
	fission.UserService = userService
	return []*swagger.Endpoint{
		endpoint.Post(
			"/fission/mp/exchange-prize-list",
			"裂变奖品兑换列表",
			endpoint.Handler(exchangePrizeList),
			endpoint.Body(RequestFissionUser{}, "", true),
			endpoint.Tags("裂变活动"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{Data: fission.ExchangeList{}}, "exchangePrizeList", "接口调用成功"),
		),

		endpoint.Post(
			"/fission/mp/exchange-prize-detail",
			"裂变奖品兑换详情",
			endpoint.Handler(exchangePrizeDetail),
			endpoint.Body(common.RequestID{}, "", true),
			endpoint.Tags("裂变活动"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{Data: fission.ExchangeDetail{}}, "exchangePrizeDetail", "接口调用成功"),
		),

		endpoint.Post(
			"/fission/mp/prize-list",
			"奖品列表",
			endpoint.Handler(prizeList),
			endpoint.Body(RequestActivityID{}, "", true),
			endpoint.Tags("裂变活动"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{Data: fission.PrizeList{}}, "prizeList", "接口调用成功"),
		),

		endpoint.Post(
			"/fission/mp/prize-info",
			"奖品详情",
			endpoint.Handler(prizeInfo),
			endpoint.Body(common.RequestID{}, "", true),
			endpoint.Tags("裂变活动"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{Data: fission.PrizeDetail{}}, "prizeInfo", "接口调用成功"),
		),

		endpoint.Post(
			"/fission/mp/assistant-list",
			"助力记录",
			endpoint.Handler(assistantList),
			endpoint.Body(RequestFissionUser{}, "", true),
			endpoint.Tags("裂变活动"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{Data: fission.AssistantList{}}, "assistantList", "接口调用成功"),
		),

		endpoint.Post(
			"/fission/mp/mp-home",
			"活动首页",
			endpoint.Handler(mpHome),
			endpoint.Body(RequestFissionUser{}, "", true),
			endpoint.Tags("裂变活动"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{Data: fission.MpHomeInfo{}}, "mpHome", "接口调用成功"),
		),

		endpoint.Post(
			"/fission/mp/join-activity",
			"用户参加活动",
			endpoint.Handler(joinActivity),
			endpoint.Body(RequestFissionUser{}, "", true),
			endpoint.Tags("裂变活动"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{Data: fission.MpHomeInfo{}}, "mpHome", "接口调用成功"),
		),

		endpoint.Post(
			"/fission/mp/exchange-prize",
			"兑换奖品",
			endpoint.Handler(exchangePrize),
			endpoint.Body(RequestFissionUser{}, "", true),
			endpoint.Tags("裂变活动"),
			endpoint.Response(http.StatusOK, "", "mpHome", "接口调用成功"),
		),

		endpoint.Post(
			"/fission/mp/exchange-check",
			"兑换奖品检查",
			endpoint.Handler(exchangePrizeCheck),
			endpoint.Body(RequestFissionUser{}, "", true),
			endpoint.Tags("裂变活动"),
			endpoint.Response(http.StatusOK, "", "mpHome", "接口调用成功"),
		),

		endpoint.Post(
			"/fission/mp/scroll-list",
			"滚动条列表",
			endpoint.Handler(scrollList),
			endpoint.Body(fission.ScrollListRequest{}, "", true),
			endpoint.Tags("裂变活动"),
			endpoint.Response(http.StatusOK, "", "scrollList", "接口调用成功"),
		),

		endpoint.Post(
			"/fission/activity-list", "裂变活动列表",
			endpoint.Handler(activityList),
			endpoint.Tags("activity-list"),
			endpoint.Response(http.StatusOK, "",
				"activity-list", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/activity-add", "裂变活动添加",
			endpoint.Handler(activitySet),
			endpoint.Tags("activity-add"),
			endpoint.Response(http.StatusOK, "",
				"activity-add", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/activity-detail", "裂变活动明细-统计信息",
			endpoint.Handler(activityDetail),
			endpoint.Tags("activity-detail"),
			endpoint.Response(http.StatusOK, "",
				"activity-detail", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/activity-detail-list", "裂变活动明细-列表",
			endpoint.Handler(activityDetailList),
			endpoint.Tags("activity-detail-list"),
			endpoint.Response(http.StatusOK, "",
				"activity-detail-list", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/invite-detail", "发起人邀请明细",
			endpoint.Handler(inviteDetail),
			endpoint.Tags("inviteDetail"),
			endpoint.Response(http.StatusOK, "",
				"inviteDetail", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/exchange-record", "发起人兑换记录",
			endpoint.Handler(exchangeRecord),
			endpoint.Tags("exchange-record"),
			endpoint.Response(http.StatusOK, "",
				"exchange-detail", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/prize-exchange-record", "兑换奖品记录",
			endpoint.Handler(prizeExchangeRecord),
			endpoint.Tags("prize-exchange-record"),
			endpoint.Response(http.StatusOK, "",
				"prize-exchange-detail", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/open-exchange-warn", "确认开启奖品兑换提醒",
			endpoint.Handler(openExchangeWarn),
			endpoint.Tags("open-exchange-warn"),
			endpoint.Response(http.StatusOK, "",
				"open-exchange-warn", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/exchange-total", "奖品兑换记录总数",
			endpoint.Handler(ExchangeTotal),
			endpoint.Tags("exchange-total"),
			endpoint.Response(http.StatusOK, "",
				"exchange-total", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/activity-list-simple", "活动列表全部有效活动",
			endpoint.Handler(activityListSimple),
			endpoint.Tags("activity-list-simple"),
			endpoint.Response(http.StatusOK, "",
				"activity-list-simple", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/mutual-activity-list", "选择互斥活动",
			endpoint.Handler(mutualActivityList),
			endpoint.Tags("mutual-activity-list"),
			endpoint.Response(http.StatusOK, "",
				"mutual-activity-list", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/activity-update", "裂变活动编辑",
			endpoint.Handler(activitySet),
			endpoint.Tags("activity-update"),
			endpoint.Response(http.StatusOK, "",
				"activity-update", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/chat-group-list", "社群选择列表",
			endpoint.Handler(chatGroupList),
			endpoint.Tags("chat-group-list"),
			endpoint.Response(http.StatusOK, "",
				"chat-group-list", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/activity-delete", "删除活动",
			endpoint.Handler(activityDelete),
			endpoint.Tags("activity-delete"),
			endpoint.Response(http.StatusOK, "",
				"activity-delete", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/update-status", "更新活动状态",
			endpoint.Handler(updateStatus),
			endpoint.Tags("update-status"),
			endpoint.Response(http.StatusOK, "",
				"update-status", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/activity-info", "活动详情",
			endpoint.Handler(activityInfo),
			endpoint.Tags("activity-info"),
			endpoint.Response(http.StatusOK, "",
				"activity-info", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/qw-list", "企微活动列表",
			endpoint.Handler(fissionActivityList),
			endpoint.Tags("list"),
			endpoint.Response(http.StatusOK, "",
				"activity-info", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/send", "发送兑换通知接口测试",
			endpoint.Handler(sendMessage),
			endpoint.Tags("send"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功"),
		),

		endpoint.Post(
			"/fission/mp-user", "获取公众号用户信息",
			endpoint.Handler(mpUserInfo),
			endpoint.Tags("send"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/get-exchange-config-member", "获取兑换提醒的额外员工",
			endpoint.Handler(getExchangeConfigMember),
			endpoint.Tags("get-exchange-config-member"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/get-style-config", "配色列表",
			endpoint.Handler(styleConfigList),
			endpoint.Tags("get-style-config"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/add-config", "自定义添加配色",
			endpoint.Handler(addStyleConfig),
			endpoint.Tags("get-style-config"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/del-config", "删除自定义配色",
			endpoint.Handler(delStyleConfig),
			endpoint.Tags("get-style-config"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/set-style-config", "设置配色",
			endpoint.Handler(setStyleConfig),
			endpoint.Tags("get-style-config"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功"),
		),
	}
}
