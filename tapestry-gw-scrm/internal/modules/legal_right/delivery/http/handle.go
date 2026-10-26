package legal_right_http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	api "gw-scrm/internal/modules/legal_right/api"
	"gw-scrm/internal/pkg/common"
)

// save 权益保存
func save(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	params["enterprise_id"] = tokenInfo.EnterpriseId
	params["operation_id"] = tokenInfo.AdminId
	act := &api.LegalRightApi{}
	res, rErr := act.Save(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// List 权益列表
func list(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	params["operation_id"] = tokenInfo.AdminId
	act := &api.LegalRightApi{}
	res, rErr := act.List(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// UpdateStatus 更新权益状态
func UpdateStatus(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.LegalRightApi{}
	res, rErr := act.UpdateRightStatus(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}
