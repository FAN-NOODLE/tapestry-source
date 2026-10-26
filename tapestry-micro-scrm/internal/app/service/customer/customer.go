package customerSvc

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	customerDao "micro-scrm/internal/app/dao/customer"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	externalUserDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	tagsDao "micro-scrm/internal/app/dao/tags"
	customerReq "micro-scrm/internal/app/http/request/customer"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	customerResp "micro-scrm/internal/app/http/response/customer"
	"micro-scrm/internal/app/models"
	membersModel "micro-scrm/internal/app/models/members"
	workExternalModel "micro-scrm/internal/app/models/micro_scrm"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/app/service"
	tagsSvc "micro-scrm/internal/app/service/tags"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/work"
	"strings"
	"time"
)

type Customer struct {
	service.Base
	wxUserCorpExternalDao *externalUserDao.WxUserCorpExternal
	workExternalDao       *externalUserDao.UserRelation
}

func NewCustomer() *Customer {
	return &Customer{
		wxUserCorpExternalDao: externalUserDao.NewWxUserCorpExternal(),
		workExternalDao:       externalUserDao.NewUserRelation(),
	}
}

func (s *Customer) UpdateRemark(c app.Ctx, request customerReq.UpdateRemarkRequest) error {
	// 查询wx_user_corp_external 的id
	userInfo, err := s.wxUserCorpExternalDao.FindOne(c, request.EnterpriseID, request.ExternalUserID)
	if err != nil {
		return err
	}
	// 修改备注
	err = s.workExternalDao.UpdateRemark(c, request.EnterpriseID, userInfo.ID, request.MemberID, request.Remark)
	if err != nil {
		return err
	}
	relationInfo, err := s.workExternalDao.GetOne(c, map[string]interface{}{
		"enterprise_id":         request.EnterpriseID,
		"member_id":             request.MemberID,
		"work_external_user_id": userInfo.ID,
	})
	if err != nil {
		return err
	}
	// 同步备注到企微
	enterpriseInfo, _ := microBossDao.NewEnterprise().GetOne(c, request.EnterpriseID)
	memberUserInfo, _ := membersDao.NewMembers().GetOne(c, fmt.Sprintf("enterprise_id = %d and id = %d", request.EnterpriseID, request.MemberID))

	qyWork := new(work.QyWorkImpl).SetDomain(enterpriseInfo.Domain)
	updateData := work.EditExternalDescriptionRequest{
		UserID:         memberUserInfo.Userid,
		ExternalUserid: request.ExternalUserID,
		Remark:         request.Remark,
		Description:    relationInfo.Description,
	}
	err = qyWork.UpdateExternalRelationDescription(c, updateData)

	return err

}

func (s *Customer) SimpleInfo(c app.Ctx, request customerReq.SimpleInfoRequest) (workExternalModel.WxUserCorpExternalModel, error) {

	// 查询wx_user_corp_external 的id
	userInfo, err := s.wxUserCorpExternalDao.FindOne(c, request.EnterpriseID, request.ExternalUserID)
	if err != nil {
		return workExternalModel.WxUserCorpExternalModel{}, err
	}

	return userInfo, err

}

func (s *Customer) SetInfo(c app.Ctx, request customerReq.OpenIdInfoReq) error {
	data := workExternalModel.WxUserOpenidRelationModel{
		BaseModel: models.CreateBase(app.UIDContext(c)),
		WxOpenid:  request.OpenID,
		Unionid:   request.UnionId,
	}
	if err := externalUserDao.NewOpenidRelation().SetInfo(c, request.EnterpriseID, data); err != nil {
		return err
	}
	return nil
}

func (s *Customer) GetMemberListsByUnionId(c app.Ctx, request customerReq.MemberListsReq) ([]string, error) {
	var result []string
	res, err := s.workExternalDao.GetMemberIdsByUnionId(c, request.EnterpriseID, request.UnionId)
	if err != nil {
		return nil, err
	}
	result, err = membersDao.NewMembers().GetMobilesByMemberIds(c, request.EnterpriseID, res)
	return result, err
}

// GetExternalUserList 顾客列表
func (s *Customer) GetExternalUserList(c app.Ctx, request customerReq.ExternalUserListsReq, enterpriseConfigName string) ([]workExternalModel.ExternalUserList, int64, error) {

	// 判断当前员工类型 店员或店长
	memberLists, userids, err := s.GetMembersByManageType(c, request.EnterpriseID, request.MemberId, request.Userid, request.MemberType, request.ThirdIds, request.MemberList)
	if err != nil {
		return nil, 0, err
	}
	// 管理的或选择的
	request.MemberList = memberLists
	request.Userids = userids

	// 通用查询列表
	result, total, err := s.workExternalDao.GetUserList(c, request)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return result, total, nil
	}
	// 处理数据
	if enterpriseConfigName == "coach" {
		result, err = s.HandleUserListCoach(c, request, result, enterpriseConfigName)
	}
	if enterpriseConfigName == "sw" {
		result, err = s.HandleUserListSW(c, request, result, enterpriseConfigName)
	}
	if enterpriseConfigName == "ks" {
		result, err = s.HandleUserListKS(c, request, result, enterpriseConfigName)
	}
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil

}

// GetMembersByManageType 根据导购或店长类型获取店员
func (s *Customer) GetMembersByManageType(c app.Ctx, enterpriseID, memberId int64, userid string, memberType int, thirdIds []int, memberList []int64) (memberLists []int64, userids []string, err error) {

	// SM 管理
	if memberType == membersModel.MembersTypeSM {
		memberLists = memberList
		// 未选中店员
		if len(memberList) == 0 {
			// 查询店长管理的门店下店员 GetMemberByThirdIds
			memberIds, err1 := membersDao.NewMemberDepartments().GetMemberIdsByThirdIds(c, enterpriseID, thirdIds, []int64{})
			if err1 != nil {
				return memberLists, userids, err1
			}
			memberIds = helper.ArrayUniq(memberIds)
			// 店长和管理的店员
			memberLists = append(memberLists, memberId)
			memberLists = append(memberLists, memberIds...)
		}
		userids, _ = membersDao.NewMembers().GetUserIdsById(c, enterpriseID, memberLists)
	} else {
		// SA 店员
		memberLists = append(memberLists, memberId)
		userids = append(userids, userid)
	}

	return memberLists, userids, nil

}

