package microactivity

import (
	"fmt"
	"gw-scrm/internal/config"
)

var urlEnum = struct {
	ActivityLists string
}{
	// 获取活动列表
	ActivityLists: "%s/micro-activity/v1/common/activity-list",
}

// 获取活动请求链接
func GetUrl(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.MicroActivity.Domain)
}
