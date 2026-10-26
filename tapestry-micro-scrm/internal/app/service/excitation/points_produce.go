package excitationSvc

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/http/request"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	memberReq "micro-scrm/internal/app/http/request/member"
	reserveReq "micro-scrm/internal/app/http/request/reserve"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/app/service/department"
	"micro-scrm/internal/pkg/helper"
	pkgKafka "micro-scrm/internal/pkg/kafka"
	"micro-scrm/internal/pkg/sdk/microwork"
	"sort"
	"strconv"
	"strings"
	"time"
)

type PointsProduceService struct {
	service.Base
	AchievementDao       *excitationDao.Achievement
	PointsRecordDao      *excitationDao.PointsRecord
	IntegralRuleDao      *excitationDao.IntegralRuleDetail
	memberSalesVolumeDao *microScrmDao.MemberSalesVolume
	storeTypeMap         map[int64]string
}

func NewPointsProduceService() *PointsProduceService {
	return &PointsProduceService{
		AchievementDao:       excitationDao.NewAchievement(),
		PointsRecordDao:      excitationDao.NewPointsRecord(),
		IntegralRuleDao:      excitationDao.NewIntegralRuleDetail(),
		memberSalesVolumeDao: microScrmDao.NewMemberSalesVolume(),
		storeTypeMap: map[int64]string{
			0:                               "",
			microWorkDao.STORE_TYPE_OUTLETS: "奥莱",
			microWorkDao.STORE_TYPE_RETAIL:  "零售",
		},
	}
}

// MedalPoints 奖章添加积分
func (s *PointsProduceService) MedalPoints(c ctx.Context, enterpriseId int64, dataList []excitationReq.ManualPointsSave) error {
	var recordData []microScrmModel.ExcitationPointsRecordModel
	memberDataSvc := NewMemberDataService()
	for i, data := range dataList {
		if data.Points == 0 {
			// 如果积分为0 则直接跳过
			continue
		}
		pointsRemark := data.PointsRemark
		var fyTimeStr string
		strArr := strings.Split(data.PointsRemark, "|")
		if len(strArr) > 1 {
			fyTimeStr = strArr[0]
			pointsRemark = strArr[1]
		}
		tmp := microScrmModel.ExcitationPointsRecordModel{
			BaseModel: models.BaseModel{
				ID:        app.UID(),
				CreatedAt: data.NowTime,
				UpdatedAt: data.NowTime,
			},
			MemberID:         data.MemberId,
			IntegralRuleName: data.RuleName,
			Points:           strconv.FormatInt(data.Points, 10),
			SourceType:       microScrmModel.ExcitationPointsRecordAutoType_Medal,
			PointsRemark:     pointsRemark,
			FyTimeStr:        fyTimeStr,
		}
		recordData = append(recordData, tmp)
		err := memberDataSvc.savePoints(c, excitationReq.UpdatePointsReq{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: enterpriseId,
			},
			MemberId: data.MemberId,
			Points:   float64(data.Points),
		})
		if err != nil {
			return err
		}
		//添加积分数据组装
		dataList[i].PointsRecordId = tmp.ID
		//奖章增加积分均增加金币
		dataList[i].HasCoins = true
	}
	err := s.PointsRecordDao.CreateInBatchesData(c, enterpriseId, recordData)
	if err != nil {
		return err
	}
	s.SaveCoins(c, enterpriseId, microScrmModel.ExcitationCoinsRecordAutoType_Points, dataList)

	return nil
}

