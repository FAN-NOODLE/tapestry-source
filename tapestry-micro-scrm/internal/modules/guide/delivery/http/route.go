package guideHttp

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"micro-scrm/internal/modules/guide"
	"micro-scrm/internal/modules/guide/respository/mysql"
	"micro-scrm/internal/modules/members"
	mysql2 "micro-scrm/internal/modules/members/respository/mysql"
)

func Endpoints() []*swagger.Endpoint {
	membersService = members.New(mysql2.New())
	guideService = guide.NewInside(mysql.NewGuide(), membersService)
	return []*swagger.Endpoint{}
}
