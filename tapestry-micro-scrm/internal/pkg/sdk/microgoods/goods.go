package microgoods

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

// SpuList 获取spu列表
func SpuList(c ctx.Context, request *SpuListRequest) (*SpuListResponse, error) {
	var response SpuListResponse
	err := httpClient.Post(c, GetURL(urlEnum.SpuList), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
