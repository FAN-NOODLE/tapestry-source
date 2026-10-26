package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/callback"
	"micro-scrm/internal/pkg/db"
)

//
func New() callback.Repository {
	return &CallbackRepo{SQLRepo: repos.NewSQLRepo("", db.DB())}
}
