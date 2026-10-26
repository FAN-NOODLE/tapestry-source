package work_external_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/work-external"
	"micro-scrm/internal/modules/work-external/respository/click_house"
	"net/http"
)

// Endpoints 路由配置
func Endpoints() []*swagger.Endpoint {

	workService = workexternal.New(click_house.NewRepo())
	enterpriseService = enterprise.New(mysql.New())
	return []*swagger.Endpoint{
		endpoint.Post("/work-external/calc", "外部联系人计算",
			endpoint.Handler(Calc),
			endpoint.Tags("外部联系人计算"),
			endpoint.Response(http.StatusOK, "",
				"calc", "接口调用成功")),
	}
}
