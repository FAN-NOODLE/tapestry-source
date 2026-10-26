package marketingPromotionhttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/marketingPromotion"
	"micro-scrm/internal/modules/marketingPromotion/respository/mysql"
	"net/http"
)

//Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {
	marketingPromotionService = marketingPromotion.New(mysql.New())
	enterpriseService = enterprise.New(enterpriseMysql.New())
	return []*swagger.Endpoint{
		endpoint.Post("/marketing-promotion/poster", "获取海报",
			endpoint.Handler(getPoster),
			endpoint.Tags("获取海报"),
			endpoint.Response(http.StatusOK, "", "poster", "接口调用成功")),
	}
}
