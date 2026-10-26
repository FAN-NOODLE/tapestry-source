package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	planReq "micro-scrm/internal/app/http/request/plan"
	excitationModels "micro-scrm/internal/app/models/excitation"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	memberSvc "micro-scrm/internal/app/service/member"
	planService "micro-scrm/internal/app/service/plan"
	workTagSvc "micro-scrm/internal/app/service/tags"
	"micro-scrm/internal/modules/ws_trace"
	"sort"
	"strconv"
	"time"

	"micro-scrm/internal/modules/members"

	"micro-scrm/internal/modules/plan_task/delivery/http_entity"

	"micro-scrm/internal/orm/model/micro_scrm_enterprise"

	"github.com/jinzhu/copier"
	"github.com/shakinm/xlsReader/xls"
	"github.com/shakinm/xlsReader/xls/structure"
	"github.com/tealeg/xlsx"

	plan_task2 "micro-scrm/internal/orm/service/plan_task"
	"micro-scrm/internal/pkg/helper"

	"github.com/jinzhu/now"

	response "micro-scrm/internal/pkg/resoponse"

	planTaskService "micro-scrm/internal/modules/plan_task"

	"micro-scrm/internal/modules/plan_task"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"

	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errs"
)

var wsTraceService *ws_trace.Service
var memberService *members.Service
var pService *plan_task.PTaskService

/**
*新增子任务
 */

const (
	LeaveGroup = 7
)

func planTasksAdd(c *gin.Context) {
	params := plan_task.PlanTasksRequest{}
	err := c.ShouldBindJSON(&params)
	if err != nil || len(params.PlanTaskConditions) == 0 {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if params.ConditionType == 0 {
		params.ConditionType = microScrmModel.TaskConditionTypeTag
	}
	if (params.Dimension == microScrmModel.TaskDimensionImportUnionId || params.Dimension == microScrmModel.TaskDimensionImportSaUnionId) && params.ConsumerFileId == 0 {
		ginplus.ResError(c, errors.New("导入的客户模板不存在"))
		return
	}

	if err := planTaskService.PlanTasksAdd(enterpriseID.(int64), ctx.Wrap(c), params); err != nil {
		ginplus.ResError(c, errs.OperationsBackError(err.Error()))
		return
	}

	ginplus.ResSuccess(c, nil)
}

/**
*修改子任务
 */
func planTasksUpdate(c *gin.Context) {
	params := plan_task.PlanTasksRequest{}
	err := c.ShouldBind(&params)
	if err != nil || params.TaskId == 0 {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if params.ConditionType == 0 {
		params.ConditionType = 1
	}
	/**
	*获取计划内容
	 */
	if err := planTaskService.PlanTasksUpdate(enterpriseID.(int64), ctx.Wrap(c), params); err != nil {
		ginplus.ResError(c, errs.OperationsBackError(err.Error()))
		return
	}

	ginplus.ResSuccess(c, nil)

}

/**
*子任务详情
 */
func planTasksInfo(c *gin.Context) {

	params := plan_task.PlanTasksInfoRequst{}
	err := c.BindJSON(&params)
	if err != nil || params.TaskId == 0 {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	ginplus.ResSuccess(c, planTaskService.PlanTasksInfo(ctx.Wrap(c), enterpriseID.(int64), params))
}

/*
*删除子任务
 */
func planTasksDelet(c *gin.Context) {
	params := plan_task.PlanTasksInfoRequst{}
	err := c.BindJSON(&params)
	if err != nil || params.TaskId == 0 {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if err := planTaskService.PlanTasksDelet(enterpriseID.(int64), params, ctx.Wrap(c)); err != nil {
		ginplus.ResError(c, errs.OperationsBackError(err.Error()))
		return
	}

	ginplus.ResSuccess(c, nil)
}

func GetPlanList(c *gin.Context) {
	// 获取body参数
	var body plan_task.PlanListRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		ginplus.ResError(c, err)
		return
	}

	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}

	_ctx := ctx.Wrap(c)
	data, total, err := plan_task.GetPlanList(_ctx, enterpriseID.(int64), body)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	var operationUserId []int64
	if data != nil {
		for _, value := range data {
			operationUserId = append(operationUserId, value.OperationUserId)
		}
	}

	operationMap := make(map[int64]string)
	if len(operationUserId) > 0 {
		enterpriseAdminList, err := orm.BossService.Enterprise.GetEnterpriseAdminListByOperationUserId(operationUserId)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}

		for _, value := range enterpriseAdminList {
			operationMap[value.Id] = value.Username
		}
	}

	var result = make([]*plan_task.PlanListDetail, 0)
	for _, value := range data {
		var detail plan_task.PlanListDetail

		detail.Id = value.ID
		detail.Name = value.Name
		detail.OperationUserId = value.OperationUserId
		detail.Remark = value.Remark
		detail.Status = value.Status
		detail.StartTime = value.StartTime
		detail.EndTime = value.EndTime
		detail.TaskNum = value.TaskNum
		detail.CreatedAt = value.CreatedAt
		detail.UpdatedAt = value.UpdatedAt
		detail.DeletedAt = value.DeletedAt
		detail.TemplateId = value.PlanTemplateID
		detail.Category = value.Category

		if value.IsOpen == 1 {
			detail.IsOpen = true
		} else {
			detail.IsOpen = false
		}

		if val, ok := operationMap[value.OperationUserId]; ok {
			detail.CreateName = val
		}

		result = append(result, &detail)
	}
	response.SuccessList(c, total, body.Page, body.Size, result)
	return
}

func GetTaskList(c *gin.Context) {
	// 获取body参数
	var body planReq.TaskListRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		ginplus.ResError(c, err)
		return
	}

	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}

	data, total, err := plan_task.GetTaskList(ctx.Wrap(c), enterpriseID.(int64), body)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	var result = make([]*plan_task.TaskListDetail, 0)
	for _, value := range data {
		var detail plan_task.TaskListDetail

		detail.Id = value.ID
		detail.Name = value.Name
		detail.PlanId = value.PlanId
		detail.Remark = value.Remark
		detail.Status = value.Status
		detail.ConditionType = value.ConditionType
		detail.Dimension = value.Dimension
		detail.CreatedAt = value.CreatedAt
		detail.UpdatedAt = value.UpdatedAt
		detail.DeletedAt = value.DeletedAt
		detail.TaskCategory = value.TaskCategory
		detail.GroupSendType = value.GroupSendType
		detail.FinishRate = value.FinishRate
		detail.TaskStartTime = value.TaskStartTime
		detail.TaskEndTime = value.TaskEndTime
		detail.PlanTemplateId = value.PlanTemplateId
		if value.IsOpen == 1 {
			detail.IsOpen = true
		} else {
			detail.IsOpen = false
		}
		detail.OriginType = value.OriginType
		result = append(result, &detail)
	}
	response.SuccessList(c, total, body.Page, body.Size, result)

	return
}

