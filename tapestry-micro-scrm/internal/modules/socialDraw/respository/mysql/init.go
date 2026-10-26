package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/socialDraw"
	"micro-scrm/internal/pkg/db"
)

func NewRepo() socialDraw.Repository {
	return &GroupChatRepo{repos.NewSQLRepo("t_work_group_chat", db.DB())}
}
