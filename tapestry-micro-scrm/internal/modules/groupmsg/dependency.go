package groupmsg

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/ws_trace"
	"micro-scrm/internal/pkg/common"
)

type WorkApiInterface interface {
	SendGroupMsg(ctx ctx.Context, sender string, content string, chatType string, domain string, WsTraceService *ws_trace.Service, MemberID int64, taskID int64, enterpriseID int64) (string, error)
	GroupMsgSendResult(c ctx.Context, sender string, msgID string, domain string) ([]*ChatGroupMsgSendRecord, error)
	SetMediaID(c ctx.Context, content string, enterpriseID int64, domain string) (string, error)
}

type MemberInterface interface {
	GetMemberUserID(c ctx.Context, memberID int64) (string, error)
}

type EnterpriseInterface interface {
	GetInfo(c ctx.Context, enterpriseID int64) (domain string, corpID string, hash string, err error)
}

type UCInterface interface {
	GetMemUID(c ctx.Context, userID string, domain string) (string, error)
}

type ShortLinkInterface interface {
	CreateShortParamID(c ctx.Context, param common.Params) (int64, error)
}
