package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/sensitive"
	"micro-scrm/internal/pkg/db"
)

//New 初始化存储
func New() sensitive.Repository {
	return &SensitiveRepo{SQLRepo: repos.NewSQLRepo("t_chat_sensitive", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &sensitive.Detail{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]sensitive.Detail{}
	}))}
}
