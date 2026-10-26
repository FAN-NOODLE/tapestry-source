package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/version"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

func getCurrentVersionInfo(c *gin.Context) {
	var param version.GetCurrentVersionInfoRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		ginplus.ResError(c, err)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	param.EnterpriseID = tokenInfo.EnterpriseId
	versionInfo, err := microscrm.GetCurrentVersionInfo(ctx.Wrap(c), param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}
