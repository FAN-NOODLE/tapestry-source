package plan_task

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	membersDao "micro-scrm/internal/app/dao/members"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	planReq "micro-scrm/internal/app/http/request/plan"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/modules/plan_task/delivery/http_entity"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/orm/model/micro_scrm/group_chat"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/dl"
	"micro-scrm/internal/pkg/helper"
	"sort"
)

func CustomerDataService(cc ctx.Context, enterpriseID int64, req CustomerDataReq) (CustomerDataResp, error) {
	var resp CustomerDataResp

	// 任务总件数
	taskCount, TaskErr := orm.Service.PlanTask.TaskCount(enterpriseID, req.TaskID)
	if TaskErr != nil {
		return resp, TaskErr
	}
	resp.TaskCount = taskCount

	// 任务完成数
	taskCompleteCount, TaskCompleteErr := orm.Service.PlanTask.TaskComplateCount(enterpriseID,
		req.TaskID, req.StartTime, req.EndTime)
	if TaskCompleteErr != nil {
		return resp, TaskCompleteErr
	}
	resp.TaskCompleteCount = taskCompleteCount

	// 任务完成率
	var rate string
	if taskCount == 0 {
		rate = "0%"
	} else {
		rate = mathrounded(float64(taskCompleteCount), float64(taskCount))
	}
	resp.TaskCompleteRate = rate

	// 触达ba数
	baCount, baErr := orm.Service.PlanTask.TouchBaCount(enterpriseID, req.TaskID)
	if baErr != nil {
		return resp, baErr
	}
	resp.TouchBaCount = baCount

	ckDao := planDao.NewPlanTaskConditionDataCkDao()
	// 匹配的用户数
	customerCount, customerErr := ckDao.CustomerCount(cc, enterpriseID, req.TaskID, req.PlanID)
	if customerErr != nil {
		return resp, customerErr
	}
	resp.CustomerCount = customerCount
	// 匹配用户数去重
	DisCount, DisErr := ckDao.CustomerDistinctCount(cc, enterpriseID, req.PlanID, map[string]interface{}{
		"task_id":    req.TaskID,
		"deleted_at": 0,
	})

	if DisErr != nil {
		return resp, DisErr
	}
	resp.CustomerDistictCount = DisCount
	// 累计触达客户数去重，沟通任务统计is_finished,朋友圈任务统计is_touch
	var totalReachCustomerNumber int64
	taskInfo, err := planDao.NewPlanTaskDao().GetOne(cc, enterpriseID, map[string]interface{}{"id": req.TaskID})
	if err != nil {
		return resp, err
	}
	whereMap := map[string]interface{}{
		"task_id": taskInfo.ID,
	}
	if taskInfo.TaskCategory == microScrmModel.PLANTASKCONTACT {
		whereMap["is_finished"] = microScrmModel.FinishedYes
	} else {
		whereMap["is_touch"] = microScrmModel.TouchStatusYes
	}
	totalReachCustomerNumber, err = ckDao.GetCount(cc, enterpriseID, taskInfo.PlanID, whereMap)
	if err != nil {
		return resp, err
	}
	resp.TotalReachCustomerNumber = totalReachCustomerNumber
	// 查询统计记录
	// planStatistic, _ := pService.GetTaskStatistics(cc, enterpriseID, req.TaskID, req.EndTime-24*60*60)
	// if planStatistic != nil {
	// 	resp.TotalReachCustomerNumber = planStatistic.TotalReachCustomerNumber
	// 	resp.TotalOrderCustomerNumber = planStatistic.TotalOrderCustomerNumber
	// 	resp.TotalOrderNumber = planStatistic.TotalOrderNumber
	// 	resp.TotalOrderFee = planStatistic.TotalOrderFee
	// }

	return resp, nil
}

