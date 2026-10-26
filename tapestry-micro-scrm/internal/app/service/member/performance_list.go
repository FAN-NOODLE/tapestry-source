package membersSvc

import (
	"fmt"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/dao/micro_scrm"
	memberReq "micro-scrm/internal/app/http/request/member"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"sort"
	"strconv"
)

type PerformanceListService struct {
	service.Base
	PerformanceListDao *membersDao.MemberExternalBehaviorData
}

func NewNewPerformanceListService() *PerformanceListService {
	return &PerformanceListService{
		PerformanceListDao: membersDao.NewMemberExternalBehaviorData(),
	}
}

// ActiveCustomers 活跃榜顾客榜单

func (s *PerformanceListService) ActiveCustomers(c app.Ctx, request memberReq.PerformanceListReq) ([]memberReq.PerformanceListRes, int64, error) {

	var (
		total  int64
		result []memberReq.PerformanceListRes
		err    error
	)

	members, err := membersDao.NewMembers().GetMemberByThirdIds(c, request.EnterpriseId, request.ThirdId, []int64{})
	if err != nil {
		return nil, total, err
	}

	memberIdMap := make(map[int64]memberReq.MemberStaffRes, 10)
	var membersID []int64
	for _, value := range members {
		membersID = append(membersID, value.MembersId)
		memberIdMap[value.MembersId] = value
	}

	res, total, err := s.PerformanceListDao.ActiveCustomers(c, request, membersID)
	if err != nil {
		return nil, total, err
	}
	for _, item := range res {
		if v, ok := memberIdMap[item.MemberID]; ok {
			replyPercentageFloat, _ := strconv.ParseFloat(item.ReplyPercentage, 10)
			newResult := memberReq.PerformanceListRes{
				EnterpriseId:    request.EnterpriseId,
				Avatar:          v.MembersAvatar,
				Name:            v.MembersName,
				StoreName:       v.StoreName,
				MemberID:        v.MembersId,
				ChatCnt:         item.ChatCnt,
				MessageCnt:      item.MessageCnt,
				ReplyPercentage: fmt.Sprintf("%0.f", replyPercentageFloat),
				NewContactCnt:   item.NewContactCnt,
			}
			result = append(result, newResult)
		}
	}

	return result, total, err
}

// InvitationV2 优化邀约榜单

