package plan_tasks

import (
	"github.com/panjf2000/ants/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	membersDao "micro-scrm/internal/app/dao/members"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/xjob"
	sync2 "sync"
)

type MaTaskDataEnd struct {
	Name string
}

func NewMaTaskDataEnd() MaTaskDataEnd {
	return MaTaskDataEnd{
		Name: "MA任务结束的推送(三方根据获取任务模板创建的任务)",
	}
}

func (job MaTaskDataEnd) Run(c ctx.Context, handler *xjob.JobHandler) error {

	var (
		err            error
		enterpriseList []boss.EnterpriseModel
	)
	enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}

	for _, v := range enterpriseList {
		enterpriseId := v.ID
		err = job.handle(c, enterpriseId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (job MaTaskDataEnd) handle(c ctx.Context, enterpriseId int64) error {
	lists, err := planDao.NewPlanTaskDao().GetMaPlanEndLists(c, enterpriseId)
	if err != nil {
		return err
	}
	if len(lists) == 0 {
		return nil
	}
	memberInfoMap := make(map[string]string)
	memberLists, err2 := membersDao.NewMembers().GetMany(c, map[string]interface{}{"enterprise_id": enterpriseId, "deleted_at": 0})
	if err2 != nil {
		return err2
	}
	for _, v := range memberLists {
		item := v
		memberInfoMap[item.Userid] = item.Name
	}
	p, err := ants.NewPool(5)
	defer p.Release()
	wg := sync2.WaitGroup{}
	for _, value := range lists {
		info := value
		wg.Add(1)
		taskId := info.ID
		err = p.Submit(func() {
			defer wg.Done()
			planService.NewPlanTask().MaFinishedPush(c, enterpriseId, taskId, memberInfoMap)
		})
	}
	wg.Wait()
	return nil
}
