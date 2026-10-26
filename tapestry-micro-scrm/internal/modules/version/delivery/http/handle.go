package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/version"
)

func getCurrentVersionInfo(c *gin.Context) {
	var param version.GetCurrentVersionInfoRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		ginplus.ResError(c, err)
		return
	}

	// 获取企业ID, 不检测
	//enterpriseID, _ := c.Get(common.ContextKeyEnum.EnterpriseID)
	//param.EnterpriseID = enterpriseID.(int64)

	result, resultErr := version.GetCurrentVersionInfo(ctx.Wrap(c), param)
	if resultErr != nil {
		ginplus.ResError(c, resultErr)
		return
	}

	ginplus.ResSuccess(c, result)
}
