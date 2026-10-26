package legal_right_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/legal-right/save", "权益保存",
			endpoint.Handler(save),
			endpoint.Tags("权益保存")),

		endpoint.Post("/legal-right/list", "权益列表",
			endpoint.Handler(list),
			endpoint.Tags("权益列表")),

		endpoint.Post("/legal-right/update-status", "修改权益状态",
			endpoint.Handler(UpdateStatus),
			endpoint.Tags("修改权益状态")),
	}
}
