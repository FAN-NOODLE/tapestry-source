package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/marketingPromotion"
	"micro-scrm/internal/pkg/db"
)

func New() marketingPromotion.Repository {
	return &MarketingPromotionRepo{SQLRepo: repos.NewSQLRepo("t_marketing_promotion", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &marketingPromotion.GetPosterRequest{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]marketingPromotion.GetPosterRequest{}
	}))}
}
