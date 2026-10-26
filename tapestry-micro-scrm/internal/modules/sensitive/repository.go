package sensitive

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

type Repository interface {
	repos.IRepository
	SelectList(ctx ctx.Context, req ListRequest) (int64, *[]Detail, error)
	CountBySensitive(ctx ctx.Context, ids []int64) (*[]CountBySensitive, error)
	FindSensitive(ctx ctx.Context, request SaveRequest) (*Detail, error)
	DeleteSensitive(ctx ctx.Context, request DeleteRequest) (int64, error)
	DeleteSensitiveRecords(ctx ctx.Context, request DeleteRequest) (int64, error)
	SelectRecordsList(ctx ctx.Context, req DetailListRequest) (int64, *[]Records, error)
	SelectChatMsgList(ctx ctx.Context, ids []int64) (*[]ChatMsg, error)
	SelectUserRelation(ctx ctx.Context, request ChatMsgMembers) (*DetailListResponse, error)
	InsertSensitiveRecords(ctx ctx.Context, request Records) (int64, error)
	GetMembersList(c ctx.Context, req []string, enterpriseId int64) (*[]UserInfo, error)
}
