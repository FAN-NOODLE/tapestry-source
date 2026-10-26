package community_task

import (
	"encoding/json"
	"errors"
	"fmt"
	"micro-scrm/internal/app"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	messageSvc "micro-scrm/internal/app/service/message"
	"strings"
	"sync"
	"time"

	planTaskSftp "micro-scrm/internal/orm/service/plan_task"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/wx/wxoffiaccount/cache"
	"micro-scrm/internal/config"
	"micro-scrm/internal/events/entity"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	micro_redis "micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk/microwork"
)

/**
*添加企业客户事件
 */
func ChangeExternalChatAddMember(ctx ctx.Context, data entity.ExternalUserInfoEntity) error {
	defer func() {
		if err := recover(); err != nil {
			log.Info(ctx).Msgf("定时推送公司编号%v的任务的任务:%v", data.EnterpriseID, err)
		}
	}()
	// 查询该企业是否开通了计划任务
	result, err := orm.Service.BrandConfig.GetSFTPConfig([]int64{data.EnterpriseID}, planTaskSftp.SFTP_TYPE_PLAN_TASK)
	if err != nil || len(result) == 0 {
		return errors.New("企业未开通计划任务")
	}
	// 查询群内新增人数的任务任务
	if planTasks, err := orm.Service.PlanTask.GetPlanCommunityTask(data.EnterpriseID, time.Now().Unix(), 3); err == nil {
		if len(planTasks) > 0 {
			var taskWG sync.WaitGroup
			for _, task := range planTasks {
				taskWG.Add(1)
				go func(newtask micro_scrm_enterprise.ResponCondition) {
					defer taskWG.Done()
					if err := handworkGroupChatAddNumber(ctx, data, newtask); err != nil {
						log.Info(ctx).Msgf("企业信息为%v的社群任务%v的群内新增人数条件%v执行失败： %v", data.EnterpriseID, newtask.TaskId, newtask.ConditionId, err)
					}
				}(task)
			}
			taskWG.Wait()
		}
	} else {
		log.Info(ctx).Msgf("公司编号:%v,获取社区任务群内新增人数的任务失败:%v", data.EnterpriseID, err)
	}
	// 群添加人后查询群人数是否满足计划任务的条件
	// HandworkGroupChatTotalNumber(ctx, data)
	return nil
}
func ChangeExternalChatDeleteMember(ctx ctx.Context, data entity.ExternalUserInfoEntity) error {
	defer func() {
		if err := recover(); err != nil {
			log.Info(ctx).Msgf("定时推送公司编号%v的任务的任务:%v", data.EnterpriseID, err)
		}
	}()
	// 查询该企业是否开通了计划任务
	result, err := orm.Service.BrandConfig.GetSFTPConfig([]int64{data.EnterpriseID}, planTaskSftp.SFTP_TYPE_PLAN_TASK)
	if err != nil || len(result) == 0 {
		return errors.New("企业未开通计划任务")
	}
	// 查询群内删减人数的任务任务
	if planTasks, err := orm.Service.PlanTask.GetPlanCommunityTask(data.EnterpriseID, time.Now().Unix(), 7); err == nil {
		if len(planTasks) > 0 {
			var taskWG sync.WaitGroup
			for _, task := range planTasks {
				taskWG.Add(1)
				go func(newtask micro_scrm_enterprise.ResponCondition) {
					defer taskWG.Done()
					if err := handworkGroupChatrefundNumber(ctx, data, newtask); err != nil {
						log.Info(ctx).Msgf("企业信息为%v的社群任务%v的群内新增人数条件%v执行失败： %v", data.EnterpriseID, newtask.TaskId, newtask.ConditionId, err)
					}
				}(task)
			}
			taskWG.Wait()
		}
	} else {
		log.Info(ctx).Msgf("公司编号:%v,获取社区任务群内新增人数的任务失败:%v", data.EnterpriseID, err)
	}
	// 群减员后查询群人数是否满足计划任务的条件
	// HandworkGroupChatTotalNumber(ctx, data)
	return nil
}

