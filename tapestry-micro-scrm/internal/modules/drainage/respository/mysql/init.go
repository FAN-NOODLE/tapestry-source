package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/drainage"
	"micro-scrm/internal/pkg/db"
)

func NewDrainage() drainage.Repository {
	return &DrainageRepo{SQLRepo: repos.NewSQLRepo("t_drainage_welcomes", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &drainage.Welcome{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]drainage.Welcome{}
	}))}
}
