package customerSvc

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	workGroupDao "micro-scrm/internal/app/dao/work"
	customerReq "micro-scrm/internal/app/http/request/customer"
	"micro-scrm/internal/app/service"
)

type DataCompassSvc struct {
	service.Base
	workGroupMemberLogDao *workGroupDao.WorkGroupChatMemberLog
}

var _ DataCompassSvc

func NewDataCompassSvc() *DataCompassSvc {
	return &DataCompassSvc{
		workGroupMemberLogDao: workGroupDao.NewWorkGroupChatMemberLog(),
	}
}

// ChatGroupCompass 社群数据的统计
func (s *DataCompassSvc) ChatGroupCompass(c ctx.Context, request customerReq.CompassReq) (customerReq.ChatGroupCompassRes, error) {
	var res customerReq.ChatGroupCompassRes
	memberInfo, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{"enterprise_id": request.EnterpriseId, "mobile": request.Mobile})
	if err != nil {
		return res, err
	}
	var memberIds []int64
	memberIds = append(memberIds, memberInfo.ID)
	if request.IsSM == 1 { //店长身份查看数据罗盘统计
		if request.SearchAll == 1 { //查询管理的全部导购信息
			memberInfos, err := membersDao.NewMembers().GetMemberByThirdIds(c, request.EnterpriseId, request.ThirdIds, []int64{})
			if err != nil {
				return res, err
			}
			for _, info := range memberInfos {
				memberIds = append(memberIds, info.MembersId)
			}
		} else {
			if len(request.MemberIds) > 0 {
				memberIds = request.MemberIds
			}
			if len(request.ThirdIds) > 0 {
				memberInfos, err := membersDao.NewMembers().GetMemberByThirdIds(c, request.EnterpriseId, request.ThirdIds, []int64{})
				if err != nil {
					return res, err
				}
				for _, info := range memberInfos {
					memberIds = append(memberIds, info.MembersId)
				}
			}
		}
	}
	groupChatIds, err := workGroupDao.NewWorkGroupChatMember().GetChatIdsByMemberIds(c, request.EnterpriseId, request.StartTime, request.EndTime, memberIds)
	if err != nil {
		return res, err
	}
	if len(groupChatIds) == 0 {
		return res, err
	}
	res.GroupNumber = int64(len(groupChatIds))
	//累计的群成员数
	GroupMemberAddFriend, err := s.workGroupMemberLogDao.GetJoinGroupUserCounts(c, request.EnterpriseId, request.StartTime, request.EndTime, groupChatIds)
	if err != nil {
		return res, err
	}
	res.CumulativeGroupNumber = GroupMemberAddFriend
	//群成员加好友数
	addFriendNumber, err := s.workGroupMemberLogDao.GetChatFriendsRelationCount(c, request.EnterpriseId, request.StartTime, request.EndTime, groupChatIds, memberIds)
	if err != nil {
		return res, err
	}
	res.AddFriendNumber = addFriendNumber
	if addFriendNumber == 0 || GroupMemberAddFriend == 0 {
		res.AddFriendPer = "0%"
	} else {
		addFriendPerFloat, _ := decimal.NewFromInt(addFriendNumber).Div(decimal.NewFromInt(GroupMemberAddFriend)).Float64()
		addFriendPer := addFriendPerFloat * 100
		res.AddFriendPer = fmt.Sprintf("%.2f", addFriendPer) + "%"
	}
	//退群人数
	leaveGroupNumber, err := s.workGroupMemberLogDao.GetLeaveGroupCounts(c, request.EnterpriseId, request.StartTime, request.EndTime, groupChatIds)
	if err != nil {
		return res, err
	}
	res.LeaveGroupNumber = leaveGroupNumber
	if leaveGroupNumber == 0 || GroupMemberAddFriend == 0 {
		res.LeaveGroupPer = "0%"
	} else {
		leaveGroupPerFloat, _ := decimal.NewFromInt(leaveGroupNumber).Div(decimal.NewFromInt(GroupMemberAddFriend)).Float64()
		leaveGroupPer := leaveGroupPerFloat * 100
		res.LeaveGroupPer = fmt.Sprintf("%.2f", leaveGroupPer) + "%"
	}
	//入群人数
	joinGroupNumber, err := s.workGroupMemberLogDao.GetJoinGroupCounts(c, request.EnterpriseId, request.StartTime, request.EndTime, groupChatIds)
	if err != nil {
		return res, err
	}
	res.JoinGroupNumber = joinGroupNumber
	return res, nil
}

