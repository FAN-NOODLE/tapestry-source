package externalUserJob

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	externalUserService "micro-scrm/internal/app/service/external_user"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/xjob"
)

type syncExternalUser struct {
	name string
}

func NewSyncExternalUser() *syncExternalUser {
	return &syncExternalUser{
		name: "同步外部联系人信息",
	}
}

func (job syncExternalUser) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

	var (
		err            error
		enterpriseList []boss.EnterpriseModel
	)
	// 获取全部企业
	enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}
	for _, enterprise := range enterpriseList {

		err = job.handleEnterpriseUser(ctx, enterprise.ID)
		if err != nil {
			log.Error(ctx).Msgf("处理企业%d异常,错误:%v", enterprise.ID, err)
		}

	}
	return nil
}

func (job syncExternalUser) handleEnterpriseUser(ctx ctx.Context, enterpriseId int64) error {

	// 处理全部用户
	err := externalUserService.NewWxUserCorpExternal().SyncExternalUserJob(ctx, enterpriseId)
	if err != nil {
		return err

	}
	return nil
}
