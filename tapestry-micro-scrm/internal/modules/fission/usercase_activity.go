package fission

import (
	"fmt"
	"micro-scrm/internal/pkg/helper"
	"sort"
	"strings"
	"sync"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
)

// ActivityList 裂变活动列表
func (s *Service) ActivityList(c ctx.Context, request *ActivityListRequest) (int64, []*ActivityEntity, error) {
	total, list, err := s.repo.ActivityList(c, request)
	if err != nil {
		return 0, nil, err
	}
	for key, item := range list {
		(list)[key].Status = GetActivityStatus(item.StartTime, item.EndTime, item.Status)
	}
	return total, list, nil
}

// ActivitySet 裂变活动设置
func (s *Service) ActivitySet(c ctx.Context, enterpriseID int64, activity *ActivityEntity, memberIDs []int64, chatGroups []*ChatGroupEntity, welcomeContent []map[string]interface{}, prize []*PrizeEntity, mutualActivityIDs []int64) (int64, error) {
	info := &ActivityEntity{}
	action := ActionAdd
	// 查询企业domain,hash
	domain, _, hash, err := s.Enterprise.GetInfo(c, enterpriseID)
	if err != nil {
		return 0, errs.EnterpriseNotFund
	}
	if domain == "" || hash == "" {
		return 0, errs.EnterpriseNotFund
	}
	if activity.ID != 0 {
		// 查询活动详情
		info, err = s.ActivityDetail(c, enterpriseID, activity.ID)
		if err != nil {
			return 0, err
		}
		action = ActionUpdate
		activity.DrainageWelcomesID = info.DrainageWelcomesID
	}
	wg := sync.WaitGroup{}
	prizes := &PrizeSet{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 判断奖品修改
		prizes = s.fissionPrizeCheck(c, prize, activity.ID, enterpriseID, activity.StartTime)
	}()

	groups := &ChatGroupSet{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 判断社群修改
		groups = s.fissionGroupCheck(c, chatGroups, activity.ID, enterpriseID)
	}()

	mutual := &MutualSet{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 判断互斥修改
		mutual = s.mutualCheck(c, mutualActivityIDs, activity.ID, enterpriseID)
	}()
	wg.Wait()

	// 设置活动
	if action == ActionAdd {
		activity.ID = uid.ID(c)
		// 设置默认联系我二维码
		ID, err := s.GenerateContact(c, enterpriseID, activity.ChannelID, domain, memberIDs, ContactTypeActivityDefault)
		if err != nil {
			return 0, err
		}
		activity.ContactId = ID
	}
	welcomeStr := ""
	switch activity.FissionType {
	case GroupFissionType:
		for index, vMap := range welcomeContent {
			for k, v := range vMap {
				if k == "text" {
					var vText Text
					vj, _ := json.Marshal(v)
					err = json.Unmarshal(vj, &vText)
					if err != nil {
						log.Error(c).Msgf("ActivitySettext err %+v", err)
					} else {
						vText.Content = vText.Content + "活动链接:" + GetFissionActivityUrl(hash, activity.ID)
						welcomeContent[index][k] = vText
					}

				}
			}
		}
		jsons, _ := json.Marshal(welcomeContent)
		welcomeStr = string(jsons)
	case MemberFissionType:
		// 添加员工，拼接link
		var vData = Link{}
		vData.Desc = ActivityStr
		vData.Title = activity.Name
		vData.Url = GetFissionActivityUrl(hash, activity.ID)
		vData.PicUrl = ActivityCover

		var vMap = make(map[string]interface{})
		vMap["link"] = vData
		welcomeContent = append(welcomeContent, vMap)
		jsons, _ := json.Marshal(welcomeContent)
		welcomeStr = string(jsons)
	case VipFissionType:
		contentJson, _ := json.Marshal(welcomeContent)
		welcomeStr = string(contentJson)
	}

	// 判断员工修改
	members := s.fissionMemberCheck(c, memberIDs, activity.ID)
	// 判断是否需重新生成联系我二维码
	if ActionUpdate == action && (0 != len(members.AddMembers) || 0 != len(members.DelMembers)) {
		err := s.UpdateContact(c, enterpriseID, info.ContactId, domain, memberIDs)
		if err != nil {
			return 0, err
		}
	}
	//添加欢迎语
	activityID, err := s.repo.ActivitySet(c,
		enterpriseID,
		action,
		activity,
		members,
		groups,
		welcomeStr,
		prizes,
		mutual)
	if err != nil {
		return 0, err
	}

	return activityID, nil
}

