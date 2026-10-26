package planService

import (
	"encoding/json"
	"errors"
	"fmt"
	"micro-scrm/internal/app"
	customerDao "micro-scrm/internal/app/dao/customer"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	membersDao "micro-scrm/internal/app/dao/members"
	externalUserDao "micro-scrm/internal/app/dao/micro_scrm"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	tagsDao "micro-scrm/internal/app/dao/tags"
	workGroupDao "micro-scrm/internal/app/dao/work"
	customerReq "micro-scrm/internal/app/http/request/customer"
	customerSearchReq "micro-scrm/internal/app/http/request/customer_search"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	openReq "micro-scrm/internal/app/http/request/open"
	planReq "micro-scrm/internal/app/http/request/plan"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	customerSearchResp "micro-scrm/internal/app/http/response/customer_search"
	planResp "micro-scrm/internal/app/http/response/plan"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	workModels "micro-scrm/internal/app/models/work"
	customerSvc "micro-scrm/internal/app/service/customer"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	workTagSvc "micro-scrm/internal/app/service/tags"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/helper"
	pkgKafka "micro-scrm/internal/pkg/kafka"
	"micro-scrm/internal/pkg/sdk/macallback"
	"micro-scrm/internal/pkg/sdk/microwork"
	"micro-scrm/internal/pkg/sdk/work"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"gorm.io/gorm"
)

type PlanTask struct {
	planTaskDao *planDao.PlanTaskDao
}

func NewPlanTask() *PlanTask {
	return &PlanTask{
		planTaskDao: planDao.NewPlanTaskDao(),
	}
}

func (p PlanTask) UpdatePlanTaskPushStatus(c app.Ctx, enterpriseID, taskID int64) error {
	return p.planTaskDao.UpdatePlanTaskPushStatus(c, enterpriseID, taskID)
}

func (p PlanTask) UpdatePlanTaskPush(c app.Ctx, enterpriseID, taskID int64) error {
	return p.planTaskDao.UpdatePlanTaskPush(c, enterpriseID, taskID)
}

func (p PlanTask) UpdatePlanTaskPushFailed(c app.Ctx, enterpriseID, taskID int64) error {
	return p.planTaskDao.UpdatePlanTaskPushFailed(c, enterpriseID, taskID)
}

// GetOrderMonthLists 获取任务日历数据
func (p PlanTask) GetOrderMonthLists(c ctx.Context, request planReq.MonthListsReq) ([]planReq.MonthListsData, error) {
	monthMap := make(map[int64]planReq.MonthListsStatus)
	var monthSplice []planReq.MonthListsData
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: request.EnterpriseId},
		Mobile:         request.Mobile,
	})
	if err != nil || memberInfo.ID == 0 {
		return monthSplice, errorx.NewErrorMsg(1001, "查询导购信息异常")
	}
	lists, err := externalUserDao.NewReservesOrder().GetOrderMonthLists(c, request, memberInfo.ID)
	if err != nil {
		return monthSplice, err
	}
	for _, v := range lists {
		zeroTime := helper.GetZeroDateByParams(int64(v.ReserveTime))
		if item, ok := monthMap[zeroTime]; !ok {
			if v.Status == 1 || v.Status == 2 { // 只处理进行中，完成和成功的状态
				monthMap[zeroTime] = planReq.MonthListsStatus{Finished: true, Afoot: false}
			} else if v.Status == 0 {
				monthMap[zeroTime] = planReq.MonthListsStatus{Finished: false, Afoot: true}
			}
		} else {
			if v.Status != 1 && v.Status != 2 && item.Finished == true {
				monthMap[zeroTime] = planReq.MonthListsStatus{Finished: false, Afoot: true}
			}
		}
	}

	for k, value := range monthMap {
		info := planReq.MonthListsData{
			Time:             k,
			MonthListsStatus: planReq.MonthListsStatus{Finished: value.Finished, Afoot: value.Afoot},
		}
		monthSplice = append(monthSplice, info)
	}
	return monthSplice, err
}

// GetPlanOrderMonthLists 获取计划任务的日历列表数据
func (p PlanTask) GetPlanOrderMonthLists(c ctx.Context, request planReq.MonthListsReq) ([]planReq.MonthListsData, error) {
	monthMap := make(map[int64]planReq.MonthListsStatus)
	// planMap := make(map[int64]bool)
	var monthSplice []planReq.MonthListsData
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: request.EnterpriseId},
		Mobile:         request.Mobile,
	})
	if err != nil || memberInfo.ID == 0 {
		return monthSplice, errorx.NewErrorMsg(1001, "查询导购信息异常")
	}
	data, err := planDao.NewPlanTaskAllotLogDao().GetPlanOrderMonthLists(c, request, memberInfo.Userid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return monthSplice, err
	}
	for _, value := range data {
		zeroTime := helper.GetZeroDateByParams(value.CreatedAt)
		if item, ok := monthMap[zeroTime]; !ok {
			if value.IsFinished == 1 {
				monthMap[zeroTime] = planReq.MonthListsStatus{Finished: true}
			} else if value.IsFinished == 0 {
				monthMap[zeroTime] = planReq.MonthListsStatus{Afoot: true}
			}
		} else {
			if value.IsFinished == 0 && item.Finished == true {
				monthMap[zeroTime] = planReq.MonthListsStatus{Afoot: true}
			}
		}
	}
	for k, value := range monthMap {
		info := planReq.MonthListsData{
			Time:             k,
			MonthListsStatus: value,
		}
		monthSplice = append(monthSplice, info)
	}
	return monthSplice, err
}

// GetPlanLists 查询任务中心的主任务列表
func (p PlanTask) GetPlanLists(c ctx.Context, request planReq.ListsReq) ([]planReq.PlansListInfo, error) {
	var result []planReq.PlansListInfo
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: request.EnterpriseId},
		Mobile:         request.Mobile,
	})
	if err != nil || memberInfo.ID == 0 {
		return result, errorx.NewErrorMsg(1001, "查询导购信息异常")
	}
	mapList := make(map[int64]planReq.PlansListInfo)

	data, err := planDao.NewPlanTaskAllotLogDao().GetPlanLists(c, request, memberInfo.Userid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return result, err
	}
	unShowIds := make(map[int64]bool)
	nowTime := time.Now().Unix()
	var planIds []int64
	// 1、当天内所有任务均已完成的，显示绿点
	// 2、当天所有任务均待完成的，展示红点
	// 3、当天所有任务均已过期的，不展示点
	// 4、当天内同时存在已完成、已过期、待完成的，展示红点
	// 3、当天内仅存在已完成、已过期的，展示绿点
	// 4、当天内仅存在待完成、已过期的，展示红点
	// 绿点 --已完成，红点---待完成，
	for _, value := range data {
		var info planReq.PlansListInfo
		info = value
		if item, ok := mapList[value.Id]; !ok { // 先判断是否完成，在判断是否到期
			if info.IsFinished == microScrmModel.PLAN_FINISHED && info.Status == microScrmModel.PLAN_FISISHING { // 已完成
				info.FinishNum = 1
				info.Status = microScrmModel.PLAN_FINISHED
			}
			if info.EndTime < nowTime && (info.Status == microScrmModel.PLAN_FISISHING || info.Status == microScrmModel.PLAN_NOFINISH) { // 未完成，已结束,
				info.Status = microScrmModel.PLAN_NOFINISH
			}
			info.AllNum = 1
			mapList[value.Id] = info
		} else {
			info.AllNum = item.AllNum
			info.FinishNum = item.FinishNum
			info.AllNum++
			if item.Status == microScrmModel.PLAN_NOFINISH {
				continue
			}
			if info.IsFinished == microScrmModel.PLAN_FINISHED && item.Status == microScrmModel.PLAN_FINISHED { // 如果子任务都是完成的状态，就跳过
				info.FinishNum++
				continue
			}
			if info.IsFinished == microScrmModel.PLAN_FISISHING && item.Status == microScrmModel.PLAN_FINISHED { // 子任务待完成，有已完成的任务，则更改任务状态为待完成
				info.Status = microScrmModel.PLAN_FISISHING
			}
			if info.IsFinished == microScrmModel.PLAN_NOFINISH && item.Status == microScrmModel.PLAN_FINISHED { // 子任务已过期(未完成)，任务状态是已完成，则状态更新为已完成
				info.Status = microScrmModel.PLAN_FINISHED
			}
			mapList[value.Id] = info
		}
		if !helper.InArray(value.Id, planIds) {
			planIds = append(planIds, value.Id)
		}
	}
	// 持续任务的列表查询
	taskTypeIds, err := planDao.NewPlanTaskDao().GetTaskTypePlanIdByPlanIds(c, request.EnterpriseId, planIds)
	if err != nil {
		return nil, err
	}
	// 查询是企业群发类型的任务
	qySendIds, err := planDao.NewPlanTaskDao().GetMany(c, request.EnterpriseId, map[string]interface{}{
		"task_category":   microScrmModel.PLANTASKCONTACT,
		"group_send_type": microScrmModel.PLANTASKQYSEND,
		"plan_id":         planIds,
	})
	if err != nil {
		return nil, err
	}

	planTaskIds := make(map[int64][]int64, 0)
	for _, v := range qySendIds {
		planTaskIds[v.PlanID] = append(planTaskIds[v.PlanID], v.ID)
	}
	var noFinishedPlanIds []int64
	for k, model := range planTaskIds {
		planId := k
		taskIds := model
		// 获取有企业群发的沟通任务类型的主任务
		num, err := planDao.NewPlanTaskConditionDataDao().GetNoFinishedPlanIdsByTaskId(c, request.EnterpriseId, planId, memberInfo.Userid, taskIds, request.StartTime, request.EndTime)
		if err != nil {
			continue
		}
		if num > 1 {
			noFinishedPlanIds = append(noFinishedPlanIds, planId)
		}
	}

	// 查询已将关闭触达提示的主任务id
	closePlanIds, err := planDao.NewPlanTipsRelationDao().GetClosePlanIds(c, request.EnterpriseId, memberInfo.Userid, noFinishedPlanIds)
	if err != nil {
		return nil, err
	}
	if request.Status >= 0 {
		for _, v := range mapList {
			if v.Status != request.Status { // 如果不相等，那么记录进入到不展示的数据中
				if _, ok := unShowIds[v.Id]; !ok {
					unShowIds[v.Id] = true
				}
			}
		}
	}
	for _, v := range mapList {
		info := v
		if _, ok := unShowIds[v.Id]; !ok {
			if helper.InArray(info.Id, taskTypeIds) {
				info.TaskType = 2
			}
			info.CloseTips = 1
			if info.Status == microScrmModel.PLAN_FINISHED {
				if helper.InArray(info.Id, noFinishedPlanIds) && !helper.InArray(info.Id, closePlanIds) {
					info.CloseTips = 0
				}
			}
			result = append(result, info)
		}
	}
	// 数组排序,asc案例
	t := result
	sort.Slice(t, func(i, j int) bool {
		if t[i].Status == t[j].Status {
			if t[i].EndTime == t[j].EndTime {
				return t[i].Id < t[j].Id
			} else {
				return t[i].EndTime < t[j].EndTime
			}
		} else {
			return t[i].Status < t[j].Status
		}
	})
	return result, nil
}

// GetPlanDetail 获取主任务详情
func (p PlanTask) GetPlanDetail(c ctx.Context, request planReq.PlanInfo) (data planReq.PlanInfoRes, err error) {
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: request.EnterpriseId},
		Mobile:         request.Mobile,
	})
	if err != nil || memberInfo.ID == 0 {
		return planReq.PlanInfoRes{}, errorx.NewErrorMsg(1001, "查询导购信息异常")
	}
	data, err = p.planTaskDao.GetPlanDetail(c, request)
	if err != nil {
		return planReq.PlanInfoRes{}, err
	}
	data.TaskType = 0
	taskType, err := p.planTaskDao.GetTaskTypeInfo(c, request)
	if err != nil {
		return data, err
	}
	if taskType > 0 {
		data.TaskType = 1
	}
	finishNum, taskNum, err := planDao.NewPlanTaskAllotLogDao().GetFinishNumByPlanId(c, request.EnterpriseId, memberInfo.Userid, request.Id)
	data.FinishNum = finishNum
	data.TaskNum = strconv.Itoa(int(taskNum))
	return data, err
}

// GetPlanTipsCount 获取当天的未完成任务数
func (p PlanTask) GetPlanTipsCount(c ctx.Context, request planReq.TipsReq) (planReq.TipsRes, error) {
	var result planReq.TipsRes
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: request.EnterpriseId},
		Mobile:         request.Mobile,
	})
	if err != nil || memberInfo.ID == 0 {
		return result, errorx.NewErrorMsg(1001, "查询导购信息异常")
	}
	planNum, err := planDao.NewPlanTaskAllotLogDao().GetPlanTipsCount(c, request, memberInfo.Userid)
	if err != nil {
		return result, err
	}
	result.PlanNums = planNum
	reserveNum, err := externalUserDao.NewReservesOrder().GetPlanTipsCount(c, request.EnterpriseId, request.StartTime, request.EndTime, memberInfo.ID)
	if err != nil {
		return result, err
	}
	result.ReserveNums = reserveNum
	sowTaskNum, err := customerDao.NewSowTask().GetPlanTipsCount(c, request.EnterpriseId, request.StartTime, request.EndTime, memberInfo.ID)
	if err != nil {
		return result, err
	}
	result.SowTaskNums = sowTaskNum
	return result, err
}

// LeaveGroupChat 退群人数
func (p PlanTask) LeaveGroupChat(c ctx.Context, request planReq.LeaveGroupChatReq) (total int64, result []customerReq.ExternalUserTagsInfoListsRes, err error) {
	taskAllotInfo, err := planDao.NewPlanTaskAllotLogDao().GetTaskAllotGroupInfoById(c, request.EnterpriseId, request.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil, err
	}
	one, err := planDao.NewPlanTaskConditionDataDao().GetOne(c, request.EnterpriseId, taskAllotInfo.PlanId, map[string]interface{}{
		"task_allot_log_id": request.Id,
	})
	if err != nil {
		return 0, nil, err
	}
	leaveUnionIds, err := workGroupDao.NewWorkGroupChatMember().GetChatUnionIds(c, request.EnterpriseId, one.GroupChatId, taskAllotInfo.StartTime, taskAllotInfo.EndTime, 2)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil, err
	}
	if len(leaveUnionIds) == 0 {
		return 0, nil, err
	}

	total = int64(len(leaveUnionIds))
	// todo 调整处理退群用户列表查询逻辑

	userTagsInfo, err := planDao.NewUserRelation().GetUserTagInfoByUnionIds(c, taskAllotInfo.UserId, leaveUnionIds, request)
	return total, userTagsInfo, err
}

