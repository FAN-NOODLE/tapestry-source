package work

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

type Api interface {
	SetDomainValue(string)
	SendMsgTemplate(c ctx.Context, req SendMsgTemplateReq) (*QyMassHairResponse, error)
	GetGroupMsgSendResult(c ctx.Context, req GroupMsgSendResultRequest) (*GroupMsgSendResultResponse, error)
	SetMediaID(c ctx.Context, content string, enterpriseID int64) (string, error)
}