// fissionGroupCheck 判断裂变群
func (s *Service) fissionGroupCheck(c ctx.Context, groups []*ChatGroupEntity, activityID int64, enterpriseID int64) *ChatGroupSet {
	var fissionChatIDs, groupChatIDs []int64
	resp := &ChatGroupSet{}
	var groupMap = make(map[int64]*ChatGroupEntity)
	if activityID != 0 {
		fissionGroup, err := s.repo.FissionGroupIDs(c, activityID, enterpriseID, false)
		if err != nil {
			return resp
		}
		for _, item := range groups {
			groupChatIDs = append(groupChatIDs, item.ChatID)
		}
		if len(fissionGroup) > 0 {
			for _, group := range fissionGroup {
				fissionChatIDs = append(fissionChatIDs, group.ChatID)
				if false == helper.InArray[int64](group.ChatID, groupChatIDs) {
					resp.DelGroups = append(resp.DelGroups, group.ID)
				}
				groupMap[group.ChatID] = group
			}
		}
	}
	for _, group := range groups {
		if false == helper.InArray[int64](group.ChatID, fissionChatIDs) {
			resp.AddGroups = append(resp.AddGroups, group)
		} else {
			if items, ok := groupMap[group.ChatID]; ok {
				group.ID = items.ID
			}
			resp.UpdateGroups = append(resp.UpdateGroups, group)
		}
	}
	return resp
}

// fissionPrizeCheck 判断裂变奖品
func (s *Service) fissionPrizeCheck(c ctx.Context, prizes []*PrizeEntity, activityID int64, enterpriseID int64, startTime int64) *PrizeSet {
	// 查询活动奖品
	var fissionPrizeIDs, prizeIDs []int64
	resp := &PrizeSet{}
	if activityID != 0 {
		total, prizeList, err := s.repo.PrizeList(c, enterpriseID, activityID, 1, 10)
		if err != nil {
			return resp
		}
		for _, item := range prizes {
			prizeIDs = append(prizeIDs, item.ID)
		}
		if total > 0 {
			for _, prize := range *prizeList {
				fissionPrizeIDs = append(fissionPrizeIDs, prize.ID)
				if false == helper.InArray[int64](prize.ID, prizeIDs) {
					resp.DelPrizes = append(resp.DelPrizes, prize.ID)
				}
			}
		}
	}
	for _, prize := range prizes {
		if false == helper.InArray[int64](prize.ID, fissionPrizeIDs) {
			resp.AddPrizes = append(resp.AddPrizes, prize)
		} else {
			// 如果不是未开始的活动
			if startTime < time.Now().Unix() {
				prize.Stocks = 0
			}
			resp.UpdatePrizes = append(resp.UpdatePrizes, prize)
		}
	}
	return resp
}

// fissionMembersCheck 判断裂变员工
func (s *Service) fissionMemberCheck(c ctx.Context, memberIDs []int64, activityID int64) *MembersSet {
	var fissionMemberIDs []int64
	resp := &MembersSet{}
	if activityID != 0 {
		// 查询活动关联员工
		fissionMembers, err := s.repo.FissionMemberIDs(c, activityID)
		if err != nil {
			return resp
		}
		if len(fissionMembers) > 0 {
			for _, member := range fissionMembers {
				fissionMemberIDs = append(fissionMemberIDs, member.MemberID)
				if false == helper.InArray[int64](member.MemberID, memberIDs) {
					resp.DelMembers = append(resp.DelMembers, member.ID)
				}
			}
		}
	}
	for _, memberID := range memberIDs {
		if false == helper.InArray[int64](memberID, fissionMemberIDs) {
			resp.AddMembers = append(resp.AddMembers, memberID)
		}
	}
	return resp
}

// mutualCheck 判断互斥活动
func (s *Service) mutualCheck(c ctx.Context, mutualIDs []int64, activityID int64, enterpriseID int64) *MutualSet {
	var mutualActivityIDs []int64
	resp := &MutualSet{}
	if activityID != 0 {
		mutualActivity, err := s.repo.MutualActivityIDs(c, activityID, enterpriseID)
		if err != nil {
			return resp
		}
		if len(mutualActivity) > 0 {
			for _, mutual := range mutualActivity {
				mutualActivityIDs = append(mutualActivityIDs, mutual.ActivityID)
				if false == helper.InArray[int64](mutual.ActivityID, mutualIDs) {
					resp.DelMutual = append(resp.DelMutual, mutual.ID)
				}
			}
		}
	}
	for _, mutualID := range mutualIDs {
		if false == helper.InArray[int64](mutualID, mutualActivityIDs) {
			resp.AddMutual = append(resp.AddMutual, mutualID)
		}
	}
	return resp
}

