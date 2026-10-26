package activityCampaign

import (
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/common"
	"sort"
)

//活动的逻辑处理

// GetSignInInfo 签到列表
func GetSignInInfo(enterpriseID, activityID, startTime, endTime int64, pageSize, pageNum, typeInt int) (result []*micro_scrm_enterprise.SignInList, count int64, err error) {
	if pageSize == 0 {
		pageSize = 10
	}
	if pageNum == 0 {
		pageNum = 1
	}
	list, count, err := orm.Service.ActivityCampaign.GetSignInList(enterpriseID, activityID, startTime, endTime, pageSize, pageNum, typeInt)
	if err != nil {
		return nil, 0, err
	}
	if len(list) != 0 {
		sort.SliceIsSorted(list, func(i, j int) bool {
			return list[i].LastSignTime > list[j].LastSignTime
		})
	}
	if typeInt == 1 {
		for _, inList := range list {
			if inList.Name != "" {
				name, _ := common.AesDecrypt(inList.Name)
				inList.Name = name
			}
			if inList.Phone != "" {
				shrPhone, _ := common.AesDecrypt(inList.Phone)
				inList.Phone = shrPhone
			}
		}
		return list, 0, nil
	}
	for _, inList := range list {
		if inList.Name != "" {
			name, _ := common.AesDecrypt(inList.Name)
			inList.Name = hiddenName(name)
		}
		if inList.Phone != "" {
			shrPhone, _ := common.AesDecrypt(inList.Phone)
			inList.Phone = hiddenPhone(shrPhone)
		} else {
			inList.Phone = ""
		}
	}
	return list, count, nil

}
func hiddenPhone(phone string) string {
	str := []rune(phone)
	return string(str[0:3]) + "****" + string(str[7:11])
}
func hiddenName(name string) string {
	nameRune := []rune(name)
	var result string
	lens := len(nameRune)
	if lens <= 1 {
		return name
	}
	if lens == 2 {
		result = string(nameRune[:1]) + "*"
	} else if lens == 3 {
		result = string(nameRune[:1]) + "**"
	} else if lens == 4 {
		result = string(nameRune[:1]) + "***"
	} else if lens > 4 {
		result = string(nameRune[:2]) + "***"
	}
	return result
}

// GetSignInInfoDetail 签到明细
func GetSignInInfoDetail(enterpriseID, activityID, startTime, endTime int64, unionID string, pageSize, pageNum int) (result micro_scrm_enterprise.NewResultSignInDetail, err error) {
	if pageSize == 0 {
		pageSize = 10
	}
	if pageNum == 0 {
		pageNum = 1
	}
	var resultInfo micro_scrm_enterprise.NewResultSignInDetail
	info, err := orm.Service.ActivityCampaign.UserInfo(enterpriseID, activityID, unionID)
	if err != nil {
		return resultInfo, err
	}
	shrPhone, err := common.AesDecrypt(info.ShrPhone)
	if err != nil {
		return resultInfo, err
	}
	shrName, err := common.AesDecrypt(info.ShrName)
	if err != nil {
		return resultInfo, err
	}
	avatarInfo, err := orm.Service.ActivityCampaign.GetAvatarInfo(enterpriseID, unionID)
	if err != nil {
		return resultInfo, err
	}
	resultInfo.Avatar = avatarInfo.Avatar
	detail, err := orm.Service.ActivityCampaign.GetSignDetail(enterpriseID, activityID, startTime, endTime, unionID, pageSize, pageNum)
	if err != nil {
		return resultInfo, err
	}

	resultInfo.Name = info.Name
	if info.Phone != "" {
		resultInfo.Phone = hiddenPhone(info.Phone)
	} else {
		resultInfo.Phone = ""
	}

	resultInfo.Day = int64(len(detail))
	resultInfo.SignInInfoDetail = detail
	resultInfo.Count = int64(len(detail))
	resultInfo.ShrName = hiddenName(shrName)
	if info.ShrPhone != "" {
		resultInfo.ShrPhone = hiddenPhone(shrPhone)
	} else {
		resultInfo.ShrPhone = ""
	}

	resultInfo.ShrAddress = info.ShrAddress
	return resultInfo, nil
}
