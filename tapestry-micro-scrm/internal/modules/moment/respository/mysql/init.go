package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/department"
	"micro-scrm/internal/modules/department/respository/mysql"
	"micro-scrm/internal/modules/members"
	mysql2 "micro-scrm/internal/modules/members/respository/mysql"
	"micro-scrm/internal/modules/moment"
	"micro-scrm/internal/pkg/db"
)

func NewTopic() moment.Repository {
	return &MomentRepo{SQLRepo: repos.NewSQLRepo("t_moment_topics", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &moment.MomentTopic{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]moment.MomentTopic{}
	}))}
}

func NewMember() moment.Repository {
	return &MomentRepo{SQLRepo: repos.NewSQLRepo("t_moment_members", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &moment.MomentMember{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]moment.MomentMember{}
	}))}
}

func NewTopicMember() moment.Repository {
	return &MomentRepo{SQLRepo: repos.NewSQLRepo("t_moment_topic_members", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &moment.MomentTopicMember{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]moment.MomentTopicMember{}
	}))}
}

func NewDepartment() department.Repository {
	return &mysql.DepartmentRepo{SQLRepo: repos.NewSQLRepo("t_departments", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &department.Department{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]department.Department{}
	}))}
}

func NewOriginalMember() members.Repository {
	return &mysql2.MemberRepo{SQLRepo: repos.NewSQLRepo("t_members", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &members.Member{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]members.Member{}
	}))}
}