// getPlanListByFront 获取小程序任务计划列表
func getPlanListByFront(c *gin.Context) {
	var request plan_task.PlanListRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	// 指定日期
	dayTime := request.DateTime
	// 区间(月初月末或者某天的开始和结束)
	startTime := now.BeginningOfMonth().Unix()
	endTime := now.EndOfMonth().Unix()
	// 月初月末
	startTimeOfMonth := startTime
	endTimeOfMonth := endTime
	// 是否按天的维度统计
	isDay := false
	if dayTime != 0 {
		startTime = now.With(time.Unix(dayTime, 0)).BeginningOfMonth().Unix()
		if startTime == dayTime {
			// 指定某月的初始和末尾时间
			endTime = now.With(time.Unix(dayTime, 0)).EndOfMonth().Unix()
			startTimeOfMonth = startTime
			endTimeOfMonth = endTime
		} else {
			// 指定某天的初始和末尾时间
			startTime, endTime = helper.StartAndEndByDayTime(dayTime)
			isDay = true
		}
	}
	// 指定BA
	userId := getUserId(request.EnterpriseId, request.Mobile)
	if userId == "" {
		ginplus.ResError(c, errors.New("BA信息有误"))
		return
	}
	// 获取BA的所有计划id
	ids, _ := orm.Service.PlanTask.GetAllPlanIDS(request.EnterpriseId, startTime, endTime)
	if len(ids) == 0 {
		ginplus.ResError(c, errors.New("暂无任务计划"))
		return
	}
	var list []plan_task2.ResponseFrontPlanList
	var err error
	list, err = orm.Service.PlanTask.GetPlanListByFront(request.EnterpriseId, startTime, endTime, request.Page, request.Size, userId, request.Name, request.Status)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	allDays := make([]int, 0)
	allPlanId := make(map[int64]bool)
	newList := make([]plan_task.ResponseFrontPlanList, 0)
	for k, v := range list {
		allDays = append(allDays, helper.GetDayByUnix(v.CreatedAt))
		if allPlanId[v.Id] {
			continue
		}
		allPlanId[v.Id] = true
		if isDay {
			// 查指定日的所有子任务是否全部完成
			count := orm.Service.PlanTask.IsFinishByPlanTask(request.EnterpriseId, v.Id, startTime, endTime, userId)
			if count > 0 {
				list[k].StatusTex = "已完成"
			} else {
				list[k].StatusTex = "待完成"
				if time.Now().Unix() > v.EndTime {
					list[k].StatusTex = "已结束" // 未完成
				}
			}
		} else {
			switch list[k].Status {
			case 0:
				list[k].StatusTex = "未开始"
			case 1:
				list[k].StatusTex = "进行中"
			case 2:
				list[k].StatusTex = "已结束"
			}
		}
		temp := plan_task.ResponseFrontPlanList{
			Id:        list[k].Id,
			Name:      list[k].Name,
			Remark:    list[k].Remark,
			Status:    list[k].Status,
			StatusTex: list[k].StatusTex,
			IsOpen:    list[k].IsOpen,
			StartTime: list[k].StartTime,
			EndTime:   list[k].EndTime,
			CreatedAt: list[k].CreatedAt,
		}
		newList = append(newList, temp)
	}
	allDays = helper.ArrayUniq[int](allDays)
	resData := make(map[string]interface{})
	resData["list"] = newList
	resData["all_days"] = allDays
	finishDays := make([]int, 0)
	fDays := orm.Service.PlanTask.GetFinishDays(request.EnterpriseId, startTimeOfMonth, endTimeOfMonth, userId)
	for _, v1 := range fDays {
		if v1 >= startTimeOfMonth && v1 <= endTimeOfMonth {
			finishDays = append(finishDays, helper.GetDayByUnix(v1))
		}
	}
	resData["finish_days"] = helper.ArrayUniq[int](finishDays)
	ginplus.ResSuccess(c, resData)
	return
}

