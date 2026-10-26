package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {

	return []*swagger.Endpoint{
		endpoint.Post("/version/current-version-info", "获取当前版本信息",
			endpoint.Handler(getCurrentVersionInfo),
			endpoint.Tags("获取当前版本信息"),
			endpoint.Response(http.StatusOK, "",
				"getCurrentVersionInfo", "接口调用成功")),
	}
}