func (s *PerformanceListService) Invitation(c app.Ctx, request memberReq.PerformanceListInvitationReq) ([]memberReq.PerformanceListInvitationRes, int64, error) {
	var (
		total int64

		err error
	)
	result := make([]memberReq.PerformanceListInvitationRes, 0)

	members, err := membersDao.NewMembers().GetMemberByThirdIds(c, request.EnterpriseId, request.ThirdId, []int64{})
	if err != nil {
		return result, total, err
	}

	memberIdMap := make(map[int64]memberReq.MemberStaffRes)
	var membersID []int64
	for _, value := range members {
		membersID = append(membersID, value.MembersId)
		memberIdMap[value.MembersId] = value
	}

	many, err := departmentsDao.NewDepartments().GetMany(c, map[string]interface{}{
		"enterprise_id": request.EnterpriseId,
		"third_id":      request.ThirdId,
		"deleted_at":    0,
	})
	if err != nil {
		return result, 0, err
	}
	storeCodes := make([]string, 0)
	for _, model := range many {
		storeCodes = append(storeCodes, model.Name)
	}

	res, total, err := membersDao.NewMemberInvitationList().InvitationCk(c, request, membersID, storeCodes)
	if err != nil {
		res, total, err = membersDao.NewMemberInvitationList().Invitation(c, request, membersID, storeCodes)
	}

	// res, total, err := membersDao.NewMemberInvitationList().Invitation(c, request, membersID)
	// 返回数据

	resultMemberIds := make([]int64, 0)
	for _, item := range res {
		resultMemberIds = append(resultMemberIds, item.MemberID)
	}
	resMember, err := membersDao.NewMembers().GetMemberByThirdIds(c, request.EnterpriseId, nil, resultMemberIds)
	if err != nil {
		return result, 0, err
	}
	resMemberMap := make(map[int64]memberReq.MemberStaffRes)
	for _, value := range resMember {
		resMemberMap[value.MembersId] = value
	}
	notInMemberMapNum := 0
	for _, item := range res {
		var newResult memberReq.PerformanceListInvitationRes
		isView := 0
		if v, ok := memberIdMap[item.MemberID]; ok { // 不存在门店的导购数据
			if request.UserType == 4 {
				isView = 1
			}
			if v.Mobile == request.Mobile {
				isView = 1
			}
			newResult = memberReq.PerformanceListInvitationRes{
				EnterpriseId: request.EnterpriseId,
				Avatar:       v.MembersAvatar,
				Name:         v.MembersName,
				StoreName:    v.StoreName,
				MemberID:     item.MemberID,
				Num:          item.Num,
				SuccessNum:   item.SuccessNum,
				PriceTotal:   item.PriceTotal,
				IsView:       isView,
			}
			result = append(result, newResult)
		} else {
			notInMemberMapNum++
			isView = 1
			if v, ok = resMemberMap[item.MemberID]; ok {
				// 查询导购的交易邀约单号
				reserveNumbers, err := microScrmDao.NewReservesOrderSales().GetReserveNumberForSales(c, request.EnterpriseId, "member_user_id=? and pay_time BETWEEN ? AND ? ", v.MembersUserid, request.FinancialTime[0], request.FinancialTime[1])
				if err != nil {
					continue
				}
				newResult = memberReq.PerformanceListInvitationRes{
					EnterpriseId:   request.EnterpriseId,
					Avatar:         v.MembersAvatar,
					Name:           v.MembersName,
					StoreName:      v.StoreName,
					MemberID:       item.MemberID,
					Num:            item.Num,
					SuccessNum:     item.SuccessNum,
					PriceTotal:     item.PriceTotal,
					IsView:         isView,
					ReserveNumbers: reserveNumbers,
				}
				result = append(result, newResult)
			}
		}
	}

	// 如果res里的数据长度比Size小说明分页数据不够了直接查邀约榜规定时间内的导购id
	if len(res) < request.Size {

		// 声明补充的用户信息
		var resultAdd []memberReq.PerformanceListInvitationRes

		orderMemberId, err := membersDao.NewMemberInvitationList().InvitationMemberIdsCk(c, request, membersID, storeCodes)
		if err != nil {
			orderMemberId, err = membersDao.NewMemberInvitationList().InvitationMemberIds(c, request, membersID, storeCodes)
		}
		// 循环导购memberId找到没有出现在邀约单里的导购memberId没有出现在导购里的
		differMemberId := helper.ArrayDiff(membersID, orderMemberId)
		// 先把这些用户全部加到结果里排完序在处理分页
		for _, item := range members {
			if helper.InArray(item.MembersId, differMemberId) {
				isView := 0
				if request.UserType == 4 {
					isView = 1
				}
				if item.Mobile == request.Mobile {
					isView = 1
				}
				resultAdd = append(resultAdd, memberReq.PerformanceListInvitationRes{
					EnterpriseId: request.EnterpriseId,
					Avatar:       item.MembersAvatar,
					Name:         item.MembersName,
					StoreName:    item.StoreName,
					MemberID:     item.MembersId,
					IsView:       isView,
				})
			}
		}

		// 对差级的导购memberID由大到小排序
		sort.Slice(resultAdd, func(i, j int) bool {
			return resultAdd[i].MemberID < resultAdd[j].MemberID
		})

		// 按照分页返回数据开始
		startKey := request.Page*request.Size - int(total) - request.Size
		if startKey < 0 {
			startKey = 0
		}
		if startKey < len(resultAdd) {
			// 按照分页返回数据结束
			endKey := request.Page*request.Size - int(total)
			if endKey > len(resultAdd) {
				endKey = len(resultAdd)
			}
			resultAdd = resultAdd[startKey:endKey]
			result = append(result, resultAdd...)
		}
	}

	return result, int64(len(members) + notInMemberMapNum), err

}

