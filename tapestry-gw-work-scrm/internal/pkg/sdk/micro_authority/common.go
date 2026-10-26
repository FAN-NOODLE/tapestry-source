package microAuthority

import (
	"fmt"
	"gw-work-scrm/internal/config"
	"gw-work-scrm/internal/pkg/sdk"
)

// ResponseBase 服务接口返回数据结构体头
type ResponseBase struct {
	sdk.ResponseBase
}

// GetURL 获取shop请求链接
func GetURL(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.MicroAuthority.Domain)
}
