package plan_task

import (
	"database/sql"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	planReq "micro-scrm/internal/app/http/request/plan"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"strings"
	"sync"
	"time"

	"micro-scrm/internal/modules/members"

	"micro-scrm/internal/pkg/assist"

	"micro-scrm/internal/modules/ws_trace"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

// PlanTaskService 计划任务服务
type PTaskService struct {
	repo Repository
}

// NewService 初始化服务
func NewService(repo Repository) *PTaskService {
	return &PTaskService{
		repo: repo,
	}
}

// 获取指定日期的统计数据
func (s *PTaskService) SelectStatisticsByDay(c ctx.Context, enterpriseID int64, date string) (map[int64]TaskEverydayStatistics, error) {
	statistics := make(map[int64]TaskEverydayStatistics)
	list, _ := s.repo.SelectEverydayStatisticsByDay(c, enterpriseID, date)
	if len(list) > 0 {
		for _, v := range list {
			statistics[v.TaskID] = v
		}
	}
	return statistics, nil
}

func (s *PTaskService) GetTaskStatistics(c ctx.Context, enterpriseID int64, taskID int64, yesterdayTs int64) (*TaskEverydayStatistics, error) {
	return s.repo.GetTaskEverydayStatistics(c, enterpriseID, taskID, yesterdayTs)
}

type taskTmp struct {
	CustomerNum int64
	IsMain      bool
}

func (s *PTaskService) GetSubTaskStatistics(c ctx.Context, enterpriseID int64, memberID int64, chatID int64, mainTaskID int64, childTaskID int64, subTaskID int64, taskCategory int, subTaskStatistic *TaskEverydayStatistics) (*TaskEverydayStatistics, error) {
	// 查询触达人数统计数据
	if chatID > 0 {
		subTaskStatistic.TotalReachCustomerNumber, _ = s.repo.GetSubTasksChatCustomerNumber(c, enterpriseID, chatID)
	} else {
		subTaskStatistic.TotalReachCustomerNumber, _ = s.repo.GetSubTasksPortraitCustomerNumber(c, enterpriseID, subTaskStatistic.TaskID, mainTaskID, taskCategory)
	}
	return subTaskStatistic, nil
}

func (s *PTaskService) StatisticsOrderData(c ctx.Context, corp *enterprise.Corp, wsTraceService *ws_trace.Service, memberService *members.Service) {
	_, endTime := helper.PeriodOfTimestamp()
	// todo 结束时间测试时写成当天
	// yesterdayTs := endTime
	yesterdayTs := endTime - 24*60*60
	date := time.Unix(yesterdayTs, 0).Format("2006-01-02")
	dataTs := yesterdayTs
	planIDm := make(map[int64]int64)             // 主计划任务的子计划任务id集合
	planChatIDTasks := make(map[int64][]int64)   // 主计划任务的子计划群任务id集合
	planPortraitTasks := make(map[int64][]int64) // 主计划任务的子计划客户画像任务id集合
	taskIDM := make(map[int64][]int64)           // 任务的群聊id集合
	chatM := make(map[int64]int64)

	taskCustomerNum := make(map[int64]taskTmp)

	// 1、查找截至到今日凌晨为止且没有过期计划任务中完成的任务
	list, err := s.repo.GetTasksFinishDay(c, corp.EnterpriseID, yesterdayTs)
	if err != nil {
		return
	}
	for _, v := range list {
		tmp := TaskEverydayStatistics{
			Date:     date,
			DataTime: dataTs,
			TaskID:   v.ID, // task_allot_log表的id
			TaskType: 3,
		}
		// 根据userid 查询member_id
		memberID, err := memberService.GetMemberId(ctx.Wrap(c), corp.EnterpriseID, v.UserId)
		if err != nil || memberID == 0 {
			log.Err(c, err).Msgf("查询员工错误:%+v", v)
			continue
		}
		one, err := planDao.NewPlanTaskDao().GetOne(c, corp.EnterpriseID, map[string]interface{}{
			"id": v.TagId,
		})
		if err != nil {
			log.Err(c, err).Msgf("查询子任务错误:%+v", v)
			continue
		}
		subTaskSta, _ := s.GetSubTaskStatistics(c, corp.EnterpriseID, memberID, v.GroupChatId, v.PlanId, v.TaskId, v.ID, one.TaskCategory, &tmp)
		// 保存入库
		if subTaskSta != nil {
			err = s.repo.InsertTaskEveryDayStatistic(c, corp.EnterpriseID, *subTaskSta, false)
			if err != nil {
				log.Err(c, err).Msgf("查询截至到昨天重复客户去重统计数据保存到数据库错误:%+v", tmp)
			}
		}

		if _, ok := planIDm[v.PlanId]; !ok {
			planIDm[v.PlanId] = v.PlanId
			taskIDM[v.TaskId] = []int64{}
			if v.GroupChatId > 0 {
				taskIDM[v.TaskId] = append(taskIDM[v.TaskId], v.GroupChatId)
				chatM[v.GroupChatId] = v.GroupChatId
				planChatIDTasks[v.PlanId] = append(planChatIDTasks[v.PlanId], v.TaskId)
			} else {
				planPortraitTasks[v.PlanId] = append(planPortraitTasks[v.PlanId], v.TaskId)
			}
		} else {
			if _, ok := taskIDM[v.TaskId]; !ok {
				taskIDM[v.TaskId] = []int64{}
				if v.GroupChatId > 0 {
					planChatIDTasks[v.PlanId] = append(planChatIDTasks[v.PlanId], v.TaskId)
				} else {
					planPortraitTasks[v.PlanId] = append(planPortraitTasks[v.PlanId], v.TaskId)
				}
			}
			if v.GroupChatId > 0 {
				taskIDM[v.TaskId] = append(taskIDM[v.TaskId], v.GroupChatId)
				chatM[v.GroupChatId] = v.GroupChatId
			}
		}

	}

	planIds, err := planDao.NewPlanTaskDao().GetPlanMany(c, corp.EnterpriseID, "status=1 and deleted_at=0")
	if err != nil {
		return
	}

	for _, v := range planIds {
		planId := v.ID
		taskIds, err := planDao.NewPlanTaskDao().GetMany(c, corp.EnterpriseID, map[string]interface{}{
			"plan_id":    planId,
			"deleted_at": 0,
		})
		if err != nil {
			return
		}
		for _, task := range taskIds {
			taskId := task.ID
			customerNum, err := planDao.NewPlanTaskConditionDataDao().GetGroupCount(c, corp.EnterpriseID, planId, map[string]interface{}{
				"task_id":     taskId,
				"is_finished": microScrmModel.FinishedYes,
			})
			if err != nil {
				return
			}
			tmp := TaskEverydayStatistics{
				Date:                     date,
				DataTime:                 dataTs,
				TaskID:                   taskId,
				TotalReachCustomerNumber: customerNum,
				TaskType:                 2,
			}
			// 保存入库
			err = s.repo.InsertTaskEveryDayStatistic(c, corp.EnterpriseID, tmp, false)

			if err != nil {
				continue
			}

		}
	}
	return

	for k := range planIDm {
		allUnionsM := make(map[string]string)
		taskUnionM := make(map[string]string)
		if portraitTasks, ok := planPortraitTasks[k]; ok {
			if len(portraitTasks) > 0 {
				list, err := s.repo.GetTasksPortraitCustomers(c, corp.EnterpriseID, k, portraitTasks, yesterdayTs)
				if err != nil {
					log.Err(c, err).Msgf("查询截至到昨天重复客户去重统计数据保存到数据库错误:%+v", portraitTasks)
				}
				for _, ls := range list {
					allUnionsM[ls.UnionID] = ls.UnionID
					taskUnionKey := strconv.FormatInt(ls.TaskID, 10) + "-" + ls.UnionID
					taskUnionM[taskUnionKey] = ls.UnionID
				}
			}
		}

		// 2、查询群任务触达的客户集合
		if chatTasks, ok := planChatIDTasks[k]; ok {
			if len(chatTasks) > 0 {
				for _, ct := range chatTasks {
					chatUnions, _ := s.repo.GetTasksChatMembers(c, corp.EnterpriseID, taskIDM[ct])
					for _, union := range chatUnions {
						if union == "" {
							continue
						}
						taskUnionKey := strconv.FormatInt(ct, 10) + "-" + union
						taskUnionM[taskUnionKey] = union
						if _, ok := allUnionsM[union]; !ok {
							allUnionsM[union] = union
						}
					}
				}
			}
		}
		taskCustomerNum[k] = taskTmp{
			CustomerNum: int64(len(allUnionsM)),
			IsMain:      true,
		}
		// 匹配计算每个子任务的触达客户数
		for i := range taskUnionM {
			arr := strings.Split(i, "-")
			taskId, _ := strconv.ParseInt(arr[0], 10, 64)
			if taskId == 0 {
				continue
			}
			if _, ok := taskCustomerNum[taskId]; !ok {
				taskCustomerNum[taskId] = taskTmp{
					CustomerNum: 1,
				}
			} else {
				tmpCustomerNum := taskCustomerNum[taskId]
				tmpCustomerNum.CustomerNum += 1
				taskCustomerNum[taskId] = tmpCustomerNum
			}
		}
	}

	for k, y := range taskCustomerNum {
		var taskType int64
		taskType = 1 // 1主 2子
		if y.IsMain == false {
			taskType = 2
		}
		tmp := TaskEverydayStatistics{
			Date:                     date,
			DataTime:                 dataTs,
			TaskID:                   k,
			TotalReachCustomerNumber: y.CustomerNum,
			TaskType:                 taskType,
		}
		// 保存入库
		err = s.repo.InsertTaskEveryDayStatistic(c, corp.EnterpriseID, tmp, false)
		if err != nil {
			log.Err(c, err).Msgf("查询截至到昨天重复客户去重统计数据保存到数据库错误:%+v", tmp)
			continue
		}
	}

	// 2 查询有trace_id的计划任务统计其订单
	traceList, err := wsTraceService.SelectWsTracesBySourceType(c, corp.EnterpriseID, ws_trace.SourceTypeTask)
	if err != nil {
		return
	}
	planIDs := make(map[int64]int64)
	taskIDs := make(map[int64]int64)
	wg := &sync.WaitGroup{}
	ch := make(chan bool, 3) // 控制并发的数量
	for _, v := range traceList {
		wg.Add(1)
		ch <- true
		planIDs[v.MainTaskID] = v.MainTaskID
		taskIDs[v.ChildTaskID] = v.ChildTaskID
		go func(subTask ws_trace.WsTrace) {
			defer EndGroup(wg, ch)
			// 查询任务统计是否存在，如果存在更新
			tmp := TaskEverydayStatistics{
				Date:     date,
				DataTime: dataTs,
				TaskID:   subTask.SubTaskID,
				TaskType: 3,
			}
			// 获取聚合key
			wst := ws_trace.WsTrace{
				MemberID:    subTask.MemberID,
				MainTaskID:  subTask.MainTaskID,
				ChildTaskID: subTask.ChildTaskID,
				SubTaskID:   subTask.SubTaskID,
				GroupChatID: subTask.GroupChatID,
				SourceType:  ws_trace.SourceTypeTask,
			}
			keyM := assist.StructToMap(wst)
			ukVoucher, err := ws_trace.GetWsTracePrimary(corp.EnterpriseID, ws_trace.SourceTypeTask, keyM)
			if err != nil {
				log.Err(c, err).Msgf("查询每日子任务聚合key错误:%+v", v)
				return
			}
			// 查询订单统计数据
			orderStatistic, err := s.repo.GetSubTotalOrderStatistic(c, corp.EnterpriseID, ukVoucher, yesterdayTs)
			if err != nil {
				log.Err(c, err).Msgf("查询每日子任务订单统计数据错误:%+v", subTask.SubTaskID)
				return
			}
			if orderStatistic == nil {
				return
			}

			tmp.TotalOrderCustomerNumber = orderStatistic.TotalOrderCustomerNumber
			tmp.TotalOrderNumber = orderStatistic.TotalOrderNumber
			tmp.TotalOrderFee = orderStatistic.TotalOrderFee

			planStatistic, _ := s.GetTaskStatistics(c, corp.EnterpriseID, subTask.SubTaskID, yesterdayTs)
			if planStatistic != nil {
				tmp.TotalReachCustomerNumber = planStatistic.TotalReachCustomerNumber
			}

			// 保存入库
			err = s.repo.InsertTaskEveryDayStatistic(c, corp.EnterpriseID, tmp, true)
			if err != nil {
				log.Err(c, err).Msgf("查询每日子任务订单统计数据保存到数据库错误:%+v", tmp)
				return
			}
		}(v)
	}
	wg.Wait()
	// 处理plan主任务
	for _, taskID := range planIDs {
		// 查询任务统计是否存在，如果存在更新
		tmp := TaskEverydayStatistics{
			Date:     date,
			DataTime: dataTs,
			TaskID:   taskID,
			TaskType: 1,
		}
		// 获取查询任务的订单统计数据统计
		orderStatistic, err := s.repo.GetTotalOrderStatistic(c, corp.EnterpriseID, taskID, true, yesterdayTs)

		if err != nil {
			log.Err(c, err).Msgf("查询截至到昨天订单统计数据保存到数据库错误:%+v", taskID)
			continue
		}
		if orderStatistic == nil {
			continue
		}
		tmp.TotalOrderCustomerNumber = orderStatistic.TotalOrderCustomerNumber
		tmp.TotalOrderNumber = orderStatistic.TotalOrderNumber
		tmp.TotalOrderFee = orderStatistic.TotalOrderFee

		planStatistic, err := s.GetTaskStatistics(c, corp.EnterpriseID, taskID, yesterdayTs)
		if err != nil && err != sql.ErrNoRows {
			log.Err(c, err).Msgf("查询截至到昨天订单统计数据保存到数据库错误:%+v", tmp)
		}

		if planStatistic != nil {
			tmp.TotalReachCustomerNumber = planStatistic.TotalReachCustomerNumber
		}

		// 保存入库
		err = s.repo.InsertTaskEveryDayStatistic(c, corp.EnterpriseID, tmp, true)
		if err != nil {
			log.Err(c, err).Msgf("查询截至到昨天订单统计数据保存到数据库错误:%+v", tmp)
			continue
		}
	}

	for _, taskID := range taskIDs {
		// 查询任务统计是否存在，如果存在更新
		tmp := TaskEverydayStatistics{
			Date:     date,
			DataTime: dataTs,
			TaskID:   taskID,
			TaskType: 2,
		}
		// 获取查询任务的订单统计数据统计
		orderStatistic, err := s.repo.GetTotalOrderStatistic(c, corp.EnterpriseID, taskID, false, yesterdayTs)

		if err != nil {
			log.Err(c, err).Msgf("查询每日订单统计数据错误:%+v", taskID)
			continue
		}
		if orderStatistic == nil {
			continue
		}
		tmp.TotalOrderCustomerNumber = orderStatistic.TotalOrderCustomerNumber
		tmp.TotalOrderNumber = orderStatistic.TotalOrderNumber
		tmp.TotalOrderFee = orderStatistic.TotalOrderFee

		planStatistic, _ := s.GetTaskStatistics(c, corp.EnterpriseID, taskID, yesterdayTs)
		if planStatistic != nil {
			tmp.TotalReachCustomerNumber = planStatistic.TotalReachCustomerNumber
		}

		// 保存入库
		err = s.repo.InsertTaskEveryDayStatistic(c, corp.EnterpriseID, tmp, true)
		if err != nil {
			log.Err(c, err).Msgf("查询每日订单统计数据保存到数据库错误:%+v", tmp)
			continue
		}
	}
}

// GetPlanList 获取计划列表
func GetPlanList(ctx ctx.Context, enterpriseID int64, param PlanListRequest) (planListResponse []micro_scrm_enterprise.PlanModel, total int64, err error) {
	if param.StartTime == param.EndTime && param.StartTime != 0 {
		param.EndTime = param.StartTime + (60*60*24 - 1)
	}

	var operationUserId []int64
	// 查询当前传过来的创建人名字
	if param.CreateName != "" {
		enterpriseAdmin, err := orm.BossService.Enterprise.GetEnterpriseAdminList(boss.
			EnterpriseAdminModel{Username: param.CreateName})
		if err != nil {
			return planListResponse, total, err
		}

		for _, value := range enterpriseAdmin {
			operationUserId = append(operationUserId, value.Id)
		}
	}

	infoList, total, infoErr := orm.Service.PlanTask.GetPlanList(ctx, enterpriseID, param.SearchType, operationUserId,
		param.CreateName,
		param.PageParam.Page, param.PageParam.Size,
		micro_scrm_enterprise.PlanModel{
			BranchModel: model.BranchModel{
				BaseModel: model.BaseModel{
					ID: param.Id,
				},
			},
			Name:           param.Name,
			Status:         param.Status,
			StartTime:      param.StartTime,
			EndTime:        param.EndTime,
			PlanTemplateID: param.PlanTemplateId,
		},
		param.TaskLabel,
	)
	if infoErr != nil {
		return planListResponse, total, infoErr
	}

	return infoList, total, nil
}

// GetTaskList 获取任务列表
func GetTaskList(ctx ctx.Context, enterpriseID int64, param planReq.TaskListRequest) (taskListResponse []micro_scrm_enterprise.
	PlanTaskModel, total int64, err error) {

	infoList, total, infoErr := orm.Service.PlanTask.GetTaskList(enterpriseID, param.Page, param.Size, param)
	if infoErr != nil {
		return taskListResponse, total, infoErr
	}

	return infoList, total, nil
}

func EndGroup(wg *sync.WaitGroup, ch chan bool) {
	<-ch
	wg.Done()
}
