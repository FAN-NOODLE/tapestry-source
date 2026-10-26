package ck

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/socialDraw"
	"micro-scrm/internal/pkg/db"
)

type GroupChatRepo struct {
	*repos.SQLRepo
}

func NewRepo() socialDraw.Ckrepository {
	return &GroupChatRepo{repos.NewSQLRepo("t_work_group_chat", db.ClickHouse())}
}
