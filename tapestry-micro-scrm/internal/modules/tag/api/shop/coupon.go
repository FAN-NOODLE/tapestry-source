package shop

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/pkg/sdk/shop"
)

// GetWxCouponListByIds 通过id查询微信优惠券列表
func GetWxCouponListByIds(c ctx.Context, ids []int64) (map[int64]*entity.TagCouponInfo, error) {
	var request = shop.WxCouponListByIdsRequest{
		CouponIds: ids,
	}
	response, err := shop.WxCouponListByIds(c, &request)
	if err != nil {
		return nil, err
	}

	var ret = make(map[int64]*entity.TagCouponInfo)

	for _, v := range response.Data {
		ret[v.ID] = &entity.TagCouponInfo{
			Name:    v.CouponName,
			StockID: v.CouponStockID,
		}
	}
	return ret, nil
}
