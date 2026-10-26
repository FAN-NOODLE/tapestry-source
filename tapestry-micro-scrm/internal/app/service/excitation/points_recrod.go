package excitationSvc

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/http/request"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	excitationModels "micro-scrm/internal/app/models/excitation"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"strconv"
)

type PointsRecordService struct {
	service.Base
	PointsRecordDao *excitationDao.PointsRecord
}

func NewPointsRecordService() *PointsRecordService {
	return &PointsRecordService{
		PointsRecordDao: excitationDao.NewPointsRecord(),
	}
}

// GetRecordList 获取导购积分获取记录列表
func (p *PointsRecordService) GetRecordList(c ctx.Context, params excitationReq.RecordListReq) (resList []excitationResp.PointsRecordList, count int64, err error) {

	searchExists, memberIdArr := p.Base.GetSearchMemberIds(c, params.MemberSearchReq)
	if !searchExists {
		// 如果检索条件查不到用户直接返回空数据
		return resList, count, nil
	}
	params.MemberIdArr = memberIdArr
	params.HasCoins = true
	// 查询结果
	list, count, err := p.PointsRecordDao.GetPagination(c, params)
	if err != nil {
		return resList, 0, err
	}
	// 整理结果中的memberId用户查询导购信息
	var memberIds []int64
	for _, record := range list {
		memberIds = append(memberIds, record.MemberID)
	}
	memberIds = helper.ArrayUniq(memberIds)
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
	// 整理最终返回的数据
	for _, record := range list {
		if _, ok := memberMap[record.MemberID]; ok {
			points, _ := strconv.ParseFloat(record.Points, 64)
			tmpCoins, _ := strconv.ParseFloat(record.CoinsRecord.Coins, 64)
			resList = append(resList, excitationResp.PointsRecordList{
				Id:               record.ID,
				IntegralRuleName: record.IntegralRuleName,
				IncentiveType:    record.IncentiveType,
				RuleCategoryID:   record.Type,
				RuleCategoryName: helper.RuleThree(record.IntegralRuleDetailName != "", record.IntegralRuleDetailName, p.GetRuleCategoryName(record.Type)),
				PointsRemark:     record.PointsRemark,
				FyTimeStr:        record.FyTimeStr,
				Points:           points,
				SourceType:       record.SourceType,
				CreatedAt:        record.CreatedAt,
				MemberInfo:       memberMap[record.MemberID],
				CreatedTimeStr:   helper.GetFormatDateByParams(record.CreatedAt),
				IncentiveTypeStr: getIncentiveTypeStr(record.IncentiveType),
				SourceTypeStr:    getPointsSourceTypeStr(record.SourceType),
				Coins:            tmpCoins,
			})
		}
	}
	return resList, count, nil
}

func (p *PointsRecordService) GetPointsList(c ctx.Context, params excitationReq.GetPointsListReq) (resList []excitationResp.GetPointsListResp, count int64, err error) {
	// 查询结果
	paramsWhere := excitationReq.RecordListReq{
		MemberSearchReq: request.MemberSearchReq{
			EnterpriseId: params.EnterpriseId,
			MemberIdArr: []int64{
				params.MemberId,
			},
		},
	}
	list, count, err := p.PointsRecordDao.GetPagination(c, paramsWhere)
	if err != nil {
		return resList, 0, err
	}

	// 整理最终返回的数据
	for _, record := range list {
		points, _ := strconv.ParseFloat(record.Points, 64)
		resList = append(resList, excitationResp.GetPointsListResp{
			Id:            record.ID,
			PointsRemark:  record.PointsRemark,
			Points:        fmt.Sprintf("%0.1f", points),
			CreatedAt:     record.CreatedAt,
			IncentiveType: record.IncentiveType,
		})
	}
	return resList, count, nil

}

func (p *PointsRecordService) GetRuleCategoryName(categoryId int64) string {
	var ruleTypeMap = make(map[int64]string)
	ruleTypeMap[excitationModels.RuleCategoryId1] = "邀约数（每次）"
	ruleTypeMap[excitationModels.RuleCategoryId2] = "邀约到店率（每周）"
	ruleTypeMap[excitationModels.RuleCategoryId3] = "邀约成功率（每月）"
	ruleTypeMap[excitationModels.RuleCategoryId4] = "邀约成功（每次）"
	ruleTypeMap[excitationModels.RuleCategoryId5] = "邀约占比（每月）"
	ruleTypeMap[excitationModels.RuleCategoryId6] = "邀约到店确认（每次）"
	ruleTypeMap[excitationModels.RuleCategoryId7] = "任务完成（每次）"
	ruleTypeMap[excitationModels.RuleCategoryId8] = "任务完成率（每月）"
	ruleTypeMap[excitationModels.RuleCategoryId9] = "任务触达（每人）"
	ruleTypeMap[excitationModels.RuleCategoryId10] = "新增好友数（每人）"
	ruleTypeMap[excitationModels.RuleCategoryId11] = "客户流失数（每月）"
	ruleTypeMap[excitationModels.RuleCategoryId12] = "聊天总数（每月）"
	ruleTypeMap[excitationModels.RuleCategoryId13] = "已回复聊天占比（每月）"
	ruleTypeMap[excitationModels.RuleCategoryId14] = "完成播种（每次）"
	ruleTypeMap[excitationModels.RuleCategoryId15] = "培训积分等比兑换"
	ruleTypeMap[excitationModels.RuleCategoryId16] = "线上销售业绩率（每月）"
	ruleTypeMap[excitationModels.RuleCategoryId17] = "点击新老客按钮"
	ruleTypeMap[excitationModels.RuleCategoryId18] = "点击邀约否到店购买"
	ruleTypeMap[excitationModels.RuleCategoryId19] = "线上销售金额（每月）"
	if _, ok := ruleTypeMap[categoryId]; ok {
		return ruleTypeMap[categoryId]
	}
	return ""
}

func getIncentiveTypeStr(incentiveType int64) string {
	var incentiveTypeMap = make(map[int64]string)
	incentiveTypeMap[excitationModels.IncentiveConfigTypeInvite] = "邀约维度"
	incentiveTypeMap[excitationModels.IncentiveConfigTypeTask] = "任务维度"
	incentiveTypeMap[excitationModels.IncentiveConfigTypeInteract] = "客户互动维度"
	incentiveTypeMap[excitationModels.IncentiveConfigTypeTrain] = "培训维度"
	incentiveTypeMap[excitationModels.IncentiveConfigTypeOnline] = "线上能力"
	if _, ok := incentiveTypeMap[incentiveType]; ok {
		return incentiveTypeMap[incentiveType]
	}
	return ""
}

func getPointsSourceTypeStr(SourceType int64) string {
	var sourceTypeMap = make(map[int64]string)
	sourceTypeMap[microScrmModel.ExcitationPointsRecordAutoType_Auto] = "规则统计"
	sourceTypeMap[microScrmModel.ExcitationPointsRecordAutoType_Manual] = "手动补发"
	sourceTypeMap[microScrmModel.ExcitationPointsRecordAutoType_Medal] = "奖章奖励"
	if _, ok := sourceTypeMap[SourceType]; ok {
		return sourceTypeMap[SourceType]
	}
	return ""
}
