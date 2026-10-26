package drainageContactJob

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	appJob "micro-scrm/internal/app/job"
	cacheSvc "micro-scrm/internal/app/service/cache"
	"micro-scrm/internal/app/service/drainage_contacts"
	"micro-scrm/internal/pkg/xjob"
)

type DrainageContactsUpdate struct {
	name                string
	drainageContactsSvc *drainageContactsSvc.MainSvc
}

func NewDrainageContactsUpdate() *DrainageContactsUpdate {
	return &DrainageContactsUpdate{
		name:                "累计数据获得积分任务",
		drainageContactsSvc: drainageContactsSvc.NewMainSvc(),
	}
}

func (job *DrainageContactsUpdate) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}
	var wg threading.RoutineGroup

	for _, id := range enterpriseIds {
		cacheRes, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(c, id, "store_name")
		if cacheRes.Name != "coach" {
			continue
		}
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

func (job *DrainageContactsUpdate) Handle(c ctx.Context, enterpriseId int64) error {
	err := job.drainageContactsSvc.UpdateContact(c, enterpriseId)
	if err != nil {
		return err
	}
	return nil
}
