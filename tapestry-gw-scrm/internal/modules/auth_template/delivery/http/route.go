package auth_template_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/auth-template/detail", "授权页配置详情",
			endpoint.Handler(Detail),
			endpoint.Tags("授权页配置详情")),
		endpoint.Post("/auth-template/save", "授权页配置保存",
			endpoint.Handler(Save),
			endpoint.Tags("授权页配置保存")),
	}
}
