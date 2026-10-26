package shop

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

// WxCouponListByIds 根据id获取微信优惠券列表
func WxCouponListByIds(c ctx.Context, request *WxCouponListByIdsRequest) (*WxCouponListByIdsResponse, error) {
	var response WxCouponListByIdsResponse
	err := httpClient.Post(c, GetURL(urlEnum.WxCouponListByIds), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
