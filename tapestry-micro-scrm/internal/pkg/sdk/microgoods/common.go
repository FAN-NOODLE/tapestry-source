package microgoods

import (
	"fmt"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	SpuList string
}{
	// 优惠券列表
	SpuList: "%s/micro-goods/v1/goods/spu-list",
}

// ResponseBase 服务接口返回数据结构体头
type ResponseBase struct {
	sdk.ResponseBase
}

// httpClient http请求代理
var httpClient = sdk.NewHTTPClient(sdk.ParseJSON{})

// GetURL 获取shop请求链接
func GetURL(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.MicroGoods.Domain)
}
