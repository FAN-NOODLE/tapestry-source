package microactivities

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// Activities 获取活动列表
func Activities(c ctx.Context, request *ActivitiesByIdsRequest) (*ActivitiesByIdsResponse, error) {
	var response ActivitiesByIdsResponse
	err := httpClient.Post(c, GetURL(urlEnum.ActivityListByIdsUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Activities 获取活动列表
func ActivityPages(c ctx.Context, request *ActivityPagesRequest) (*ActivityPagesResponse, error) {
	var response ActivityPagesResponse
	err := httpClient.Post(c, GetURL(urlEnum.ActivityListByIdsUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
