package plan_tasks

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/xjob"
)

type MomentComments struct {
	Name string
	Days int64
}

func NewMomentComments() MomentComments {
	return MomentComments{
		Name: "计划任务朋友圈任务顾客评论互动查询推送MA",
		Days: 3,
	}
}

func (job MomentComments) Run(c ctx.Context, handler *xjob.JobHandler) error {
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

func (job MomentComments) handle(c ctx.Context, enterpriseId int64) error {
	var err error
	// 查询活动有效期内进行中的企业群发的子任务
	taskList, err := planService.NewPlanTask().GetMaMomentTaskByDay(c, enterpriseId, job.Days)
	if err != nil || len(taskList) == 0 {
		return err
	}
	enterpriseInfo, err := common.GetEnterpriseInfo(c, enterpriseId)
	if err != nil {
		return err
	}
	qyWorkApi := new(work.QyWorkImpl).SetDomain(enterpriseInfo.Domain)
	for _, item := range taskList {
		if item.MomentJobId == "" {
			continue
		}
		err = planService.NewPlanTask().QyMomentComments(c, enterpriseId, item, qyWorkApi)
		if err != nil {
			fmt.Printf("异步处理朋友圈互动结果失败，err:%v", err)
			continue
		}
	}
	fmt.Printf("执行企微朋友圈互动查询完成")
	return nil
}
