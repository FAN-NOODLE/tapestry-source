package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// OvertimeList 超时列表
func OvertimeList(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.OvertimeListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
