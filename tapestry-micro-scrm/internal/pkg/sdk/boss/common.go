package boss

import (
	"fmt"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	BossEnterpriseDetail     string
	EnterpriseWorkPluginList string
	BossEnterpriseSmsTotal   string
	EnterpriseSmsDeplete     string
	BossEnterpriseList       string
}{
	// 优惠券列表
	BossEnterpriseDetail:     "%s/service/v1/enterprise/detail",
	BossEnterpriseList:       "%s/service/v1/enterprise/list",
	EnterpriseWorkPluginList: "%s/service/v1/enterprise-work-plugin/list",
	BossEnterpriseSmsTotal:   "%s/service/v1/enterprise-payment/get-sms",
	EnterpriseSmsDeplete:     "%s/service/v1/enterprise-payment/deplete-sms",
}

// ResponseBase 服务接口返回数据结构体头
type ResponseBase struct {
	sdk.ResponseBase
}

// httpClient http请求代理
//var httpClient = sdk.NewHTTPClient(sdk.ParseJSON{})

// GetURL 获取shop请求链接
func GetURL(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.Boss.Domain)
}
