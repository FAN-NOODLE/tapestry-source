package excitation

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	customerDao "micro-scrm/internal/app/dao/customer"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	memberReq "micro-scrm/internal/app/http/request/member"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/app/service/department"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	fiscalYearSvc "micro-scrm/internal/app/service/fiscal_year"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
	"sort"
	"strconv"
	"time"
)

type SendMedal struct {
	name             string
	nowTime          int64
	retailMemberIds  []int64
	outletsMemberIds []int64
	isDiffStore      bool
}

func NewSendMedal() *SendMedal {
	return &SendMedal{
		name: "发送奖章任务",
	}
}

func (job *SendMedal) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	var wg threading.RoutineGroup

	for _, id := range enterpriseIds {
		tmpId := id
		wg.Run(func() {
			err = job.Handle(c, tmpId)
			if err != nil {
				app.Logger().Error(c).Msgf("error : %v", err)
			}
		})
	}

	wg.Wait()

	return nil
}

func (job *SendMedal) Handle(c ctx.Context, enterpriseId int64) error {
	nowTime := time.Now()
	//nowTime := time.Unix(1730390400, 0)
	job.nowTime = nowTime.Unix()
	yesterdayTime := nowTime.AddDate(0, 0, -1).Unix()

	// 获取昨天的财月数据
	fiscalMonth, err := microWorkDao.NewFiscalMonth().GetDataByYearAndTime(c, enterpriseId, yesterdayTime)
	if err != nil {
		return err
	}

	// 通过财月数据获财年的所有数据
	fiscalMonthYearAll, err := microWorkDao.NewFiscalMonth().GetByYear(c, enterpriseId, fiscalMonth.FiscalYear)
	if err != nil || len(fiscalMonthYearAll) == 0 {
		return err
	}

	formatResult := make(map[int64]microWorkModel.FiscalMonthModel)
	for _, model := range fiscalMonthYearAll {
		formatResult[model.Period] = model
	}

	// 获取已存在的奖章规则（排除手动上传类型）
	medalList, err := excitationDao.NewMedal().GetMany(c, "enterprise_id = ? and source = ? and status = ?", enterpriseId, excitationModels.MedalSourceRuleStatistics, excitationModels.MedalStatusEnable)
	if err != nil {
		return err
	}
	if len(medalList) == 0 {
		return err
	}

	for _, medal := range medalList {
		// 判断当前时间是否可以执行
		startTime, endTime, ok := fiscalYearSvc.NewFiscalYear().JobIsRun(medal.MedalRule.TimeType, formatResult, fiscalMonth, yesterdayTime)
		if !ok {
			continue
		}

		// 判断是否为指定时间并且在时间范围内
		if medal.RuleType == excitationModels.MedalRuleTypeSpecify {
			if int64(medal.RuleStartTime) > yesterdayTime || int64(medal.RuleEndTime) < yesterdayTime {
				continue
			}
		}

		// 判断任务是否执行
		uniqueStr := nowTime.Format("2006-01-02") + fmt.Sprintf("-%d", medal.ID)
		exists, err := excitationDao.NewTaskRecord().GetOne(c, "enterprise_id = ? AND unique_str = ?", enterpriseId, uniqueStr)
		if exists.ID > 0 || err != nil {
			continue
		}

		// 获取奥莱所有员工
		outletsMemberIds, err := department.NewMemberDepartmentsService().GetMemberIdsByStoreType(c, enterpriseId, microWorkDao.STORE_TYPE_OUTLETS)
		job.outletsMemberIds = outletsMemberIds
		// 获取零售所有员工
		retailMemberIds, err := department.NewMemberDepartmentsService().GetMemberIdsByStoreType(c, enterpriseId, microWorkDao.STORE_TYPE_RETAIL)
		job.retailMemberIds = retailMemberIds
		var topList []models.TopList
		job.isDiffStore = false
		// 计算奖章规则维度的top人员
		switch medal.MedalRule.ShortName {
		// 获取邀请top榜单
		case excitationModels.MedalRuleInviteMonthlyList, excitationModels.MedalRuleInviteSeasonList, excitationModels.MedalRuleInviteYearList:
			job.isDiffStore = true
			topList, err = job.getInviteTopList(c, enterpriseId, medal.InviteTop, startTime, endTime)
			// 获取任务top榜单
		case excitationModels.MedalRuleTaskFinishList, excitationModels.MedalRuleTaskTouchList:
			topList, err = job.getPlanTaskTopList(c, enterpriseId, medal.InviteTop, medal.MedalRule.ShortName, startTime, endTime)
			// 获取沟通top榜单
		case excitationModels.MedalRuleResponseTimeList, excitationModels.MedalRuleFriendNumberList,
			excitationModels.MedalRuleInteractiveList, excitationModels.MedalRuleSownList:
			topList, err = job.getCommunicateTopList(c, enterpriseId, medal.InviteTop, medal.MedalRule.ShortName, startTime, endTime)
			// 销售top榜
		case excitationModels.MedalRuleSaleMonth,
			excitationModels.MedalRuleSaleSeason,
			excitationModels.MedalRuleSaleYear,
			excitationModels.MedalRuleOnlineSaleMonth,
			excitationModels.MedalRuleOnlineSaleSeason,
			excitationModels.MedalRuleOnlineSaleYear:
			job.isDiffStore = true
			topList, err = job.getSalesTopList(c, enterpriseId, medal.InviteTop, medal.MedalRule.ShortName, startTime, endTime)
			// 培训top榜
		case excitationModels.MedalRuleGoodStudent, excitationModels.MedalRuleExamManiac, excitationModels.MedalRuleStudyMember:
			topList, err = job.getTrainingList(c, enterpriseId, medal.InviteTop, medal.MedalRule.ShortName, startTime, endTime)
		case excitationModels.MedalStudyCoursesList, excitationModels.MedalStudyTimeList, excitationModels.MedalNewCoursesList:
			topList, err = job.getLearningRankList(c, enterpriseId, medal.InviteTop, medal.MedalRule.ShortName, int(fiscalMonth.Period))
		}

		if err != nil {
			app.Logger().Error(c).AnErr("计算奖章规则维度的top人员异常", err)
			continue
		}

		if len(topList) == 0 {
			app.Logger().Error(c).Msg("计算奖章规则维度的人员为空")
			continue
		}

		// 对应 medal.MedalStyle map坐标
		var timeValue int64
		var timeStr string
		// 取当前时间的月份或季度或年
		switch medal.MedalRule.TimeType {
		case excitationModels.MedalRuleTimeTypeMonth:
			monthInt := fiscalMonth.Period
			timeStr = "FY" + strconv.FormatInt(fiscalMonth.FiscalYear, 10) + "P" + strconv.FormatInt(monthInt, 10)
			timeValue = monthInt - 1
			if monthInt <= 0 {
				timeValue = 11
			}
		case excitationModels.MedalRuleTimeTypeQuarter:
			period := fiscalMonth.Period
			// 判断当前月份是第几季度
			if helper.InArray(period, []int64{1, 2, 3}) {
				timeValue = 0
			}
			if helper.InArray(period, []int64{4, 5, 6}) {
				timeValue = 1
			}
			if helper.InArray(period, []int64{7, 8, 9}) {
				timeValue = 2
			}
			if helper.InArray(period, []int64{10, 11, 12}) {
				timeValue = 3
			}
			timeStr = strconv.FormatInt(fiscalMonth.FiscalYear, 10) + "第" + strconv.FormatInt(timeValue+1, 10) + "财季"
		case excitationModels.MedalRuleTimeTypeYear:
			// 年份默认为1
			timeValue = 0
			timeStr = "FY" + strconv.FormatInt(fiscalMonth.FiscalYear, 10)
		}

		// 发放奖励
		err = job.sendRewardHandle(c, enterpriseId, topList, medal, timeValue, timeStr)
		if err != nil {
			app.Logger().Error(c).AnErr("发放奖励", err)
			continue
		}

		// 记录执行任务
		err = excitationDao.NewTaskRecord().CreateData(c, excitationModels.ExcitationTaskRecordModel{
			MicroBaseModel: models.CreateMicroBase(enterpriseId, app.UID()),
			UniqueStr:      uniqueStr,
			Type:           excitationModels.TaskTypeMedal,
		})

		if err != nil {
			app.Logger().Error(c).AnErr("记录执行任务", err)
			continue
		}
	}

	return nil
}

