package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/enterprise/respository/mysql"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	// 初始化服务
	service = enterprise.New(mysql.New())

	return []*swagger.Endpoint{
		endpoint.Post(
			"/js/config",
			"企业签名",
			endpoint.Handler(jsConfig),
			endpoint.Tags("js-config"),
			endpoint.Response(
				http.StatusOK,
				"",
				"MemberDetail",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/js/agent-config",
			"应用签名",
			endpoint.Handler(jsAgentConfig),
			endpoint.Tags("js-config"),
			endpoint.Response(
				http.StatusOK,
				"",
				"MemberDetail",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/enterprise/official/redirect-url",
			"公众号信息",
			endpoint.Handler(redirectURL),
			endpoint.Tags("official-info"),
			endpoint.Response(
				http.StatusOK,
				"",
				"official",
				"接口调用成功",
			),
		),

		endpoint.Post(
			"/enterprise/official/auth",
			"公众号登录",
			endpoint.Handler(officialAuth),
			endpoint.Tags("official-auth"),
			endpoint.Response(
				http.StatusOK,
				"",
				"official",
				"接口调用成功",
			),
		),
	}
}
