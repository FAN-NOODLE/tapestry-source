package overtimehttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/overtime"
	"micro-scrm/internal/modules/overtime/respository/mysql"
	"net/http"
)

// Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {
	service = overtime.New(mysql.New())
	return []*swagger.Endpoint{
		endpoint.Post("/overtime/list", "获取聊天超时列表",
			endpoint.Handler(overtimeList),
			endpoint.Body(overtime.ListRequest{}, "/overtime/list", true),
			endpoint.Tags("获取聊天超时列表"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: overtime.PageListResponse{},
			}, "overtimeList", "接口调用成功")),
	}

}
