package fissionhttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{

		endpoint.Post("/fission/export-exchange-record", "兑奖记录导出",
			endpoint.Handler(exportExchangeRecord),
			endpoint.Tags("兑奖记录导出"),
			endpoint.Response(http.StatusOK, "", "exportExchangeRecord", "接口调用成功")),
		endpoint.Post("/fission/activity-detail", "裂变活动详情",
			endpoint.Handler(activityDetail),
			endpoint.Tags("裂变活动详情"),
			endpoint.Response(http.StatusOK, "", "activityDetail", "接口调用成功")),
		endpoint.Post("/fission/activity-detail-list", "裂变活动详情列表",
			endpoint.Handler(activityDetailList),
			endpoint.Tags("裂变活动详情列表"),
			endpoint.Response(http.StatusOK, "", "activityDetailList", "接口调用成功")),
		endpoint.Post("/fission/exchange-record", "发起人兑换记录",
			endpoint.Handler(exchangeRecord),
			endpoint.Tags("发起人兑换记录"),
			endpoint.Response(http.StatusOK, "", "exchangeRecord", "接口调用成功")),
		endpoint.Post("/fission/prize-exchange-record", "兑换奖品记录",
			endpoint.Handler(prizeExchangeRecord),
			endpoint.Tags("兑换奖品记录"),
			endpoint.Response(http.StatusOK, "", "prizeExchangeRecord", "接口调用成功")),
		endpoint.Post("/fission/open-exchange-warn", "确认开启奖品兑换提醒",
			endpoint.Handler(openExchangeWarn),
			endpoint.Tags("确认开启奖品兑换提醒"),
			endpoint.Response(http.StatusOK, "", "openExchangeWarn", "接口调用成功")),
		endpoint.Post("/fission/exchange-total", "奖品兑换记录总数",
			endpoint.Handler(exchangeTotal),
			endpoint.Tags("exchange-total"),
			endpoint.Response(http.StatusOK, "", "exchangeTotal", "接口调用成功")),
		endpoint.Post("/fission/invite-detail", "发起人邀请明细",
			endpoint.Handler(inviteDetail),
			endpoint.Tags("发起人邀请明细"),
			endpoint.Response(http.StatusOK, "", "inviteDetail", "接口调用成功")),
		endpoint.Post("/fission/activity-list-simple", "裂变活动列表 全部",
			endpoint.Handler(activityListSimple),
			endpoint.Tags("裂变活动列表 全部"),
			endpoint.Response(http.StatusOK, "", "activityListSimple", "接口调用成功")),
		endpoint.Post("/fission/activity-list", "裂变活动列表",
			endpoint.Handler(activityList),
			endpoint.Tags("裂变活动列表"),
			endpoint.Response(http.StatusOK, "", "fission", "接口调用成功")),

		endpoint.Post("/fission/activity-add", "裂变活动新增",
			endpoint.Handler(activityAdd),
			endpoint.Tags("裂变活动新增"),
			endpoint.Response(http.StatusOK, "", "fission", "接口调用成功")),

		endpoint.Post("/fission/activity-update", "裂变活动编辑",
			endpoint.Handler(activityUpdate),
			endpoint.Tags("裂变活动编辑"),
			endpoint.Response(http.StatusOK, "", "fission", "接口调用成功")),

		endpoint.Post("/fission/activity-delete", "裂变活动删除",
			endpoint.Handler(activityDelete),
			endpoint.Tags("裂变活动删除"),
			endpoint.Response(http.StatusOK, "", "fission", "接口调用成功")),

		endpoint.Post("/fission/update-status", "裂变活动更新状态",
			endpoint.Handler(updateStatus),
			endpoint.Tags("裂变活动更新状态"),
			endpoint.Response(http.StatusOK, "", "fission", "接口调用成功")),

		endpoint.Post("/fission/activity-info", "裂变活动详情",
			endpoint.Handler(activityInfo),
			endpoint.Tags("裂变活动详情"),
			endpoint.Response(http.StatusOK, "", "fission", "接口调用成功")),

		endpoint.Post("/fission/chat-group-list", "裂变活动-选择社群列表",
			endpoint.Handler(chatGroupList),
			endpoint.Tags("裂变活动-选择社群列表"),
			endpoint.Response(http.StatusOK, "", "fission", "接口调用成功")),

		endpoint.Post("/fission/mutual-activity-list", "裂变活动-互斥活动选择列表",
			endpoint.Handler(mutualActivityList),
			endpoint.Tags("裂变活动-互斥活动选择列表"),
			endpoint.Response(http.StatusOK, "", "fission", "接口调用成功")),
		endpoint.Post(
			"/fission/get-exchange-config-member", "获取兑换提醒的额外员工列表",
			endpoint.Handler(getExchangeConfigMember),
			endpoint.Tags("get-exchange-config-member"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/get-style-config", "获取兑换提醒的额外员工列表",
			endpoint.Handler(getStyleConfig),
			endpoint.Tags("get-exchange-config-member"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/set-style-config", "获取兑换提醒的额外员工列表",
			endpoint.Handler(setStyleConfig),
			endpoint.Tags("get-exchange-config-member"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/add-style-config", "获取兑换提醒的额外员工列表",
			endpoint.Handler(addStyleConfig),
			endpoint.Tags("get-exchange-config-member"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功"),
		),
		endpoint.Post(
			"/fission/del-style-config", "获取兑换提醒的额外员工列表",
			endpoint.Handler(delStyleConfig),
			endpoint.Tags("get-exchange-config-member"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功"),
		),
	}
}
