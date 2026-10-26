package excitation

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	appJob "micro-scrm/internal/app/job"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	"micro-scrm/internal/pkg/xjob"
)

type CoinsExpired struct {
	name            string
	coinsExpiredSvc *excitationSvc.CoinsService
}

func NewCoinsExpired() *CoinsExpired {
	return &CoinsExpired{
		name:            "累计数据获得积分任务",
		coinsExpiredSvc: excitationSvc.NewCoinsService(),
	}
}

func (job *CoinsExpired) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	var wg threading.RoutineGroup

	for _, id := range enterpriseIds {
		tmpId := id
		wg.Run(func() {
			err := job.Handle(c, tmpId)
			if err != nil {
				app.Logger().Error(c).Msgf("error : %v", err)
			}
		})
	}

	wg.Wait()

	return nil
}

func (job *CoinsExpired) Handle(c ctx.Context, enterpriseId int64) error {
	err := excitationSvc.NewCoinsService().CoinsExpiredProduce(c, enterpriseId)
	if err != nil {
		return err
	}
	return nil
}
