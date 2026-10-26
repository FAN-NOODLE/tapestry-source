package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/marketing"
	"micro-scrm/internal/pkg/db"
)

//New 初始化存储
func New() marketing.Repository {
	return &MarketingRepo{SQLRepo: repos.NewSQLRepo("t_marketings", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &marketing.MarketingStrategy{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]marketing.MarketingStrategy{}
	}))}
}
