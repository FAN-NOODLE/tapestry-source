package dataCompassSvc

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	workDao "micro-scrm/internal/app/dao/work"
	dataCompassReq "micro-scrm/internal/app/http/request/datacompass"
	dataCompassResp "micro-scrm/internal/app/http/response/datacompass"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
)

type ChatGroup struct {
	service.Base
}

func NewChatGroup() *ChatGroup {
	return &ChatGroup{}
}

func (s *ChatGroup) List(c app.Ctx, req dataCompassReq.ChatGroupReq) (any, int64, error) {
	// 查询所有所有群组
	groupChatList, total, err := workDao.NewGroupChat().GetList(c, req.PageNum, req.PageSize, map[string]any{
		"enterprise_id": req.EnterpriseID,
		"status":        1,
		"type":          1,
	})
	if err != nil {
		return nil, 0, err
	}

	var respChatGroup []dataCompassResp.ChatGroupResp

	var chatIds []int64
	for _, groupChat := range groupChatList {
		tmpResp := dataCompassResp.ChatGroupResp{
			GroupChatId:          groupChat.ID,
			GroupName:            groupChat.Name,
			LeaveGroupPercentage: "0%",
		}
		respChatGroup = append(respChatGroup, tmpResp)
		chatIds = append(chatIds, groupChat.ID)
	}
	if err != nil {
		return nil, 0, err
	}

	var (
		totalNumbers  map[int64]int64
		addNumbers    map[int64]int64
		leaveNumbers  map[int64]int64
		memberNumbers map[int64]int64
	)

	wg := threading.NewRoutineGroup()

	// 查询群累计人数
	wg.RunSafe(func() {
		totalNumbers, err = workDao.NewWorkGroupChatMemberCk().SelectTotalCount(c, chatIds, req)
	})

	// 群新增人数
	wg.RunSafe(func() {
		addNumbers, err = workDao.NewWorkGroupChatMemberCk().SelectAddCount(c, chatIds, req)
	})

	// 退群人数
	wg.RunSafe(func() {
		leaveNumbers, err = workDao.NewWorkGroupChatMemberCk().SelectLossCount(c, chatIds, req)
	})

	// 获取群组导购数量
	wg.RunSafe(func() {
		memberNumbers, err = workDao.NewWorkGroupChatMemberCk().SelectMemberCount(c, chatIds, req)
	})

	wg.Wait()

	for k, groupChat := range respChatGroup {
		// 群成员总数
		if val, ok := totalNumbers[groupChat.GroupChatId]; ok {
			respChatGroup[k].CumulativeGroupNumber = val
			groupChat.CumulativeGroupNumber = val
		}

		// 群新增人数
		if val, ok := addNumbers[groupChat.GroupChatId]; ok {
			respChatGroup[k].NumberOfNewcomers = val
		}

		// 退群人数
		if val, ok := leaveNumbers[groupChat.GroupChatId]; ok {
			respChatGroup[k].LeaveGroupNumber = val
			groupChat.LeaveGroupNumber = val
		}

		// 获取群组导购数量
		if val, ok := memberNumbers[groupChat.GroupChatId]; ok {
			respChatGroup[k].BANumber = val
		}

		leaTotal := groupChat.CumulativeGroupNumber + groupChat.LeaveGroupNumber
		if leaTotal != 0 && groupChat.LeaveGroupNumber != 0 {
			respChatGroup[k].LeaveGroupPercentage = mathRounded(float64(groupChat.LeaveGroupNumber), float64(leaTotal))
		} else {
			respChatGroup[k].LeaveGroupPercentage = "0%"
		}
	}

	return respChatGroup, total, nil
}

func mathRounded(params1, params2 float64) string {
	if params1 == 0 {
		return "0%"
	}
	params3 := helper.MathRounded((params1/params2)*100, 2)
	if params3 == float64(int64(params3)) {
		return fmt.Sprintf("%.0f", params3) + "%"
	}
	return fmt.Sprintf("%.2f", params3) + "%"
}