// SaCustomerLists sacode&&unionid|customerid的列表
func (p PlanTask) SaCustomerLists(c ctx.Context, request planReq.CustomerListsReq) (total int64, result []planReq.CustomerListsRes, err error) {
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id": request.Id,
	})
	if err != nil {
		return 0, nil, err
	}
	relationList, err := planDao.NewUserRelation().GetMany(c, map[string]interface{}{
		"enterprise_id": request.EnterpriseId,
		"user_id":       one.UserId,
		"deleted_at":    0,
	})
	if err != nil {
		return 0, nil, err
	}
	var relationUnionIds []string
	var searchRelationUnionIds []string
	var allRelationUnionIds []string
	var remarkMap = make(map[string]string)
	for _, relation := range relationList {
		if relation.RelationStatus == customer.WorkUserRelationStatusNormal {
			relationUnionIds = append(relationUnionIds, relation.Unionid)
			remarkMap[relation.Unionid] = relation.Remark
		}
		if request.Keywords != "" {
			if strings.Contains(relation.Remark, request.Keywords) {
				searchRelationUnionIds = append(searchRelationUnionIds, relation.Unionid)
			}
		}
		allRelationUnionIds = append(allRelationUnionIds, relation.Unionid)
	}
	var searchNickUnionid []string
	if request.Keywords != "" {
		corpExternalList, _ := planDao.NewWxUserCorpExternalCk().GetMany(c, request.EnterpriseId, "nick like ? AND unionid in ?", "%"+request.Keywords+"%", allRelationUnionIds)
		for _, external := range corpExternalList {
			searchNickUnionid = append(searchNickUnionid, external.Unionid)
		}
		if len(searchNickUnionid) == 0 && len(searchRelationUnionIds) == 0 {
			return 0, nil, err
		}
	}

	total, lists, err := planDao.NewPlanTaskConditionDataDao().GetCustomerLists(c, request, one.PlanId, relationUnionIds, helper.ArrayUniq(helper.ArrayMerge(searchNickUnionid, searchRelationUnionIds)))
	if err != nil {
		return 0, nil, err
	}
	unionIdInfo := make(map[string]struct {
		Name   string
		Avatar string
	}, 0)
	var unionIds []string
	for _, value := range lists {
		info := value
		res := planReq.CustomerListsRes{
			TaskConditionDataModel: info,
			CustomerInfoRes: planReq.CustomerInfoRes{
				Name:   info.UserName,
				Avatar: "https://cdn-cos-uat.coach.com.cn/plan_task/default.png",
			},
		}

		unionIds = append(unionIds, info.UnionId)
		if remark, ok := remarkMap[info.UnionId]; ok {
			res.RelRemark = remark
		}
		result = append(result, res)
	}

	qyTagIdMap := make(map[string]string)
	tagLists, _ := workTagSvc.NewWorkTags().GetChoiceTagList(c, request.EnterpriseId)
	for _, model := range tagLists {
		qyTagIdMap[model.TagID] = model.Name
	}
	if len(unionIds) > 0 {
		external, err := planDao.NewWxUserCorpExternal().GetMany(c, request.EnterpriseId, map[string]interface{}{
			"unionid": unionIds,
		})
		if err != nil {
			return 0, nil, err
		}
		for _, model := range external {
			item := model
			if _, ok := unionIdInfo[item.Unionid]; !ok {
				unionIdInfo[item.Unionid] = struct {
					Name   string
					Avatar string
				}{Name: item.Nick, Avatar: item.Avatar}
			}
		}
		unionIdTags := make(map[string][]string)
		unionIdTagIds := make(map[string][]string)
		if request.IsTag == 1 { // 查询的用户展示标签
			TagsList, err := planDao.NewPlanTaskConditionDataDao().GetConditionTagsByUnionIds(c, request.EnterpriseId, one.PlanId, request.Id, unionIds)
			if err != nil {
				return 0, nil, err
			}
			for _, v := range TagsList {
				tag := v
				if tag.WwTagId != "" {
					if !helper.InArray[string](tag.WwTagId, unionIdTagIds[tag.UnionId]) {
						tagName := qyTagIdMap[tag.WwTagId]
						unionIdTags[tag.UnionId] = append(unionIdTags[tag.UnionId], tagName)
					}
				}
			}
		}
		for i, res := range result {
			v := res
			if v.UnionId != "" {
				if item, ok := unionIdInfo[v.UnionId]; ok {
					v.Avatar = item.Avatar
					v.Name = item.Name
				}
				if !helper.InArray(v.UnionId, relationUnionIds) {
					v.IsUnion = 0
				} else {
					v.IsUnion = 1
				}
			}
			if request.IsTag == 1 {
				if _, ok := unionIdTags[v.UnionId]; ok {
					v.TagName = unionIdTags[v.UnionId]
				}
			} else {
				if v.TaskConditionDataModel.WwTagId != "" {
					if value, ok := qyTagIdMap[v.TaskConditionDataModel.WwTagId]; ok {
						v.TagName = append(v.TagName, value)
					}
				}
			}
			result[i] = v
		}
	}
	return total, result, nil
}

// AddMaPlan 添加MA任务
func (p PlanTask) AddMaPlan(c ctx.Context, enterpriseId int64, request openReq.AddPlanReq) (planId, taskId int64, err error) {
	var info microScrmModel.PlanTaskTemplateModel
	if request.PlanTemplateId > 0 {
		info, err = planDao.NewPlanTemplateDao().GetTemplateInfo(c, enterpriseId, request.PlanTemplateId)
		if err != nil {
			return planId, taskId, errors.New("任务模板不存在")
		}
		if info.Status != 1 {
			return planId, taskId, errors.New("任务模板已禁用")
		}
	}
	planId, taskId, err = p.planTaskDao.CreateMaPlanTask(c, enterpriseId, request, info)
	if err != nil {
		return planId, taskId, err
	}
	return planId, taskId, err
}

func (p PlanTask) CreateMaPlanByContent(c ctx.Context, enterpriseId int64, request openReq.CreatePlanByContentReq) (int64, error) {
	startTime := helper.GetDateToTimes(request.PlanStartTime)
	endTime := helper.GetDateToTimes(request.PlanEndTime)
	if startTime <= 0 {
		return 0, errors.New("任务开始时间有误")
	}
	if endTime <= 0 {
		return 0, errors.New("任务结束时间有误")
	}
	if startTime >= endTime {
		return 0, errors.New("任务开始时间必须早于结束时间")
	}
	if endTime <= time.Now().Unix() {
		return 0, errors.New("任务结束时间必须晚于当前时间")
	}
	return p.planTaskDao.CreatePlanByContent(c, enterpriseId, request, startTime, endTime)
}

// CreateMaTask creates a Shenyan MA sub task from request content.
func (p PlanTask) CreateMaTask(c ctx.Context, enterpriseId int64, request openReq.CreateTaskReq) (int64, error) {
	planInfo, err := p.planTaskDao.GetPlanDetail(c, planReq.PlanInfo{
		BaseReq: planReq.BaseReq{EnterpriseId: enterpriseId},
		Id:      request.PlanId,
	})
	if err != nil {
		return 0, errors.New("任务不存在")
	}
	if planInfo.Id == 0 {
		return 0, errors.New("主任务不存在")
	}
	if taskId, err := p.getTaskIdByPlanAndCrowd(c, enterpriseId, request.PlanId, request.CrowdId); err == nil && taskId > 0 {
		if err := p.pushCreateTaskCrowd(c, enterpriseId, request, taskId); err != nil {
			return 0, err
		}
		return taskId, nil
	} else if err != nil {
		return 0, err
	}

	taskStartTime := helper.GetDateToTimes(request.TaskStartTime)
	taskEndTime := helper.GetDateToTimes(request.TaskEndTime)
	if taskStartTime <= 0 {
		return 0, errors.New("子任务开始时间格式不正确")
	}
	if taskEndTime <= 0 {
		return 0, errors.New("子任务结束时间格式不正确")
	}
	if taskStartTime >= taskEndTime {
		return 0, errors.New("子任务开始时间必须早于结束时间")
	}
	if taskEndTime <= time.Now().Unix() {
		return 0, errors.New("子任务结束时间必须晚于当前时间")
	}
	if taskStartTime < planInfo.StartTime {
		return 0, errors.New("子任务开始时间不能早于主任务开始时间")
	}
	if taskEndTime > planInfo.EndTime {
		return 0, errors.New("子任务结束时间不能晚于主任务结束时间")
	}

	templateContent, err := p.validateShenyanTemplateContent(c, request)
	if err != nil {
		return 0, err
	}

	lockKey := fmt.Sprintf("shenyan:task:lock:%d:%s", request.PlanId, request.CrowdId)
	lockValue := strconv.FormatInt(app.UID(), 10)
	locked := app.Redis().SetNX(c, lockKey, lockValue, 30*time.Second).Val()
	if !locked {
		time.Sleep(200 * time.Millisecond)
		if taskId, err := p.getTaskIdByPlanAndCrowd(c, enterpriseId, request.PlanId, request.CrowdId); err == nil && taskId > 0 {
			if err := p.pushCreateTaskCrowd(c, enterpriseId, request, taskId); err != nil {
				return 0, err
			}
			return taskId, nil
		}
		return 0, errors.New("任务创建中，请稍后重试")
	}
	defer app.Redis().Del(c, lockKey)

	if taskId, err := p.getTaskIdByPlanAndCrowd(c, enterpriseId, request.PlanId, request.CrowdId); err == nil && taskId > 0 {
		if err := p.pushCreateTaskCrowd(c, enterpriseId, request, taskId); err != nil {
			return 0, err
		}
		return taskId, nil
	} else if err != nil {
		return 0, err
	}

	var tagName, wxTagId string
	var tagType, tagExpireType int
	var tagId int64
	if request.TaskCategory == microScrmModel.PLANTASKMOMENT {
		tagExpireType = microScrmModel.TAG_TASK_SEND_EXPIRE
		if request.TaskTag == "" {
			return 0, errors.New("朋友圈任务标签不能为空")
		}
		tagName = request.TaskTag
		wxTagId, tagId, err = workTagSvc.NewWorkTags().CreateTaskTag(c, enterpriseId, "任务标签", tagName, 0, tagType)
		if err != nil {
			return 0, errors.New("朋友圈任务标签创建失败:" + err.Error())
		}
	}

	taskId := app.UID()
	if err = p.planTaskDao.CreateTask(c, request, enterpriseId, taskId, templateContent, tagExpireType, tagName, wxTagId); err != nil {
		return 0, err
	}
	if request.TaskCategory == microScrmModel.PLANTASKMOMENT {
		relation := workModels.PlanTasksTagRelationModel{
			ID:           app.UID(),
			EnterpriseID: enterpriseId,
			PlanID:       taskId,
			TagID:        tagId,
			StartTime:    int(helper.GetDateToTimes(request.TaskStartTime)),
			EndTime:      int(helper.GetDateToTimes(request.TaskEndTime)),
		}
		if _, err = workGroupDao.NewPlanTasksTagRelation().CreateData(c, relation); err != nil {
			return 0, fmt.Errorf("创建标签关系失败: %w", err)
		}
	}
	if err := p.pushCreateTaskCrowd(c, enterpriseId, request, taskId); err != nil {
		return 0, err
	}
	return taskId, nil
}

func (p PlanTask) pushCreateTaskCrowd(c ctx.Context, enterpriseId int64, request openReq.CreateTaskReq, taskId int64) error {
	if len(request.Crowd) == 0 {
		return nil
	}
	return p.PushMaTaskCrowd(c, enterpriseId, openReq.PushTaskCrowdReq{
		TaskId:  taskId,
		CrowdId: request.CrowdId,
		Crowd:   request.Crowd,
	})
}

func (p PlanTask) PushMaTaskCrowd(c ctx.Context, enterpriseId int64, request openReq.PushTaskCrowdReq) error {
	if len(request.Crowd) == 0 {
		return errors.New("人群包不能为空")
	}
	for _, item := range request.Crowd {
		if item.UnionID == "" || item.SaCode == "" {
			return errors.New("人群包必须传unionid和sa_code")
		}
	}
	taskInfo, err := p.planTaskDao.GetOne(c, enterpriseId, map[string]interface{}{"id": request.TaskId})
	if err != nil {
		return errors.New("子任务不存在")
	}
	if taskInfo.ID == 0 {
		return errors.New("子任务不存在")
	}
	if taskInfo.CrowdID != request.CrowdId {
		return errors.New("人群包ID与子任务不匹配")
	}
	crowdReq := openReq.TaskCrowdReq{
		TaskId:       taskInfo.PlanID,
		CrowdId:      request.CrowdId,
		PlanSendType: int(taskInfo.SendType),
		Crowd:        request.Crowd,
	}
	return p.maTaskCrowdInsert(c, enterpriseId, crowdReq, request.TaskId)
}

