package notification

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	appJob "micro-scrm/internal/app/job"
	notificationSvc "micro-scrm/internal/app/service/notification"
	"micro-scrm/internal/pkg/xjob"
)

type TimedRelease struct {
	name            string
	timedReleaseSvc *notificationSvc.InfoMergeSvc
}

func NewTimedRelease() *TimedRelease {
	return &TimedRelease{
		name:            "累计数据获得积分任务",
		timedReleaseSvc: notificationSvc.NewInfoMergeSvc(),
	}
}

func (job *TimedRelease) Run(c ctx.Context, handler *xjob.JobHandler) error {
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

func (job *TimedRelease) Handle(c ctx.Context, enterpriseId int64) error {
	err := job.timedReleaseSvc.TimedRelease(c, enterpriseId)
	if err != nil {
		return err
	}
	return nil
}
