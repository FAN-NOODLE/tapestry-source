package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/fission"
	"micro-scrm/internal/pkg/db"
)

func New() fission.Repository {
	return &FissionRepo{repos.NewSQLRepo("t_fission", db.DB())}
}
