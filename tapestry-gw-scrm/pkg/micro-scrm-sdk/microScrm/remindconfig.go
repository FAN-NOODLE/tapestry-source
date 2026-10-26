package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// RemindConfigApi 提醒设置 SDK
type RemindconfigApi struct {
	BaseSDK
}

func NewRemindconfig(host string) *RemindconfigApi {
	urlHost = host
	return &RemindconfigApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// ConfigInfo 获取提醒设置信息
func (api *RemindconfigApi) ConfigInfo(c ctx.Context, remindconfig *ConfigInfoRequest) (*ConfigInfoResponse, error) {
	path := "/micro-scrm/v1/remind-config/config-info"
	url := urlHost + path
	res, err := api.Post(c, url, remindconfig)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &ConfigInfoResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	return result.Data.(*ConfigInfoResponse), err
}

// ConfigSave 保存提醒设置信息
func (api *RemindconfigApi) ConfigSave(c ctx.Context, remindconfig *ConfigSaveRequest) error {
	path := "/micro-scrm/v1/remind-config/config-save"
	url := urlHost + path
	res, err := api.Post(c, url, remindconfig)
	if err != nil {
		return err
	}
	result := &struct{}{}
	_, err = res.Read(clients.JSON(result))
	return err
}
