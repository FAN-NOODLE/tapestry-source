package dataComPass

import (
	"fmt"
	"strconv"
	"sync"

	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

func GetNewFriendsCount(enterpriseId int64, memberId []string, startTime int64, endTime, searchAll int64) (data int64, err error) {
	data, err = orm.Service.DataCompass.GetNewGoodFriendsCount(enterpriseId, memberId, startTime, endTime, searchAll)
	return
}

func GetFriendsAll(enterpriseId int64, memberId []string, endTime, searchAll int64) (total int64, err error) {
	total, err = orm.Service.DataCompass.GetGoodFriendsAll(enterpriseId, memberId, endTime, searchAll)
	return
}

func GetFriendsVIP(enterpriseId int64, startTime int64, endTime int64, memberIds []string) (data int64, err error) {
	data, err = orm.Service.DataCompass.GetFriendsVIP(enterpriseId, startTime, endTime, memberIds)
	return
}

func GetFriendsVIPAll(enterpriseId int64, endTime int64, memberIds []string, searchAll int64) (data int64, err error) {
	data, err = orm.Service.DataCompass.GetFriendsVIPAll(enterpriseId, endTime, memberIds, searchAll)
	return
}

func GetUnionId(enterpriseId int64, Ids []int64) (data []micro_scrm_enterprise.WxUserCorpExternalModel, err error) {
	data, err = orm.Service.DataCompass.GeUnionIds(enterpriseId, Ids)
	return
}

func GetJoinGroupCount(enterpriseId int64, baGroupIds []int64, memberIds []string, startTime int64, endTime, searchAll int64) (data int64, err error) {
	data, err = orm.Service.DataCompass.GetJoinCroupCount(enterpriseId, baGroupIds, memberIds, startTime, endTime, searchAll)
	return
}

func GetJoinGroupCountAll(enterpriseId int64, baGroupIds []int64, memberIds []string, endTime, searchAll int64) (data int64, err error) {
	data, err = orm.Service.DataCompass.GetJoinCroupCountAll(enterpriseId, baGroupIds, memberIds, endTime, searchAll)
	return
}

func GetDrainFriends(enterpriseId int64, memberId []string, startTime int64, endTime, searchAll int64) (data int64, err error) {
	data, err = orm.Service.DataCompass.GetDrainFriends(enterpriseId, memberId, startTime, endTime, searchAll)
	return
}

func GetMemberId(enterpriseId int64, memberUserId []string) (data []micro_scrm_enterprise.MembersModel, err error) {
	data, err = orm.Service.DataCompass.GetMemberId(enterpriseId, memberUserId)
	return
}

func GetBaData(c ctx.Context, enterpriseId int64, baGroupIds []int64, memberUserId []string, startTime int64, endTime, searchAll int64) BuddyData {

	memberData, MemberErr := GetMemberId(enterpriseId, memberUserId)
	if MemberErr != nil {
		log.Info(ctx.Wrap(c)).Msgf("好友数据：新增好友数据报错：%v", MemberErr.Error())
		return BuddyData{}
	}
	var memberIds []string
	for _, item := range memberData {
		memberIds = append(memberIds, strconv.Itoa(int(item.ID)))
	}
	wg := sync.WaitGroup{}
	var newFriendDataCount int64
	var NewJoinGroupCount int64
	wg.Add(1)
	go func() {
		defer wg.Done()

		//获取新增好友数据
		newFriendDataRes, newFriendDataErr := GetNewFriendsCount(enterpriseId, memberIds, startTime, endTime, searchAll)
		if newFriendDataErr != nil {
			log.Info(ctx.Wrap(c)).Msgf("好友数据：新增好友数据报错：%v", newFriendDataErr.Error())
			return
		}
		newFriendDataCount = newFriendDataRes
		//获取新增好友入群人数
		NewJoinGroupCountRes, JoinGroupCountErr := GetJoinGroupCount(enterpriseId, baGroupIds, memberIds, startTime, endTime, searchAll)
		if JoinGroupCountErr != nil {
			log.Info(ctx.Wrap(c)).Msgf("好友数据：通过unionid获取入群人数报错：%v", JoinGroupCountErr.Error())
			return
		}
		NewJoinGroupCount = NewJoinGroupCountRes

	}()

	wg.Add(1)
	var AllFriendsData int64
	//var allFriendsVip int64
	var AllJoinGroupCount int64
	go func() {
		defer wg.Done()
		//获取累计好友数据
		AllFriendsTotal, FriendsAllDataErr := GetFriendsAll(enterpriseId, memberIds, endTime, searchAll)
		if FriendsAllDataErr != nil {
			log.Info(ctx.Wrap(c)).Msgf("好友数据：获取累计好友数据报错：%v", FriendsAllDataErr.Error())
			return
		}
		AllFriendsData = AllFriendsTotal

		//累计会员绑定人数
		//allFriendsVipRes, allFriendsVipeErr := GetFriendsVIPAll(enterpriseId, endTime, memberIds, searchAll)
		//if allFriendsVipeErr != nil {
		//	log.Info(ctx.Wrap(c)).Msgf("好友数据：累计会员绑定人数报错：%v", FriendsAllDataErr.Error())
		//	return
		//}
		//allFriendsVip = allFriendsVipRes

		//通过unionid获取入群人数
		AllJoinGroupCountRes, AllJoinGroupCountErr := GetJoinGroupCountAll(enterpriseId, baGroupIds, memberIds, endTime, searchAll)
		if AllJoinGroupCountErr != nil {
			log.Info(ctx.Wrap(c)).Msgf("好友数据：通过unionid获取入群人数报错：%v", AllJoinGroupCountErr.Error())
			return
		}
		AllJoinGroupCount = AllJoinGroupCountRes
	}()

	wg.Add(1)
	var FriendsDrain int64
	go func() {
		defer wg.Done()
		//好友流失人数
		FriendsDrainRes, FriendsDrainErr := GetDrainFriends(enterpriseId, memberIds, startTime, endTime, searchAll)
		if FriendsDrainErr != nil {
			log.Info(ctx.Wrap(c)).Msgf("好友数据：好友流失人数报错：%v", FriendsDrainErr.Error())
			return
		}
		FriendsDrain = FriendsDrainRes
	}()
	wg.Wait()
	intoNum := ""
	if newFriendDataCount == 0 {
		intoNum = "0%"
	} else {
		intoNum = mathrounded(float64(NewJoinGroupCount), float64(newFriendDataCount))
	}

	joinGroup := ""
	if AllFriendsData == 0 {
		joinGroup = "0%"
	} else {
		joinGroup = mathrounded(float64(AllJoinGroupCount), float64(AllFriendsData))
	}

	diain := ""
	if (float64(AllFriendsData) + float64(FriendsDrain)) == 0 {
		diain = "0%"
	} else {
		diain = mathrounded(float64(FriendsDrain), float64(AllFriendsData)+float64(FriendsDrain))
	}
	return BuddyData{
		BANumber:             fmt.Sprintf("%d", len(memberIds)),
		NewAddBuddyNumber:    IntToString(newFriendDataCount),
		GrandTotalNumber:     IntToString(AllFriendsData),
		NewBuddyIntoNumber:   intoNum,
		BuddyChurnNumber:     IntToString(FriendsDrain),
		BuddyChurnPercentage: diain,
		BuddyIntoNumber:      IntToString(AllJoinGroupCount),
		BuddyIntoPercentage:  joinGroup,
	}
}

//func IntToString(num int64) string {
//	var value string
//	if num > 10000 {
//		value = fmt.Sprintf("%.1fk", float64(num)/float64(1000))
//	} else {
//		value = fmt.Sprintf("%d", num)
//	}
//	return value
//}

// GetBaFirstDepartment 获取BA的一级门店
func GetBaFirstDepartment(enterpriseID int64, userID string) (info micro_scrm.MemberDepartmentsModel, err error) {
	info, err = orm.Service.Departments.GetUserDepartmentInfo(enterpriseID, userID)
	return
}
