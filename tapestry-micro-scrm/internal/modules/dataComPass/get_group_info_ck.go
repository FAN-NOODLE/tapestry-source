package dataComPass

import (
	"encoding/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/wx/wxoffiaccount/cache"
	"micro-scrm/internal/orm"
	"sync"
)

type CompassCacheCk struct {
	Cache *cache.RedisCache
	Key   string
}

func (c CompassCacheCk) DataInfoResultCk(ctx ctx.Context, request DataRequest) (result PanelStatisticsData, err error) {

	// redis缓存
	var resultInfo string
	err = c.Cache.Get(ctx, c.Key, &resultInfo)
	if err == nil {
		var newResultInfo PanelStatisticsData
		if errs := json.Unmarshal([]byte(resultInfo), &newResultInfo); err != nil {
			err = errs
			return
		}
		return newResultInfo, nil
	}

	if request.SearchAll != 1 && len(request.UserIds) == 0 {
		request.SearchAll = 1
	}

	group2, errs := GetBAExternalGroupCk(ctx, request.EnterpriseID, request.EndTime, request.UserIds, request.SearchAll)
	if errs != nil {
		err = errs
		return
	}
	personInfo, errs := GetExternalContactPersonCk(ctx, request.EnterpriseID, request.StartTime, request.EndTime, group2, request.UserIds, request.SearchAll)
	if errs != nil {
		err = errs
		return
	}
	newInfo := GetBaDataCk(ctx, request.EnterpriseID, group2, request.UserIds, request.StartTime, request.EndTime, request.SearchAll)
	personInfo.BuddyData = newInfo

	//添加缓存
	marshal, errs := json.Marshal(personInfo)
	if errs != nil {
		err = errs
		return
	}
	c.Cache.Set(ctx, c.Key, marshal, 1800)

	return personInfo, nil
}

// GetBAExternalGroupCk 获取外部群
func GetBAExternalGroupCk(ctx ctx.Context, enterpriseID, endTime int64, baInfo []string, searchAll int64) (newResult []int64, err error) {

	//所有群信息
	var resultMap2 = make(map[int64]struct{}, 0)
	groupInfo2, err := orm.Service.DataCompass.GetBAGroupInfoCk(ctx, enterpriseID, endTime, baInfo, searchAll)
	if err != nil {
		return nil, err
	}

	for _, info2 := range groupInfo2 {
		resultMap2[info2.GroupChatID] = struct{}{}
	}
	for k2 := range resultMap2 {
		newResult = append(newResult, k2)
	}
	return
}

