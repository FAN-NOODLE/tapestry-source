package dataComPass

import (
	"encoding/json"
	"fmt"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/pkg/db/gormx"
	"strconv"
	"sync"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/wx/wxoffiaccount/cache"

	"micro-scrm/internal/logic/smart_shopping_guide"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/work"
	"micro-scrm/internal/orm/service/group_chat"
	"micro-scrm/internal/pkg/helper"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

type CompassCommunityCache struct {
	Cache *cache.RedisCache
	Key   string
	Count string
}

// GetDataCompassCommunityData 数据缓存
func (d CompassCommunityCache) GetDataCompassCommunityData(ctx2 ctx.Context, params *DataRequest) (result []*CommunityDataDetail, count int64, err error) {
	var cacheString string
	if cacheErr := d.Cache.Get(ctx2, d.Key, &cacheString); cacheErr != nil {
		data, counts := params.DatacompassCommunityData(ctx2)
		marshal, jsonErr := json.Marshal(&data)
		if jsonErr != nil {
			err = jsonErr
			return
		}
		if cacheErr = d.Cache.Set(ctx2, d.Key, string(marshal), 1800); cacheErr != nil {
			err = cacheErr
			return
		}
		if cacheErr = d.Cache.Set(ctx2, d.Count, strconv.FormatInt(counts, 10), 1800); cacheErr != nil {
			err = cacheErr
			return
		}
		return data, counts, nil
	} else {
		if cacheErr := json.Unmarshal([]byte(cacheString), &result); cacheErr != nil {
			err = cacheErr
			return
		}
		var countString string
		if cacheErr := d.Cache.Get(ctx2, d.Count, &countString); cacheErr != nil {
			err = cacheErr
			return
		}
		count, _ = strconv.ParseInt(countString, 10, 64)
		return result, count, nil
	}
}

func (params *DataRequest) DatacompassCommunityData(ctx ctx.Context) ([]*CommunityDataDetail, int64) {

	baList, total, err := smart_shopping_guide.NewSmartShoppingGuideLogic(params.EnterpriseID).GetPagingBaUserIDList(ctx, []uint8{uint8(work.SW_TYPE_ENTERPRISE), uint8(work.SW_TYPE_ENTERPRISE_MANAGE)}, params.ShoppingGuide, 0, 0)

	if err != nil {
		return []*CommunityDataDetail{}, 0
	}
	BAUserIds := []string{}
	for _, value := range baList {
		BAUserIds = append(BAUserIds, value.MemUserid)
	}
	if len(BAUserIds) == 0 {
		return []*CommunityDataDetail{}, 0
	}
	// 根据电话查询随意BA的menberIds
	members, _ := orm.Service.Member.FindBas(params.EnterpriseID, BAUserIds)
	// 根据导购查询对应的群
	requst := group_chat.GropuNameQueryconditions{
		EnterpriseID: params.EnterpriseID,
		StartTime:    params.StartTime,
		EndTime:      params.EndTime,
		Page: group_chat.Page{
			PageSize: params.Page.PageSize,
			PageNum:  params.Page.PageNum,
		},
	}
	groupChats, total, err := orm.Service.GroupChat.FindGroupsBymemberIds(requst, BAUserIds)
	groupChatList := make([]*CommunityDataDetail, 0)
	if err != nil {
		return groupChatList, total
	}
	chatIds := []int64{}
	for _, groupChat := range groupChats {
		NEWgroupChat := &CommunityDataDetail{
			GroupChatId:               groupChat.ID,
			GroupName:                 groupChat.Name,
			CumulativeGroupNumber:     0,
			NumberOfNewcomers:         0,
			LeaveGroupNumber:          0,
			LeaveGroupPercentage:      "0%",
			AddFriendNumber:           0,
			AddFriendPercentage:       "0%",
			WxMemberBindingNumber:     0,
			WxMemberBindingPercentage: "0%",
			ParticipateNumber:         0,
			ParticipatePercentage:     "0%",
			IntoGroupNumber:           0,
			IntoGroupPercentage:       "0%",
			BANumber:                  0,
		}
		groupChatList = append(groupChatList, NEWgroupChat)
		chatIds = append(chatIds, groupChat.ID)
	}
	//------------------
	var sync sync.WaitGroup
	sync.Add(8)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Info(ctx).Msgf("查询群号:%v的群累计人数错误%v", chatIds, err)
			}
		}()
		defer sync.Done()
		cumulativeGroupNumbers, err := orm.Service.GroupChat.FindGroupsMembersByChatIds(chatIds, requst)
		if err == nil {
			for _, groupChat := range groupChatList {
				for _, value := range cumulativeGroupNumbers {
					if groupChat.GroupChatId == value.GroupChatId {
						groupChat.CumulativeGroupNumber = value.Number
						break
					}
				}
			}
		}
	}()
	//群新增人数
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Info(ctx).Msgf("查询群号:%v的群新增人数错误%v", chatIds, err)

			}
		}()
		defer sync.Done()
		numberOfNewcomers, err := orm.Service.GroupChat.FindGroupsAddMembersByChatIds(chatIds, requst)
		if err == nil {
			for _, groupChat := range groupChatList {
				for _, value := range numberOfNewcomers {
					if groupChat.GroupChatId == value.GroupChatId {
						groupChat.NumberOfNewcomers = value.Number
						break
					}
				}
			}
		}

	}()

	// 退群人数
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Info(ctx).Msgf("查询群号:%v的群退群人数错误%v", chatIds, err)

			}
		}()
		defer sync.Done()
		leaveGroupNumbers, err := orm.Service.GroupChat.FindGroupsCancleMembersByChatIds(chatIds, requst)
		if err == nil {
			for _, groupChat := range groupChatList {
				for _, value := range leaveGroupNumbers {
					if groupChat.GroupChatId == value.GroupChatId {
						groupChat.LeaveGroupNumber = value.Number
						break
					}
				}
			}
		}

	}()

	// 群成员加好友数
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Info(ctx).Msgf("查询群号:%v的群成员加好友数错误%v", chatIds, err)

			}
		}()
		defer sync.Done()
		addFriendNumbers, err := orm.Service.GroupChat.FindGroupsMembersAddhaoyouByChatIds(chatIds, requst, members)

		if err == nil {
			for _, groupChat := range groupChatList {
				for _, value := range addFriendNumbers {
					if groupChat.GroupChatId == value.GroupChatId {
						groupChat.AddFriendNumber = value.Number
						break
					}
				}
			}
		}
	}()
	// 微信会员绑定人数
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Info(ctx).Msgf("查询群号:%v的微信会员绑定人数错误%v", chatIds, err)

			}
		}()
		defer sync.Done()
		wxMemberBindingNumbers, err := orm.Service.GroupChat.FindGroupsMembershuiyuanByChatIds(chatIds, requst)
		if err == nil {
			for _, groupChat := range groupChatList {
				for _, value := range wxMemberBindingNumbers {
					if groupChat.GroupChatId == value.GroupChatId {
						groupChat.WxMemberBindingNumber = value.Number
						break
					}
				}
			}
		}
	}()
	// 群活动参与人数
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Info(ctx).Msgf("查询群号:%v的群活动参与人数错误%v", chatIds, err)

			}
		}()
		defer sync.Done()
		fission_users, err := orm.Service.GroupChat.FindGroupszhutaizuliByChatIds(chatIds, requst)

		if err == nil {
			for _, groupChat := range groupChatList {
				for _, value := range fission_users {
					if groupChat.GroupChatId == value.GroupChatId {
						groupChat.ParticipateNumber = groupChat.ParticipateNumber + value.Number
						break
					}
				}
			}
		}

		// participateNumber, err := orm.Service.GroupChat.FindGroupsMembersfissionassistantByChatIds(chatIds, requst)

		// if err == nil {
		// 	for _, groupChat := range groupChatList {
		// 		for _, value := range participateNumber {
		// 			if groupChat.GroupChatId == value.GroupChatId {
		// 				groupChat.ParticipateNumber = groupChat.ParticipateNumber + value.Number
		// 				break
		// 			}
		// 		}
		// 	}
		// }

	}()
	// 群活动进群人数
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Info(ctx).Msgf("查询群号:%v的群活动进群人数错误%v", chatIds, err)

			}
		}()
		defer sync.Done()
		intoGroupNumbers, err := orm.Service.GroupChat.FindGroupsMembersfissionassistantByChatIds(chatIds, requst)
		if err == nil {
			for _, groupChat := range groupChatList {
				for _, value := range intoGroupNumbers {
					if groupChat.GroupChatId == value.GroupChatId {
						groupChat.IntoGroupNumber = value.Number
						break
					}
				}
			}
		}
	}()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Info(ctx).Msgf("查询群号:%v的BA数量错误%v", chatIds, err)

			}
		}()
		defer sync.Done()
		bANumbers, err := orm.Service.GroupChat.FindGroupsBaByChatIds(chatIds, requst, BAUserIds)
		if err == nil {
			for _, groupChat := range groupChatList {
				for _, value := range bANumbers {
					if groupChat.GroupChatId == value.GroupChatId {
						groupChat.BANumber = value.Number
						break
					}
				}
			}
		}
	}()
	sync.Wait()
	for _, groupChat := range groupChatList {
		leaTotal := groupChat.CumulativeGroupNumber + groupChat.LeaveGroupNumber
		if leaTotal != 0 && groupChat.LeaveGroupNumber != 0 {
			groupChat.LeaveGroupPercentage = mathrounded(float64(groupChat.LeaveGroupNumber), float64(leaTotal))
		} else {
			groupChat.LeaveGroupPercentage = "0%"
		}

		if groupChat.CumulativeGroupNumber == 0 {
			groupChat.AddFriendPercentage = "0%"
			groupChat.WxMemberBindingPercentage = "0%"
			groupChat.ParticipatePercentage = "0%"
			groupChat.IntoGroupPercentage = "0%"
		} else {
			groupChat.AddFriendPercentage = mathrounded(float64(groupChat.AddFriendNumber), float64(groupChat.CumulativeGroupNumber))
			groupChat.WxMemberBindingPercentage = mathrounded(float64(groupChat.WxMemberBindingNumber), float64(groupChat.CumulativeGroupNumber))
			groupChat.ParticipatePercentage = mathrounded(float64(groupChat.ParticipateNumber), float64(groupChat.CumulativeGroupNumber))
			groupChat.IntoGroupPercentage = mathrounded(float64(groupChat.IntoGroupNumber), float64(groupChat.CumulativeGroupNumber))
		}

	}
	return groupChatList, total

}

