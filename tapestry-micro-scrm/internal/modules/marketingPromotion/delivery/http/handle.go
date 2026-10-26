package marketingPromotionhttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/marketingPromotion"
	"micro-scrm/internal/pkg/errs"
)

var marketingPromotionService *marketingPromotion.Service
var enterpriseService *enterprise.Service

// getPoster 获取海报
func getPoster(c *gin.Context) {
	request := &marketingPromotion.GetPosterRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	domain, _, _, _ := enterpriseService.GetInfo(ctx.Wrap(c), request.EnterpriseID)
	if domain == "" {
		ginplus.ResError(c, errs.EnterpriseNotFund)
		return
	}
	response, _ := marketingPromotionService.GetPoster(ctx.Wrap(c), domain, *request)
	if response.Object == "" {
		ginplus.ResError(c, errs.GetPosterErr, errs.ErrorStatus)

		return
	}
	ginplus.ResSuccess(c, response)
	return
}
