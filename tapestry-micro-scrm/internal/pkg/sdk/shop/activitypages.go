package shop

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/common"
)

// ActivityPageList 根据id获取活动页列表
func ActivityPageList(c ctx.Context, request *ActivityPageListRequest) (*ActivityPageListResponse, error) {
	var response ActivityPageListResponse
	enterpriseHash := c.Value(common.ContextKeyEnum.EnterpriseHash)
	err := httpClient.Post(c, GetURL(urlEnum.ActivityPageListUrl), request, &response, clients.Header(common.HeaderKeyEnum.EnterpriseHash, enterpriseHash.(string)))
	if err != nil {
		return nil, err
	}
	return &response, nil
}