// HandleUserListCoach Coach顾客列表
func (s *Customer) HandleUserListCoach(c app.Ctx, request customerReq.ExternalUserListsReq, result []workExternalModel.ExternalUserList, enterpriseConfigName string) (resultList []workExternalModel.ExternalUserList, err error) {
	// 列表的外部联系人WorkExternalUserID集合
	// 列表memberID集合
	var workExternalUserIds []int64
	var memberIds []int64
	var storeCodes []string
	var keySaCode []string
	for _, v := range result {
		workExternalUserIds = append(workExternalUserIds, v.WorkExternalUserID)
		memberIds = append(memberIds, v.MemberID)
		storeCodes = append(storeCodes, v.LastPurchaseStoreCode)
		keySaCode = append(keySaCode, v.SaUserId)
	}
	// 根据列表查询是否是自己好友的数据
	userSARelationsMap, err := s.GetUserSARelationsMap(c, request.EnterpriseID, request.MemberId, workExternalUserIds)
	if err != nil {
		return result, err
	}
	// 整理relation里的member_id导购信息
	memberNameListMap, err := s.GetAllByIdMap(c, request.EnterpriseID, memberIds)
	if err != nil {
		return result, err
	}
	// 整理keysa 的导购信息
	keySaMap := make(map[string]membersModel.MemberModel, 0)
	keySaLists, err := membersDao.NewMembers().GetManyByUserIds(c, request.EnterpriseID, keySaCode)
	if err != nil {
		return result, err
	}
	for _, v := range keySaLists {
		keySaMap[v.Userid] = v
	}

	// 查出外部联系人的key_sa导购
	// userKeySARelationsMap, memberNameMap, err := s.GetUserKeySARelationsMap(c, request.EnterpriseID, workExternalUserIds)
	// if err != nil {
	// 	return result, err
	// }

	// var cardDescriptionMap = make(map[string]string, 0)
	// cardDescriptionMap[""] = "0"
	// cardDescriptionMap["PROSPECT"] = "1"
	// cardDescriptionMap["C"] = "2"
	// cardDescriptionMap["D"] = "3"
	// cardDescriptionMap["GOLD"] = "4"
	// cardDescriptionMap["SILVER"] = "5"
	// cardDescriptionMap["DIAMOND"] = "6"
	cardDescriptionMap, _ := s.GetCustomerCardMapForType(enterpriseConfigName)

	for k, v := range result {

		if _, ok := cardDescriptionMap[v.CardLevel]; ok {
			result[k].CustomerCard = cardDescriptionMap[v.CardLevel]
		}
		result[k].Last12mPurchaseAmount = float32(v.AmountTotal / 100)
		result[k].ClientType = v.CltType
		result[k].CustomerCode = v.UnionCode

		// 添加导购姓名
		if _, ok := memberNameListMap[v.MemberID]; ok {
			result[k].MemberName = memberNameListMap[v.MemberID].Name
		}

		// 所属顾问是 key_sa导购 并且是当前导购为专属
		result[k].IsKeySa = 0
		// 有好友关系
		if _, ok := userSARelationsMap[v.WorkExternalUserID]; ok {
			result[k].IsRelation = 1
		}

		if request.ManageType != workExternalModel.ManageNotSA {
			if _, ok := keySaMap[v.SaUserId]; ok {
				result[k].KeySaMemberID = keySaMap[v.SaUserId].ID
				result[k].KeySaMemberName = keySaMap[v.SaUserId].Name
				if result[k].KeySaMemberID == request.MemberId {
					result[k].IsKeySa = 1 // 自己专属标识
				}
			}
		}
	}

	return result, nil
}

// HandleUserListSW SW顾客列表
func (s *Customer) HandleUserListSW(c app.Ctx, request customerReq.ExternalUserListsReq, result []workExternalModel.ExternalUserList, enterpriseConfigName string) (resultList []workExternalModel.ExternalUserList, err error) {
	var workExternalUserIds []int64
	var memberIds []int64
	var storeCodes []string
	for _, v := range result {
		workExternalUserIds = append(workExternalUserIds, v.WorkExternalUserID)
		memberIds = append(memberIds, v.MemberID)

		if v.LastPurchaseStoreCode != "" {
			storeCodes = append(storeCodes, v.LastPurchaseStoreCode)
		}
	}
	memberNameListMap, err := s.GetAllByIdMap(c, request.EnterpriseID, memberIds)
	if err != nil {
		return result, err
	}

	userSARelationsMap, err := s.GetUserSARelationsMap(c, request.EnterpriseID, request.MemberId, workExternalUserIds)
	if err != nil {
		return result, err
	}
	storeDepartmentInfoMap, err := s.GetStoreDepartmentInfoMap(c, request.EnterpriseID, storeCodes)
	if err != nil {
		return result, err
	}

	// var cardDescriptionMap = make(map[string]string, 0)
	// cardDescriptionMap[""] = "0"
	// cardDescriptionMap["Z"] = "1"
	// cardDescriptionMap["B"] = "2"
	// cardDescriptionMap["C"] = "3"
	// cardDescriptionMap["W"] = "4"
	cardDescriptionMap, _ := s.GetCustomerCardMapForType(enterpriseConfigName)

	for k, v := range result {

		// sw
		if _, ok := cardDescriptionMap[v.CardLevel]; ok {
			result[k].CustomerCard = cardDescriptionMap[v.CardLevel]
		}
		result[k].Last12mPurchaseAmount = float32(v.AmountTotal / 100)
		result[k].ClientType = v.ClientType
		result[k].CustomerCode = v.UnionCode

		// 生日解密
		result[k].Birthday, _ = helper.DecodeString(v.Birthday)

		// 添加导购姓名
		if _, ok := memberNameListMap[v.MemberID]; ok {
			result[k].MemberName = memberNameListMap[v.MemberID].Name
		}

		result[k].AmountTotal = v.AmountTotal / 100
		if v.LastPurchaseDate > 0 {
			result[k].LastPurchaseDateStr = time.Unix(int64(v.LastPurchaseDate), 0).Format("2006-01-02 15:04:05")
		}
		if v.LastPurchaseStoreCode != "" {
			lastPurchaseStoreCode := strings.ToUpper(v.LastPurchaseStoreCode)
			if _, ok := storeDepartmentInfoMap[lastPurchaseStoreCode]; ok {
				result[k].LastPurchaseStoreName = storeDepartmentInfoMap[lastPurchaseStoreCode].StoreName
			}
		}
		result[k].IsKeySa = 0
		// 有好友关系
		if _, ok := userSARelationsMap[v.WorkExternalUserID]; ok {
			result[k].IsRelation = 1
		}

	}

	return result, nil
}

