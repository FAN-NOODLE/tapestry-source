package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// SettingsApi 会话设置 SDK
type SettingsApi struct {
	BaseSDK
}

func NewSettings(host string) *SettingsApi {
	urlHost = host
	return &SettingsApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// SettingsInfo 获取会话设置信息
func (api *SettingsApi) SettingsInfo(c ctx.Context, settings *SettingsInfoRequest) (*ChatSettings, error) {
	path := "/micro-scrm/v1/settings/info"
	url := urlHost + path
	res, err := api.Post(c, url, settings)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &ConfigInfoResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	return result.Data.(*ChatSettings), err
}

// SettingsUpdate 保存会话设置信息
func (api *SettingsApi) SettingsUpdate(c ctx.Context, settings *ChatSettings) error {
	path := "/micro-scrm/v1/settings/update"
	url := urlHost + path
	res, err := api.Post(c, url, settings)
	if err != nil {
		return err
	}
	result := &struct{}{}
	_, err = res.Read(clients.JSON(result))
	return err
}
