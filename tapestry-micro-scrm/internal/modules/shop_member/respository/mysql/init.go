package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/shop_member"
	"micro-scrm/internal/pkg/db"
)

// New 初始化存储
func New() shop_member.Repository {
	return &ShopMemberRepository{repos.NewSQLRepo("t_shop_members", db.DB())}
}