// getPlanTaskListByFront 获取小程序任务计划子任务列表
func getPlanTaskListByFront(c *gin.Context) {
	var request plan_task.PlanListRequest
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	// if request.Id == 0 {
	//	ginplus.ResError(c, errors.New("任务计划id不能为空"))
	//	return
	// }
	if request.EnterpriseId == 0 {
		ginplus.ResError(c, errors.New("企业id不能为空"))
		return
	}
	resData := make(map[string]interface{})
	resData["status"] = "已完成"
	resData["latest_time"] = time.Now().Unix()
	// 指定天的初始值
	startTime, endTime := helper.StartAndEndByDayTime(int64(request.StartTime))
	planInfo, _ := orm.Service.PlanTask.GetPlanInfoByFront(request.EnterpriseId, request.Id)
	userId := getUserId(request.EnterpriseId, request.Mobile)
	// userId := "YangMinYue"
	if userId == "" {
		ginplus.ResError(c, errors.New("BA信息有误"))
		return
	}
	list, err := orm.Service.PlanTask.GetPlanTaskListByFront(request.EnterpriseId, request.Id, startTime, endTime, request.Page, request.Size, userId, request.Status)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if len(list) == 0 {
		model, _ := orm.Service.PlanTask.GetLatestTaskByUserId(request.EnterpriseId, request.Id, userId)
		resData["latest_time"] = model.CreatedAt
		resData["status"] = ""
	}
	// 子任务是否全部完成
	isAllOk := true
	var newConditionIDList = make(map[int64]struct{}, 0)
	var newmap2 = make(map[int64]map[int64]struct{}, 0)
	var taskAllotLogIds []int64
	var newStartTime, newEndTime int64
	var groupChatAllotIds []int64
	for _, v := range list {
		newStartTime = v.TaskStartTime
		newEndTime = v.TaskEndTime
		newConditionIDList[v.ConditionID] = struct{}{}
		if val, ok := newmap2[v.ConditionID]; ok {
			val[v.GroupChatID] = struct{}{}
		} else {
			var mapInfo = make(map[int64]struct{}, 0)
			mapInfo[v.GroupChatID] = struct{}{}
			newmap2[v.ConditionID] = mapInfo
		}
		if v.Status != 1 {
			isAllOk = false
		}
		if v.TaskCategory == microScrmModel.PLANTASKCONTACT && v.GroupSendType == microScrmModel.PLANTASKQYSEND {
			taskAllotLogIds = append(taskAllotLogIds, v.Id)
		}
		if v.Dimension == 1 { // 社群维度
			groupChatAllotIds = append(groupChatAllotIds, v.Id)
		}
	}
	var newConList = make([]int64, 0)
	for k := range newConditionIDList {
		newConList = append(newConList, k)
	}
	switch planInfo.Status {
	case 0:
		resData["status"] = "未开始"
	case 1:
		if !isAllOk {
			resData["status"] = "进行中"
		}
	case 2:
		resData["status"] = "已结束"
		// if !isAllOk {
		//	resData["status"] = "未完成"
		// }
	}
	resData["status_code"] = planInfo.Status
	// groupChatMap := make(map[int64][]int64, 0)
	var chatIdList = make([]int64, 0)
	if len(groupChatAllotIds) > 0 {
		many, err := microScrmDao.NewPlanTaskConditionDataDao().GetMany(ctx.Wrap(c), request.EnterpriseId, request.Id, map[string]interface{}{
			"id":         groupChatAllotIds,
			"deleted_at": 0,
		})
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		for _, model := range many {
			chatIdList = append(chatIdList, model.GroupChatId)
		}
	}

	var newList = make([]*micro_scrm_enterprise.NewResponseFrontPlanTaskList, 0)
	var conMap = make(map[int64]micro_scrm_enterprise.ConditionInfo, 0)
	info, _ := orm.Service.PlanTask.GetConditionInfo(request.EnterpriseId, newConList)
	// statusInfo := orm.Service.PlanTask.GetStatus(request.EnterpriseId, newStartTime, newEndTime, newConList)
	info2 := orm.Service.PlanTask.GetConditionDataInfo2(request.EnterpriseId, newStartTime, newEndTime, chatIdList)
	// var statusMap = make(map[int64]int, 0)
	var leaveNumberMap = make(map[int64]map[string]struct{}, 0)
	// for _, v2 := range statusInfo {
	// 	if val, ok := statusMap[v2.ConditionID]; !ok {
	// 		statusMap[v2.ConditionID] = v2.IsFinished
	// 	} else {
	// 		if v2.IsFinished != 1 {
	// 			statusMap[v2.ConditionID] = 0
	// 		} else if v2.IsFinished == 1 && val != 0 {
	// 			statusMap[v2.ConditionID] = 1
	// 		} else {
	// 			continue
	// 		}
	// 	}
	// }
	for _, detailInfo := range info2 {
		var newDetail = make(map[string]struct{}, 0)
		if _, ok := leaveNumberMap[detailInfo.GroupChatID]; !ok {
			newDetail[detailInfo.WorkExternalUserId] = struct{}{}
			leaveNumberMap[detailInfo.GroupChatID] = newDetail
		} else {
			leaveNumberMap[detailInfo.GroupChatID][detailInfo.WorkExternalUserId] = struct{}{}
		}
	}

	for _, conditionInfo := range info {
		conMap[conditionInfo.ID] = conditionInfo
	}
	chqun := map[int64]*micro_scrm_enterprise.NewResponseFrontPlanTaskList{}
	for _, taskList := range list {
		var newInfo micro_scrm_enterprise.NewResponseFrontPlanTaskList
		var newConditionDetail micro_scrm_enterprise.ConditionDetail
		var newContent micro_scrm_enterprise.Content
		_ = copier.Copy(&newInfo, &taskList)

		// m2 := newmap2[taskList.ConditionID]
		// var sList = make([]int64, 0)
		// for k := range m2 {
		// 	sList = append(sList, k)
		// }
		// if item, ok := groupChatMap[taskList.Id]; ok {
		// 	sList = item
		// }
		if val, ok := conMap[taskList.ConditionID]; ok {
			if err := json.Unmarshal([]byte(val.Conditions), &newConditionDetail); err != nil {
				ginplus.ResError(c, err)
				return
			}
			if err := json.Unmarshal([]byte(val.Content), &newContent); err != nil {
				ginplus.ResError(c, err)
				return
			}
			newInfo.ConditionInfo = newConditionDetail
			newInfo.ContentInfo = newContent
		}
		newInfo.GroupNumber = int64(taskList.ChatNum)
		if val, ok := leaveNumberMap[taskList.GroupChatID]; ok {
			newInfo.LeaveNumber = int64(len(val))
		}
		if newInfo.ConditionType == microScrmModel.TaskConditionTypeGroupNum || newInfo.ConditionType == microScrmModel.TaskConditionTypeAddGroupNum {
			chqun[taskList.ConditionID] = &newInfo
		}
		statusCode := 0 // 0为进行中，1为已完成，2为已结束
		newInfo.StatusCode = statusCode
		if taskList.Dimension == 1 {
			newInfo.StatusCode = newInfo.Status
			if planInfo.Status == 2 && newInfo.Status != 1 {
				newInfo.StatusCode = 2
			}
		} else {
			newInfo.StatusCode = taskList.Status
			if planInfo.Status == 2 && taskList.Status != 1 {
				newInfo.StatusCode = 2
			}
		}
		newInfo.Status = taskList.Status
		newInfo.FinishRate = taskList.FinishRate
		newList = append(newList, &newInfo)
	}
	resultList := make([]*micro_scrm_enterprise.NewResponseFrontPlanTaskList, 0)
	nofinishNums, err := microScrmDao.NewPlanTaskConditionDataDao().GetNoFinishedCountByIds(ctx.Wrap(c), request.EnterpriseId, request.Id, taskAllotLogIds)
	if err != nil {
		return
	}
	nofinishedMap := make(map[int64]int64, 0)
	for _, v := range nofinishNums {
		nofinishedMap[v.Id] = v.Num
	}
	for _, taskList := range newList {
		item := taskList
		if item.ConditionType == 2 || item.ConditionType == 3 {
			if value, ok := chqun[item.ConditionID]; ok {
				resultList = append(resultList, value)
				delete(chqun, item.ConditionID)

			}
			item.CloseTips = 1
		} else {
			item.NoFinishedNum = nofinishedMap[item.Id]
			if item.CloseTips == 0 && item.NoFinishedNum == 0 {
				item.CloseTips = 1
			}
			resultList = append(resultList, item)
		}

	}
	sort.Slice(resultList, func(i, j int) bool {
		if resultList[i].Status == resultList[j].Status {
			return resultList[i].Id < resultList[j].Id
		} else {
			return resultList[i].Status < resultList[j].Status
		}

	})

	resData["list"] = resultList
	resData["all_days"] = orm.Service.PlanTask.GetPlanTaskAllDays(request.EnterpriseId, planInfo.Id, planInfo.StartTime, planInfo.EndTime, userId)
	resData["finish_days"] = orm.Service.PlanTask.GetPlanTimeFinishDays(request.EnterpriseId, planInfo.Id, planInfo.StartTime, planInfo.EndTime, userId)
	ginplus.ResSuccess(c, resData)
	return
}