// 获取邀请top榜单
// inviteTop 邀请top数
func (job *SendMedal) getInviteTopList(c ctx.Context, enterpriseId int64, inviteTop int32, startTime, endTime int64) (topList []models.TopList, err error) {
	//邀约数据要往前再找一个月
	startTime = startTime - 86400*30
	endTime = endTime - 86400*30
	// 获取奥莱邀请榜名单
	memberList, err := microScrmDao.NewReservesOrder().GetTopList(c, enterpriseId, int(inviteTop), job.outletsMemberIds, startTime, endTime)
	if err != nil {
		return nil, err
	}
	if len(memberList) > 0 {
		// 按照榜单排行降序
		sort.Slice(memberList, func(i, j int) bool {
			return memberList[i].Total > memberList[j].Total
		})
		for i, info := range memberList {
			topList = append(topList, models.TopList{
				MemberId:     info.MemberId,
				Total:        info.Total,
				Ranking:      fmt.Sprintf("奥莱第%d名", i+1),
				RankingValue: fmt.Sprintf("当月邀约金额%d", info.Total/100),
			})
		}

	}

	// 获取零售邀请榜名单
	memberList, err = microScrmDao.NewReservesOrder().GetTopList(c, enterpriseId, int(inviteTop), job.retailMemberIds, startTime, endTime)
	if err != nil {
		return nil, err
	}
	if len(memberList) > 0 {
		// 按照榜单排行降序
		sort.Slice(memberList, func(i, j int) bool {
			return memberList[i].Total > memberList[j].Total
		})
		for i, info := range memberList {
			topList = append(topList, models.TopList{
				MemberId:     info.MemberId,
				Total:        info.Total,
				Ranking:      fmt.Sprintf("零售第%d名", i+1),
				RankingValue: fmt.Sprintf("当月邀约金额%d", info.Total/100),
			})
		}
	}
	if len(topList) == 0 {
		return nil, errors.New(fmt.Sprintf("邀请,获取企业%d,时间段数据为空", enterpriseId))
	}

	return
}

