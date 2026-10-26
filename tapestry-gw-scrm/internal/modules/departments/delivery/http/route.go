package departments_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/departments"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/departments/detail", "组织详情",
			endpoint.Handler(departmentDetail),
			endpoint.Tags("组织管理"),
			endpoint.Body(departments.Department{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: departments.Department{},
			}, "departments", "接口调用成功")),
		endpoint.Post("/departments/list", "组织列表",
			endpoint.Handler(departmentList),
			endpoint.Tags("组织管理"),
			endpoint.Body(departments.DepartmentListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: []departments.Department{},
			}, "departments", "接口调用成功")),

		endpoint.Post("/departments/update-record", "更新记录",
			endpoint.Handler(updateRecord),
			endpoint.Tags("组织管理"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: departments.DepartmentUpdateRecord{},
			}, "updateRecord", "接口调用成功")),

		endpoint.Post("/departments/moment-list", "部门展示",
			endpoint.Handler(departmentMomentList),
			endpoint.Body(departments.MomentListRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: departments.MomentListResponse{},
			}, "departmentList", "接口调用成功")),

		endpoint.Post("/departments/user-lost-config", "导购丢失客户推送消息配置", endpoint.Handler(userLostConfig)),
		endpoint.Post("/departments/user-lost-info", "导购丢失客户推送消息配置", endpoint.Handler(userLostInfo)),
	}
}
