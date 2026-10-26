package excitationSvc

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"math"
	"micro-scrm/internal/app"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	membersDao "micro-scrm/internal/app/dao/members"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	openReq "micro-scrm/internal/app/http/request/open"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	openResp "micro-scrm/internal/app/http/response/open"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/app/service/department"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"time"
)

type MemberDataService struct {
	service.Base
	memberDataDao       *excitationDao.MemberData
	memberDataRecordDao *excitationDao.MemberDataRecord
}

func NewMemberDataService() *MemberDataService {
	return &MemberDataService{
		memberDataDao:       excitationDao.NewMemberData(),
		memberDataRecordDao: excitationDao.NewMemberDataRecord(),
	}
}

// DataIndex 小程序导购积分&等级
func (s *MemberDataService) DataIndex(c ctx.Context, params excitationReq.MemberDataIndex) (excitationResp.MemberDataIndexResp, error) {

	var (
		result      excitationResp.MemberDataIndexResp
		isLastLevel int
	)
	resultRule := make([]excitationResp.RelationRule, 0)
	// 用户信息
	member, err := membersDao.NewMembers().GetMemberByMobile(c, params.EnterpriseId, params.Mobile)
	if err != nil {
		return result, err
	}

	// 统计数据
	memberData, _ := s.memberDataDao.DataIndex(c, params.EnterpriseId, member.Id)
	// 下个等级别数据
	maxLevel, _ := excitationDao.NewLevelRule().GetMaxDetail(c, params.EnterpriseId, memberData.LevelType)
	// 查询等级图标
	levelIcon := excitationDao.NewLevelRule().GetLevelIcon(c, params.EnterpriseId, memberData.LevelID)

	// 能力图谱规则查询开启的
	incentiveRelation, _ := excitationDao.NewIncentiveRelation().GetMany(c, "enterprise_id = ? and status = ?", params.EnterpriseId, excitationModels.IncentiveConfigStatusOn)
	detailMap := make(map[int64]excitationModels.ExcitationIncentiveRelationModel)
	var relationIds []int64
	for _, detail := range incentiveRelation {
		detailMap[detail.ID] = detail
		relationIds = append(relationIds, detail.ID)
	}
	relationRuleResult, _ := excitationDao.NewIncentiveRelationRule().GetMany(c, "enterprise_id = ? and incentive_relation_id in(?)", params.EnterpriseId, relationIds)

	if err != nil {
		return excitationResp.MemberDataIndexResp{}, err
	}

	//获取当前导购所属奥莱或者零售
	storeType, _ := department.NewMemberDepartmentsService().GetMemberStoreType(c, params.EnterpriseId, member.Id)
	for _, value := range relationRuleResult {
		if detail, ok := detailMap[value.IncentiveRelationID]; ok {
			// 根据规则是否区分奥莱零售以及当前用户的奥莱零售类型过滤获取到的规则
			if detail.IsDiffStore == excitationModels.IsDiffStore_Yes {
				if value.StoreType == 0 || int(storeType) != value.StoreType {
					//规则要求区分奥莱零售,过滤与用户身份不符的规则
					continue
				}
			} else {
				//规则要求不区分奥莱零售，过滤不符合的规则
				if value.StoreType != 0 {
					continue
				}
			}
		}

		if value.IncentiveConfigType == excitationModels.IncentiveConfigTypeInvite && memberData.InviteValue == "" {
			continue
		}
		if value.IncentiveConfigType == excitationModels.IncentiveConfigTypeTask && memberData.TaskValue == "" {
			continue
		}
		if value.IncentiveConfigType == excitationModels.IncentiveConfigTypeInteract && memberData.InteractValue == "" {
			continue
		}
		if value.IncentiveConfigType == excitationModels.IncentiveConfigTypeTrain && memberData.TrainValue == "" {
			continue
		}
		if value.IncentiveConfigType == excitationModels.IncentiveConfigTypeOnline && memberData.OnlineValue == "" {
			continue
		}

		resultRule = append(resultRule, excitationResp.RelationRule{
			IncentiveConfigType: value.IncentiveConfigType,
			RuleCondition:       value.RuleCondition,
			RuleValue:           value.RuleValue,
			RuleGive:            value.RuleGive,
		})
	}

	//查看下一个等级是不是当前等级如果是那就是最后一个等级
	if maxLevel.ID == memberData.LevelID {
		isLastLevel = 1
	}
	// 信息数据
	result = excitationResp.MemberDataIndexResp{
		MemberId:      member.Id,
		Avatar:        member.Avatar,
		Name:          member.Name,
		Level:         memberData.Level,
		LevelIcon:     levelIcon,
		PointValue:    memberData.PointValue,
		MedalNum:      memberData.MedalNum,
		MaxLevel:      maxLevel.LevelName,
		MaxLevelValue: maxLevel.MaxValue,
		RelationRule:  resultRule,
		IsLastLevel:   isLastLevel,
		Coins:         memberData.Coins,
	}

	// 查询数据图谱开启状态如果没开启统计能力图谱数据就设置成空
	for _, v := range incentiveRelation {
		switch v.IncentiveConfigType {
		case excitationModels.IncentiveConfigTypeInvite: // 邀约维度
			if v.Status == excitationModels.IncentiveConfigStatusOn {
				result.InviteValue = memberData.InviteValue
				result.InviteRuleValue = memberData.InviteRuleValue
				result.InviteUpdatedAt = memberData.InviteUpdatedAt
				result.InviteRuleStr = memberData.InviteRuleStr

			}
		case excitationModels.IncentiveConfigTypeTask: // 任务维度
			if v.Status == excitationModels.IncentiveConfigStatusOn {
				result.TaskValue = memberData.TaskValue
				result.TaskRuleValue = memberData.TaskRuleValue
				result.TaskUpdatedAt = memberData.TaskUpdatedAt
				result.TaskRuleStr = memberData.TaskRuleStr
			}
		case excitationModels.IncentiveConfigTypeInteract: // 客户互动维度
			if v.Status == excitationModels.IncentiveConfigStatusOn {
				result.InteractValue = memberData.InteractValue
				result.InteractRuleValue = memberData.InteractRuleValue
				result.InteractUpdatedAt = memberData.InteractUpdatedAt
				result.InteractRuleStr = memberData.InteractRuleStr
			}
		case excitationModels.IncentiveConfigTypeTrain: // 培训维度
			if v.Status == excitationModels.IncentiveConfigStatusOn {
				result.TrainValue = memberData.TrainValue
				result.TrainRuleValue = memberData.TrainRuleValue
				result.TrainUpdatedAt = memberData.TrainUpdatedAt
				result.TrainRuleStr = memberData.TrainRuleStr
			}
		case excitationModels.IncentiveConfigTypeOnline: // 线上能力
			if v.Status == excitationModels.IncentiveConfigStatusOn {
				result.OnlineValue = memberData.OnlineValue
				result.OnlineRuleValue = memberData.OnlineRuleValue
				result.OnlineUpdatedAt = memberData.OnlineUpdatedAt
				result.OnlineRuleStr = memberData.OnlineRuleStr
			}
		}
	}

	return result, nil
}