func GetExternalContactPersonCk(ctx ctx.Context, enterpriseID, startTime, endTime int64, groupChatInfo2 []int64, baInfo []string, searchAll int64) (info PanelStatisticsData, err error) {
	var wg sync.WaitGroup
	//最终返回数据
	var resultInfo PanelStatisticsData
	// 累计群成员
	var GrandTotalNumber int64
	// 新入群人
	var NewAddGroupPersonCount int64
	// 退群人
	var LeaveGroupCount int64
	//群成员加好友数
	var GroupMemberAddFriend int64
	//获取群加好友数
	GroupMemberAddFriend, err = orm.Service.DataCompass.GetGroupCustomCountCk(ctx, enterpriseID, endTime, groupChatInfo2, baInfo, searchAll)
	//获取累计群成员信息
	GrandTotalNumber, err = orm.Service.DataCompass.GetGroupCustomerNumberCountCk(ctx, enterpriseID, endTime, groupChatInfo2)
	//获取ba的memberID
	NewAddGroupPersonCount, err = orm.Service.DataCompass.GetContactEnterPersonCk(ctx, enterpriseID, startTime, endTime, groupChatInfo2)
	if err != nil {
		return resultInfo, err
	}
	LeaveGroupCount, err = orm.Service.DataCompass.GetContactLeavePersonCk(ctx, enterpriseID, startTime, endTime, groupChatInfo2)
	if err != nil {
		return resultInfo, err
	}
	//----------------------------------社群活动--------------------------------------------------
	//活动群的id列表
	var activityInfo = make(map[int64]struct{}, 0)
	var activityGroup = make([]int64, 0)
	var activityGroupID = make([]int64, 0)
	//社群活动的累计群成员数
	var activityCommunityNumber = make(map[string]struct{}, 0)
	//获取列变活动的群
	wg.Add(1)
	var errs2 *error
	go func(errs2 *error) {
		defer wg.Done()
		activityInfoData, err := orm.Service.DataCompass.GetFissionActivityInfoCk(ctx, enterpriseID, endTime, baInfo, searchAll)
		if err != nil {
			errs2 = &err
			return
		}
		for _, v := range activityInfoData {
			activityGroup = append(activityGroup, v.ActivityID)
			activityGroupID = append(activityGroupID, v.ChatID)
			activityInfo[v.ActivityID] = struct{}{}
		}
		newPerson, _ := orm.Service.DataCompass.GetContactPersonCk(ctx, enterpriseID, endTime, activityGroupID)
		for _, model := range newPerson {
			activityCommunityNumber[model.UnionID] = struct{}{}
		}
	}(errs2)

	// 群活动的参与人数
	var activityParticipateNumber int64
	var activityGrandTotal1 = make([]string, 0)
	var activityGrandTotal2 = make(map[string]struct{}, 0)
	//群活动的参与人数
	//客态
	wg.Add(1)
	var errs3 *error
	go func(errs3 *error) {
		defer wg.Done()

		activityPerson, err := orm.Service.DataCompass.GetContactActivityPersonCk(ctx, enterpriseID, startTime, endTime, activityGroup)
		if err != nil {
			errs3 = &err
			return
		}
		for _, unionID := range activityPerson {
			if _, ok := activityCommunityNumber[unionID]; ok {
				activityGrandTotal1 = append(activityGrandTotal1, unionID)
			}
		}
	}(errs3)

	wg.Add(1)
	var errs4 *error
	go func(errs4 *error) {
		defer wg.Done()

		person2, err := orm.Service.DataCompass.GetContactActivityPerson2Ck(ctx, enterpriseID, startTime, endTime, activityGroup)
		if err != nil {
			errs4 = &err
			return
		}
		for _, unionID2 := range person2 {
			if _, ok := activityCommunityNumber[unionID2]; ok {
				activityGrandTotal2[unionID2] = struct{}{}
			}
		}
		for _, v := range activityGrandTotal1 {
			activityGrandTotal2[v] = struct{}{}
		}
		activityParticipateNumber = int64(len(activityGrandTotal2))
	}(errs4)

	//群活动进群人数
	var activityIntoGroupNumber int64
	wg.Add(1)
	var errs5 *error
	go func(errs5 *error) {
		defer wg.Done()
		groupInfo, err := orm.Service.DataCompass.GetIntoActivityGroupInfoCk(ctx, enterpriseID, startTime, endTime, activityGroup)
		if err != nil {
			errs5 = &err
			return
		}
		activityIntoGroupNumber = int64(len(groupInfo))
	}(errs5)

	//助力且为会员的人
	var boostAndMember int64
	wg.Add(1)
	go func() {
		defer wg.Done()
		memberInfo := orm.Service.DataCompass.GetMemberInfoCk(ctx, enterpriseID, activityGrandTotal1)
		boostAndMember = memberInfo
	}()
	wg.Wait()

	if errs5 != nil {
		return resultInfo, nil
	}

	//-------------------------------------数据组装--------------------------------------------
	//社群数据
	var commuitData CommunityData
	commuitData.GroupNumber = IntToString(int64(len(groupChatInfo2)))   //群数量
	commuitData.CumulativeGroupNumber = IntToString(GrandTotalNumber)   //累计群成员
	commuitData.AddFriendNumber = IntToString(GroupMemberAddFriend)     //群加好友数
	commuitData.LeaveGroupNumber = IntToString(LeaveGroupCount)         //退群人数
	commuitData.NumberOfNewcomers = IntToString(NewAddGroupPersonCount) //新入群人数
	var adPercentage, lgPercentage float64
	if float64(GrandTotalNumber) != 0 {
		adPercentage = (float64(GroupMemberAddFriend) / float64(GrandTotalNumber)) * 100
	}
	commuitData.AddFriendPercentage = PercentageRounding(adPercentage) //群成员加好友率
	if float64(int(LeaveGroupCount)+int(GrandTotalNumber)) != 0 {
		lgPercentage = (float64(LeaveGroupCount) / float64(int(LeaveGroupCount)+int(GrandTotalNumber))) * 100
	}
	commuitData.LeaveGroupPercentage = PercentageRounding(lgPercentage) // 退群率

	//活动社群数据
	var communityActivityData CommunityActivityData
	communityActivityData.ActivityGroupNumber = IntToString(int64(len(activityInfo))) //活动群数量
	communityActivityData.IntoGroupNumber = IntToString(activityIntoGroupNumber)      //活动进群人数
	communityActivityData.ParticipateNumber = IntToString(activityParticipateNumber)  //活动参与人数
	var partPercentage, intoPercentage float64
	if float64(len(activityCommunityNumber)) != 0 {
		partPercentage = (float64(activityParticipateNumber) / float64(len(activityCommunityNumber))) * 100
	} else {
		partPercentage = 0
	}
	communityActivityData.ParticipatePercentage = PercentageRounding(partPercentage) // 群活动的参与率
	if float64(activityIntoGroupNumber) != 0 {
		intoPercentage = (float64(boostAndMember) / float64(activityIntoGroupNumber)) * 100
	} else {
		intoPercentage = 0
	}
	communityActivityData.IntoGroupPercentage = PercentageRounding(intoPercentage) //活动进群绑定率
	resultInfo.CommunityData = commuitData
	resultInfo.CommunityActivityData = communityActivityData
	return resultInfo, nil
}
