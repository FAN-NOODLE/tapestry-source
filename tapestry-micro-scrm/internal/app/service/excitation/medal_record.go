package excitationSvc

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/http/request"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
)

type MedalRecordService struct {
	service.Base
	medalRecordDao *excitationDao.MedalRecord
}

func NewMedalRecordService() *MedalRecordService {
	return &MedalRecordService{
		medalRecordDao: excitationDao.NewMedalRecord(),
	}
}

func (s *MedalRecordService) GetList(c ctx.Context, params excitationReq.MedalRecordList) ([]excitationResp.MedalRecordList, int64, error) {
	var resp []excitationResp.MedalRecordList
	// 判断是否查部门
	if params.DepartmentId > 0 {
		department, err := departmentsDao.NewDepartments().GetOne(c, map[string]any{
			"enterprise_id": params.EnterpriseId,
			"third_id":      params.DepartmentId,
		})
		if err != nil {
			return nil, 0, err
		}
		if department.ID == 0 {
			return nil, 0, errors.New("部门未找到")
		}

		memberIds, err := membersDao.NewMemberDepartments().GetMemberIds(c, params.EnterpriseId, []int64{department.ID})
		if err != nil {
			return nil, 0, err
		}
		if len(memberIds) == 0 {
			return resp, 0, nil
		}
		params.MemberIds = append(params.MemberIds, memberIds...)
	}

	// 员工筛选
	if params.MemberUserId != "" {
		members, err := membersDao.NewMembers().GetUserId(c, params.EnterpriseId, params.MemberUserId)
		if err != nil {
			return nil, 0, err
		}
		if len(members) > 0 {
			for _, member := range members {
				params.MemberIds = append(params.MemberIds, member.ID)
			}
		} else {
			return resp, 0, nil
		}
	}

	// 员工姓名筛选
	if params.MemberName != "" {
		memberIds, _ := membersDao.NewMembers().GetMemberIdsByName(c, params.EnterpriseId, params.MemberName)
		if len(memberIds) > 0 {
			params.MemberIds = append(params.MemberIds, memberIds...)
		} else {
			return resp, 0, nil
		}
	}

	result, total, err := s.medalRecordDao.GetPagination(c, params)
	if err != nil {
		return resp, total, err
	}

	var medalRuleIds []int
	var memberIds []int64
	var ids []int64 //主键id
	for _, model := range result {
		medalRuleIds = append(medalRuleIds, model.RuleDimensionID, model.RuleCategoryID)
		memberIds = append(memberIds, model.MemberId)
		ids = append(ids, model.Id)
	}
	//查询奖章相关金币记录
	var medalCoinsMap = make(map[int64]string)
	if len(ids) > 0 {
		coinsList, err := excitationDao.NewCoinsRecord().GetMany(c, params.EnterpriseId, "medal_record_id in ?", ids)
		if err != nil {
			return nil, 0, err
		}
		for _, val := range coinsList {
			medalCoinsMap[val.MedalRecordID] = val.Coins
		}
	}

	// 查询导购信息
	memberDepartmentList, err := membersDao.NewMembers().GetMemberInfoAndDepartmentByMemberIds(c, params.EnterpriseId, memberIds)
	if err != nil {
		return nil, 0, err
	}

	// 整理导购id和导购信息的key=>value
	var memberMap = make(map[int64]excitationResp.MemberInfo)
	for _, info := range memberDepartmentList {
		dpName := ""
		if len(info.MemberDepartments) > 0 && info.MemberDepartments[0].Department.ID > 0 {
			dpName = info.MemberDepartments[0].Department.Name
		}
		memberMap[info.ID] = excitationResp.MemberInfo{
			MemberUserId:   info.Userid,
			MemberName:     info.Name,
			DepartmentName: dpName,
		}
	}

	medalRuleNameMap := make(map[int64]string)
	medalRule, err := excitationDao.NewMedalRule().GetMany(c, "id in (?) and status = ?", medalRuleIds, excitationModels.MedalRuleStatusEnable)
	if err != nil {
		return resp, total, err
	}
	for _, rule := range medalRule {
		medalRuleNameMap[rule.ID] = rule.RuleName
	}

	for _, model := range result {
		model.RuleDimensionName = medalRuleNameMap[int64(model.RuleDimensionID)]
		model.RuleCategoryName = medalRuleNameMap[int64(model.RuleCategoryID)]
		model.MemberUserId = memberMap[model.MemberId].MemberUserId
		model.MemberName = memberMap[model.MemberId].MemberName
		model.DepartmentName = memberMap[model.MemberId].DepartmentName
		if _, ok := medalCoinsMap[model.Id]; ok {
			model.Coins = medalCoinsMap[model.Id]
		}
		resp = append(resp, model)
	}

	return resp, total, err
}

