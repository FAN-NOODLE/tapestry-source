package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/sms"
	"micro-scrm/internal/pkg/db"
)

func New() sms.Repository {
	return &SmsTemplateRepo{SQLRepo: repos.NewSQLRepo("t_sms_sign", db.DB(), repos.EntityFactory(func(context ctx.Context) interface{} {
		return &sms.SmsSignInfo{}
	}), repos.ListFactory(func(context ctx.Context, i int64) interface{} {
		return &[]sms.SmsSignInfo{}
	}))}
}
