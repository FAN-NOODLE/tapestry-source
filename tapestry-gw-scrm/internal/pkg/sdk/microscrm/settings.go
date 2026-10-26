package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// SettingsUpdate 消息风控secret保存
func SettingsUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.SettingsUpdateUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// SettingsInfo 获取会话存档配置信息
func SettingsInfo(c ctx.Context, request interface{}) (interface{}, error) {
	var response ScrmResponse
	err := httpClient.Post(c, GetUrl(urlEnum.SettingsInfoUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
