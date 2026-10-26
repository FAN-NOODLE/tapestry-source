package clickhouse

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	groupMember "micro-scrm/internal/modules/group-member"
	"micro-scrm/internal/pkg/db"
)

// NewRepo 初始化仓库
func NewRepo() groupMember.Repository {
	return &GroupRepo{repos.NewSQLRepo("t_work_group_chat", db.ClickHouse())}
}
