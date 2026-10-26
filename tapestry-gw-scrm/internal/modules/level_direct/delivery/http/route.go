package level_direct_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/level-direct/list", "等级直达列表",
			endpoint.Handler(List),
			endpoint.Tags("等级直达列表")),
		endpoint.Post("/level-direct/add", "等级直达添加",
			endpoint.Handler(Add),
			endpoint.Tags("等级直达添加")),
		endpoint.Post("/level-direct/update", "等级直达添加",
			endpoint.Handler(Update),
			endpoint.Tags("等级直达添加")),
		endpoint.Post("/level-direct/delete", "等级直达添加",
			endpoint.Handler(Delete),
			endpoint.Tags("等级直达添加")),
		endpoint.Post("/level-direct/detail", "等级直达详情",
			endpoint.Handler(Detail),
			endpoint.Tags("等级直达详情")),
	}
}
