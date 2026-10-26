package open_http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/pkg/common"
	scrm_errs "gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

func openOrder(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	if tokenInfo == nil {
		ginplus.ResError(c, errorx.New("令牌无效,请重新登录"))
		return
	}
	var request = make(map[string]interface{})
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	request["enterprise_id"] = tokenInfo.EnterpriseId
	res, rErr := microscrm.SaveOpenOrder(ctx.Wrap(c), request)
	if rErr != nil {
		ginplus.ResError(c, rErr)
		return
	}
	ginplus.ResSuccess(c, res)
}
