package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gw-work-scrm/internal/config"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"net/http"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {
	// 初始化服务
	posterSdk = microScrm.NewPoster(config.CfgData.Service.MicroScrm.Domain)

	return []*swagger.Endpoint{
		endpoint.Post(
			"/marketing-promotion/poster",
			"下载海报",
			endpoint.Handler(poster),
			endpoint.Tags("下载海报"),
			endpoint.Response(
				http.StatusOK,
				"",
				"marketing-promotion/poster",
				"接口调用成功",
			),
		),
	}
}
