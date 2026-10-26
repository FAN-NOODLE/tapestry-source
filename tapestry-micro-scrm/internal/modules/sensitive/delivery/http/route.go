package sensitivehttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/sensitive"
	"micro-scrm/internal/modules/sensitive/respository/mysql"
	"net/http"
)

// Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {
	service = sensitive.New(mysql.New())
	return []*swagger.Endpoint{
		endpoint.Post("/sensitive/list", "获取敏感词列表",
			endpoint.Handler(sensitiveList),
			endpoint.Body(sensitive.ListRequest{}, "", true),
			endpoint.Tags("敏感词"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: sensitive.ResponseList{},
			}, "sensitiveList", "接口调用成功")),

		endpoint.Post("/sensitive/save", "敏感词保存",
			endpoint.Handler(sensitiveSave),
			endpoint.Body(sensitive.SaveRequest{}, "", true),
			endpoint.Tags("敏感词"),
			endpoint.Response(http.StatusOK, "", "sensitiveSave", "接口调用成功")),

		endpoint.Post("/sensitive/delete", "敏感词删除",
			endpoint.Handler(sensitiveDelete),
			endpoint.Body(sensitive.DeleteRequest{}, "", true),
			endpoint.Tags("敏感词"),
			endpoint.Response(http.StatusOK, "", "sensitiveDelete", "接口调用成功")),

		endpoint.Post("/sensitive/detail", "敏感词详情",
			endpoint.Handler(sensitiveDetail),
			endpoint.Body(sensitive.DeleteRequest{}, "", true),
			endpoint.Tags("敏感词"),
			endpoint.Response(http.StatusOK, "", "sensitiveDetail", "接口调用成功")),

		endpoint.Post("/sensitive/detail-list", "敏感词明细列表",
			endpoint.Handler(sensitiveDetailList),
			endpoint.Body(sensitive.DetailListRequest{}, "", true),
			endpoint.Tags("敏感词"),
			endpoint.Response(http.StatusOK, "", "sensitiveDetailList", "接口调用成功")),
	}

}
