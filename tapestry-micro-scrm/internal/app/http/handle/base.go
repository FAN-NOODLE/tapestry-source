package handle

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/app"
	"micro-scrm/internal/pkg/logger"
	"micro-scrm/internal/pkg/resoponse"
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

func (b BaseHandle) ResList(c *gin.Context, total int64, page, size int, list interface{}) {
	response.SuccessList(c, total, page, size, list)
	return
}

func (b BaseHandle) Log() *logger.Logger {
	return app.Logger()
}

func (b BaseHandle) WarpContext(c *gin.Context) ctx.Context {
	return ctx.Wrap(c)
}
