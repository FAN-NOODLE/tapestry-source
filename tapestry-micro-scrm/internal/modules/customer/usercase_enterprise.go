package customer

import (
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gorm.io/gorm"
	workExternalDao "micro-scrm/internal/app/dao/micro_scrm"
	workExternalModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/helper"
	"strings"
	"time"
)

// QueryList 管理后台 客户列表
func (s *Service) QueryList(c ctx.Context, cond *ListRequest) (*ListResponse, error) {
	cond.EnterpriseID = common.EnterpriseID(c)

	var sortKeyListMap = make(map[string]string)
	sortKeyListMap["pay_order_num"] = "pay_order"
	sortKeyListMap["pay_total"] = "total_fee"
	sortKeyListMap["refund_order_num"] = "refund_order"
	sortKeyListMap["refund_total"] = "refund_fee"
	sortKeyListMap["per_customer_transaction"] = "per_customer_transaction"
	sortKeyListMap["latest_marketing_time"] = "latest_market_time"

	sortKey, found := sortKeyListMap[cond.SortKey]
	if !found {
		cond.SortKey = ""
	} else {
		cond.SortKey = sortKey
		cond.SortType = strings.ToUpper(cond.SortType)
		if cond.SortType != "ASC" && cond.SortType != "DESC" {
			cond.SortType = "DESC"
		}
	}

	total, userList, err := s.esRepo.QueryUserList(c, cond)

	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	timeCond := cond
	timeCond.Size = cond.Size
	timeCond.SortKey = "id"
	timeCond.SortType = "DESC"
	_, timeList, err := s.esRepo.QueryUserList(c, timeCond)
	list := &ListResponse{
		List:     []ListItem{},
		Total:    total,
		Page:     cond.Page,
		Size:     cond.Size,
		DtUpdate: time.Now().Format("2006-01-02 15:04:05"),
	}
	if len(timeList) > 0 {
		for _, lu := range timeList {
			list.DtUpdate = lu.DtUpdate
		}
	}
	var workExternalUserIDs []int64
	//tagListMap := make(map[int64][]int64, 0)
	var unionIds []string
	for _, eu := range userList {
		unionIds = append(unionIds, eu.UnionId)
		for _, ext := range eu.WxExternal {
			workExternalUserIDs = append(workExternalUserIDs, ext.ID)
		}
	}

	userSaInfos, errs := workExternalDao.NewUserRelation().GetUserSAInfoByUnionIds(c, cond.EnterpriseID, unionIds)
	if errs != nil && !errors.Is(errs, gorm.ErrRecordNotFound) {
		return nil, errs
	}
	userSaListMap := make(map[string]workExternalModel.UserSaInfo, 0)

	for _, v := range userSaInfos {
		userSaListMap[v.UnionId] = v
	}

	relationsListMap := make(map[int64][]*WorkExternalUserRemark, 0)
	if len(workExternalUserIDs) > 0 {
		relations, err1 := s.workRepo.QueryExternalRelationRemarks(c, cond.EnterpriseID, workExternalUserIDs)
		if err1 != nil {
			log.Err(c, err1).Msgf("获取外部联系人关系失败: %+v", workExternalUserIDs)
		}
		for _, v := range relations {
			relationsListMap[v.WorkExternalUserID] = append(relationsListMap[v.WorkExternalUserID], v)
		}
	}

	for _, eu := range userList {
		item := ListItem{UserID: eu.ID, Name: eu.Name, Avatar: eu.Avatar, Mobile: eu.Mobile, RemarkList: []string{}, Relations: []*WorkExternalUserRemark{}}
		//sr := ListItem{Relations: []*WorkExternalUserRemark{}}
		if item.Name == "" {
			item.Name = eu.Name
		}
		if item.Nick == "" {
			item.Nick = eu.Nick
		}
		if item.Grade == "" {
			item.Grade = eu.Grade
		}
		if item.UnionId == "" {
			item.UnionId = eu.UnionId
		}
		if item.Point == "" {
			item.Point = eu.Point
		}
		if item.Birthday == "" {
			item.Birthday = eu.Birthday
		}
		if item.UnionCode == "" {
			item.UnionCode = eu.UnionCode
		}
		if eu.Remark != "" {
			item.RemarkList = helper.DelItems[string](strings.Split(eu.Remark, ","), []string{""})
		}
		sr := ListItem{Relations: []*WorkExternalUserRemark{}}
		for _, ext := range eu.WxExternal {
			if _, ok := relationsListMap[ext.ID]; ok {
				for _, v := range relationsListMap[ext.ID] {
					v.RemarkList = []string{}
					if v.Remark == "" {
						continue
					}
					v.RemarkList = helper.DelItems[string](strings.Split(v.Remark, ","), []string{""})
					sr.Relations = append(sr.Relations, v)
				}
				item.Relations = relationsListMap[ext.ID]
			}
		}

		//var workExternalUserIDs []int64
		//for _, ext := range eu.WxExternal {
		//	workExternalUserIDs = append(workExternalUserIDs, ext.ID)
		//}
		//if len(workExternalUserIDs) > 0 {
		//	item.Relations, err = s.workRepo.QueryExternalRelationRemarks(c, cond.EnterpriseID, workExternalUserIDs)
		//	if err != nil {
		//		log.Err(c, err).Msgf("获取外部联系人关系失败: %+v", workExternalUserIDs)
		//		continue
		//	}
		//
		//	for _, v := range item.Relations {
		//		v.RemarkList = []string{}
		//		if v.Remark == "" {
		//			continue
		//		}
		//		v.RemarkList = assist.DelStringFromSlice(strings.Split(v.Remark, ","), "")
		//		sr.Relations = append(sr.Relations, v)
		//	}
		//
		//	item.Relations = sr.Relations
		//}
		if len(eu.Tags) > 0 {
			tags, _ := s.tagRepo.QueryTags(c, eu.Tags)
			for _, tag := range tags {
				item.Tags = append(item.Tags, &TagItem{ID: tag.ID, Name: tag.TagName})
			}
		}
		if len(eu.Groups) > 0 {
			groups, _ := s.tagRepo.QueryUserTagGroups(c, eu.Groups)
			for _, tag := range groups {
				item.CustomerGroups = append(item.CustomerGroups, &TagItem{ID: tag.ID, Name: tag.Name})
			}
		}

		if _, ok := userSaListMap[item.UnionId]; ok {
			//fmt.Println(userSaListMap[item.UnionId])
			userSaInfo := userSaListMap[item.UnionId]
			item.MemberId = userSaInfo.MemberId
			item.Description = userSaInfo.Description
			item.DepartmentId = userSaInfo.DepartmentId
			item.DepartmentName = userSaInfo.DepartmentName
			item.IsKeySa = userSaInfo.IsKeySa // 是否专属导购：0否，1是
			item.SaUserId = userSaInfo.SaUserId
			item.SaUserName = userSaInfo.SaUserName
		}

		// 查询外部联系人专属SA信息
		//userSaInfo, errs := workExternalDao.NewUserRelation().GetUserSAInfoByUnionId(cond.EnterpriseID, item.UnionId)
		//if errs != nil && !errors.Is(errs, gorm.ErrRecordNotFound) {
		//	return nil, errs
		//}
		//
		//item.MemberId = userSaInfo.MemberId
		//item.Description = userSaInfo.Description
		//item.DepartmentId = userSaInfo.DepartmentId
		//item.DepartmentName = userSaInfo.DepartmentName
		//item.IsKeySa = userSaInfo.IsKeySa // 是否专属导购：0否，1是
		//item.SaUserId = userSaInfo.SaUserId
		//item.SaUserName = userSaInfo.SaUserName

		item.PayOrderNum = eu.PayOrder
		item.PayTotal = eu.TotalFee
		item.RefundOrderNum = eu.RefundOrder
		item.RefundTotal = eu.RefundFee
		item.LatestMarketingTime = eu.LastMarketTime
		item.UnionId = eu.UnionId

		list.List = append(list.List, item)

	}

	if len(list.List) == 0 {
		return list, nil
	}
	return list, nil
}

// QueryEnterpriseUser 查询用户信息 （euid or unionid or appid+openid)
func (s *Service) QueryEnterpriseUser(c ctx.Context, enterpriseID int64, source *SourceUserBase) (*EnterpriseUser, error) {
	return s.repo.QueryEnterpriseUser(c, enterpriseID, source)
}

// UpdateEnterpriseUserBase 更新用户基本信息
func (s *Service) UpdateEnterpriseUserBase(c ctx.Context, enterpriseID int64, user *EnterpriseUser, base *EnterpriseUserBase) error {
	return s.repo.UpdateEnterpriseUserBase(c, enterpriseID, user.ID, base)
}

// 无需排序简洁版查找用户列表
func (s *Service) SimpleQueryList(c ctx.Context, cond *ListRequest) (int64, []*EnterpriseUser, error) {
	return s.esRepo.QueryUserList(c, cond)
}

// 无需排序简洁版查找用户列表
func (s *Service) QueryByUnionID(c ctx.Context, cond *CondUnionID) (int64, []*EnterpriseUser, error) {
	return s.esRepo.QueryByUnionID(c, cond)
}