// Achievement 业绩榜单
func (s *PerformanceListService) Achievement(c app.Ctx, request memberReq.PerformanceListAchievementReq, start int) ([]memberReq.PerformanceListAchievementRes, int64, error) {

	var (
		total  int64
		result []memberReq.PerformanceListAchievementRes
		err    error
	)

	members, err := membersDao.NewMembers().GetMemberByThirdIds(c, request.EnterpriseId, request.ThirdId, []int64{})
	if err != nil {
		return nil, total, err
	}

	memberIdMap := make(map[int64]memberReq.MemberStaffRes, 10)
	var membersID []int64
	var addMembersID []int64
	for _, value := range members {
		// 这个字段用来手动分页提取用户信息
		addMembersID = append(addMembersID, value.MembersId)
		// start为1的时候返回的是明星榜单
		if start == 0 {
			membersID = append(membersID, value.MembersId)
		}
		memberIdMap[value.MembersId] = value
	}
	res, total, err := microScrmDao.NewOrderAchievement().Achievement(c, request, membersID)
	if err != nil {
		return result, total, err
	}
	for _, item := range res {
		if v, ok := memberIdMap[item.MemberID]; ok {
			newResult := memberReq.PerformanceListAchievementRes{
				EnterpriseId:     request.EnterpriseId,
				Avatar:           v.MembersAvatar,
				Name:             v.MembersName,
				StoreName:        v.StoreName,
				MemberID:         v.MembersId,
				OrderTotal:       item.OrderTotal,
				OrderPriceTotal:  item.OrderPriceTotal,
				RefundOrderTotal: item.RefundOrderTotal,
				PayOrderTotal:    item.PayOrderTotal,
			}
			result = append(result, newResult)
		}
	}

	// 如果res里的数据长度比Size小说明分页数据不够了直接查邀约榜规定时间内的导购id
	if len(res) < request.Size && start == 0 {
		// 声明补充的用户信息
		var resultAdd []memberReq.PerformanceListAchievementRes

		orderMemberId, _ := microScrmDao.NewOrderAchievement().InsertAchievementMemberIds(c, request, membersID)
		// 循环导购memberId找到没有出现在邀约单里的导购memberId没有出现在导购里的
		differMemberId := helper.ArrayDiff(addMembersID, orderMemberId)
		// 先把这些用户全部加到结果里排完序在处理分页
		for _, item := range members {
			if helper.InArray(item.MembersId, differMemberId) {
				resultAdd = append(resultAdd, memberReq.PerformanceListAchievementRes{
					EnterpriseId: request.EnterpriseId,
					Avatar:       item.MembersAvatar,
					Name:         item.MembersName,
					StoreName:    item.StoreName,
					MemberID:     item.MembersId,
				})
			}
		}

		// 对差级的导购memberID由大到小排序
		sort.Slice(resultAdd, func(i, j int) bool {
			return resultAdd[i].MemberID < resultAdd[j].MemberID
		})

		// 按照分页返回数据开始
		startKey := request.Page*request.Size - int(total) - request.Size
		if startKey < 0 {
			startKey = 0
		}
		if startKey < len(resultAdd) {
			// 按照分页返回数据结束
			endKey := request.Page*request.Size - int(total)
			if endKey > len(resultAdd) {
				endKey = len(resultAdd)
			}
			resultAdd = resultAdd[startKey:endKey]
			result = append(result, resultAdd...)
		}
	}

	return result, int64(len(members)), err
}

