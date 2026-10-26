package overtimehttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gw-scrm/internal/config"
	"gw-scrm/pkg/micro-scrm-sdk/microScrm"
)

var customerSDK *microScrm.CustomerApi

func Endpoints() []*swagger.Endpoint {
	customerSDK = microScrm.NewCustomer(config.CfgData.Service.MicroScrm.Domain)
	return []*swagger.Endpoint{
		endpoint.Post("/overtime/list", "超时记录列表",
			endpoint.Handler(overtimeList),
			endpoint.Tags("超时记录列表")),
	}
}
