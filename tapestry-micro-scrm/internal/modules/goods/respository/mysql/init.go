package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/goods"
	"micro-scrm/internal/pkg/db"
)

func New() goods.Repository {
	return &GoodsRepo{repos.NewSQLRepo("t_goods", db.DB())}
}