// ManualPoints 手动添加积分
func (s *PointsProduceService) ManualPoints(c ctx.Context, enterpriseId int64, dataList []excitationReq.ManualPointsSave) error {
	ruleList, err := s.IntegralRuleDao.GetRule(c, "enterprise_id = ?", enterpriseId)
	if err != nil {
		return err
	}
	var recordData []microScrmModel.ExcitationPointsRecordModel
	memberDataSvc := NewMemberDataService()
	for i, data := range dataList {
		if data.Points == 0 {
			// 如果积分为0 则直接跳过
			continue
		}
		tmp := microScrmModel.ExcitationPointsRecordModel{
			BaseModel:              models.CreateBase(app.UID()),
			MemberID:               data.MemberId,
			IntegralRuleName:       data.RuleName,
			IntegralRuleDetailName: data.RuleDetailName,
			Points:                 strconv.FormatInt(data.Points, 10),
			SourceType:             microScrmModel.ExcitationPointsRecordAutoType_Manual,
			PointsRemark:           data.PointsRemark,
		}
		for _, ruleDetail := range ruleList {
			if ruleDetail.RuleID == data.RuleId {
				tmp.IncentiveType = ruleDetail.Rule.Type
				tmp.IntegralRuleName = ruleDetail.Rule.RuleName
				tmp.IntegralRuleID = data.RuleId
			}
		}
		recordData = append(recordData, tmp)
		err := memberDataSvc.savePoints(c, excitationReq.UpdatePointsReq{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: enterpriseId,
			},
			MemberId: data.MemberId,
			Points:   float64(data.Points),
		})
		if err != nil {
			return err
		}
		//添加积分数据组装
		dataList[i].PointsRecordId = tmp.ID
	}
	err = s.PointsRecordDao.CreateInBatchesData(c, enterpriseId, recordData)
	if err != nil {
		return err
	}
	//处理加金币(仍然算积分加的金币)
	s.SaveCoins(c, enterpriseId, microScrmModel.ExcitationCoinsRecordAutoType_Points, dataList)
	return nil
}

// ProduceQueue  队列生产写入（用于批量插入且需要排队的情况）
func (s *PointsProduceService) ProduceQueue(c ctx.Context, params excitationReq.PointsProduceSave) {
	str, _ := json.Marshal(params)
	err := pkgKafka.GetAppKafka().Producer(string(str), "points_produce")
	if err != nil {
		log.Info(c).Msgf("新增积分数据累计队列失败,积分类型:%d,错误原因：%+v", params.Type, err)
	}
}

// Produce 用户数据积累
func (s *PointsProduceService) Produce(c ctx.Context, unixTime int64, params excitationReq.PointsProduceSave) error {
	if params.MemberId == 0 && params.Mobile != "" {
		// 兼容没有memberId但是有mobile的情况
		memberInfo, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{"enterprise_id": params.EnterpriseId, "mobile": params.Mobile})
		if err != nil {
			return err
		}
		params.MemberId = memberInfo.ID
	}
	if params.Value == 0 {
		log.Info(c).Msgf("积分获取跳出，原因：数据为0")
		// 如果数据为0则无需记录
		return nil
	}
	unique, err := s.AchievementDao.GetOne(c, "unique_hash = ? AND enterprise_id = ?", params.UniqueHash, params.EnterpriseId)
	if err != nil {
		return err
	}
	if unique.ID > 0 {
		log.Info(c).Msgf("积分获取跳出，原因：achievementHash存在")
		return nil
	}

	// 判断是否需要积累数据,未配置规则时不累计
	refuse, err := s.checkRule(c, params)
	if err != nil {
		return err
	}
	if refuse {
		log.Info(c).Msgf("积分获取跳出，原因：未配置规则")
		// 没有配置规则直接跳过
		return nil
	}
	// 用户数据积累
	err = s.saveAchievement(c, &params)
	if err != nil {
		return err
	}
	// 消费用户数据积累获得积分
	err = s.consumePointsRecord(c, unixTime, params)
	if err != nil {
		return err
	}
	return nil
}

// ConsumeTest 消费用户数据积累
func (s *PointsProduceService) ConsumeTest(c ctx.Context, unixTime int64, id int64, enterpriseId int64) error {
	achievementData, err := s.AchievementDao.GetOne(c, enterpriseId, "id = ?", id)
	// 消费用户数据积累获得积分
	err = s.consumePointsRecord(c, unixTime, excitationReq.PointsProduceSave{
		MemberId:     achievementData.MemberID,
		Type:         achievementData.Type,
		Value:        achievementData.Value,
		UniqueHash:   achievementData.UniqueHash,
		EnterpriseId: achievementData.EnterpriseID,
		SourceRemark: achievementData.SourceRemark,
	})
	if err != nil {
		return err
	}
	return nil
}

