package gomicrouser

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// Save 权益保存
func Save(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.RightSave), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// List 权益列表
func List(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.RightList), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// UpdateRightStatus 更新权益状态
func UpdateRightStatus(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.UpdateRightStatusUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