// 获取计划任务top榜单
func (job *SendMedal) getPlanTaskTopList(c ctx.Context, enterpriseId int64, inviteTop int32, statisticType string, startTime, endTime int64) (topList []models.TopList, err error) {
	var memberList []microScrmModel.PlanTopList
	if statisticType == excitationModels.MedalRuleTaskFinishList {
		memberList, err = microScrmDao.NewPlanTaskAllotLogDao().GetTopList(c, enterpriseId, inviteTop, startTime, endTime)
		// 按照榜单排行降序
		sort.Slice(memberList, func(i, j int) bool {
			return memberList[i].Total > memberList[j].Total
		})
		for i, info := range memberList {
			memberList[i].Ranking = fmt.Sprintf("第%d名", i+1)
			memberList[i].RankingValue = fmt.Sprintf("当月导购任务完成数量%d", info.Total)
		}
	}
	if statisticType == excitationModels.MedalRuleTaskTouchList {
		memberList, err = microScrmDao.NewPlanTaskAllotLogDao().GetTopListByTouch(c, enterpriseId, int(inviteTop), startTime, endTime)
		// 按照榜单排行降序
		sort.Slice(memberList, func(i, j int) bool {
			return memberList[i].Total > memberList[j].Total
		})
		for i, info := range memberList {
			memberList[i].Ranking = fmt.Sprintf("第%d名", i+1)
			memberList[i].RankingValue = fmt.Sprintf("当月导购沟通任务触达人数%d", info.Total)
		}
	}

	if err != nil {
		return nil, err
	}
	if len(memberList) == 0 {
		return nil, errors.New(fmt.Sprintf("计划任务,获取企业%d,时间段数据为空", enterpriseId))
	}

	var membersUserIds []string
	for _, m := range memberList {
		membersUserIds = append(membersUserIds, m.MemberUserId)
	}

	members, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{"enterprise_id": enterpriseId, "userid": membersUserIds})
	if err != nil {
		return nil, err
	}

	var memberMap = make(map[string]int64)
	for _, member := range members {
		memberMap[member.Userid] = member.ID
	}
	for _, info := range memberList {
		if memberId, ok := memberMap[info.MemberUserId]; ok {
			topList = append(topList, models.TopList{
				MemberId:     memberId,
				Total:        info.Total,
				Ranking:      info.Ranking,
				RankingValue: info.RankingValue,
			})
		}
	}

	return
}

