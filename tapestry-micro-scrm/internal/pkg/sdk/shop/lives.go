package shop

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/common"
)

// LiveList  根据ID获取直播列表
func LiveList(c ctx.Context, request *LiveListRequest) (*LiveListResponse, error) {
	var response LiveListResponse
	enterpriseHash := c.Value(common.ContextKeyEnum.EnterpriseHash)
	err := httpClient.Post(c, GetURL(urlEnum.LiveListUrl), request, &response, clients.Header(common.HeaderKeyEnum.EnterpriseHash, enterpriseHash.(string)))
	if err != nil {
		return nil, err
	}
	return &response, nil
}
