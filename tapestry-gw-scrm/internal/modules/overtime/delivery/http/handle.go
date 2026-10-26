package overtimehttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/overtime/api"
	"gw-scrm/internal/pkg/common"
)

// overtimeList 超时列表
func overtimeList(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.OvertimeApi{}
	res, rErr := act.OvertimeList(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}