// 获取沟通top榜单
// inviteTop 邀请top数
func (job *SendMedal) getCommunicateTopList(c ctx.Context, enterpriseId int64, inviteTop int32, statisticType string, startTime, endTime int64) (resTopList []models.TopList, err error) {
	var topList []models.TopList
	var memberTopList []microScrmModel.TopList
	switch statisticType {
	case excitationModels.MedalRuleResponseTimeList:
		memberTopList, err = microScrmDao.NewMemberExternalBehaviorData().GetAverageReplyTimeTopList(c, enterpriseId, int(inviteTop), startTime, endTime)
		if err != nil {
			return nil, err
		}
		if len(memberTopList) == 0 {
			return nil, errors.New(fmt.Sprintf("获取沟通top榜单,获取企业%d,时间段数据为空", enterpriseId))
		}
		// 按照榜单排行降序
		sort.Slice(memberTopList, func(i, j int) bool {
			return memberTopList[i].Total < memberTopList[j].Total
		})
		for i, info := range memberTopList {
			resTopList = append(resTopList, models.TopList{
				MemberId:     info.MemberId,
				Total:        int64(info.Total),
				Ranking:      fmt.Sprintf("第%d名", i+1),
				RankingValue: fmt.Sprintf("当月平均回复时长%f秒", info.Total),
			})
		}

	case excitationModels.MedalRuleFriendNumberList:
		topList, err = microScrmDao.NewMemberExternalBehaviorData().GetContactCountTopList(c, enterpriseId, int(inviteTop), startTime, endTime)
		// 按照榜单排序
		sort.Slice(topList, func(i, j int) bool {
			return topList[i].Total > topList[j].Total
		})
		for i, info := range topList {
			resTopList = append(resTopList, models.TopList{
				MemberId:     info.MemberId,
				Total:        info.Total,
				Ranking:      fmt.Sprintf("第%d名", i+1),
				RankingValue: fmt.Sprintf("当月新增客户数%d", info.Total),
			})
		}
	case excitationModels.MedalRuleInteractiveList:
		topList, err = microScrmDao.NewMemberExternalBehaviorData().GetChatCountTopList(c, enterpriseId, int(inviteTop), startTime, endTime)
		// 按照榜单排序
		sort.Slice(topList, func(i, j int) bool {
			return topList[i].Total > topList[j].Total
		})
		for i, info := range topList {
			resTopList = append(resTopList, models.TopList{
				MemberId:     info.MemberId,
				Total:        info.Total,
				Ranking:      fmt.Sprintf("第%d名", i+1),
				RankingValue: fmt.Sprintf("当月聊天总数%d", info.Total),
			})
		}
	case excitationModels.MedalRuleSownList:
		topList, err = customerDao.NewSowTask().GetTopList(c, enterpriseId, int(inviteTop), startTime, endTime)
		// 按照榜单排序
		sort.Slice(topList, func(i, j int) bool {
			return topList[i].Total > topList[j].Total
		})
		for i, info := range topList {
			resTopList = append(resTopList, models.TopList{
				MemberId:     info.MemberId,
				Total:        info.Total,
				Ranking:      fmt.Sprintf("第%d名", i+1),
				RankingValue: fmt.Sprintf("当月播种总数%d", info.Total),
			})
		}
	}

	if err != nil {
		return nil, err
	}
	if len(topList) == 0 {
		return nil, errors.New(fmt.Sprintf("获取沟通top榜单,获取企业%d,时间段数据为空", enterpriseId))
	}
	return
}

