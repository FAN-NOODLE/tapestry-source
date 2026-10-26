package wxcouponhttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	wxcoupon "gw-scrm/internal/modules/wxCoupon"
	"gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/http"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

//wxCouponAvailableList 可用的微信优惠券列表
func wxCouponAvailableList(c *gin.Context) {
	var request wxcoupon.AvailableWxCouponRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microscrm.WxCouponAvailableList(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}
