package settings_http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-scrm/internal/config"
	"gw-scrm/internal/modules/enterprise"
	"gw-scrm/internal/modules/settings"
	"gw-scrm/internal/modules/settings/api"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/sdk/microboss"
	"strings"
)

// getSettings 企微设置信息
func getSettings(c *gin.Context) {
	var requestInfo settings.AuthRequest
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	requestInfo.EnterpriseID = tokenInfo.EnterpriseId
	requestInfo.AdminID = tokenInfo.AdminId
	act := &api.ActiveApiImpl{}
	var settingResp settings.SettingResponse
	//外部联系人接口调用凭证
	authRes, aerr := act.AuthContact(requestInfo)
	if aerr != nil {
		ginplus.ResError(c, aerr)
		return
	}
	// 自建应用信息
	pluginRes, perr := act.WorkPlugin(requestInfo)
	if perr != nil {
		ginplus.ResError(c, perr)
		return
	}
	// 客户联系信息
	callbackRes, cerr := act.Callback(requestInfo)
	if cerr != nil {
		ginplus.ResError(c, cerr)
		return
	}

	// 查询企业信息
	enterpriseData := enterprise.Info(ctx.Wrap(c), tokenInfo.EnterpriseId)
	enterpriseDataRes := enterpriseData.(microboss.Enterprise)
	EnterpriseHashValue := enterpriseDataRes.Data.Hash

	// 客户资料
	CustomerPortraitUrl := config.CfgData.Service.WebH5.Domain + "/work/" + EnterpriseHashValue + "/customer-portrait"
	// 素材库
	MaterialLibraryUrl := config.CfgData.Service.WebH5.Domain + "/work/" + EnterpriseHashValue + "/material-library"
	// 营销推广
	MarketingPromotionUrl := config.CfgData.Service.WebH5.Domain + "/work/" + EnterpriseHashValue + "/marketing-promotion"
	// 裂变活动
	FissionActivityUrl := config.CfgData.Service.WebH5.Domain + "/work/" + EnterpriseHashValue + "/fission-activities"

	settingResp.ExternalSecret = authRes.ExternalSecret
	settingResp.AgentID = pluginRes.AgentId
	// settingResp.Domain = strings.TrimPrefix(pluginRes.Domain, "https://")
	settingResp.Domain = strings.TrimPrefix(config.CfgData.Service.WebH5.Domain, "https://")
	settingResp.Secret = pluginRes.Secret
	settingResp.CheckFile = pluginRes.CheckFile
	settingResp.FileUrl = pluginRes.FileUrl
	settingResp.RedPacketsUrl = pluginRes.RedPacketsUrl
	// settingResp.CustomerInformationUrl = pluginRes.CustomerInformationUrl
	settingResp.CustomerInformationUrl = CustomerPortraitUrl
	// settingResp.MaterialLibraryUrl = pluginRes.MaterialLibraryUrl
	settingResp.MaterialLibraryUrl = MaterialLibraryUrl
	// settingResp.MarketingPromotionUrl = pluginRes.MarketingPromotionUrl
	settingResp.MarketingPromotionUrl = MarketingPromotionUrl
	// 裂变活动
	settingResp.FissionActivityUrl = FissionActivityUrl
	//settingResp.UcCallbackUrl = callbackRes.UcCallbackUrl
	settingResp.Token = callbackRes.Token
	settingResp.Key = callbackRes.Key
	settingResp.UcCallbackUrl = config.CfgData.Service.UcCallBack.Domain + "/app/index?plIdentifier=0&corpid=" + enterpriseDataRes.Data.CorpId
	settingResp.IsSuper = callbackRes.IsSuper
	ginplus.ResSuccess(c, settingResp)
}

// authContactUpdate 外部联系人调用凭证更新
func authContactUpdate(c *gin.Context) {
	var requestInfo settings.ContactRequest
	err := c.BindJSON(&requestInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	requestInfo.EnterpriseID = tokenInfo.EnterpriseId
	act := &api.ActiveApiImpl{}
	// 外部联系人接口调用凭证更新
	res := act.AuthContactUpdate(requestInfo)
	ginplus.ResSuccess(c, res)
}

func setWxqyToken(c *gin.Context) {
	var requestInfo settings.WxqyTokenRequest
	err := c.BindJSON(&requestInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	requestInfo.EnterpriseID = tokenInfo.EnterpriseId
	act := &api.ActiveApiImpl{}
	res := act.SetWxqyToken(requestInfo)
	ginplus.ResSuccess(c, res)
}

// workPluginUpdate 自建应用信息更新
func workPluginUpdate(c *gin.Context) {
	var requestInfo settings.PluginRequest
	err := c.BindJSON(&requestInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	requestInfo.EnterpriseID = tokenInfo.EnterpriseId
	act := &api.ActiveApiImpl{}
	// 自建应用新增
	rerr := act.WorkPluginUpdate(ctx.Wrap(c), requestInfo)
	if rerr != nil {
		ginplus.ResError(c, rerr)
		return
	}
	ginplus.ResSuccess(c, nil)
}

// uploadCheckFile 上传文件
func uploadCheckFile(c *gin.Context) {
	_, err := c.FormFile("file")
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	path := config.CfgData.Path.Upload
	upload := common.UploadCfg{}
	_, err = upload.UploadFile(c, path)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("上传文件错误: err:%s, path:%s", err.Error(), path)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
}

// sessionArchivingUpdate 会话存档自建应用信息更新
func sessionArchivingUpdate(c *gin.Context) {
	var requestInfo settings.PluginRequest
	err := c.BindJSON(&requestInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	requestInfo.EnterpriseID = tokenInfo.EnterpriseId
	act := &api.ActiveApiImpl{}
	// 自建应用新增
	rerr := act.SessionArchivingUpdate(ctx.Wrap(c), requestInfo)
	if rerr != nil {
		ginplus.ResError(c, rerr)
		return
	}
	rerr = act.SettingsUpdate(ctx.Wrap(c), requestInfo)
	if rerr != nil {
		ginplus.ResError(c, rerr)
		return
	}

	ginplus.ResSuccess(c, nil)
}

// getChatSettings 企微设置信息
func getChatSettings(c *gin.Context) {
	request := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.ActiveApiImpl{}
	res := act.SettingsInfo(ctx.Wrap(c), request)
	ginplus.ResSuccess(c, res)

}
