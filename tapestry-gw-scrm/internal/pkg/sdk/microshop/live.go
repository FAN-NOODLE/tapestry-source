package microshop

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

//LivesRequest 直播列表
type LivesRequest struct {
	EnterpriseID int64      `json:"enterprise_id"`
	Page         int64      `json:"page"`
	Size         int64      `json:"size"`
	Search       SearchLive `json:"search"`
}

//SearchLive 直播搜索条件
type SearchLive struct {
	ID        string `json:"id"`
	LiveTitle string `json:"live_title"`
}

// 直播列表
func LiveList(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.LiveListUrl), request)
	return response, err
}