// SalesTopList 销售top榜
// inviteTop top数
func (job *SendMedal) getSalesTopList(c ctx.Context, enterpriseId int64, inviteTop int32, statisticType string, startTime, endTime int64) (topList []models.TopList, err error) {

	memberMap := make(map[int64]float64)
	type MemberTotal struct {
		memberId int64
		total    float64
	}

	if helper.InArray(statisticType, []string{excitationModels.MedalRuleSaleMonth, excitationModels.MedalRuleSaleSeason, excitationModels.MedalRuleSaleYear}) {
		// 获取导购线上销售金额
		memberMap, err = microScrmDao.NewMemberSalesVolume().GetMemberSalesVolume(c, memberReq.MemberSalesVolumeReq{
			EnterpriseId: enterpriseId,
			StartTime:    startTime,
			EndTime:      endTime,
			SearchType:   microScrmModel.ALL_SALES_VOLUME,
		})
	}

	if helper.InArray(statisticType, []string{excitationModels.MedalRuleOnlineSaleMonth, excitationModels.MedalRuleOnlineSaleSeason, excitationModels.MedalRuleOnlineSaleYear}) {
		// 获取导购销售总金额
		memberMap, err = microScrmDao.NewMemberSalesVolume().GetMemberSalesVolume(c, memberReq.MemberSalesVolumeReq{
			EnterpriseId: enterpriseId,
			StartTime:    startTime,
			EndTime:      endTime,
			SearchType:   microScrmModel.ONLINE_SALES_VOLUME,
		})
	}

	if err != nil {
		return
	}

	// 分别获取零售和奥莱的top
	outletsTopList := make([]MemberTotal, 0)
	retailTopList := make([]MemberTotal, 0)
	for k, v := range memberMap {
		if helper.InArray(k, job.outletsMemberIds) {
			outletsTopList = append(outletsTopList, MemberTotal{memberId: k, total: v})
		}
		if helper.InArray(k, job.retailMemberIds) {
			retailTopList = append(retailTopList, MemberTotal{memberId: k, total: v})
		}
	}

	// 按总数从大到小重新排序
	sort.Slice(outletsTopList, func(i, j int) bool {
		return outletsTopList[i].total > outletsTopList[j].total
	})
	sort.Slice(retailTopList, func(i, j int) bool {
		return retailTopList[i].total > retailTopList[j].total
	})

	//获取top以及和top最后一位数量相同的其他导购
	var breakTotal float64
	var i int
	for key, member := range outletsTopList {
		i++
		if key == int(inviteTop)-1 {
			breakTotal = member.total //获取最后一位的数量
		}
		if key >= int(inviteTop) {
			if breakTotal > 0 && breakTotal > member.total {
				break
			}
		}

		topList = append(topList, models.TopList{
			MemberId:     member.memberId,
			Total:        int64(member.total),
			Ranking:      fmt.Sprintf("奥莱第%d名", i),
			RankingValue: fmt.Sprintf("当月销售金额%.2f", member.total),
		})
	}
	breakTotal = 0
	i = 0
	for key, member := range retailTopList {
		i++
		if key == int(inviteTop)-1 {
			breakTotal = member.total //获取最后一位的数量
		}
		if key >= int(inviteTop) {
			if breakTotal > 0 && breakTotal > member.total {
				break
			}
		}
		topList = append(topList, models.TopList{
			MemberId:     member.memberId,
			Total:        int64(member.total),
			Ranking:      fmt.Sprintf("零售第%d名", i),
			RankingValue: fmt.Sprintf("当月销售金额%.2f", member.total),
		})
	}

	return
}

// getTrainingList 培训top榜
// inviteTop top数
func (job *SendMedal) getTrainingList(c ctx.Context, enterpriseId int64, inviteTop int32, statisticType string, startTime, endTime int64) (topList []models.TopList, err error) {
	var memberList []excitationResp.MemberSumPointsList
	switch statisticType {
	case excitationModels.MedalRuleGoodStudent:
		memberList, err = excitationDao.NewPointsRecord().GetTopListByRemark(c, enterpriseId, excitationModels.RuleCategoryId15, int(inviteTop), startTime, endTime, nil)
	case excitationModels.MedalRuleExamManiac:
		memberList, err = excitationDao.NewPointsRecord().GetTopListByRemark(c, enterpriseId, excitationModels.RuleCategoryId15, int(inviteTop), startTime, endTime, []string{"考试积分"})
	case excitationModels.MedalRuleStudyMember:
		memberList, err = excitationDao.NewPointsRecord().GetTopListByRemark(c, enterpriseId, excitationModels.RuleCategoryId15, int(inviteTop), startTime, endTime, []string{"内容学习"})
	}
	// 按照榜单排行降序
	sort.Slice(memberList, func(i, j int) bool {
		return memberList[i].SumPoints > memberList[j].SumPoints
	})
	for i, info := range memberList {
		topList = append(topList, models.TopList{
			MemberId:     info.MemberId,
			Total:        int64(info.SumPoints),
			Ranking:      fmt.Sprintf("第%d名", i+1),
			RankingValue: fmt.Sprintf("获得WeLearning积分%f", info.SumPoints),
		})
	}
	return
}