func (p PlanTask) getTaskIdByPlanAndCrowd(c ctx.Context, enterpriseId, planId int64, crowdId string) (int64, error) {
	result, err := p.planTaskDao.GetOne(c, enterpriseId, map[string]interface{}{
		"plan_id":  planId,
		"crowd_id": crowdId,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	return result.ID, nil
}

func (p PlanTask) validateShenyanTemplateContent(c ctx.Context, request openReq.CreateTaskReq) (string, error) {
	if request.GroupSendType != microScrmModel.PLANTASKQYSEND && request.GroupSendType != microScrmModel.PLANTASK1V1SEND {
		return "", errors.New("group_send_type must be 1 or 2")
	}
	if request.TaskCategory != microScrmModel.PLANTASKCONTACT && request.TaskCategory != microScrmModel.PLANTASKMOMENT {
		return "", errors.New("invalid task_category")
	}
	if request.TaskCategory == microScrmModel.PLANTASKMOMENT && request.GroupSendType != microScrmModel.PLANTASKQYSEND {
		return "", errors.New("moment task does not support group_send_type 2")
	}
	if len(request.TemplateContent) == 0 {
		request.TemplateContent = map[string]interface{}{}
	}
	content := helper.MarshalToString(request.TemplateContent)
	var conditionContent planReq.ConditionContent
	if err := json.Unmarshal([]byte(content), &conditionContent); err != nil {
		return "", errors.New("invalid template_content")
	}
	if request.GroupSendType != microScrmModel.PLANTASK1V1SEND && conditionContent.Text.Content == "" {
		return "", errors.New("text.content is required")
	}
	if err := validateShenyanSingleMainMaterial(conditionContent, request.TemplateContentType); err != nil {
		return "", err
	}

	switch request.TemplateContentType {
	case 1:
		return content, nil
	case 2:
		if len(conditionContent.Image.ImgList) == 0 {
			return "", errors.New("image.imgList is required")
		}
		if len(conditionContent.Image.ImgList) > 9 {
			return "", errors.New("image.imgList count cannot exceed 9")
		}
		for index, item := range conditionContent.Image.ImgList {
			if !isShenyanHTTPURL(item.PicUrl) {
				return "", errors.New("image.pic_url must be an http or https url")
			}
			finalURL, _ := normalizeShenyanImageURL(item.PicUrl)
			if err := validateShenyanImageMedia(c, finalURL, request.TaskCategory == microScrmModel.PLANTASKMOMENT); err != nil {
				return "", fmt.Errorf("image.pic_url invalid: %w", err)
			}
			conditionContent.Image.ImgList[index].PicUrl = finalURL
		}
	case 4:
		if !isShenyanHTTPURL(conditionContent.Video.VideoUrl) {
			return "", errors.New("video.video_url must be an http or https url")
		}
		if err := validateShenyanVideoMedia(c, conditionContent.Video.VideoUrl); err != nil {
			return "", fmt.Errorf("video.video_url invalid: %w", err)
		}
	case 6:
		if conditionContent.Link.Title == "" {
			return "", errors.New("link.title is required")
		}
		if !isShenyanHTTPURL(conditionContent.Link.LinkUrl) {
			return "", errors.New("link.link_url must be an http or https url")
		}
		if !isShenyanHTTPURL(conditionContent.Link.PicUrl) {
			return "", errors.New("link.pic_url must be an http or https url")
		}
		finalURL, _ := normalizeShenyanImageURL(conditionContent.Link.PicUrl)
		if err := validateShenyanImageMedia(c, finalURL, false); err != nil {
			return "", fmt.Errorf("link.pic_url invalid: %w", err)
		}
		conditionContent.Link.PicUrl = finalURL
		if conditionContent.Link.CoverUrl != "" && !isShenyanHTTPURL(conditionContent.Link.CoverUrl) {
			return "", errors.New("link.cover_url must be an http or https url")
		}
		if conditionContent.Link.CoverUrl != "" {
			finalURL, _ := normalizeShenyanImageURL(conditionContent.Link.CoverUrl)
			if err := validateShenyanImageMedia(c, finalURL, false); err != nil {
				return "", fmt.Errorf("link.cover_url invalid: %w", err)
			}
			conditionContent.Link.CoverUrl = finalURL
		}
	case 7:
		if request.TaskCategory == microScrmModel.PLANTASKMOMENT {
			return "", errors.New("moment task does not support mini program")
		}
		if conditionContent.Applets.Title == "" || conditionContent.Applets.AppId == "" || conditionContent.Applets.Path == "" {
			return "", errors.New("applets.title, app_id and path are required")
		}
		if !isShenyanHTTPURL(conditionContent.Applets.CoverUrl) {
			return "", errors.New("applets.cover_url must be an http or https url")
		}
		finalURL, _ := normalizeShenyanImageURL(conditionContent.Applets.CoverUrl)
		if err := validateShenyanImageMedia(c, finalURL, false); err != nil {
			return "", fmt.Errorf("applets.cover_url invalid: %w", err)
		}
		conditionContent.Applets.CoverUrl = finalURL
	default:
		return "", errors.New("unsupported template_content_type")
	}

	return helper.MarshalToString(conditionContent), nil
}

func validateShenyanSingleMainMaterial(content planReq.ConditionContent, contentType int32) error {
	if contentType != 2 && len(content.Image.ImgList) > 0 {
		return errors.New("only one main material type is allowed")
	}
	if contentType != 4 && content.Video.VideoUrl != "" {
		return errors.New("only one main material type is allowed")
	}
	if contentType != 6 && (content.Link.Title != "" || content.Link.LinkUrl != "" || content.Link.PicUrl != "") {
		return errors.New("only one main material type is allowed")
	}
	if contentType != 7 && (content.Applets.AppId != "" || content.Applets.Path != "" || content.Applets.CoverUrl != "") {
		return errors.New("only one main material type is allowed")
	}
	return nil
}

func isShenyanHTTPURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}
func (p PlanTask) MaTaskCrowd(c ctx.Context, enterpriseId int64, request openReq.TaskCrowdReq) error {
	if helper.GetDateToTimes(request.PlanSendTime) <= 0 {
		return errors.New("任务发送时间有误")
	}
	planInfo, err := p.planTaskDao.GetPlanDetail(c, planReq.PlanInfo{
		BaseReq: planReq.BaseReq{EnterpriseId: enterpriseId},
		Id:      request.TaskId,
	})
	if err != nil {
		return errors.New("任务不存在")
	}
	info, err := planDao.NewPlanTemplateDao().GetTemplateInfo(c, enterpriseId, request.PlanTemplateId)
	if err != nil {
		return errors.New("任务模板不存在")
	}
	if info.Status != 1 {
		return errors.New("任务模板已禁用")
	}
	if planInfo.OperationUserID == 0 { // 主任务的创建人，默认为子任务模版的创建人
		_ = p.planTaskDao.UpdatePlanData(c, enterpriseId, map[string]interface{}{
			"id": planInfo.Id,
		}, map[string]interface{}{
			"operation_user_id": info.AdminID,
		})
	}

	taskStartTime := helper.GetDateToTimes(request.TaskStartTime)
	taskEndTime := helper.GetDateToTimes(request.TaskEndTime)
	if taskStartTime > 0 && taskStartTime < planInfo.StartTime {
		return errors.New("子任务开始时间不能早于主任务开始时间")
	}
	if taskEndTime > planInfo.EndTime {
		return errors.New("子任务结束时间不能晚于主任务结束时间")
	}
	if taskStartTime == 0 || taskEndTime == 0 { // 如果MA任务传入的结束时间是0，则引用主任务的开始时间和结束时间，（兼容历史的MA任务）
		request.TaskStartTime = helper.GetFormatDateByParams(planInfo.StartTime)
		request.TaskEndTime = helper.GetFormatDateByParams(planInfo.EndTime)
	}

	taskId, err := p.GetTaskIdByCrowdId(c, enterpriseId, request, info)
	if err != nil {
		return err
	}
	err = p.maTaskCrowdInsert(c, enterpriseId, request, taskId)
	if err != nil {
		log.Info(c).Msgf("MA人群包任务写入失败,任务id:%s", request.CrowdId)
	}
	return nil
}

// GetTaskIdByCrowdId 持续任务根据CrowdId 查询子任务是否创建
func (p PlanTask) GetTaskIdByCrowdId(c ctx.Context, enterpriseId int64, request openReq.TaskCrowdReq, info microScrmModel.PlanTaskTemplateModel) (taskId int64, err error) {
	crowdId := request.CrowdId
	redisKey := fmt.Sprintf("mataskId:%s", crowdId)

	// 1. 先查缓存
	taskIdStr := app.Redis().Get(c, redisKey).Val()
	if taskIdStr != "" {
		return strconv.ParseInt(taskIdStr, 10, 64)
	}

	// 2. 使用分布式锁保护业务创建过程
	lockKey := fmt.Sprintf("mataskId:lock:%s", crowdId)
	lockValue := strconv.FormatInt(app.UID(), 10) // 转为字符串

	locked := app.Redis().SetNX(c, lockKey, lockValue, 30*time.Second).Val()
	if !locked {
		// 其他请求正在创建，等待后重试
		time.Sleep(200 * time.Millisecond)
		taskIdStr = app.Redis().Get(c, redisKey).Val()
		if taskIdStr != "" {
			return strconv.ParseInt(taskIdStr, 10, 64)
		}
		return 0, fmt.Errorf("任务创建中，请稍后重试")
	}
	defer app.Redis().Del(c, lockKey)

	// 3. 双重检查：再次检查缓存和数据库
	taskIdStr = app.Redis().Get(c, redisKey).Val()
	if taskIdStr != "" {
		return strconv.ParseInt(taskIdStr, 10, 64)
	}

	// 查询数据库
	result, err := p.planTaskDao.GetOne(c, enterpriseId, map[string]interface{}{
		"plan_id":  request.TaskId,
		"crowd_id": request.CrowdId,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("查询任务失败: %w", err)
	}

	// 如果数据库已有记录，直接返回并更新缓存
	if result.ID != 0 {
		// 更新缓存
		app.Redis().Set(c, redisKey, strconv.FormatInt(result.ID, 10), 24*time.Hour)
		return result.ID, nil
	}

	// 4. 数据库没有记录，开始创建新任务
	var backUpInfo microScrmModel.PlanTaskTemplateBackupModel
	if info.OpenBackUp == microScrmModel.OpenBackYes { //有开启备用模板
		backUpInfo, err = planDao.NewPlanTemplateDao().GetTemplateBackUpInfo(c, enterpriseId, info.ID)
		if err != nil {
			return 0, fmt.Errorf("获取备用模板失败: %w", err)
		}
	}

	var tagName, wxTagId string
	var tagType, tagExpireType int

	// 判断任务是不是要生成标签的任务
	if info.TaskCategory == microScrmModel.TaskCategoryMoment { // 如果标签任务群发,并且要添加标签
		tagExpireType = microScrmModel.TAG_TASK_SEND_EXPIRE
		if request.TaskTag == "" {
			return 0, errors.New("朋友圈任务标签不能为空:" + helper.Int64ToString(request.TaskId))
		}
		tagName = request.TaskTag
		var tagId int64
		wxTagId, tagId, err = workTagSvc.NewWorkTags().CreateTaskTag(c, enterpriseId, "任务标签", tagName, 0, tagType)
		if err != nil {
			return 0, errors.New("朋友圈任务标签创建失败:" + err.Error())
		}

		// 保存标签关系
		tagRelationId := app.UID()
		tasksTagsRelation := workModels.PlanTasksTagRelationModel{
			ID:           tagRelationId,
			EnterpriseID: enterpriseId,
			PlanID:       info.ID,
			TagID:        tagId,
			StartTime:    int(helper.GetDateToTimes(request.TaskStartTime)),
			EndTime:      int(helper.GetDateToTimes(request.TaskEndTime)),
		}
		_, err = workGroupDao.NewPlanTasksTagRelation().CreateData(c, tasksTagsRelation)
		if err != nil {
			return 0, fmt.Errorf("创建标签关系失败: %w", err)
		}
	}

	// 生成新任务ID
	newTaskId := app.UID()

	// 创建任务
	err = p.planTaskDao.CreateTaskByCrowdId(c, request, enterpriseId, newTaskId, info, backUpInfo, tagExpireType, tagName, wxTagId)
	if err != nil {
		return 0, fmt.Errorf("创建任务失败: %w", err)
	}

	// 5. 业务成功后才缓存
	app.Redis().Set(c, redisKey, strconv.FormatInt(newTaskId, 10), 24*time.Hour)

	return newTaskId, nil
}

// ma人群包的插入
func (p PlanTask) maTaskCrowdInsert(c ctx.Context, enterpriseId int64, request openReq.TaskCrowdReq, id int64) error {
	taskId := id
	var unionIds, userIds []string
	var customerIds []string
	externalMaps := make(map[string]string)
	externalUserIdMaps := make(map[string]string)
	membersMaps := make(map[string]string)
	var data []*planReq.UnionIDData
	var datalist []planReq.ImportInsertData
	lists := make(map[string]planReq.ImportInsertData, 0)
	for _, v := range request.Crowd {
		info := v
		if request.PlanSendType == 2 {
			if info.UnionID == "" && info.CustomerId != "" {
				customerIds = append(customerIds, info.CustomerId)
			}
		}
		data = append(data, &info)
	}
	if len(customerIds) > 0 { // 计算出unionid为空，customerid有值的数据
		externalResult, err := planDao.NewWxUserCorpExternal().GetMany(c, enterpriseId, map[string]interface{}{
			"union_code": customerIds,
		})
		if err != nil {
			return err
		}
		if len(externalResult) > 0 { // 根据customerid找到库中的unionid用户
			for _, model := range externalResult {
				if _, ok := externalMaps[model.UnionCode]; !ok {
					externalMaps[model.UnionCode] = model.Unionid
				}
				externalUserIdMaps[model.UnionCode] = model.Openid
			}
		}
	}
	memberLists, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"deleted_at":    0,
	})
	if err != nil {
		return err
	}
	for _, info := range memberLists { // 获取未离职的导购列表
		membersMaps[info.Userid] = info.Name
	}
	logDataList := make([]microScrmModel.PlanTaskCrowdDataLogModel, 0)
	importData := make([]microScrmModel.PlanTaskCrowdDataModel, 0)
	if request.PlanSendType == 2 { // sacode&&unionid的关系匹配

		for _, v := range data {
			unionIds = append(unionIds, v.UnionID)
			if !helper.InArray[string](v.SaCode, userIds) {
				userIds = append(userIds, v.SaCode)
			}
		}
		list, err := planDao.NewPlanImportDao().GetInsertImportDataInfoBySaCodeAndUnionId(c, enterpriseId, data)
		if err != nil {
			return err
		}
		datalist = list
		for _, item := range datalist {
			key := item.UnionID + "_" + item.UserID
			lists[key] = item
		}
		for _, v := range data {
			ImportRes := "匹配失败"
			var saName, workExternalUserId string
			said := v.SaCode
			unionId := v.UnionID
			customerId := v.CustomerId
			if unionId == "" && customerId != "" {
				isMatch := 1
				customerInfo, err2 := microwork.NewCrm().GetCustomerInfoByUnionId(c, microwork.UserInfoRequest{
					EnterpriseID: enterpriseId,
					CustomerCode: customerId,
				})
				if err2 != nil {
					isMatch = 2
				}
				if customerInfo.Data.Mobile == "" || said == "" {
					isMatch = 2
				}
				if v, ok := externalMaps[customerId]; ok {
					unionId = v
				}
				if v, ok := externalUserIdMaps[customerId]; ok {
					workExternalUserId = v
				}
				userName := customerInfo.Data.Name
				if item, ok := membersMaps[said]; ok { // said存在插入数据
					saName = item
					newData := microScrmModel.PlanTaskCrowdDataModel{
						Name:               saName,
						SaCode:             v.SaCode,
						Unionid:            v.UnionID,
						CrowdID:            request.CrowdId,
						WorkExternalUserID: workExternalUserId,
						CustomerID:         v.CustomerId,
						TaskID:             taskId,
						BaseModel:          models.CreateBase(app.UID()),
						UserName:           userName,
						IsMatch:            isMatch,
					}
					if isMatch == 1 {
						ImportRes = "匹配成功"
					}
					importData = append(importData, newData)
				}
			} else {
				pickKey := unionId + "_" + said
				if _, ok := lists[pickKey]; ok {
					workExternalUserId = lists[pickKey].WorkExternalUserId
					if item, ok := membersMaps[said]; ok { // said存在插入数据
						saName = item
						newData := microScrmModel.PlanTaskCrowdDataModel{
							Name:               saName,
							SaCode:             v.SaCode,
							Unionid:            v.UnionID,
							CrowdID:            request.CrowdId,
							WorkExternalUserID: workExternalUserId,
							CustomerID:         v.CustomerId,
							TaskID:             taskId,
							BaseModel:          models.CreateBase(app.UID()),
							UserName:           lists[pickKey].UserName,
							IsMatch:            1,
						}
						ImportRes = "匹配成功"
						importData = append(importData, newData)
					}
				}
			}
			logData := microScrmModel.PlanTaskCrowdDataLogModel{
				Name:       saName,
				SaCode:     v.SaCode,
				Unionid:    v.UnionID,
				ImportRes:  ImportRes,
				CrowdID:    request.CrowdId,
				TaskID:     taskId,
				BaseModel:  models.CreateBase(app.UID()),
				CustomerID: v.CustomerId,
			}
			logDataList = append(logDataList, logData)
		}
		err = planDao.NewPlanTaskCrowdDataDao().InsertCrowdDataBatch(c, importData, enterpriseId)
		err = planDao.NewPlanTaskCrowdDataLogDao().InsertICrowdDataLogBatch(c, logDataList, enterpriseId)
	}
	if request.PlanSendType == 1 { // 根据unionid的人群包导入
		lists, err := planDao.NewPlanImportDao().GetInsertImportDataInfoByUnionId(c, enterpriseId, data, 2)
		if err != nil {
			return err
		}
		keySaRelation := make(map[string]planReq.ImportInsertData, 0)
		for _, value := range lists {
			item := value
			keySaIndex := item.UserID + item.UnionID
			keySaRelation[keySaIndex] = item
		}
		for _, v := range data {
			ImportRes := "匹配失败"
			item := v
			keySaIndex := item.SaCode + item.UnionID
			saName := ""
			if value, ok := keySaRelation[keySaIndex]; ok {
				info := value
				if memberName, ok := membersMaps[info.UserID]; ok {
					saName = memberName
				}
				newData := microScrmModel.PlanTaskCrowdDataModel{
					Name:               saName,
					SaCode:             info.UserID,
					Unionid:            info.UnionID,
					CrowdID:            request.CrowdId,
					WorkExternalUserID: info.WorkExternalUserId,
					TaskID:             taskId,
					BaseModel:          models.CreateBase(app.UID()),
					CustomerID:         "",
					IsMatch:            1,
				}
				ImportRes = "匹配成功"
				importData = append(importData, newData)
			}
			logData := microScrmModel.PlanTaskCrowdDataLogModel{
				Name:       saName,
				SaCode:     item.SaCode,
				Unionid:    item.UnionID,
				ImportRes:  ImportRes,
				CrowdID:    request.CrowdId,
				TaskID:     taskId,
				BaseModel:  models.CreateBase(app.UID()),
				CustomerID: item.CustomerId,
			}
			logDataList = append(logDataList, logData)
		}
		err = planDao.NewPlanTaskCrowdDataDao().InsertCrowdDataBatch(c, importData, enterpriseId)
		err = planDao.NewPlanTaskCrowdDataLogDao().InsertICrowdDataLogBatch(c, logDataList, enterpriseId)
	}
	// 如果需要给人群包打标签，开始处理标签
	taskInfo, err := planDao.NewPlanTaskDao().GetOne(c, enterpriseId, map[string]interface{}{"id": taskId})
	if err != nil {
		return err
	}
	if taskInfo.TaskCategory == microScrmModel.PLANTASKMOMENT && taskInfo.GroupSendType == microScrmModel.PLANTASKQYSEND {
		var externalUserTags []tagsReq.ExternalUserTag
		for _, value := range importData {
			item := value
			if item.Unionid == "" {
				continue
			}
			externalUserTags = append(externalUserTags, tagsReq.ExternalUserTag{
				TagId:   taskInfo.QyTagId,
				UnionId: item.Unionid,
				UserId:  item.SaCode,
			})

		}
		if len(externalUserTags) > 0 {
			chunks := lo.Chunk(externalUserTags, 500)
			for _, chunk := range chunks {
				req := tagsReq.BatchMarkingTagReq{
					ExternalUserTagList: chunk,
					TaskID:              taskId,
					EnterpriseID:        enterpriseId,
				}

				str, _ := json.Marshal(req)
				err = pkgKafka.GetAppKafka().Producer(string(str), "batch_marking_tag")
			}
		}
	}
	return err
}

