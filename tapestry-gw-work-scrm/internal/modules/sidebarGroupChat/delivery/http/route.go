package sidebarGroupChathttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gw-work-scrm/internal/config"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"net/http"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {
	sidebarGroupChatSdk = microScrm.NewSidebarGroupChat(config.CfgData.Service.MicroScrm.Domain)
	return []*swagger.Endpoint{
		endpoint.Post(
			"/group-chat/info",
			"侧边栏社群详情",
			endpoint.Handler(groupChatInfo),
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
			"/group-chat/group-chat-member",
			"侧边栏社群成员列表",
			endpoint.Handler(groupChatMember),
			endpoint.Tags("侧边栏社群成员列表"),
			endpoint.Body(microScrm.GroupChatInfoRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"groupChatMember",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/group-chat/group-users-add-employees",
			"侧边栏社群成员添加员工",
			endpoint.Handler(groupUsersAddEmployees),
			endpoint.Tags("侧边栏社群成员添加员工"),
			endpoint.Body(microScrm.GroupUsersAddEmployeesRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"groupUsersAddEmployees",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/group-chat/view-join-activity",
			"侧边栏社群成员参加活动列表",
			endpoint.Handler(groupUsersJoinActivity),
			endpoint.Tags("侧边栏社群成员添加员工"),
			endpoint.Body(microScrm.GroupUsersJoinActivityRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"groupUsersAddEmployees",
				"接口调用成功",
			),
		),
	}
}
