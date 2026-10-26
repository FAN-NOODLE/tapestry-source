package gomicrouser

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// LevelDirectList 等级直达列表
func LevelDirectList(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.LevelDirectListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// LevelDirectAdd 等级直达添加
func LevelDirectAdd(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.LevelDirectAddUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// LevelDirectUpdate 等级直达修改
func LevelDirectUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.LevelDirectUpdateUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// LevelDirectDelete 等级直达删除
func LevelDirectDelete(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.LevelDirectDeleteUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// LevelDirectDetail 等级直达详情
func LevelDirectDetail(c ctx.Context, request interface{}) (interface{}, error) {
	var response GoUserResponse
	err := httpClient.Post(c, GetUrl(urlEnum.LevelDirectDetailUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