// MaFinishedPush MA任务完成之后，任务完成结果推送
func (p PlanTask) MaFinishedPush(c ctx.Context, enterpriseId, taskId int64, infoMap map[string]string) {
	_ = planDao.NewPlanTaskDao().UpdateMaPlanFinishPush(c, enterpriseId, taskId, microScrmModel.FINISHPUSHMAING)
	taskInfo, err := p.planTaskDao.GetOne(c, enterpriseId, map[string]interface{}{
		"id": taskId,
	})
	originType := taskInfo.OriginType
	if err != nil {
		return
	}
	err = planDao.NewPlanTaskAllotLogDao().UpdateData(c, enterpriseId, map[string]interface{}{
		"is_finished": microScrmModel.PLAN_FISISHING,
		"task_id":     taskId,
	}, map[string]interface{}{
		"is_finished": microScrmModel.PLAN_NOFINISH,
	})

	planInfo, err := planDao.NewPlanTaskDao().GetPlanDetail(c, planReq.PlanInfo{
		BaseReq: planReq.BaseReq{EnterpriseId: enterpriseId},
		Id:      taskInfo.PlanID,
	})
	if err != nil {
		return
	}
	var templateName string
	if taskInfo.PlanTemplateId > 0 {
		templateInfo, err := planDao.NewPlanTemplateDao().GetTemplateInfo(c, enterpriseId, taskInfo.PlanTemplateId)
		if err != nil {
			return
		}
		templateName = templateInfo.TemplateName
	}

	conditionDataDao := planDao.NewPlanTaskConditionDataDao()
	page := 1
	size := 100
	for {
		lists, err := conditionDataDao.GetListsByTaskId(c, enterpriseId, taskInfo.PlanID, taskInfo.ID, page, size)
		if err != nil {
			log.Info(c).Msgf("获取任务数据推送给MA的任务完成数据失败%s", err)
			continue
		}
		if len(lists) == 0 {
			break
		}
		var events []microwork.MaEvents
		for _, v := range lists {
			item := v
			date := time.Unix(taskInfo.TaskEndTime, 0).UTC().Format("2006-01-02T15:04:05Z")
			finishStr := "已过期"
			if v.IsFinished == 1 {
				finishStr = "已完成"
				date = time.Unix(item.UpdatedAt, 0).UTC().Format("2006-01-02T15:04:05Z")
			}
			saName := ""
			if memberName, ok := infoMap[v.UserId]; ok {
				saName = memberName
			}
			event := microwork.MaEvents{
				IdentityType:      "wechat-unionid",
				IdentityValue:     item.UnionId,
				Date:              date,
				CResult:           finishStr,
				CSaName:           saName,
				CSaCode:           v.UserId,
				CTaskId:           strconv.FormatInt(v.TaskId, 10),
				CTaskName:         taskInfo.Name,
				CCrowdId:          taskInfo.CrowdID,
				CTaskRemark:       taskInfo.Remark,
				FlowId:            taskInfo.FlowId,
				Campaign:          taskInfo.Campaign,
				CPlanTemplateId:   strconv.FormatInt(taskInfo.PlanTemplateId, 10),
				CPlanTemplateName: templateName,
				CTaskClassifation: planInfo.Name,
			}
			events = append(events, event)
		}
		sendMAPlanResultData(c, enterpriseId, originType, events)
		page++
	}
	err = planDao.NewPlanTaskDao().UpdateMaPlanFinishPush(c, enterpriseId, taskId, microScrmModel.FINISHPUSHMAEND)
	if err != nil {
		fmt.Printf("更新MA推送结果回执报错%s\n", err.Error())
	}

}

// GetQyGroupListsByPlanId 根据主任务id获取有效期内的群发任务类型的子任务列表
func (p PlanTask) GetQyGroupListsByPlanId(c ctx.Context, enterpriseId int64, groupSendType int, planIds []int64) ([]microScrmModel.PlanTasksModel, error) {
	return p.planTaskDao.GetQyGroupListsByPlanId(c, enterpriseId, groupSendType, planIds)
}

// GetPlanIdByStatus 根据任务状态获取主任务id
func (p PlanTask) GetPlanIdByStatus(c ctx.Context, enterpriseId int64, status int) ([]int64, error) {
	return p.planTaskDao.GetPlanIdByStatus(c, enterpriseId, status)
}

// QyGroupSendResult 企业任务群发的结果触达
func (p PlanTask) QyGroupSendResult(ctx ctx.Context, enterpriseId, taskId int64, info microScrmModel.TaskAllotLogModel, qyWorkApi *work.QyWorkImpl, taskMap map[int64]int) {
	memberInfo, err := membersDao.NewMembers().GetOne(ctx, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"userid":        info.UserId,
		"deleted_at":    0,
	})
	if err != nil {
		return
	}
	var msgIdList []string
	if info.GroupMsgId == "-1" {
		msgList, errs := planDao.NewPlanTaskAllotLogDao().GetTaskAllotMsg(ctx, enterpriseId, map[string]interface{}{
			"task_allot_log_id": info.ID,
		})
		if errs != nil {
			return
		}
		for _, model := range msgList {
			msgIdList = append(msgIdList, model.MsgId)
		}
	} else {
		msgIdList = append(msgIdList, info.GroupMsgId)
	}
	var groupSendMap []bool
	contactNum := 0
	for _, msgId := range msgIdList {
		isGroupSend := false
		var massHairReq = work.GroupMsgSendResultRequest{}
		massHairReq.MsgID = msgId
		massHairReq.UserID = info.UserId
		massHairReq.Limit = 100
		nextCursor := ""
		originType := taskMap[info.TaskId]
		for {
			if nextCursor != "" {
				massHairReq.Cursor = nextCursor
			}
			response, err := qyWorkApi.GetGroupMsgSendResult(ctx, massHairReq)
			if err != nil {
				resErr := err.(*errorx.Error)
				app.Logger().Error(ctx).Msgf("获取群发消息失败, err:%s", err.Error())
				if resErr.Code == 60111 || resErr.Code == 41093 { // 导购已离职，群发消息的状态更新为已取消，消息已取消，也需要更新
					err = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseId, map[string]interface{}{
						"id": info.ID,
					}, map[string]interface{}{"is_group_send": microScrmModel.QyGroupMsgSendCancel})
				}
				break
			}
			app.Logger().Error(ctx).Err(err).Interface("params", massHairReq).Msg("qyGroupSend.qyGroupSend")
			updateExternalStatusMap := make(map[int][]string, 0)
			maSendExternalUserId := make([]string, 0)
			var sendTime int64
			if response.SendList != nil {
				for _, item := range response.SendList {
					if item.Status > 0 && !isGroupSend {
						isGroupSend = true
					}
					if item.Status == 0 {
						continue
					}
					if item.Status == microScrmModel.TouchStatusYes && sendTime == 0 {
						sendTime = item.SendTime
					}
					updateExternalStatusMap[item.Status] = append(updateExternalStatusMap[item.Status], item.ExternalUserID)
					if originType > 0 && item.Status == microScrmModel.TouchStatusYes { // 如果是MA任务，则需要推送任务完成结果给MA
						taskConditionInfo, err := planDao.NewPlanTaskConditionDataDao().GetOne(ctx, enterpriseId, info.PlanId, map[string]interface{}{
							"task_allot_log_id":     info.ID,
							"work_external_user_id": item.ExternalUserID,
						})
						if err != nil {
							app.Logger().Error(ctx).Msgf("MA任务:%d查询导购%s的客户%s触达失败，错误为:%s", taskId, info.UserId, item.ExternalUserID, err.Error())
						}
						if taskConditionInfo.IsFinished == microScrmModel.FinishedNo { // 查询还是未完成，则需要更新完成的消息通知
							maSendExternalUserId = append(maSendExternalUserId, item.ExternalUserID)
							md5Str := helper.Md5V(taskConditionInfo.UnionId + strconv.FormatInt(info.ID, 10) + taskConditionInfo.CustomerID)
							excitationSvc.NewPointsProduceService().ProduceQueue(ctx, excitationReq.PointsProduceSave{
								MemberId:     memberInfo.ID,
								Type:         excitationModels.RuleCategoryId9,
								Value:        1,
								UniqueHash:   md5Str + fmt.Sprintf("-%d-%d", memberInfo.ID, excitationModels.RuleCategoryId9),
								EnterpriseId: enterpriseId,
							})
						}
					}
				}
			}
			for status, externalIds := range updateExternalStatusMap {
				updateMap := map[string]interface{}{
					"is_touch": status,
				}
				if status == microScrmModel.TouchStatusYes { // 推送成功触达的；
					updateMap["is_finished"] = microScrmModel.FinishedYes
					updateMap["finished_at"] = sendTime
					contactNum += len(externalIds)
				}
				err := planDao.NewPlanTaskConditionDataDao().UpdateData(ctx, enterpriseId, info.PlanId, map[string]interface{}{
					"task_allot_log_id":     info.ID,
					"work_external_user_id": externalIds,
				}, updateMap)
				if err != nil {
					app.Logger().Error(ctx).Msgf("子任务:%d的导购%s群发给客户消息触达更新失败，错误为:%s", taskId, info.UserId, err.Error())
					continue
				}
			}
			if len(maSendExternalUserId) > 0 {
				queueData := planReq.MaPushFinishQueueData{
					EnterpriseId:    enterpriseId,
					TaskId:          info.TaskId,
					PlanId:          info.PlanId,
					ExternalUserIds: maSendExternalUserId,
					UserId:          info.UserId,
					TaskAllotLogId:  info.ID,
				}
				NewPlanTaskQueueSvc().MaPushFinishProduce(ctx, queueData)
			}
			nextCursor = response.NextCursor
			if nextCursor == "" {
				break
			}
		}
		groupSendMap = append(groupSendMap, isGroupSend)
	}
	isGroupAllSend := lo.EveryBy(groupSendMap, func(send bool) bool {
		return send
	})

	if isGroupAllSend { // 当导购点击了【群发助手】中的【发送】后，子任务状态自动变更为【已完成】，任务名单中，如果群发失败的客户增加群发失败的原因提醒，同时客户状态为【标记已联系】可手动点击
		updateMap := map[string]interface{}{"is_group_send": microScrmModel.QyGroupMsgSendFinish}
		err = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseId, map[string]interface{}{
			"id": info.ID,
		}, updateMap)
		excitationSvc.NewPointsProduceService().TaskFinished(ctx, excitationReq.TaskFinishedReq{
			TaskAllotLogId: info.ID,
			EnterpriseId:   enterpriseId,
			MemberId:       memberInfo.ID,
		})
		err = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseId, map[string]interface{}{
			"id":          info.ID,
			"is_finished": microScrmModel.PLAN_FISISHING,
		}, map[string]interface{}{"is_finished": microScrmModel.PLAN_FINISHED, "contacted_count": contactNum, "end_time": time.Now().Unix()})
		if err != nil {
			app.Logger().Error(ctx).Msgf("子任务:%d的导购%s群发结果状态更新失败，错误为:%s", taskId, info.UserId, err.Error())
		}
	}
	// 更新任务的完成率
	finishCount, err := planDao.NewPlanTaskConditionDataDao().GetGroupCount(ctx, enterpriseId, info.PlanId, map[string]interface{}{
		"task_allot_log_id": info.ID,
		"is_finished":       microScrmModel.FinishedYes,
	})
	if err != nil {
		return
	}
	finishNum := decimal.NewFromInt(finishCount)
	customerNum := decimal.NewFromInt(int64(info.CustomerNum))
	divNum := finishNum.Div(customerNum).Round(2).Mul(decimal.NewFromInt(100)).IntPart()
	err = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseId, map[string]interface{}{
		"id": info.ID,
	}, map[string]interface{}{"finish_rate": divNum, "contacted_count": finishCount})
}