// GenerateContact 生成联系我二维码
func (s *Service) GenerateContact(c ctx.Context, enterpriseID int64, channelID int64, domain string, memberIDs []int64, contactType int) (int64, error) {
	ID := uid.ID(ctx.Wrap(c))
	workImpl := new(work.QyWorkImpl).SetDomain(domain)
	// 查询员工的的userID
	resp, err := s.repo.MemberList(c, memberIDs, enterpriseID)
	if err != nil {
		return 0, err
	}
	var userID []string
	for _, item := range resp {
		userID = append(userID, item.UserID)
	}
	res, err := workImpl.GetWxEexternalContact(c, work.ContactCreateRequest{
		ID:           ID,
		CodeCame:     "裂变活动",
		ChannelId:    channelID,
		Type:         TypeFissionDefault,
		AutoFriend:   1,
		UserIds:      userID,
		EnterpriseId: enterpriseID,
	})
	if err != nil {
		log.Error(c).Msgf("GenerateContact-GetWxEexternalContact err %+v", err)
		return 0, errs.GetWxExternalContactErr
	}
	if res == nil {
		e, ok := err.(*errorx.Error)
		// 判断是不是企业微信错误码
		if ok && e.Code == 41054 {
			return 0, err
		}
		return 0, errs.GetWxExternalContactErr
	}
	err = s.Contact.CreateContact(c, enterpriseID, ID, "裂变活动", channelID, res.QrCode, res.ConfigId, contactType)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

// UpdateContact 编辑联系我二维码
func (s *Service) UpdateContact(c ctx.Context, enterpriseID int64, ID int64, domain string, memberIDs []int64) error {
	workImpl := new(work.QyWorkImpl).SetDomain(domain)
	// 获取二维码详情
	_, configID, err := s.Contact.GetContactQrCode(c, ID, enterpriseID)
	if configID == "" {
		return err
	}
	// 查询员工的的userID
	resp, err := s.repo.MemberList(c, memberIDs, enterpriseID)
	if err != nil {
		return err
	}
	var userID []string
	for _, item := range resp {
		userID = append(userID, item.UserID)
	}
	res, err := workImpl.UpdateWxEexternalContact(c, work.ContactCreateRequest{
		ID:       ID,
		CodeCame: "裂变活动",
		UserIds:  userID,
	}, configID)

	if err != nil {
		return errs.GetWxExternalContactErr
	}
	if res == nil {
		e, ok := err.(*errorx.Error)
		// 判断是不是企业微信错误码
		if ok && e.Code == 41054 {
			return err
		}
		return errs.GetWxExternalContactErr
	}
	return nil
}

// MutualActivity 选择互斥活动
func (s *Service) MutualActivityList(c ctx.Context, request *MutualActivityListRequest) (int64, []*MutualActivityEntity, error) {
	// ** request.Action
	//  detail 活动详情，查询此活动关联的互斥活动
	//  list   查询互斥活动列表
	// **
	if request.Action == "list" {
		var filter []int64
		filter = append(filter, request.Filter...)
		mutualActivity, err := s.repo.MutualActivityIDs(c, 0, request.EnterpriseID)
		if err != nil {
			return 0, nil, err
		}
		for _, item := range mutualActivity {
			filter = append(filter, item.MutualActivityID)
		}
		request.Filter = filter
		request.FilterCount = len(filter)
	}
	total, list, err := s.repo.MutualActivityList(c, request)
	if err != nil {
		return 0, nil, err
	}
	for key, item := range list {
		(list)[key].Status = GetActivityStatus(item.StartTime, item.EndTime, item.Status)
	}
	return total, list, nil
}

// GetActivityStatus 获取活动状态
func GetActivityStatus(startTime int64, endTime int64, activityStatus int) int {
	nowTime := time.Now().Unix()
	status := 0
	if startTime > nowTime && activityStatus == StatusOnShelf {
		status = StatusNotStarted
	} else if startTime < nowTime && endTime > nowTime && activityStatus == StatusOnShelf {
		status = StatusOnGoing
	} else if endTime < nowTime {
		status = StatusEnd
	} else if activityStatus == StatusOffShelf {
		status = StatusOffShelf
	}
	return status
}

// ActivityDelete 活动删除
func (s *Service) ActivityDelete(c ctx.Context, enterpriseID int64, activityID int64) error {
	// 查询活动信息
	activityInfo, err := s.ActivityDetail(c, enterpriseID, activityID)
	if err != nil {
		return err
	}
	if activityInfo.ID == 0 {
		return errs.FissionActivityNotExistErr
	}
	if activityInfo.Status == StatusOnGoing {
		return errs.FissionDeleteErr
	}
	// 删除活动
	err = s.repo.ActivityDelete(c, activityID)
	// 删除该活动互斥活动
	err = s.repo.MutualDeleteByActivityID(c, activityID)
	// 删除该活动社群
	err = s.repo.ChatDeleteByActivityID(c, activityID)
	// 删除该活动的奖品
	err = s.repo.PrizeDeleteByActivityID(c, activityID)
	// 删除该活动的欢迎语
	err = s.repo.WelcomeDeleteByID(c, activityInfo.DrainageWelcomesID)
	// 删除该活动的联系我二维码
	err = s.repo.ContactDeleteByID(c, activityInfo.ContactId)
	if err != nil {
		return err
	}
	return nil
}

// ActivityInfo 活动详情
func (s *Service) ActivityInfo(c ctx.Context, activityID int64, enterpriseID int64) (*ActivityInfoResponse, error) {
	resp := &ActivityInfoResponse{}
	activityInfo, _ := s.ActivityDetail(c, enterpriseID, activityID)
	// 查询企业domain,hash
	_, _, hash, err := s.Enterprise.GetInfo(c, enterpriseID)
	if err != nil {
		return resp, err
	}
	if activityInfo.FissionType == GroupFissionType {
		if hash == "" {
			return resp, errs.EnterpriseNotFund
		}
	}
	resp.Activity = *activityInfo
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 查询裂变员工列表
		fissionMember, _ := s.GetFissionMember(c, activityID, enterpriseID)
		resp.Members = fissionMember
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 查询互斥活动
		resp.MutualActivity = make([]*MutualActivityEntity, 0)
		if activityInfo.MutualType == IsMutualType {
			mutualActivity, _ := s.GetMutualActivity(c, activityID, enterpriseID)
			resp.MutualActivity = mutualActivity
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 查询活动奖品
		prizeList, _ := s.PrizeList(c, enterpriseID, activityID, 1, 10)
		for _, item := range prizeList.List {
			// 编辑展示剩余库存
			remain := item.Stocks - item.ExchangeNum
			if remain >= 0 {
				item.Stocks = remain
			}
		}
		resp.Prize = prizeList.List
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 查询活动社群
		if activityInfo.FissionType == GroupFissionType {
			fissionGroup, _ := s.GetFissionGroup(c, activityID, enterpriseID)
			resp.ChatGroups = fissionGroup
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 查询活动欢迎语
		welComeContent, _ := s.repo.WelComeInfo(c, enterpriseID, activityInfo.DrainageWelcomesID)
		contentJson := make([]map[string]interface{}, 0)
		// 替换欢迎语链接
		if activityInfo.FissionType == GroupFissionType {
			activityLink := GetFissionActivityUrl(hash, activityID)
			welComeContent = strings.Replace(welComeContent, "活动链接:"+activityLink, " ", -1)
		}
		_ = json.Unmarshal([]byte(welComeContent), &contentJson)
		resp.WelcomeContent = contentJson
	}()
	wg.Wait()
	return resp, nil
}

// UpdateStatus 更新活动状态
func (s *Service) UpdateStatus(c ctx.Context, enterpriseID, activityID int64, status int) error {
	// 查询活动信息
	activityInfo, err := s.ActivityDetail(c, enterpriseID, activityID)
	if err != nil {
		return err
	}
	if activityInfo.ID == 0 {
		return errs.FissionActivityNotExistErr
	}
	// 判断更新状态与活动状态
	if activityInfo.Status == StatusEnd {
		return errs.FissionUpdateErr
	}
	// 更新活动状态
	err = s.repo.UpdateStatus(c, activityID, status)
	if err != nil {
		return err
	}
	return nil
}

// GetFissionActivityUrl 获取裂变活动
func GetFissionActivityUrl(enterpriseHash string, activityId int64) string {
	domain := config.CfgData.Service.WebH5.Domain
	path := config.CfgData.WebHtmlUrl.FissionActivityUrl
	path = fmt.Sprintf(path, enterpriseHash, activityId)
	return domain + path
}

// ActivityList 裂变活动列表
func (s *Service) QwActivityList(c ctx.Context, request *QwActivityListRequest) (int64, *[]QwActivityList, error) {
	total, list, err := s.repo.QwActivityList(c, request)
	if err != nil {
		return 0, nil, err
	}

	if len(*list) > 0 {
		for key, value := range *list {
			(*list)[key].Status = GetActivityStatus(value.StartTime, value.EndTime, value.Status)
			(*list)[key].HtmlUrl = GetFissionActivityUrl(request.EnterpriseHash, value.Id)
			(*list)[key].ImgUrl = ActivityCover
		}
	}

	return total, list, nil
}

// StatisticsInfo 裂变活动统计数据
func (s *Service) StatisticsInfo(c ctx.Context, QwActivityListRequest *QwActivityListRequest) (*Statistics, error) {
	// 查询活动信息
	activityInfo, err := s.repo.StatisticsInfo(c, QwActivityListRequest)
	if err != nil {
		return activityInfo, err
	}
	return activityInfo, nil
}

// ChatGroupList 获取社群列表
func (s *Service) ChatGroupList(c ctx.Context, request *ChatGroupListRequest) (int64, []*WorkChatGroupEntity, error) {
	// 查询已经关联活动的社群
	var filter []int64
	filter = append(filter, request.Filter...)
	activityGroups, err := s.repo.FissionGroupIDs(c, 0, request.EnterpriseID, false)
	for _, item := range activityGroups {
		filter = append(filter, item.ChatID)
	}
	request.Filter = filter
	request.FilterCount = len(filter)
	total, list, err := s.repo.ChatGroupList(c, request)
	if err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

// GetFissionGroup 查询活动社群列表
func (s *Service) GetFissionGroup(c ctx.Context, activityID int64, enterpriseID int64) ([]WorkChatGroupEntity, error) {
	var list []WorkChatGroupEntity
	resp := make([]*WorkChatGroupEntity, 0)
	fissionGroups, err := s.repo.FissionGroupIDs(c, activityID, enterpriseID, false)
	if err != nil {
		return list, err
	}
	if len(fissionGroups) == 0 {
		return list, nil
	}
	var groupIDs []int64
	var groupMap = make(map[int64]*ChatGroupEntity)
	for _, group := range fissionGroups {
		groupMap[group.ChatID] = group
		groupIDs = append(groupIDs, group.ChatID)
	}
	_, resp, err = s.repo.ChatGroupList(c, &ChatGroupListRequest{
		ChatIDs:      groupIDs,
		EnterpriseID: enterpriseID,
		ChatIDCount:  len(groupIDs),
	})
	if err != nil {
		return list, err
	}
	if len(resp) > 0 {
		for _, item := range resp {
			if group, ok := groupMap[item.GroupID]; ok {
				item.QrCode = group.QrCode
				item.Sort = group.Sort
			}
			list = append(list, *item)
		}
	}
	// 按照sort排序
	sort.Sort(SortChatGroup(list))
	return list, nil
}

type SortChatGroup []WorkChatGroupEntity

func (s SortChatGroup) Len() int           { return len(s) }
func (s SortChatGroup) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortChatGroup) Less(i, j int) bool { return s[i].Sort < s[j].Sort }

// GetMutualActivity 查询活动关联互斥活动列表
func (s *Service) GetMutualActivity(c ctx.Context, activityID int64, enterpriseID int64) ([]*MutualActivityEntity, error) {
	// 查询该活动的互斥活动IDs
	resp := make([]*MutualActivityEntity, 0)
	mutualActivityIDs, err := s.repo.MutualActivityIDs(c, activityID, enterpriseID)
	if err != nil {
		return resp, err
	}
	if len(mutualActivityIDs) == 0 {
		return resp, nil
	}
	var activityIDs []int64
	for _, activity := range mutualActivityIDs {
		activityIDs = append(activityIDs, activity.MutualActivityID)
	}
	_, resp, err = s.MutualActivityList(c, &MutualActivityListRequest{
		MutualActivityIDs:   activityIDs,
		EnterpriseID:        enterpriseID,
		MutualActivityCount: len(activityIDs),
		Action:              "detail",
	})
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetFissionMember 查询活动关联裂变员工列表
func (s *Service) GetFissionMember(c ctx.Context, activityID int64, enterpriseID int64) ([]*Member, error) {
	resp := make([]*Member, 0)
	fissionMemberIDs, err := s.repo.FissionMemberIDs(c, activityID)
	if err != nil {
		return resp, err
	}
	if len(fissionMemberIDs) == 0 {
		return resp, nil
	}
	var memberIDs []int64
	for _, member := range fissionMemberIDs {
		memberIDs = append(memberIDs, member.MemberID)
	}
	resp, err = s.repo.MemberList(c, memberIDs, enterpriseID)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
