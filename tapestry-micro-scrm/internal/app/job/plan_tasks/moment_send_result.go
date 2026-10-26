package plan_tasks

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	workDao "micro-scrm/internal/app/dao/work"
	planReq "micro-scrm/internal/app/http/request/plan"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/xjob"
)

type MomentSendResult struct {
	Name string
}

func NewMomentSendResult() MomentSendResult {
	return MomentSendResult{
		Name: "朋友圈任务脚本创建是否完成结果查询",
	}
}

func (job MomentSendResult) Run(c ctx.Context, handler *xjob.JobHandler) error {
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

func (job MomentSendResult) handle(c ctx.Context, enterpriseId int64) error {
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

	for _, list := range taskLists {
		item := list
		taskId := item.ID
		if item.TaskCategory != microScrmModel.PLANTASKMOMENT {
			continue
		}
		if item.MomentJobId == "" {
			continue
		}

		if item.MomentStatus != 3 {
			var massHairReq = work.MomentTaskResultRequest{}
			massHairReq.JobId = item.MomentJobId
			res, err := qyWorkApi.GetMomentTaskResult(c, massHairReq)
			if err != nil {
				app.Logger().Error(c).Msgf("获取朋友圈的发表结果失败, err:%s", err.Error())
				return err
			}
			if res.ErrCode != 0 {
				app.Logger().Error(c).Msgf("获取朋友圈的发表结果返回错误, res:%v", res)
				return err
			}
			if res.Status != 3 {
				continue
			}
			momentId := res.Result.MomentId
			err = planDao.NewPlanTaskDao().UpdateData(c, enterpriseId, map[string]interface{}{
				"id": taskId,
			}, map[string]interface{}{"moment_status": res.Status, "moment_job_id": momentId})

			if item.TagExpireType != microScrmModel.TAG_TASK_SEND_EXPIRE {
				continue
			}
			if item.QyTagId != "" { // 有企业标签的朋友圈任务需要删除标签
				err = workDao.NewPlanTasksTagRelation().DeleteData(c, map[string]interface{}{
					"enterprise_id": enterpriseId,
					"plan_id":       taskId,
				})
				data := planReq.MomentDelTagQueueData{
					EnterpriseId: enterpriseId,
					TaskId:       taskId,
				}
				planService.NewPlanTaskQueueSvc().MomentDelTagProduce(c, data)
				if err != nil {
					return err
				}
			}

		}
	}
	fmt.Printf("朋友圈任务脚本创建是否完成结果查询")
	return nil
}
