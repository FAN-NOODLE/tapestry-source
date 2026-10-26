package remindconfig_http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/remindconfig/api"
	"gw-scrm/internal/pkg/common"
)

// configInfo 获取超时提醒设置
func configInfo(c *gin.Context) {
	request := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	request["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.RemindconfigApi{}
	res := act.ConfigInfo(ctx.Wrap(c), request)
	ginplus.ResSuccess(c, res)
}

// configSave 保存超时提醒设置
func configSave(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	act := &api.RemindconfigApi{}
	res, rErr := act.ConfigSave(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)

	return
}
