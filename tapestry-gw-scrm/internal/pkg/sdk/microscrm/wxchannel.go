package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

//WxChannelList 微信场景值列表
func WxChannelList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.WxChannelUrl), request)
}
