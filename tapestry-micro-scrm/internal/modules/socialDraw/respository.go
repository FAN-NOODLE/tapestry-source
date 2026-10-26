package socialDraw

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

type Repository interface {
	repos.IRepository

	FindWorkGroupChatList(ctx ctx.Context, req WorkGroupListRequest) ([]*WorkGroupChatResponse, error)
	WorkGroupChatSetting(ctx ctx.Context, req WorkGroupChatSettingRequest) error
	ManageWorkGroupType(ctx ctx.Context, req WorkGroupChatType) (*WorkGroupChatType, error)
	FindWorkGroupTypeList(ctx ctx.Context, req WorkGroupChatType) ([]*WorkGroupChatType, error)
	GetWorkGroupTypeInfo(ctx ctx.Context, req WorkGroupChatType) (*WorkGroupChatType, error)
	GetWorkGroupMembers(ctx ctx.Context, req WorkGroupChatMembersRequest) ([]*WorkGroupChatMember, int, error)
	GetWorkGroupExternalUsers(ctx ctx.Context, enterpriseId int64, userIds []string) (map[string]*WxUserCorpExternal, error)
	CountWorkGroupByType(ctx ctx.Context, req WorkGroupChatType) (map[int64]int, error)

	WorkGroupFilters(ctx ctx.Context, req WorkGroupFilterRequest) ([]*WorkGroupFilterResponse, error)
	CheckWorkGroupChatTypeUsed(ctx ctx.Context, enterpriseId int64, chatTypeId int64) (bool, error)
	GetWorkGroupChatListFilters(ctx ctx.Context, req WorkGroupListRequest) []string
	GetMarsID(ctx ctx.Context, enterpriseId int64, UserID string) (string, error)
	GetWorkGroupMembersInfo(ctx ctx.Context, req WorkGroupChatMembersRequest) ([]*WorkGroupChatMemberList, int, error)
	GetWorkGroupUserActivityInfo(ctx ctx.Context, enterpriseId int64, UserID string, cid int64) ([]*ActInfo, error)
}

type Ckrepository interface {
	repos.IRepository
	GroupUserUvStatistics(ctx ctx.Context, req GroupInfoReq) (*WorkGroupChatStatistic, error)
	GroupMemberStatistics(ctx ctx.Context, req GroupInfoReq) (*WorkGroupChatStatistic, error)
	GroupMemberStatisticsForCrm(ctx ctx.Context, req GroupInfoReq) (*WorkGroupChatStatistic, error)
	NewGroupUserChat(ctx ctx.Context, req GroupInfoReq) ([]*Chart, error)
	NewGroupUserChatForCrm(ctx ctx.Context, req GroupInfoReq) ([]*Chart, error)
	ValidGroupUserChat(ctx ctx.Context, req GroupInfoReq) ([]*Chart, error)
	ValidGroupUserChatForCrm(ctx ctx.Context, req GroupInfoReq) ([]*Chart, error)
	GetWorkGroupStatics(ctx ctx.Context, req WorkGroupChatStaticsRequest) (map[int64]*WorkGroupChatStatistic, error)
	GetWorkGroupMembersList(ctx ctx.Context, req WorkGroupChatMembersRequest) ([]*WorkGroupChatMemberList, int, error)
}
