package userGroupSvc

import (
	"encoding/json"
	"math"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	userGroupDao "micro-scrm/internal/app/dao/user_group"
	memberReq "micro-scrm/internal/app/http/request/member"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	userGroupModel "micro-scrm/internal/app/models/user_group"
	"micro-scrm/internal/app/service"
	modulesSvc "micro-scrm/internal/app/service/modules"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"strings"
	"time"
)

type GroupSvc struct {
	service.Base
	groupDetail *userGroupDao.GroupDetail
}

func NewGroupSvc() *GroupSvc {
	return &GroupSvc{
		groupDetail: userGroupDao.NewGroupDetail(),
	}
}

// TagsHandle 标签类型的分组处理
func (s GroupSvc) TagsHandle(c app.Ctx, enterpriseId, groupId int64) error {

	list, err := s.groupDetail.GetList(c, enterpriseId, groupId)
	if err != nil {
		return err
	}
	checkTag := make(map[int64]bool, 0)
	groupTags := make([]int64, 0)
	IncludeTagGroupKey := make(map[string]bool, 0)
	FilterTagGroupKey := make(map[string]bool, 0)
	FilterTagGroupMap := make(map[int64]int64, 0)
	IncludeTagGroupMap := make(map[int64]int64, 0)
	for _, v := range list {
		if _, ok := checkTag[v.TagID]; !ok {
			checkTag[v.TagID] = true
			groupTags = append(groupTags, v.TagID)
		}
		key := strconv.FormatInt(v.TagGroupID, 10) + "-" + strconv.FormatInt(v.TagID, 10)
		if userGroupModel.UGD_IncludeType_Include == v.IncludeType {
			IncludeTagGroupKey[key] = true
			IncludeTagGroupMap[v.TagGroupID] = v.TagGroupID
		} else {
			FilterTagGroupKey[key] = true
			FilterTagGroupMap[v.TagGroupID] = v.TagGroupID
		}
	}
	request := customer.ListRequest{
		ListCond: customer.ListCond{
			EnterpriseID: enterpriseId,
			TagCond: struct {
				Include   bool    `json:"include"`
				TagIDList []int64 `json:"tag_id_list"`
			}{
				true, groupTags,
			},
		},
		Page: 1,
		Size: 1,
	}
	// 先查询出总数
	total, _, err := modulesSvc.NewCustomerService().EsRepo().QueryUserList(c, &request)
	if err != nil {
		s.Log().Error(c).Msgf("根据标签id集合，查找enterprise_user, 错误：%+v", err)
		return err
	}
	maxPage := int(math.Ceil(float64(total / 10000)))
	if maxPage == 0 {
		maxPage += 1
	}
	for i := 1; i <= maxPage; i++ {
		req := customer.ListRequest{
			ListCond: customer.ListCond{
				EnterpriseID: enterpriseId,
				TagCond: struct {
					Include   bool    `json:"include"`
					TagIDList []int64 `json:"tag_id_list"`
				}{
					true, groupTags,
				},
			},
			Page: i,
			Size: 10000,
		}
		_, lists, err := modulesSvc.NewCustomerService().EsRepo().QueryUserList(c, &req)

		if err != nil {
			s.Log().Error(c).Msgf("根据标签id集合，查找enterprise_user, 错误：%+v", err)
			continue
		}
		userList := lists
		for _, users := range userList {
			includeTagGroup := make(map[int64]int64)
			FilterTagGroup := make(map[int64]int64)
			for _, t := range users.Tags {
				for y := range FilterTagGroupKey { // 同一标签组内的标签为“或”关系，标签组互相为“且”关系
					isContain := strings.Contains(y, strconv.FormatInt(t, 10))
					if !isContain {
						continue
					}
					arr := strings.Split(y, "-")
					tagGroupID, _ := strconv.ParseInt(arr[0], 10, 64)
					FilterTagGroup[tagGroupID] = tagGroupID
				}
				if len(FilterTagGroup) > 0 && len(FilterTagGroup) == len(FilterTagGroupMap) {
					break
				}
				for y := range IncludeTagGroupKey {
					isContain := strings.Contains(y, strconv.FormatInt(t, 10))
					if !isContain {
						continue
					}
					arr := strings.Split(y, "-")
					tagGroupID, _ := strconv.ParseInt(arr[0], 10, 64)
					includeTagGroup[tagGroupID] = tagGroupID
				}
			}
			if len(includeTagGroup) > 0 && len(includeTagGroup) == len(IncludeTagGroupMap) {
				users.Groups = append(users.Groups, groupId)
				users.Groups = helper.ArrayUniq[int64](users.Groups)
				var updateData microScrmModel.EnterpriseUserModel
				jsonGroups, err := json.Marshal(users.Groups)
				if err != nil {
					return err
				}
				updateData.Groups = string(jsonGroups)
				if err = microScrmDao.NewEnterpriseUser().UpdateById(c, updateData, enterpriseId, users.ID); err != nil {
					s.Log().Info(c).Msgf("更新用户的分组:%+v, 错误：%v", users.ID, err)
				}
			}
		}
	}
	return nil
}