// HandleUserListKS KS顾客列表
func (s *Customer) HandleUserListKS(c app.Ctx, request customerReq.ExternalUserListsReq, result []workExternalModel.ExternalUserList, enterpriseConfigName string) (resultList []workExternalModel.ExternalUserList, err error) {
	var workExternalUserIds []int64
	var memberIds []int64
	var storeCodes []string
	for _, v := range result {
		workExternalUserIds = append(workExternalUserIds, v.WorkExternalUserID)
		memberIds = append(memberIds, v.MemberID)

		if v.LastPurchaseStoreCode != "" {
			storeCodes = append(storeCodes, v.LastPurchaseStoreCode)
		}
	}
	memberNameListMap, err := s.GetAllByIdMap(c, request.EnterpriseID, memberIds)
	if err != nil {
		return result, err
	}

	userSARelationsMap, err := s.GetUserSARelationsMap(c, request.EnterpriseID, request.MemberId, workExternalUserIds)
	if err != nil {
		return result, err
	}
	storeDepartmentInfoMap, err := s.GetStoreDepartmentInfoMap(c, request.EnterpriseID, storeCodes)
	if err != nil {
		return result, err
	}

	// var cardDescriptionMap = make(map[string]string, 0)
	// cardDescriptionMap[""] = ""
	// cardDescriptionMap["0"] = "0"
	// cardDescriptionMap["1"] = "1"
	// cardDescriptionMap["2"] = "2"
	// cardDescriptionMap["3"] = "3"
	cardDescriptionMap, _ := s.GetCustomerCardMapForType(enterpriseConfigName)

	for k, v := range result {

		// ks 为空 和 有值
		if _, ok := cardDescriptionMap[v.CardLevel]; ok {
			result[k].CustomerCard = cardDescriptionMap[v.CardLevel]
		}

		result[k].Last36mPurchaseAmount = float32(v.AmountTotal / 100)
		result[k].CustomerCode = v.UnionCode

		// 添加导购姓名
		if _, ok := memberNameListMap[v.MemberID]; ok {
			result[k].MemberName = memberNameListMap[v.MemberID].Name
		}

		result[k].AmountTotal = v.AmountTotal / 100
		if v.LastPurchaseDate > 0 {
			result[k].LastPurchaseDateStr = time.Unix(int64(v.LastPurchaseDate), 0).Format("2006-01-02 15:04:05")
		}
		if v.LastPurchaseStoreCode != "" {
			lastPurchaseStoreCode := strings.ToUpper(v.LastPurchaseStoreCode)
			if _, ok := storeDepartmentInfoMap[lastPurchaseStoreCode]; ok {
				result[k].LastPurchaseStoreName = storeDepartmentInfoMap[lastPurchaseStoreCode].StoreName
			}
		}
		result[k].IsKeySa = 0
		// 有好友关系
		if _, ok := userSARelationsMap[v.WorkExternalUserID]; ok {
			result[k].IsRelation = 1
		}
	}

	return result, nil
}

// GetUserSARelationsMap 当前导购是否和顾客有好友关系
func (s *Customer) GetUserSARelationsMap(c app.Ctx, enterpriseID, memberId int64, workExternalUserIds []int64) (map[int64]workExternalModel.UserSaInfo, error) {

	var userSARelationsMap = make(map[int64]workExternalModel.UserSaInfo, 0)
	userSARelations, err := s.workExternalDao.GetUserSARelations(c, enterpriseID, memberId, workExternalUserIds)
	if err != nil {
		return userSARelationsMap, err
	}
	for _, v := range userSARelations {
		userSARelationsMap[v.WorkExternalUserId] = v
	}

	return userSARelationsMap, nil
}

// GetAllByIdMap 查询添加导购信息
func (s *Customer) GetAllByIdMap(c app.Ctx, enterpriseID int64, memberIds []int64) (map[int64]membersModel.MemberModel, error) {

	var memberNameListMap = make(map[int64]membersModel.MemberModel, 0)
	memberInfoLists, err := membersDao.NewMembers().GetAllById(c, enterpriseID, memberIds)
	if err != nil {
		return memberNameListMap, err
	}
	for _, item := range memberInfoLists {
		memberNameListMap[item.ID] = item
	}
	return memberNameListMap, nil
}

