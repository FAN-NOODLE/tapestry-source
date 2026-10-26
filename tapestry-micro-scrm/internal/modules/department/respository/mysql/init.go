package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/department"
	"micro-scrm/internal/pkg/db"
)

//New 初始化存储
func New() department.Repository {
	return &DepartmentRepo{repos.NewSQLRepo("t_departments", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &department.Department{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]department.Department{}
	}))}
}
