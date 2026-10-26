package microshop

import (
	"fmt"
	"gw-scrm/internal/config"
)

var urlEnum = struct {
	CouponListUrl       string
	LiveListUrl         string
	ActivityPages       string
	StoreList           string
	AuthorizationDetail string
	AuthorizationSave   string
}{
	// 优惠券列表
	CouponListUrl: "%s/service/v1/template/home/coupon-list",
	//直播列表
	LiveListUrl:   "%s/service/v1/live-broadcast/list",
	StoreList:     "%s/service/v1/store/store/list-by-enterprise",
	ActivityPages: "%s/service/v1/new-template/list",

	// 授权页配置
	AuthorizationDetail: "%s/service/v1/new-template/authorization-detail",
	AuthorizationSave:   "%s/service/v1/new-template/authorization-save",
}

// 获取shop请求链接
func GetUrl(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.Shop.Domain)
}