// sendMAPlanResultData 任务发送结果回执给到MA任务
func sendMAPlanResultData(c ctx.Context, enterpriseId int64, originType int, lists []microwork.MaEvents) {
	params := microwork.SendMaFinishReq{
		EnterpriseID: enterpriseId,
		OriginType:   originType,
		Events:       lists,
		CommonProps: microwork.MACommonProps{
			Event: "c_sa_completed_task_result",
		},
	}
	err := macallback.SendFinish(c, params)
	app.Logger().Info(c).Err(err).Interface("params", params).Msg("请求MA的任务发送结果回执")

}

// QyMomentSendResult 朋友圈任务群发的结果触达
func (p PlanTask) QyMomentSendResult(ctx ctx.Context, enterpriseId int64, taskDetail microScrmModel.PlanTasksModel, qyWorkApi *work.QyWorkImpl) error {
	nowTime := time.Now().Unix()
	momentId := taskDetail.MomentJobId
	req := work.GetMomentTaskRequest{}
	req.MomentId = momentId
	for {
		res, err := qyWorkApi.GetMomentTask(ctx, req)
		if err != nil {
			app.Logger().Error(ctx).Msgf("获取群发消息失败, err:%s", err.Error())
			return err
		}
		if res.ErrCode != 0 {
			app.Logger().Error(ctx).Msgf("获取群发消息企微返回错误, response:%v", res)
			return err
		}
		req.Cursor = &res.NextCursor
		for _, v := range res.TaskList {
			info := v
			if info.PublishStatus == 1 {
				// 更新task_allot_log 表数据
				one, err := planDao.NewPlanTaskAllotLogDao().GetOne(ctx, enterpriseId, map[string]interface{}{
					"task_id": taskDetail.ID,
					"user_id": info.UserId,
				})
				if err != nil {
					app.Logger().Error(ctx).Msgf("查询导购任务的记录失败, err:%v", err)
					continue
				}

				if one.IsFinished == microScrmModel.PLAN_FINISHED {
					if one.EndTime+7200 > int(time.Now().Unix()) {
						_ = p.GetMomentSendResult(ctx, enterpriseId, taskDetail, momentId, one, info, qyWorkApi)
					}
					continue
				} else if one.IsFinished == microScrmModel.PLAN_FISISHING {
					// 进行中状态：获取结果，
					_ = p.GetMomentSendResult(ctx, enterpriseId, taskDetail, momentId, one, info, qyWorkApi)
				}
				// 任务完成加积分
				memberInfo, err := membersDao.NewMembers().GetOne(ctx, map[string]interface{}{"enterprise_id": enterpriseId, "userid": info.UserId})
				if err != nil {
					continue
				}
				excitationSvc.NewPointsProduceService().TaskFinished(ctx, excitationReq.TaskFinishedReq{
					TaskAllotLogId: one.ID,
					EnterpriseId:   enterpriseId,
					MemberId:       memberInfo.ID,
				})
				err = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseId, map[string]interface{}{
					"id":          one.ID,
					"is_finished": microScrmModel.PLAN_FISISHING,
				}, map[string]interface{}{
					"is_finished": microScrmModel.PLAN_FINISHED,
					"end_time":    nowTime,
				})
				if err != nil {
					app.Logger().Error(ctx).Msgf("更新导购任务完成状态失败, err:%v", err)
					continue
				}

			}
		}
		if res.NextCursor == "" {
			break
		}
	}
	return nil

}
func (p PlanTask) GetMomentSendResult(ctx ctx.Context, enterpriseId int64, taskDetail microScrmModel.PlanTasksModel, momentId string, one microScrmModel.TaskAllotLogModel, info work.GetMomentTaskInfo, qyWorkApi *work.QyWorkImpl) error {
	// 查询朋友圈任务的可见范围列表客户，更新顾客状态为已联系
	sendRequest := work.GetMomentSendTaskRequest{}
	sendRequest.UserId = info.UserId
	sendRequest.MomentId = momentId
	for {
		result, err := qyWorkApi.GetMomentSendResult(ctx, sendRequest)
		if err != nil {
			app.Logger().Error(ctx).Msgf("查询朋友圈任务的可见范围列表客户报错, err:%s", err.Error())
			return err
		}
		if result.ErrCode != 0 {
			app.Logger().Error(ctx).Msgf("查询朋友圈任务的可见范围列表客户企微返回错误, response:%v", result)
			return err
		}
		var externalUserIds []string
		for _, value := range result.CustomerList {
			externalUserIds = append(externalUserIds, value.ExternalUserId)
		}
		//查询朋友圈任务的可见范围列表客户的状态是否为已联系
		touchOpenIds, _ := planDao.NewPlanTaskConditionDataDao().GetOpenIds(ctx, enterpriseId, one.PlanId,
			"task_allot_log_id = ? AND work_external_user_id in ? AND is_touch = ?",
			one.ID, externalUserIds, microScrmModel.TouchStatusYes)
		// 如果朋友圈任务的可见范围列表客户的状态为已联系，则从外部用户id中去除
		if len(touchOpenIds) > 0 && len(externalUserIds) > 0 {
			externalUserIds = helper.GenericArrayDiff(externalUserIds, touchOpenIds)
		}
		//更新朋友圈任务的可见范围列表客户的状态为已联系
		err = planDao.NewPlanTaskConditionDataDao().UpdateData(ctx, enterpriseId, one.PlanId, map[string]interface{}{
			"task_allot_log_id":     one.ID,
			"work_external_user_id": externalUserIds,
		}, map[string]interface{}{
			"is_touch": microScrmModel.TouchStatusYes,
		})
		if err != nil {
			app.Logger().Error(ctx).Msgf("更新朋友圈可见范围用户为可展示, err:%v", err)
			continue
		}
		// 发送朋友圈任务的可见范围列表客户的回执给到MA
		p.SendMAMomentSendResultData(ctx, enterpriseId, taskDetail, info.UserId, externalUserIds)
		sendRequest.Cursor = &result.NextCursor
		if result.NextCursor == "" {
			break
		}
	}
	// 朋友圈任务完成，计算可见人数，计算朋友圈任务的触达率，moment_look_count/customer_num
	momentCount, err := planDao.NewPlanTaskConditionDataDao().GetGroupCount(ctx, enterpriseId, one.PlanId, map[string]interface{}{
		"task_allot_log_id": one.ID,
		"is_touch":          microScrmModel.TouchStatusYes,
	})
	if err != nil {
		return err
	}
	momentLookCount := decimal.NewFromInt(momentCount)
	customerNum := decimal.NewFromInt(int64(one.CustomerNum))
	divNum := momentLookCount.Div(customerNum).Round(2).Mul(decimal.NewFromInt(100)).IntPart()
	err = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseId, map[string]interface{}{
		"id": one.ID,
	}, map[string]interface{}{"finish_rate": divNum, "moment_look_count": momentCount})
	return nil
}

// SendMAMomentSendResultData  朋友圈任务发送结果及可见范围回执给到MA
func (p PlanTask) SendMAMomentSendResultData(c ctx.Context, enterpriseId int64, taskDetail microScrmModel.PlanTasksModel, SaCode string, externalUserIds []string) {
	sendMaMomentResultReq := microwork.SendMaMomentResultReq{
		EnterpriseID:        enterpriseId,
		SaCode:              SaCode,
		TaskId:              taskDetail.ID,
		CrowdId:             taskDetail.CrowdID,
		Campaign:            taskDetail.Campaign,
		CustomerUnionidList: make([]string, 0),
	}
	if len(externalUserIds) > 0 {
		unionIds, _ := planDao.NewWxUserCorpExternal().GetUnionIdsByOpenIds(c, enterpriseId, externalUserIds)
		if len(unionIds) > 0 {
			sendMaMomentResultReq.CustomerUnionidList = unionIds
		}
	}
	err := macallback.SendMomentResult(c, taskDetail.OriginType, sendMaMomentResultReq)
	log.Info(c).Err(err).Interface("params", sendMaMomentResultReq).Msg("请求MA的任务发送结果回执")
}

// GetPlanStat 获取主任务数据
func (p PlanTask) GetPlanStat(c ctx.Context, request planReq.NewStatReq) (planReq.NewStatRes, error) {
	var result planReq.NewStatRes
	var err error
	detail, err := planDao.NewPlanTaskDao().GetPlanDetail(c, planReq.PlanInfo{
		BaseReq: planReq.BaseReq{
			EnterpriseId: request.EnterpriseId,
		},
		Id: request.PlanId})
	if err != nil {
		return result, err
	}
	startTime := request.StartTime
	endTime := request.EndTime
	result.PlanName = detail.Name
	result.PlanStartTime = detail.StartTime
	result.PlanEndTime = detail.EndTime
	taskNum, err := planDao.NewPlanTaskDao().GetCount(c, request.EnterpriseId, request.PlanId, startTime, endTime)
	if err != nil {
		return result, err
	}
	result.PeriodSubtasksNumber = taskNum
	subtasksFinishNumber := int64(0)
	if detail.Status > microScrmModel.PlanNotStart {
		many, err := planDao.NewPlanTaskDao().GetMany(c, request.EnterpriseId, map[string]interface{}{"plan_id": request.PlanId})
		if err != nil {
			return result, err
		}
		var taskIds []int64
		for _, v := range many {
			taskIds = append(taskIds, v.ID)
		}
		taskFinishLists, err := planDao.NewPlanTaskAllotLogDao().GetTaskNotFinishNumByPlanId(c, request.EnterpriseId, request.PlanId)
		if err != nil {
			return result, err
		}
		for _, list := range taskFinishLists {
			if list.Num == 1 {
				subtasksFinishNumber++
			}
		}
		result.SubtasksFinishNumber = subtasksFinishNumber
	}
	subtasksFinishPercentage := 0.00
	if result.PeriodSubtasksNumber > 0 {
		subtasksFinishPercentage, _ = decimal.NewFromInt(subtasksFinishNumber).DivRound(decimal.NewFromInt(result.PeriodSubtasksNumber), 2).Float64()
	}
	result.SubtasksFinishPercentage = fmt.Sprintf("%2.f%%", subtasksFinishPercentage*100)
	result.TaskNumber, err = planDao.NewPlanTaskAllotLogDao().GetCount(c, request.EnterpriseId, map[string]interface{}{
		"plan_id":    request.PlanId,
		"deleted_at": 0,
	})
	result.TaskFinishNumber, err = planDao.NewPlanTaskAllotLogDao().GetCount(c, request.EnterpriseId, map[string]interface{}{
		"plan_id":     request.PlanId,
		"deleted_at":  0,
		"is_finished": microScrmModel.PLAN_FINISHED,
	})
	taskFinishPercentage := 0.00
	if result.TaskNumber > 0 {
		taskFinishPercentage, _ = decimal.NewFromInt(result.TaskFinishNumber).DivRound(decimal.NewFromInt(result.TaskNumber), 2).Float64()
	}
	result.TaskFinishPercentage = fmt.Sprintf("%2.f%%", taskFinishPercentage*100)
	result.ReachGuideNumber, err = planDao.NewPlanTaskAllotLogDao().GetSaCount(c, request.EnterpriseId, map[string]interface{}{
		"plan_id":    request.PlanId,
		"deleted_at": 0,
	})
	ckDao := planDao.NewPlanTaskConditionDataCkDao()

	result.MatchCustomersNumber, err = ckDao.GetCount(c, request.EnterpriseId, request.PlanId, map[string]interface{}{
		"deleted_at": 0,
	})
	result.MatchCustomersNumber2, err = ckDao.CustomerDistinctCount(c, request.EnterpriseId, request.PlanId, map[string]interface{}{
		"deleted_at": 0,
	})
	result.ReachLordNumber, err = ckDao.GetCountByArg(c, request.EnterpriseId, request.PlanId, "chat_num>0 and deleted_at=0")

	result.MatchGroupNumber, err = ckDao.GetCountByArg(c, request.EnterpriseId, request.PlanId, "group_chat_id>0 and deleted_at=0")

	return result, err
}

