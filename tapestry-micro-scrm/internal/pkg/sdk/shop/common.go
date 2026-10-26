package shop

import (
	"fmt"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	WxCouponListByIds   string
	LiveListUrl         string
	ActivityPageListUrl string
}{
	// 优惠券列表
	WxCouponListByIds:   "%s/service/v1/wx-coupon/ids-list",
	LiveListUrl:         "%s/service/v1/live-broadcast/range-live-broadcast",
	ActivityPageListUrl: "%s/service/v1/new-template/range-template",
}

// ResponseBase shop接口返回结构体头
type ResponseBase struct {
	sdk.ResponseBase
}

// httpClient http请求代理
var httpClient = sdk.NewHTTPClient(sdk.ParseJSON{})

// GetURL 获取shop请求链接
func GetURL(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.Shop.Domain)
}