// @author :LiKun
// function : 计划增删查改
//start:2021-07-28
//update:2021-07-29

func AddPlan(c *gin.Context) {
	var receive plan_task.ReceivePlan
	if err := c.ShouldBindJSON(&receive); err != nil {
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if err := plan_task.PlanAdd(ctx.Wrap(c), enterpriseID.(int64), receive); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

func PlanEdit(c *gin.Context) {
	var receiveInfo plan_task.ReceivePlanModel
	if err := c.ShouldBindJSON(&receiveInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if err := plan_task.PlanEdit(ctx.Wrap(c), enterpriseID.(int64), receiveInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

func RemovePlan(c *gin.Context) {
	var planId plan_task.PlanId
	if err := c.ShouldBindJSON(&planId); err != nil {
		// log.Fatal("Data parsing error")
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	_ = plan_task.PlanRemove(ctx.Wrap(c), planId.ID, enterpriseID.(int64))
	ginplus.ResOK(c)
}

func PlanCopy(c *gin.Context) {
	var receiveInfo plan_task.ReceivePlanModel
	if err := c.ShouldBindJSON(&receiveInfo); err != nil {
		// log.Fatal("Data parsing error")
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if err := plan_task.PlanCopy(ctx.Wrap(c), enterpriseID.(int64), receiveInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

func planTasksConditionsCompletion(c *gin.Context) {
	params := plan_task.PlanTasksConditionsRequst{}
	if err := c.BindJSON(&params); err != nil || params.Id == 0 {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	result, total := planTaskService.PlanTasksConditionsCompletion(enterpriseID.(int64), params)
	response.SuccessList(c, total, params.Page, params.Size, result)
	return
}

func planTasksConditionsUserCompletion(c *gin.Context) {
	params := plan_task.PlanTasksConditionsBaRequst{}
	if err := c.BindJSON(&params); err != nil || params.Id == 0 || params.User_id == "" {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	result, total := planTaskService.PlanTasksConditionsUserCompletion(ctx.Wrap(c), enterpriseID.(int64), params)
	response.SuccessList(c, total, params.Page, params.Size, result)
	return
}

func UpdateStatus(c *gin.Context) {
	var receive plan_task.PlanStatus
	if err := c.ShouldBindJSON(&receive); err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if err := plan_task.PlanUpdateStatus(ctx.Wrap(c), receive, enterpriseID.(int64)); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// getPlanTaskInfoByFront 获取小程序子任务详情
func getPlanTaskInfoByFront(c *gin.Context) {
	var request plan_task.PlanTasksConRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	// 退群信息
	listsa := orm.Service.PlanTask.GetChatId(request.EnterpriseId, request.Id)
	// 条件
	var conditionInfo micro_scrm_enterprise.ConditionDetail
	var contentInfo micro_scrm_enterprise.Content
	if err := json.Unmarshal([]byte(listsa.Conditions), &conditionInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if err := json.Unmarshal([]byte(listsa.Content), &contentInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	var customerChurnRate string
	var result plan_task.NewResponseFrontPlanTaskConList

	idInfo := orm.Service.PlanTask.GetConditionIdInfo(request.EnterpriseId, request.Id)
	if listsa.Dimension == microScrmModel.TaskDimensionCommunity { // 如果是社群的查社群逻辑
		// 查询退群的逻辑
		one, err := microScrmDao.NewPlanTaskConditionDataDao().GetOne(ctx.Wrap(c), request.EnterpriseId, listsa.PlanID, map[string]interface{}{
			"task_allot_log_id": request.Id,
		})
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		dataInfo := orm.Service.PlanTask.GetConditionDataInfo(request.EnterpriseId, listsa.StartTime, listsa.EndTime, one.GroupChatId)
		numCount := orm.Service.PlanTask.GetNumCount(request.EnterpriseId, one.GroupChatId)
		count := int64(len(dataInfo))
		if count == 0 {
			customerChurnRate = "0%"
		} else {
			zCount := numCount + count
			newNumber := float64(count) / float64(zCount)
			newNumber2 := math.Floor(newNumber*100 + 0.5)
			customerChurnRate = fmt.Sprintf("%.0f", newNumber2) + "%"
		}
		result.LeaveNumber = count

		var nameMap = make(map[int64]string, 0)
		var chatIdMap = make(map[int64]string, 0)
		for _, names := range orm.Service.PlanTask.GetChatName(request.EnterpriseId) {
			nameMap[names.ID] = names.Name
			chatIdMap[names.ID] = names.ChatID
		}
		// 查询群任务相关数据
		getGroupInfo, _ := microScrmDao.NewPlanTaskConditionDataDao().GetMany(ctx.Wrap(c), request.EnterpriseId, idInfo.PlanId, map[string]interface{}{
			"task_allot_log_id": request.Id,
		})
		var groupInfo []*plan_task.GroupInfo
		// flag := false
		var groupNums int
		for _, info := range getGroupInfo {
			// if _, ok := newGroupMap[info.GroupChatID]; !ok {
			var groupInfoa plan_task.GroupInfo
			groupInfoa.ID = info.ID
			groupInfoa.IsFinish = int(info.IsFinished)
			groupInfoa.GroupChatID = info.GroupChatId
			if val, ok := nameMap[info.GroupChatId]; ok {
				groupInfoa.GroupName = val
			}
			if val, ok := chatIdMap[info.GroupChatId]; ok {
				groupInfoa.ChatID = val
			}
			if info.IsFinished == 1 {
				groupNums++
			}
			groupInfo = append(groupInfo, &groupInfoa)
		}
		result.GroupInfo = groupInfo
	}

	result.Id = request.Id
	result.GroupChatID = listsa.GroupChatId
	result.ConditionID = listsa.ConditionID
	result.Remark = listsa.Remark
	result.CustomerChurnRate = customerChurnRate

	result.StatusCode = idInfo.IsFinished                                                             // 非群任务，任务完成，则子任务完成
	if idInfo.Status == microScrmModel.PlanEnd && idInfo.IsFinished != microScrmModel.PLAN_FINISHED { // 任务已结束，子任务未完成，则为过期
		result.StatusCode = microScrmModel.PlanEnd
	}
	if idInfo.IsFinished == microScrmModel.PLAN_NOFINISH {
		result.IsFinished = microScrmModel.PLAN_NOFINISH
	}
	if listsa.ConditionType == LeaveGroup {
		result.IsFinished = listsa.IsFinished
	}
	qyTagIdMap := make(map[string]string)
	if listsa.Dimension == microScrmModel.TaskDimensionTag {
		tagLists, _ := workTagSvc.NewWorkTags().GetChoiceTagList(ctx.Wrap(c), request.EnterpriseId)
		for _, model := range tagLists {
			qyTagIdMap[model.TagID] = model.Name
		}
	}
	cc := ctx.Wrap(c)
	taskInfo, _ := microScrmDao.NewPlanTaskDao().GetOne(cc, request.EnterpriseId, map[string]interface{}{
		"id": idInfo.TaskId,
	})
	result.TaskCategory = taskInfo.TaskCategory
	result.GroupSendType = taskInfo.GroupSendType
	result.IsAddTag = taskInfo.IsAddTag
	result.PlanId = taskInfo.PlanID
	result.PlanName = idInfo.PlanName
	tagNames := make([]string, 0)
	if len(conditionInfo.Tag) > 0 {
		for _, value := range conditionInfo.Tag {
			if tagName, ok := qyTagIdMap[value.QwTagID]; ok {
				tagNames = append(tagNames, tagName)
			}
		}
	}
	if taskInfo.QyTagId != "" {
		if tagName, ok := qyTagIdMap[taskInfo.QyTagId]; ok {
			tagNames = append(tagNames, tagName)
		}
	}
	result.Tags = tagNames
	result.ExpireStatus = taskInfo.Status
	if taskInfo.TaskCategory == microScrmModel.PLANTASKCONTACT && taskInfo.GroupSendType == microScrmModel.PLANTASKQYSEND {
		noFinishedNum, _ := microScrmDao.NewPlanTaskConditionDataDao().GetDistinctUserCount(ctx.Wrap(c), request.EnterpriseId, idInfo.PlanId, map[string]interface{}{
			"task_allot_log_id": request.Id,
			"is_finished":       0,
		})
		result.NoFinishedNum = noFinishedNum
	}

	if contentInfo.Applets.AppId != "" {
		// sw小程序路径添加参数
		resConfig, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(cc, request.EnterpriseId, "store_name")
		if resConfig.Name == "sw" {
			// 修改小程序
			pathMap, err := memberSvc.NewMembersService().SetMiniProgramLinkParams(cc, request.EnterpriseId, request.MemberId, map[int64]string{0: contentInfo.Applets.Path})
			if err != nil {
				return
			}
			contentInfo.Applets.Path = pathMap[0]
		}
	}
	taskAllotInfo, err := microScrmDao.NewPlanTaskAllotLogDao().GetOne(cc, request.EnterpriseId, map[string]interface{}{
		"id": request.Id,
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	result.Dimension = listsa.Dimension
	result.ConditionType = listsa.ConditionType
	result.ContentType = listsa.ContentType
	result.ConditionInfo = conditionInfo
	result.ContentInfo = contentInfo
	result.IsForward = idInfo.IsForward
	result.TaskId = idInfo.TaskId
	result.FinishRate = listsa.FinishRate
	result.TaskStartTime = listsa.TaskStartTime
	result.TaskEndTime = listsa.TaskEndTime
	result.TaskExpireType = listsa.TaskExpireType
	result.CustomerNum = taskAllotInfo.CustomerNum
	ginplus.ResSuccess(c, &result)
	return
}

func UpdateGroupStatus(c *gin.Context) {
	var receive plan_task.UpdateReceive
	if err := c.ShouldBindJSON(&receive); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if receive.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业id不能为空"))
		return
	}
	if err := orm.Service.PlanTask.UpdateGroupStatus(receive.EnterpriseID, receive.Status, receive.ID); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

func setPlanTaskStatusByFront(c *gin.Context) {
	var request plan_task.TaskConditionReq
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	excitationSvc.NewPointsProduceService().TaskFinished(ctx.Wrap(c), excitationReq.TaskFinishedReq{
		TaskAllotLogId: request.ID,
		EnterpriseId:   request.EnterpriseId,
		Mobile:         request.Mobile,
	})
	userId := getUserId(request.EnterpriseId, request.Mobile)
	err := orm.Service.PlanTask.SetStatusByFront(request.EnterpriseId, request.Dimension, request.ID, request.ConditionId, request.GroupChatID, userId)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	// 积分队列埋点
	res, err := microBossDao.NewEnterpriseConfig().GetInfo(ctx.Wrap(c), map[string]interface{}{"enterprise_id": request.EnterpriseId, "key": "store_name"})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if request.Dimension != 1 {
		if res.Name != "coach" {
			err := planService.NewPlanTaskConditionData().FinishALL(ctx.Wrap(c), request.EnterpriseId, request.ID, userId)
			if err != nil {
				ginplus.ResError(c, err)
				return
			}
		}
		excitationSvc.NewPointsProduceService().ProduceQueue(ctx.Wrap(c), excitationReq.PointsProduceSave{
			Mobile:       request.Mobile,
			Type:         excitationModels.RuleCategoryId7,
			Value:        1,
			UniqueHash:   strconv.FormatInt(request.ID, 10) + fmt.Sprintf("-%s-%d", request.Mobile, excitationModels.RuleCategoryId7),
			EnterpriseId: request.EnterpriseId,
		})

	}

	ginplus.ResSuccess(c, []interface{}{})
	return
}

func getUserId(eid int64, mobile string) string {
	model := orm.Service.Member.GetUserIdByMobile(eid, mobile)
	return model.Userid
}

// 获取标签的csv文件保存并上传
func getTagCSVFile(c *gin.Context) {
	var request plan_task.GetFileInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if request.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业ID为空!"))
		return
	}

	if err := plan_task.CSVUploadLogin(request, request.EnterpriseID); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
}

func DownloadFileInfo(c *gin.Context) {
	var request plan_task.DownLoadFileInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if request.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业ID为空!"))
		return
	}
	var idInfo plan_task.IDInfo
	ID, err := plan_task.SaveUploadInfo(request, request.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var newTagInfo plan_task.TagInfo
	newTagInfo.TagName = request.TagName
	newTagInfo.GroupName = request.TagGroupName
	if request.BAType == "1" {
		newTagInfo.BaType = 1
	} else if request.BAType == "2" {
		newTagInfo.BaType = 2
	} else {
		newTagInfo.BaType = 0
	}
	// 创建标签组
	if err := plan_task.CreateTag(ctx.Wrap(c), request.EnterpriseID, newTagInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}

	idInfo.ID = ID
	ginplus.ResSuccess(c, &idInfo)
}

// 文件存在校验
func GetAllFileName(c *gin.Context) {
	var fileDetection plan_task.FileDetection
	if err := c.ShouldBindJSON(&fileDetection); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if fileDetection.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业ID为空!"))
		return
	}
	detection, err := plan_task.GetFileNameDetection(fileDetection.EnterpriseID, fileDetection.FileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var result plan_task.CountInfo
	result.Count = detection
	ginplus.ResSuccess(c, &result)
}

func UpdateIsExtract(c *gin.Context) {
	var request plan_task.UpdateIdInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if request.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业ID为空!"))
		return
	}
	if err := plan_task.GetUpdateInfo(request.EnterpriseID, request.ID); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
}

func customerData(c *gin.Context) {
	var params plan_task.CustomerDataReq
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}

	res, err := plan_task.CustomerDataService(ctx.Wrap(c), enterpriseID.(int64), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func taskListExport(c *gin.Context) {
	params := http_entity.ExportTasksRequst{}
	if err := c.BindJSON(&params); err != nil || params.Id == 0 {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if params.CompletionStartTime == params.CompletionEndTime && params.CompletionStartTime != 0 {
		params.CompletionEndTime = params.CompletionStartTime + (60*60*24 - 1)
	}
	if params.PlanStartTime == params.PlanEndTime && params.PlanStartTime != 0 {
		params.PlanEndTime = params.PlanStartTime + (60*60*24 - 1)
	}

	result, err := plan_task.ExportCustomerTaskData(ctx.Wrap(c), enterpriseID.(int64), params, pService)
	if err != nil {
		response.Success(c, err)
		return
	}
	response.Success(c, result)
	return
}

func groupTaskData(c *gin.Context) {
	var params plan_task.CustomerDataReq
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}

	res, err := plan_task.GroupTaskDataService(enterpriseID.(int64), params.PlanID, params.TaskID, params.StartTime, params.EndTime, pService, ctx.Wrap(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func groupTaskDataExport(c *gin.Context) {
	params := http_entity.ExportTasksRequst{}
	if err := c.BindJSON(&params); err != nil || params.Id == 0 {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if params.CompletionStartTime == params.CompletionEndTime && params.CompletionStartTime != 0 {
		params.CompletionEndTime = params.CompletionStartTime + (60*60*24 - 1)
	}
	if params.PlanStartTime == params.PlanEndTime && params.PlanStartTime != 0 {
		params.PlanEndTime = params.PlanStartTime + (60*60*24 - 1)
	}

	result, err := plan_task.ExportGroupTaskData(ctx.Wrap(c), enterpriseID.(int64), params, pService, ctx.Wrap(c))
	if err != nil {
		response.Success(c, err)
		return
	}
	response.Success(c, result)
	return
}

func customerTaskListData(c *gin.Context) {
	params := plan_task.PlanTasksConditionsRequst{}
	if err := c.BindJSON(&params); err != nil || params.Id == 0 {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	result, total := planTaskService.PlanTasksConditionsCompletionV2(enterpriseID.(int64), params, pService, ctx.Wrap(c))
	response.SuccessList(c, total, params.Page, params.Size, result)
	return
}

func groupTaskListData(c *gin.Context) {
	params := http_entity.TaskListRequst{}
	if err := c.BindJSON(&params); err != nil || params.Id == 0 {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	enterpriseID, isExist := c.Get(common.ContextKeyEnum.EnterpriseID)
	if !isExist || enterpriseID.(int64) <= 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if params.CompletionStartTime == params.CompletionEndTime && params.CompletionStartTime != 0 {
		params.CompletionEndTime = params.CompletionStartTime + (60*60*24 - 1)
	}
	if params.PlanStartTime == params.PlanEndTime && params.PlanStartTime != 0 {
		params.PlanEndTime = params.PlanStartTime + (60*60*24 - 1)
	}
	result, total, err := plan_task.GroupTaskListData(ctx.Wrap(c), enterpriseID.(int64), params, pService, ctx.Wrap(c))
	if err != nil {
		response.Success(c, err)
		return
	}
	response.SuccessList(c, total, params.Page, params.Size, result)
	return
}

// uploadTaskUsers 上传执行任务计划的员工
func uploadTaskUsers(c *gin.Context) {
	fileType := c.PostForm("file_type")
	enterpriseIdStr := c.PostForm("enterprise_id")
	enterPriseID, _ := strconv.ParseInt(enterpriseIdStr, 10, 64)
	fh, err := c.FormFile("file")
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	var (
		ft        string
		totalData []*plan_task.UploadUserResponse
	)
	fileData, err := common.ReadExcelFile(fh, &ft)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if ft == ".xlsx" {
		totalData = readXlsxData(fileData.(*xlsx.File).Sheets)
	} else {
		totalData = readXlsData(*fileData.(*xls.Workbook))
	}

	data := readUploadUserList(ctx.Wrap(c), enterPriseID, totalData, fileType)
	ginplus.ResSuccess(c, data)
	return
}

func readXlsxData(sheets []*xlsx.Sheet) []*plan_task.UploadUserResponse {
	var data []*plan_task.UploadUserResponse
	for index, sheet := range sheets {
		if index > 0 {
			continue
		}
		for ri, row := range sheet.Rows {
			if ri < 1 {
				continue
			}
			if len(row.Cells) == 0 {
				continue
			}
			if row.Cells[0].Value == "" {
				continue
			}
			var upUser = &plan_task.UploadUserResponse{
				LineNum: ri + 1,
				Item:    row.Cells[0].String(),
			}
			data = append(data, upUser)
		}
	}
	return data
}

func readXlsData(wb xls.Workbook) []*plan_task.UploadUserResponse {
	var data []*plan_task.UploadUserResponse
	for index, sheet := range wb.GetSheets() {
		if index > 0 {
			continue
		}
		for ri, row := range sheet.GetRows() {
			if ri < 1 {
				continue
			}
			if len(row.GetCols()) == 0 {
				continue
			}
			if cell, err := row.GetCol(0); err != nil || cell.GetString() == "" {
				continue
			}
			getColStr := func(c structure.CellData, err error) string {
				if err != nil {
					return ""
				}
				return c.GetString()
			}

			var upUser = &plan_task.UploadUserResponse{
				LineNum: ri + 1,
				Item:    getColStr(row.GetCol(0)),
			}
			data = append(data, upUser)
		}
	}
	return data
}

func readUploadUserList(ctx ctx.Context, enterpriseId int64, totalData []*plan_task.UploadUserResponse, fileType string) []*plan_task.UploadUserResponse {
	var (
		Items     []string
		data      []*plan_task.UploadUserResponse
		upDataMap = make(map[string]string)
	)
	for i := 0; i <= len(totalData)-1; i++ {
		var upUser = totalData[i]
		data = append(data, upUser)
		if _, ok := upDataMap[upUser.Item]; ok {
			upUser.ErrType = 2
			upUser.Error = "数据重复"
		} else {
			Items = append(Items, upUser.Item)
		}
		upDataMap[upUser.Item] = upUser.Item
	}

	if fileType == "userid" {
		var memberMap = make(map[string]*members.Member)
		memberUsers, _ := memberService.GetMembersByUserIdsOrMobiles(ctx, enterpriseId, members.MemberByUserID, Items)
		for _, v := range memberUsers {
			memberMap[v.Userid] = v
		}
		for _, v := range data {
			if m, ok := memberMap[v.Item]; !ok {
				v.ErrType = 1
				v.Error = "User ID 错误"
			} else {
				v.UserID = m.Userid
				v.Name = m.Name
			}
		}
		return data
	}
	var memberMap = make(map[string]*members.Member)
	memberUsers, _ := memberService.GetMembersByMarsIds(ctx, enterpriseId, Items)
	for _, v := range memberUsers {
		memberMap[v.MarsId] = v
	}
	for _, v := range data {
		if m, ok := memberMap[v.Item]; !ok {
			v.ErrType = 1
			v.Error = "Mars Code 错误"
		} else {
			v.UserID = m.Userid
			v.Name = m.Name
		}
	}
	return data
}
