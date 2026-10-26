package plan_tasks

import (
	"fmt"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/xjob"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type GroupMsgResult struct {
	Name string
}

func NewGroupMsgResult() GroupMsgResult {
	return GroupMsgResult{
		Name: "计划任务企业群发的结果查询",
	}
}

func (job GroupMsgResult) Run(c ctx.Context, handler *xjob.JobHandler) error {
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

func (job GroupMsgResult) handle(c ctx.Context, enterpriseId int64) error {
	var err error
	// 查询活动有效期内进行中的企业群发的子任务
	planIds, err := planService.NewPlanTask().GetPlanIdByStatus(c, enterpriseId, microScrmModel.PlanIng)
	if err != nil {
		return err
	}
	taskLists, err := planService.NewPlanTask().GetQyGroupListsByPlanId(c, enterpriseId, microScrmModel.PLANTASKQYSEND, planIds)
	if err != nil {
		return err
	}
	var taskIds []int64
	taskMap := make(map[int64]int)
	for _, list := range taskLists {
		item := list
		taskIds = append(taskIds, item.ID)
		taskMap[item.ID] = item.OriginType
	}
	taskIds = helper.ArrayUniq(taskIds)
	if len(taskIds) == 0 {
		return err
	}
	// 查询task_allot_log表中企业群发状态为未发送的导购信息
	sendUserInfo := make(map[int64][]microScrmModel.TaskAllotLogModel, 0)
	for _, v := range taskIds {
		taskId := v
		lists, err := planService.NewTaskAllotLog().GetListsByTaskId(c, enterpriseId, taskId)
		if err != nil {
			continue
		}
		sendUserInfo[taskId] = lists
	}
	enterpriseInfo, err := common.GetEnterpriseInfo(c, enterpriseId)
	if err != nil {
		return err
	}
	memberLists, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"deleted_at":    0,
		"status":        1,
	})
	if err != nil {
		return err
	}
	var userIds []string
	for _, list := range memberLists {
		item := list
		userIds = append(userIds, item.Userid)
	}
	qyWorkApi := new(work.QyWorkImpl).SetDomain(enterpriseInfo.Domain)
	// 调用企微的发送结果结果，根据返回结果中的值，如果有一条的结果值大于0，则修改task_allot_log表中的is_group_send状态为1，已推送
	wg := sync.WaitGroup{}
	p, _ := ants.NewPool(5)
	defer p.Release()
	for k, lists := range sendUserInfo {
		item := lists
		if len(item) == 0 {
			continue
		}
		taskId := k
		for _, v := range item {
			info := v
			if !helper.InArray(info.UserId, userIds) {
				continue
			}
			wg.Add(1)
			fn := func(ctx ctx.Context, enterpriseId, taskId int64, allotLog microScrmModel.TaskAllotLogModel, qyWorkApi *work.QyWorkImpl, taskMap map[int64]int) func() {
				return func() {
					defer wg.Done()
					planService.NewPlanTask().QyGroupSendResult(c, enterpriseId, taskId, allotLog, qyWorkApi, taskMap)
				}
			}
			err := p.Submit(fn(c, enterpriseId, taskId, info, qyWorkApi, taskMap))
			time.Sleep(200 * time.Millisecond) // 防止企微接口限频
			if err != nil {
				fmt.Printf("异步处理企业群发结果失败，err:%v", err)
				continue
			}
		}

	}
	wg.Wait()
	fmt.Printf("执行企微群发任务信息结果查询完成")
	return nil
}