// FriendsCompass 好友数据的统计
func (s *DataCompassSvc) FriendsCompass(c ctx.Context, request customerReq.CompassReq) (customerReq.FriendsCompassRes, error) {
	var res customerReq.FriendsCompassRes
	memberInfo, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{"enterprise_id": request.EnterpriseId, "mobile": request.Mobile})
	if err != nil {
		return res, err
	}
	var memberIds []int64
	memberIds = append(memberIds, memberInfo.ID)
	if request.IsSM == 1 { //店长身份查看数据罗盘统计
		if request.SearchAll == 1 { //查询管理的全部导购信息
			memberInfos, err := membersDao.NewMembers().GetMemberByThirdIds(c, request.EnterpriseId, request.ThirdIds, []int64{})
			if err != nil {
				return res, err
			}
			for _, info := range memberInfos {
				memberIds = append(memberIds, info.MembersId)
			}
		} else {
			if len(request.MemberIds) > 0 {
				memberIds = request.MemberIds
			}
			if len(request.ThirdIds) > 0 {
				memberInfos, err := membersDao.NewMembers().GetMemberByThirdIds(c, request.EnterpriseId, request.ThirdIds, []int64{})
				if err != nil {
					return res, err
				}
				for _, info := range memberInfos {
					memberIds = append(memberIds, info.MembersId)
				}
			}
		}
	}
	// 新增好友数
	addFriendsNumber, err := microScrmDao.NewUserRelation().GetAddFriendsCount(c, request.EnterpriseId, request.StartTime, request.EndTime, memberIds)
	if err != nil {
		return res, err
	}
	res.AddFriendNumber = addFriendsNumber
	// 累计好友数
	allFriendsNumber, err := microScrmDao.NewUserRelation().GetALLFriendsCount(c, request.EnterpriseId, request.EndTime, memberIds)
	if err != nil {
		return res, err
	}
	res.AllFriendsNumber = allFriendsNumber
	lossFriendsNumber, err := microScrmDao.NewUserRelation().GetLossFriendsCount(c, request.EnterpriseId, request.StartTime, request.EndTime, memberIds)
	if err != nil {
		return res, err
	}
	res.LossFriendsNumber = lossFriendsNumber
	if lossFriendsNumber == 0 || allFriendsNumber == 0 {
		res.LossFriendsPer = "0%"
	} else {
		lossFriendsPer, _ := decimal.NewFromInt(lossFriendsNumber).Div(decimal.NewFromInt(allFriendsNumber)).Float64()
		lossFriendsPerFloat := lossFriendsPer * 100
		res.LossFriendsPer = fmt.Sprintf("%.2f", lossFriendsPerFloat) + "%"
	}
	groupChatIds, err := workGroupDao.NewWorkGroupChatMember().GetChatIdsByMemberIds(c, request.EnterpriseId, request.StartTime, request.EndTime, memberIds)
	if err != nil {
		return res, err
	}
	if len(groupChatIds) == 0 {
		return res, err
	}
	// 新增好友的入群统计
	newJoinChatFriendsNumber, err := s.workGroupMemberLogDao.GetChatNewFriendsRelationCount(c, request.EnterpriseId, request.StartTime, request.EndTime, groupChatIds, memberIds)
	if err != nil {
		return res, err
	}
	if newJoinChatFriendsNumber == 0 || addFriendsNumber == 0 {
		res.NewJoinGroupPer = "0%"
	} else {
		addJoinChat, _ := decimal.NewFromInt(newJoinChatFriendsNumber).Div(decimal.NewFromInt(addFriendsNumber)).Float64()
		addJoinChatFloat := addJoinChat * 100
		res.NewJoinGroupPer = fmt.Sprintf("%.2f", addJoinChatFloat) + "%"
	}
	//入群人数
	joinGroupNumber, err := s.workGroupMemberLogDao.GetJoinGroupCounts(c, request.EnterpriseId, request.StartTime, request.EndTime, groupChatIds)
	if err != nil {
		return res, err
	}
	if joinGroupNumber == 0 || allFriendsNumber == 0 {
		res.JoinGroupPer = "0%"
	} else {
		joinChat, _ := decimal.NewFromInt(joinGroupNumber).Div(decimal.NewFromInt(allFriendsNumber)).Float64()
		joinChatFloat := joinChat * 100
		res.JoinGroupPer = fmt.Sprintf("%.2f", joinChatFloat) + "%"
	}
	res.SaNumber = len(memberIds)
	return res, err
}
