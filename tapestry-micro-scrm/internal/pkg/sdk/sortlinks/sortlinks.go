package sortlinks

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

func GetSortLinkInfo(c ctx.Context, sortLinkUrl string) (*SortLinkInfo, error) {
	var response SortLinkInfo
	err := sendGet(c, sortLinkUrl, &response)
	return &response, err
}
