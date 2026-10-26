package excitation

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	memberReq "micro-scrm/internal/app/http/request/member"
	reserveReq "micro-scrm/internal/app/http/request/reserve"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service/department"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
	"sort"
	"strconv"
	"time"
)

type RelationRule struct {
	name          string
	memberDataSvc *excitationSvc.MemberDataService
	taskRecordDao *excitationDao.TaskRecord
}

func NewRelationRule() *RelationRule {
	return &RelationRule{
		name:          "关联能力图谱规则处理",
		memberDataSvc: excitationSvc.NewMemberDataService(),
		taskRecordDao: excitationDao.NewTaskRecord(),
	}
}

func (job *RelationRule) Run(c ctx.Context, handler *xjob.JobHandler) error {

	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	var wg threading.RoutineGroup
	for _, id := range enterpriseIds {
		tmpId := id
		wg.Run(func() {
			err := job.Handle(c, tmpId)
			if err != nil {
				app.Logger().Error(c).Msgf("error : %v", err)
			}
		})
	}

	wg.Wait()

	return nil
}

func (job *RelationRule) Handle(c ctx.Context, enterpriseId int64) error {
	// 查询配置的企业
	resConfig, err := microBossDao.NewEnterpriseConfig().
		GetInfo(c, map[string]interface{}{"enterprise_id": enterpriseId, "key": "update_member_level"})
	if err != nil {
		return err
	}
	if resConfig.ID == 0 {
		return nil
	}

	paramsTime := time.Now()
	//paramsTime := time.Unix(1693677600, 0)

	nowUnix := paramsTime.Unix() // 自定义执行时间
	yesterdayUnix := nowUnix - 86400

	// 获取定时任务执行前一天的相关财月（也就是符合数据统计的财月）
	fiscalMonth, err := microWorkDao.NewFiscalMonth().GetDataByYearAndTime(c, enterpriseId, yesterdayUnix)
	if err != nil {
		return err
	}

	// 查询关联能力图谱规则
	incentiveRelation, err1 := excitationDao.NewIncentiveRelation().GetMany(c, "enterprise_id = ?", enterpriseId)

	if err1 != nil {
		return err1
	}

	if len(incentiveRelation) == 0 {
		return nil
	}

	for _, item := range incentiveRelation {
		// 任务是否执行
		// uniqueStr 拼接 日期-3-规则时间类型-规则类型
		uniqueStr := paramsTime.Format("2006-01-02") + fmt.Sprintf("-%d-%d", excitationModels.TimeRule3, item.IncentiveConfigType)
		exists, _ := job.taskRecordDao.GetOne(c, "enterprise_id = ? AND unique_str = ?", enterpriseId, uniqueStr)
		if exists.ID > 0 {
			continue
		}

		resArr := job.getRelationOneArr(c, enterpriseId, incentiveRelation, item.IncentiveConfigType)

		var exeSign = false
		switch item.TimeRule {
		case excitationModels.TimeRule3: // 每周

			if int(paramsTime.Weekday()) > 0 {
				// 财周任务应该是周天执行
				continue
			}

			// 财周起止时间
			endTime := helper.GetLastDateOfWeek(helper.GetZeroTimeByTime(paramsTime)).Unix()
			startTime := endTime - 86400*7

			// 获取各类数据
			resTopMap, err2 := job.getRelationDataArr(c, enterpriseId, startTime, endTime, item.IncentiveConfigType)

			if err2 != nil {
				return err2
			}

			errs := job.handleMemberDataLimit(c, enterpriseId, item, resArr, resTopMap)
			if errs != nil {
				return errs
			}
			exeSign = true
		case excitationModels.TimeRule4: // 每月

			if fiscalMonth.EndTime-yesterdayUnix > 86400 {
				// 财月任务应该是财月结束后一天内执行
				continue
			}

			// 获取各类数据
			resTopMap, err2 := job.getRelationDataArr(c, enterpriseId, fiscalMonth.StartTime, fiscalMonth.EndTime, item.IncentiveConfigType)

			if err2 != nil {
				return err2
			}

			errs := job.handleMemberDataLimit(c, enterpriseId, item, resArr, resTopMap)

			if errs != nil {
				return errs
			}
			exeSign = true

		}

		if exeSign {
			// 记录执行任务
			err = job.taskRecordDao.CreateData(c, excitationModels.ExcitationTaskRecordModel{
				MicroBaseModel: models.CreateMicroBase(enterpriseId, app.UID()),
				UniqueStr:      uniqueStr,
				Type:           excitationModels.TaskTypeIncentive,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (job *RelationRule) getRelationDataArr(c ctx.Context, enterpriseId, startTime, endTime int64, incentiveConfigType int) (map[int64]models.MemberTopSortList, error) {

	var res map[int64]models.MemberTopSortList
	if incentiveConfigType == excitationModels.IncentiveConfigTypeInvite {
		var inviteTopArr []models.MemberTopSortList
		var inviteTopMap = make(map[int64]models.MemberTopSortList)

		// 邀约目标达成率 = 邀约金额 / 个人销售金额 x 100%
		achievementParams := reserveReq.AchievementParams{
			EnterpriseId: enterpriseId,
			Type:         excitationModels.RuleCategoryId5,
			StartTime:    startTime,
			EndTime:      endTime,
		}
		// 获取所有导购邀约-邀约金额
		achievementMap, err := microScrmDao.NewReservesOrder().GetAllMemberAchievement(c, achievementParams)
		if err != nil {
			return inviteTopMap, err
		}
		// 获取所有导购-个人总销售金额
		volumeMap, err := microScrmDao.NewMemberSalesVolume().GetMemberSalesVolumeByReserve(c,
			memberReq.MemberSalesVolumeReq{
				EnterpriseId: enterpriseId,
				StartTime:    startTime,
				EndTime:      endTime,
			})
		if err != nil {
			return inviteTopMap, err
		}

		// 获取邀约金额在个人总销售金额的占比
		for memberId, tmpData := range achievementMap {
			var temp models.MemberTopSortList
			temp.MemberId = memberId
			temp.Value = tmpData
			temp.Ratio = 0.00
			if _, ok := volumeMap[memberId]; ok {
				if volumeMap[memberId] > 0 {
					temp.Ratio = (tmpData / 100) / volumeMap[memberId]
				}
			}
			inviteTopArr = append(inviteTopArr, temp)
		}

		// 按照邀约目标达成率降序排序
		sort.Slice(inviteTopArr, func(i, j int) bool {
			return inviteTopArr[i].Ratio > inviteTopArr[j].Ratio
		})
		for k, v := range inviteTopArr {
			var tempMap models.MemberTopSortList
			tempMap.MemberId = v.MemberId
			tempMap.Value = v.Value
			tempMap.Ratio = v.Ratio * 100 //邀约金额 / 个人销售金额 x 100%
			tempMap.SortTop = k + 1
			inviteTopMap[v.MemberId] = tempMap
		}

		return inviteTopMap, nil

	}

	// 任务 -- 百分比
	if incentiveConfigType == excitationModels.IncentiveConfigTypeTask {
		var planTaskTopArr []models.MemberTopSortList
		var planTaskMap = make(map[int64]models.MemberTopSortList)

		// 查询所有导购的-任务完成总数
		planTaskRate, err := microScrmDao.NewPlanTaskAllotLogDao().GetAllMemberAchievement(c, enterpriseId, excitationModels.RuleCategoryId8, startTime, endTime)
		if err != nil {
			return planTaskMap, err
		}
		var userIds []string
		for userId, ratio := range planTaskRate {
			userIds = append(userIds, userId)
			var tempArr models.MemberTopSortList
			tempArr.MemberUserid = userId
			tempArr.Ratio = ratio // 任务完成率 百分比值
			planTaskTopArr = append(planTaskTopArr, tempArr)

		}
		// 获取userId对应的memberId
		memberList, err := membersDao.NewMembers().GetManyByUserIds(c, enterpriseId, userIds)
		if err != nil {
			return planTaskMap, err
		}
		var memberMap = make(map[string]int64)
		for _, member := range memberList {
			memberMap[member.Userid] = member.ID
		}

		// 按照任务完成率降序排序
		sort.Slice(planTaskTopArr, func(i, j int) bool {
			return planTaskTopArr[i].Ratio > planTaskTopArr[j].Ratio
		})
		for k, v := range planTaskTopArr {
			var tempMap models.MemberTopSortList
			// 匹配导购id
			if _, ok := memberMap[v.MemberUserid]; ok {
				tempMap.MemberId = memberMap[v.MemberUserid]
			}

			tempMap.Ratio = v.Ratio
			tempMap.SortTop = k + 1
			planTaskMap[tempMap.MemberId] = tempMap
		}
		return planTaskMap, nil
	}

	// 客户互动 -- 排行
	if incentiveConfigType == excitationModels.IncentiveConfigTypeInteract {

		var interactTopArr []models.MemberTopSortList
		var interactTopMap = make(map[int64]models.MemberTopSortList)
		// 互动维度的积分类型
		var categoryIds = []int{
			excitationModels.RuleCategoryId10,
			excitationModels.RuleCategoryId11,
			excitationModels.RuleCategoryId12,
			excitationModels.RuleCategoryId13,
			excitationModels.RuleCategoryId14,
		}
		// 获取所有导购时间段内聊天互动排行榜
		interactMemberList, err := excitationDao.NewPointsRecord().GetPointsTopByCategoryIds(c, enterpriseId, categoryIds, startTime, endTime)
		if err != nil {
			return interactTopMap, err
		}

		for _, tmpValue := range interactMemberList {
			var temp models.MemberTopSortList
			temp.MemberId = tmpValue.MemberId
			temp.Total = tmpValue.Total
			interactTopArr = append(interactTopArr, temp)
		}

		// Top 排行map 按照客户互动数排序
		sort.Slice(interactTopArr, func(i, j int) bool {
			return interactTopArr[i].Total > interactTopArr[j].Total
		})
		for k, v := range interactTopArr {
			interactTopArr[k].SortTop = k + 1

			var tempMap models.MemberTopSortList
			tempMap.MemberId = v.MemberId
			tempMap.Total = v.Total
			tempMap.SortTop = k + 1
			interactTopMap[v.MemberId] = tempMap
		}

		return interactTopMap, nil

	}

	// 培训互动积分 -- 排行
	if incentiveConfigType == excitationModels.IncentiveConfigTypeTrain {

		var trainTopArr []models.MemberTopSortList
		var trainTopMap = make(map[int64]models.MemberTopSortList)

		// 查询导购培训互动积分排行
		memberPointsList, err := excitationDao.NewPointsRecord().GetTopList(c, enterpriseId, excitationModels.RuleCategoryId15, startTime, endTime)
		if err != nil {
			return trainTopMap, err
		}

		// 排行
		for memberId, tmpValue := range memberPointsList {
			var temp models.MemberTopSortList
			temp.MemberId = memberId
			temp.Value = tmpValue // 积分
			trainTopArr = append(trainTopArr, temp)
		}

		// 按照积分降序排序
		sort.Slice(trainTopArr, func(i, j int) bool {
			return trainTopArr[i].Value > trainTopArr[j].Value
		})
		for k, v := range trainTopArr {

			var tempMap models.MemberTopSortList
			tempMap.MemberId = v.MemberId
			tempMap.Value = v.Value
			tempMap.SortTop = k + 1
			trainTopMap[v.MemberId] = tempMap
		}

		return trainTopMap, nil
	}

	// 线上能力 -- 排行和百分比
	if incentiveConfigType == excitationModels.IncentiveConfigTypeOnline {
		// 目标达成率 = （线上平台销售 / 个人销售金额）* 100

		var onlineTopArr []models.MemberTopSortList
		var onlineTopMap = make(map[int64]models.MemberTopSortList)

		// 获取导购线上销售金额
		// memberOnlineListMap, err := microScrmDao.NewMemberSalesVolume().GetMemberSalesVolume(c, memberReq.MemberSalesVolumeReq{
		//	EnterpriseId: enterpriseId,
		//	StartTime:    startTime,
		//	EndTime:      endTime,
		//	SearchType:   microScrmModel.ONLINE_SALES_VOLUME,
		// })
		// if err != nil {
		//	return onlineTopMap, err
		// }

		// 获取导购线上销售占比
		volumeMap, err := microScrmDao.NewMemberSalesVolume().GetOnlineRatio(c, memberReq.MemberSalesVolumeReq{
			EnterpriseId: enterpriseId,
			StartTime:    startTime,
			EndTime:      endTime,
			SearchType:   microScrmModel.ALL_SALES_VOLUME,
		})
		if err != nil {
			return onlineTopMap, err
		}
		for memberId, tmpValue := range volumeMap {
			var temp models.MemberTopSortList
			temp.MemberId = memberId
			temp.Ratio = tmpValue * 100
			onlineTopArr = append(onlineTopArr, temp)
		}

		// 导购线上销售金额排行
		// for memberId, tmpValue := range memberOnlineListMap {
		//	var temp models.MemberTopSortList
		//	temp.MemberId = memberId
		//	temp.Value = tmpValue
		//	temp.Ratio = 0
		//	if _, ok := volumeMap[memberId]; ok {
		//		temp.Ratio = volumeMap[memberId] * 100
		//	}
		//	onlineTopArr = append(onlineTopArr, temp)
		// }

		// 按照线上销售降序排序
		sort.Slice(onlineTopArr, func(i, j int) bool {
			return onlineTopArr[i].Ratio > onlineTopArr[j].Ratio
		})
		for k, v := range onlineTopArr {

			var tempMap models.MemberTopSortList
			tempMap.MemberId = v.MemberId
			tempMap.SortTop = k + 1
			tempMap.Ratio = v.Ratio
			onlineTopMap[v.MemberId] = tempMap
		}

		return onlineTopMap, nil

	}

	return res, nil

}

// 按照时间段出来导购关联能力图谱数据
func (job *RelationRule) handleMemberDataLimit(c ctx.Context, enterpriseId int64, incentiveRelation excitationModels.ExcitationIncentiveRelationModel, resArr []excitationResp.RelationRuleDataResp, resTopMap map[int64]models.MemberTopSortList) error {
	//获取全部员工
	memberList, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"status":        membersModel.MembersStatusEnable,
	})
	if err != nil {
		return err
	}
	if len(memberList) > 0 {
		for _, row := range memberList {
			memberDataInfo, err := excitationSvc.NewMemberDataService().CheckOrCreateMemberData(c, enterpriseId, row.ID)

			if err != nil {
				return err
			}

			// 当前导购信息
			var params excitationReq.UpdateDimensionFieldReq
			params.Id = memberDataInfo.ID
			params.EnterpriseId = enterpriseId
			params.MemberId = row.ID
			params.DimensionValue = ""
			params.DimensionRuleValue = 0
			params.IncentiveConfigType = incentiveRelation.IncentiveConfigType
			topTotal := int64(len(memberList)) //需要排行的员工总数
			//判断是否区分奥莱零售
			storeType := int64(excitationModels.IsDiffStore_No)
			if incentiveRelation.IsDiffStore == excitationModels.IsDiffStore_Yes {
				//获取当前导购所属奥莱或者零售
				storeType, _ = department.NewMemberDepartmentsService().GetMemberStoreType(c, enterpriseId, params.MemberId)
				if storeType == 0 {
					//不属于奥莱也不属于零售则直接跳过
					continue
				}
				topTotal, _ = department.NewMemberDepartmentsService().GetMemberCountByStoreType(c, enterpriseId, storeType)

			}
			//筛选符合当前导购的规则
			var tmpArr []excitationResp.RelationRuleDataResp
			for _, v := range resArr {
				if v.StoreType == int(storeType) {
					tmp := v
					tmpArr = append(tmpArr, tmp)
				}
			}
			if len(tmpArr) == 0 || incentiveRelation.Status == excitationModels.IncentiveConfigStatusOff {
				// 没有配置关联能力图谱更新导购数据
				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				continue
			}
			// 匹配每个导购的 各维度关联能力指标 由高到低 匹配到即可完成执行
			if incentiveRelation.IncentiveConfigType == excitationModels.IncentiveConfigTypeInvite {
				errs := job.handleRelationRule1(c, enterpriseId, params, tmpArr, resTopMap)

				if errs != nil {
					return errs
				}
				continue

			}

			if incentiveRelation.IncentiveConfigType == excitationModels.IncentiveConfigTypeTask {
				errs := job.handleRelationRule2(c, enterpriseId, params, tmpArr, resTopMap)

				if errs != nil {
					return errs
				}
				continue

			}
			if incentiveRelation.IncentiveConfigType == excitationModels.IncentiveConfigTypeInteract {
				errs := job.handleRelationRule3(c, topTotal, params, tmpArr, resTopMap)

				if errs != nil {
					return errs
				}
				continue

			}
			if incentiveRelation.IncentiveConfigType == excitationModels.IncentiveConfigTypeTrain {
				errs := job.handleRelationRule4(c, topTotal, params, tmpArr, resTopMap)

				if errs != nil {
					return errs
				}
				continue

			}
			if incentiveRelation.IncentiveConfigType == excitationModels.IncentiveConfigTypeOnline {
				errs := job.handleRelationRule5(c, enterpriseId, params, tmpArr, resTopMap)

				if errs != nil {
					return errs
				}
				continue

			}
		}
	}

	return nil

}

func (job *RelationRule) handleRelationRule1(c ctx.Context, enterpriseId int64, params excitationReq.UpdateDimensionFieldReq, inviteArr []excitationResp.RelationRuleDataResp, inviteTopMap map[int64]models.MemberTopSortList) error {

	// 邀约 -- 排行和百分比
	// 月度 -- 邀约目标达成率 = 邀约金额 / 个人销售金额 x 100%
	// 判断指标
	for _, value := range inviteArr {
		// 排序和总比赋值
		inviteMemberTopRatio := 0.00
		inviteMemberTopSortTop := 0
		haveOk := true
		if _, ok := inviteTopMap[params.MemberId]; ok {
			inviteMemberTopRatio = inviteTopMap[params.MemberId].Ratio
			inviteMemberTopSortTop = inviteTopMap[params.MemberId].SortTop
		} else {
			haveOk = false
		}

		switch value.RuleCondition { // 指标条件: 1top榜  2高于 3达到 4低于
		case excitationModels.RuleCondition1:
			// 1top榜
			if inviteMemberTopSortTop > 0 && inviteMemberTopSortTop <= value.RuleValue {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				params.DimensionRuleStr = fmt.Sprintf("TOP %d", inviteMemberTopSortTop)
				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		case excitationModels.RuleCondition2:

			// 2 高于
			if inviteMemberTopRatio > 0 && inviteMemberTopRatio > float64(value.RuleValue) {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				tempValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", inviteMemberTopRatio), 64)
				params.DimensionRuleStr = fmt.Sprintf("目标达成率 %v", tempValue) + "%"
				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		case excitationModels.RuleCondition3:
			// 3 达到
			if inviteMemberTopRatio > 0 && inviteMemberTopRatio >= float64(value.RuleValue) {

				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				tempValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", inviteMemberTopRatio), 64)
				params.DimensionRuleStr = fmt.Sprintf("目标达成率 %v", tempValue) + "%"

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		case excitationModels.RuleCondition4:

			// 4 低于
			if inviteMemberTopRatio < float64(value.RuleValue) ||
				(inviteMemberTopRatio == 0.00 && haveOk == false) {

				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				if inviteMemberTopRatio > 0 {
					tempValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", inviteMemberTopRatio), 64)
					params.DimensionRuleStr = fmt.Sprintf("目标达成率 %v", tempValue) + "%"
				} else {
					params.DimensionRuleStr = fmt.Sprintf("目标达成率 低于 %v", value.RuleValue) + "%"
				}

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		default:
			err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
			if err != nil {
				return err
			}
			continue

		}
	}

	return nil
}

func (job *RelationRule) handleRelationRule2(c ctx.Context, enterpriseId int64, params excitationReq.UpdateDimensionFieldReq, taskArr []excitationResp.RelationRuleDataResp, planTaskMap map[int64]models.MemberTopSortList) error {

	// 任务 -- 百分比
	// 判断指标
	for _, value := range taskArr {

		// 排序和总比赋值
		planTaskMemberTopRatio := 0.00
		planTaskMemberTopSortTop := 0
		haveOk := true
		if _, ok := planTaskMap[params.MemberId]; ok {
			planTaskMemberTopRatio = planTaskMap[params.MemberId].Ratio
			planTaskMemberTopSortTop = planTaskMap[params.MemberId].SortTop
		} else {
			haveOk = false
		}

		switch value.RuleCondition { // 指标条件: 1top榜  2高于 3达到 4低于
		case excitationModels.RuleCondition1:
			if planTaskMemberTopSortTop > 0 && planTaskMemberTopSortTop <= value.RuleValue {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				params.DimensionRuleStr = fmt.Sprintf("TOP %d", planTaskMemberTopSortTop)
				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue

		case excitationModels.RuleCondition2:
			if planTaskMemberTopRatio > float64(value.RuleValue) {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				tempValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", planTaskMemberTopRatio), 64)
				params.DimensionRuleStr = fmt.Sprintf("任务完成率 %v", tempValue) + "%"

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue

		case excitationModels.RuleCondition3:

			if planTaskMemberTopRatio >= float64(value.RuleValue) {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				tempValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", planTaskMemberTopRatio), 64)
				params.DimensionRuleStr = fmt.Sprintf("任务完成率 %v", tempValue) + "%"

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue

		case excitationModels.RuleCondition4:

			if planTaskMemberTopRatio < float64(value.RuleValue) ||
				(planTaskMemberTopRatio == 0.00 && haveOk == false) {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue

				if planTaskMemberTopRatio > 0 {
					tempValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", planTaskMemberTopRatio), 64)
					params.DimensionRuleStr = fmt.Sprintf("任务完成率 %v", tempValue) + "%"
				} else {
					params.DimensionRuleStr = fmt.Sprintf("任务完成率 低于 %v", value.RuleValue) + "%"
				}

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		default:
			err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
			if err != nil {
				return err
			}
			continue

		}
	}

	return nil
}

func (job *RelationRule) handleRelationRule3(c ctx.Context, topTotal int64, params excitationReq.UpdateDimensionFieldReq, interactArr []excitationResp.RelationRuleDataResp, interactTopMap map[int64]models.MemberTopSortList) error {

	// 客户互动 -- 排行
	// 判断指标
	for _, value := range interactArr {

		haveOk := true
		var interactMemberTopSortTop = 0
		if _, ok := interactTopMap[params.MemberId]; ok {
			interactMemberTopSortTop = interactTopMap[params.MemberId].SortTop
		} else {
			haveOk = false
		}
		switch value.RuleCondition { // 指标条件: 1top榜  2高于 3达到 4低于 5 top百分比
		case excitationModels.RuleCondition2:
			// 排名高于
			if interactMemberTopSortTop > 0 && interactMemberTopSortTop < value.RuleValue {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				params.DimensionRuleStr = fmt.Sprintf("TOP %d", interactMemberTopSortTop)
				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		//达到和top一样
		case excitationModels.RuleCondition1,
			excitationModels.RuleCondition3:
			// 排名达到
			if interactMemberTopSortTop > 0 && interactMemberTopSortTop <= value.RuleValue {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				params.DimensionRuleStr = fmt.Sprintf("TOP %d", interactMemberTopSortTop)

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		case excitationModels.RuleCondition4:
			if interactMemberTopSortTop > value.RuleValue ||
				(interactMemberTopSortTop == 0 && haveOk == false) {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				if interactMemberTopSortTop > 0 {
					params.DimensionRuleStr = fmt.Sprintf("TOP %d", interactMemberTopSortTop)
				} else {
					params.DimensionRuleStr = fmt.Sprintf("TOP 低于 %d", value.RuleValue)
				}

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		case excitationModels.RuleCondition5: //top百分比
			topPercent := float64(interactMemberTopSortTop) / float64(topTotal) * 100
			if interactMemberTopSortTop == 0 {
				//不存在业务数据时，获取到的业务top为0，top百分比为100
				topPercent = 100
			}
			// 排名高于
			if topPercent <= float64(value.RuleValue) {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				params.DimensionRuleStr = fmt.Sprintf("TOP %d%%", value.RuleValue)
				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		default:
			err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
			if err != nil {
				return err
			}
			continue
		}
	}

	return nil
}

func (job *RelationRule) handleRelationRule4(c ctx.Context, topTotal int64, params excitationReq.UpdateDimensionFieldReq, trainArr []excitationResp.RelationRuleDataResp, trainTopMap map[int64]models.MemberTopSortList) error {

	// 培训积分 -- 排行
	// 判断指标
	for _, value := range trainArr {

		haveOk := true
		var trainMemberTopSortTop = 0
		if _, ok := trainTopMap[params.MemberId]; ok {
			trainMemberTopSortTop = trainTopMap[params.MemberId].SortTop
		} else {
			haveOk = false
		}

		switch value.RuleCondition { // 指标条件: 1top榜  2高于 3达到 4低于
		case excitationModels.RuleCondition2:
			// 2 高于
			if trainMemberTopSortTop > 0 && trainMemberTopSortTop < value.RuleValue {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				params.DimensionRuleStr = fmt.Sprintf("TOP %d", trainMemberTopSortTop)

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
			//纯名次，top和达到一样
		case excitationModels.RuleCondition1,
			excitationModels.RuleCondition3:

			// 3 达到
			if trainMemberTopSortTop > 0 && trainMemberTopSortTop <= value.RuleValue {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				params.DimensionRuleStr = fmt.Sprintf("TOP %d", trainMemberTopSortTop)

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		case excitationModels.RuleCondition4:

			// 4 低于
			if trainMemberTopSortTop > value.RuleValue ||
				(trainMemberTopSortTop == 0 && haveOk == false) {

				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue

				if trainMemberTopSortTop > 0 {
					params.DimensionRuleStr = fmt.Sprintf("TOP %d", trainMemberTopSortTop)
				} else {
					params.DimensionRuleStr = fmt.Sprintf("TOP 低于 %d", value.RuleValue)
				}

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		case excitationModels.RuleCondition5: //top百分比
			topPercent := float64(trainMemberTopSortTop) / float64(topTotal) * 100
			if trainMemberTopSortTop == 0 {
				//不存在业务数据时，获取到的业务top为0，top百分比为100
				topPercent = 100
			}
			// 排名高于
			if topPercent <= float64(value.RuleValue) {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				params.DimensionRuleStr = fmt.Sprintf("TOP %d%%", value.RuleValue)
				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		default:
			err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
			if err != nil {
				return err
			}
			continue

		}
	}

	return nil
}

func (job *RelationRule) handleRelationRule5(c ctx.Context, enterpriseId int64, params excitationReq.UpdateDimensionFieldReq, onlineArr []excitationResp.RelationRuleDataResp, onlineTopMap map[int64]models.MemberTopSortList) error {

	// 线上能力 -- 排行和百分比
	// 目标达成率 = （线上平台销售 / 个人销售金额）* 100
	// 判断指标
	for _, value := range onlineArr {

		haveOk := true
		var onlineMemberTopSortTop = 0
		var onlineMemberTopRatio = 0.00
		if _, ok := onlineTopMap[params.MemberId]; ok {
			onlineMemberTopSortTop = onlineTopMap[params.MemberId].SortTop
			onlineMemberTopRatio = onlineTopMap[params.MemberId].Ratio
		} else {
			haveOk = false
		}

		switch value.RuleCondition { // 指标条件: 1top榜  2高于 3达到 4低于
		case excitationModels.RuleCondition1:
			// 1top榜 按照金额
			if onlineMemberTopSortTop > 0 && onlineMemberTopSortTop <= value.RuleValue {
				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				params.DimensionRuleStr = fmt.Sprintf("TOP %d", onlineMemberTopSortTop)

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue

		case excitationModels.RuleCondition2:

			if onlineMemberTopRatio > float64(value.RuleValue) {

				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				tempValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", onlineMemberTopRatio), 64)
				params.DimensionRuleStr = fmt.Sprintf("目标完成率 %v", tempValue) + "%"

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		case excitationModels.RuleCondition3:
			if onlineMemberTopRatio >= float64(value.RuleValue) {

				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue
				tempValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", onlineMemberTopRatio), 64)
				params.DimensionRuleStr = fmt.Sprintf("目标完成率 %v", tempValue) + "%"

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue

		case excitationModels.RuleCondition4:

			if onlineMemberTopRatio < float64(value.RuleValue) ||
				(onlineMemberTopRatio == 0.00 && haveOk == false) {

				// 更新导购数据
				params.IncentiveConfigType = value.IncentiveConfigType
				params.DimensionValue = value.RuleGive
				params.DimensionRuleValue = value.RuleValue

				if onlineMemberTopRatio > 0 {
					tempValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", onlineMemberTopRatio), 64)
					params.DimensionRuleStr = fmt.Sprintf("目标完成率 %v", tempValue) + "%"
				} else {
					params.DimensionRuleStr = fmt.Sprintf("目标完成率 低于 %v", value.RuleValue) + "%"
				}

				err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
				if err != nil {
					return err
				}
				return nil
			}
			continue
		default:
			err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
			if err != nil {
				return err
			}
			continue

		}
	}

	return nil
}

// 根据时间循环处理各维度关联能力图谱规则
func (job *RelationRule) handleRelationRuleForTime(c ctx.Context, enterpriseId, startTime, endTime int64, memberInfo membersModel.MemberModel, inviteArr, taskArr, interactArr, trainArr, onlineArr []excitationResp.RelationRuleDataResp) error {

	memberDataInfo, err := excitationSvc.NewMemberDataService().CheckOrCreateMemberData(c, enterpriseId, memberInfo.ID)

	if err != nil {
		return err
	}

	var params excitationReq.UpdateDimensionFieldReq

	// 当前导购信息
	params.Id = memberDataInfo.ID
	params.EnterpriseId = memberDataInfo.EnterpriseID
	params.MemberId = memberDataInfo.MemberID

	// 邀约 -- 排行和百分比
	if len(inviteArr) > 0 {
		// 月度 -- 邀约目标达成率 = 邀约金额 / 个人销售金额 x 100%
		achievementParams := reserveReq.AchievementParams{
			EnterpriseId: enterpriseId,
			Type:         excitationModels.RuleCategoryId5,
			StartTime:    startTime,
			EndTime:      endTime,
		}
		// 获取所有导购邀约-邀约金额
		achievementMap, err := microScrmDao.NewReservesOrder().GetAllMemberAchievement(c, achievementParams)
		if err != nil {
			return err
		}
		// 获取所有导购-个人总销售金额
		volumeMap, err := microScrmDao.NewMemberSalesVolume().GetMemberSalesVolumeByReserve(c,
			memberReq.MemberSalesVolumeReq{
				EnterpriseId: enterpriseId,
				StartTime:    startTime,
				EndTime:      endTime,
			})
		if err != nil {
			return err
		}
		var inviteTopArr []models.MemberTopSortList
		// 获取邀约金额在个人总销售金额的占比
		for memberId, tmpData := range achievementMap {
			var temp models.MemberTopSortList
			temp.Value = tmpData
			temp.Ratio = 0.00
			if _, ok := volumeMap[memberId]; ok {
				if volumeMap[memberId] > 0 {
					temp.Ratio = tmpData / volumeMap[memberId]
				}
			}
			inviteTopArr = append(inviteTopArr, temp)
		}

		// 1 按照邀约目标达成率降序排序
		sort.Slice(inviteTopArr, func(i, j int) bool {
			return inviteTopArr[i].Ratio > inviteTopArr[j].Ratio
		})
		var inviteTopMap = make(map[int64]models.MemberTopSortList)
		for k, v := range inviteTopArr {
			var tempMap models.MemberTopSortList
			tempMap.MemberId = v.MemberId
			tempMap.Value = v.Value
			tempMap.Ratio = v.Ratio
			tempMap.SortTop = k + 1
			inviteTopMap[v.MemberId] = tempMap
		}

		// 判断指标
		for _, value := range inviteArr {

			// 排序和总比赋值
			inviteMemberTopRatio := 0.00
			inviteMemberTopSortTop := 0
			if _, ok := inviteTopMap[params.MemberId]; ok {
				inviteMemberTopRatio = inviteTopMap[params.MemberId].Value
				inviteMemberTopSortTop = inviteTopMap[params.MemberId].SortTop
			}

			switch value.RuleCondition { // 指标条件: 1top榜  2高于 3达到 4低于
			case excitationModels.RuleCondition1:

				// 1top榜
				if inviteMemberTopSortTop > 0 && inviteMemberTopSortTop <= value.RuleValue {
					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break

			case excitationModels.RuleCondition2:

				// 2 高于
				if inviteMemberTopRatio > 0 && inviteMemberTopRatio > float64(value.RuleValue) {
					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break

			case excitationModels.RuleCondition3:
				// 3 高于
				if inviteMemberTopRatio > 0 && inviteMemberTopRatio >= float64(value.RuleValue) {

					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break

			case excitationModels.RuleCondition4:

				// 4 低于
				if inviteMemberTopRatio > 0 && inviteMemberTopRatio < float64(value.RuleValue) {

					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break

			}
		}
	} else {
		// 没有配置关联能力图谱更新导购数据
		params.IncentiveConfigType = excitationModels.IncentiveConfigTypeInvite
		params.DimensionValue = ""
		params.DimensionRuleValue = 0
		err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
		if err != nil {
			return err
		}
	}

	// 任务 -- 百分比
	if len(taskArr) > 0 {

		// 查询所有导购的-任务完成率
		planTaskRate, err := microScrmDao.NewPlanTaskAllotLogDao().GetAllMemberAchievement(c, enterpriseId, excitationModels.RuleCategoryId8, startTime, endTime)
		if err != nil {
			return err
		}
		var userIds []string
		var planTaskTopArr []models.MemberTopSortList
		for userId, ratio := range planTaskRate {
			userIds = append(userIds, userId)
			var tempArr models.MemberTopSortList
			tempArr.MemberUserid = userId
			tempArr.Ratio = ratio // 任务完成率 百分比值
			planTaskTopArr = append(planTaskTopArr, tempArr)

		}
		// 获取userId对应的memberId
		memberList, err := membersDao.NewMembers().GetManyByUserIds(c, enterpriseId, userIds)
		if err != nil {
			return err
		}
		var memberMap = make(map[string]int64)
		for _, member := range memberList {
			memberMap[member.Userid] = member.ID
		}
		if err != nil {
			return err
		}

		// 2 按照任务完成率降序排序
		sort.Slice(planTaskTopArr, func(i, j int) bool {
			return planTaskTopArr[i].Ratio > planTaskTopArr[j].Ratio
		})
		var planTaskMap = make(map[int64]models.MemberTopSortList)
		for k, v := range planTaskTopArr {
			var tempMap models.MemberTopSortList
			// 匹配导购id
			if _, ok := memberMap[v.MemberUserid]; ok {
				tempMap.MemberId = memberMap[v.MemberUserid]
			}

			tempMap.Ratio = v.Ratio
			tempMap.SortTop = k + 1
			planTaskMap[v.MemberId] = tempMap
		}

		// 判断指标
		for _, value := range taskArr {

			// 排序和总比赋值
			planTaskMemberTopRatio := 0.00
			// planTaskMemberTopSortTop := 0
			if _, ok := planTaskMap[params.MemberId]; ok {
				planTaskMemberTopRatio = planTaskMap[params.MemberId].Ratio
				// planTaskMemberTopSortTop = planTaskMap[params.MemberId].SortTop
			}

			switch value.RuleCondition { // 指标条件: 1top榜  2高于 3达到 4低于
			case excitationModels.RuleCondition1:

				break

			case excitationModels.RuleCondition2:
				if planTaskMemberTopRatio > float64(value.RuleValue) {
					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}

				break

			case excitationModels.RuleCondition3:

				if planTaskMemberTopRatio >= float64(value.RuleValue) {
					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break

			case excitationModels.RuleCondition4:

				if planTaskMemberTopRatio < float64(value.RuleValue) {
					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}

				}
				break

			}
		}
	} else {
		// 没有配置关联能力图谱更新导购数据
		params.IncentiveConfigType = excitationModels.IncentiveConfigTypeTask
		params.DimensionValue = ""
		params.DimensionRuleValue = 0
		err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
		if err != nil {
			return err
		}
	}

	// 客户互动 -- 排行
	if len(interactArr) > 0 {
		// 获取所有导购时间段内聊天互动排行榜
		interactMemberList, err := microScrmDao.NewMemberExternalBehaviorData().GetChatCountAllTop(c, enterpriseId, startTime, endTime)
		if err != nil {
			return err
		}
		var interactTopArr []models.MemberTopSortList
		for _, tmpValue := range interactMemberList {
			var temp models.MemberTopSortList
			temp.MemberId = tmpValue.MemberId
			temp.Total = tmpValue.Total
			interactTopArr = append(interactTopArr, temp)
		}

		// Top 排行map 按照客户互动数排序
		sort.Slice(interactTopArr, func(i, j int) bool {
			return interactTopArr[i].Total > interactTopArr[j].Total
		})
		var interactTopMap = make(map[int64]models.MemberTopSortList)
		for k, v := range interactTopArr {
			interactTopArr[k].SortTop = k + 1

			var tempMap models.MemberTopSortList
			tempMap.MemberId = v.MemberId
			tempMap.Total = v.Total
			tempMap.SortTop = k + 1
			interactTopMap[v.MemberId] = tempMap
		}

		// 判断指标
		for _, value := range interactArr {

			var interactMemberTopSortTop = 0
			if _, ok := interactTopMap[params.MemberId]; ok {
				interactMemberTopSortTop = interactTopMap[params.MemberId].SortTop
			}
			switch value.RuleCondition { // 指标条件: 1top榜  2高于 3达到 4低于
			case excitationModels.RuleCondition1:

				break

			case excitationModels.RuleCondition2:
				if interactMemberTopSortTop > 0 && interactMemberTopSortTop < value.RuleValue {
					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break

			case excitationModels.RuleCondition3:
				if interactMemberTopSortTop > 0 && interactMemberTopSortTop <= value.RuleValue {
					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break

			case excitationModels.RuleCondition4:
				if interactMemberTopSortTop > 0 && interactMemberTopSortTop > value.RuleValue {
					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break

			}
		}
	} else {
		// 没有配置关联能力图谱更新导购数据
		params.IncentiveConfigType = excitationModels.IncentiveConfigTypeInteract
		params.DimensionValue = ""
		params.DimensionRuleValue = 0
		err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
		if err != nil {
			return err
		}
	}

	// 培训互动积分 -- 排行
	if len(trainArr) > 0 {

		// 查询导购培训互动积分排行
		memberPointsList, err := excitationDao.NewPointsRecord().GetTopList(c, enterpriseId, excitationModels.RuleCategoryId15, startTime, endTime)
		if err != nil {
			return err
		}

		// 排行
		var trainTopArr []models.MemberTopSortList
		for memberId, tmpValue := range memberPointsList {
			var temp models.MemberTopSortList
			temp.MemberId = memberId
			temp.Value = tmpValue // 积分
			trainTopArr = append(trainTopArr, temp)
		}

		// 4 按照积分降序排序
		sort.Slice(trainTopArr, func(i, j int) bool {
			return trainTopArr[i].Value > trainTopArr[j].Value
		})
		var trainTopMap = make(map[int64]models.MemberTopSortList)
		for k, v := range trainTopArr {

			var tempMap models.MemberTopSortList
			tempMap.MemberId = v.MemberId
			tempMap.Value = v.Value
			tempMap.SortTop = k + 1
			trainTopMap[v.MemberId] = tempMap
		}

		// 判断指标
		for _, value := range trainArr {

			var inviteMemberTopSortTop = 0
			if _, ok := trainTopMap[params.MemberId]; ok {
				inviteMemberTopSortTop = trainTopMap[params.MemberId].SortTop
			}

			switch value.RuleCondition { // 指标条件: 1top榜  2高于 3达到 4低于
			case excitationModels.RuleCondition1:

				break

			case excitationModels.RuleCondition2:

				// 2 高于
				if inviteMemberTopSortTop > 0 && inviteMemberTopSortTop < value.RuleValue {

					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break

			case excitationModels.RuleCondition3:

				// 3 达到
				if inviteMemberTopSortTop > 0 && inviteMemberTopSortTop <= value.RuleValue {

					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break

			case excitationModels.RuleCondition4:

				// 4 低于
				if inviteMemberTopSortTop > 0 && inviteMemberTopSortTop > value.RuleValue {

					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break
			}
		}
	} else {
		// 没有配置关联能力图谱更新导购数据
		params.IncentiveConfigType = excitationModels.IncentiveConfigTypeTrain
		params.DimensionValue = ""
		params.DimensionRuleValue = 0
		err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
		if err != nil {
			return err
		}
	}

	// 线上能力 -- 排行和百分比
	if len(onlineArr) > 0 {
		// 目标达成率 = （线上平台销售 / 个人销售金额）* 100
		// 获取导购线上销售金额
		memberOnlineListMap, err := microScrmDao.NewMemberSalesVolume().GetMemberSalesVolume(c, memberReq.MemberSalesVolumeReq{
			EnterpriseId: enterpriseId,
			StartTime:    startTime,
			EndTime:      endTime,
			SearchType:   microScrmModel.ONLINE_SALES_VOLUME,
		})
		if err != nil {
			return err
		}

		// 获取导购线上销售占比
		volumeMap, err := microScrmDao.NewMemberSalesVolume().GetOnlineRatio(c, memberReq.MemberSalesVolumeReq{
			EnterpriseId: enterpriseId,
			StartTime:    startTime,
			EndTime:      endTime,
			SearchType:   microScrmModel.ALL_SALES_VOLUME,
		})
		if err != nil {
			return err
		}

		// 导购线上销售金额排行
		var onlineTopArr []models.MemberTopSortList
		for memberId, tmpValue := range memberOnlineListMap {
			var temp models.MemberTopSortList
			temp.MemberId = memberId
			temp.Value = tmpValue
			temp.Ratio = 0
			if _, ok := volumeMap[memberId]; ok {
				temp.Ratio = volumeMap[memberId] * 100
			}
			onlineTopArr = append(onlineTopArr, temp)
		}

		// 按照线上销售金额降序排序
		sort.Slice(onlineTopArr, func(i, j int) bool {
			return onlineTopArr[i].Value > onlineTopArr[j].Value
		})
		var onlineTopMap = make(map[int64]models.MemberTopSortList)
		for k, v := range onlineTopArr {

			var tempMap models.MemberTopSortList
			tempMap.MemberId = v.MemberId
			tempMap.Value = v.Value
			tempMap.SortTop = k + 1
			tempMap.Ratio = v.Ratio
			onlineTopMap[v.MemberId] = tempMap
		}

		// 判断指标
		for _, value := range onlineArr {

			var onlineMemberTopSortTop = 0
			var onlineMemberTopRatio = 0.00
			if _, ok := onlineTopMap[params.MemberId]; ok {
				onlineMemberTopSortTop = onlineTopMap[params.MemberId].SortTop
				onlineMemberTopRatio = onlineTopMap[params.MemberId].Ratio
			}

			switch value.RuleCondition { // 指标条件: 1top榜  2高于 3达到 4低于
			case excitationModels.RuleCondition1:
				// 1top榜 按照金额
				if onlineMemberTopSortTop > 0 && onlineMemberTopSortTop <= value.RuleValue {
					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break

			case excitationModels.RuleCondition2:

				if onlineMemberTopRatio > float64(value.RuleValue) {

					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break

			case excitationModels.RuleCondition3:
				if onlineMemberTopRatio >= float64(value.RuleValue) {

					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break

			case excitationModels.RuleCondition4:

				if onlineMemberTopRatio < float64(value.RuleValue) {

					// 更新导购数据
					params.IncentiveConfigType = value.IncentiveConfigType
					params.DimensionValue = value.RuleGive
					params.DimensionRuleValue = value.RuleValue
					err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
					if err != nil {
						return err
					}
				}
				break
			}
		}
	} else {
		// 没有配置关联能力图谱更新导购数据
		params.IncentiveConfigType = excitationModels.IncentiveConfigTypeOnline
		params.DimensionValue = ""
		params.DimensionRuleValue = 0
		err := excitationSvc.NewMemberDataService().SaveDimensionField(c, params)
		if err != nil {
			return err
		}
	}

	return nil
}

// 各维度关联能力图谱规则数组处理
func (job *RelationRule) getRelationOneArr(c ctx.Context, enterpriseId int64, incentiveRelation []excitationModels.ExcitationIncentiveRelationModel, incentiveConfigType int) (resArr []excitationResp.RelationRuleDataResp) {
	var incentiveRelationIds []int64
	for _, v := range incentiveRelation {
		incentiveRelationIds = append(incentiveRelationIds, v.ID)
	}
	// 能力图谱规则
	relationRuleResult, _ := excitationDao.NewIncentiveRelationRule().GetMany(c, "enterprise_id = ? and incentive_relation_id in (?) and incentive_config_type = ?", enterpriseId, incentiveRelationIds, incentiveConfigType)

	// 指标规则等级排序
	relationRuleMap := excitationSvc.NewMemberDataService().RelationRuleMap()

	// 处理各维度关联图谱规则排序
	resArr = job.sortRelationRule(relationRuleResult, relationRuleMap)

	return
}

// 指标能力由高到低排序
func (job *RelationRule) sortRelationRule(relationRules []excitationModels.ExcitationIncentiveRelationRuleModel, relationRuleMap map[string]excitationReq.RelationRuleArray) (relationRuleArr []excitationResp.RelationRuleDataResp) {

	for _, item := range relationRules {
		sortKey := 0
		if _, ok := relationRuleMap[item.RuleGive]; ok {
			sortKey = relationRuleMap[item.RuleGive].Key
		}
		relationRuleArr = append(relationRuleArr, excitationResp.RelationRuleDataResp{

			Id:                  item.ID,
			IncentiveConfigType: item.IncentiveConfigType,
			IncentiveRelationId: item.IncentiveRelationID,
			RuleCondition:       item.RuleCondition,
			RuleValue:           item.RuleValue,
			RuleGive:            item.RuleGive,
			SortKey:             sortKey,
			StoreType:           item.StoreType,
		})
	}

	sort.Slice(relationRuleArr, func(i, j int) bool {
		return relationRuleArr[i].SortKey > relationRuleArr[j].SortKey
	})

	return
}