func ExportCustomerTaskData(cc ctx.Context, enterpriseID int64, params http_entity.ExportTasksRequst, pService *PTaskService) (newResult []micro_scrm_enterprise.ExportTaskDataResp, err error) {
	if params.CompletionStartTime == params.CompletionEndTime && params.CompletionStartTime != 0 {
		params.CompletionEndTime = params.CompletionStartTime + (60*60*24 - 1)
	}
	if params.PlanStartTime == params.PlanEndTime && params.PlanStartTime != 0 {
		params.PlanEndTime = params.PlanStartTime + (60*60*24 - 1)
	}
	// 根据部门ID，获取所有下级部门id集合
	// DepartmentIds := orm.Service.Departments.GetDepartmentsChildrenIds(enterpriseID, params.DepartmentId)
	var searchUserIds []string
	if params.BaName != "" {
		searchUserIds, _ = membersDao.NewMembers().GetUserIdsByName(cc, enterpriseID, params.BaName)
	}
	if params.DepartmentId != 0 {
		memberLists, err := membersDao.NewMembers().GetMemberByDepartmentId(cc, enterpriseID, []int64{params.DepartmentId})
		if err != nil {
			return nil, err
		}
		for _, v := range memberLists {
			info := v
			if !helper.InArray[string](info.MembersUserid, searchUserIds) {
				searchUserIds = append(searchUserIds, info.MembersUserid)
			}
		}
	}
	info, err := planDao.NewPlanTaskDao().GetConditionInfo(cc, enterpriseID, map[string]interface{}{"id": params.Id})
	if err != nil {
		return nil, err
	}
	taskInfo, err := planDao.NewPlanTaskDao().GetOne(cc, enterpriseID, map[string]interface{}{"id": info.TaskID})
	if err != nil {
		return nil, err
	}
	planInfo, err := planDao.NewPlanTaskDao().GetPlanDetail(cc, planReq.PlanInfo{
		BaseReq: planReq.BaseReq{
			EnterpriseId: enterpriseID,
		},
		Id: info.PlanID,
	})
	if err != nil {
		return nil, err
	}

	result, _ := orm.Service.PlanTask.ExportPlanTasksCompletionV2(enterpriseID, params.Id, params.PlanStartTime, params.PlanEndTime, params.CompletionStartTime, params.CompletionEndTime, searchUserIds, params.IsFinished)
	userIds := make([]string, 0)
	for _, value := range result {
		userIds = append(userIds, value.UserId)
	}
	memberInfoLists, err := membersDao.NewMembers().GetUnscopedMemberInfoAndDepartmentsByUserId(cc, enterpriseID, userIds)
	if err != nil {
		return nil, err
	}
	sort.Slice(memberInfoLists, func(i, j int) bool {
		return memberInfoLists[i].DepartmentId > memberInfoLists[j].DepartmentId
	})
	memberInfoMap := make(map[string]membersModel.MemberDepartment)
	for _, v := range memberInfoLists {
		if _, ok := memberInfoMap[v.MembersUserid]; !ok {
			memberInfoMap[v.MembersUserid] = v
		}

	}
	// count = total
	chatNameMap := make(map[int64]group_chat.WorkGroupChatModel, 0)
	for _, v := range orm.Service.PlanTask.GetGroupChatName(enterpriseID) {
		chatNameMap[int64(v.ID)] = v
	}
	if len(userIds) > 0 {
		for _, v := range result {
			item := micro_scrm_enterprise.ExportTaskDataResp{
				DataBase: micro_scrm_enterprise.DataBase{
					UserId:        v.UserId,
					UserName:      memberInfoMap[v.UserId].MembersName,
					PlanName:      planInfo.Name,
					TaskName:      taskInfo.Name,
					IsFinished:    v.IsFinished,
					StartTime:     v.StartTime,
					CompletedTime: v.EndTime,
				},
				DepartmentName:  memberInfoMap[v.UserId].DepartmentName,
				CustomerNum:     v.CustomerNum,
				ID:              v.ID,
				WhetherToResign: 0,
				Category:        planInfo.Category,
			}
			v.TaskCategory = taskInfo.TaskCategory
			v.GroupSendType = taskInfo.GroupSendType
			// // 查询统计记录
			// planStatistic, _ := p.GetTaskStatistics(cc, enterpriseID, v.ID, time.Now().Unix())
			// if planStatistic != nil {
			// 	v.TotalReachCustomerNumber = planStatistic.TotalReachCustomerNumber
			// 	v.TotalOrderCustomerNumber = planStatistic.TotalOrderCustomerNumber
			// 	v.TotalOrderNumber = planStatistic.TotalOrderNumber
			// 	v.TotalOrderFee = planStatistic.TotalOrderFee
			// }
			if taskInfo.TaskCategory == microScrmModel.PLANTASKCONTACT {
				item.TotalReachCustomerNumber = v.ContactedCount
			}
			if taskInfo.TaskCategory == microScrmModel.PLANTASKMOMENT {
				item.TotalReachCustomerNumber = v.MomentLookCount
			}
			item.FinishRate = fmt.Sprintf("%s%%", v.FinishRate)
			newResult = append(newResult, item)
		}
	}
	return

	// if err != nil {
	//	return
	// }
	// userIds := make([]string, 0)
	// for _, value := range result {
	//	userIds = append(userIds, value.UserId)
	// }
	// BaInfoList, infoErr := orm.WorkService.Work.GetBAInfoListWithUserIdPart(enterpriseId, userIds)
	// if infoErr != nil {
	//	err = infoErr
	//	return
	// }
	// //优化逻辑
	// var baInfoMap = make(map[string]string, 0)
	// for _, baInfo := range BaInfoList {
	//	baInfoMap[baInfo.MemUserid] = baInfo.WorkerNo
	// }
	// for index, res := range result {
	//	if val, find := baInfoMap[res.UserId]; find {
	//		result[index].WorkNo = val
	//	}
	//	// 查询统计记录
	//	planStatistic, _ := pService.GetTaskStatistics(cc, enterpriseId, res.ID, param.EndTime-24*60*60)
	//	if planStatistic != nil {
	//		result[index].TotalReachCustomerNumber = planStatistic.TotalReachCustomerNumber
	//		result[index].TotalOrderCustomerNumber = planStatistic.TotalOrderCustomerNumber
	//		result[index].TotalOrderNumber = planStatistic.TotalOrderNumber
	//		result[index].TotalOrderFee = planStatistic.TotalOrderFee
	//	}
	// }
	// return

	// bossRes, errs := orm.BossService.Enterprise.GetInfo(boss.EnterpriseModel{
	//	ID: enterpriseId,
	// })
	// if errs != nil {
	//	err = errs
	//	return
	// }
	// d, err := dl.NewDL().FindBaListWithUserId(ctx, dl.BAListRequest{
	//	EnterpriseID: enterpriseId,
	//	RankLength:   0,
	//	QyDomain:     bossRes.Domain,
	//	AdminId:      "0",
	//	Type:         1,
	//	StartTime:    1537790425,
	//	EndTime:      1853409625,
	//	RankColumn:   "cumulative_customer_num",
	// })
	// baDeptes := []dl.BAinfoExt{}
	// if err == nil {
	//	baDeptes = d.List
	// }
	// for index, value := range result {
	//	for _, baDept := range baDeptes {
	//		if value.WorkNo == baDept.WorkerNo {
	//			areaLevelListLen := len(baDept.AreaLevelList)
	//			if areaLevelListLen > 1 {
	//				area := ""
	//				for i := areaLevelListLen - 2; i >= 0; i-- {
	//					if i == areaLevelListLen-2 {
	//						area = baDept.AreaLevelList[i]
	//					} else {
	//						area += "/" + baDept.AreaLevelList[i]
	//					}
	//				}
	//				result[index].BelongingArea = area
	//			} else if areaLevelListLen == 1 {
	//				result[index].BelongingArea = baDept.AreaLevelList[0]
	//			}
	//			break
	//		} else {
	//			result[index].BelongingArea = "-"
	//		}
	//	}
	//
	// }
	// return
}

