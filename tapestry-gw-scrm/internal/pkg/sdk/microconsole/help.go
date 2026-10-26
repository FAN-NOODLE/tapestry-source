package microconsole

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/modules/index"
)

// GetAuthContact 获取外部联系人凭证
func HelpList(c ctx.Context, request interface{}) (*index.HelpListRes, error) {
	var response index.HelpListRes
	err := httpClient.Post(c, GetUrl(urlEnum.HelpListUri), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
