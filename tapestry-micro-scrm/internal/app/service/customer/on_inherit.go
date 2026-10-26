package customerSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	membersDao "micro-scrm/internal/app/dao/members"
	externalUserDao "micro-scrm/internal/app/dao/micro_scrm"
	customerReq "micro-scrm/internal/app/http/request/customer"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
)

type OnInherit struct {
	service.Base
	workExternalDao *externalUserDao.UserRelation
}

func NewOnInherit() *OnInherit {
	return &OnInherit{
		workExternalDao: externalUserDao.NewUserRelation()}
}

func (s OnInherit) GetTakeoverMemberList(c ctx.Context, req customerReq.OnInheritReq) ([]int64, error) {
	memberIds := make([]int64, 0)
	if len(req.ThirdIds) > 0 {
		list, err := membersDao.NewMemberDepartments().GetMemberIdsByThirdIds(c, req.EnterpriseID, req.ThirdIds, []int64{})
		if err != nil {
			return []int64{}, err
		}
		memberIds = list
	}
	memberIds = append(memberIds, req.MemberId)
	if len(req.MemberList) > 0 {
		memberIds = helper.ArrayIntersect(memberIds, req.MemberList)
	}
	userIds, err := membersDao.NewMembers().GetUserIdsById(c, req.EnterpriseID, memberIds)
	if err != nil {
		return nil, err
	}
	canOnInheritMemberIds, err := s.workExternalDao.GetOnInheritMemberIds(c, req, userIds)
	if err != nil {
		return []int64{}, err
	}
	intersect := helper.ArrayIntersect(canOnInheritMemberIds, memberIds)
	return intersect, err
}

// GetSelectAllNum 获取在职继承全选的人数
func (s OnInherit) GetSelectAllNum(c ctx.Context, request customerReq.OnInheritReq) (int64, error) {

	externalIds, memberIds := s.getSelectAllExternalIdAndMemberId(c, request)
	if len(request.MemberList) > 0 {
		memberIds = helper.ArrayIntersect(memberIds, request.MemberList)
	}
	userIds, err := membersDao.NewMembers().GetUserIdsById(c, request.EnterpriseID, memberIds)
	if err != nil {
		return 0, err
	}
	num, err := s.workExternalDao.GetOnInheritSelectAllNum(c, request, externalIds, memberIds, userIds)
	if err != nil {
		return num, err
	}
	return num, err
}

// TransferCount    获取在职继承全选的人数
func (s OnInherit) TransferCount(c ctx.Context, request customerReq.OnInheritSearchReq) (int64, error) {
	var total int64
	request.FilterExternalUserId = s.getSelectAllExternalId(c, request)
	query := externalSearchDao.NewExternalExtraCk().GetRelationShipQuery(c, request.SearchList)
	err := query.Count(&total).Error
	return total, err
}

// getSelectAllExternalIdAndMemberId 在职继承全选的统一代码封装
func (s OnInherit) getSelectAllExternalIdAndMemberId(c ctx.Context, request customerReq.OnInheritReq) ([]int64, []int64) {
	// 1、先获取最近90天内分配次数超过两次的外部联系人
	// 2、获取在职分配的状态为等待接替的外部联系人
	externalIds := make([]int64, 0)
	memberIds := make([]int64, 0)
	startTime, _ := helper.YesterdayPeriodOfTimestamp(90)
	limitExternalIds, err := externalUserDao.NewCustomerOnInheritRecord().GetLimitExternalIds(c, request.EnterpriseID, startTime)
	if err != nil {
		return externalIds, memberIds
	}
	for _, v := range limitExternalIds {
		externalIds = append(externalIds, v.ExternalID)
	}
	startTime, _ = helper.YesterdayPeriodOfTimestamp(1)
	waitExternalIds, err := externalUserDao.NewCustomerOnInheritRecord().GetWaitExternalIds(c, request.EnterpriseID, startTime)
	if err != nil {
		return externalIds, memberIds
	}
	externalIds = append(externalIds, waitExternalIds...)

	if len(request.ThirdIds) > 0 {
		list, err := membersDao.NewMemberDepartments().GetMemberIdsByThirdIds(c, request.EnterpriseID, request.ThirdIds, []int64{})
		if err != nil {
			return externalIds, memberIds
		}
		memberIds = list
	}
	memberIds = append(memberIds, request.MemberId)
	return externalIds, memberIds
}

// getSelectAllExternalIdAndMemberId 在职继承全选的统一代码封装
func (s OnInherit) getSelectAllExternalId(c ctx.Context, request customerReq.OnInheritSearchReq) []int64 {
	// 1、先获取最近90天内分配次数超过两次的外部联系人
	// 2、获取在职分配的状态为等待接替的外部联系人
	externalIds := make([]int64, 0)
	startTime, _ := helper.YesterdayPeriodOfTimestamp(90)
	limitExternalIds, err := externalUserDao.NewCustomerOnInheritRecord().GetLimitExternalIds(c, request.EnterpriseID, startTime)
	if err != nil {
		return externalIds
	}
	for _, v := range limitExternalIds {
		externalIds = append(externalIds, v.ExternalID)
	}
	startTime, _ = helper.YesterdayPeriodOfTimestamp(1)
	waitExternalIds, err := externalUserDao.NewCustomerOnInheritRecord().GetWaitExternalIds(c, request.EnterpriseID, startTime)
	if err != nil {
		return externalIds
	}
	return append(externalIds, waitExternalIds...)
}

func (s OnInherit) GetSelectAllList(c ctx.Context, request customerReq.OnInheritReq) ([]customerReq.OnInheritExternalInfoRes, error) {
	// 1、先获取最近90天内分配次数超过两次的外部联系人
	// 2、获取在职分配的状态为等待接替的外部联系人
	externalIds, memberIds := s.getSelectAllExternalIdAndMemberId(c, request)
	if len(request.MemberList) > 0 {
		memberIds = helper.ArrayIntersect(memberIds, request.MemberList)
	}
	userIds, err := membersDao.NewMembers().GetUserIdsById(c, request.EnterpriseID, memberIds)
	if err != nil {
		return nil, err
	}
	result, err := s.workExternalDao.GetOnInheritSelectList(c, request, externalIds, memberIds, userIds)
	if err != nil {
		return result, err
	}
	return result, err
}

func (s OnInherit) TransferList(c ctx.Context, request customerReq.OnInheritSearchReq) ([]customerReq.OnInheritExternalInfoRes, error) {
	var data = make([]customerReq.OnInheritExternalInfoRes, 0)
	request.FilterExternalUserId = s.getSelectAllExternalId(c, request)
	query := externalSearchDao.NewExternalExtraCk().GetRelationShipQuery(c, request.SearchList)
	err := query.Select("enterprise_id,member_id,work_external_user_id,user_id").Find(&data).Error
	return data, err
}
