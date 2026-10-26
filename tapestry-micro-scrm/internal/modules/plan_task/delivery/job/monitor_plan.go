package job

import (
	"fmt"
	"log"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	workDao "micro-scrm/internal/app/dao/work"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	marketingService "micro-scrm/internal/app/service/marketing"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/sdk/work"
	"regexp"
	"sync"
	"time"

	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/xjob"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// TurnOnPlan 计划任务开启
func TurnOnPlan(ctx ctx.Context, jobHandler *xjob.JobHandler) error {
	data, _ := orm.BossService.Enterprise.GetAllInfo()
	var wg sync.WaitGroup
	for _, v := range data {
		wg.Add(1)
		go func(enterpriseId int64) {
			defer wg.Done()
			StartPaln(ctx, enterpriseId)
		}(v.ID)
	}
	wg.Wait()
	return nil
}

func StartPaln(c ctx.Context, enterpriseId int64) { // 开始任务以子任务的时间来判断
	result, err := orm.Service.PlanTask.GetPlanStartTime(micro_scrm_enterprise.PlanModel{}, enterpriseId)
	if err != nil {
		log.Println(err)
		return
	}

	var idList = make([]int64, 0)
	for _, v := range result {
		idList = append(idList, v.ID)
	}
	if len(idList) > 0 {
		err = microScrmDao.NewPlanTaskDao().UpdatePlanData(c, enterpriseId, map[string]interface{}{
			"id": idList,
		}, map[string]interface{}{
			"status": microScrmModel.PlanIng,
		})
		if err != nil {
			return
		}
	}
	_ = microScrmDao.NewPlanTaskDao().UpdateTaskStatus(c, enterpriseId, microScrmModel.PlanIng, "task_start_time<=? and task_end_time>? and deleted_at=0 and status=?", time.Now().Unix(), time.Now().Unix(), microScrmModel.PlanNotStart)
}

// ClosurePlan 计划结束
func ClosurePlan(ctx ctx.Context, jobHandler *xjob.JobHandler) error {
	data, _ := orm.BossService.Enterprise.GetAllInfo()
	var wg sync.WaitGroup
	for _, v := range data {
		wg.Add(1)
		go func(Id int64) {
			defer wg.Done()
			closeTask(ctx, Id)
		}(v.ID)
		// result, err := orm.Service.PlanTask.GetPlanStartTime(micro_scrm_enterprise.PlanModel{}, v.ID)
	}
	wg.Wait()
	return nil
}

func closeTask(ctx ctx.Context, ID int64) error {
	result, err := orm.Service.PlanTask.GetPerformPlan(micro_scrm_enterprise.PlanModel{}, ID)
	if err != nil {
		log.Println(err)
		return err
	}

	var planIds = make([]int64, 0)

	for _, v := range result {
		planIds = append(planIds, v.ID)
	}
	err = microScrmDao.NewPlanTaskDao().UpdatePlanData(ctx, ID, map[string]interface{}{
		"id": planIds,
	}, map[string]interface{}{
		"status": microScrmModel.PlanEnd,
	})
	if err != nil {
		return err
	}
	var momentTaskIds, contactTaskIds, tagTaskIds []int64
	// 查询所有已过期的子任务
	taskLists, err2 := microScrmDao.NewPlanTaskDao().GetMany(ctx, ID, "task_end_time<? and status=?", time.Now().Unix(), microScrmModel.PlanIng)
	if err2 != nil {
		return err2
	}
	taskId := make([]int64, 0)
	for _, v := range taskLists {
		if v.GroupSendType == microScrmModel.PLANTASKQYSEND {
			if v.TaskCategory == microScrmModel.PLANTASKMOMENT {
				momentTaskIds = append(momentTaskIds, v.ID)
			}
			if v.TaskCategory == microScrmModel.PLANTASKCONTACT {
				contactTaskIds = append(contactTaskIds, v.ID)
			}
		}
		taskId = append(taskId, v.ID)
	}
	err = microScrmDao.NewPlanTaskDao().UpdateTaskStatus(ctx, ID, microScrmModel.PlanEnd, map[string]interface{}{
		"id": taskId,
	})
	if err != nil {
		return err
	}
	err = microScrmDao.NewPlanTaskAllotLogDao().UpdateData(ctx, ID, map[string]interface{}{
		"is_finished": microScrmModel.PLAN_FISISHING,
		"task_id":     taskId,
	}, map[string]interface{}{
		"is_finished": microScrmModel.PLAN_NOFINISH,
	})
	if err != nil {
		return err
	}
	// 更新有关贺卡的任务关系
	if err := marketingService.NewCustomCard().UpdatePlanCardInfo(ctx, ID, planIds); err != nil {
		return err
	}
	taskTagLists, err2 := microScrmDao.NewPlanTaskDao().GetMany(ctx, ID, map[string]interface{}{
		"plan_id": planIds,
	})
	for _, v := range taskTagLists {
		tagTaskIds = append(tagTaskIds, v.ID)
	}
	if len(tagTaskIds) > 0 {
		err := workDao.NewPlanTasksTagRelation().DeleteData(ctx, map[string]interface{}{
			"enterprise_id": ID,
			"plan_id":       tagTaskIds,
		})
		if err != nil {
			return err
		}
	}
	if len(momentTaskIds) > 0 {

		dbDSN := config.CfgData.DB.MasterDSN
		// 正则表达式匹配
		re := regexp.MustCompile(`^([^:]+):([^@]+)@(?:tcp\()?([^:]+):([0-9]+)\)?(?:/([^?]+))?(?:\?(.*))?$`)
		matches := re.FindStringSubmatch(dbDSN)
		var host, port, username, password string
		if matches != nil && len(matches) >= 5 {
			username = matches[1]
			password = matches[2]
			host = matches[3]
			port = matches[4]
		}
		// 停止未发送的朋友圈任务
		info, err := common.GetEnterpriseInfo(ctx, ID)
		if err != nil {
			return err
		}
		qyWorkApi := new(work.QyWorkImpl).SetDomain(info.Domain)
		var cacelMomentReq work.CancelMomentTaskRequest
		var taskIds []int64
		for _, v := range momentTaskIds {
			taskIds = append(taskIds, v)
		}
		taskLists, err2 := microScrmDao.NewPlanTaskDao().GetMany(ctx, ID, map[string]interface{}{
			"id":            taskIds,
			"moment_status": 3,
		})
		if err2 != nil {
			return err2
		}
		for _, v := range taskLists {
			taskInfo := v
			cacelMomentReq.MomentId = taskInfo.MomentJobId
			_, err = qyWorkApi.CancelMomentTask(ctx, cacelMomentReq)
			if host != "" { // 朋友圈任务结束的时候，再次同步ck数据
				sql := fmt.Sprintf("INSERT INTO micro_scrm_enterprise_%d_ck.task_condition_data_local "+
					" SELECT id,task_allot_log_id,tag_name,ww_tag_id,union_id,work_external_user_id,is_union,mobile,mars_id,user_id,condition_id,task_id,plan_id,created_at,updated_at,deleted_at,customer_id,is_finished,finished_at,remark,is_contact,is_call,is_add,is_match,is_touch,user_name,group_chat_id,customer_status,forward FROM"+
					" mysql('%s:%s', 'micro_scrm_enterprise_%d', 'task_condition_data_%d', '%s', '%s') where task_id=%d", ID, host, port, ID, taskInfo.PlanID, username, password, taskInfo.ID)
				err = gormx.GetCkInstance().Exec(sql).Error
			}
			if err != nil {
				continue
			}

		}

	}

	if len(contactTaskIds) > 0 {
		// 停止未发送的企业群发任务
		many, err := microScrmDao.NewPlanTaskAllotLogDao().GetMany(ctx, ID, map[string]interface{}{
			"task_id":       contactTaskIds,
			"is_group_send": microScrmModel.QyGroupMsgSend,
		})
		if err != nil {
			return err
		}
		info, err := common.GetEnterpriseInfo(ctx, ID)
		if err != nil {
			return err
		}
		qyWorkApi := new(work.QyWorkImpl).SetDomain(info.Domain)
		var request work.CancelGroupMsgSendRequest
		var cancelIds []int64
		for _, v := range many {
			info := v
			if info.GroupMsgId == "" {
				continue
			}
			request.MsgId = info.GroupMsgId
			_, err = qyWorkApi.CancelGroupMsgSend(ctx, request)
			if err != nil {
				fmt.Printf("停止企业群发报错，err:%v\n", err)
				continue
			}
			cancelIds = append(cancelIds, v.ID)
		}
		err = microScrmDao.NewPlanTaskAllotLogDao().UpdateData(ctx, ID, map[string]interface{}{
			"id": cancelIds,
		}, map[string]interface{}{
			"is_group_send": microScrmModel.QyGroupMsgSendCancel,
		})

	}

	// startTime, endTime := helper.PeriodOfTimestamp()
	// startTime = startTime - 60*60*24
	// endTime = endTime - 60*60*24
	// results, err := orm.Service.PlanTask.GetOncePushInfo(ID, danci, micro_scrm_enterprise.TaskConditionDataModel{})
	// if err != nil {
	// 	return err
	// }
	// newClient, _ := workwechat.NewWorkWeChat(ctx, ID)
	// for _, v := range results {
	// 	if err := orm.Service.PlanTask.UpdateTagStatus(ID, v.ExternalTagId, v.RelationId); err != nil {
	// 		return err
	// 	}
	// 	newClient.ManagerCustomerTag(work.UserTagHandlerRequest{
	// 		Userid:         v.UserId,
	// 		ExternalUserid: v.WorkExternalUserId,
	// 		RemoveTag:      []string{v.TagId},
	// 	})
	// }
	// result2, err := orm.Service.PlanTask.UpdateYesterdayInfo(ID, startTime, endTime, 1, micro_scrm_enterprise.TaskConditionDataModel{})
	// if err != nil {
	// 	return err
	// }
	// newClient2, _ := workwechat.NewWorkWeChat(ctx, ID)
	// for _, v := range result2 {
	// 	if err := orm.Service.PlanTask.UpdateTagStatus(ID, v.ExternalTagId, v.RelationId); err != nil {
	// 		return err
	// 	}
	// 	newClient2.ManagerCustomerTag(work.UserTagHandlerRequest{
	// 		Userid:         v.UserId,
	// 		ExternalUserid: v.WorkExternalUserId,
	// 		RemoveTag:      []string{v.TagId},
	// 	})
	// }

	return nil
}
