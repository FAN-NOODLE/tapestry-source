package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/activity"
	"micro-scrm/internal/pkg/db"
)

func New() activity.Repository {
	return &ActivityRepo{repos.NewSQLRepo("t_shop_activity", db.DB())}
}