func (s *MemberDataService) GradePoints(c ctx.Context, params excitationReq.GradePoints) (excitationResp.GradePointsResp, error) {
	var (
		result     excitationResp.GradePointsResp
		percentage int64
	)
	// 用户信息
	member, err := membersDao.NewMembers().GetMemberByMobile(c, params.EnterpriseId, params.Mobile)
	if err != nil {
		return result, err
	}
	// 统计数据
	memberData, _ := s.memberDataDao.DataIndex(c, params.EnterpriseId, member.Id)
	// 等级列表
	LevelRuleDetailReq := excitationReq.LevelRuleDetailReq{
		EnterpriseId: params.EnterpriseId,
	}
	// 等级数据
	LevelRule := result.LevelRule
	levelRuleRes, err := excitationDao.NewLevelRule().GetDetail(c, LevelRuleDetailReq)
	if err != nil {
		return result, err
	}
	// 能力图谱规则
	relationRuleResult, _ := excitationDao.NewIncentiveRelationRule().GetMany(c, "enterprise_id = ?", params.EnterpriseId)
	for _, value := range levelRuleRes {
		// 等级规则附加条件
		var (
			rule     []excitationResp.RuleResp
			selected bool
		)
		for _, relatValue := range relationRuleResult {
			if value.RuleInviteID == relatValue.ID {
				rule = append(rule, excitationResp.RuleResp{
					IncentiveConfigType: relatValue.IncentiveConfigType,
					RuleCondition:       value.RuleConditionInvite,
					RuleValue:           relatValue.RuleValue,
					RuleGive:            relatValue.RuleGive,
				})
			}
			if value.RuleTaskID == relatValue.ID {
				rule = append(rule, excitationResp.RuleResp{
					IncentiveConfigType: relatValue.IncentiveConfigType,
					RuleCondition:       value.RuleConditionTask,
					RuleValue:           relatValue.RuleValue,
					RuleGive:            relatValue.RuleGive,
				})
			}
			if value.RuleInteractID == relatValue.ID {
				rule = append(rule, excitationResp.RuleResp{
					IncentiveConfigType: relatValue.IncentiveConfigType,
					RuleCondition:       value.RuleConditionInteract,
					RuleValue:           relatValue.RuleValue,
					RuleGive:            relatValue.RuleGive,
				})
			}
			if value.RuleTrainID == relatValue.ID {
				rule = append(rule, excitationResp.RuleResp{
					IncentiveConfigType: relatValue.IncentiveConfigType,
					RuleCondition:       value.RuleConditionTrain,
					RuleValue:           relatValue.RuleValue,
					RuleGive:            relatValue.RuleGive,
				})
			}
			if value.RuleOnlineID == relatValue.ID {
				rule = append(rule, excitationResp.RuleResp{
					IncentiveConfigType: relatValue.IncentiveConfigType,
					RuleCondition:       value.RuleConditionOnline,
					RuleValue:           relatValue.RuleValue,
					RuleGive:            relatValue.RuleGive,
				})
			}
		}

		// 当前所在的等级
		if value.ID == memberData.LevelID {
			selected = true
		} else {
			selected = false
		}

		// 等级
		LevelRule = append(LevelRule, excitationResp.LevelRuleResp{
			Id:        value.ID,
			LevelName: value.LevelName,
			LevelIcon: value.LevelIcon,
			MinValue:  value.MinValue,
			MaxValue:  value.MaxValue,
			Selected:  selected,
			Rule:      rule,
		})
	}

	//拿到当前等级在循环等级取出
	// 百分比
	memberCount := s.memberDataDao.GetMemberCount(c, params.EnterpriseId)
	levelCount := s.memberDataDao.GetLevelCount(c, params.EnterpriseId, memberData.LevelType)
	if levelCount <= 0 {
		percentage = 0
	} else {
		percentage = int64(float64(levelCount) / float64(memberCount) * 100)
	}

	// 信息数据
	result = excitationResp.GradePointsResp{
		LevelName:      memberData.Level,
		LevelUpdatedAt: memberData.LevelUpdatedAt,
		Percentage:     percentage,
		LevelRule:      LevelRule,
	}

	return result, nil

}

