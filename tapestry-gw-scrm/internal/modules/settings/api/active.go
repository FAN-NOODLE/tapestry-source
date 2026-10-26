package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/modules/settings"
	"gw-scrm/internal/pkg/sdk/microboss"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

// PluginIsAdd 自建应用添加
const PluginIsAdd = 1

type ActiveApiImpl struct{}

// AuthContactDataResponse 外部联系人调用凭证信息
type AuthContactDataResponse struct {
	ExternalSecret string `json:"external_secret""`
}

// AuthContact 外部联系人调用凭证
func (a ActiveApiImpl) AuthContact(req settings.AuthRequest) (*AuthContactDataResponse, error) {
	res, err := microboss.GetAuthContact(ctx.New(), req)
	if err != nil {
		return nil, err
	}
	var response = AuthContactDataResponse{
		ExternalSecret: res.Data.ExternalSecret,
	}
	return &response, nil
}

// WorkPluginResponse 自建应用信息
type WorkPluginDataResponse struct {
	AgentId                int64  `json:"agent_id"`
	Secret                 string `json:"secret"`
	Domain                 string `json:"domain"`
	RedPacketsUrl          string `json:"red_packets_url"`
	CustomerInformationUrl string `json:"customer_information_url"`
	MaterialLibraryUrl     string `json:"material_library_url"`
	MarketingPromotionUrl  string `json:"marketing_promotion_url"`
	CheckFile              int8   `json:"check_file"`
	FileUrl                string `json:"file_url"`
}

// WorkPlugin 自建应用信息
func (a ActiveApiImpl) WorkPlugin(req settings.AuthRequest) (*WorkPluginDataResponse, error) {
	res, err := microboss.GetWorkPlugin(ctx.New(), req)
	if err != nil {
		return nil, err
	}
	var response = WorkPluginDataResponse{
		AgentId:                res.Data.AgentId,
		Secret:                 res.Data.Secret,
		Domain:                 res.Data.Domain,
		RedPacketsUrl:          res.Data.RedPacketsUrl,
		CustomerInformationUrl: res.Data.CustomerInformationUrl,
		MaterialLibraryUrl:     res.Data.MaterialLibraryUrl,
		MarketingPromotionUrl:  res.Data.MarketingPromotionUrl,
		CheckFile:              res.Data.CheckFile,
		FileUrl:                res.Data.FileUrl,
	}
	return &response, nil
}

// CallBackDataResponse 客户联系回调配置
type CallBackDataResponse struct {
	Key           string `json:"key"`
	Token         string `json:"token"`
	UcCallbackUrl string `json:"uc_callback_url"`
	IsSuper       int    `json:"is_super""`
}

// Callback 客户联系配置
func (a ActiveApiImpl) Callback(req settings.AuthRequest) (*CallBackDataResponse, error) {
	res, err := microboss.GetCallBack(ctx.New(), req)
	if err != nil {
		return nil, err
	}
	var response = CallBackDataResponse{
		Key:           res.Data.Key,
		Token:         res.Data.Token,
		UcCallbackUrl: res.Data.UcCallbackUrl,
		IsSuper:       res.Data.IsSuper,
	}
	return &response, nil
}

// AuthContactUpdate 调用凭证更新
func (a ActiveApiImpl) AuthContactUpdate(req settings.ContactRequest) error {
	return microboss.AuthContactUpdate(ctx.New(), req)
}

// SetWxqyToken 客户联系修改token
func (a ActiveApiImpl) SetWxqyToken(req settings.WxqyTokenRequest) error {
	return microboss.SetWxqyToken(ctx.New(), req)
}

// WorkPluginUpdate 自建应用信息更新
func (a ActiveApiImpl) WorkPluginUpdate(c ctx.Context, req settings.PluginRequest) error {
	if req.Type == PluginIsAdd {
		return microboss.WorkPluginAdd(c, req)
	} else {
		return microboss.WorkPluginUpdate(c, req)
	}
}

// SessionArchivingUpdate 会话存档自建应用信息更新
func (a ActiveApiImpl) SessionArchivingUpdate(c ctx.Context, req settings.PluginRequest) error {
	if req.Type == PluginIsAdd {
		return microboss.SessionArchivingAdd(c, req)
	} else {
		return microboss.SessionArchivingUpdate(c, req)
	}
}

// SettingsUpdate 消息风控secret保存
func (a ActiveApiImpl) SettingsUpdate(c ctx.Context, req settings.PluginRequest) error {
	_, err := microscrm.SettingsUpdate(c, req)
	if err != nil {
		return err
	}
	return nil
}

// SettingsInfo 获取会话存档信息
func (a ActiveApiImpl) SettingsInfo(c ctx.Context, req interface{}) interface{} {
	res, err := microscrm.SettingsInfo(c, req)
	if err != nil {
		return nil
	}
	return res
}
