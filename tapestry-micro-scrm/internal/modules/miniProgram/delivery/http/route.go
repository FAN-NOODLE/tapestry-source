package miniProgramHttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-scrm/internal/modules/miniProgram"
	"micro-scrm/internal/modules/miniProgram/respository/mysql"
	"net/http"
)

// Endpoints
func Endpoints() []*swagger.Endpoint {

	s = miniProgram.New(mysql.NewMiniProgram())
	return []*swagger.Endpoint{
		endpoint.Post("/mini-program/add", "新增小程序",
			endpoint.Handler(add),
			endpoint.Tags("新增小程序"),
			endpoint.Response(http.StatusOK, "", "mini-program", "接口调用成功")),

		endpoint.Post("/mini-program/authorize", "授权小程序",
			endpoint.Handler(authorize),
			endpoint.Tags("授权小程序"),
			endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),

		endpoint.Post("/mini-program/list", "授权小程序",
			endpoint.Handler(list),
			endpoint.Tags("小程序列表"),
			endpoint.Response(http.StatusOK, "", "mini-program", "接口调用成功")),

		endpoint.Post("/mini-program/del", "删除小程序",
			endpoint.Handler(del),
			endpoint.Tags("删除小程序"),
			endpoint.Response(http.StatusOK, "", "mini-program", "接口调用成功")),

		endpoint.Post("/mini-program/edit", "编辑小程序",
			endpoint.Handler(edit),
			endpoint.Tags("编辑小程序"),
			endpoint.Response(http.StatusOK, "", "source-group", "接口调用成功")),
	}
}
