package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// GoodsApi 商品SDK
type GoodsApi struct {
	BaseSDK
}

// NewGoods 返回API实例
func NewGoods(host string) *GoodsApi {
	urlHost = host
	return &GoodsApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// GoodsList 商品列表
func (g GoodsApi) GoodsSkuList(c ctx.Context, request *GoodsListReq) (*SpuListRes, error) {
	path := "/micro-scrm/v1/goods/sku-list"
	url := urlHost + path
	res, err := g.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &SpuListRes{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*SpuListRes), err
}

// GoodsSpuList 商品列表
func (g GoodsApi) GoodsSpuList(c ctx.Context, request *GoodsSpuListReq) (interface{}, error) {
	path := "/micro-scrm/v1/goods/spu-list"
	url := urlHost + path
	res, err := g.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &ResponseList{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}