// GetUserKeySARelationsMap 查询Key_sa
func (s *Customer) GetUserKeySARelationsMap(c app.Ctx, enterpriseID int64, workExternalUserIds []int64) (map[int64]workExternalModel.UserSaInfo, map[int64]membersModel.MemberModel, error) {

	var userKeySARelationsMap = make(map[int64]workExternalModel.UserSaInfo, 0)
	var memberNameMap = make(map[int64]membersModel.MemberModel, 0)
	userKeySARelations, err := s.workExternalDao.GetUserKeySARelations(c, enterpriseID, workExternalUserIds)
	if err != nil {
		return userKeySARelationsMap, memberNameMap, err
	}
	var keySAMemberList []int64
	for _, v := range userKeySARelations {
		userKeySARelationsMap[v.WorkExternalUserId] = v
		keySAMemberList = append(keySAMemberList, v.MemberId)
	}

	// 查询key_sa导购信息
	if len(keySAMemberList) > 0 {
		memberLists, err := membersDao.NewMembers().GetAllById(c, enterpriseID, keySAMemberList)
		if err != nil {
			return userKeySARelationsMap, memberNameMap, err
		}
		for _, item := range memberLists {
			memberNameMap[item.ID] = item
		}
	}
	return userKeySARelationsMap, memberNameMap, nil
}

// GetStoreDepartmentInfoMap 查询门店信息-sw
func (s *Customer) GetStoreDepartmentInfoMap(c app.Ctx, enterpriseID int64, storeCodes []string) (map[string]microWorkModel.StoreDepartmentModel, error) {

	storeDepartmentInfo, err := microWorkDao.NewStoreDepartment().GetStoreInfoByStoreCode(c, enterpriseID, storeCodes)
	var storeDepartmentInfoMap = make(map[string]microWorkModel.StoreDepartmentModel, 0)
	if err != nil {
		return storeDepartmentInfoMap, err
	}
	for _, item := range storeDepartmentInfo {
		storeDepartmentInfoMap[strings.ToUpper(item.ThirdStoreCode)] = item
	}

	return storeDepartmentInfoMap, nil
}

func (s *Customer) GetInfoByUnionId(c app.Ctx, request customerReq.UserInfoRequest) (workExternalModel.WxUserCorpExternalModel, error) {

	userInfo, err := s.wxUserCorpExternalDao.GetInfoByUnionId(c, request.EnterpriseID, request.UnionId)
	if err != nil {
		return workExternalModel.WxUserCorpExternalModel{}, err
	}

	return userInfo, err

}

func (s *Customer) GetUserSAByExternalUserIds(c app.Ctx, request customerReq.UserSaRequest) (info []workExternalModel.UserSaData, err error) {

	info, err = s.workExternalDao.GetUserSAByExternalUserIds(c, request.EnterpriseID, request.WorkExternalUserID)
	if err != nil {
		return []workExternalModel.UserSaData{}, err
	}

	return info, err

}

func (s *Customer) GetUserDetailTips(c app.Ctx, request customerReq.UserDetailTipsRequest, memberID int64, userId string) (res customerReq.UserDetailTipsData, err error) {

	reserveNum, err1 := externalUserDao.NewReservesOrder().GetReserveTipsCount(c, request.EnterpriseID, request.ExternalUserID)
	if err1 != nil {
		return res, err1
	}
	res.ReserveNum = reserveNum
	// 外部联系人的 WorkExternalUserID
	userInfo, errs := externalUserDao.NewWxUserCorpExternal().FindOne(c, request.EnterpriseID, request.ExternalUserID)
	if errs != nil {
		return res, errs
	}
	var memberIDTmp []int64
	memberIDTmp = append(memberIDTmp, memberID)
	sowTaskNum, err2 := customerDao.NewSowTask().GetSowTaskTipsCount(c, request.EnterpriseID, userInfo.ID, memberIDTmp)
	if err2 != nil {
		return res, err2
	}

	res.SowTaskNum = sowTaskNum
	// 获取未完成的任务数量 三个月的
	taskIds, _ := externalUserDao.NewPlanTaskDao().GetNoFinishTaskNums(c, request.EnterpriseID, userId)
	if len(taskIds) > 0 {
		// 再获取与该导购有关的
		planTaskNum, _ := externalUserDao.NewPlanTaskConditionDataCkDao().GetTaskAllotLogCount(c, request.EnterpriseID, map[string]interface{}{
			"task_id":  taskIds,
			"user_id":  userId,
			"union_id": userInfo.Unionid,
		})
		res.PlanTaskNum = planTaskNum
	}
	return res, nil

}

func (s *Customer) GetUserSAListByExternalUserIds(c app.Ctx, request customerReq.UserBindRelationRequest, enterpriseConfigName string) (list []workExternalModel.UserSaListData, err error) {

	// coach 不验证管理权限；sw和ks只有管理可以看绑定关系，并且验证管理范围
	if enterpriseConfigName == "sw" || enterpriseConfigName == "ks" {
		// 判断当前员工类型 店员或店长
		memberList, _, err := s.GetMembersByManageType(c, request.EnterpriseID, request.MemberId, "", request.MemberType, request.ThirdIds, request.MemberList)
		if err != nil {
			return list, err
		}
		request.MemberList = memberList
	}

	// 外部联系人的 WorkExternalUserID
	userInfo, errs := externalUserDao.NewWxUserCorpExternal().FindOne(c, request.EnterpriseID, request.ExternalUserID)
	if errs != nil {
		return list, errs
	}
	list, err = externalUserDao.NewUserRelation().GetUserSAListByExternalUserIds(c, request.EnterpriseID, userInfo.ID, request.MemberList)

	if err != nil {
		return list, err
	}
	return list, nil

}

// GetUserSARelation 查询外部联系人和导购的关系
func (s *Customer) GetUserSARelation(c app.Ctx, enterpriseId, memberId, workExternalUserID int64) (info workExternalModel.UserSaInfo, err error) {

	info, err = externalUserDao.NewUserRelation().GetUserSARelation(c, enterpriseId, memberId, workExternalUserID)

	if err != nil {
		return info, err
	}
	return info, nil

}

// GetUserSARelations 查询外部联系人和导购的关系
func (s *Customer) GetUserSARelations(c app.Ctx, enterpriseId, memberId int64, workExternalUserIDs []int64) (list []workExternalModel.UserSaInfo, err error) {

	list, err = s.workExternalDao.GetUserSARelations(c, enterpriseId, memberId, workExternalUserIDs)

	if err != nil {
		return list, err
	}
	return list, nil

}