// getTrainingList 培训top榜
// inviteTop top数
func (job *SendMedal) getLearningRankList(c ctx.Context, enterpriseId int64, inviteTop int32, statisticType string, pMonth int) (topList []models.TopList, err error) {
	//获取自然月的月份

	rankList, err := microScrmDao.NewMemberLearnRank().GetTopList(c, enterpriseId, statisticType, pMonth, int(inviteTop))
	if err != nil {
		return
	}
	var typeStr string
	if statisticType == excitationModels.MedalStudyCoursesList {
		typeStr = "当月学习课程数量"
	} else if statisticType == excitationModels.MedalStudyTimeList {
		typeStr = "当月平台课程学习时长"
	} else if statisticType == excitationModels.MedalNewCoursesList {
		typeStr = "当月平台新品课程完成率"
	}

	for i, info := range rankList {
		topList = append(topList, models.TopList{
			MemberId:     info.MemberID,
			Total:        int64(helper.StringToInt(info.Value)),
			Ranking:      fmt.Sprintf("第%d名", i+1),
			RankingValue: fmt.Sprintf("%s%s", typeStr, info.Value),
		})
	}
	//消费数据
	err = microScrmDao.NewMemberLearnRank().ConsumeData(c, enterpriseId, statisticType, pMonth)
	if err != nil {
		return nil, err
	}
	return
}

// 发送奖励逻辑
func (job *SendMedal) sendRewardHandle(c ctx.Context, enterpriseId int64, topList []models.TopList, medal excitationModels.ExcitationMedalModel, timeValue int64, timeStr string) error {
	var result []excitationModels.ExcitationMedalRecordModel
	for _, topInfo := range topList {
		result = append(result, excitationModels.ExcitationMedalRecordModel{
			MicroBaseModel: models.MicroBaseModel{
				BaseModel: models.BaseModel{
					ID:        app.UIDContext(c),
					CreatedAt: job.nowTime,
					UpdatedAt: job.nowTime,
				},
				EnterpriseID: enterpriseId,
			},
			MedalID:      medal.ID,
			MedalStyleID: medal.MedalStyle[timeValue].ID,
			MemberID:     topInfo.MemberId,
			Ranking:      topInfo.Ranking,
			RankingValue: topInfo.RankingValue,
		})
	}
	if medal.Integral > 0 {
		var addMedalPoints []excitationReq.ManualPointsSave
		for _, medalRecord := range result {
			memberTypeStr := ""
			if job.isDiffStore {
				if helper.InArray(medalRecord.MemberID, job.outletsMemberIds) {
					memberTypeStr = "(奥莱)"
				} else if helper.InArray(medalRecord.MemberID, job.retailMemberIds) {
					memberTypeStr = "(零售)"
				}
			}
			tmp := excitationReq.ManualPointsSave{
				MemberId:      medalRecord.MemberID,
				RuleName:      medal.Name,
				RuleDetailId:  0,
				PointsRemark:  timeStr + "|" + fmt.Sprintf("奖章-%s-TOP%d", medal.MedalStyle[timeValue].MedalName, medal.InviteTop) + memberTypeStr,
				Points:        int64(medal.Integral),
				NowTime:       job.nowTime,
				MedalRecordId: medalRecord.ID, //奖章记录组装
			}
			addMedalPoints = append(addMedalPoints, tmp)
		}

		err := excitationSvc.NewPointsProduceService().MedalPoints(c, enterpriseId, addMedalPoints)
		if err != nil {
			return err
		}
	}
	err := excitationDao.NewMedalRecord().CreateInBatchesData(c, result)

	return err
}
