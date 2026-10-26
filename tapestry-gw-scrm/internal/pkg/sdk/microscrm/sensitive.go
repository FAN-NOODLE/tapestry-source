package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// SensitiveList 敏感词列表
func SensitiveList(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.SensitiveListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// SensitiveSave 敏感词保存
func SensitiveSave(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.SensitiveSaveUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// SensitiveDelete 敏感词删除
func SensitiveDelete(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.SensitiveDeleteUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// SensitiveDetail 敏感词详情
func SensitiveDetail(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.SensitiveDetailUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// SensitiveDetailList 敏感词明细列表
func SensitiveDetailList(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.SensitiveDetailListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
