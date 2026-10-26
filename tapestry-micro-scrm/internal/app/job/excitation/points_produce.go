package excitation

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	memberReq "micro-scrm/internal/app/http/request/member"
	reserveReq "micro-scrm/internal/app/http/request/reserve"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
	"time"
)

type PointsProduce struct {
	name                 string
	pointsProduceSvc     *excitationSvc.PointsProduceService
	taskRecordDao        *excitationDao.TaskRecord
	memberSalesVolumeDao *microScrmDao.MemberSalesVolume
}

func NewPointsProduce() *PointsProduce {
	return &PointsProduce{
		name:                 "累计数据获得积分任务",
		pointsProduceSvc:     excitationSvc.NewPointsProduceService(),
		taskRecordDao:        excitationDao.NewTaskRecord(),
		memberSalesVolumeDao: microScrmDao.NewMemberSalesVolume(),
	}
}

func (job *PointsProduce) Run(c ctx.Context, handler *xjob.JobHandler) error {
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

func (job *PointsProduce) Handle(c ctx.Context, enterpriseId int64) error {
	// 获取已配置的积分规则
	paramsTime := time.Now()
	//paramsTime := time.Unix(1685815200, 0)
	nowUnix := paramsTime.Unix() // 自定义执行时间
	yesterdayUnix := nowUnix - 86400
	query := "enterprise_id = ? and (validity_type = ? or (validity_type = ? AND start_time < ? AND end_time > ?))"
	// 获取已配置的积分规则（没配置规则的不需要积累数据）
	ruleList, err := excitationDao.NewIntegralRuleDetail().
		GetOnCategoryList(c, query, enterpriseId, excitationModels.ValidityTypeLong, excitationModels.ValidityTypeDeadline, nowUnix, nowUnix)
	if err != nil {
		return err
	}
	if len(ruleList) == 0 {
		return nil
	}
	// 获取定时任务执行前一天的相关财月（也就是符合数据统计的财月）
	fiscalMonth, err := microWorkDao.NewFiscalMonth().GetDataByYearAndTime(c, enterpriseId, yesterdayUnix)
	if err != nil {
		return err
	}
	// 循环所有的任务列表
	for _, rule := range ruleList {
		uniqueStr := paramsTime.Format("2006-01-02") + fmt.Sprintf("-%d", rule.RuleCategoryID)
		exists, _ := job.taskRecordDao.GetOne(c, "enterprise_id = ? AND unique_str = ?", enterpriseId, uniqueStr)
		if exists.ID > 0 {
			continue
		}
		exeSign := false
		switch rule.RuleCategoryID {
		// 每周
		case excitationModels.RuleCategoryId2:
			if int(paramsTime.Weekday()) > 0 {
				// 财周任务应该是周天执行
				continue
			}
			err := job.weekDayJob(c, enterpriseId, rule.RuleCategoryID, paramsTime)
			if err != nil {
				log.Info(c).Msgf("%s::积分定时任务执行失败，执行类目：%d,每周执行的: %v", job.name, rule.RuleCategoryID, err)
			}
			exeSign = true
		// 每月
		case excitationModels.RuleCategoryId3,
			excitationModels.RuleCategoryId5,
			excitationModels.RuleCategoryId8,
			excitationModels.RuleCategoryId11,
			excitationModels.RuleCategoryId12,
			excitationModels.RuleCategoryId13,
			excitationModels.RuleCategoryId16,
			excitationModels.RuleCategoryId19:
			if fiscalMonth.EndTime-yesterdayUnix > 86400 {
				// 财月任务应该是财月结束后一天内执行
				continue
			}
			err := job.monthJob(c, enterpriseId, rule.RuleCategoryID, paramsTime, fiscalMonth)
			if err != nil {
				log.Info(c).Msgf("%s::积分定时任务执行失败，执行类目：%d,每月执行的: %v", job.name, rule.RuleCategoryID, err)
			}
			exeSign = true
			// 每天
			//case 0:
			//	err := job.dayJob(c, enterpriseId, rule.RuleCategoryID, paramsTime)
			//	if err != nil {
			//		log.Info(c).Msgf("%s::积分定时任务执行失败，执行类目：%d,每日执行的: %v", job.name, rule.RuleCategoryID, err)
			//	}
			//	exeSign = true
		}
		if exeSign {
			// 记录执行任务
			err = job.taskRecordDao.CreateData(c, excitationModels.ExcitationTaskRecordModel{
				MicroBaseModel: models.CreateMicroBase(enterpriseId, app.UID()),
				UniqueStr:      uniqueStr,
				Type:           excitationModels.TaskTypePoint,
			})
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// 周任务加积分
func (job *PointsProduce) weekDayJob(c ctx.Context, enterpriseId int64, categoryId int64, paramsTime time.Time) error {
	endTime := helper.GetLastDateOfWeek(helper.GetZeroTimeByTime(paramsTime)).Unix()
	startTime := endTime - 86400*7
	switch categoryId {
	// 每周
	case excitationModels.RuleCategoryId2,
		excitationModels.RuleCategoryId3:
		achievementParams := reserveReq.AchievementParams{
			EnterpriseId: enterpriseId,
			Type:         categoryId,
			StartTime:    startTime,
			EndTime:      endTime,
		}
		// 获取邀约相关统计
		achievement, err := microScrmDao.NewReservesOrder().GetAllMemberAchievement(c, achievementParams)
		if err != nil {
			return err
		}
		for memberId, resData := range achievement {
			produceParams := excitationReq.PointsProduceSave{
				MemberId:     memberId,
				Type:         categoryId,
				Value:        int64(resData),
				UniqueHash:   paramsTime.Format("2006-01-02") + fmt.Sprintf("-%d-%d", memberId, categoryId),
				EnterpriseId: enterpriseId,
			}
			err := job.pointsProduceSvc.Produce(c, paramsTime.Unix(), produceParams)
			if err != nil {
				log.Info(c).Msgf("%s::积分定时任务执行失败，执行类目：%d,每周任务，data:%+v，错误: %+v", job.name, categoryId, produceParams, err)
			}
		}
	}
	return nil
}

// 月任务加积分
func (job *PointsProduce) monthJob(c ctx.Context, enterpriseId int64, categoryId int64, nowTime time.Time, monthData microWorkModel.FiscalMonthModel) error {
	// 财月的统计时间根据参数时间查询开始和结束时间
	startTime := monthData.StartTime
	endTime := monthData.EndTime
	switch categoryId {
	case excitationModels.RuleCategoryId3:
		// 邀约成功率的统计时间范围需要再往前多找一个月（由于邀约的成功匹配有最多一个月的延迟）
		fiscalInfo, err := microWorkDao.NewFiscalMonth().GetFiscalMonth(c, enterpriseId, startTime-86400)
		startTime = fiscalInfo.StartTime
		endTime = fiscalInfo.EndTime
		achievementParams := reserveReq.AchievementParams{
			EnterpriseId: enterpriseId,
			Type:         categoryId,
			StartTime:    startTime,
			EndTime:      endTime,
		}
		// 获取邀约相关统计
		achievement, err := microScrmDao.NewReservesOrder().GetAllMemberAchievement(c, achievementParams)
		if err != nil {
			return err
		}
		for memberId, resData := range achievement {
			produceParams := excitationReq.PointsProduceSave{
				MemberId:     memberId,
				Type:         categoryId,
				Value:        int64(resData),
				UniqueHash:   nowTime.Format("2006-01-02") + fmt.Sprintf("-%d-%d", memberId, categoryId),
				EnterpriseId: enterpriseId,
			}
			err := job.pointsProduceSvc.Produce(c, nowTime.Unix(), produceParams)
			if err != nil {
				log.Info(c).Msgf("%s::积分定时任务执行失败，执行类目：%d,邀约成功率任务，data:%+v，错误: %+v", job.name, categoryId, produceParams, err)
			}
		}
	case excitationModels.RuleCategoryId5: //邀约占比每月
		achievementParams := reserveReq.AchievementParams{
			EnterpriseId: enterpriseId,
			Type:         categoryId,
			StartTime:    startTime,
			EndTime:      endTime,
		}
		// 获取邀约相关统计
		achievement, err := microScrmDao.NewReservesOrder().GetAllMemberAchievement(c, achievementParams)
		if err != nil {
			return err
		}
		// 获取个人总销售金额
		volumeMap, err := job.memberSalesVolumeDao.GetMemberSalesVolumeByReserve(c, memberReq.MemberSalesVolumeReq{
			EnterpriseId: enterpriseId,
			StartTime:    startTime,
			EndTime:      endTime,
		})
		if err != nil {
			return err
		}
		// 获取邀约金额在个人总销售金额的占比
		for memberId, tmpData := range achievement {
			achievement[memberId] = 0
			if _, ok := volumeMap[memberId]; ok {
				if volumeMap[memberId] > 0 {
					achievement[memberId] = tmpData / 100 / volumeMap[memberId] // 金额是分需要除100
				}
			}
		}
	// 任务完成率
	case excitationModels.RuleCategoryId8:
		// 查询所有用户的任务完成率
		achievement, err := microScrmDao.NewPlanTaskAllotLogDao().GetAllMemberAchievement(c, enterpriseId, categoryId, startTime, endTime)
		var userIds []string
		for userId := range achievement {
			userIds = append(userIds, userId)
		}
		userIds = helper.ArrayUniq(userIds)
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
		for userId, resData := range achievement {
			produceParams := excitationReq.PointsProduceSave{
				MemberId:     memberMap[userId],
				Type:         categoryId,
				Value:        int64(resData),
				UniqueHash:   nowTime.Format("2006-01-02") + fmt.Sprintf("-%d-%d", memberMap[userId], categoryId),
				EnterpriseId: enterpriseId,
			}
			err := job.pointsProduceSvc.Produce(c, nowTime.Unix(), produceParams)
			if err != nil {
				log.Info(c).Msgf("%s::积分定时任务执行失败，执行类目：%d,每月任务，data:%+v，错误: %+v", job.name, categoryId, produceParams, err)
			}
		}
	case excitationModels.RuleCategoryId11:
		list, err := microScrmDao.NewUserRelation().GetGroupMemberLossFriendsCount(c, enterpriseId, startTime, endTime)
		if err != nil {
			return err
		}
		for _, member := range list {
			produceParams := excitationReq.PointsProduceSave{
				MemberId:     member.MemberId,
				Type:         categoryId,
				Value:        member.Total,
				UniqueHash:   nowTime.Format("2006-01-02") + fmt.Sprintf("-%d-%d", member.MemberId, categoryId),
				EnterpriseId: enterpriseId,
				SourceRemark: "",
			}
			err := job.pointsProduceSvc.Produce(c, nowTime.Unix(), produceParams)
			if err != nil {
				log.Info(c).Msgf("%s::积分定时任务执行失败，执行类目：%d,每日任务，data:%+v，错误: %+v", job.name, categoryId, produceParams, err)
			}
		}
	case excitationModels.RuleCategoryId12: // 聊天数
		list, err := microScrmDao.NewMemberExternalBehaviorData().
			GetChatCount(c, enterpriseId, startTime, endTime)
		if err != nil {
			return err
		}
		for _, member := range list {
			produceParams := excitationReq.PointsProduceSave{
				MemberId:     member.MemberId,
				Type:         categoryId,
				Value:        member.Total,
				UniqueHash:   nowTime.Format("2006-01-02") + fmt.Sprintf("-%d-%d", member.MemberId, categoryId),
				EnterpriseId: enterpriseId,
			}
			err := job.pointsProduceSvc.Produce(c, nowTime.Unix(), produceParams)
			if err != nil {
				log.Info(c).Msgf("%s::积分定时任务执行失败，执行类目：%d,每月任务，data:%+v，错误: %+v", job.name, categoryId, produceParams, err)
			}
		}
	case excitationModels.RuleCategoryId13: // 聊天回复率
		list, err := microScrmDao.NewMemberExternalBehaviorData().
			GetReplyPercentageCount(c, enterpriseId, startTime, endTime)
		if err != nil {
			return err
		}
		for _, member := range list {
			produceParams := excitationReq.PointsProduceSave{
				MemberId:     member.MemberId,
				Type:         categoryId,
				Value:        member.Total,
				UniqueHash:   nowTime.Format("2006-01-02") + fmt.Sprintf("-%d-%d", member.MemberId, categoryId),
				EnterpriseId: enterpriseId,
			}
			err := job.pointsProduceSvc.Produce(c, nowTime.Unix(), produceParams)
			if err != nil {
				log.Info(c).Msgf("%s::积分定时任务执行失败，执行类目：%d,每月任务，data:%+v，错误: %+v", job.name, categoryId, produceParams, err)
			}
		}
	case excitationModels.RuleCategoryId16:
		volumeMap, err := job.memberSalesVolumeDao.GetOnlineRatio(c, memberReq.MemberSalesVolumeReq{
			EnterpriseId: enterpriseId,
			StartTime:    startTime,
			EndTime:      endTime,
			SearchType:   microScrmModel.ALL_SALES_VOLUME,
		})
		if err != nil {
			return err
		}
		for memberId, resData := range volumeMap {
			produceParams := excitationReq.PointsProduceSave{
				MemberId:     memberId,
				Type:         categoryId,
				Value:        int64(resData * 100),
				UniqueHash:   nowTime.Format("2006-01-02") + fmt.Sprintf("-%d-%d", memberId, categoryId),
				EnterpriseId: enterpriseId,
			}
			err := job.pointsProduceSvc.Produce(c, nowTime.Unix(), produceParams)
			if err != nil {
				log.Info(c).Msgf("%s::积分定时任务执行失败，执行类目：%d,每周任务，data:%+v，错误: %+v", job.name, categoryId, produceParams, err)
			}
		}
	case excitationModels.RuleCategoryId19:
		volumeMap, err := job.memberSalesVolumeDao.GetMemberSalesVolume(c, memberReq.MemberSalesVolumeReq{
			EnterpriseId: enterpriseId,
			StartTime:    startTime,
			EndTime:      endTime,
			SearchType:   microScrmModel.ALL_SALES_VOLUME,
		})
		if err != nil {
			return err
		}
		for memberId, resData := range volumeMap {
			produceParams := excitationReq.PointsProduceSave{
				MemberId:     memberId,
				Type:         categoryId,
				Value:        int64(resData * 100),
				UniqueHash:   nowTime.Format("2006-01-02") + fmt.Sprintf("-%d-%d", memberId, categoryId),
				EnterpriseId: enterpriseId,
			}
			err := job.pointsProduceSvc.Produce(c, nowTime.Unix(), produceParams)
			if err != nil {
				log.Info(c).Msgf("%s::积分定时任务执行失败，执行类目：%d,每周任务，data:%+v，错误: %+v", job.name, categoryId, produceParams, err)
			}
		}
	}
	return nil
}

// 每日任务
func (job *PointsProduce) dayJob(c ctx.Context, enterpriseId int64, categoryId int64, paramsTime time.Time) error {
	//// 每日任务是后一日凌晨执行前一日的任务
	//endTime := helper.GetZeroTimeByTime(paramsTime).Unix()
	//startTime := endTime - 86400
	//switch categoryId {
	//// 客户流失数
	//
	//}
	return nil
}
