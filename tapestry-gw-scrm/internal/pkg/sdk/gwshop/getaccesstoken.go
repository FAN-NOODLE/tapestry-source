package gwshop

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

// 获取accesstoken
func GetAccessTokenQuery(c ctx.Context, request interface{}) (interface{}, error) {
	resp, err := sdk.SendPost(c, GetUrl(urlEnum.GetAccessToken), request)
	return resp, err
}