// SwCustomers sw顾客榜单
func (s *PerformanceListService) SwCustomers(c app.Ctx, request memberReq.SwCustomersListReq, start int) ([]memberReq.SwCustomersListReqRes, int64, error) {
	var (
		total  int64
		result []memberReq.SwCustomersListReqRes
		err    error
	)

	members, err := membersDao.NewMembers().GetMemberByThirdIds(c, request.EnterpriseId, request.ThirdId, []int64{})
	if err != nil {
		return result, total, err
	}

	memberIdMap := make(map[int64]memberReq.MemberStaffRes, 10)
	var membersID []int64
	for _, value := range members {
		// start为1的时候返回的是明星榜单
		if start == 0 {
			membersID = append(membersID, value.MembersId)
		}
		memberIdMap[value.MembersId] = value
	}

	// 调整顾客榜的查询方式
	// res, total, err := microScrmDao.NewExternalEverydayStatisticsDao().GetSwCustomers(c, request, membersID)
	// 查询业务 TODO增加ck查询
	res, total, err := microScrmDao.NewUserRelation().SwCustomersCk(c, request, membersID)
	if err != nil {
		res, total, err = microScrmDao.NewUserRelation().SwCustomers(c, request, membersID)
	}
	if err != nil {
		return result, total, err
	}
	for _, item := range res {
		if v, ok := memberIdMap[item.MemberID]; ok {
			newResult := memberReq.SwCustomersListReqRes{
				EnterpriseId:   request.EnterpriseId,
				Avatar:         v.MembersAvatar,
				Name:           v.MembersName,
				StoreName:      v.StoreName,
				MemberID:       v.MembersId,
				NewCustomers:   item.NewCustomers,
				OldCustomers:   item.OldCustomers,
				SleepCustomers: item.SleepCustomers,
				VipCustomers:   item.VipCustomers,
			}
			result = append(result, newResult)
		}
	}

	return result, total, err
}

// SwActive sw活跃榜
func (s *PerformanceListService) SwActive(c app.Ctx, request memberReq.SwActiveListReqReq, start int) ([]memberReq.SwActiveListReqRes, int64, error) {
	var (
		total  int64
		result []memberReq.SwActiveListReqRes
		err    error
		// isPush int64
	)

	startTime := request.FinancialTime[0]
	endTime := request.FinancialTime[1]

	// 查询导购信息
	members, err := membersDao.NewMembers().GetMemberByThirdIds(c, request.EnterpriseId, request.ThirdId, []int64{})
	if err != nil {
		return result, total, err
	}

	var membersID []int64
	var membersUserId []string
	for _, value := range members {
		// start为1的时候返回的是明星榜单
		if start == 0 {
			membersID = append(membersID, value.MembersId)
			membersUserId = append(membersUserId, value.MembersUserid)
		}
	}

	// 先查询任务记录

	taskRes, err := microScrmDao.NewPlanTaskAllotLogDao().GetAllMemberTaskRate(c, request.EnterpriseId, membersUserId, startTime, endTime)
	if err != nil {
		return result, total, err
	}
	// 在查询企微统计数据
	res, err := s.PerformanceListDao.SwActive(c, request, membersID)
	if err != nil {
		return result, total, err
	}
	// 循环组装数据
	for _, item := range members {
		newResult := memberReq.SwActiveListReqRes{
			EnterpriseId:    request.EnterpriseId,
			Avatar:          item.MembersAvatar,
			Name:            item.MembersName,
			StoreName:       item.StoreName,
			MemberID:        item.MembersId,
			ReplyPercentage: "0", // 默认的回复占比为0
		}

		// 找到了任务数据
		if v, ok := taskRes[item.MembersUserid]; ok {
			newResult.CompletionRate = v
			// isPush = 1
		}
		// 找到了企微统计数据
		if v, ok := res[item.MembersId]; ok {
			replyPercentageFloat, _ := strconv.ParseFloat(v.ReplyPercentage, 10)
			newResult.ChatCnt = v.ChatCnt
			newResult.MessageCnt = v.MessageCnt
			newResult.ReplyPercentage = fmt.Sprintf("%0.f", replyPercentageFloat)
			newResult.NewContactCnt = v.NewContactCnt
			// isPush = 1
		}
		// 追加数据 这里修改了要展示所有的导购所有注释了判断条件
		/*if isPush > 0 {
			result = append(result, newResult)
		}*/
		result = append(result, newResult)
	}
	// 排序
	sort.Slice(result, func(i, j int) bool {
		if result[i].CompletionRate == result[j].CompletionRate {
			return result[i].ChatCnt > result[j].ChatCnt
		} else {
			return result[i].CompletionRate > result[j].CompletionRate
		}
	})
	total = int64(len(result))
	// 手动分页
	// 按照分页返回数据开始
	startKey := (request.Page - 1) * request.Size
	endKey := request.Page * request.Size
	if startKey > int(total) {
		return []memberReq.SwActiveListReqRes{}, total, err
	}
	if endKey > int(total) {
		result = result[startKey:]
	} else {
		result = result[startKey:endKey]
	}
	return result, total, err
}

