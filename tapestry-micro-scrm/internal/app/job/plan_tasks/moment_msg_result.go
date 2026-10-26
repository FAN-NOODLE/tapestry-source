package plan_tasks

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/xjob"
)

type MomentMsgResult struct {
	Name string
}

func NewMomentMsgResult() MomentMsgResult {
	return MomentMsgResult{
		Name: "计划任务朋友圈任务的结果查询",
	}
}

func (job MomentMsgResult) Run(c ctx.Context, handler *xjob.JobHandler) error {
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

func (job MomentMsgResult) handle(c ctx.Context, enterpriseId int64) error {
	var err error
	// 查询活动有效期内进行中的企业群发的子任务
	planIds, err := planService.NewPlanTask().GetPlanIdByStatus(c, enterpriseId, microScrmModel.PlanIng)
	if err != nil {
		return err
	}
	if len(planIds) == 0 {
		return nil
	}
	taskLists, err := planService.NewPlanTask().GetQyGroupListsByPlanId(c, enterpriseId, microScrmModel.PLANTASKQYSEND, planIds)
	if err != nil {
		return err
	}
	if len(taskLists) == 0 {
		return nil
	}
	enterpriseInfo, err := common.GetEnterpriseInfo(c, enterpriseId)
	if err != nil {
		return err
	}
	qyWorkApi := new(work.QyWorkImpl).SetDomain(enterpriseInfo.Domain)
	// 调用企微的发送结果结果，根据返回结果中的值，如果有一条的结果值大于0，则修改task_allot_log表中的is_group_send状态为1，已推送
	for _, item := range taskLists {
		if item.TaskCategory != microScrmModel.PLANTASKMOMENT {
			continue
		}
		if item.MomentJobId == "" {
			continue
		}
		exist, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, enterpriseId, map[string]interface{}{
			"task_id":     item.ID,
			"is_finished": microScrmModel.PLAN_FISISHING,
		})
		if exist.ID == 0 {
			continue
		}
		err = planService.NewPlanTask().QyMomentSendResult(c, enterpriseId, item, qyWorkApi)
		if err != nil {
			fmt.Printf("异步处理朋友圈群发结果失败，err:%v", err)
			continue
		}
	}
	fmt.Printf("执行企微朋友圈群发结果查询完成")
	return nil
}
