package dataComPass

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/wx/wxoffiaccount/cache"
	"math"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"strings"
	"sync"
)

const (
	timeLayout = "2006-01-02 15:04:05"
)

// GetBAExternalGroup 获取ba的外部群
func GetBAExternalGroup(ctx ctx.Context, enterpriseID, startTime, endTime int64, baInfo []string, searchAll int64) (newResult []int64, err error) {
	//所有群信息
	//var resultMap1 = make(map[int64]struct{}, 0)
	var resultMap2 = make(map[int64]struct{}, 0)
	//memberIDList, _ := orm.Service.DataCompass.GetMemberIDList(enterpriseID, baInfo)
	//groupInfo, err := orm.Service.DataCompass.GetBAGroupInfo(enterpriseID, startTime, endTime, baInfo)
	groupInfo2, err := orm.Service.DataCompass.GetBAGroupInfo2(ctx, enterpriseID, endTime, baInfo, searchAll)
	if err != nil {
		return nil, err
	}
	//for _, info := range groupInfo {
	//	resultMap1[info.GroupChatID] = struct{}{}
	//}
	for _, info2 := range groupInfo2 {
		resultMap2[info2.GroupChatID] = struct{}{}
	}
	for k2 := range resultMap2 {
		newResult = append(newResult, k2)
	}
	return
}

