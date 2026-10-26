package shorturl

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/modules/marketing"
	"gw-scrm/internal/pkg/sdk/cdp"
)

// 获取短链
func GetShortUrlQuery(c ctx.Context, request marketing.ShortUrlRequest) (interface{}, error) {
	resp, err := cdp.SendPost(c, GetUrl(urlEnum.ShortUrl), request)
	return resp, err

}

//添加短链
func AddShortUrlQuery(c ctx.Context, request marketing.ShortAddReq) (interface{}, error) {
	resp, err := cdp.SendPost(c, GetUrl(urlEnum.ShortUrlAdd), request)
	return resp, err

}