// 查询群人数的任务任务
func HandworkGroupChatTotalNumber(ctx ctx.Context, data entity.ExternalUserInfoEntity) {
	if planTasks, err := orm.Service.PlanTask.GetPlanCommunityTask(data.EnterpriseID, time.Now().Unix(), 2); err == nil {
		if len(planTasks) > 0 {
			var taskWG sync.WaitGroup
			for _, task := range planTasks {
				taskWG.Add(1)
				go func(newtask micro_scrm_enterprise.ResponCondition) {
					defer taskWG.Done()
					if err := handworkGroupChatNumber(ctx, data, newtask); err != nil {
						log.Info(ctx).Msgf("企业信息为%v的社群任务%v的群内新增人数条件%v执行失败： %v", data.EnterpriseID, newtask.TaskId, newtask.ConditionId, err)
					}
				}(task)
			}
			taskWG.Wait()
		}
	} else {
		log.Info(ctx).Msgf("公司编号:%v,获取ctx ctx.Context, data entity.ExternalUserCallbackEntity社区任务群内新增人数的任务失败:%v", data.EnterpriseID, err)

	}
}
func handworkGroupChatNumber(ctx ctx.Context, data entity.ExternalUserInfoEntity, task micro_scrm_enterprise.ResponCondition) error {
	// 将未开启的任务开启
	if task.TaskStatus == 0 {
		orm.Service.PlanTask.UpdateplanTaskStatusByID(data.EnterpriseID, task.PlanId)
	}
	if task.PlanStatus == 0 {
		orm.Service.PlanTask.UpdateplanStatusByID(data.EnterpriseID, task.PlanId)
	}
	conditions := micro_scrm_enterprise.ConditionDetail{}
	if err := json.Unmarshal([]byte(task.Conditions), &conditions); err != nil {
		return err
	}
	if conditions.Logical.Expression == "" || conditions.Logical.Nubmer == 0 {
		return errors.New("错误的表达式")
	}
	/*
	*判断task是否指定人群 task.IsSpecifyCrowd=2表示指定人群,否则未指定人群
	 */
	specifyCrowdUser := []string{}
	if task.IsSpecifyCrowd == 2 {
		specifyCrowdUserFlag := false
		specifyCrowdUser, specifyCrowdUserFlag = getSpecifyCrowdUser(ctx, data.EnterpriseID, task.TaskId, task.EndTime, task.TaskType)
		if !specifyCrowdUserFlag {
			log.Info(ctx).Msgf("企业信息为%v的社群任务%v:指定人群但未获取到", data.EnterpriseID, task.TaskId)
			return errors.New("计划任务指定人群但未获取到")
		}
	}
	if result, err := orm.Service.PlanTask.WorkGroupChatTotalNumber(data.EnterpriseID, data.GroupChatID, specifyCrowdUser); err != nil {
		return err
	} else {
		if result.ID != 0 {
			if isFlag := dataProcessing(conditions.Logical.Expression, result.Number, conditions.Logical.Nubmer); isFlag {
				go createLog(ctx, data, task, result.UserId, result.Number)
			}
		}
	}
	return nil
}
func handworkGroupChatAddNumber(ctx ctx.Context, data entity.ExternalUserInfoEntity, task micro_scrm_enterprise.ResponCondition) error {
	// 将未开启的任务开启
	if task.TaskStatus == 0 {
		orm.Service.PlanTask.UpdateplanTaskStatusByID(data.EnterpriseID, task.PlanId)
	}
	if task.PlanStatus == 0 {
		orm.Service.PlanTask.UpdateplanStatusByID(data.EnterpriseID, task.PlanId)
	}
	conditions := micro_scrm_enterprise.ConditionDetail{}
	if err := json.Unmarshal([]byte(task.Conditions), &conditions); err != nil {
		return err
	}
	if conditions.Logical.Expression == "" || conditions.Logical.Nubmer == 0 {
		return errors.New("错误的表达式")
	}
	/*
	*判断task是否指定人群 task.IsSpecifyCrowd=2表示指定人群,否则未指定人群
	 */
	specifyCrowdUser := []string{}
	if task.IsSpecifyCrowd == 2 {
		specifyCrowdUserFlag := false
		specifyCrowdUser, specifyCrowdUserFlag = getSpecifyCrowdUser(ctx, data.EnterpriseID, task.TaskId, task.EndTime, task.TaskType)
		if !specifyCrowdUserFlag {
			log.Info(ctx).Msgf("企业信息为%v的社群任务%v:指定人群但未获取到", data.EnterpriseID, task.TaskId)
			return errors.New("计划任务指定人群但未获取到")
		}
	}

	if result, err := orm.Service.PlanTask.WorkGroupChatAddNumber(data.EnterpriseID, data.GroupChatID, task.StartTime, task.EndTime, specifyCrowdUser); err != nil {
		return err
	} else {
		if result.ID != 0 {
			if isFlag := dataProcessing(conditions.Logical.Expression, result.Number, conditions.Logical.Nubmer); isFlag {
				go createLog(ctx, data, task, result.UserId, result.Number)
			}
		}
	}
	return nil
}
func handworkGroupChatrefundNumber(ctx ctx.Context, data entity.ExternalUserInfoEntity, task micro_scrm_enterprise.ResponCondition) error {
	// 将未开启的任务开启
	if task.TaskStatus == 0 {
		orm.Service.PlanTask.UpdateplanTaskStatusByID(data.EnterpriseID, task.PlanId)
	}
	if task.PlanStatus == 0 {
		orm.Service.PlanTask.UpdateplanStatusByID(data.EnterpriseID, task.PlanId)
	}
	conditions := micro_scrm_enterprise.ConditionDetail{}
	if err := json.Unmarshal([]byte(task.Conditions), &conditions); err != nil {
		return err
	}
	if conditions.Logical.Expression == "" || conditions.Logical.Nubmer == 0 {
		return errors.New("错误的表达式")
	}
	/*
	*判断task是否指定人群 task.IsSpecifyCrowd=2表示指定人群,否则未指定人群
	 */
	specifyCrowdUser := []string{}
	if task.IsSpecifyCrowd == 2 {
		specifyCrowdUserFlag := false
		specifyCrowdUser, specifyCrowdUserFlag = getSpecifyCrowdUser(ctx, data.EnterpriseID, task.TaskId, task.EndTime, task.TaskType)
		if !specifyCrowdUserFlag {
			log.Info(ctx).Msgf("企业信息为%v的社群任务%v:指定人群但未获取到", data.EnterpriseID, task.TaskId)
			return errors.New("计划任务指定人群但未获取到")
		}
	}
	if result, err := orm.Service.PlanTask.WorkGroupChatRefundNumber(data.EnterpriseID, data.GroupChatID, task.StartTime, task.EndTime, specifyCrowdUser); err != nil {
		return err
	} else {
		if result.ID != 0 {
			if isFlag := dataProcessing(conditions.Logical.Expression, result.Number, conditions.Logical.Nubmer); isFlag {
				go createLog(ctx, data, task, result.UserId, result.Number)

			}
		}
	}
	return nil
}
func createLog(ctx ctx.Context, data entity.ExternalUserInfoEntity, task micro_scrm_enterprise.ResponCondition, userId string, number int64) {
	if userId == "" {
		return
	}
	taskAllotLog := micro_scrm_enterprise.TaskAllotLogModel{
		PlanId:      task.PlanId,
		ConditionId: task.ConditionId,
		TaskId:      task.TaskId,
		UserId:      userId,
		StartTime:   int(time.Now().Unix()),
		ChatNum:     1,
	}
	var logid int64
	// 查询该任务是否给这个人已经推送过
	taskdata, err := orm.Service.PlanTask.FindShequnTaskPushIS(data.EnterpriseID, taskAllotLog)
	if err == nil && taskdata.ID != 0 {
		logid = taskdata.ID
		// 更新时间
		orm.Service.PlanTask.UpdatetaskAllotLogTime(data.EnterpriseID, taskdata.ID)
		return
	} else {
		logid, err = orm.Service.PlanTask.CreatetaskAllotLog(data.EnterpriseID, taskAllotLog)
		if err != nil {
			return
		}
	}
	conditionData := microScrmModel.TaskConditionDataModel{
		BaseModel: models.BaseModel{
			ID: app.UID(),
		},
		UserId:         userId,
		ConditionId:    task.ConditionId,
		TaskId:         task.TaskId,
		PlanId:         task.PlanId,
		TaskAllotLogId: logid,
		GroupChatId:    data.GroupChatID,
	}

	err = planDao.NewPlanTaskConditionDataDao().CreateConditionData(ctx, data.EnterpriseID, task.PlanId, conditionData)
	if err != nil {
		return
	}

	nowDate := strings.Split(time.Now().Format("2006-01-02"), "-")
	content := make([]map[string]string, 0)
	content = append(content, map[string]string{"key": "任务计划名称", "value": task.PlanName})
	content = append(content, map[string]string{"key": "子任务名称", "value": task.TaskName})
	content = append(content, map[string]string{"key": "社群人数变动时间", "value": fmt.Sprintf("%s年%s月%s日", nowDate[0], nowDate[1], nowDate[2])})
	content = append(content, map[string]string{"key": "说明", "value": "群人数变动触发任务条件，点击查看详情"})

	params := microwork.SendMiniMessageRequest{
		EnterpriseID: data.EnterpriseID,
		MemUserIds:   []string{userId},
		Title:        "客户群变动任务计划提醒",
		Content:      content,
		Page:         fmt.Sprintf("/pagesA/task-detail/task-detail?id=%d&name=%s&num=%d", logid, task.TaskName, number),
	}
	_, err = messageSvc.NewMessageService().SendMiniMessage(ctx, params)
	log.Info(ctx).Err(err).Interface("params", params).Msg("createLog.SendMiniMessageRequest")
	orm.Service.PlanTask.UpdateTaskPush(data.EnterpriseID, task.TaskId)
}
func dataProcessing(expression string, resultNum, expressionNum int64) bool {
	isFlag := false
	switch expression {
	case "eq":
		if resultNum == expressionNum {
			isFlag = true
		}
	case "gt":
		if resultNum > expressionNum {
			isFlag = true
		}
	case "ge":
		if resultNum >= expressionNum {
			isFlag = true
		}
	case "lt":
		if resultNum < expressionNum {
			isFlag = true
		}
	case "le":
		if resultNum <= expressionNum {
			isFlag = true
		}
	default:
		isFlag = false
	}
	return isFlag
}
func getSpecifyCrowdUser(ctx ctx.Context, enterpriseID, taskId, endTime int64, taskType int) ([]string, bool) {
	var err error
	specifyCrowdUserFlag := false
	specifyCrowdUser := getPlanTaskRids(ctx, enterpriseID, taskId)
	if len(specifyCrowdUser) > 0 { // 说明缓存存在，否则不存在
		specifyCrowdUserFlag = true
		return specifyCrowdUser, specifyCrowdUserFlag
	}
	// 如果缓存不存在查询数据库
	if specifyCrowdUser, err = orm.Service.PlanTask.PlanTasksIsSpecifyCrowdByTaskId(ctx, enterpriseID, taskId); err == nil {
		if len(specifyCrowdUser) > 0 && taskType == 2 { // 如果是持续任务了重新设置缓存
			setPlanTaskRids(ctx, enterpriseID, taskId, endTime, 2, specifyCrowdUser)
		}
	}
	if len(specifyCrowdUser) > 0 {
		specifyCrowdUserFlag = true
	}
	return specifyCrowdUser, specifyCrowdUserFlag
}
func getPlanTaskRids(ctx ctx.Context, enterpriseID, taskId int64) []string {
	planTaskSpecifyCrowds := make([]string, 0)
	key := fmt.Sprintf("%s::%d::plan_task::specify_crowd::%d", config.CfgData.App, enterpriseID, taskId)
	redisCache := cache.NewRedisCache(micro_redis.Client())
	redisByte := []byte{}
	redisCache.Get(ctx, key, &redisByte)
	json.Unmarshal(redisByte, &planTaskSpecifyCrowds)
	return planTaskSpecifyCrowds
}
func setPlanTaskRids(ctx ctx.Context, enterpriseID, taskId, endTime int64, isSpecifyCrowd int, userIds []string) {
	key := fmt.Sprintf("%s::%d::plan_task::specify_crowd::%d", config.CfgData.App, enterpriseID, taskId)
	redisCache := cache.NewRedisCache(micro_redis.Client())
	if isSpecifyCrowd == 1 {
		redisCache.Set(ctx, key, "", 0)
		return
	}
	redisTime := endTime - time.Now().Unix()
	if redisTime <= 0 {
		return
	}
	jsonData, jsonErr := json.Marshal(userIds)
	if jsonErr != nil {
		return
	}
	redisCache.Set(ctx, key, jsonData, redisTime)
}
