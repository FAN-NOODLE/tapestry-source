// Package callback 对接Open 平台
package coupon

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

//go:generate mockgen -source=${GOFILE} -destination=mocks/mock_${GOFILE}

// Repository 存储定义
type Repository interface {
	repos.IRepository
	AddCoupon(ctx.Context, *BasicInfoEntity) (int64, error)
	UpdateCoupon(ctx.Context, *BasicInfoEntity) (int64, error)
	CountCoupon(ctx.Context, *BasicInfoEntity) int64
	DeleteCoupon(ctx.Context, *BasicInfoEntity) (int64, error)
	AddCouponUser(ctx.Context, *BasicCouponUserInfoEntity) (int64, error)
	UpdateCouponUser(ctx.Context, *BasicCouponUserInfoEntity) (int64, error)
	UpdateCouponUserUserId(ctx.Context, *BasicCouponUserInfoEntity) (int64, error)
	CountCouponUser(ctx.Context, *BasicCouponUserInfoEntity) int64
	CouponUserList(c ctx.Context, request *UserCouponListRequest) (int64, *[]*UserCouponListResponse, error)
	SelectCouponUser(ctx.Context, *BatchCouponUserEntity) (*[]BasicCouponUserInfoEntity, error)
	DetailByCoupon(c ctx.Context, request *WithGoodsRequest) (*BasicInfoEntity, error)
	SelectCouponUserUserIdList(c ctx.Context) (*[]BasicCouponUserInfoEntity, error)
	UpdateUnionIdCouponUserUserId(ctx.Context, *BasicCouponUserInfoEntity) (int64, error)
	//批量查询优惠券是否被其他策略占用
	GetMarketingStrategyIsUsedBatch(c ctx.Context, WxCouponIDS []int64) ([]int64, error)
	//批量查询优惠券列表
	GetAllCouponAvailableList(c ctx.Context, request WxCouponAvailableRequest) ([]*WxCouponAvailable, error)
}
