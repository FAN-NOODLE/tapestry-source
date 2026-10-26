package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/remindconfig"
	"micro-scrm/internal/pkg/db"
)

//New 初始化存储
func New() remindconfig.Repository {
	return &RemindConfigRepo{SQLRepo: repos.NewSQLRepo("t_chat_remind_config", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &remindconfig.ChatRemindConfig{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]remindconfig.ChatRemindConfig{}
	}))}
}
