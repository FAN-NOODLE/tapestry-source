package microScrm

import (
	"fmt"
	"gw-work-scrm/internal/config"
)

// GetUrl 获取cdp请求链接
func GetUrl(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.MicroScrm.Domain)
}

