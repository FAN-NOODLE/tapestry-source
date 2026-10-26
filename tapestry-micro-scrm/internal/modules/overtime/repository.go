package overtime

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

type Repository interface {
	repos.IRepository
	SelectRecordsList(ctx ctx.Context, req ListRequest) (int64, *[]Records, error)
	SelectChatMsgList(ctx ctx.Context, ids []int64) (*[]ChatMsg, error)
	SelectUserRelation(ctx ctx.Context, request ChatMsgMembers) (*ChatMsgMembers, error)
	SelectMembersList(c ctx.Context, req []string, enterpriseId int64) (*[]MemberInfo, error)
	MemberMultiDepartments(ctx ctx.Context, memberIds []int64) (map[int64]MemberDepartsList, error)
	SelectReplyMsg(ctx ctx.Context, request ChatMsg) (*[]ChatMsg, error)
	FindOvertime(ctx ctx.Context, request Records) (*Records, error)
	SelectMemberByWhere(ctx ctx.Context, request ListRequest) (*[]SearchMemberResponse, error)
}
