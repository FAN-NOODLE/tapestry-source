package click_house

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	workexternal "micro-scrm/internal/modules/work-external"
	"micro-scrm/internal/pkg/db"
)

//New 初始化存储
func NewRepo() workexternal.Repository {
	return &WorkRepo{repos.NewSQLRepo("t_work", db.ClickHouse())}
	//return nil
}
