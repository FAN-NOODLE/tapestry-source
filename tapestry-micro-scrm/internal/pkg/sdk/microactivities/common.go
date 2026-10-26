package microactivities

import (
	"fmt"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	ActivityListByIdsUrl string
	ActivityPagesUrl     string
}{
	//活动列表
	ActivityListByIdsUrl: "%s/micro-activity/v1/common/activity-list-by-ids",
	//活动页列表
	ActivityPagesUrl: "%s/micro-activity/v1/common/activity-list-by-ids",
}

// ResponseBase 服务接口返回数据结构体头
type ResponseBase struct {
	sdk.ResponseBase
}

// httpClient http请求代理
var httpClient = sdk.NewHTTPClient(sdk.ParseJSON{})

// GetURL 获取shop请求链接
func GetURL(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.MicroActivity.Domain)
}
