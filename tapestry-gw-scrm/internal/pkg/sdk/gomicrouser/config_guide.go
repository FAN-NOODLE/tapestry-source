package gomicrouser

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// ConfigGuideList 配置引导列表
func ConfigGuideList(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ConfigGuideListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