// GetDataCompassFriendData 数据缓存
func (d CompassCommunityCache) GetDataCompassFriendData(ctx2 ctx.Context, params *DataRequest) (result []BuddyDataDetail, count int64, err error) {
	//var countCache = new(int64)
	//countCache = &count
	//result = params.DataCompassFriendData(ctx2, countCache)
	//return

	//获取列表缓存是否存在
	var listString, countString string
	var listCacheSign, countCacheSign bool
	if cacheErr := d.Cache.Get(ctx2, d.Key, &listString); cacheErr == nil {
		if cacheErr = json.Unmarshal([]byte(listString), &result); cacheErr != nil {
			err = cacheErr
			return
		}
		listCacheSign = true
	}
	//查看总条数缓存是否存在
	var countCache = new(int64)
	if cacheErr := d.Cache.Get(ctx2, d.Count, &countString); cacheErr == nil {
		count, _ = strconv.ParseInt(countString, 10, 64)
		countCache = &count
		countCacheSign = true
	}
	if !listCacheSign || !countCacheSign {
		//列表或者总数缓存不存在则直接查询(把总数缓存指针传进去，如果有总数缓存则不再重复查总数)
		result = params.DataCompassFriendData(ctx2, countCache)
		marshal, jsonErr := json.Marshal(&result)
		if jsonErr != nil {
			err = jsonErr
			return
		}
		if cacheErr := d.Cache.Set(ctx2, d.Key, string(marshal), 3600); cacheErr != nil {
			err = cacheErr
			return
		}
		if !countCacheSign {
			//总数缓存不存在则设置总数缓存
			if cacheErr := d.Cache.Set(ctx2, d.Count, strconv.FormatInt(*countCache, 10), 3600); cacheErr != nil {
				err = cacheErr
				return
			}
		}
	}
	return result, *countCache, nil

}

