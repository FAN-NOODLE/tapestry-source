package couponHttp

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/coupon"
	"micro-scrm/internal/modules/goods"
	"micro-scrm/internal/modules/user"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errs"
	response "micro-scrm/internal/pkg/resoponse"
)

var couponService *coupon.Service
var userService *user.Service
var goodService *goods.Service

// SyncCoupon 同步优惠券相关接口
func SyncCoupon(c *gin.Context) {
	method := c.GetHeader(common.DDCMethodHeader)
	var err error
	switch method {
	case common.DdcEventEnum.WxCouponAdd.String(),
		common.DdcEventEnum.WxCouponUpdate.String(),
		common.DdcEventEnum.WxCouponDelete.String():
		var request coupon.SyncCouponRequest
		err = c.BindJSON(&request)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		_, err = couponService.SyncCoupon(ctx.Wrap(c), &request.CouponInfo, common.DdcEvent(request.Method))
	case common.DdcEventEnum.CouponUserAdd.String(),
		common.DdcEventEnum.CouponUserUpdate.String():
		var request coupon.SyncCouponUserRequest
		err = c.BindJSON(&request)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		_, err = couponService.SyncCouponUser(ctx.Wrap(c), &request.DataList.CouponUserInfoList, common.DdcEvent(request.Method))

	default:
		ginplus.ResError(c, errs.MethodHeaderErr)
		return
	}

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// customerCoupon 用户优惠券
func customerCoupon(c *gin.Context) {
	req := &coupon.UserCouponListRequest{}
	if err := c.ShouldBind(req); err != nil {
		ginplus.ResError(c, err)
		return
	}

	total, list, err := couponService.CouponUserList(ctx.Wrap(c), req)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	response.SuccessList(c, total, req.Page, req.Size, list)
}

// wxCouponList 微信优惠券列表
func wxCouponList(c *gin.Context) {
	req := coupon.WxCouponRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	cmp, err := couponService.WxCouponList(ctx.Wrap(c), req)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	response.SuccessList(c, cmp.Total, int(cmp.Page), int(cmp.Size), cmp.List)
}

// couponGoods 查询选中商品
func couponGoods(c *gin.Context) {
	req := &coupon.WithGoodsRequest{}
	if err := c.ShouldBind(req); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 10
	}

	// 查询一个活动
	cx := ctx.Wrap(c)
	info, err := couponService.CouponInfo(cx, req)
	if err != nil || info == nil {
		ginplus.ResError(c, err)
		return
	}

	spuLists := make([]string, 0, 10)
	err = json.Unmarshal([]byte(info.Range), &spuLists)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	// 无商品情况
	if len(spuLists) == 0 {
		ginplus.ResSuccess(c, coupon.WithGoodsResponse{
			List:  []coupon.WithGoodsSpuItem{},
			Page:  req.Page,
			Size:  req.Size,
			Total: 0,
		})
		return
	}

	// 转化类型为INT64
	spuIntLists := make([]int64, 0, len(spuLists))
	for _, tmp := range spuLists {
		sku, err := strconv.ParseInt(tmp, 10, 64)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		spuIntLists = append(spuIntLists, sku)
	}

	// 查询SpU列表
	goodsReq := new(goods.SpuListRequest)
	goodsReq.Size = req.Size
	goodsReq.Page = req.Page
	goodsReq.SpuIDs = spuIntLists
	goodsReq.EnterpriseID = req.EnterpriseID
	goodsReq.Status = goods.SPuNormalStatus
	total, list, err := goodService.SpuList(cx, goodsReq)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	// 未查询到商品情况
	if list == nil || len(list) == 0 {
		ginplus.ResSuccess(c, coupon.WithGoodsResponse{
			List:  []coupon.WithGoodsSpuItem{},
			Page:  req.Page,
			Size:  req.Size,
			Total: 0,
		})
		return
	}

	// 进行数据过滤
	data := make([]coupon.WithGoodsSpuItem, 0, len(list))

	for _, item := range list {
		data = append(data, coupon.WithGoodsSpuItem{
			SpuID:   item.SpuID,
			SpuName: item.SpuName,
			Image:   item.Picture,
		})
	}

	result := coupon.WithGoodsResponse{
		List:  data,
		Page:  req.Page,
		Size:  req.Size,
		Total: total,
	}

	ginplus.ResSuccess(c, result)
}

// wxCouponAvailableList 可用的微信优惠券列表
func wxCouponAvailableList(c *gin.Context) {
	req := coupon.WxCouponAvailableRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	req.NowTime = time.Now().Unix()
	cmp, err := couponService.GetAllWxCouponAvailableList(ctx.Wrap(c), req)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	response.SuccessList(c, cmp.Total, int(cmp.Page), int(cmp.Size), cmp.List)
}
