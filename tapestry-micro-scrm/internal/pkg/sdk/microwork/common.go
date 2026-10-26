package microwork

import (
	"fmt"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/sdk"
)

// ResponseBase shop接口返回结构体头
type ResponseBase struct {
	sdk.ResponseBase
}

// httpClient http请求代理
var httpClient = sdk.NewHTTPClient(sdk.ParseJSON{})

// GetURL 获取shop请求链接
func GetURL(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.MicroWork.Domain)
}