func (params *DataRequest) DataCompassFriendData(ctx ctx.Context, total *int64) (list []BuddyDataDetail) {
	var sTime = strconv.FormatInt(params.StartTime, 10)
	offset := (params.PageNum - 1) * params.PageSize
	gormx := gormx.GetCkInstance().WithContext(ctx).Table("micro_scrm_ck.work_external_user_relation weur").
		Select("member_id shopping_member_id,"+
			"uniqExact(IF(weur.createtime > "+sTime+", weur.unionid,NULL)) AS new_add_buddy_number,"+ //新增好友数
			"uniqExact(weur.unionid) AS grand_total_number,"+ //累计好友数
			//"uniqExact(IF(wgcm.status=1, weur.unionid, NULL)) AS buddy_into_percen_number,"+ //好友入群数
			//"ROUND(uniqExact(IF(wgcm.status=1 AND createtime > "+sTime+", weur.unionid, NULL)) / uniqExact(IF(createtime > "+sTime+", weur.unionid, NULL)) * 100,2) AS new_buddy_into_percentage,"+ //新增好友入群率
			"ROUND(uniqExact(IF(relation_status > 0 ,weur.unionid, NULL)) / uniqExact(weur.unionid) * 100,2) AS buddy_churn_percentage"). //好友流失率
		//"ROUND(uniqExact(IF(wgcm.status=1, weur.unionid, NULL)) / uniqExact(weur.unionid) * 100,2) AS buddy_into_percentage"). //好友入群率
		//Joins("left join micro_scrm_ck.work_group_chat_member wgcm on weur.unionid = wgcm.unionid").
		Where("weur.enterprise_id", params.EnterpriseID).
		Where("weur.deleted_at", 0).
		Where("weur.createtime <= ?", int(params.EndTime)).
		Group("member_id")
	if *total == 0 {
		gormx.Count(total)
	}
	err := gormx.Order("new_add_buddy_number desc,shopping_member_id desc").
		Offset(offset).
		Limit(params.PageSize).
		Find(&list).Error

	var memberIds []int64
	for _, val := range list {
		memberIds = append(memberIds, val.ShoppingMemberId)
	}
	// 查询导购信息
	memberDepartmentList, err := membersDao.NewMembers().GetMemberInfoAndDepartmentByMemberIds(ctx, params.EnterpriseID, memberIds)
	if err != nil {
		return list
	}

	// 整理导购id和导购信息的key=>value
	var memberMap = make(map[int64]MemberInfo)
	for _, info := range memberDepartmentList {
		dpName := ""
		if len(info.MemberDepartments) > 0 && info.MemberDepartments[0].Department.ID > 0 {
			dpName = info.MemberDepartments[0].Department.Name
		}
		memberMap[info.ID] = MemberInfo{
			MemberUserId:   info.Userid,
			MemberName:     info.Name,
			DepartmentName: dpName,
			Avatar:         info.Avatar,
		}
	}
	for i, val := range list {
		if val.NewBuddyIntoPercentage == "NaN" {
			list[i].NewBuddyIntoPercentage = "0%"
		} else {
			list[i].NewBuddyIntoPercentage += "%"
		}
		if val.BuddyChurnPercentage == "NaN" {
			list[i].BuddyChurnPercentage = "0%"
		} else {
			list[i].BuddyChurnPercentage += "%"
		}
		if val.BuddyIntoPercentage == "NaN" {
			list[i].BuddyIntoPercentage = "0%"
		} else {
			list[i].BuddyIntoPercentage += "%"
		}
		if info, ok := memberMap[val.ShoppingMemberId]; ok {
			list[i].BelongingStore = info.DepartmentName
			list[i].ShoppingGuideName = info.MemberName
			list[i].Avatar = info.Avatar
			list[i].ShoppingUserId = info.MemberUserId
			list[i].ShoppingGuideID = info.MemberUserId
		}
	}
	return list

}