// DepartmentHandle 门店类型的分组处理
func (s GroupSvc) DepartmentHandle(c app.Ctx, enterpriseId, groupId int64) error {

	list, err := s.groupDetail.GetList(c, enterpriseId, groupId)
	if err != nil {
		return err
	}
	departmentIds := make([]int64, 0)
	for _, v := range list {
		departmentIds = append(departmentIds, v.DepartmentID)
	}
	if len(departmentIds) < 1 {
		return nil
	}
	pageSize := 100
	offset := 0
	total, dataErr := membersDao.NewMemberDepartments().SelectEUidBYDepartmentCount(c, enterpriseId, departmentIds)
	if dataErr != nil || total == 0 {
		s.Log().Info(c).Msgf("当前Department_iderr: %v，Department_iderr", dataErr)
		return dataErr
	}
	datalist, err := membersDao.NewMemberDepartments().SelectEUidBYDepartment(c, enterpriseId, departmentIds, offset, pageSize)
	if err != nil {
		return err
	}
	if total > int64(pageSize) {
		maxPage := int(math.Ceil(float64(total / int64(pageSize))))
		if maxPage == 0 {
			maxPage += 1
		}
		for i := 1; i <= maxPage; i++ {
			offset = i * pageSize
			list, errs := membersDao.NewMemberDepartments().SelectEUidBYDepartment(c, enterpriseId, departmentIds, offset, pageSize)
			if errs != nil {
				break
			}
			if errs = s.updateGroupsBYEUid(c, list, enterpriseId, groupId); errs != nil {
				continue
			}
		}
	} else {
		upErrs := s.updateGroupsBYEUid(c, datalist, enterpriseId, groupId)
		if upErrs != nil {
			return upErrs
		}
	}
	return nil
}