// GetPointsRules 小程序端获取积分规则信息接
func (s *MemberDataService) GetPointsRules(c ctx.Context, params excitationReq.GetPointsRulesReq) ([]excitationResp.GetPointsRulesResp, error) {
	var resp []excitationResp.GetPointsRulesResp
	result, err := excitationDao.NewIntegralRuleDetail().GetMany(c, "enterprise_id = ? and rule_category_id in ?", params.EnterpriseId, params.RuleCategoryId)
	if err != nil {
		return resp, err
	}

	//获取当前导购所属奥莱或者零售
	storeType, _ := department.NewMemberDepartmentsService().GetMemberStoreType(c, params.EnterpriseId, params.MemberId)
	for _, value := range result {
		// 根据规则是否区分奥莱零售以及当前用户的奥莱零售类型过滤获取到的规则
		if value.IsDiffStore == excitationModels.IsDiffStore_Yes {
			if value.StoreType == 0 || storeType != value.StoreType {
				//规则要求区分奥莱零售,过滤与用户身份不符的规则
				continue
			}
		} else {
			//规则要求不区分奥莱零售，过滤不符合的规则
			if value.StoreType != 0 {
				continue
			}
		}
		var points float64
		if value.IsAstrict == excitationModels.IsAstrict {
			startTime, endTime, err := microWorkDao.NewFiscalMonth().GetTimeByAstrictType(c, params.EnterpriseId, value.AstrictType, time.Now())
			if err != nil {
				return nil, err
			}
			// 获取当前用户该类型当月积分累计
			points, _ = excitationDao.NewPointsRecord().GetSumPoints(c, excitationReq.SumPointsReq{
				MemberId:     params.MemberId,
				EnterpriseId: params.EnterpriseId,
				Type:         value.RuleCategoryID,
				StartTime:    startTime,
				EndTime:      endTime,
			})
		}

		resp = append(resp, excitationResp.GetPointsRulesResp{
			EnterpriseId:   value.EnterpriseID,
			RuleCategoryId: value.RuleCategoryID,
			Integral:       value.Integral,
			IntegralType:   value.IntegralType,
			RuleValue:      value.RuleValue,
			AstrictType:    value.AstrictType,
			MaxIntegral:    value.MaxIntegral,
			NowIntegral:    int64(points),
		})
	}
	return resp, err
}

