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
	activitySdk = microScrm.NewActivity(config.CfgData.Service.MicroScrm.Domain)

	return []*swagger.Endpoint{
		endpoint.Post(
			"/activity/goods-list",
			"活动商品列表",
			endpoint.Handler(goodsList),
			endpoint.Tags("商品"),
			endpoint.Response(
				http.StatusOK,
				"",
				"goods-list",
				"接口调用成功",
			),
		),
	}
}