func GetExternalContactPerson(ctx ctx.Context, enterpriseID, startTime, endTime int64, groupChatInfo2 []int64, baInfo []string, searchAll int64) (info PanelStatisticsData, err error) {
	//获取当日零时,末时 时间戳
	//zeroHour, endHour := ws_time.DayStartAndEndTimestamp()
	var wg sync.WaitGroup
	//最终返回数据
	var resultInfo PanelStatisticsData
	// 累计群成员
	var GrandTotalNumber int64
	// 新入群人
	//var NewAddGroupPerson = make([]string, 0)
	var NewAddGroupPersonCount int64
	// 退群人
	//var LeaveGroup = make([]string, 0)
	var LeaveGroupCount int64
	//var unionList = make([]string, 0)
	//群成员加好友数
	var GroupMemberAddFriend int64
	//获取群加好友数
	GroupMemberAddFriend, err = orm.Service.DataCompass.GetGroupCustomCount(ctx, enterpriseID, endTime, groupChatInfo2, baInfo, searchAll)
	//获取累计群成员信息
	GrandTotalNumber, err = orm.Service.DataCompass.GetGroupCustomerNumberCount(ctx, enterpriseID, endTime, groupChatInfo2)
	//获取ba的memberID
	NewAddGroupPersonCount, err = orm.Service.DataCompass.GetContactEnterPerson(ctx, enterpriseID, startTime, endTime, groupChatInfo2)
	if err != nil {
		return resultInfo, err
	}
	//NewAddGroupPerson = enterPerson
	LeaveGroupCount, err = orm.Service.DataCompass.GetContactLeavePerson(ctx, enterpriseID, startTime, endTime, groupChatInfo2)
	if err != nil {
		return resultInfo, err
	}
	//LeaveGroup = leavePerson
	////获取群成员加好友的人信息
	//wg.Add(1)
	//var errs1 *error
	//go func(errs1 *error) {
	//	defer wg.Done()
	//	relationInfo, err := orm.Service.DataCompass.GetRelationInfo(enterpriseID, unionList, baInfo)
	//	if err != nil {
	//		errs1 = &err
	//		return
	//	}
	//	for _, newIdInfo := range relationInfo {
	//		if _, ok := GrandTotalNumber[newIdInfo]; ok {
	//			GroupMemberAddFriend[newIdInfo] = struct{}{}
	//		}
	//	}
	//
	//}(errs1)

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

		activityInfoData, err := orm.Service.DataCompass.GetFissionActivityInfo(enterpriseID, endTime, baInfo)
		if err != nil {
			errs2 = &err
			return
		}
		for _, v := range activityInfoData {
			activityGroup = append(activityGroup, v.ActivityID)
			activityGroupID = append(activityGroupID, v.ChatID)
			activityInfo[v.ActivityID] = struct{}{}
		}
		newPerson, _ := orm.Service.DataCompass.GetContactPerson(ctx, enterpriseID, endTime, activityGroupID)
		for _, model := range newPerson {
			activityCommunityNumber[model.UnionID] = struct{}{}
		}
	}(errs2)
	//fmt.Println("111", activityInfoData)

	//
	//activityInfo, err := orm.Service.DataCompass.GetFissionActivityInfo(enterpriseID, endTime, baInfo)
	//if err != nil {
	//	return resultInfo, err
	//}
	//for _, v := range activityInfo {
	//	activityGroup = append(activityGroup, v.ActivityID)
	//	activityGroupID = append(activityGroupID, v.ChatID)
	//}
	//log.Println("%v",activityGroupID)
	////社群活动的累计群成员数
	//var activityCommunityNumber = make(map[string]struct{}, 0)
	//newPerson, _ := orm.Service.DataCompass.GetContactPerson(enterpriseID, endTime, activityGroupID)
	//
	//for _, model := range newPerson {
	//	activityCommunityNumber[model.UnionID] = struct{}{}
	//}

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

		activityPerson, err := orm.Service.DataCompass.GetContactActivityPerson(enterpriseID, startTime, endTime, activityGroup)
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

	//activityPerson, err := orm.Service.DataCompass.GetContactActivityPerson(enterpriseID, startTime, endTime, activityGroup)
	//if err != nil {
	//	return resultInfo, err
	//}
	//for _, unionID := range activityPerson {
	//	if _, ok := activityCommunityNumber[unionID]; ok {
	//		activityGrandTotal1 = append(activityGrandTotal1, unionID)
	//	}
	//}
	//主态
	wg.Add(1)
	var errs4 *error
	go func(errs4 *error) {
		defer wg.Done()

		person2, err := orm.Service.DataCompass.GetContactActivityPerson2(enterpriseID, startTime, endTime, activityGroup)
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

	//person2, err := orm.Service.DataCompass.GetContactActivityPerson2(enterpriseID, startTime, endTime, activityGroup)
	//if err != nil {
	//	return resultInfo, err
	//}
	//for _, unionID2 := range person2 {
	//	if _, ok := activityCommunityNumber[unionID2]; ok {
	//		activityGrandTotal2[unionID2] = struct{}{}
	//	}
	//}
	//for _, v := range activityGrandTotal1 {
	//	activityGrandTotal2[v] = struct{}{}
	//}
	//activityParticipateNumber = int64(len(activityGrandTotal2))

	//群活动进群人数
	var activityIntoGroupNumber int64
	wg.Add(1)
	var errs5 *error
	go func(errs5 *error) {
		defer wg.Done()
		groupInfo, err := orm.Service.DataCompass.GetIntoActivityGroupInfo(enterpriseID, startTime, endTime, activityGroup)
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
		memberInfo := orm.Service.DataCompass.GetMemberInfo(enterpriseID, activityGrandTotal1)
		boostAndMember = memberInfo
	}()
	wg.Wait()

	if errs5 != nil {
		return resultInfo, nil
	}
	//-------------------------------------数据组装--------------------------------------------
	//社群数据
	var commuitData CommunityData
	commuitData.GroupNumber = IntToString(int64(len(groupChatInfo2))) //群数量
	commuitData.CumulativeGroupNumber = IntToString(GrandTotalNumber) //累计群成员
	commuitData.AddFriendNumber = IntToString(GroupMemberAddFriend)   //群加好友数
	//commuitData.LeaveGroupNumber = IntToString(int64(len(LeaveGroup)))  //退群人数
	commuitData.LeaveGroupNumber = IntToString(LeaveGroupCount) //退群人数
	//commuitData.NumberOfNewcomers = IntToString(int64(len(NewAddGroupPerson))) //新入群人数
	commuitData.NumberOfNewcomers = IntToString(NewAddGroupPersonCount) //新入群人数
	var adPercentage, lgPercentage float64
	if float64(GrandTotalNumber) != 0 {
		adPercentage = (float64(GroupMemberAddFriend) / float64(GrandTotalNumber)) * 100
	}
	commuitData.AddFriendPercentage = PercentageRounding(adPercentage) //群成员加好友率
	//if float64(len(LeaveGroup)+int(GrandTotalNumber)) != 0 {
	//	lgPercentage = (float64(len(LeaveGroup)) / float64(len(LeaveGroup)+int(GrandTotalNumber))) * 100
	//}
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

// IntToString k值转换
func IntToString(num int64) string {
	var value string
	if num > 10000 {
		s := float64(num) / float64(1000)
		s0 := fmt.Sprintf("%.1f", s)
		s1 := fmt.Sprintf("%.3f", s)
		splitList := strings.Split(s1, ".")
		lastNumber := splitList[len(splitList)-1]
		var newString string
		for i, s2 := range []rune(lastNumber) {
			if i == 1 {
				newString += "."
			}
			newString += string(s2)
		}
		float, err := strconv.ParseFloat(newString, 64)
		if err != nil {
			return ""
		}
		floor := math.Floor(float + 0.5)
		sprintfs := fmt.Sprintf("%.1f", floor)
		s00 := strings.Split(s0, ".")[0]
		s01 := strings.Split(sprintfs, ".")[0]
		value = s00 + "." + s01 + "k"
		//value = fmt.Sprintf("%.1fk",math.Floor(s) + 0.5)
	} else {
		value = fmt.Sprintf("%d", num)
	}
	return value
}

// PercentageRounding 百分比四舍五入
func PercentageRounding(infos float64) string {
	if infos == 0 {
		return "0%"
	} else {
		s := fmt.Sprintf("%.2f", helper.MathRounded(infos, 2))
		return s + "%"
	}

}

type CompassCache struct {
	Cache *cache.RedisCache
	Key   string
}

func (c CompassCache) DataInfoResult(ctx ctx.Context, request DataRequest) (result PanelStatisticsData, err error) {
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

	//if request.SearchAll != 1 {
	//	part, errs := smart_shopping_guide.NewSmartShoppingGuideLogic(request.EnterpriseID).GetBaUserIDListPart(ctx, []uint8{work.SW_TYPE_ENTERPRISE}, request.ShoppingGuide)
	//	if errs != nil {
	//		err = errs
	//		return
	//	}
	//	request.ShoppingGuide = part
	//}
	group2, errs := GetBAExternalGroup(ctx, request.EnterpriseID, request.StartTime, request.EndTime, request.UserIds, request.SearchAll)
	if errs != nil {
		err = errs
		return
	}
	personInfo, errs := GetExternalContactPerson(ctx, request.EnterpriseID, request.StartTime, request.EndTime, group2, request.UserIds, request.SearchAll)
	if errs != nil {
		err = errs
		return
	}
	newInfo := GetBaData(ctx, request.EnterpriseID, group2, request.UserIds, request.StartTime, request.EndTime, request.SearchAll)
	personInfo.BuddyData = newInfo
	marshal, errs := json.Marshal(personInfo)
	if errs != nil {
		err = errs
		return
	}
	c.Cache.Set(ctx, c.Key, marshal, 86400)
	return personInfo, nil
}
