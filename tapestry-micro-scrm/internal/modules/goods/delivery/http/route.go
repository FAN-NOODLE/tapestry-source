package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-scrm/internal/modules/goods"
	"micro-scrm/internal/modules/goods/respository/mysql"
	"net/http"
)

//Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {
	// 初始化服务
	goodsService = goods.NewGoods(mysql.New())

	return []*swagger.Endpoint{

		endpoint.Post(
			"goods/sku-list",
			"商品列表",
			endpoint.Handler(skuList),
			endpoint.Tags("商品"),
			endpoint.Response(http.StatusOK, "", "goods", "接口调用成功"),
		),

		endpoint.Post(
			"goods/sync",
			"商品数据同步",
			endpoint.Handler(syncGoods),
			endpoint.Tags("商品"),
			endpoint.Response(http.StatusOK, "", "syncGoods", "接口调用成功"),
		),

		endpoint.Post(
			"goods/spu-list",
			"Spu商品列表",
			endpoint.Handler(goodSpuList),
			endpoint.Tags("商品"),
			endpoint.Response(http.StatusOK, "", "goods", "接口调用成功"),
		),
	}
}
