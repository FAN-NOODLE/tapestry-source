package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/overtime"
	"micro-scrm/internal/pkg/db"
)

//New 初始化存储
func New() overtime.Repository {
	return &OvertimeRepo{SQLRepo: repos.NewSQLRepo("t_chat_overtime_records", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &overtime.Records{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]overtime.Records{}
	}))}
}
