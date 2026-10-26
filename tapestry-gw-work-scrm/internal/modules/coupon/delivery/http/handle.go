package couponhttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-work-scrm/internal/pkg/errs"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"net/http"
)

var couponSdk *microScrm.CouponApi

// groupChatInfo 获取社群侧边栏基础信息
func wxCouponList(c *gin.Context) {
	var req microScrm.WxCouponRequest
	jsonError := c.BindJSON(&req)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	memberInfo := c.Value("memberInfo").(*microScrm.Member)
	req.EnterpriseID = c.Value("enterpriseId").(int64)
	req.MemberUserId = memberInfo.Userid

	data, err := couponSdk.WxCouponList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, data)
}

// couponGoods 查询优惠券关联的商品
func couponGoods(c *gin.Context) {
	request := new(microScrm.WithGoodsRequest)
	if err := c.BindJSON(request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	enterpriseId, ok := c.Value("enterpriseId").(int64)
	if !ok {
		ginplus.ResError(c, errs.TokenEnterpriseIdErr)
		return
	}

	request.EnterpriseID = enterpriseId
	data, err := couponSdk.CouponGoods(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, data)
}
