package microauthority

import (
	"fmt"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	WxQrcode, SendMsg string
}{
	// 生成小程序太阳码
	WxQrcode: "%s/micro-authority/v1/wx/qr-code",
	SendMsg:  "%s/micro-authority/v1/common/send-message",
}

// ResponseBase 服务接口返回数据结构体头
type ResponseBase struct {
	sdk.ResponseBase
}

// httpClient http请求代理
var httpClient = sdk.NewHTTPClient(sdk.ParseJSON{})

// GetURL 获取shop请求链接
func GetURL(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.MicroAuthority.Domain)
}
