package sidebarGroupChathttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-scrm/internal/modules/sidebarGroupChat"
	sidebarGroupChatMysql "micro-scrm/internal/modules/sidebarGroupChat/respository/mysql"
	"net/http"
)

//Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {
	sidebarGroupChatService = sidebarGroupChat.New(sidebarGroupChatMysql.New())
	return []*swagger.Endpoint{

		endpoint.Post("/sidebar/group-chat/info", "侧边栏社群详情",
			endpoint.Handler(groupChatInfo),
			endpoint.Tags("侧边栏社群详情"),
			endpoint.Response(http.StatusOK, "",
				"groupChatInfo", "接口调用成功")),

		endpoint.Post("/sidebar/group-chat/group-chat-member", "侧边栏社群成员列表",
			endpoint.Handler(groupChatMember),
			endpoint.Tags("侧边栏社群成员列表"),
			endpoint.Response(http.StatusOK, "",
				"groupChatMember", "接口调用成功")),

		endpoint.Post("/sidebar/group-chat/group-users-add-employees", "群成员添加员工列表",
			endpoint.Handler(groupUsersAddEmployees),
			endpoint.Tags("侧边栏社群成员添加员工"),
			endpoint.Response(http.StatusOK, "",
				"groupUsersAddEmployees", "接口调用成功")),

		endpoint.Post("/sidebar/group-chat/group-users-join-activity", "群成员参加裂变活动列表",
			endpoint.Handler(groupUsersJoinActivity),
			endpoint.Tags("侧边栏社群成员添加员工"),
			endpoint.Response(http.StatusOK, "",
				"groupUsersAddEmployees", "接口调用成功")),
	}
}
