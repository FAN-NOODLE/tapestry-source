package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/groupmsg"
	"micro-scrm/internal/pkg/db"
)

func New() groupmsg.Repository {
	return &Repo{
		repos.NewSQLRepo("", db.DB()),
	}
}
