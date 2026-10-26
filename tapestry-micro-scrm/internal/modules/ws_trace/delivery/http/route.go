package http

import (
	"net/http"

	"micro-scrm/internal/modules/ws_trace"
	"micro-scrm/internal/modules/ws_trace/respository/mysql"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	wsTraceService = ws_trace.NewService(mysql.NewWsTrace())
	return []*swagger.Endpoint{
		endpoint.Post("/ws-trace/add", "插入或更新ws-trace记录",
			endpoint.Handler(wsAdd),
			endpoint.Tags("插入或更新ws-trace记录"),
			endpoint.Response(http.StatusOK, "",
				"getCurrentVersionInfo", "接口调用成功")),
	}
}
