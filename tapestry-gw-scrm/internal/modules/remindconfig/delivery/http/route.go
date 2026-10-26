package remindconfig_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/remindconfig"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {

	return []*swagger.Endpoint{
		endpoint.Post("/remind-config/config-info", "获取超时提醒设置信息",
			endpoint.Handler(configInfo),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: remindconfig.ConfigInfoResponse{},
			}, "OvertimeInfo", "接口调用成功")),
		endpoint.Post("/remind-config/config-save", "获取超时提醒设置保存",
			endpoint.Handler(configSave),
			endpoint.Tags("获取超时提醒设置保存"),
			endpoint.Response(http.StatusOK, "", "configSave", "接口调用成功")),
	}
}
