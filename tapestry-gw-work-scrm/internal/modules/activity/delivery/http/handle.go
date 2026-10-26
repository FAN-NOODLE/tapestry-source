package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-work-scrm/internal/pkg/errs"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

// activitySdk 请求API
var activitySdk *microScrm.ActivityApi

// 商品列表
func goodsList(c *gin.Context) {

	enterpriseId, ok := c.Value("enterpriseId").(int64)
	if !ok {
		ginplus.ResError(c, errs.TokenEnterpriseIdErr)
		return
	}

	request := new(microScrm.ActivityGoodsListRequest)
	if err := c.BindJSON(request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	memberInfo := c.Value("memberInfo").(*microScrm.Member)
	request.MemberUserId = memberInfo.Userid
	request.EnterpriseID = enterpriseId
	result, err := activitySdk.ActivityGoodsList(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, result)
}
