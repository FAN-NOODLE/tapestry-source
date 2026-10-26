package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/ws_trace"
	"micro-scrm/internal/pkg/common"
)

var wsTraceService *ws_trace.Service

func wsAdd(c *gin.Context) {
	var param ws_trace.WsTrace
	if err := c.ShouldBindJSON(&param); err != nil {
		ginplus.ResError(c, err)
		return
	}
	// 获取企业ID, 不检测
	enterpriseID, _ := c.Get(common.ContextKeyEnum.EnterpriseID)
	result, resultErr := wsTraceService.WsAdd(ctx.Wrap(c), enterpriseID.(int64), param)
	if resultErr != nil {
		ginplus.ResError(c, resultErr)
		return
	}

	ginplus.ResSuccess(c, result)
}
