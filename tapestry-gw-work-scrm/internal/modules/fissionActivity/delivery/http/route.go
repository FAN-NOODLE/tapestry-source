package fissionActivityhttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gw-work-scrm/internal/config"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"net/http"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {
	FissionActivitySdk = microScrm.NewFissionActivity(config.CfgData.Service.MicroScrm.Domain)
	return []*swagger.Endpoint{
		endpoint.Post(
			"/fission_activity/list",
			"裂变活动列表",
			endpoint.Handler(fissionActivityList),
			endpoint.Tags("裂变活动列表"),
			endpoint.Body(microScrm.FissionActivityListRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"groupChatInfo",
				"接口调用成功",
			),
		),
	}
}
