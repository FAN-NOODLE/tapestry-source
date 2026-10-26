package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// MemberApi 员工 SDK
type CouponApi struct {
	BaseSDK
}

// NewCoupon 返回API实例
func NewCoupon(host string) *CouponApi {
	urlHost = host
	return &CouponApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// GroupChatInfo 获取社群侧边栏基础信息
func (s CouponApi) WxCouponList(c ctx.Context, request WxCouponRequest) (*WxCouponListResponse, error) {
	path := "/micro-scrm/v1/coupon/wx-coupon-list"
	url := urlHost + path
	res, err := s.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &WxCouponListResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*WxCouponListResponse), err
}

// couponGoods 优惠券选中商品列表
func (s CouponApi) CouponGoods(c ctx.Context, request *WithGoodsRequest) (*WithGoodsResponse, error) {
	path := "/micro-scrm/v1/coupon/goods-list"
	url := urlHost + path
	res, err := s.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &WithGoodsResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*WithGoodsResponse), err
}