// GetPlanStatDetail 主任务的数据明细
func (p PlanTask) GetPlanStatDetail(c ctx.Context, request planReq.NewStatDetailReq) (int64, []*planReq.NewStatDetailResult, error) {
	var result []*planReq.NewStatDetailResult

	total, many, err := p.planTaskDao.GetStatLists(c, request)
	if err != nil {
		return total, result, err
	}

	wg := threading.NewRoutineGroup()
	for _, v := range many {
		info := planReq.NewStatDetailResult{
			Id:            v.ID,
			TaskName:      v.Name,
			TaskCategory:  v.TaskCategory,
			GroupSendType: v.GroupSendType,
			Dimension:     int(v.Dimension),
		}
		fn := func(value *planReq.NewStatDetailResult, taskId int64) func() {
			return func() {
				value.TaskNumber, err = planDao.NewPlanTaskAllotLogDao().GetCount(c, request.EnterpriseId, map[string]interface{}{
					"task_id":    taskId,
					"deleted_at": 0,
				})
				value.TaskFinishNumber, err = planDao.NewPlanTaskAllotLogDao().GetCount(c, request.EnterpriseId, map[string]interface{}{
					"task_id":     taskId,
					"deleted_at":  0,
					"is_finished": microScrmModel.PLAN_FINISHED,
				})
				taskFinishPercentage := 0.00
				if value.TaskNumber > 0 {
					taskFinishPercentage, _ = decimal.NewFromInt(value.TaskFinishNumber).DivRound(decimal.NewFromInt(value.TaskNumber), 2).Float64()
				}
				value.TaskFinishPercentage = fmt.Sprintf("%2.f%%", taskFinishPercentage*100)
				value.ReachGuideNumber, err = planDao.NewPlanTaskAllotLogDao().GetSaCount(c, request.EnterpriseId, map[string]interface{}{
					"task_id":    taskId,
					"deleted_at": 0,
				})

				ckDao := planDao.NewPlanTaskConditionDataCkDao()

				value.MatchCustomersNumber, err = ckDao.GetCount(c, request.EnterpriseId, request.PlanId, map[string]interface{}{
					"task_id":    taskId,
					"deleted_at": 0,
				})
				value.MatchCustomersNumber2, err = ckDao.CustomerDistinctCount(c, request.EnterpriseId, request.PlanId, map[string]interface{}{
					"task_id":    taskId,
					"deleted_at": 0,
				})
				value.ReachLordNumber, err = ckDao.GetCountByArg(c, request.EnterpriseId, request.PlanId, "task_id=? and chat_num>0 and deleted_at=0", taskId)

				value.MatchGroupNumber, err = ckDao.GetCountByArg(c, request.EnterpriseId, request.PlanId, "task_id=? and group_chat_id>0 and deleted_at=0", taskId)

				if err != nil {
					log.Error(c).Msgf("主任务数据明细报错，err:%v", err)
				}
			}
		}
		wg.RunSafe(fn(&info, v.ID))
		result = append(result, &info)
	}
	wg.Wait()
	return total, result, err
}

// CloseTouchTips 已完成的主任务未触达关闭
func (p PlanTask) CloseTouchTips(c ctx.Context, request planReq.PlanInfo) error {
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: request.EnterpriseId},
		Mobile:         request.Mobile,
	})
	if err != nil || memberInfo.ID == 0 {
		return errorx.NewErrorMsg(1001, "查询导购信息异常")
	}
	data := microScrmModel.PlanTipsRelationModel{
		BaseModel: models.BaseModel{
			ID: app.UID(),
		},
		PlanID: request.Id,
		UserID: memberInfo.Userid,
	}
	err = planDao.NewPlanTipsRelationDao().CreateData(c, request.EnterpriseId, data)
	return err
}

// CloseTaskTouchTips 已完成的子任务未触达关闭
func (p PlanTask) CloseTaskTouchTips(c ctx.Context, request planReq.PlanInfo) error {
	memberInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: request.EnterpriseId},
		Mobile:         request.Mobile,
	})
	if err != nil || memberInfo.ID == 0 {
		return errorx.NewErrorMsg(1001, "查询导购信息异常")
	}
	err = planDao.NewPlanTaskAllotLogDao().UpdateData(c, request.EnterpriseId, map[string]interface{}{
		"id": request.Id,
	}, map[string]interface{}{
		"close_tips": 1,
	})
	return err
}

func (p PlanTask) TaskMomentDelTag(c ctx.Context, params planReq.MomentDelTagQueueData) error {
	time.Sleep(time.Minute) // 延迟一分钟执行，针对朋友圈任务的临时标签延迟删除
	enterpriseId := params.EnterpriseId
	taskInfo, err := p.planTaskDao.GetOne(c, enterpriseId, map[string]interface{}{
		"id": params.TaskId,
	})
	if err != nil {
		return err
	}
	qyTagId := taskInfo.QyTagId
	tagsInfo, err := tagsDao.NewWorkTagsDao().GetTagInfoByQyId(c, enterpriseId, qyTagId)
	if err != nil {
		return err
	}
	err = tagsDao.NewWorkTagsDao().DelTagByQyId(c, enterpriseId, qyTagId)
	if err != nil {
		return err
	}
	info, err := common.GetEnterpriseInfo(c, enterpriseId)
	if err != nil {
		return err
	}
	qyWorkApi := new(work.QyWorkImpl).SetDomain(info.Domain)
	delQyTag := work.DelCorpTagParam{
		TagID: []string{qyTagId},
	}
	err = qyWorkApi.DelCorpTagInfo(c, delQyTag)
	if err != nil {
		return err
	}
	err = tagsDao.NewWorkExternalUserTag().DeleteData(c, map[string]interface{}{
		"enterprise_id":   enterpriseId,
		"external_tag_id": tagsInfo.ID,
	})
	if err != nil {
		return err
	}
	return err
}

// MaPushFinish 队列消费处理MA任务完成的推送
func (p PlanTask) MaPushFinish(ctx ctx.Context, info planReq.MaPushFinishQueueData) error {
	enterpriseId := info.EnterpriseId
	planInfo, err := planDao.NewPlanTaskDao().GetPlanDetail(ctx, planReq.PlanInfo{
		BaseReq: planReq.BaseReq{
			EnterpriseId: enterpriseId,
		},
		Id: info.PlanId,
	})
	if err != nil {
		log.Info(ctx).Msgf("标记联系MA回执查询主任务信息报错 %+v", err)
		return err
	}
	if !helper.InArray(planInfo.Category, []int{microScrmModel.PlanMaTypeYes, microScrmModel.PlanMaTypeNewYes}) { // 非MA任务不推送
		return nil
	}
	taskInfo, err := planDao.NewPlanTaskDao().GetOne(ctx, enterpriseId, map[string]interface{}{
		"id": info.TaskId,
	})
	if err != nil {
		log.Info(ctx).Msgf("标记联系MA回执查询子任务信息报错 %+v", err)
		return err
	}
	originType := taskInfo.OriginType
	var templateName string
	if taskInfo.PlanTemplateId > 0 {
		templateInfo, err := planDao.NewPlanTemplateDao().GetTemplateInfo(ctx, enterpriseId, taskInfo.PlanTemplateId)
		if err != nil {
			return err
		}
		templateName = templateInfo.TemplateName
	}
	var events []microwork.MaEvents
	date := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	memberInfo, err := membersDao.NewMembers().GetOne(ctx, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"userid":        info.UserId,
		"deleted_at":    0,
	})
	if err != nil {
		log.Info(ctx).Msgf("查询导购信息身份失败 %+v", err)
		return err
	}
	page := 1
	size := 100
	for {
		lists, err := planDao.NewPlanTaskConditionDataDao().GetManyByPagination(ctx, enterpriseId, info.PlanId, page, size, map[string]interface{}{
			"task_allot_log_id":     info.TaskAllotLogId,
			"work_external_user_id": info.ExternalUserIds,
		})
		if len(lists) == 0 {
			break
		}
		if err != nil {
			return err
		}
		for _, v := range lists {
			item := v
			event := microwork.MaEvents{
				IdentityType:      "wechat-unionid",
				IdentityValue:     item.UnionId,
				Date:              date,
				CResult:           "已完成",
				CSaName:           memberInfo.Name,
				CSaCode:           info.UserId,
				CTaskId:           strconv.FormatInt(info.TaskId, 10),
				CTaskName:         taskInfo.Name,
				CCrowdId:          taskInfo.CrowdID,
				CTaskRemark:       taskInfo.Remark,
				CPlanTemplateId:   strconv.FormatInt(taskInfo.PlanTemplateId, 10),
				CPlanTemplateName: templateName,
				CTaskClassifation: planInfo.Name,
				FlowId:            taskInfo.FlowId,
				Campaign:          taskInfo.Campaign,
			}
			events = append(events, event)
		}
		params := microwork.SendMaFinishReq{
			EnterpriseID: enterpriseId,
			OriginType:   originType,
			Events:       events,
			CommonProps: microwork.MACommonProps{
				Event: "c_sa_completed_task_result",
			},
		}
		err = macallback.SendFinish(ctx, params)
		if err != nil {
			log.Info(ctx).Err(err).Interface("params", params).Msg("标记联系请求MA的任务发送结果回执")
		}
		page++
	}

	return nil
}

// SearchExternalLists coach的任务计划升级，获取任务客户信息
func (p PlanTask) SearchExternalLists(c ctx.Context, request planReq.ExternalSearchListsReq, enterpriseConfigName string) (int64, []planReq.CustomerListsRes, error) {
	if request.Size == 0 {
		request.Size = 10
	}
	result := make([]planReq.CustomerListsRes, 0)
	var total int64

	// 获取任务分配日志
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id": request.Id,
	})
	if err != nil {
		return total, result, err
	}

	// 获取条件数据列表
	list, err := planDao.NewPlanTaskConditionDataDao().GetSearchExternalLists(c, request, one.PlanId)
	if err != nil || len(list) == 0 {
		return total, result, err
	}

	// 收集UnionIds
	for _, model := range list {
		request.Search.UnionIds = append(request.Search.UnionIds, model.UnionId)
	}

	if len(request.Search.UnionIds) == 0 {
		return 0, result, nil
	}

	// 根据IsSearch判断走哪个逻辑
	if request.Search.IsSearch == 1 {
		// 走CK检索逻辑
		return p.searchExternalListsWithCK(c, request, enterpriseConfigName, one.PlanId)
	} else {
		// 走非检索逻辑
		return p.getExternalLists(c, request, enterpriseConfigName, one.PlanId)
	}
}

// searchExternalListsWithCK 使用ClickHouse进行搜索
func (p PlanTask) searchExternalListsWithCK(c ctx.Context, request planReq.ExternalSearchListsReq, enterpriseConfigName string, planId int64) (int64, []planReq.CustomerListsRes, error) {
	result := make([]planReq.CustomerListsRes, 0)
	var total int64
	var nowUnix = time.Now().Unix()
	// 设置CK查询参数
	request.Search.EnterpriseID = request.EnterpriseId
	request.Search.MemberId = request.MemberId
	request.Search.IsTaskList = true

	// 从CK获取外部用户列表
	externalList, _, err := externalSearchDao.NewExternalExtraCk().ExternalSearchListCk(c, customerSearchReq.ExternalSearchListBaseReq{
		SearchList: request.Search,
		GetAll:     true,
	})
	if err != nil {
		return total, result, err
	}

	var unionIds []string
	var externalMap = make(map[string]customerSearchResp.ExternalUserList)
	for _, v := range externalList {
		unionIds = append(unionIds, v.Unionid)
		externalMap[v.Unionid] = v
	}

	if len(unionIds) == 0 {
		return 0, result, nil
	}

	// 根据CK返回的unionIds查询任务条件数据
	total, lists, err := planDao.NewPlanTaskConditionDataDao().GetSearchExternalListsByUnionIds(c, request, planId, unionIds)
	if err != nil {
		return total, result, err
	}

	cardDescriptionMap, _ := customerSvc.NewCustomer().GetCustomerCardMapForType(enterpriseConfigName)

	// 获取当前页的unionIds
	var pageUnionIds []string
	for _, model := range lists {
		pageUnionIds = append(pageUnionIds, model.UnionId)
	}

	// 获取用户备注关系
	relationList, _ := externalUserDao.NewUserRelation().GetMany(c, "enterprise_id = ? and member_id = ? AND unionid in ?",
		request.EnterpriseId, request.MemberId, pageUnionIds)

	// 构建结果
	for _, value := range lists {
		info := value
		res := planReq.CustomerListsRes{
			TaskConditionDataModel: info,
			CustomerInfoRes: planReq.CustomerInfoRes{
				Avatar: "https://cdn-cos-uat.coach.com.cn/plan_task/default.png",
			},
		}

		if item, ok := externalMap[value.UnionId]; ok {
			res.IsUnion = helper.RuleThree(item.SelfRelationStatus > 0, 0, 1)
			res.CustomerInfoRes.Name = item.Nick
			res.CustomerInfoRes.Avatar = item.Avatar
			res.CustomerInfoRes.VipType = item.VipType
			if item.TaskFinishAt > 0 {
				if nowUnix-int64(item.TaskFinishAt) < 86400*customerSearchReq.TaskHaveValidDays {
					res.CustomerInfoRes.TaskFinishAt = item.TaskFinishAt
				}
			}
			if item.ValidReserveAt > 0 {
				if nowUnix-int64(item.ValidReserveAt) < 86400*customerSearchReq.ValidDays {
					res.CustomerInfoRes.ValidReserveAt = item.ValidReserveAt
				}
			}
			if item.SuccessReserveAt > 0 {
				if nowUnix-int64(item.SuccessReserveAt) < 86400*customerSearchReq.ValidDays {
					res.CustomerInfoRes.SuccessReserveAt = item.SuccessReserveAt
				}
			}
			res.CustomerInfoRes.CouponStatus = item.CouponStatus
			res.CustomerInfoRes.CouponBy30daysStatus = item.CouponBy30daysStatus
			if cardName, o := cardDescriptionMap[item.CardLevel]; o {
				res.CustomerInfoRes.CustomerCard = cardName
			}
		}

		// 处理备注信息
		for _, relation := range relationList {
			if relation.Unionid == value.UnionId {
				res.CustomerInfoRes.RelRemark = relation.Remark
			}
		}

		result = append(result, res)
	}

	return total, result, nil
}

