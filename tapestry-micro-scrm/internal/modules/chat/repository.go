package chat

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/pkg/sdk/work"
)

type MsgRepository interface {
	repos.IRepository
	Exists(c ctx.Context, enterpriseId int64, msgId string) (bool, error)
	LastRow(c ctx.Context, result interface{}, enterpriseId int64) error
	ListByID(c ctx.Context, cond MsgListByIDRequest) ([]Msg, error)
	GetNotMediaList(c ctx.Context, req []string, enterpriseId int64, StartAt int64) (*[]Msg, error)
	GetMembersList(c ctx.Context, req []string, enterpriseId int64) (*[]UserInfo, error)
	GetExternalUsersList(c ctx.Context, req []string, enterpriseId int64) (*[]UserInfo, error)
}

type MemberRepository interface {
	repos.IRepository
	FindMember(c ctx.Context, cond interface{}, res interface{}) error
}

type ObjectRepository interface {
	repos.IRepository
	Find(c ctx.Context, cond interface{}, result interface{}) error
	GetMemberByGroup(c ctx.Context, req string, enterpriseId int64) (*[]GroupMemberInfo, error)
}

type ReceiveUserRepository interface {
	repos.IRepository
	SelectReceiveUserByMemberUserIDs(ctx ctx.Context, userIDds []string) ([]SearchReceiveResponse, error)
}

type MemberToGroupRepository interface {
	repos.IRepository
}

type MsgMediaRepository interface {
	repos.IRepository
	ListByFileId(c ctx.Context, enterpriseID int64, sdkFileidList []string) ([]MsgMedia, error)
}

type GroupChatRepository interface {
	repos.IRepository
	FindWorkGroupChatDetail(ctx ctx.Context, req WorkGroupChat) (*WorkGroupChat, error)
}

type GroupChatMemberRepository interface {
	repos.IRepository
	BatchUpdateWorkGroupChatMembers(ctx ctx.Context, where map[string]interface{}, update WorkGroupChatMember) (int64, error)
	BatchUpdateExternalWorkGroupChatMembersInfo(c ctx.Context, EnterpriseId int64, groupChatId int64, chatMembers []work.ExternalGroupMemberList) (int64, error)
}
