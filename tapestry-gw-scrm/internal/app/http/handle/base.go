package handle

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/pkg/common"
)

type BaseHandle struct {
}

func (b BaseHandle) ResSuccess(c *gin.Context, v interface{}) {
	ginplus.ResSuccess(c, v)
	return
}

func (b BaseHandle) ResError(c *gin.Context, err error) {
	ginplus.ResError(c, err)
	return
}

func (b BaseHandle) ResOk(c *gin.Context) {
	ginplus.ResOK(c)
	return
}

func (b BaseHandle) ResJson(c *gin.Context, status int, v interface{}) {
	ginplus.ResJSON(c, status, v)
	return
}

func (b BaseHandle) WarpContext(c *gin.Context) ctx.Context {
	return ctx.Wrap(c)
}

func (b BaseHandle) GetTokenInfo(c *gin.Context) *common.TokenInfo {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	return tokenInfo
}
