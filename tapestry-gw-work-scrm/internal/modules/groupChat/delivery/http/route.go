package http

import (
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gw-work-scrm/internal/config"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {
	groupChatSdk = microScrm.NewGroupChat(config.CfgData.Service.MicroScrm.Domain)
	return []*swagger.Endpoint{
		endpoint.Post(
			"/chat-group/index_list",
			"社群画像-数据分析",
			endpoint.Handler(index),
			endpoint.Tags("社群画像-数据分析"),
			endpoint.Body(microScrm.GroupChatStaticsRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "index_list", "接口调用成功"),
		),

		endpoint.Post(
			"/chat-group/members_guide",
			"未入群用户/会员列表/非会员列表（分页列表）",
			endpoint.Handler(membersGuide),
			endpoint.Tags("未入群用户/会员列表/非会员列表（分页列表）"),
			endpoint.Body(microScrm.GroupChatStaticsRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"members_guide",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/chat-group/members_info",
			"群成员列表",
			endpoint.Handler(members),
			endpoint.Tags("群成员列表"),
			endpoint.Body(microScrm.GroupChatStaticsRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"members",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/chat-group/data-update-time",
			"数据更新时间",
			endpoint.Handler(getLastDataUpdateTime),
			endpoint.Tags("数据更新时间"),
			endpoint.Response(
				http.StatusOK,
				"",
				"getDataUpdateTime",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/chat-group/activity-info",
			"活动详情",
			endpoint.Handler(getActivityInfo),
			endpoint.Tags("活动详情"),
			endpoint.Response(
				http.StatusOK,
				"",
				"activity",
				"接口调用成功",
			),
		),
	}
}
