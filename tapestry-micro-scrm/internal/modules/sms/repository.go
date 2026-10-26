package sms

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

// Repository 仓库接口
type Repository interface {
	repos.IRepository
	SmsSignFindOne(ctx ctx.Context, SmsSignRequest SmsSignInfoRequest) (*SmsSignInfo, error)
	CreateSmsSign(ctx ctx.Context, SmsSignRequest SmsSignInfo) (int64, error)
	CreateSmsTemplate(ctx ctx.Context, SmsTemplateRequest SmsTemplate) (int64, error)
	SmsTemplateFindOne(ctx ctx.Context, SmsSignRequest SmsTemplate) (*SmsTemplate, error)
	CheckDefaultSmsTemplate(ctx ctx.Context, smsTemplate SmsTemplate) (*SmsTemplate, error)
	UpdateTemplateStatus(ctx ctx.Context, SmsSignRequest UpdateSmsTemplateStatus) error
	SelectSmsTemplateSignList(ctx ctx.Context) (*[]SmsSignInfo, error)
	SelectSmsTemplateParamsList(ctx ctx.Context) (*[]SmsParams, error)
	SmsTemplateDelete(ctx ctx.Context, UpdateSmsTemplateStatus UpdateSmsTemplateStatus) error
	SmsTemplateEdit(ctx ctx.Context, UpdateSmsTemplateStatus UpdateSmsTemplateStatus) error
	UpdateSmsSign(ctx ctx.Context, SmsSignInfoRequest SmsSignInfoRequest) error
	SelectTemplateList(ctx ctx.Context, EnterpriseId int64) (*[]SmsTemplate, error)
	CreateSmsLog(ctx ctx.Context, SendLogRequest SendLog) (int64, error)
	SmsCount(ctx ctx.Context, request Request) (*CountResponse, error)
	SmsTemplateSignList(cc ctx.Context, templateSignRequest TemplateListRequest) (int64, *[]SmsSignInfo, error)
	SmsTemplateList(c ctx.Context, TemplateListRequest TemplateListRequest) (int64, *[]SmsTemplate, error)
	GetEnterpriseSmsTemplateCount(c ctx.Context, SmsTemplateRequest SmsTemplate) int64
	GetDefaultTemplate(ctx ctx.Context, EnterpriseId int64, Name string) (*SmsTemplate, error)
	SelectTemplateListByIDs(ctx ctx.Context, EnterpriseId int64, templateIDs []int64, strategyIDs []int64) (*[]SmsTemplate, error)
}