func (s *MedalRecordService) Reissue(c ctx.Context, params excitationReq.MedalRecordReissue) (any, error) {
	var memberUserIds []string
	for _, data := range params.ReissueData {
		memberUserIds = append(memberUserIds, data.MemberUserId)
	}

	memberMap := make(map[string]int64)
	members, err := membersDao.NewMembers().GetManyByUserIds(c, params.EnterpriseId, memberUserIds)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		memberMap[member.Userid] = member.ID
	}

	// 领取记录
	var saveData []excitationModels.ExcitationMedalRecordModel
	// 发送积分
	var addMedalPoints []excitationReq.ManualPointsSave

	var exitsKey []string

	for _, data := range params.ReissueData {
		tmpKey := fmt.Sprintf("%s-%d", data.MemberUserId, data.MedalRuleId)
		// 效验奖章规则和奖章是否关联
		medalStyle, err := excitationDao.NewMedalStyle().GetOne(c, map[string]any{
			"enterprise_id": params.EnterpriseId,
			"medal_id":      data.MedalRuleId,
			"id":            data.MedalId,
		})
		if err != nil || medalStyle.ID == 0 {
			return nil, errors.New(fmt.Sprintf("奖章规则Id[ %d ],奖章Id[ %d ] 不匹配", data.MedalRuleId, data.MedalId))
		}

		// 获取规则id
		medal, err := excitationDao.NewMedal().GetOne(c, map[string]any{
			"enterprise_id": params.EnterpriseId,
			"id":            data.MedalRuleId,
		})
		if err != nil || medalStyle.ID == 0 {
			return nil, errors.New(fmt.Sprintf("奖章规则Id[ %d ] 不存在", data.MedalRuleId))
		}

		// 判断是否领取过奖章
		medalRecord, err := excitationDao.NewMedalRecord().GetOne(c, map[string]any{
			"enterprise_id":  params.EnterpriseId,
			"member_id":      memberMap[data.MemberUserId],
			"medal_style_id": data.MedalId,
		})
		if err != nil || medalRecord.ID > 0 {
			continue
		}

		// 获取领取过该规则奖章
		medalRuleRecord, err := excitationDao.NewMedalRecord().GetOne(c, map[string]any{
			"enterprise_id": params.EnterpriseId,
			"member_id":     memberMap[data.MemberUserId],
			"medal_id":      data.MedalRuleId,
		})

		// 针对某些规则只能领取一个处理
		if helper.InArray(medal.MedalRule.ShortName, []string{excitationModels.MedalRuleFriendNumberList}) &&
			(helper.InArray(tmpKey, exitsKey) || medalRuleRecord.ID > 0) {
			continue
		}

		medalRecordId := s.GetChannelId()
		saveData = append(saveData, excitationModels.ExcitationMedalRecordModel{
			MicroBaseModel: models.CreateMicroBase(params.EnterpriseId, medalRecordId),
			MedalID:        data.MedalRuleId,
			MedalStyleID:   data.MedalId,
			MemberID:       memberMap[data.MemberUserId],
		})
		// 规则存在积分
		if medal.Integral > 0 {
			medalPoints := excitationReq.ManualPointsSave{
				MemberId:      memberMap[data.MemberUserId],
				RuleName:      medal.Name,
				RuleDetailId:  0,
				PointsRemark:  fmt.Sprintf("发送奖章-%s", medalStyle.MedalName),
				Points:        int64(medal.Integral),
				MedalRecordId: medalRecordId, //奖章记录组装
			}
			addMedalPoints = append(addMedalPoints, medalPoints)
		}

		exitsKey = append(exitsKey, tmpKey)
	}

	err = s.DB().Transaction(func(tx *gorm.DB) error {
		if len(addMedalPoints) > 0 {
			err = NewPointsProduceService().MedalPoints(c, params.EnterpriseId, addMedalPoints)
			if err != nil {
				return err
			}
		}
		err = s.medalRecordDao.CreateInBatchesData(c, saveData)
		return err
	})

	return nil, err
}

