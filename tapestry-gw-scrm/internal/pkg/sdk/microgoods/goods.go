package microgoods

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

// GoodsSpuList 获取商品Spu列表
func GoodsSpuList(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.GoodsSpuList), request)
	return response, err
}
