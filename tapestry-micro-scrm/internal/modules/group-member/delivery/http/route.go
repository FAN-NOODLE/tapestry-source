package groupmember_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	groupMember "micro-scrm/internal/modules/group-member"
	"micro-scrm/internal/modules/group-member/respository/clickhouse"
	workexternal "micro-scrm/internal/modules/work-external"
	"micro-scrm/internal/modules/work-external/respository/click_house"
	"net/http"
)

// Endpoints 路由配置
func Endpoints() []*swagger.Endpoint {

	workService = workexternal.New(click_house.NewRepo())
	groupService = groupMember.New(clickhouse.NewRepo())

	return []*swagger.Endpoint{
		endpoint.Post("/group-member/calc", "群成员计算",
			endpoint.Handler(Calc),
			endpoint.Tags("群成员计算"),
			endpoint.Response(http.StatusOK, "",
				"calc", "接口调用成功")),
	}
}
