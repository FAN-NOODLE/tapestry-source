package membershttp

import (
	departmentSQL "micro-scrm/internal/modules/department/respository/mysql"
	"net/http"

	"micro-scrm/internal/modules/department"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/members"
	"micro-scrm/internal/modules/members/respository/mysql"
)

//Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {
	service = members.New(mysql.New())
	service.Enterprise = enterprise.New(enterpriseMysql.New())
	departmentService = department.New(departmentSQL.New())
	return []*swagger.Endpoint{
		endpoint.Post("/members/detail", "成员详情",
			endpoint.Handler(memberDetail),
			endpoint.Tags("员工管理"),
			endpoint.Body(members.MemberDetailRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: members.Member{},
			}, "MemberDetail", "接口调用成功")),

		endpoint.Post("/members/list", "成员列表",
			endpoint.Handler(memberList),
			endpoint.Tags("员工管理"),
			endpoint.Body(members.MemberListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: []members.Member{},
			}, "MemberList", "接口调用成功")),

		endpoint.Post("/members/sync", "同步通讯录",
			endpoint.Handler(syncContactAndDept),
			endpoint.Tags("组织成员"),
			endpoint.Response(http.StatusOK, "",
				"memberSync", "接口调用成功")),

		endpoint.Post("/members/add", "成员添加",
			endpoint.Handler(memberAdd),
			endpoint.Tags("组织成员"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: members.SyncContactResponse{},
			},
				"MemberAdd", "接口调用成功")),

		endpoint.Post("/members/moment-list", "选择员工",
			endpoint.Handler(momentList),
			endpoint.Tags("朋友圈"),
			endpoint.Body(members.MomentListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: members.MomentListResponse{},
			},
				"momentList", "接口调用成功")),

		endpoint.Post("/members/simple-list", "成员列表",
			endpoint.Handler(simpleList),
			endpoint.Tags("成员管理"),
			endpoint.Body(members.MemberListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: []members.Member{},
			},
				"momentList", "接口调用成功")),

		endpoint.Post("/members/auth-info", "用户授权登录",
			endpoint.Handler(authInfo),
			endpoint.Tags("用户授权"),
			endpoint.Response(http.StatusOK, "",
				"MemberDetail", "接口调用成功")),

		endpoint.Post("/members/detail-userid", "成员详情",
			endpoint.Handler(memberDetailByUserid),
			endpoint.Tags("员工管理"),
			endpoint.Body(members.MemberDetailRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: members.Member{},
			}, "MemberDetail", "接口调用成功")),
		endpoint.Post("/members/ba-list", "选择BA员工",
			endpoint.Handler(membersBaList),
			endpoint.Tags("朋友圈"),
			endpoint.Body(members.MomentListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: members.MomentListResponse{},
			},
				"momentList", "接口调用成功")),
		endpoint.Post("/members/dp-members", "选择部门下的所有员工", endpoint.Handler(getMembersOfDepartments)),

		endpoint.Post("/members/code-session", "用户授权登录",
			endpoint.Handler(code2Session),
			endpoint.Tags("用户授权登录"),
			endpoint.Response(http.StatusOK, "",
				"MemberDetail", "接口调用成功")),
	}
}
