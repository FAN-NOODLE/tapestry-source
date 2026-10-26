package userhttp

import (
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-scrm/internal/modules/customer"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	customerMysql "micro-scrm/internal/modules/customer/respository/mysql"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/marketing"
	marketingMysql "micro-scrm/internal/modules/marketing/respository/mysql"
	"micro-scrm/internal/modules/user"
	userMysql "micro-scrm/internal/modules/user/respository/mysql"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {
	userService = user.NewUserService(userMysql.NewUser(), userMysql.NewEvent(), userEs.NewEsRepo(), enterpriseMysql.New())

	orderService = user.NewOrderService()
	enterpriseService = enterprise.New(enterpriseMysql.New())
	marketService = marketing.New(marketingMysql.New(), userEs.NewEsRepo())
	customerService = customer.New(customerMysql.NewCustomer(), customerMysql.NewWorkUser(), userEs.NewEsRepo())

	return []*swagger.Endpoint{

		endpoint.Post("/user/list", "客户列表",
			endpoint.Handler(customerList),
			endpoint.Body(customer.ListRequest{}, "", true),
			endpoint.Tags("客户列表"),
			endpoint.Response(http.StatusOK, customer.ListResponse{}, "customer", "接口调用成功")),

		endpoint.Post("/user/tag-group-user-count", "客户分组标签人数",
			endpoint.Handler(tagGroupUserCount),
			endpoint.Tags("客户标签列表"),
			endpoint.Response(http.StatusOK, "", "tagGroupList", "接口调用成功")),

		endpoint.Post("/user/list-config", "客户列表配置",
			endpoint.Handler(listConfig),
			endpoint.Body("", "", true),
			endpoint.Tags("客户列表配置"),
			endpoint.Response(http.StatusOK, user.ListConfigResponse{}, "customer", "接口调用成功")),
	}
}
