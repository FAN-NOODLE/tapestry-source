package members_http

import (
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/members"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/members/detail", "员工详情",
			endpoint.Handler(membersDetail),
			endpoint.Tags("员工管理"),
			endpoint.Body(members.MemberDetailRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: members.Member{},
			}, "memberDetail", "接口调用成功")),
		endpoint.Post("/members/list", "员工列表",
			endpoint.Handler(membersList),
			endpoint.Tags("员工管理"),
			endpoint.Body(members.MemberListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: []members.Member{},
			}, "memberList", "接口调用成功")),

		endpoint.Post("/members/sync", "同步通讯录",
			endpoint.Handler(membersSync),
			endpoint.Tags("员工管理"),
			endpoint.Response(http.StatusOK, "", "memberSync", "接口调用成功")),

		endpoint.Post("/members/momentList", "选择员工",
			endpoint.Handler(momentList),
			endpoint.Tags("朋友圈"),
			endpoint.Body(members.MomentListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: members.MomentListResponse{},
			},
				"momentList", "接口调用成功")),
		endpoint.Post("/members/ba/list", "选择BA员工",
			endpoint.Handler(momentBaList),
			endpoint.Tags("员工列表"),
			endpoint.Body(members.MomentListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: members.MomentListResponse{},
			},
				"momentList", "接口调用成功")),

		endpoint.Post("/members/dp-members", "选择组织架构员工",
			endpoint.Handler(getSimpleMembersByDepartmentIds),
			endpoint.Tags("组织架构员工列表"),
			endpoint.Response(http.StatusOK, "", "simpleMembers", "接口调用成功")),

		endpoint.Post("/members/list-by-third-ids", "根据企微部门id集合获取员工集合",
			endpoint.Handler(MemberListByThirdIds),
			endpoint.Tags("根据企微部门id集合获取员工集合"),
			endpoint.Body(members.MemberListByThirdIdsReq{}, "请求参数", true)),
	}
}