func (s *MedalRecordService) Remove(c ctx.Context, enterpriseId int64, id int64) error {
	result, err := s.medalRecordDao.GetOne(c, "id", id)
	if err != nil {
		return err
	}

	// 获取规则id
	medal, err := excitationDao.NewMedal().GetUnscopedOne(c, map[string]any{
		"enterprise_id": enterpriseId,
		"id":            result.MedalID,
	})
	if err != nil {
		return err
	}

	if medal.ID == 0 {
		return errors.New("记录数据未找到")
	}

	// 获取奖章信息
	medalStyleRes, _ := excitationDao.NewMedalStyle().GetOne(c, map[string]any{
		"enterprise_id": enterpriseId,
		"id":            result.MedalStyleID,
	})

	err = s.DB().Transaction(func(tx *gorm.DB) error {

		// 清除积分
		if helper.InArray(medal.OtherRewards, []int32{2, 4}) && medal.Integral > 0 {
			var delMedalPoints []excitationReq.ManualPointsSave
			medalPoints := excitationReq.ManualPointsSave{
				MemberId:     result.MemberID,
				RuleName:     medal.Name,
				RuleDetailId: 0,
				PointsRemark: fmt.Sprintf("收回奖章-%s", medalStyleRes.MedalName),
				Points:       int64(0 - medal.Integral),
			}
			delMedalPoints = append(delMedalPoints, medalPoints)
			err = NewPointsProduceService().MedalPoints(c, enterpriseId, delMedalPoints)
			if err != nil {
				return err
			}
		}

		// 删除领取记录
		err = excitationDao.NewMedalRecord().DeleteData(c, map[string]any{"id": id})
		if err != nil {
			return err
		}
		// 记录导购等级数据
		err = NewMemberDataService().SaveMedalNum(c, excitationReq.UpdateMedalReq{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: enterpriseId,
			},
			MemberId: result.MemberID,
			MedalNum: -1,
		})
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *MedalRecordService) AlreadyGot(c ctx.Context, enterpriseId int64, mobile string) (excitationResp.MedalRecordAlreadyGot, error) {
	var resp excitationResp.MedalRecordAlreadyGot
	membersUser, err := s.GetUserIdByMobile(c, enterpriseId, mobile)
	if err != nil {
		return resp, err
	}

	// 获取已获得的奖章数量
	records, err := s.medalRecordDao.GetGained(c, enterpriseId, membersUser.ID)

	recordsMap := make(map[int64]excitationModels.ExcitationMedalRecordModel)
	var medalIds []int64
	var styleIds []int64
	for _, record := range records {
		recordsMap[record.ID] = record
		medalIds = append(medalIds, record.MedalID)
		styleIds = append(styleIds, record.MedalStyleID)
	}

	// 获取规则信息
	medalResult, err := excitationDao.NewMedal().GetManyUnscoped(c, "enterprise_id = ? and id in ?", enterpriseId, medalIds)
	if err != nil {
		return resp, err
	}
	medalResultMap := make(map[int64]excitationModels.ExcitationMedalModel)
	for _, res := range medalResult {
		medalResultMap[res.ID] = res
	}

	// 获取奖章信息
	medalStyleResult, err := excitationDao.NewMedalStyle().GetManyUnscoped(c, "enterprise_id = ? and id in ?", enterpriseId, styleIds)
	if err != nil {
		return resp, err
	}
	medalStyleResultMap := make(map[int64]excitationModels.ExcitationMedalStyleModel)
	for _, res := range medalStyleResult {
		medalStyleResultMap[res.ID] = res
	}
	var popRecordIds []int64
	for _, record := range records {
		resp.Number++
		if record.IsReceive == excitationModels.MedalRecordIsReceiveFalse {
			resp.NewMedalSign++
			if record.IsPop == excitationModels.MedalRecordIsPopFalse {
				resp.MedalList = append(resp.MedalList, excitationResp.MedalRecordAlreadyGotMedal{
					RecordId:     record.ID,
					RuleName:     medalResultMap[record.MedalID].Name,
					ReceivedTime: record.CreatedAt,
					LightIcon:    medalStyleResultMap[record.MedalStyleID].LightIcon,
					MedalName:    medalStyleResultMap[record.MedalStyleID].MedalName,
				})
				popRecordIds = append(popRecordIds, record.ID)
			}
		}
	}
	if len(popRecordIds) > 0 {
		// 弹出一次就不再弹出
		err = s.medalRecordDao.UpdateData(c,
			map[string]any{
				"enterprise_id": enterpriseId,
				"member_id":     membersUser.ID,
				"id":            popRecordIds,
			},
			map[string]any{
				"is_pop": excitationModels.MedalRecordIsPopTrue,
			})
	}
	return resp, err
}

func (s *MedalRecordService) MemberReceived(c ctx.Context, enterpriseId int64, ids []int64, mobile string) error {
	membersUser, err := s.GetUserIdByMobile(c, enterpriseId, mobile)
	if err != nil {
		return err
	}

	err = s.medalRecordDao.UpdateData(c,
		map[string]any{
			"enterprise_id": enterpriseId,
			"member_id":     membersUser.ID,
			"id":            ids,
		},
		map[string]any{
			"is_pop":     excitationModels.MedalRecordIsPopTrue,
			"is_receive": excitationModels.MedalRecordIsReceiveTrue,
		})
	if err != nil {
		return err
	}

	// 记录更新导购等级奖章数据
	err = NewMemberDataService().SaveMedalNum(c, excitationReq.UpdateMedalReq{
		EnterpriseBaseReq: request.EnterpriseBaseReq{
			EnterpriseId: enterpriseId,
		},
		MemberId: membersUser.ID,
		MedalNum: len(ids),
	})
	if err != nil {
		return err
	}
	return nil
}