// GetList 列表
func (s *MemberDataService) GetList(c ctx.Context, params excitationReq.MemberDataListReq) ([]excitationResp.MemberDataListResp, int64, error) {
	result, total, err := s.memberDataDao.GetList(c, params)
	if err != nil {
		return result, total, err
	}

	var memberIDs []int64
	for _, v := range result {
		memberIDs = append(memberIDs, v.MemberID)
	}

	// 查询导购信息
	memberDepartmentList, err := membersDao.NewMembers().GetMemberInfoAndDepartmentByMemberIds(c, params.EnterpriseId, memberIDs)
	if err != nil {
		return nil, 0, err
	}

	// 整理导购id和导购信息的key=>value
	var memberMap = make(map[int64]excitationResp.MemberInfo)
	for _, info := range memberDepartmentList {
		dpName := ""
		thirdId := 0
		if len(info.MemberDepartments) > 0 && info.MemberDepartments[0].Department.ID > 0 {
			dpName = info.MemberDepartments[0].Department.Name
			thirdId = info.MemberDepartments[0].Department.ThirdID
		}
		memberMap[info.ID] = excitationResp.MemberInfo{
			MemberUserId:   info.Userid,
			MemberName:     info.Name,
			DepartmentName: dpName,
			ThirdId:        thirdId,
		}
	}

	//获取导购等级历史进程
	levelRecordList, err := s.memberDataRecordDao.GetMemberMaxLevelRecord(c, params.EnterpriseId, memberIDs)

	if err != nil {
		return nil, 0, nil
	}
	var levelRecordMap = make(map[int64]excitationResp.LevelRecordResp)
	for _, record := range levelRecordList {
		if _, ok := levelRecordMap[record.MemberId]; !ok {
			levelRecordMap[record.MemberId] = record
		}
	}

	for key, item := range result {
		if _, ok := memberMap[item.MemberID]; ok {
			result[key].MemberName = memberMap[item.MemberID].MemberName
			result[key].MemberUserid = memberMap[item.MemberID].MemberUserId
			result[key].ThirdId = memberMap[item.MemberID].ThirdId
			result[key].ThirdName = memberMap[item.MemberID].DepartmentName
		}
		if _, ok := levelRecordMap[item.MemberID]; ok {
			result[key].MaxLevel = levelRecordMap[item.MemberID].Level
			result[key].MaxLevelTime = levelRecordMap[item.MemberID].CreatedAt
		}
	}

	return result, total, err
}

