package open_http

import (
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/index"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/open/order", "wesocial 效果追踪订单回传接口",
			endpoint.Handler(openOrder),
			endpoint.Tags("wesocial 效果追踪订单回传接口"),
			endpoint.Body(index.PrivateDomain{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: index.PrivateDomain{},
			}, "indexData", "接口调用成功")),
	}
}