// searchExternalListsWithMySQL 使用MySQL进行搜索
func (p PlanTask) getExternalLists(c ctx.Context, request planReq.ExternalSearchListsReq, enterpriseConfigName string, planId int64) (int64, []planReq.CustomerListsRes, error) {
	result := make([]planReq.CustomerListsRes, 0)
	var nowUnix = time.Now().Unix()
	// 直接从任务条件数据中查询
	total, lists, err := planDao.NewPlanTaskConditionDataDao().GetSearchExternalListsByUnionIds(c, request, planId, request.Search.UnionIds)
	if err != nil {
		return total, result, err
	}

	// 从ck获取外部用户信息
	var pageUnionIds []string
	for _, model := range lists {
		pageUnionIds = append(pageUnionIds, model.UnionId)
	}
	userList, err := externalSearchDao.NewExternalExtraCk().GetList(c, request.EnterpriseId, "user_unionid in ?", pageUnionIds)
	if err != nil {
		return 0, nil, err
	}

	var externalMap = make(map[string]externalSearchModel.ExternalExtraModel)
	for _, model := range userList {
		externalMap[model.UserUnionid] = model
	}

	cardDescriptionMap, _ := customerSvc.NewCustomer().GetCustomerCardMapForType(enterpriseConfigName)

	// 获取用户备注关系
	relationList, _ := externalUserDao.NewUserRelation().GetMany(c, "enterprise_id = ? and member_id = ? AND unionid in ?",
		request.EnterpriseId, request.MemberId, pageUnionIds)

	// 构建结果
	for _, value := range lists {
		info := value
		res := planReq.CustomerListsRes{
			TaskConditionDataModel: info,
			CustomerInfoRes: planReq.CustomerInfoRes{
				Avatar: "https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com/plan_task/default.png",
			},
		}

		if item, ok := externalMap[value.UnionId]; ok {
			res.CustomerInfoRes.Name = item.Nick
			res.CustomerInfoRes.Avatar = item.Avatar
			res.CustomerInfoRes.VipType = item.VipType
			if item.TaskFinishAt > 0 {
				if nowUnix-int64(item.TaskFinishAt) < 86400*customerSearchReq.TaskHaveValidDays {
					res.CustomerInfoRes.TaskFinishAt = item.TaskFinishAt
				}
			}
			if item.ValidReserveAt > 0 {
				if nowUnix-int64(item.ValidReserveAt) < 86400*customerSearchReq.ValidDays {
					res.CustomerInfoRes.ValidReserveAt = item.ValidReserveAt
				}
			}
			if item.SuccessReserveAt > 0 {
				if nowUnix-int64(item.SuccessReserveAt) < 86400*customerSearchReq.ValidDays {
					res.CustomerInfoRes.SuccessReserveAt = item.SuccessReserveAt
				}
			}
			res.CustomerInfoRes.CouponStatus = item.CouponStatus
			res.CustomerInfoRes.CouponBy30daysStatus = item.CouponBy30daysStatus
			if cardName, o := cardDescriptionMap[item.CardLevel]; o {
				res.CustomerInfoRes.CustomerCard = cardName
			}
		}
		res.IsUnion = 0

		// 处理备注信息
		for _, relation := range relationList {
			if relation.Unionid == value.UnionId {
				res.CustomerInfoRes.RelRemark = relation.Remark
				res.IsUnion = helper.RuleThree(relation.RelationStatus > 0, 0, 1)
			}
		}

		result = append(result, res)
	}

	return total, result, nil
}

// SearchExternalAllOpenId   真全选，获取任务列表所有openid
func (p PlanTask) SearchExternalAllOpenId(c ctx.Context, request planReq.ExternalSearchListsReq) (openIds []string, err error) {
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id": request.Id,
	})
	if err != nil {
		return openIds, err
	}
	list, err := planDao.NewPlanTaskConditionDataDao().GetSearchExternalLists(c, request, one.PlanId)
	if err != nil || len(list) == 0 {
		return openIds, err
	}
	for _, model := range list {
		request.Search.UnionIds = append(request.Search.UnionIds, model.UnionId)
	}
	if request.Search.IsSearch == 0 {
		//如果没筛选直接mysql查询返回即可
		return externalUserDao.NewWxUserCorpExternal().GetOpenIdsByUnionIds(c, request.EnterpriseId, request.Search.UnionIds)
	}
	request.Search.EnterpriseID = request.EnterpriseId
	request.Search.MemberId = request.MemberId
	pluckIds, err := externalSearchDao.NewExternalExtraCk().ExternalSearchPluckIds(c, request.Search, "openid")
	if err != nil {
		return nil, err
	}
	for _, v := range pluckIds {
		// 使用类型断言将 interface{} 转换为 int64
		if val, ok := v.(string); ok {
			openIds = append(openIds, val)
		}
	}
	return openIds, nil
}

// SearchCustomerLists coach的任务计划升级，获取任务客户信息
func (p PlanTask) SearchCustomerLists(c ctx.Context, request planReq.CustomerSearchListsReq, enterpriseConfigName string) (int64, []planReq.CustomerListsRes, error) {
	result := make([]planReq.CustomerListsRes, 0)
	var total int64
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id": request.Id,
	})
	if err != nil {
		return total, result, err
	}
	var relationUnionIds []string
	relationUnionIds, err = planDao.NewUserRelation().GetManyUnionIds(c, map[string]interface{}{
		"enterprise_id":   request.EnterpriseId,
		"relation_status": customer.WorkUserRelationStatusNormal,
		"user_id":         one.UserId,
		"deleted_at":      0,
	})
	if err != nil {
		return total, result, err
	}
	list, err := planDao.NewPlanTaskConditionDataDao().GetSearchCustomerLists(c, request, one.PlanId, request.Search.IsUnion, relationUnionIds)
	if err != nil {
		return total, result, err
	}
	var unionIds, searchUnionIds []string
	for _, v := range list {
		unionIds = append(unionIds, v.UnionId)
		searchUnionIds = append(searchUnionIds, v.UnionId)
	}
	if request.Search.IsSearch == 1 && len(searchUnionIds) > 0 { // 是否筛选
		unionIds, err = externalUserDao.NewUserRelation().GetPlanUnionIdsBySearchListsByCk(c, request.EnterpriseId, request.Search, one.UserId, searchUnionIds)
		if err != nil {
			return total, result, err
		}
	}
	if len(unionIds) == 0 {
		return 0, result, err
	}
	total, lists, err := planDao.NewPlanTaskConditionDataDao().GetSearchCustomerListsByUnionIds(c, request, one.PlanId, unionIds)
	if err != nil {
		return total, result, err
	}
	resUnionIds := make([]string, 0)
	for _, value := range lists {
		info := value
		res := planReq.CustomerListsRes{
			TaskConditionDataModel: info,
			CustomerInfoRes: planReq.CustomerInfoRes{
				Name:   info.UserName,
				Avatar: "https://cdn-cos-uat.coach.com.cn/plan_task/default.png",
			},
		}
		resUnionIds = append(resUnionIds, info.UnionId)
		result = append(result, res)
	}
	if len(resUnionIds) > 0 {
		external, err := planDao.NewWxUserCorpExternal().GetMany(c, request.EnterpriseId, map[string]interface{}{
			"unionid": resUnionIds,
		})
		if err != nil {
			return 0, result, err
		}
		cardDescriptionMap, _ := customerSvc.NewCustomer().GetCustomerCardMapForType(enterpriseConfigName)
		extraData, err := planDao.NewUserRelation().GetExtraData(c, request.EnterpriseId, resUnionIds)
		if err != nil {
			return 0, result, err
		}
		validUnix := int(time.Now().Unix() - planDao.ValidDays*86400)
		unionIdInfo := make(map[string]planReq.CustomerInfoRes, 0)
		for _, model := range external {
			item := model
			customerInfo := planReq.CustomerInfoRes{
				Name:         item.Name,
				Avatar:       item.Avatar,
				CustomerCard: cardDescriptionMap[model.CardLevel],
				VipType:      model.VipType,
			}
			for _, extra := range extraData {
				if model.Unionid == extra.UserUnionid {
					if extra.SuccessReserveAt > validUnix {
						customerInfo.SuccessReserveAt = extra.SuccessReserveAt
					}
					if extra.ValidReserveAt > validUnix {
						customerInfo.ValidReserveAt = extra.ValidReserveAt
					}
					if extra.TaskFinishAt > validUnix {
						customerInfo.TaskFinishAt = extra.TaskFinishAt
					}
					customerInfo.CouponStatus = extra.CouponStatus
					customerInfo.CouponBy30daysStatus = extra.CouponBy30daysStatus
				}
			}
			unionIdInfo[item.Unionid] = customerInfo
		}
		qyTagIdMap := make(map[string]string)
		unionIdTags := make(map[string][]string)
		unionIdTagIds := make(map[string][]string)
		if request.IsTag == 1 { // 查询的用户展示标签
			tagLists, _ := workTagSvc.NewWorkTags().GetChoiceTagList(c, request.EnterpriseId)
			for _, model := range tagLists {
				qyTagIdMap[model.TagID] = model.Name
			}
			TagsList, err := planDao.NewPlanTaskConditionDataDao().GetConditionTagsByUnionIds(c, request.EnterpriseId, one.PlanId, request.Id, resUnionIds)
			if err != nil {
				return 0, nil, err
			}
			for _, v := range TagsList {
				tag := v
				if tag.WwTagId != "" {
					if !helper.InArray[string](tag.WwTagId, unionIdTagIds[tag.UnionId]) { // 去重多个企微标签的展示
						tagName := qyTagIdMap[tag.WwTagId]
						unionIdTagIds[tag.UnionId] = append(unionIdTagIds[tag.UnionId], tag.WwTagId)
						unionIdTags[tag.UnionId] = append(unionIdTags[tag.UnionId], tagName)
					}
				}
			}
		}
		for i, res := range result {
			v := res

			if v.UnionId != "" {
				if item, ok := unionIdInfo[v.UnionId]; ok {
					v.CustomerInfoRes = item
				}
				if !helper.InArray(v.UnionId, relationUnionIds) {
					v.IsUnion = 0
				} else {
					v.IsUnion = 1
				}
			}
			if request.IsTag == 1 {
				if _, ok := unionIdTags[v.UnionId]; ok {
					v.TagName = unionIdTags[v.UnionId]
				}
			} else {
				if v.TaskConditionDataModel.WwTagId != "" {
					if value, ok := qyTagIdMap[v.TaskConditionDataModel.WwTagId]; ok {
						v.TagName = append(v.TagName, value)
					}
				}
			}
			result[i] = v
		}
	}

	return total, result, err
}

func (p PlanTask) CustomerMove(c ctx.Context, req planReq.CustomerMoveReq) error {
	if len(req.OpenIds) > 0 {
		unionIds, err := externalUserDao.NewWxUserCorpExternal().GetUnionIdsByOpenIds(c, req.EnterpriseId, req.OpenIds)
		if err != nil {
			return err
		}
		req.UnionIds = unionIds
	}
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id": req.Id,
	})
	if err != nil {
		return err
	}
	return planDao.NewPlanTaskConditionDataDao().CustomerMove(c, one.PlanId, req)
}