func (s *Customer) GetCustomerList(c app.Ctx, request customerReq.CustomerListsReq) ([]workExternalModel.ExternalUserList, error) {

	// 判断当前员工类型 店员或店长
	if request.MemberType == membersModel.MembersTypeSM { // SM
		// 查询店长管理的门店下店员 GetMemberByThirdIds
		memberIds, err1 := membersDao.NewMemberDepartments().GetMemberIdsByThirdIds(c, request.EnterpriseID, request.ThirdIds, []int64{})
		if err1 != nil {
			return nil, err1
		}
		memberIds = helper.ArrayUniq(memberIds)
		var memberListIds []int64
		for _, v := range memberIds {
			memberListIds = append(memberListIds, v)
		}
		// 店长和管理的店员
		request.MemberList = append(request.MemberList, request.MemberId)
		request.MemberList = append(request.MemberList, memberListIds...)
	} else { // SA 店员
		request.MemberList = append(request.MemberList, request.MemberId)
	}

	// 根据三方会员卡号查询
	// result, err := s.wxUserCorpExternalDao.GetWxUserCorpExternalByUnionCode(c, request)
	result, err := s.workExternalDao.GetUserListForUnionCode(c, request)
	if err != nil {
		return nil, err
	}

	// 列表的外部联系人WorkExternalUserID集合
	// 列表memberID集合
	var workExternalUserIds []int64
	// var memberIds []int64
	for _, v := range result {
		workExternalUserIds = append(workExternalUserIds, v.WorkExternalUserID)
	}

	// 当前导购是否和会员有好友关系
	userSARelationsMap := make(map[int64]workExternalModel.UserSaInfo, 0)
	userSARelations, _ := s.workExternalDao.GetUserSARelations(c, request.EnterpriseID, request.MemberId, workExternalUserIds)
	for _, v := range userSARelations {
		userSARelationsMap[v.WorkExternalUserId] = v
	}

	// 查询添加导购信息
	// memberNameListMap := make(map[int64]membersModel.MemberModel, 0)
	// memberInfoLists, err := membersDao.NewMembers().GetAllById(c, request.EnterpriseID, memberIds)
	// for _, item := range memberInfoLists {
	// 	memberNameListMap[item.ID] = item
	// }

	// 查询Key_sa
	userKeySARelationsMap := make(map[int64]workExternalModel.UserSaInfo, 0)
	userKeySARelations, _ := s.workExternalDao.GetUserKeySARelations(c, request.EnterpriseID, workExternalUserIds)
	var keyMemberList []int64
	for _, v := range userKeySARelations {
		userKeySARelationsMap[v.WorkExternalUserId] = v
		keyMemberList = append(keyMemberList, v.MemberId)
	}

	// 查询key_sa导购信息
	memberNameMap := make(map[int64]membersModel.MemberModel, 0)
	memberLists, err := membersDao.NewMembers().GetAllById(c, request.EnterpriseID, keyMemberList)
	for _, item := range memberLists {
		memberNameMap[item.ID] = item
	}

	for k, v := range result {
		// 添加导购姓名
		// if _, ok := memberNameListMap[v.MemberID]; ok {
		// 	result[k].MemberName = memberNameListMap[v.MemberID].Name
		// }

		// 所属顾问是 key_sa导购 并且是当前导购为专属
		result[k].IsKeySa = 0
		// 有好友关系
		if _, ok := userSARelationsMap[v.WorkExternalUserID]; ok {
			result[k].IsRelation = 1
		}

		// 查询专属Key_sa
		if _, ok := userKeySARelationsMap[v.WorkExternalUserID]; ok {
			result[k].KeySaMemberID = userKeySARelationsMap[v.WorkExternalUserID].MemberId
			if _, ok := memberNameMap[result[k].KeySaMemberID]; ok {
				result[k].KeySaMemberName = memberNameMap[result[k].KeySaMemberID].Name
			}
		}
	}
	return result, nil
}

func (s *Customer) CustomerDetail(c app.Ctx, enterpriseId int64, openid string) (userInfo workExternalModel.WxUserCorpExternalModel, err error) {

	// 查询wx_user_corp_external 的id
	userInfo, err = s.wxUserCorpExternalDao.FindOne(c, enterpriseId, openid)
	if err != nil {
		return userInfo, err
	}

	return
}

func (s *Customer) GetCustomerDetail(c app.Ctx, enterpriseId int64, whereConf interface{}) (userInfo workExternalModel.WxUserCorpExternalModel, err error) {

	userInfo, err = s.wxUserCorpExternalDao.GetOne(c, whereConf, enterpriseId)
	if err != nil {
		return userInfo, err
	}

	return
}

