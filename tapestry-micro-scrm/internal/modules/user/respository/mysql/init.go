package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/user"
	"micro-scrm/internal/pkg/db"
)

// NewUser 创建用户仓库
func NewUser() user.Repository {
	return &UserRepo{repos.NewSQLRepo("t_user", db.DB())}
}

// NewEvent 创建用户事件仓库
func NewEvent() user.EventRepository {
	return &EventRepo{repos.NewSQLRepo("t_user_events", db.DB())}
}
