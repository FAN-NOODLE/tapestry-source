package microconsole

import (
	"fmt"
	"gw-scrm/internal/config"
	"gw-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	HelpListUri string
}{
	// 获取外部联系人凭证
	HelpListUri: "%s/gw-console/v1/help-center/list",
}

// GetUrl 获取请求链接
func GetUrl(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.Console.Domain)
}

var httpClient = sdk.NewHttpClient(sdk.ParseJson{})
