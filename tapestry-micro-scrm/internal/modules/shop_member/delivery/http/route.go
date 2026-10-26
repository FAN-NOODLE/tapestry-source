package shop_member_http

import (
	"github.com/go-playground/validator/v10"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-scrm/internal/modules/shop_member"
	shopMemberMysql "micro-scrm/internal/modules/shop_member/respository/mysql"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	validate = validator.New()
	shopMemberService = shop_member.New(shopMemberMysql.New())
	return []*swagger.Endpoint{
		endpoint.Post("/shop-member/level-list", "会员等级列表",
			endpoint.Handler(levelList),
			endpoint.Tags("会员等级列表"),
			endpoint.Response(http.StatusOK, "", "level-list", "接口调用成功")),
	}
}
