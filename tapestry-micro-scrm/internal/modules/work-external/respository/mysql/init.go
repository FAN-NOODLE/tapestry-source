package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	workexternal "micro-scrm/internal/modules/work-external"
	"micro-scrm/internal/pkg/db"
)

// NewTestRepo 初始化仓库
func NewRepo() workexternal.Repository {
	return &WorkRepo{repos.NewSQLRepo("t_work_external_user_relation", db.DB())}
}