// ImportDataHandle 客户名单导入类型的分组处理
func (s GroupSvc) ImportDataHandle(c app.Ctx, enterpriseId, fileId, groupId int64) error {
	nowTime := time.Now().Unix()
	count := 0
	cond := &customer.ListCond{EnterpriseID: enterpriseId}
	cond.CustomerGroupCond.Include = true
	cond.CustomerGroupCond.GroupIDList = []int64{groupId}
	totalResult, _ := modulesSvc.NewCustomerService().EsRepo().QueryUserCount(c, cond)
	totalRes := totalResult
	if totalRes != nil {
		count = int(totalRes.MemberCount)
	}
	if count > 0 { // 判断这个分组的数据已经统计过了，就不统计了
		return nil
	}
	groupImport := userGroupDao.NewGroupImportData()
	groupImportLog := userGroupDao.NewGroupImportDataLog()
	total, countErr := groupImport.GetImportDataCount(c, enterpriseId, fileId)
	if countErr != nil {
		return countErr
	}
	pageSize := 100
	maxPage := int(math.Ceil(float64(total / int64(pageSize))))
	if maxPage == 0 {
		maxPage += 1
	}
	for i := 1; i <= maxPage; i++ {
		offset := (i - 1) * pageSize
		datalist, dataErr := groupImport.GetImportData(c, enterpriseId, fileId, offset, pageSize)
		if dataErr != nil {
			continue
		}
		var unionID []string
		for _, data := range datalist {
			unionID = append(unionID, data.UnionID)
		}
		var unionIdMap = make(map[string]struct{}, 0)
		var newDataLogArr []userGroupModel.UserGroupImportDataLogModel
		for _, value := range unionID {
			unionIdMap[value] = struct{}{}
			var newDataLog userGroupModel.UserGroupImportDataLogModel
			newDataLog.CreatedAt = nowTime
			newDataLog.UpdatedAt = nowTime
			newDataLog.UnionID = value
			newDataLog.UserGroupID = groupId
			newDataLog.ImportRes = "匹配失败"
			newDataLogArr = append(newDataLogArr, newDataLog)
		}
		_ = groupImportLog.CreateInBatches(c, enterpriseId, newDataLogArr)
		condUnionId := &customer.CondUnionID{EnterpriseID: enterpriseId}
		condUnionId.UnionId = unionID
		_, userRes, userErr := modulesSvc.NewCustomerService().EsRepo().QueryByEUidPag(c, condUnionId, offset, pageSize)
		users := userRes
		if userErr != nil {
			continue
		}
		for _, esUser := range users {
			if _, ok := unionIdMap[esUser.UnionId]; !ok {
				continue
			}
			var groupsMap = make(map[int64]struct{}, 0)
			for _, item := range esUser.Groups {
				groupsMap[item] = struct{}{}
			}
			if _, ok := groupsMap[groupId]; !ok {
				esUser.Groups = append(esUser.Groups, groupId)
				esUser.Groups = helper.ArrayUniq[int64](esUser.Groups)
				var updateData microScrmModel.EnterpriseUserModel
				jsonGroups, err := json.Marshal(esUser.Groups)
				if err != nil {
					return err
				}
				updateData.Groups = string(jsonGroups)
				esErr := microScrmDao.NewEnterpriseUser().UpdateById(c, updateData, enterpriseId, esUser.ID)
				if esErr != nil {
					s.Log().Info(c).Msgf("更新用户的分组:%+v, 错误：%v", esUser.ID, esErr)
				} else {
					_ = groupImportLog.Update(c, enterpriseId, groupId, esUser.UnionId)
				}
			}
		}
	}
	_ = s.groupDetail.UpdateDetailISExtract(c, enterpriseId, groupId)
	return nil
}

func (s GroupSvc) updateGroupsBYEUid(c app.Ctx, datalist []memberReq.EUidList, enterpriseID, gid int64) error {
	var EUidArr []int64
	EUidMap := make(map[int64]struct{}, 0)
	for _, item := range datalist {
		if item.Euid == 0 {
			continue
		}
		EUidArr = append(EUidArr, item.Euid)
		if _, ok := EUidMap[item.Euid]; !ok {
			EUidMap[item.Euid] = struct{}{}
		}
	}
	condUnionId := &customer.CondUnionID{EnterpriseID: enterpriseID}
	condUnionId.EUid = EUidArr
	_, userResult, userErr := modulesSvc.NewCustomerService().EsRepo().QueryByEUidPag(c, condUnionId, 0, 100)
	users := userResult
	if userErr != nil {
		return userErr
	}
	for _, esUser := range users {
		if _, ok := EUidMap[esUser.ID]; !ok {
			continue
		}
		var groupsMap = make(map[int64]struct{}, 0)
		for _, item := range esUser.Groups {
			groupsMap[item] = struct{}{}
		}
		if _, ok := groupsMap[gid]; !ok {
			esUser.Groups = append(esUser.Groups, gid)
			esUser.Groups = helper.ArrayUniq[int64](esUser.Groups)
			var updateData microScrmModel.EnterpriseUserModel
			jsonGroups, err := json.Marshal(esUser.Groups)
			if err != nil {
				return err
			}
			updateData.Groups = string(jsonGroups)
			if err := microScrmDao.NewEnterpriseUser().UpdateById(c, updateData, enterpriseID, esUser.ID); err != nil {
				s.Log().Info(c).Msgf("更新用户的分组:%+v, 错误：%v", esUser.ID, err)
				return err
			}
		}
	}
	return nil
}
