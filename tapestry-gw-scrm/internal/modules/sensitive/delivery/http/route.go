package sensitivehttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/config"
	"gw-scrm/internal/modules/sensitive"
	"gw-scrm/pkg/micro-scrm-sdk/microScrm"
	"net/http"
)

var customerSDK *microScrm.CustomerApi

func Endpoints() []*swagger.Endpoint {
	customerSDK = microScrm.NewCustomer(config.CfgData.Service.MicroScrm.Domain)
	return []*swagger.Endpoint{
		endpoint.Post("/sensitive/list", "敏感词列表",
			endpoint.Handler(sensitiveList),
			endpoint.Tags("客户统计"),
			endpoint.Body(sensitive.ListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: sensitive.ResponseList{},
			}, "SensitiveList", "接口调用成功")),

		endpoint.Post("/sensitive/save", "敏感词保存",
			endpoint.Handler(sensitiveSave),
			endpoint.Tags("敏感词"),
			endpoint.Response(http.StatusOK, "", "sensitiveSave", "接口调用成功")),

		endpoint.Post("/sensitive/delete", "敏感词删除接口",
			endpoint.Handler(sensitiveDelete),
			endpoint.Tags("敏感词"),
			endpoint.Response(http.StatusOK, "", "sensitiveDelete", "接口调用成功")),

		endpoint.Post("/sensitive/detail", "敏感词详情接口",
			endpoint.Handler(sensitiveDetail),
			endpoint.Tags("敏感词"),
			endpoint.Response(http.StatusOK, "", "sensitiveDetail", "接口调用成功")),

		endpoint.Post("/sensitive/detail-list", "敏感词明细列表",
			endpoint.Handler(sensitiveDetailList),
			endpoint.Tags("敏感词"),
			endpoint.Response(http.StatusOK, "", "sensitiveDetailList", "接口调用成功")),
	}
}
