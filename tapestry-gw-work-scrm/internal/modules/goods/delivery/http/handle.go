package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-work-scrm/internal/pkg/errs"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

// goodsSdk goods 请求API
var goodsSdk *microScrm.GoodsApi

// 商品列表
func goodsSkuList(c *gin.Context) {

	enterpriseId, ok := c.Value("enterpriseId").(int64)
	if !ok {
		ginplus.ResError(c, errs.TokenEnterpriseIdErr)
		return
	}

	request := new(microScrm.GoodsListReq)
	if err := c.BindJSON(request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	memberInfo := c.Value("memberInfo").(*microScrm.Member)
	request.MemberUserId = memberInfo.Userid
	request.EnterpriseID = enterpriseId
	result, err := goodsSdk.GoodsSkuList(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, result)
}

// goodsSpuList 商品列表
func goodsSpuList(c *gin.Context) {

	enterpriseId, ok := c.Value("enterpriseId").(int64)
	if !ok {
		ginplus.ResError(c, errs.TokenEnterpriseIdErr)
		return
	}

	request := new(microScrm.GoodsSpuListReq)
	request.SpuIds = make([]int64, 0)
	if err := c.BindJSON(request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	request.EnterpriseID = enterpriseId
	result, err := goodsSdk.GoodsSpuList(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, result)
}
