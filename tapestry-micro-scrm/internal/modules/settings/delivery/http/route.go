package settingshttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-scrm/internal/modules/settings"
	"micro-scrm/internal/modules/settings/respository/mysql"
	"net/http"
)

// Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {
	service = settings.New(mysql.New())
	return []*swagger.Endpoint{
		endpoint.Post("/settings/info", "获取会话存档设置",
			endpoint.Handler(settingsInfo),
			endpoint.Tags("获取会话存档设置"),
			endpoint.Response(http.StatusOK, "", "info", "接口调用成功")),
		endpoint.Post("/settings/update", "更新会话存档设置",
			endpoint.Handler(settingsUpdate),
			endpoint.Tags("获取会话存档设置"),
			endpoint.Response(http.StatusOK, "", "info", "接口调用成功")),
	}

}