// 明星榜-绩效榜列表

func (s *PerformanceListService) StarActiveCustomers(c app.Ctx, request memberReq.PerformanceListReq) ([]memberReq.PerformanceListRes, error) {

	var result []memberReq.PerformanceListRes
	var membersId []int64
	res, _, err := s.PerformanceListDao.ActiveCustomers(c, request, membersId)
	if err != nil {
		return nil, err
	}

	for _, item := range res {
		membersId = append(membersId, item.MemberID)
	}

	members, err := membersDao.NewMembers().GetMemberByThirdIds(c, request.EnterpriseId, []int{}, membersId)
	if err != nil {
		return nil, err
	}

	memberIdMap := make(map[int64]memberReq.MemberStaffRes, len(members))
	for _, value := range members {
		memberIdMap[value.MembersId] = value
	}

	for _, item := range res {
		if v, ok := memberIdMap[item.MemberID]; ok {
			replyPercentageFloat, _ := strconv.ParseFloat(item.ReplyPercentage, 10)
			newResult := memberReq.PerformanceListRes{
				EnterpriseId:    request.EnterpriseId,
				Avatar:          v.MembersAvatar,
				Name:            v.MembersName,
				StoreName:       v.StoreName,
				MemberID:        v.MembersId,
				ChatCnt:         item.ChatCnt,
				MessageCnt:      item.MessageCnt,
				ReplyPercentage: fmt.Sprintf("%0.f", replyPercentageFloat),
				NewContactCnt:   item.NewContactCnt,
			}
			result = append(result, newResult)
		}
	}

	return result, err
}

// 明星榜-邀约榜

func (s *PerformanceListService) StarInvitation(c app.Ctx, request memberReq.PerformanceListInvitationReq) ([]memberReq.PerformanceListInvitationRes, error) {
	var result []memberReq.PerformanceListInvitationRes
	var membersId []int64

	// 先查ck如果ck异常就查mysql
	res, _, err := membersDao.NewMemberInvitationList().InvitationCk(c, request, membersId, nil)
	if err != nil {
		// mysql查询
		res, _, err = membersDao.NewMemberInvitationList().Invitation(c, request, membersId, nil)
	}

	if err != nil {
		return nil, err
	}

	for _, item := range res {
		membersId = append(membersId, item.MemberID)
	}

	members, err := membersDao.NewMembers().GetMemberByThirdIds(c, request.EnterpriseId, []int{}, membersId)
	if err != nil {
		return nil, err
	}

	memberIdMap := make(map[int64]memberReq.MemberStaffRes, len(members))
	for _, value := range members {
		memberIdMap[value.MembersId] = value
	}

	for _, item := range res {
		if v, ok := memberIdMap[item.MemberID]; ok {
			isView := 0
			newResult := memberReq.PerformanceListInvitationRes{
				EnterpriseId: request.EnterpriseId,
				Avatar:       v.MembersAvatar,
				Name:         v.MembersName,
				StoreName:    v.StoreName,
				MemberID:     item.MemberID,
				Num:          item.Num,
				SuccessNum:   item.SuccessNum,
				PriceTotal:   item.PriceTotal,
				IsView:       isView,
			}
			result = append(result, newResult)
		}
	}

	return result, err
}
