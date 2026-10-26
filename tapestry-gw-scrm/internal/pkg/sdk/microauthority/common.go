package microauthority

import (
	"fmt"
	"gw-scrm/internal/config"
)

// 获取cdp请求链接
func GetUrl(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.MicroAuthority.Domain)
}