func mathrounded(params1, params2 float64) string {
	if params1 == 0 {
		return "0%"
	}
	params3 := helper.MathRounded((params1/params2)*100, 2)
	if params3 == float64(int64(params3)) {
		return fmt.Sprintf("%.0f", params3) + "%"
	}
	return fmt.Sprintf("%.2f", params3) + "%"
}

func EndTimeTOCorde(endTime int64) int64 {
	if endTime != 0 {
		tm := time.Unix(endTime, 0)
		format := tm.Format("2006-01-02")
		format += " 23:59:59"
		location, _ := time.ParseInLocation(timeLayout, format, time.Local)
		return location.Unix()
	} else {
		return 0
	}
}

// GetRankData 导购顾客排行
func (d CompassCommunityCache) GetRankData(ctx2 ctx.Context, params *DataRequest) (result []RankDataCk, count int64, err error) {

	var cacheString string
	if cacheErr := d.Cache.Get(ctx2, d.Key, &cacheString); cacheErr != nil {
		data, counts := params.RankData(ctx2)
		marshal, jsonErr := json.Marshal(data)
		if jsonErr != nil {
			err = jsonErr
			return
		}
		if cacheErr = d.Cache.Set(ctx2, d.Key, string(marshal), 1800); cacheErr != nil {
			err = cacheErr
			return
		}
		if cacheErr = d.Cache.Set(ctx2, d.Count, strconv.FormatInt(counts, 10), 1800); cacheErr != nil {
			err = cacheErr
			return
		}
		return data, counts, nil
	} else {

		if cacheErr := json.Unmarshal([]byte(cacheString), &result); cacheErr != nil {
			err = cacheErr
			return
		}
		var countString string
		if cacheErr := d.Cache.Get(ctx2, d.Count, &countString); cacheErr != nil {
			err = cacheErr
			return
		}
		count, _ = strconv.ParseInt(countString, 10, 64)
		return result, count, nil
	}
}

