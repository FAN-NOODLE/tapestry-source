package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/sidebarGroupChat"
	"micro-scrm/internal/pkg/db"
)

func New() sidebarGroupChat.Repository {
	return &SidebarGroupChatRepo{SQLRepo: repos.NewSQLRepo("work_group_chat", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &sidebarGroupChat.GroupChatInfoRequest{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]sidebarGroupChat.GroupChatInfoRequest{}
	}))}
}
