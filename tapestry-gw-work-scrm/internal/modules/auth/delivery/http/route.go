package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-work-scrm/internal/config"
	"gw-work-scrm/internal/modules/auth"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"net/http"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {
	// 初始化SDK
	memberSdk = microScrm.NewMember(config.CfgData.Service.MicroScrm.Domain)
	enterpriseSdk = microScrm.NewEnterprise(config.CfgData.Service.MicroScrm.Domain)

	return []*swagger.Endpoint{
		endpoint.Post(
			"/auth/redirect_url",
			"企业微信授权地址",
			endpoint.Handler(redirectUrl),
			endpoint.Tags("用户认证"),
			endpoint.Body(auth.RedirectUrlRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				ginplus.APIResponse{Data: auth.RedirectUrlResponse{}},
				"redirect_url",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/auth/login",
			"用户授权登录",
			endpoint.Handler(authLogin),
			endpoint.Tags("用户认证"),
			endpoint.Body(auth.LoginAuthRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				ginplus.APIResponse{Data: auth.LoginAuthResponse{}},
				"auth",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/auth/config",
			"企业签名",
			endpoint.Handler(jsConfig),
			endpoint.Tags("企业签名"),
			endpoint.Body(auth.LoginAuthRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				ginplus.APIResponse{Data: auth.LoginAuthResponse{}},
				"auth",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/auth/agent-config",
			"应用签名",
			endpoint.Handler(jsAgentConfig),
			endpoint.Tags("应用签名"),
			endpoint.Body(auth.LoginAuthRequest{}, "请求参数", true),
			endpoint.Response(
				http.StatusOK,
				ginplus.APIResponse{Data: auth.LoginAuthResponse{}},
				"auth",
				"接口调用成功",
			),
		),

		endpoint.Get(
			"/auth/mock",
			"模拟登录",
			endpoint.Handler(mock),
			endpoint.Tags("模拟登录"),
		),
	}
}
