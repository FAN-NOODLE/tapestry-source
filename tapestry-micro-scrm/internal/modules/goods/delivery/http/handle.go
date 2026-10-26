package http

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/goods"
	"micro-scrm/internal/modules/marketingPromotion"
	response "micro-scrm/internal/pkg/resoponse"
	"micro-scrm/internal/pkg/sdk/boss"
	"net/http"
	"strconv"
)

var goodsService *goods.Service

// skuList 商品列表
func skuList(c *gin.Context) {
	req := new(goods.ListReq)
	if err := c.BindJSON(req); err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 10
	}

	// 查询正常状态
	req.Status = goods.SkuNormalStatus

	// 查询列表
	total, list, err := goodsService.SkuList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}

	// 获取小程序的logo+name
	appReq := &boss.AuthDetailRequest{
		EnterpriseID: req.EnterpriseID,
		AuthType:     1,
	}
	authDetailResponse, _ := boss.AuthDetail(ctx.Wrap(c), appReq)

	for key, item := range list {
		err := item.SetAttributes()
		if err != nil {
			ginplus.ResError(c, err, http.StatusOK)
			return
		}
		(*list[key]).ShareAppId = authDetailResponse.Data.AppID
		(*list[key]).SharePath = getParams(ctx.Wrap(c), *req, *item)
	}

	result := goods.SpuListRes{
		List:  list,
		Page:  req.Page,
		Size:  req.Size,
		Total: total,
	}

	ginplus.ResSuccess(c, result)
}

// 组装携带参数
func getParams(c ctx.Context, req goods.ListReq, value goods.SkuItem) string {
	params := "/" + marketingPromotion.MarketingPath[marketingPromotion.MarketingGoods] + ".html"
	params += "?spu_id=" + strconv.FormatInt(value.SpuID, 10) + "&id=" + strconv.FormatInt(value.SpuID, 10) + "&sku_id=" + strconv.FormatInt(value.SkuID, 10) + "&store_id=0&mem_uid=" + req.MemberUserId + "&share_uid=" + req.MemberUserId + "&source=scrm"
	return params
}

func syncGoods(c *gin.Context) {
	var request goods.SyncRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	err = goodsService.SyncGoodsData(ctx.Wrap(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
}

// goodSpuList 商品列表
func goodSpuList(c *gin.Context) {
	req := &goods.SpuListRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	total, list, lerr := goodsService.SpuList(ctx.Wrap(c), req)
	if lerr != nil && lerr != sql.ErrNoRows {
		ginplus.ResError(c, lerr)
		return
	}

	response.SuccessList(c, total, req.Page, req.Size, list)
}
