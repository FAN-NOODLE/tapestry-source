package v2

import (
	"encoding/json"
	"fmt"
	"github.com/chrislusf/glow/flow"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/app"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/sdk/work"
	"strings"
	"sync"
	"time"

	"github.com/emirpasic/gods/sets/hashset"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm/group_chat"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/xjob"
)

type taskCommunityPushData struct {
	baseJob
	name string
}

var _ runner = taskCommunityPushData{}

func NewTaskCommunityPushData() *taskCommunityPushData {
	return &taskCommunityPushData{
		name: "推送社群任务计划",
	}
}

func (job taskCommunityPushData) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	// 获取开启计划任务功能开启的企业信息
	var (
		err            error
		enterpriseList []boss.EnterpriseModel
	)
	n := time.Now()
	lockEight := time.Date(n.Year(), n.Month(), n.Day(), 8, 0, 0, 0, time.Local)
	if n.Before(lockEight) { // 任务推送在每天的8点之后开始推送
		return nil
	}
	enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, enterprise := range enterpriseList {
		wg.Add(1)
		go func(enterpriseID int64) {
			defer wg.Done()
			job.getCommunityTask(ctx, enterpriseID)
		}(enterprise.ID)
	}
	wg.Wait()

	return nil
}

func filterTaskConditionChat(chat *group_chat.GroupChatPoeple, task micro_scrm_enterprise.ResponCondition) bool {

	conditionDetail := micro_scrm_enterprise.ConditionDetail{}
	if err := json.Unmarshal([]byte(task.Conditions), &conditionDetail); err != nil {
		return false
	}
	if conditionDetail.Logical.Expression == "" || conditionDetail.Logical.Nubmer == 0 {
		return false
	}
	return dataProcessing(conditionDetail.Logical.Expression, chat.Number, conditionDetail.Logical.Nubmer)
}
func (job taskCommunityPushData) getCommunityTask(ctx ctx.Context, enterpriseID int64) {
	// 查询满足条件的社群群人数
	taskLists, err := orm.Service.PlanTask.GetPlanCommunityTask(enterpriseID, time.Now().Unix(), 2)
	if err != nil {
		return
	}
	chats, err := orm.Service.PlanTask.WorkGroupChatTotalNumberByenterpriseID(enterpriseID)
	if err != nil {
		return
	}
	pluginInfo, err := microBossDao.NewEnterpriseWorkPlugin().GetOne(ctx, map[string]interface{}{
		"enterprise_id": enterpriseID,
		"identifier":    work.QWCrsStoreAssistantSelf,
	})
	if err != nil {
		return
	}
	info, err := common.GetEnterpriseInfo(ctx, enterpriseID)
	if err != nil {
		return
	}
	for _, task := range taskLists {
		// taskCondMap := make(map[int64]micro_scrm_enterprise.ResponCondition)
		if task.TaskIsPush == 1 {
			continue
		}
		var taskAllotLogMap = make(map[string]micro_scrm_enterprise.TaskAllotLogModel)
		var userMap = make(map[string]int64)
		var userChatNumMap = make(map[string]int64)
		var sendTaskMap = make(map[int64][]string)
		taskId := task.TaskId
		_ = planService.NewPlanTask().UpdatePlanTaskPush(ctx, enterpriseID, taskId)
		specifyCrowdUserSet := hashset.New()
		if task.IsSpecifyCrowd == 2 {
			specifyCrowdUser, _ := orm.Service.PlanTask.PlanTasksIsSpecifyCrowdByTaskId(ctx, enterpriseID, taskId)
			if len(specifyCrowdUser) > 0 {
				for _, userid := range specifyCrowdUser {
					specifyCrowdUserSet.Add(userid)
				}
			}
		}
		taskConditionDataModelMap := make(map[string][]microScrmModel.TaskConditionDataModel, 0)
		for _, chat := range chats {
			key := chat.UserId
			if task.IsSpecifyCrowd == 2 && !specifyCrowdUserSet.Contains(chat.UserId) {
				continue
			}

			filter := filterTaskConditionChat(&chat, task)
			if !filter {
				continue
			}
			if _, ok := userMap[key]; ok {
				userMap[key] += chat.Number
				userChatNumMap[key]++
			} else {
				userMap[key] = chat.Number
				// 新增taskAllotLog表记录
				taskAllotLogMap[key] = micro_scrm_enterprise.TaskAllotLogModel{
					PlanId:      task.PlanId,
					ConditionId: task.ConditionId,
					TaskId:      taskId,
					UserId:      chat.UserId,
					StartTime:   int(time.Now().Unix()),
					EndTime:     0,
				}
				sendTaskMap[taskId] = append(sendTaskMap[taskId], chat.UserId)
				userChatNumMap[key] = 1
			}
			taskConditionDataModelMap[key] = append(taskConditionDataModelMap[key], microScrmModel.TaskConditionDataModel{
				BaseModel: models.BaseModel{
					ID: app.UID(),
				},
				UserId:         chat.UserId,
				ConditionId:    task.ConditionId,
				TaskId:         taskId,
				PlanId:         task.PlanId,
				TaskAllotLogId: 0,
				GroupChatId:    chat.ID,
			})
		}

		var taskAllotLogIdMap = make(map[string]int64)

		for index, v := range taskAllotLogMap {
			key := index
			info := v
			info.CustomerNum = int(userMap[key])
			info.ChatNum = int(userChatNumMap[key])
			logId, err := orm.Service.PlanTask.CreateAllotLog(enterpriseID, info)
			if err != nil {
				fmt.Printf("插入task_allot_log表失败")
				continue
			}
			taskAllotLogIdMap[index] = logId
		}
		for key, conditionData := range taskConditionDataModelMap {
			var taskConditionDataArr []microScrmModel.TaskConditionDataModel
			for _, v := range conditionData {
				item := v
				item.TaskAllotLogId = taskAllotLogIdMap[key]
				taskConditionDataArr = append(taskConditionDataArr, item)
			}
			err = planDao.NewPlanTaskConditionDataDao().InsertConditionData(ctx, enterpriseID, task.PlanId, taskConditionDataArr)
			if err != nil {
				log.Info(ctx).Msgf("插入task_condition_data表数据失败,err:%v", err)
				continue
			}
		}

		qyWorkApi := new(work.QyWorkImpl).SetDomain(info.Domain)
		for k, userIds := range sendTaskMap {
			planName := task.PlanName
			taskName := task.TaskName
			flow.New().Source(func(ch chan work.SendMessageRequest) {
				for _, userId := range userIds {
					key := userId
					chatNum := userChatNumMap[key]
					logId := taskAllotLogIdMap[key]
					nowDate := strings.Split(time.Now().Format("2006-01-02"), "-")
					var content []work.MiniprogramNoticeItem
					content = append(content, work.MiniprogramNoticeItem{
						Key:   "任务计划名称",
						Value: planName,
					})
					content = append(content, work.MiniprogramNoticeItem{
						Key:   "子任务名称",
						Value: taskName,
					})
					content = append(content, work.MiniprogramNoticeItem{
						Key:   "社群人数变动时间",
						Value: fmt.Sprintf("%s年%s月%s日", nowDate[0], nowDate[1], nowDate[2]),
					})
					content = append(content, work.MiniprogramNoticeItem{
						Key:   "说明",
						Value: "群人数变动触发任务条件,点击查看详情",
					})
					ch <- work.SendMessageRequest{
						QyGwRequest: work.QyGwRequest{ReqPL: work.QWCrsStoreAssistantSelf},
						Touser:      userId,
						Msgtype:     "miniprogram_notice",
						MiniprogramNotice: &work.MiniprogramNoticeContent{
							Appid:       pluginInfo.MiniAppID,
							Page:        fmt.Sprintf("/pagesA/task-detail/task-detail?id=%d&name=%s&num=%d", logId, taskName, chatNum),
							Title:       "客户群变动任务计划提醒", // 导入客户模板任务计划提醒
							ContentItem: content,
						},
					}
				}
			}, 1).Map(func(params work.SendMessageRequest) {
				response, err := qyWorkApi.SendMessage(ctx, params)
				if err != nil {
					log.Info(ctx).Msgf("发送小程序应用消息失败, err:%s", err.Error())
					return
				}
				if response.Errcode != 0 {
					log.Info(ctx).Msgf("发送小程序应用消息企微返回错误, response:%v", response)
				}
				log.Info(ctx).Err(err).Interface("params", params).Msg("qyGroupSend.qyGroupSend")
				msgId := response.Msgid
				fmt.Printf("taskId:%d----导购id:%s----任务发送小程序应用消息的--msgID:%s\n", k, params.Touser, msgId)
				err = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseID, map[string]interface{}{
					"task_id": k,
					"user_id": params.Touser,
				}, map[string]interface{}{"agent_msg_id": msgId, "agent_msg_time": time.Now().Unix()})
				log.Info(ctx).Err(err).Interface("params", params).Msg("sendPlanTaskMiniMessage.SendMiniMessageRequest")
			}).Run()
			err := orm.Service.PlanTask.UpdateTaskPush(enterpriseID, k)
			if err != nil {
				log.Info(ctx).Err(err).Msg("更新社群任务推送失败")
				continue
			}
		}
	}
	return
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
