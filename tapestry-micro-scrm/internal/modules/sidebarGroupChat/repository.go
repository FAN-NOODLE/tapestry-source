package sidebarGroupChat

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

// Repository 仓库接口
type Repository interface {
	repos.IRepository
	GroupChatInfo(ctx ctx.Context, groupChatInfoRequest GroupChatInfoRequest) (*WorkGroupChat, error)
	GroupChatLeaderInfo(ctx ctx.Context, groupChatInfoRequest GroupChatInfoRequest) (*WorkGroupChat, error)
	GroupChatMember(ctx ctx.Context, groupChatInfoRequest GroupChatInfoRequest) ([]*WorkGroupChatMember, error)
	GroupChatMemberUnionid(ctx ctx.Context, groupChatInfoRequest GroupChatInfoRequest) (*[]WorkGroupChatMember, error)
	GroupChatCorpUser(ctx ctx.Context, uid string, enterpriseId int64) (*[]WorkGroupCorpUser, error)
	GroupChatMemberMsg(ctx ctx.Context, chatMsgRequest ChatMsgRequest) (*[]ChatMsgResponse, error)
	UsersAddEmployees(ctx ctx.Context, enterpriseID int64, workExternalUserID int64, euid int64) (*[]GroupUsersAddEmployees, error)
}