// CustomerDetailCoach coach的客户详情
func (s *Customer) CustomerDetailCoach(c app.Ctx, request customerReq.ExternalUserDetailReq, userInfo workExternalModel.WxUserCorpExternalModel, userSARelation workExternalModel.UserSaInfo, memberIds []int64, memberUserIds []string) (customerReq.UserDetailInfoCoach, error) {

	var userDetail customerReq.UserDetailInfoCoach

	nickName := userInfo.Nick
	if nickName == "" {
		nickName = userInfo.Name
	}
	userDetail.UserID = userInfo.Euid
	userDetail.UnionID = userInfo.Unionid
	userDetail.ExternalUserID = userInfo.Openid
	userDetail.WorkExternalUserID = userInfo.ID
	userDetail.NickName = nickName
	userDetail.Avatar = userInfo.Avatar

	userDetail.UnionCode = userInfo.UnionCode
	userDetail.Birthday = userInfo.Birthday
	userDetail.CltType = userInfo.CltType
	userDetail.CardLevel = userInfo.CardLevel
	userDetail.AmountTotal = userInfo.AmountTotal

	// 关系和备注名
	userDetail.RelationID = userSARelation.Id
	userDetail.Remark = userSARelation.Remark
	if userSARelation.Id > 0 {
		userDetail.IsRelation = 1
	}
	var wg threading.RoutineGroup

	wg.Run(func() {

		// 企微标签
		var userTagsReq = tagsReq.UserTagsReq{}
		userTagsReq.EnterpriseID = request.EnterpriseId
		userTagsReq.WorkExternalUserId = userInfo.ID // id

		userTagsReq.MemberIds = memberIds
		userTagsList, err := tagsSvc.NewWorkTags().UserTagsAndTagGroup(c, userTagsReq)
		if err == nil {
			userDetail.UserTags = userTagsList
		}

	})

	wg.Run(func() {

		// 印象标签
		var userImpressionTagsReq = tagsReq.UserTagsReq{}
		userImpressionTagsReq.EnterpriseID = request.EnterpriseId
		userImpressionTagsReq.Userids = memberUserIds
		userImpressionTagsReq.WorkExternalUserId = userInfo.ID // euid
		userImpressionTagsList, err := tagsSvc.NewImpression().TagsAndTagGroupList(c, userImpressionTagsReq)

		if err == nil {
			userDetail.UserImpressionTags = userImpressionTagsList
		}
	})
	// 查询企微的沟通意愿标签
	wg.Run(func() {
		userChatTags, err := tagsSvc.NewWorkTags().GetUserChatTags(c, request.EnterpriseId, request.MemberId, userInfo.ID)
		if err == nil {
			userDetail.UserChatTags = userChatTags
		}
	})
	wg.Wait()

	return userDetail, nil
}

// CustomerDetailSW Sw的客户详情
func (s *Customer) CustomerDetailSW(c app.Ctx, request customerReq.ExternalUserDetailReq, userInfo workExternalModel.WxUserCorpExternalModel, userSARelation workExternalModel.UserSaInfo, memberIds []int64, memberUserIds []string) (customerReq.UserDetailInfoSW, error) {

	var userDetail customerReq.UserDetailInfoSW

	nickName := userInfo.Nick
	if nickName == "" {
		nickName = userInfo.Name
	}
	userDetail.UserID = userInfo.Euid
	userDetail.UnionID = userInfo.Unionid
	userDetail.ExternalUserID = userInfo.Openid
	userDetail.WorkExternalUserID = userInfo.ID
	userDetail.NickName = nickName
	userDetail.Avatar = userInfo.Avatar

	userDetail.UnionCode = userInfo.UnionCode
	userDetail.Birthday = userInfo.Birthday
	userDetail.CltType = userInfo.CltType
	userDetail.CardLevel = userInfo.CardLevel
	userDetail.AmountTotal = userInfo.AmountTotal

	// 关系和备注名
	userDetail.RelationID = userSARelation.Id
	userDetail.Remark = userSARelation.Remark
	if userSARelation.Id > 0 {
		userDetail.IsRelation = 1
	}

	userDetail.AmountTotal = userInfo.AmountTotal / 100
	userDetail.LastPurchaseDate = userInfo.LastPurchaseDate
	userDetail.LastPurchaseStoreCode = userInfo.LastPurchaseStoreCode
	userDetail.UpgradeNeedAmount = userInfo.UpgradeNeedAmount / 100
	userDetail.CustomerType = userInfo.CustomerType

	if userDetail.LastPurchaseDate > 0 {
		userDetail.LastPurchaseDateStr = time.Unix(int64(userDetail.LastPurchaseDate), 0).Format("2006-01-02 15:04:05")
	}

	var wg threading.RoutineGroup

	wg.Run(func() {

		// 企微标签
		var userTagsReq = tagsReq.UserTagsReq{}
		userTagsReq.EnterpriseID = request.EnterpriseId
		userTagsReq.WorkExternalUserId = userInfo.ID // id

		userTagsReq.MemberIds = memberIds
		userTagsList, err := tagsSvc.NewWorkTags().UserTagsAndTagGroup(c, userTagsReq)
		if err == nil {
			userDetail.UserTags = userTagsList
		}

	})

	wg.Run(func() {

		// 印象标签
		var userImpressionTagsReq = tagsReq.UserTagsReq{}
		userImpressionTagsReq.EnterpriseID = request.EnterpriseId
		userImpressionTagsReq.Userids = memberUserIds
		userImpressionTagsReq.WorkExternalUserId = userInfo.ID
		userImpressionTagsList, err := tagsSvc.NewImpression().TagsAndTagGroupList(c, userImpressionTagsReq)

		if err == nil {
			for _, impressionGroup := range userImpressionTagsList {
				if impressionGroup.GroupName == "沟通意愿" {
					userDetail.UserImpressionChatTags = append(userDetail.UserImpressionChatTags, impressionGroup)
				} else {
					userDetail.UserImpressionTags = append(userDetail.UserImpressionTags, impressionGroup)
				}
			}
		}
	})

	// 查询购买门店信息
	wg.Run(func() {

		if userInfo.LastPurchaseStoreCode != "" {
			var storeCodes []string
			storeDepartmentInfoMap, err := s.GetStoreDepartmentInfoMap(c, request.EnterpriseId, append(storeCodes, userInfo.LastPurchaseStoreCode))
			if err == nil {
				lastPurchaseStoreCode := strings.ToUpper(userInfo.LastPurchaseStoreCode)
				if _, ok := storeDepartmentInfoMap[lastPurchaseStoreCode]; ok {
					userDetail.LastPurchaseStoreName = storeDepartmentInfoMap[lastPurchaseStoreCode].StoreName
				}
			}

		}
	})

	wg.Run(func() {

		// 顾客鞋码生日地址
		var userInfoAddReq customerReq.GetCustomerInfoAddReq
		userInfoAddReq.EnterpriseID = request.EnterpriseId
		userInfoAddReq.UnionID = userInfo.Unionid
		userInfoAddRes, err := NewCustomerInfoAdd().GetInfoAdd(ctx.Wrap(c), userInfoAddReq)
		if err == nil {
			userDetail.UserInfoAdd = userInfoAddRes
		}

	})

	wg.Run(func() {
		// CRM标签
		extraSwModel, _ := externalSearchDao.NewExternalExtraCkSw().GetOne(c, request.EnterpriseId, "user_unionid= ?", userInfo.Unionid)
		if extraSwModel.AgeTag != "" {
			userDetail.CrmTag = append(userDetail.CrmTag, extraSwModel.AgeTag)
		}
		if extraSwModel.PriceTag != "" {
			userDetail.CrmTag = append(userDetail.CrmTag, extraSwModel.PriceTag)
		}
	})

	wg.Wait()

	return userDetail, nil
}