// 导出社群任务列表数据
func ExportGroupTaskData(ctx ctx.Context, enterpriseId int64, param http_entity.ExportTasksRequst, p *PTaskService, cc ctx.Context) (
	[]micro_scrm_enterprise.ExportGroupChatDataResp, error) {
	result := orm.Service.PlanTask.ExportGroupTasksListData(enterpriseId, param)
	userIds := []string{}
	for _, value := range result {
		userIds = append(userIds, value.UserId)
	}
	BaInfoList, infoErr := orm.WorkService.Work.GetBAInfoListWithUserIdPart(enterpriseId, userIds)
	if infoErr != nil {
		return result, infoErr
	}
	for index, res := range result {
		for _, ba := range BaInfoList {
			if res.UserId == ba.MemUserid {
				result[index].WorkNo = ba.WorkerNo
			}
		}
		// 查询统计记录
		planStatistic, _ := p.GetTaskStatistics(cc, enterpriseId, res.ID, param.EndTime-24*60*60)
		if planStatistic != nil {
			result[index].TotalReachCustomerNumber = planStatistic.TotalReachCustomerNumber
			result[index].TotalOrderCustomerNumber = planStatistic.TotalOrderCustomerNumber
			result[index].TotalOrderNumber = planStatistic.TotalOrderNumber
			result[index].TotalOrderFee = planStatistic.TotalOrderFee
		}
	}
	bossRes, err := orm.BossService.Enterprise.GetInfo(boss.EnterpriseModel{
		ID: enterpriseId,
	})
	if err != nil {
		return result, err
	}
	d, err := dl.NewDL().FindBaListWithUserId(ctx, dl.BAListRequest{
		EnterpriseID: enterpriseId,
		RankLength:   0,
		QyDomain:     bossRes.Domain,
		AdminId:      "0",
		Type:         1,
		StartTime:    1537790425,
		EndTime:      1853409625,
		RankColumn:   "cumulative_customer_num",
	})
	baDeptes := []dl.BAinfoExt{}
	if err == nil {
		baDeptes = d.List
	}
	for index, value := range result {
		for _, baDept := range baDeptes {
			if value.WorkNo == baDept.WorkerNo {
				areaLevelListLen := len(baDept.AreaLevelList)
				if areaLevelListLen > 1 {
					area := ""
					for i := areaLevelListLen - 2; i >= 0; i-- {
						if i == areaLevelListLen-2 {
							area = baDept.AreaLevelList[i]
						} else {
							area += "/" + baDept.AreaLevelList[i]
						}
					}
					result[index].BelongingArea = area
				} else if areaLevelListLen == 1 {
					result[index].BelongingArea = baDept.AreaLevelList[0]
				}
				break
			} else {
				result[index].BelongingArea = "-"
			}
		}
	}
	return result, nil
}

