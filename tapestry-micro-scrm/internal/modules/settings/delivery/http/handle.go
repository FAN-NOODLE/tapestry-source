package settingshttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/config"
	"micro-scrm/internal/modules/settings"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	response "micro-scrm/internal/pkg/resoponse"
)

var service *settings.Service

// settingsInfo 获取会话存档配置
func settingsInfo(c *gin.Context) {
	request := &settings.SettingsInfoRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	settingsInfo, err := service.SettingsInfo(ctx.Wrap(c), request)
	if errorIgnore.Ignore(err) != nil {
		ginplus.ResError(c, errs.SettingsInfoErr)
		return
	}
	settingsInfo.TrustedIps = config.CfgData.TrustedIps
	settingsInfo.MessagePrivateKey = ""
	response.Success(c, settingsInfo)
	return
}

// settingsUpdate 保存会话存档设置
func settingsUpdate(c *gin.Context) {
	request := &settings.ChatSettings{}
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	_, UpErr := service.SettingsSave(ctx.Wrap(c), request)
	if UpErr != nil {
		ginplus.ResError(c, UpErr)
		return
	}
	ginplus.ResSuccess(c, []struct{}{})
	return
}
