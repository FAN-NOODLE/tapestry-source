package microgoods

import (
	"fmt"
	"gw-scrm/internal/config"
)

var urlEnum = struct {
	GoodsSpuList string
}{
	// 获取标签组列表
	GoodsSpuList: "%s/micro-goods/v1/goods/spu-list",
}

// 获取cdp请求链接
func GetUrl(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.MicroGoods.Domain)
}