// 社群任务列表数据
func GroupTaskListData(ctx ctx.Context, enterpriseId int64, param http_entity.TaskListRequst, p *PTaskService, cc ctx.Context) (
	[]micro_scrm_enterprise.ExportGroupChatDataResp, int64, error) {

	// 根据部门ID，获取所有下级部门id集合
	param.DepartmentIds = orm.Service.Departments.GetDepartmentsChildrenIds(enterpriseId, param.DepartmentId)

	result, total := orm.Service.PlanTask.GroupTasksListData(enterpriseId, param)
	userIds := []string{}
	for index, value := range result {
		userIds = append(userIds, value.UserId)
		if value.IsFinished != 1 {
			result[index].CompletedTime = 0
		}
	}
	BaInfoList, infoErr := orm.WorkService.Work.GetBAInfoListWithUserIdPart(
		enterpriseId, userIds)
	if infoErr != nil {
		return result, 0, infoErr
	}
	for index, res := range result {
		for _, ba := range BaInfoList {
			if res.UserId == ba.MemUserid {
				result[index].WorkNo = ba.WorkerNo
			}
		}
	}
	bossRes, err := orm.BossService.Enterprise.GetInfo(boss.EnterpriseModel{
		ID: enterpriseId,
	})
	if err != nil {
		return result, 0, err
	}
	d, err := dl.NewDL().FindBaListWithUserId(ctx, dl.BAListRequest{
		EnterpriseID: enterpriseId,
		RankLength:   0,
		QyDomain:     bossRes.Domain,
		AdminId:      "0",
		Type:         1,
		StartTime:    1537790425,
		EndTime:      1853409625,
		RankColumn:   "cumulative_customer_num",
	})
	baDeptes := []dl.BAinfoExt{}
	if err == nil {
		baDeptes = d.List
	}
	for index, value := range result {
		for _, baDept := range baDeptes {
			if value.WorkNo == baDept.WorkerNo {
				areaLevelListLen := len(baDept.AreaLevelList)
				if areaLevelListLen > 1 {
					area := ""
					for i := areaLevelListLen - 2; i >= 0; i-- {
						if i == areaLevelListLen-2 {
							area = baDept.AreaLevelList[i]
						} else {
							area += "/" + baDept.AreaLevelList[i]
						}
					}
					result[index].BelongingArea = area
				} else if areaLevelListLen == 1 {
					result[index].BelongingArea = baDept.AreaLevelList[0]
				}
				break
			} else {
				result[index].BelongingArea = "-"
			}
		}

		// 查询统计记录
		planStatistic, _ := p.GetTaskStatistics(cc, enterpriseId, value.ID, param.EndTime-24*60*60)
		if planStatistic != nil {
			result[index].TotalReachCustomerNumber = planStatistic.TotalReachCustomerNumber
			result[index].TotalOrderCustomerNumber = planStatistic.TotalOrderCustomerNumber
			result[index].TotalOrderNumber = planStatistic.TotalOrderNumber
			result[index].TotalOrderFee = planStatistic.TotalOrderFee
		}
	}
	return result, total, nil
}

