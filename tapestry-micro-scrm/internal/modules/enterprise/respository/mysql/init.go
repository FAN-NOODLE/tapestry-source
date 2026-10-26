package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/pkg/db"
)

// New 初始化存储
func New() enterprise.Repository {
	return &enterpriseRepo{repos.NewSQLRepo("t_enterprise_corp", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &enterprise.Corp{}
	}))}
}