func (params *DataRequest) RankData(ctx ctx.Context) (baList []RankDataCk, count int64) {

	var rankLength int
	if params.Type == 2 {
		rankLength = 100
	} else {
		rankLength = 50
	}

	var list []RankDataCk
	err := gormx.GetCkInstance().WithContext(ctx).Table("micro_scrm_ck.work_external_user_relation").
		Select("member_id,count() total").
		Where("enterprise_id", params.EnterpriseID).
		Where("deleted_at", 0).
		//Where("relation_status", 0).
		Where("createtime >= ?", int(params.StartTime)).
		Where("createtime <= ?", int(params.EndTime)).
		Group("member_id").
		Order("total desc").
		Limit(rankLength).Find(&list).Error

	var memberIds []int64
	for _, val := range list {
		memberIds = append(memberIds, val.MemberId)
	}
	// 查询导购信息
	memberDepartmentList, err := membersDao.NewMemberDepartments().GetMembersAndDepartmentsByMemberIds(ctx, params.EnterpriseID, memberIds)
	if err != nil {
		return baList, count
	}

	// 整理导购id和导购信息的key=>value
	var memberMap = make(map[int64]MemberInfo)
	for _, info := range memberDepartmentList {
		memberMap[info.MemberID] = MemberInfo{
			MemberUserId:   info.Member.Userid,
			MemberName:     info.Member.Name,
			DepartmentName: info.Department.Name,
		}
	}
	for i, val := range list {
		if info, ok := memberMap[val.MemberId]; ok {
			list[i].DpName = info.DepartmentName
			list[i].Name = info.MemberName
			list[i].MemberId = val.MemberId
		}
	}

	return list, int64(len(list))

}