// CustomerDetailKS  Ks的客户详情
func (s *Customer) CustomerDetailKS(c app.Ctx, request customerReq.ExternalUserDetailReq, userInfo workExternalModel.WxUserCorpExternalModel, userSARelation workExternalModel.UserSaInfo, memberIds []int64, memberUserIds []string) (customerReq.UserDetailInfoKS, error) {

	var userDetail customerReq.UserDetailInfoKS

	nickName := userInfo.Nick
	if nickName == "" {
		nickName = userInfo.Name
	}
	userDetail.UserID = userInfo.Euid
	userDetail.UnionID = userInfo.Unionid
	userDetail.ExternalUserID = userInfo.Openid
	userDetail.WorkExternalUserID = userInfo.ID
	userDetail.NickName = nickName
	userDetail.Avatar = userInfo.Avatar

	userDetail.UnionCode = userInfo.UnionCode
	userDetail.Birthday = userInfo.Birthday
	userDetail.CardLevel = userInfo.CardLevel
	userDetail.AmountTotal = userInfo.AmountTotal

	// 关系和备注名
	userDetail.RelationID = userSARelation.Id
	userDetail.Remark = userSARelation.Remark
	if userSARelation.Id > 0 {
		userDetail.IsRelation = 1
	}

	userDetail.AmountTotal = userInfo.AmountTotal / 100
	userDetail.LastPurchaseDate = userInfo.LastPurchaseDate
	userDetail.LastPurchaseStoreCode = userInfo.LastPurchaseStoreCode
	userDetail.UpgradeNeedAmount = userInfo.UpgradeNeedAmount / 100
	userDetail.CustomerType = userInfo.CustomerType

	if userDetail.LastPurchaseDate > 0 {
		userDetail.LastPurchaseDateStr = time.Unix(int64(userDetail.LastPurchaseDate), 0).Format("2006-01-02 15:04:05")
	}

	var wg threading.RoutineGroup

	wg.Run(func() {

		// 企微标签
		var userTagsReq = tagsReq.UserTagsReq{}
		userTagsReq.EnterpriseID = request.EnterpriseId
		userTagsReq.WorkExternalUserId = userInfo.ID // id

		userTagsReq.MemberIds = memberIds
		userTagsList, err := tagsSvc.NewWorkTags().UserTagsAndTagGroup(c, userTagsReq)
		if err == nil {
			userDetail.UserTags = userTagsList
		}

	})

	wg.Run(func() {

		// 印象标签
		var userImpressionTagsReq = tagsReq.UserTagsReq{}
		userImpressionTagsReq.EnterpriseID = request.EnterpriseId
		userImpressionTagsReq.Userids = memberUserIds
		userImpressionTagsReq.WorkExternalUserId = userInfo.ID // euid
		userImpressionTagsList, err := tagsSvc.NewImpression().TagsAndTagGroupList(c, userImpressionTagsReq)

		if err == nil {
			userDetail.UserImpressionTags = userImpressionTagsList
		}
	})

	// 查询购买门店信息
	wg.Run(func() {

		if userInfo.LastPurchaseStoreCode != "" {
			var storeCodes []string
			storeDepartmentInfoMap, err := s.GetStoreDepartmentInfoMap(c, request.EnterpriseId, append(storeCodes, userInfo.LastPurchaseStoreCode))
			if err == nil {
				lastPurchaseStoreCode := strings.ToUpper(userInfo.LastPurchaseStoreCode)
				if _, ok := storeDepartmentInfoMap[lastPurchaseStoreCode]; ok {
					userDetail.LastPurchaseStoreName = storeDepartmentInfoMap[lastPurchaseStoreCode].StoreName
				}
			}

		}
	})

	wg.Wait()

	return userDetail, nil
}

// GetCustomerCardMapForType 各企业的会员卡等类型
func (s *Customer) GetCustomerCardMapForType(storeName string) (map[string]string, error) {
	var cardDescriptionMap = make(map[string]string, 0)

	if storeName == "coach" {
		cardDescriptionMap[""] = "0"
		cardDescriptionMap["PROSPECT"] = "1"
		cardDescriptionMap["C"] = "2"
		cardDescriptionMap["D"] = "3"
		cardDescriptionMap["GOLD"] = "4"
		cardDescriptionMap["SILVER"] = "5"
		cardDescriptionMap["DIAMOND"] = "6"
	}

	if storeName == "sw" {
		cardDescriptionMap[""] = "0"
		cardDescriptionMap["Z"] = "1"
		cardDescriptionMap["B"] = "2"
		cardDescriptionMap["C"] = "3"
		cardDescriptionMap["W"] = "4"
	}
	if storeName == "ks" {
		cardDescriptionMap[""] = ""
		cardDescriptionMap["0"] = "0"
		cardDescriptionMap["1"] = "1"
		cardDescriptionMap["2"] = "2"
		cardDescriptionMap["3"] = "3"
	}

	return cardDescriptionMap, nil
}