// SendFinish coach的企业沟通任务转发上报
func (p PlanTask) SendFinish(c ctx.Context, req planReq.SendFinishReq) error {
	if len(req.OpenIds) > 0 {
		unionIds, err := externalUserDao.NewWxUserCorpExternal().GetUnionIdsByOpenIds(c, req.EnterpriseId, req.OpenIds)
		if err != nil {
			return err
		}
		req.UnionIds = unionIds
	}
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id": req.TaskAllotLogId,
	})
	if err != nil {
		return err
	}
	taskInfo, err := p.planTaskDao.GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id": one.TaskId,
	})
	if err != nil {
		return err
	}
	if taskInfo.Status == microScrmModel.PlanEnd { // 子任务已过期
		return errors.New("子任务已结束")
	}
	nowTime := time.Now().Unix()
	// 更新顾客标记已联系
	err = planDao.NewPlanTaskConditionDataDao().UpdateData(c, req.EnterpriseId, one.PlanId, map[string]interface{}{
		"task_allot_log_id": req.TaskAllotLogId,
		"union_id":          req.UnionIds,
		"is_finished":       microScrmModel.FinishedNo,
	}, map[string]interface{}{
		"updated_at":  nowTime,
		"is_finished": microScrmModel.FinishedYes,
		"forward":     microScrmModel.AllForward, // 转发给顾客
		"finished_at": nowTime,
	})
	// 统计已触达的好友数
	finishCount, err := planDao.NewPlanTaskConditionDataDao().GetGroupCount(c, req.EnterpriseId, one.PlanId, map[string]interface{}{
		"task_allot_log_id": req.TaskAllotLogId,
		"is_finished":       microScrmModel.FinishedYes,
	})
	if err != nil {
		return err
	}
	finishNum := decimal.NewFromInt(finishCount)
	customerNum := decimal.NewFromInt(int64(one.CustomerNum))
	divNum := finishNum.Div(customerNum).Round(2).Mul(decimal.NewFromInt(100)).IntPart()
	err = planDao.NewPlanTaskAllotLogDao().UpdateData(c, req.EnterpriseId, map[string]interface{}{
		"id": one.ID,
	}, map[string]interface{}{"finish_rate": divNum, "contacted_count": finishCount})
	if int(divNum) >= taskInfo.FinishRate && one.IsFinished == microScrmModel.PLAN_FISISHING && taskInfo.GroupSendType != microScrmModel.PLANTASKQYSEND {
		// 任务完成增加积分
		excitationSvc.NewPointsProduceService().TaskFinished(ctx.Wrap(c), excitationReq.TaskFinishedReq{
			TaskAllotLogId: one.ID,
			EnterpriseId:   req.EnterpriseId,
			MemberId:       req.MemberId,
		})
		err := planDao.NewPlanTaskAllotLogDao().UpdateData(c, req.EnterpriseId, map[string]interface{}{
			"id": one.ID,
		}, map[string]interface{}{"is_finished": microScrmModel.PLAN_FINISHED, "end_time": nowTime})
		if err != nil {
			return err
		}
	}
	many, err := planDao.NewPlanTaskConditionDataDao().GetMany(c, req.EnterpriseId, one.PlanId, map[string]interface{}{
		"task_allot_log_id": req.TaskAllotLogId,
		"is_finished":       microScrmModel.FinishedYes,
		"updated_at":        nowTime,
	})
	if err != nil {
		return err
	}
	for _, v := range many { // 触达任务增加积分
		md5Str := helper.Md5V(v.UnionId + strconv.FormatInt(one.ID, 10) + v.CustomerID)
		excitationSvc.NewPointsProduceService().ProduceQueue(ctx.Wrap(c), excitationReq.PointsProduceSave{
			MemberId:     req.MemberId,
			Type:         excitationModels.RuleCategoryId9,
			Value:        1,
			UniqueHash:   md5Str + fmt.Sprintf("-%d-%d", req.MemberId, excitationModels.RuleCategoryId9),
			EnterpriseId: req.EnterpriseId,
		})
	}
	go func() {
		if len(req.UnionIds) == 0 {
			return
		}
		planInfo, err := planDao.NewPlanTaskDao().GetPlanDetail(c, planReq.PlanInfo{
			BaseReq: planReq.BaseReq{
				EnterpriseId: req.EnterpriseId,
			},
			Id: taskInfo.PlanID,
		})
		if !helper.InArray(planInfo.Category, []int{microScrmModel.PlanMaTypeYes, microScrmModel.PlanMaTypeNewYes}) { // 非MA任务不推送
			return
		}
		if err != nil {
			log.Info(c).Msgf("标记联系MA回执查询主任务信息报错 %+v", err)
			return
		}
		var templateName string
		if taskInfo.PlanTemplateId > 0 {
			templateInfo, err := planDao.NewPlanTemplateDao().GetTemplateInfo(c, req.EnterpriseId, taskInfo.PlanTemplateId)
			if err != nil {
				return
			}
			templateName = templateInfo.TemplateName
		}
		memberInfo, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{
			"enterprise_id": req.EnterpriseId,
			"userid":        one.UserId,
			"deleted_at":    0,
		})
		if err != nil {
			log.Info(c).Msgf("查询导购信息身份失败 %+v", err)
			return
		}
		originType := taskInfo.OriginType
		date := time.Now().UTC().Format("2006-01-02T15:04:05Z")
		if len(req.UnionIds) > 100 {
			chunkData := lo.Chunk(req.UnionIds, 100)
			for _, list := range chunkData {
				var events []microwork.MaEvents
				for _, v := range list {
					event := microwork.MaEvents{
						IdentityType:      "wechat-unionid",
						IdentityValue:     v,
						Date:              date,
						CResult:           "已完成",
						CSaName:           memberInfo.Name,
						CSaCode:           one.UserId,
						CTaskId:           strconv.FormatInt(one.TaskId, 10),
						CTaskName:         taskInfo.Name,
						CCrowdId:          taskInfo.CrowdID,
						CTaskRemark:       taskInfo.Remark,
						FlowId:            taskInfo.FlowId,
						Campaign:          taskInfo.Campaign,
						CPlanTemplateId:   strconv.FormatInt(taskInfo.PlanTemplateId, 10),
						CPlanTemplateName: templateName,
						CTaskClassifation: planInfo.Name,
					}
					events = append(events, event)
				}
				params := microwork.SendMaFinishReq{
					EnterpriseID: req.EnterpriseId,
					OriginType:   originType,
					Events:       events,
					CommonProps: microwork.MACommonProps{
						Event: "c_sa_completed_task_result",
					},
				}
				err = macallback.SendFinish(c, params)
				if err != nil {
					log.Info(c).Err(err).Interface("params", params).Msg("标记联系请求MA的任务发送结果回执")
				}
			}
		} else {
			var events []microwork.MaEvents
			for _, v := range req.UnionIds {
				event := microwork.MaEvents{
					IdentityType:      "wechat-unionid",
					IdentityValue:     v,
					Date:              date,
					CResult:           "已完成",
					CSaName:           memberInfo.Name,
					CSaCode:           one.UserId,
					CTaskId:           strconv.FormatInt(one.TaskId, 10),
					CTaskName:         taskInfo.Name,
					CCrowdId:          taskInfo.CrowdID,
					CTaskRemark:       taskInfo.Remark,
					FlowId:            taskInfo.FlowId,
					Campaign:          taskInfo.Campaign,
					CPlanTemplateId:   strconv.FormatInt(taskInfo.PlanTemplateId, 10),
					CPlanTemplateName: templateName,
					CTaskClassifation: planInfo.Name,
				}
				events = append(events, event)
			}

			params := microwork.SendMaFinishReq{
				EnterpriseID: req.EnterpriseId,
				OriginType:   originType,
				Events:       events,
				CommonProps: microwork.MACommonProps{
					Event: "c_sa_completed_task_result",
				},
			}
			err = macallback.SendFinish(c, params)
			if err != nil {
				log.Info(c).Err(err).Interface("params", params).Msg("标记联系请求MA的任务发送结果回执")
			}
		}

	}()
	return err
}

// ShowTips 判断主任务下有非1v1沟通的子任务
func (p PlanTask) ShowTips(c ctx.Context, req planReq.NewStatReq) (bool, error) {
	var isShow bool
	one, err := p.planTaskDao.GetOne(c, req.EnterpriseId, "plan_id=? and task_category=? and group_send_type<? and deleted_at=0", req.PlanId, microScrmModel.PLANTASKCONTACT, microScrmModel.PLANTASK1V1SEND)
	if err != nil {
		return isShow, err
	}
	if one.ID > 0 {
		isShow = true
	}
	return isShow, err
}

func (p PlanTask) TaskPlanNameList(c ctx.Context, req planReq.NewStatReq) ([]string, error) {
	return p.planTaskDao.GetPlanNameList(c, req.EnterpriseId, req.PlanId)
}

// GetCustomerTask 获取顾客详情页中，导购和顾客的任务列表
func (p PlanTask) GetCustomerTask(c ctx.Context, req planReq.CustomerTaskReq) ([]planResp.CustomerTaskResp, int64, error) {
	result := make([]planResp.CustomerTaskResp, 0)
	var total int64
	// 先获取task_condition_data表中的task_allot_log
	taskAllotLogIds, err := planDao.NewPlanTaskConditionDataCkDao().GetTaskAllotLogIds(c, req.EnterpriseId, map[string]interface{}{
		"user_id":  req.UserId,
		"union_id": req.UnionId,
	})
	if err != nil {
		return result, total, err
	}
	result, total, err = p.planTaskDao.GetTaskList(c, req.EnterpriseId, taskAllotLogIds, req)
	if err != nil {
		return result, total, err
	}
	return result, total, err
}

// GetMaMomentTaskByDay   获取MA朋友圈任务列表
func (p PlanTask) GetMaMomentTaskByDay(c ctx.Context, enterpriseId, days int64) ([]microScrmModel.PlanTasksModel, error) {
	unix := time.Now().Unix() - days*86400
	return planDao.NewPlanTaskDao().GetMany(c, enterpriseId, "task_is_push=1 AND task_start_time>=? AND task_category=? AND origin_type in ?",
		unix, microScrmModel.PLANTASKMOMENT, []int{microScrmModel.START_DT_TASK_ORIGIN_TYPE, microScrmModel.SHENYAN_TASK_ORIGIN_TYPE})
}

// QyMomentComments  获取朋友圈评论点赞
func (p PlanTask) QyMomentComments(ctx ctx.Context, enterpriseId int64, taskDetail microScrmModel.PlanTasksModel, qyWorkApi *work.QyWorkImpl) error {
	userList, err := planDao.NewPlanTaskAllotLogDao().GetMany(ctx, enterpriseId, "task_id=? AND is_finished = ?", taskDetail.ID, microScrmModel.PLAN_FINISHED)
	if err != nil || len(userList) == 0 {
		return err
	}
	for _, v := range userList {
		// 查询朋友圈任务的可见范围列表客户，更新顾客状态为已联系
		sendRequest := work.GetMomentCommentsRequest{}
		sendRequest.UserId = v.UserId
		sendRequest.MomentId = taskDetail.MomentJobId
		result, err := qyWorkApi.GetMomentComments(ctx, sendRequest)
		if err != nil {
			app.Logger().Error(ctx).Msgf("查询朋友圈任务的可见范围列表客户报错, err:%s", err.Error())
			continue
		}
		if result.ErrCode != 0 {
			app.Logger().Error(ctx).Msgf("查询朋友圈任务的可见范围列表客户企微返回错误, response:%v", result)
			continue
		}
		var externalUserIds []string
		for _, value := range result.CommentList {
			externalUserIds = append(externalUserIds, value.ExternalUserid)
		}
		for _, value := range result.LikeList {
			externalUserIds = append(externalUserIds, value.ExternalUserid)
		}
		if len(externalUserIds) == 0 {
			app.Logger().Error(ctx).Msgf("没有可推送的互动数据, response:%v", result)
			continue
		}
		externalUserIds = helper.ArrayUniq(externalUserIds)
		unionIdMap, err := planDao.NewWxUserCorpExternal().GetUnionIdMapByOpenIds(ctx, enterpriseId, externalUserIds)
		if len(unionIdMap) == 0 {
			app.Logger().Error(ctx).Msgf("没有可推送的数据, response:%v", result)
			continue
		}
		// 发送朋友圈任务的可见范围列表客户的回执给到MA
		p.SendMAMomentCommentsData(ctx, enterpriseId, taskDetail, v.UserId, unionIdMap, *result)
	}
	return nil

}

// SendMAMomentCommentsData   朋友圈任务评论回执给到MA
func (p PlanTask) SendMAMomentCommentsData(c ctx.Context, enterpriseId int64, taskDetail microScrmModel.PlanTasksModel, saCode string, unionIdMap map[string]string, resp work.GetMomentCommentResp) {
	redisKey := fmt.Sprintf("matask:moment-comments-send-time:%d-%d-%s", enterpriseId, taskDetail.ID, saCode)
	//获取上次推送时间，获取到的互动数据如果早于上次推送时间，则不推送
	lastSendTimeStr := app.Redis().Get(c, redisKey).Val()
	var maxCreateTime int64
	var lastSendUnix int64
	if lastSendTimeStr != "" {
		lastSendUnix, _ = strconv.ParseInt(lastSendTimeStr, 10, 64)
	}
	sendMaMomentCommentReq := microwork.SendMaMomentCommentsReq{
		EnterpriseID: enterpriseId,
		SaCode:       saCode,
		TaskId:       taskDetail.ID,
		CrowdId:      taskDetail.CrowdID,
		Campaign:     taskDetail.Campaign,
		CommentList:  make([]microwork.SendMaMomentCommentsInfo, 0),
		LikeList:     make([]microwork.SendMaMomentCommentsInfo, 0),
	}
	if len(unionIdMap) > 0 {
		for _, v := range resp.CommentList {
			if unionIdMap[v.ExternalUserid] == "" {
				continue
			}
			if v.CreateTime <= lastSendUnix {
				continue
			}
			if v.CreateTime > maxCreateTime {
				maxCreateTime = v.CreateTime
			}
			sendMaMomentCommentReq.CommentList = append(sendMaMomentCommentReq.CommentList, microwork.SendMaMomentCommentsInfo{
				UnionId:   unionIdMap[v.ExternalUserid],
				CreateTme: v.CreateTime,
			})
		}
		for _, v := range resp.LikeList {
			if unionIdMap[v.ExternalUserid] == "" {
				continue
			}
			if v.CreateTime <= lastSendUnix {
				continue
			}
			if v.CreateTime > maxCreateTime {
				maxCreateTime = v.CreateTime
			}
			sendMaMomentCommentReq.LikeList = append(sendMaMomentCommentReq.LikeList, microwork.SendMaMomentCommentsInfo{
				UnionId:   unionIdMap[v.ExternalUserid],
				CreateTme: v.CreateTime,
			})
		}
	}
	if len(sendMaMomentCommentReq.LikeList) == 0 && len(sendMaMomentCommentReq.CommentList) == 0 {
		return
	}
	err := macallback.SendMomentComments(c, taskDetail.OriginType, sendMaMomentCommentReq)
	//推送完记录上次推送时间
	app.Redis().Set(c, redisKey, strconv.FormatInt(maxCreateTime, 10), 7*24*time.Hour)

	log.Info(c).Err(err).Interface("params", sendMaMomentCommentReq).Msg("朋友圈任务评论回执给到MA")
}
