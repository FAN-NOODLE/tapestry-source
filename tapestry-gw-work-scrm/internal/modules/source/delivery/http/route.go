package sourcehttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gw-work-scrm/internal/config"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"net/http"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {
	sourceSdk = microScrm.NewScource(config.CfgData.Service.MicroScrm.Domain)
	return []*swagger.Endpoint{
		endpoint.Post(
			"/source/source-list",
			"素材列表",
			endpoint.Handler(sourceList),
			endpoint.Tags("侧边栏-客户画像"),
			endpoint.Response(
				http.StatusOK,
				"",
				"source-list",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/source/source-group-list",
			"素材分组列表",
			endpoint.Handler(sourceGroupList),
			endpoint.Tags("侧边栏-客户画像"),
			endpoint.Response(
				http.StatusOK,
				"",
				"source-group-list",
				"接口调用成功",
			),
		),
		endpoint.Post(
			"/source/make-media",
			"生成临时素材",
			endpoint.Handler(makeMedia),
			endpoint.Tags("侧边栏-客户画像"),
			endpoint.Response(
				http.StatusOK,
				"",
				"source-group-list",
				"接口调用成功",
			),
		),
	}
}