func (s *Customer) UserList(c app.Ctx, request customerReq.UserList) ([]customerResp.UserList, int64, error) {
	var resp []customerResp.UserList
	userList, total, err := externalUserDao.NewWxUserCorpExternalCk().GetPagination(c, request)

	if err != nil {
		return resp, 0, err
	}

	for _, user := range userList {
		tmpResp := customerResp.UserList{
			Id:        user.ID,
			Avatar:    user.Avatar,
			Birthday:  user.Birthday,
			Grade:     user.Grade,
			Mobile:    user.Mobile,
			Name:      user.Name,
			Nick:      user.Nick,
			UnionCode: user.UnionCode,
			UnionId:   user.Unionid,
		}
		resp = append(resp, tmpResp)
	}

	return resp, total, err
}

func (s *Customer) UserDetail(c app.Ctx, enterpriseId, userId int64) (customerResp.UserDetail, error) {
	var resDetail customerResp.UserDetail
	enterpriseUser, err := externalUserDao.NewWxUserCorpExternalCk().GetOne(c, map[string]any{"id": userId}, enterpriseId)
	if err != nil {
		return resDetail, err
	}

	resDetail = customerResp.UserDetail{
		UserID:        enterpriseUser.ID,
		UnionID:       enterpriseUser.Unionid,
		UserName:      enterpriseUser.Name,
		Nickname:      enterpriseUser.Nick,
		Avatar:        enterpriseUser.Avatar,
		Gender:        int(enterpriseUser.Sex),
		Birthday:      enterpriseUser.Birthday,
		UnionCode:     enterpriseUser.UnionCode,
		Point:         enterpriseUser.Point,
		Grade:         enterpriseUser.Grade,
		StoreCode:     enterpriseUser.StoreCode,
		BACode:        enterpriseUser.BaCode,
		Mobile:        enterpriseUser.Mobile,
		TagList:       make([]string, 0),
		ImpressionTag: make([]string, 0),
		GroupList:     make([]string, 0),
	}

	relationIds, _ := externalUserDao.NewUserRelation().GetUserRelationIds(c, enterpriseId, []int64{userId})

	// 获取用户企业标签
	tagIds, _ := tagsDao.NewWorkExternalUserTag().GetTagIdsByUserId(c, enterpriseId, relationIds)
	if len(tagIds) > 0 {
		tagRes, _ := tagsDao.NewWorkTagsDao().GetListByIds(c, enterpriseId, tagIds)
		for _, re := range tagRes {
			resDetail.TagList = append(resDetail.TagList, re.Name)
		}
	}

	// 获取印象标签
	impressionTagIds, _ := tagsDao.NewImpressionRelation().GetTagIdsByUserId(c, enterpriseId, userId)
	if len(impressionTagIds) > 0 {
		impressionRes, _ := tagsDao.NewImpression().GetListByTagId(c, impressionTagIds)
		for _, re := range impressionRes {
			resDetail.ImpressionTag = append(resDetail.ImpressionTag, re.Name)
		}
	}

	return resDetail, nil
}

// GetExternalUserTotalCk 获取全选的客户关系数
func (s *Customer) GetExternalUserTotalCk(c ctx.Context, request customerReq.ExternalUserListBaseReq) (int64, error) {
	memberLists, userids, err := s.GetMembersByManageType(c, request.EnterpriseID, request.MemberId, request.Userid, request.MemberType, request.ThirdIds, request.MemberList)
	if err != nil {
		return 0, err
	}
	request.Userids = userids
	request.MemberLists = memberLists
	return s.workExternalDao.GetUserTotalCk(c, request)
}

// GetExternalUserIdsCk 获取全选的客户关系的external_user_id
func (s *Customer) GetExternalUserIdsCk(c ctx.Context, request customerReq.ExternalUserListBaseReq) ([]string, error) {
	memberLists, userids, err := s.GetMembersByManageType(c, request.EnterpriseID, request.MemberId, request.Userid, request.MemberType, request.ThirdIds, request.MemberList)
	if err != nil {
		return nil, err
	}
	request.MemberLists = memberLists
	return s.workExternalDao.GetExternalUserIdsCk(c, request, userids)
}

func (s *Customer) UpdateDescription(c app.Ctx, request customerReq.UpdateDescriptionRequest) error {
	// 查询wx_user_corp_external 的id
	userInfo, err := s.wxUserCorpExternalDao.FindOne(c, request.EnterpriseID, request.ExternalUserID)
	if err != nil {
		return err
	}
	// 修改备注
	err = s.workExternalDao.UpdateData(c, map[string]interface{}{
		"enterprise_id":         request.EnterpriseID,
		"member_id":             request.MemberID,
		"work_external_user_id": userInfo.ID,
	}, map[string]interface{}{
		"description": request.Description,
	})
	if err != nil {
		return err
	}
	relationInfo, err := s.workExternalDao.GetOne(c, map[string]interface{}{
		"enterprise_id":         request.EnterpriseID,
		"member_id":             request.MemberID,
		"work_external_user_id": userInfo.ID,
	})
	if err != nil {
		return err
	}
	// 同步备注到企微
	enterpriseInfo, _ := microBossDao.NewEnterprise().GetOne(c, request.EnterpriseID)
	memberUserInfo, _ := membersDao.NewMembers().GetOne(c, fmt.Sprintf("enterprise_id = %d and id = %d", request.EnterpriseID, request.MemberID))

	qyWork := new(work.QyWorkImpl).SetDomain(enterpriseInfo.Domain)
	updateData := work.EditExternalDescriptionRequest{
		UserID:         memberUserInfo.Userid,
		ExternalUserid: request.ExternalUserID,
		Remark:         relationInfo.Remark,
		Description:    request.Description,
	}
	err = qyWork.UpdateExternalRelationDescription(c, updateData)

	return err

}
func (s *Customer) GetRemarkByUnionIdAndMemberId(c app.Ctx, enterpriseId int64, unionId string, memberId int64) (string, error) {
	// 查询wx_user_corp_external 的id
	return s.workExternalDao.GetRemarkByUnionIdAndMemberId(c, enterpriseId, unionId, memberId)
}
