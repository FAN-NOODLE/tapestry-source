package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// ConfigInfo 获取提醒设置信息
func ConfigInfo(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ConfigInfoUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// ConfigSave 保存获取提醒设置信息
func ConfigSave(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.ConfigSaveUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