// 积分更新
func (s *MemberDataService) savePoints(c ctx.Context, params excitationReq.UpdatePointsReq) error {
	info, err := s.memberDataDao.GetOne(c, "enterprise_id = ? and member_id = ?", params.EnterpriseId, params.MemberId)
	if err != nil {
		return err
	}
	if info.ID > 0 {
		oldPoints, _ := strconv.ParseFloat(info.PointValue, 64)
		newPoints := oldPoints + params.Points
		updateData := map[string]interface{}{
			"point_value":      fmt.Sprintf("%0.1f", helper.RuleThree(newPoints > 0, newPoints, 0)),
			"point_updated_at": time.Now().Unix(),
		}
		err = s.memberDataDao.UpdateData(c,
			map[string]interface{}{
				"enterprise_id": params.EnterpriseId,
				"id":            info.ID,
			}, updateData)
	} else {
		// 创建的时候默认给 最低等级 按照每天处理更新等级
		memberLevelRule, err1 := excitationDao.NewLevelRule().GetLevelInfoForPoint(c, params.EnterpriseId, "0")
		if err1 != nil {
			return err1
		}
		if memberLevelRule.ID == 0 {
			return nil
		}
		err := s.MemberLevelInit(c, excitationModels.ExcitationMemberDataModel{
			MicroBaseModel: models.CreateMicroBase(params.EnterpriseId, app.UID()),
			MemberID:       params.MemberId,
			PointValue:     fmt.Sprintf("%0.1f", params.Points),
			PointUpdatedAt: int(time.Now().Unix()),
			Level:          memberLevelRule.LevelName,
			LevelType:      memberLevelRule.LevelType,
			LevelID:        memberLevelRule.ID,
			LevelUpdatedAt: int(time.Now().Unix()),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveLevel 更新等级
func (s *MemberDataService) SaveLevel(c ctx.Context, params excitationReq.UpdateLevelReq) error {
	info, err := s.memberDataDao.GetOne(c, "enterprise_id = ? and member_id = ?", params.EnterpriseId, params.MemberId)
	if err != nil {
		return err
	}
	if info.ID > 0 {
		updateData := map[string]interface{}{
			"level":            params.Level,
			"level_id":         params.LevelId,
			"level_updated_at": time.Now().Unix(),
		}
		err = s.memberDataDao.UpdateData(c,
			map[string]interface{}{
				"enterprise_id": params.EnterpriseId,
				"id":            info.ID,
			}, updateData)
	} else {
		err := s.MemberLevelInit(c, excitationModels.ExcitationMemberDataModel{
			MicroBaseModel: models.CreateMicroBase(params.EnterpriseId, app.UID()),
			MemberID:       params.MemberId,
			Level:          params.Level,
			LevelID:        params.LevelId,
			LevelUpdatedAt: int(time.Now().Unix()),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveMedalNum 更新奖章
func (s *MemberDataService) SaveMedalNum(c ctx.Context, params excitationReq.UpdateMedalReq) error {
	info, err := s.memberDataDao.GetOne(c, "enterprise_id = ? and member_id = ?", params.EnterpriseId, params.MemberId)
	if err != nil {
		return err
	}
	if info.ID > 0 {
		sign := helper.RuleThree(params.MedalNum >= 0, "+", "-")
		// 回收奖章 更新奖章 时间为最近获取的奖章时间
		medalUpdatedAt := time.Now().Unix()
		if params.MedalNum < 0 {
			// 如果总数最后为0或负数
			if info.MedalNum-int(math.Abs(float64(params.MedalNum))) <= 0 {
				updateData := map[string]interface{}{
					"medal_num":        0,
					"medal_updated_at": 0,
				}
				err = s.memberDataDao.UpdateData(c,
					map[string]interface{}{
						"enterprise_id": params.EnterpriseId,
						"member_id":     params.MemberId,
					}, updateData)
			} else {
				// 最近获取奖章时间
				lastMemberReceived, errs := excitationDao.NewMedalRecord().GetLastMemberReceived(c, params.EnterpriseId, params.MemberId)
				if errs != nil {
					return errs
				}
				medalUpdatedAt = lastMemberReceived.CreatedAt
				updateData := map[string]interface{}{
					"medal_num":        gorm.Expr("medal_num "+sign+" ?", math.Abs(float64(params.MedalNum))),
					"medal_updated_at": medalUpdatedAt,
				}
				err = s.memberDataDao.UpdateData(c,
					map[string]interface{}{
						"enterprise_id": params.EnterpriseId,
						"member_id":     params.MemberId,
					}, updateData)
			}

		} else {
			updateData := map[string]interface{}{
				"medal_num":        gorm.Expr("medal_num "+sign+" ?", math.Abs(float64(params.MedalNum))),
				"medal_updated_at": medalUpdatedAt,
			}
			err = s.memberDataDao.UpdateData(c,
				map[string]interface{}{
					"enterprise_id": params.EnterpriseId,
					"member_id":     params.MemberId,
				}, updateData)
		}

	} else {
		// 创建的时候默认给 最低等级 按照每天处理更新等级
		memberLevelRule, err1 := excitationDao.NewLevelRule().GetLevelInfoForPoint(c, params.EnterpriseId, "0")
		if err1 != nil {
			return err1
		}
		if memberLevelRule.ID == 0 {
			return nil
		}
		err := s.MemberLevelInit(c, excitationModels.ExcitationMemberDataModel{
			MicroBaseModel: models.CreateMicroBase(params.EnterpriseId, app.UID()),
			MemberID:       params.MemberId,
			MedalNum:       params.MedalNum,
			MedalUpdatedAt: int(time.Now().Unix()),
			Level:          memberLevelRule.LevelName,
			LevelType:      memberLevelRule.LevelType,
			LevelID:        memberLevelRule.ID,
			LevelUpdatedAt: int(time.Now().Unix()),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveDimensionField 导购各维度字段更新
func (s *MemberDataService) SaveDimensionField(c ctx.Context, params excitationReq.UpdateDimensionFieldReq) error {

	// 各维度类型
	var fieldStr string
	switch params.IncentiveConfigType {
	case excitationModels.IncentiveConfigTypeInvite:
		fieldStr = "invite"
		break
	case excitationModels.IncentiveConfigTypeTask:
		fieldStr = "task"
		break
	case excitationModels.IncentiveConfigTypeInteract:
		fieldStr = "interact"
		break
	case excitationModels.IncentiveConfigTypeTrain:
		fieldStr = "train"
		break
	case excitationModels.IncentiveConfigTypeOnline:
		fieldStr = "online"
		break
	}
	if fieldStr == "" {
		return nil
	}

	timeNow := time.Now().Unix()
	if params.Id > 0 {
		updateData := map[string]interface{}{
			fieldStr + "_value":      params.DimensionValue,
			fieldStr + "_rule_value": params.DimensionRuleValue,
			fieldStr + "_rule_str":   params.DimensionRuleStr,
			fieldStr + "_updated_at": timeNow,
		}
		err := s.memberDataDao.UpdateData(c,
			map[string]interface{}{
				"enterprise_id": params.EnterpriseId,
				"id":            params.Id,
			}, updateData)
		if err != nil {
			return err
		}

	} else {
		memberLevelRule, err1 := excitationDao.NewLevelRule().GetLevelInfoForPoint(c, params.EnterpriseId, "0")

		if err1 != nil {
			return err1
		}
		if memberLevelRule.ID == 0 {
			return nil
		}

		// 新增默认数据
		createData := map[string]interface{}{
			"id":                     app.UID(),
			"enterprise_id":          params.EnterpriseId,
			"member_id":              params.MemberId,
			"point_value":            0,
			"point_updated_at":       timeNow,
			"level":                  memberLevelRule.LevelName,
			"level_type":             memberLevelRule.LevelType,
			"level_id":               memberLevelRule.ID,
			"level_updated_at":       timeNow,
			"created_at":             timeNow,
			"updated_at":             timeNow,
			fieldStr + "_value":      params.DimensionValue,
			fieldStr + "_rule_value": params.DimensionRuleValue,
			fieldStr + "_rule_str":   params.DimensionRuleStr,
			fieldStr + "_updated_at": timeNow,
		}
		err := s.memberDataDao.CreateOne(c, createData)

		if err != nil {
			return err
		}
	}

	// 添加等级或关联图谱信息操作记录
	recordData := map[string]interface{}{
		"id":                app.UID(),
		"enterprise_id":     params.EnterpriseId,
		"member_id":         params.MemberId,
		"created_at":        timeNow,
		"updated_at":        timeNow,
		fieldStr + "_value": params.DimensionValue,
	}
	err := s.memberDataRecordDao.CreateOne(c, recordData)
	if err != nil {
		return err
	}

	return nil
}

// CheckOrCreateMemberData 验证导购等级数据或创建
func (s *MemberDataService) CheckOrCreateMemberData(c ctx.Context, enterpriseId, memberId int64) (memberDataInfo excitationModels.ExcitationMemberDataModel, err error) {

	memberDataInfo, err = s.memberDataDao.GetOne(c, "enterprise_id = ? and member_id = ?", enterpriseId, memberId)
	if err != nil {
		return memberDataInfo, err
	}
	timeNow := time.Now().Unix()
	if memberDataInfo.ID == 0 {
		memberLevelRule, err1 := excitationDao.NewLevelRule().GetLevelInfoForPoint(c, memberDataInfo.EnterpriseID, "0")

		if err1 != nil {
			return memberDataInfo, err1
		}
		if memberLevelRule.ID == 0 {
			return memberDataInfo, nil
		}

		// 新增默认数据
		createData := map[string]interface{}{
			"id":               app.UID(),
			"enterprise_id":    enterpriseId,
			"member_id":        memberId,
			"point_value":      0,
			"point_updated_at": timeNow,
			"level":            memberLevelRule.LevelName,
			"level_type":       memberLevelRule.LevelType,
			"level_id":         memberLevelRule.ID,
			"level_updated_at": timeNow,
			"created_at":       timeNow,
			"updated_at":       timeNow,
		}
		err = s.memberDataDao.CreateOne(c, createData)

		if err != nil {
			return memberDataInfo, err
		}
		memberDataInfo, err = s.memberDataDao.GetOne(c, "enterprise_id = ? and member_id = ?", enterpriseId, memberId)
		if err != nil {
			return memberDataInfo, err
		}
	}

	return memberDataInfo, err
}

// UpdateLevel 定时任务处理导购等级变更
func (s *MemberDataService) UpdateLevel(c ctx.Context, params excitationReq.UpdateLevelReq) error {

	memberDataInfo, err := s.CheckOrCreateMemberData(c, params.EnterpriseId, params.MemberId)

	if err != nil {
		return err
	}
	timeNow := time.Now().Unix()

	// 查询导购积分所在等级规则
	memberLevelRule, err1 := excitationDao.NewLevelRule().GetLevelInfoForPoint(c, params.EnterpriseId, memberDataInfo.PointValue)

	if err1 != nil {
		return err1
	}
	if memberLevelRule.ID == 0 {
		return nil
	}

	// 处理导购等级数据
	isCheck, errs := s.handleLevelRule(c, params, memberDataInfo, memberLevelRule)
	if errs != nil {
		return errs
	}

	var levelName = memberLevelRule.LevelName
	var levelId = memberLevelRule.ID
	var levelType = memberLevelRule.LevelType

	// 未满足条件
	if isCheck == false {
		// 积分跨等级情况 积分满足但是附加条件不满足 降一个等级验证
		for i := 1; i < memberLevelRule.LevelType; i++ {

			checkMemberLevelRule, err2 := excitationDao.NewLevelRule().GetLevelInfoForLevel(c, memberDataInfo.EnterpriseID, memberLevelRule.LevelType-i)

			if err2 != nil {
				return err2
			}
			if memberLevelRule.ID == 0 || checkMemberLevelRule.ID == 0 {
				return nil
			}
			isCheck, err = s.handleLevelRule(c, params, memberDataInfo, checkMemberLevelRule)
			if err != nil {
				return err
			}

			if isCheck == true {
				levelName = checkMemberLevelRule.LevelName
				levelId = checkMemberLevelRule.ID
				levelType = checkMemberLevelRule.LevelType
				break
			}
		}

	}
	if isCheck == true {
		// 默认添加初始等级
		if memberDataInfo.ID == 0 {
			err := s.MemberLevelInit(c, excitationModels.ExcitationMemberDataModel{
				MicroBaseModel: models.CreateMicroBase(params.EnterpriseId, app.UID()),
				MemberID:       params.MemberId,
				PointValue:     "0",
				PointUpdatedAt: int(timeNow),
				Level:          levelName,
				LevelType:      levelType,
				LevelID:        levelId,
				LevelUpdatedAt: int(timeNow),
			})
			if err != nil {
				return err
			}
		} else {

			if levelId != memberDataInfo.LevelID || levelType != memberDataInfo.LevelType {

				// 满足条件则升级
				updateLevelData := map[string]interface{}{
					"level":            levelName,
					"level_type":       levelType,
					"level_id":         levelId,
					"level_updated_at": timeNow,
				}
				err = s.memberDataDao.UpdateData(c,
					map[string]interface{}{
						"enterprise_id": params.EnterpriseId,
						"member_id":     params.MemberId,
					}, updateLevelData)
				if err != nil {
					return err
				}

				// 添加等级或关联图谱信息操作记录
				recordData := map[string]interface{}{
					"id":            app.UID(),
					"enterprise_id": params.EnterpriseId,
					"member_id":     params.MemberId,
					"level":         levelName,
					"created_at":    timeNow,
					"updated_at":    timeNow,
				}
				err = s.memberDataRecordDao.CreateOne(c, recordData)
				if err != nil {
					return err
				}

			}

		}
	}
	return nil

}

// handleLevelRule 处理导购等级数据
func (s *MemberDataService) handleLevelRule(c ctx.Context, params excitationReq.UpdateLevelReq, memberDataInfo excitationModels.ExcitationMemberDataModel, memberLevelRule excitationResp.LevelRuleDetailInfoResp) (bool, error) {
	// 判断附加条件
	var levelRuleIDs []int64
	var tempDiff = []int64{0}
	levelRuleIDs = append(levelRuleIDs, memberLevelRule.RuleInviteID, memberLevelRule.RuleTaskID, memberLevelRule.RuleInteractID, memberLevelRule.RuleTrainID, memberLevelRule.RuleOnlineID)
	levelRuleIDs, _ = helper.ArrayDifference(levelRuleIDs, tempDiff)

	// 默认满足
	var isCheck = true
	if len(levelRuleIDs) > 0 {
		// 查询规则信息
		var relationRuleIdsWhere = map[string]interface{}{
			"enterprise_id": params.EnterpriseId,
			"id":            levelRuleIDs,
		}
		levelRules, err2 := excitationDao.NewIncentiveRelationRule().GetMany(c, relationRuleIdsWhere, "id > 0")
		if err2 != nil {
			return true, err2
		}

		if len(levelRules) > 0 { // 有附加条件
			// 能力指标
			relationRuleMap := s.RelationRuleMap()
			// 导购要满足的指标值类型
			var ruleConditionKey int
			// 导购要满足的指标值
			var memberRuleGiveValue string
			var memberRuleValue string
			// 判断指标能力条件 达成>= 所有附加条件都要满足
			for _, levelRule := range levelRules {

				memberRuleGiveValue = levelRule.RuleGive
				if levelRule.IncentiveConfigType == excitationModels.IncentiveConfigTypeInvite {
					// 导购所在等级指标值
					memberRuleValue = memberDataInfo.InviteValue
					// 导购积分达到的等级规则
					ruleConditionKey = memberLevelRule.RuleConditionInvite
					isCheck = s.checkLevelRule(isCheck, memberRuleValue, memberRuleGiveValue, ruleConditionKey, relationRuleMap, levelRule)
					if isCheck == false {
						break
					}
				}
				if levelRule.IncentiveConfigType == excitationModels.IncentiveConfigTypeTask {

					// 导购所在等级指标值
					memberRuleValue = memberDataInfo.TaskValue
					// 导购积分达到的等级规则
					ruleConditionKey = memberLevelRule.RuleConditionTask
					isCheck = s.checkLevelRule(isCheck, memberRuleValue, memberRuleGiveValue, ruleConditionKey, relationRuleMap, levelRule)
					if isCheck == false {
						break
					}
				}
				if levelRule.IncentiveConfigType == excitationModels.IncentiveConfigTypeInteract {
					// 导购所在等级指标值
					memberRuleValue = memberDataInfo.InteractValue
					// 导购积分达到的等级规则
					ruleConditionKey = memberLevelRule.RuleConditionInteract

					isCheck = s.checkLevelRule(isCheck, memberRuleValue, memberRuleGiveValue, ruleConditionKey, relationRuleMap, levelRule)
					if isCheck == false {
						break
					}
				}
				if levelRule.IncentiveConfigType == excitationModels.IncentiveConfigTypeTrain {

					// 导购所在等级指标值
					memberRuleValue = memberDataInfo.TrainValue
					// 导购积分达到的等级规则
					ruleConditionKey = memberLevelRule.RuleConditionTrain
					isCheck = s.checkLevelRule(isCheck, memberRuleValue, memberRuleGiveValue, ruleConditionKey, relationRuleMap, levelRule)

					if isCheck == false {
						break
					}
				}
				if levelRule.IncentiveConfigType == excitationModels.IncentiveConfigTypeOnline {
					//if _, ok := relationRuleMap[memberDataInfo.OnlineValue]; ok {
					//	memberInviteKey = relationRuleMap[memberDataInfo.OnlineValue].Key
					//}

					// 导购所在等级指标值
					memberRuleValue = memberDataInfo.OnlineValue
					// 导购积分达到的等级规则
					ruleConditionKey = memberLevelRule.RuleConditionOnline
					isCheck = s.checkLevelRule(isCheck, memberRuleValue, memberRuleGiveValue, ruleConditionKey, relationRuleMap, levelRule)

					if isCheck == false {
						break
					}
				}
			}
		}
	}

	return isCheck, nil

}

// CheckLevelRule 判断等级附加条件 达到
func (s *MemberDataService) checkLevelRule(isCheck bool, memberRuleValue, memberRuleGiveValue string, ruleConditionKey int, relationRuleMap map[string]excitationReq.RelationRuleArray, levelRule excitationModels.ExcitationIncentiveRelationRuleModel) bool {
	// 根据指标值
	var memberRuleKey, memberRuleGiveKey int
	if _, ok := relationRuleMap[memberRuleValue]; ok {

		// 导购所获得的指标
		memberRuleKey = relationRuleMap[memberRuleValue].Key
		if _, ok1 := relationRuleMap[levelRule.RuleGive]; ok1 {
			// 等级规则所要求的的指标
			memberRuleGiveKey = relationRuleMap[memberRuleGiveValue].Key

			// 达到
			if ruleConditionKey == excitationModels.RuleCondition3 {
				if memberRuleGiveKey > memberRuleKey {
					// 不满足
					isCheck = false
				}
			}
		}
	} else {
		// 有附加条件但是导购获得值没有
		isCheck = false
	}

	return isCheck

}

// RelationRuleMap 能力指标
func (s *MemberDataService) RelationRuleMap() map[string]excitationReq.RelationRuleArray {
	var relationRuleArray []excitationReq.RelationRuleArray
	relationRuleArray = append(relationRuleArray,
		excitationReq.RelationRuleArray{
			Key:   1,
			Value: "A-",
		},
		excitationReq.RelationRuleArray{
			Key:   2,
			Value: "A",
		},
		excitationReq.RelationRuleArray{
			Key:   3,
			Value: "A+",
		},
		excitationReq.RelationRuleArray{
			Key:   4,
			Value: "S",
		},
		excitationReq.RelationRuleArray{
			Key:   5,
			Value: "SS",
		},
	)

	var relationRuleMap = make(map[string]excitationReq.RelationRuleArray)
	for _, v := range relationRuleArray {
		relationRuleMap[v.Value] = v
	}

	return relationRuleMap
}

// MemberLevelRecord  导购等级历史进程
func (s *MemberDataService) MemberLevelRecord(c ctx.Context, enterpriseId int64, memberId int64) (res []excitationResp.LevelRecordResp, err error) {

	allRecord, err := s.memberDataRecordDao.GetMemberLevelRecord(c, enterpriseId, []int64{memberId}, "created_at asc")
	var signLevel string //标记等级变更
	for _, record := range allRecord {
		if signLevel != record.Level {
			res = append(res, record)
		}
		signLevel = record.Level
	}

	if err != nil {
		return res, err
	}

	return res, nil
}

// MemberLevelInit 初始化导购等级
func (s *MemberDataService) MemberLevelInit(c ctx.Context, data excitationModels.ExcitationMemberDataModel) error {
	err := s.memberDataDao.CreateData(c, data)
	if err != nil {
		return err
	}
	// 添加等级或关联图谱信息操作记录
	recordData := map[string]interface{}{
		"id":            app.UID(),
		"enterprise_id": data.EnterpriseID,
		"member_id":     data.MemberID,
		"level":         data.Level,
		"created_at":    data.CreatedAt,
		"updated_at":    data.UpdatedAt,
	}
	return s.memberDataRecordDao.CreateOne(c, recordData)
}

// MemberDataInfo  获取导购金币详情
func (s *MemberDataService) MemberDataInfo(c ctx.Context, params openReq.MemberInfoReq) (res openResp.MemberInfoResp, err error) {
	memberUserInfo, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{"enterprise_id": params.EnterpriseId, "userid": params.UserId})
	if err != nil {
		return res, err
	}
	if memberUserInfo.ID == 0 {
		return res, errors.New("导购不存在")
	}
	info, err := s.memberDataDao.GetOne(c, "enterprise_id = ? and member_id = ?", params.EnterpriseId, memberUserInfo.ID)
	if err != nil {
		return res, err
	}
	return openResp.MemberInfoResp{
		Coins:  info.Coins,
		Avatar: memberUserInfo.Avatar,
		Name:   memberUserInfo.Name,
	}, nil
}
