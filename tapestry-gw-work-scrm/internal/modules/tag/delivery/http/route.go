package taghttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gw-work-scrm/internal/config"
	"gw-work-scrm/internal/modules/auth"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"net/http"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {
	tagSdk = microScrm.NewTag(config.CfgData.Service.MicroScrm.Domain)
	return []*swagger.Endpoint{
		endpoint.Post(
			"/tag/list",
			"企业标签库-黄晓敏",
			endpoint.Handler(customerTagList),
			endpoint.Tags("侧边栏-客户画像"),
			endpoint.Body(auth.RedirectUrlRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				"",
				"tag-list",
				"接口调用成功",
			),
		),
	}
}
