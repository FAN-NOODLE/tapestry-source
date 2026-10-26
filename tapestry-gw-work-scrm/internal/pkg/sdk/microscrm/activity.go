package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// ActivityApi 商品SDK
type ActivityApi struct {
	BaseSDK
}

// NewGoods 返回API实例
func NewActivity(host string) *ActivityApi {
	urlHost = host
	return &ActivityApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// ActivityGoodsList 商品列表
func (g ActivityApi) ActivityGoodsList(c ctx.Context, request *ActivityGoodsListRequest) (*ActivityGoodsListResponse, error) {
	path := "/micro-scrm/v1/activity/goods/list"
	url := urlHost + path
	res, err := g.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &ActivityGoodsListResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*ActivityGoodsListResponse), err
}