// 消费用户数据积累获得积分
func (s *PointsProduceService) consumePointsRecord(c ctx.Context, unixTime int64, params excitationReq.PointsProduceSave) error {
	// 查询消费条件
	list, err := s.IntegralRuleDao.GetRule(c, "enterprise_id = ? AND rule_category_id = ?", params.EnterpriseId, params.Type)
	if err != nil {
		return err
	}
	var trueRule []excitationModels.ExcitationIntegralRuleDetailModel
	var average, memberStoreType int64 // 标记导购平均值的奥莱零售类型
	for _, rule := range list {
		storeType := int64(0) // 默认通用规则
		if rule.IsDiffStore == excitationModels.IsDiffStore_Yes {
			// 如果不是通用规则需要根据当前导购属于奥莱或者零售来过滤规则
			storeType, err = department.NewMemberDepartmentsService().GetMemberStoreType(c, params.EnterpriseId, params.MemberId)
			if err != nil || storeType == 0 {
				// 获取该导购奥莱或者零售类型失败 则直接结束
				return err
			}
		}
		if rule.StoreType != storeType {
			// 不符合规则类型的直接跳过
			continue
		}
		if rule.ValidityType == excitationModels.ValidityTypeDeadline && (rule.StartTime > unixTime || rule.EndTime < unixTime) {
			// 如果是限时规则，检查时间是否合法，不合法则忽略
			continue
		}
		if rule.RuleIf1 == excitationModels.RuleIf1Average {
			// 如果条件是平均值,获取平均值
			average, memberStoreType, err = s.GetAverage(c, unixTime, excitationReq.GetAverageReq{
				EnterpriseId:   params.EnterpriseId,
				MemberId:       params.MemberId,
				RuleCategoryId: params.Type,
			})
			if err != nil {
				return err
			}
			if average == -1 {
				// 如果用户不属于奥莱或者零售，平均值条件忽略掉
				continue
			}
			rule.RuleValue = average
		}
		trueRule = append(trueRule, rule)
	}
	// 按照规则梯度条件降序(达到难度高的排在前面)
	sort.Slice(trueRule, func(i, j int) bool {
		return trueRule[i].RuleValue > trueRule[j].RuleValue
	})
	strMap := map[int64]string{
		excitationModels.RuleIfGT:  "高于",
		excitationModels.RuleIfEgt: "达到",
	}
	remarkName := s.getPointsConditionStr(params.Type).Name
	for _, rule := range trueRule {
		if params.Value > rule.RuleValue || // 条件满足
			params.Type == excitationModels.RuleCategoryId15 || // 积分等比兑换类不考虑满足与否
			(params.Value == rule.RuleValue && (rule.RuleIf == excitationModels.RuleIfTimes || rule.RuleIf == excitationModels.RuleIfEgt)) {
			// 用户积累数据更新
			updateAchievementData := map[string]interface{}{
				"updated_at": unixTime,
				"value":      0, // 加积分直接消费掉数据
			}
			remark := fmt.Sprintf(s.getPointsConditionStr(params.Type).Name, rule.RuleValue)
			if strings.Contains(remarkName, "%s") {
				ruleIfStr := ""
				if value, ok := strMap[rule.RuleIf]; ok {
					ruleIfStr = value
				}
				remark = fmt.Sprintf(s.getPointsConditionStr(params.Type).Name, ruleIfStr, rule.RuleValue)
				if rule.RuleIf1 == excitationModels.RuleIf1Average {
					remark += "(" + s.storeTypeMap[memberStoreType] + "平均值)"
				}
			}
			// 如果任务值满足条件则消费(次数和达到都是大于等于，只有高于是大于)
			recordData := microScrmModel.ExcitationPointsRecordModel{
				BaseModel: models.BaseModel{
					ID:        app.UID(),
					CreatedAt: unixTime,
					UpdatedAt: unixTime,
				},
				IncentiveType:    rule.Rule.Type,
				MemberID:         params.MemberId,
				Type:             params.Type,
				Value:            params.Value,
				ConditionValue:   rule.RuleIf1,
				SourceType:       microScrmModel.ExcitationPointsRecordAutoType_Auto,
				IntegralRuleID:   rule.ID,
				IntegralRuleName: rule.Rule.RuleName,
				PointsRemark:     remark,
			}

			points := float64(rule.Integral) // 实际积分计算结果
			// 在次数的条件下 有可能出现需要增加倍数积分的情况（当前业务场景不存在倍数情况20230620）
			if rule.RuleIf == excitationModels.RuleIfTimes {
				// 添加积分的条件可能无法消耗掉所有的次数类型的积累（只有积分等比兑换的情况下必须全部消耗，可以有小数）
				if params.Type == excitationModels.RuleCategoryId15 {
					// 四舍五入保留1位小数
					points = helper.MathRounded(float64(params.Value)/float64(rule.RuleValue)*float64(rule.Integral), 1)
					// 当前业务场景不存在倍数情况20230620
					// } else {
					// multiple := params.Value / rule.RuleValue // 两个int64相除 默认向下取整（计算倍数积分）
					// intPoints := multiple * rule.Integral
					// updateAchievementData["value"] = params.Value - multiple*rule.RuleValue // 用户数据消耗也和获得积分的倍数相同
					// points = float64(intPoints)
					// // 处理添加积分记录的文字说明
					// recordData.PointsRemark = fmt.Sprintf(s.getPointsConditionStr(params.Type).Name, rule.RuleValue*multiple)
					// if strings.Contains(remarkName, "%s") {
					//	ruleIfStr := ""
					//	if value, ok := strMap[rule.RuleIf]; ok {
					//		ruleIfStr = value
					//	}
					//	remark = fmt.Sprintf(s.getPointsConditionStr(params.Type).Name, ruleIfStr, rule.RuleValue*multiple)
					// }
				}
			}
			// 减积分的情况
			if rule.IntegralType == excitationModels.IntegralTypeDec {
				points = 0 - points
			}
			recordData.Points = fmt.Sprintf("%0.1f", points)
			if params.SourceRemark != "" {
				// 特殊渠道数据需要记录来源
				recordData.PointsRemark = params.SourceRemark
			}
			// 处理超限逻辑
			astrict, err := s.checkAstrict(c, excitationReq.CheckAstrictReq{
				PointsProduceSave: params,
				NowTime:           unixTime,
				AddPoints:         points,
			}, rule)
			if err != nil {
				return err
			}
			updateAchievementDataCon := fmt.Sprintf("enterprise_id = %d AND member_id = %d AND type = %d", params.EnterpriseId, params.MemberId, params.Type)
			if astrict {
				// 达到月上限直接清空累计并跳过
				log.Info(c).Msgf("积分获取跳出，原因：达到积分规则上限")
				err = excitationDao.NewAchievement().UpdateData(c, updateAchievementDataCon, map[string]interface{}{
					"updated_at": unixTime,
					"value":      0, // 超限直接清空累计
				})
				return nil
			}
			err = s.Base.DB().Transaction(func(tx *gorm.DB) error {
				// 创建积分记录
				err := excitationDao.NewPointsRecord().WithTx(tx).CreateData(c, params.EnterpriseId, recordData)
				if err != nil {
					return err
				}
				// 更新用户数据积累
				err = excitationDao.NewAchievement().WithTx(tx).UpdateData(c, updateAchievementDataCon, updateAchievementData)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
			// 更新导购积分
			err = NewMemberDataService().savePoints(c, excitationReq.UpdatePointsReq{
				EnterpriseBaseReq: request.EnterpriseBaseReq{
					EnterpriseId: params.EnterpriseId,
				},
				MemberId: recordData.MemberID,
				Points:   points,
			})
			if err != nil {
				return err
			}
			//处理金币变更
			coinsParams := excitationReq.ManualPointsSave{
				MemberId:       recordData.MemberID,
				PointsRecordId: recordData.ID,
				Points:         int64(points),
				PointsRemark:   recordData.PointsRemark,
				HasCoins:       true,
			}
			s.SaveCoins(c, params.EnterpriseId, microScrmModel.ExcitationCoinsRecordAutoType_Points, []excitationReq.ManualPointsSave{coinsParams})
			break
		}
	}
	return nil
}

// 用户数据积累
func (s *PointsProduceService) saveAchievement(c ctx.Context, params *excitationReq.PointsProduceSave) error {
	// 查询是否存在需更新的数据
	achievementData, err := s.AchievementDao.GetOne(c, "enterprise_id = ? AND member_id = ? AND type = ?", params.EnterpriseId, params.MemberId, params.Type)
	if err != nil {
		return err
	}
	err = s.Base.DB().Transaction(func(tx *gorm.DB) error {
		excitationTransaction := excitationDao.NewAchievement()
		if achievementData.ID > 0 {
			newValue := params.Value
			if s.getPointsConditionStr(achievementData.Type).Cumulative {
				// 存在未被积分条件消费的数据且属于需要累计数据的类型则更新该数据
				newValue = achievementData.Value + params.Value
			}
			updateDataMap := map[string]interface{}{
				"unique_hash":   params.UniqueHash,
				"value":         newValue,
				"updated_at":    time.Now().Unix(),
				"source_remark": params.SourceRemark,
			}
			err := excitationTransaction.WithTx(tx).UpdateData(c, fmt.Sprintf("id = %d", achievementData.ID), updateDataMap)
			if err != nil {
				return err
			}
			params.Value = newValue
		} else {
			// 不存在未被积分条件消费的数据,新增一条
			achievementData = excitationModels.ExcitationAchievementModel{
				MicroBaseModel: models.CreateMicroBase(params.EnterpriseId, app.UID()),
				MemberID:       params.MemberId,
				Type:           params.Type,
				Value:          params.Value,
				UniqueHash:     params.UniqueHash,
				SourceRemark:   params.SourceRemark,
			}
			err := excitationTransaction.WithTx(tx).CreateData(c, achievementData)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 获取积分记录说明
func (s *PointsProduceService) getPointsConditionStr(recordType int64) excitationReq.RuleCategory {
	ruleStr := make(map[int64]excitationReq.RuleCategory)
	ruleStr[excitationModels.RuleCategoryId1] = excitationReq.RuleCategory{
		Name:       "创建邀约%d次",
		Rate:       "次",
		Cumulative: true,
	}
	ruleStr[excitationModels.RuleCategoryId2] = excitationReq.RuleCategory{
		Name: "邀约到店率(周) %s %d%%",
		Rate: "周",
	}
	ruleStr[excitationModels.RuleCategoryId3] = excitationReq.RuleCategory{
		Name: "邀约成功率(月) %s %d%%",
		Rate: "月",
	}
	ruleStr[excitationModels.RuleCategoryId4] = excitationReq.RuleCategory{
		Name:       "邀约成功%d次",
		Rate:       "次",
		Cumulative: true,
	}
	ruleStr[excitationModels.RuleCategoryId5] = excitationReq.RuleCategory{
		Name: "邀约占比(月) %s %d%%",
		Rate: "月",
	}
	ruleStr[excitationModels.RuleCategoryId6] = excitationReq.RuleCategory{
		Name:       "邀约到店确认%d次",
		Rate:       "次",
		Cumulative: true,
	}
	ruleStr[excitationModels.RuleCategoryId7] = excitationReq.RuleCategory{
		Name:       "任务完成%d次",
		Rate:       "次",
		Cumulative: true,
	}
	ruleStr[excitationModels.RuleCategoryId8] = excitationReq.RuleCategory{
		Name: "任务完成率(月) %s %d%%",
		Rate: "月",
	}
	ruleStr[excitationModels.RuleCategoryId9] = excitationReq.RuleCategory{
		Name:       "任务触达%d人",
		Rate:       "人",
		Cumulative: true,
	}
	ruleStr[excitationModels.RuleCategoryId10] = excitationReq.RuleCategory{
		Name:       "新增好友%d人",
		Rate:       "次",
		Cumulative: true,
	}
	ruleStr[excitationModels.RuleCategoryId11] = excitationReq.RuleCategory{
		Name: "客户流失%d人",
		Rate: "人",
	}
	ruleStr[excitationModels.RuleCategoryId12] = excitationReq.RuleCategory{
		Name: "聊天总数(月) %s %d",
		Rate: "月",
	}
	ruleStr[excitationModels.RuleCategoryId13] = excitationReq.RuleCategory{
		Name: "已回复聊天占比(月) %s %d%%",
		Rate: "月",
	}
	ruleStr[excitationModels.RuleCategoryId14] = excitationReq.RuleCategory{
		Name:       "完成播种%d次",
		Rate:       "次",
		Cumulative: true,
	}
	ruleStr[excitationModels.RuleCategoryId15] = excitationReq.RuleCategory{
		Name: "积分等比兑换(积分) %d分",
		Rate: "次",
	}
	ruleStr[excitationModels.RuleCategoryId16] = excitationReq.RuleCategory{
		Name: "线上销售业绩率(月) %s %d%%",
		Rate: "次",
	}
	ruleStr[excitationModels.RuleCategoryId17] = excitationReq.RuleCategory{
		Name:       "点击新老客按钮%d次",
		Rate:       "次",
		Cumulative: true,
	}
	ruleStr[excitationModels.RuleCategoryId18] = excitationReq.RuleCategory{
		Name:       "邀约是否到店购买点击%d次",
		Rate:       "次",
		Cumulative: true,
	}
	ruleStr[excitationModels.RuleCategoryId19] = excitationReq.RuleCategory{
		Name:       "线上销售金额(月) %s %d",
		Rate:       "次",
		Cumulative: true,
	}

	if value, ok := ruleStr[recordType]; ok {
		return value
	}
	return excitationReq.RuleCategory{}
}

// GetAverage 获取各种统计的平均值
func (s *PointsProduceService) GetAverage(c ctx.Context, unixTime int64, averageReq excitationReq.GetAverageReq) (average, memberStoreType int64, err error) {
	for storeType := range s.storeTypeMap {
		memberIds, err := department.NewMemberDepartmentsService().GetMemberIdsByStoreType(c, averageReq.EnterpriseId, storeType)
		if err != nil {
			return 0, storeType, err
		}
		// 判断用户是否是奥莱用户或者零售用户，获取相关平均值
		if helper.InArray(averageReq.MemberId, memberIds) {
			// 获取平均值缓存
			averageRedisKey := fmt.Sprintf("excitation:average:%d-%d-%d-%d", averageReq.EnterpriseId, averageReq.RuleCategoryId, storeType, unixTime)
			averageStr := s.Redis().Get(c, averageRedisKey).Val()
			if averageStr != "" {
				average, err = strconv.ParseInt(averageStr, 10, 64)
				if err != nil {
					return 0, storeType, err
				}
				return average, storeType, nil
			}
			// 没有平均值缓存则查库
			rate := s.getPointsConditionStr(averageReq.RuleCategoryId).Rate // 周 月
			var startTime, endTime int64
			// 需要执行的数据查询统计都是今天之前的数据
			yesterday := time.Unix(unixTime, 0).AddDate(0, 0, -1)
			switch rate {
			case "周":
				// 财周的结束时间为周日0点
				endTime = helper.GetLastDateOfWeek(helper.GetZeroTimeByTime(yesterday)).Unix()
				startTime = endTime - 86400*7
			case "月":
				// 财月的统计时间根据参数时间查询开始和结束时间
				fiscalInfo, _ := microWorkDao.NewFiscalMonth().GetFiscalMonth(c, averageReq.EnterpriseId, yesterday.Unix())
				startTime = fiscalInfo.StartTime
				endTime = fiscalInfo.EndTime
				if averageReq.RuleCategoryId == excitationModels.RuleCategoryId3 {
					// 邀约成功率的统计时间范围需要再往前多找一个月（由于邀约的成功匹配有最多一个月的延迟）
					fiscalInfo, _ = microWorkDao.NewFiscalMonth().GetFiscalMonth(c, averageReq.EnterpriseId, startTime-86400)
					startTime = fiscalInfo.StartTime
					endTime = fiscalInfo.EndTime
				}
			}
			achievementParams := reserveReq.AchievementParams{
				EnterpriseId: averageReq.EnterpriseId,
				MemberIds:    memberIds,
				Type:         averageReq.RuleCategoryId,
				StartTime:    startTime,
				EndTime:      endTime,
			}
			// 平均值查库
			switch averageReq.RuleCategoryId {
			case excitationModels.RuleCategoryId2,
				excitationModels.RuleCategoryId3,
				excitationModels.RuleCategoryId5:
				res, err := microScrmDao.NewReservesOrder().GetAverage(c, achievementParams)
				if err != nil {
					return 0, storeType, err
				}
				average = res
				if averageReq.RuleCategoryId == excitationModels.RuleCategoryId5 {
					salesTotal, err := s.memberSalesVolumeDao.GetAllVolumeByMemberIds(c, memberReq.AchievementParams{
						EnterpriseId: averageReq.EnterpriseId,
						MemberIds:    memberIds,
						SearchType:   microScrmModel.ALL_SALES_VOLUME,
						StartTime:    startTime,
						EndTime:      endTime,
					})
					if err != nil {
						return 0, storeType, err
					}
					average = 0
					if res > 0 && salesTotal > 0 {
						average = res / salesTotal
					}

				}
			case excitationModels.RuleCategoryId12,
				excitationModels.RuleCategoryId13:
				average, err = microScrmDao.NewMemberExternalBehaviorData().GetAverage(c, achievementParams)
				if err != nil {
					return 0, storeType, err
				}
			}
			// 平均值写缓存
			s.Redis().Set(c, averageRedisKey, average, time.Hour*12)
			return average, storeType, nil
		}
	}

	// 如果没有查到
	return -1, 0, nil
}

// 检查是否配置了规则
func (s *PointsProduceService) checkRule(c ctx.Context, params excitationReq.PointsProduceSave) (bool, error) {
	var ruleWhere = map[string]interface{}{
		"enterprise_id":    params.EnterpriseId,
		"rule_category_id": params.Type,
	}
	ruleInfo, _ := s.IntegralRuleDao.GetOne(c, ruleWhere)
	if ruleInfo.ID == 0 {
		// 没有配置该分类
		log.Info(c).Msgf("该企业没有配置此维度的获取积分业务，企业id:%d,分类id：%d", params.EnterpriseId, params.Type)
	}
	return false, nil
}

// 检查积分是否超限
func (s *PointsProduceService) checkAstrict(c ctx.Context, params excitationReq.CheckAstrictReq, ruleInfo excitationModels.ExcitationIntegralRuleDetailModel) (bool, error) {
	if ruleInfo.IsAstrict == excitationModels.IsAstrict && ruleInfo.MaxIntegral > 0 {
		startTime, endTime, err := microWorkDao.NewFiscalMonth().GetTimeByAstrictType(c, params.EnterpriseId, ruleInfo.AstrictType, time.Unix(params.NowTime, 0))
		if err != nil {
			return false, err
		}
		// 获取当前用户该类型积分时间范围内的累计
		points, _ := s.PointsRecordDao.GetSumPoints(c, excitationReq.SumPointsReq{
			MemberId:     params.MemberId,
			EnterpriseId: params.EnterpriseId,
			Type:         params.Type,
			StartTime:    startTime,
			EndTime:      endTime,
		})
		if (points + params.AddPoints) > float64(ruleInfo.MaxIntegral) {
			return true, nil
		}
	}
	return false, nil
}

// TaskFinished 完成任务用方法
func (s *PointsProduceService) TaskFinished(c ctx.Context, req excitationReq.TaskFinishedReq) {
	one, _ := microScrmDao.NewPlanTaskAllotLogDao().GetOne(c, req.EnterpriseId, "id = ?", req.TaskAllotLogId)
	hashStr := fmt.Sprintf("-%s-%d", helper.RuleThree(req.MemberId > 0, strconv.FormatInt(req.MemberId, 10), req.Mobile), excitationModels.RuleCategoryId7)
	if one.IsFinished == microScrmModel.PLAN_FISISHING {
		s.ProduceQueue(ctx.Wrap(c), excitationReq.PointsProduceSave{
			MemberId:     req.MemberId,
			Mobile:       req.Mobile,
			Type:         excitationModels.RuleCategoryId7,
			Value:        1,
			UniqueHash:   strconv.FormatInt(req.TaskAllotLogId, 10) + hashStr,
			EnterpriseId: req.EnterpriseId,
		})
	}

	// 神策数据上报-完成任务
	memberInfo, _ := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel:    models.BaseModel{ID: req.MemberId},
			EnterpriseID: req.EnterpriseId,
		},
		Mobile: req.Mobile,
	})

	_, _ = microwork.NewShenCe().TaskComplete(c, req.EnterpriseId, memberInfo.Userid, one.TaskId)
}

// SaveCoins 加金币方法
func (s *PointsProduceService) SaveCoins(c ctx.Context, enterpriseID int64, sourceType int, dataList []excitationReq.ManualPointsSave) {
	if len(dataList) <= 0 {
		return
	}
	coinsSetting, err := excitationDao.NewCoinsSettings().GetOne(c, "enterprise_id=?", enterpriseID)
	if err != nil {
		log.Info(c).Msgf("未配置加金币规则，企业:%d", enterpriseID)
		return
	}
	if coinsSetting.Ratio <= 0 {
		log.Info(c).Msgf("加金币规则配置不正确，比率为：%f,企业:%d", coinsSetting.Ratio, enterpriseID)
		return
	}

	var remarkMap = map[int]string{
		1: "积分同步",
		2: "兑礼中心",
		3: "手动补发",
		4: "金币过期",
	}
	var coinsRecordList []microScrmModel.ExcitationCoinsRecordModel
	memberDataDao := excitationDao.NewMemberData()
	for _, record := range dataList {
		if !record.HasCoins {
			log.Info(c).Msgf("此条不加金币，企业:%d,数据%v", enterpriseID, record)
			continue
		}
		//计算金币，四舍五入保留一位小数
		coins := helper.MathRounded(float64(record.Points)*coinsSetting.Ratio, 1)
		if coins == 0 {
			log.Info(c).Msgf("四舍五入后金币为0，企业:%d,数据%v", enterpriseID, record)
			continue
		}
		if sourceType == microScrmModel.ExcitationCoinsRecordAutoType_Points {
			if coins < 0 {
				log.Info(c).Msgf("通过积分和奖章的方式均不扣减金币，企业:%d,数据%v", enterpriseID, record)
				continue
			}
		}
		// 处理导购金币更新
		coins, err = memberDataDao.SaveCoins(c, excitationReq.UpdateCoinsReq{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: enterpriseID,
			},
			MemberId: record.MemberId,
			Coins:    coins,
		})
		if err != nil {
			log.Info(c).Msgf("更新导购金币数据失败，企业:%d,错误:%v", enterpriseID, err)
			return
		}
		ratioSlice := strings.Split(coinsSetting.RatioRemark, ":")
		tmp := microScrmModel.ExcitationCoinsRecordModel{
			BaseModel:      models.CreateBase(app.UID()),
			MemberID:       record.MemberId,
			RatioRemark:    fmt.Sprintf("%s积分兑换%s金币", ratioSlice[0], ratioSlice[1]),
			Coins:          fmt.Sprintf("%0.1f", coins),
			SourceType:     sourceType,
			PointsRecordID: record.PointsRecordId,
			MedalRecordID:  record.MedalRecordId,
			CoinsRemark:    helper.RuleThree(record.PointsRemark != "", record.PointsRemark, remarkMap[sourceType]),
			OrderNum:       record.OrderNum,
		}
		coinsRecordList = append(coinsRecordList, tmp)
	}
	err = excitationDao.NewCoinsRecord().CreateInBatchesData(c, enterpriseID, coinsRecordList)
	if err != nil {
		log.Info(c).Msgf("批量写入金币记录失败，企业:%d,错误:%v", enterpriseID, err)
		return
	}

}
