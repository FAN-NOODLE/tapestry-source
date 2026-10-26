package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/department"
	"micro-scrm/internal/modules/department/respository/mysql"
	"micro-scrm/internal/modules/members"
	"micro-scrm/internal/pkg/db"
)

//New 初始化存储
func New() members.Repository {
	return &MemberRepo{SQLRepo: repos.NewSQLRepo("t_members", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &members.Member{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]members.Member{}
	}))}
}

//NewDepartment 初始化组织存储
func NewDepartment() department.Repository {
	return &mysql.DepartmentRepo{SQLRepo: repos.NewSQLRepo("t_departments", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &department.Department{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]department.Department{}
	}))}
}
