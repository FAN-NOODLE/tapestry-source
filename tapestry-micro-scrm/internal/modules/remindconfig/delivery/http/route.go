package remindconfigthttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-scrm/internal/modules/members"
	membersMysql "micro-scrm/internal/modules/members/respository/mysql"
	"micro-scrm/internal/modules/remindconfig"
	"micro-scrm/internal/modules/remindconfig/respository/mysql"
	"net/http"
)

// Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {
	service = remindconfig.New(mysql.New())
	memberService = members.New(membersMysql.New())

	return []*swagger.Endpoint{
		endpoint.Post("/remind-config/config-info", "获取提醒设置",
			endpoint.Handler(configInfo),
			endpoint.Tags("获取提醒设置"),
			endpoint.Response(http.StatusOK, "", "config-info", "接口调用成功")),

		endpoint.Post("/remind-config/config-save", "保存提醒设置",
			endpoint.Handler(configSave),
			endpoint.Tags("保存提醒设置"),
			endpoint.Response(http.StatusOK, "", "config-save", "接口调用成功")),
	}

}
