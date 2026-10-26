package departmentshttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/department"
	"micro-scrm/internal/modules/department/respository/mysql"
	"net/http"
)

//Endpoints 路由配置
func Endpoints() []*swagger.Endpoint {
	departmentService = department.New(mysql.New())
	return []*swagger.Endpoint{
		endpoint.Post("/departments/list", "组织列表",
			endpoint.Handler(listDept),
			endpoint.Tags("组织管理"),
			endpoint.Body(department.DepartmentDetailRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: []department.Department{},
			}, "listDept", "接口调用成功")),

		endpoint.Post("/departments/detail", "组织详情",
			endpoint.Handler(detailDept),
			endpoint.Tags("组织管理"),
			endpoint.Body(department.DepartmentListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: department.Department{},
			}, "detailDept", "接口调用成功")),

		endpoint.Post("/departments/update-record", "组织更新记录",
			endpoint.Handler(updateRecord),
			endpoint.Tags("组织管理"),
			endpoint.Body(department.DepartmentListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: department.Department{},
			}, "updateRecord", "接口调用成功")),

		endpoint.Post("/departments/moment-list", "选择部门",
			endpoint.Handler(departmentMomentList),
			endpoint.Body(department.MomentListRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: department.MomentListResponse{},
			}, "departmentList", "接口调用成功")),
		endpoint.Post("/departments/user-lost-config", "导购丢失客户推送消息配置", endpoint.Handler(userLostConfig)),
		endpoint.Post("/departments/user-lost-info", "导购丢失客户推送消息详情", endpoint.Handler(userLostInfo)),
	}
}