func GroupTaskDataService(enterpriseID, planId, taskId, startTime,
	endTime int64, pService *PTaskService, cc ctx.Context) (CustomerDataResp, error) {
	var resp CustomerDataResp

	// 任务总件数
	taskCount, TaskErr := orm.Service.PlanTask.TaskCount(enterpriseID, taskId)
	if TaskErr != nil {
		return resp, TaskErr
	}
	resp.TaskCount = taskCount

	// 任务完成数
	taskCompleteCount, TaskCompleteErr := orm.Service.PlanTask.TaskComplateCount(
		enterpriseID, taskId, startTime, endTime)
	if TaskCompleteErr != nil {
		return resp, TaskCompleteErr
	}
	resp.TaskCompleteCount = taskCompleteCount

	// 任务完成率
	var rate string
	if taskCount == 0 {
		rate = "0%"
	} else {
		rate = mathrounded(float64(taskCompleteCount), float64(taskCount))
	}
	resp.TaskCompleteRate = rate

	// 触达群主
	masterCount, masterErr := orm.Service.PlanTask.TouchMasterCount(enterpriseID,
		taskId, startTime, endTime)
	if masterErr != nil {
		return resp, masterErr
	}
	resp.TouchMasterCount = masterCount

	// 匹配的社群数
	groupCount, groupErr := orm.Service.PlanTask.GroupChatCount(enterpriseID,
		planId, taskId, startTime, endTime)
	if groupErr != nil {
		return resp, groupErr
	}
	resp.GroupChatCount = groupCount

	// 匹配的社群数
	newCount, newErr := orm.Service.PlanTask.GroupNewTaskCount(enterpriseID,
		taskId, startTime, endTime)
	if newErr != nil {
		return resp, newErr
	}
	resp.NewTaskCount = newCount

	// 查询统计记录
	planStatistic, _ := pService.GetTaskStatistics(cc, enterpriseID, taskId, endTime-24*60*60)
	if planStatistic != nil {
		resp.TotalReachCustomerNumber = planStatistic.TotalReachCustomerNumber
		resp.TotalOrderCustomerNumber = planStatistic.TotalOrderCustomerNumber
		resp.TotalOrderNumber = planStatistic.TotalOrderNumber
		resp.TotalOrderFee = planStatistic.TotalOrderFee
	}

	return resp, nil
}

func mathrounded(params1, params2 float64) string {
	if params1 == 0 {
		return "0%"
	}
	params3 := helper.MathRounded((params1/params2)*100, 2)
	if params3 == float64(int64(params3)) {
		return fmt.Sprintf("%.0f", params3) + "%"
	}
	return fmt.Sprintf("%.2f", params3) + "%"
}
