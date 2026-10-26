package dataComPass

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"strconv"
	"sync"
)

func GetMemberIdCk(c ctx.Context, enterpriseId int64, memberUserId []string) (data []micro_scrm_enterprise.MembersModel, err error) {
	data, err = orm.Service.DataCompass.GetMemberIdCk(c, enterpriseId, memberUserId)
	return
}

func GetNewFriendsCountCk(c ctx.Context, enterpriseId int64, memberId []string, startTime int64, endTime int64, searchAll int64) (data int64, err error) {
	data, err = orm.Service.DataCompass.GetNewGoodFriendsCountCk(c, enterpriseId, memberId, startTime, endTime, searchAll)
	return
}

func GetJoinGroupCountCk(c ctx.Context, enterpriseId int64, baGroupIds []int64, memberIds []string, startTime int64, endTime, searchAll int64) (data int64, err error) {
	data, err = orm.Service.DataCompass.GetJoinCroupCountCk(c, enterpriseId, baGroupIds, memberIds, startTime, endTime, searchAll)
	return
}

func GetFriendsAllCk(c ctx.Context, enterpriseId int64, memberId []string, endTime, searchAll int64) (total int64, err error) {
	total, err = orm.Service.DataCompass.GetGoodFriendsAllCk(c, enterpriseId, memberId, endTime, searchAll)
	return
}

//GetJoinGroupCountAllCk 通过unionid获取入群人数
func GetJoinGroupCountAllCk(c ctx.Context, enterpriseId int64, baGroupIds []int64, memberIds []string, endTime, searchAll int64) (data int64, err error) {
	data, err = orm.Service.DataCompass.GetJoinCroupCountAllCk(c, enterpriseId, baGroupIds, memberIds, endTime, searchAll)
	return
}

//GetDrainFriendsCk 好友流失人数
func GetDrainFriendsCk(c ctx.Context, enterpriseId int64, memberId []string, startTime int64, endTime, searchAll int64) (data int64, err error) {
	data, err = orm.Service.DataCompass.GetDrainFriendsCk(c, enterpriseId, memberId, startTime, endTime, searchAll)
	return
}

func GetBaDataCk(c ctx.Context, enterpriseId int64, baGroupIds []int64, memberUserId []string, startTime int64, endTime, searchAll int64) BuddyData {

	memberData, MemberErr := GetMemberIdCk(c, enterpriseId, memberUserId)
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
		newFriendDataRes, newFriendDataErr := GetNewFriendsCountCk(c, enterpriseId, memberIds, startTime, endTime, searchAll)
		if newFriendDataErr != nil {
			log.Info(ctx.Wrap(c)).Msgf("好友数据：新增好友数据报错：%v", newFriendDataErr.Error())
			return
		}
		newFriendDataCount = newFriendDataRes
		//获取新增好友入群人数
		NewJoinGroupCountRes, JoinGroupCountErr := GetJoinGroupCountCk(c, enterpriseId, baGroupIds, memberIds, startTime, endTime, searchAll)
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
		AllFriendsTotal, FriendsAllDataErr := GetFriendsAllCk(c, enterpriseId, memberIds, endTime, searchAll)
		if FriendsAllDataErr != nil {
			log.Info(ctx.Wrap(c)).Msgf("好友数据：获取累计好友数据报错：%v", FriendsAllDataErr.Error())
			return
		}
		AllFriendsData = AllFriendsTotal

		//通过unionid获取入群人数
		AllJoinGroupCountRes, AllJoinGroupCountErr := GetJoinGroupCountAllCk(c, enterpriseId, baGroupIds, memberIds, endTime, searchAll)
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
		FriendsDrainRes, FriendsDrainErr := GetDrainFriendsCk(c, enterpriseId, memberIds, startTime, endTime, searchAll)
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
