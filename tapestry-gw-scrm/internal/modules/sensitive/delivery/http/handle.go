package sensitivehttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/sensitive/api"
	"gw-scrm/internal/pkg/common"
)

// sensitiveList 敏感词列表
func sensitiveList(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	act := &api.SensitiveApi{}
	res, rErr := act.SensitiveList(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)

	return
}

// sensitiveSave 敏感词保存
func sensitiveSave(c *gin.Context) {
	params := common.GetRawMap(c)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId

	act := &api.SensitiveApi{}
	res, rErr := act.SensitiveSave(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)

	return
}

// sensitiveDelete 敏感词删除
func sensitiveDelete(c *gin.Context) {
	params := common.GetRawMap(c)

	act := &api.SensitiveApi{}
	res, rErr := act.SensitiveDelete(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)

	return
}

// sensitiveDetail 敏感词详情
func sensitiveDetail(c *gin.Context) {
	params := common.GetRawMap(c)

	act := &api.SensitiveApi{}
	res, rErr := act.SensitiveDetail(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)
		return
	}
	ginplus.ResSuccess(c, res)

	return
}

// sensitiveDetailList 敏感词明细列表获取
func sensitiveDetailList(c *gin.Context) {
	params := common.GetRawMap(c)

	act := &api.SensitiveApi{}
	res, rErr := act.SensitiveDetailList(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}
