package membersSvc

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	memberReq "micro-scrm/internal/app/http/request/member"
	notificationResp "micro-scrm/internal/app/http/response/notification"
	"micro-scrm/internal/app/models"
	departmentsModels "micro-scrm/internal/app/models/departments"
	membersModel "micro-scrm/internal/app/models/members"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"net/url"
	"strings"
)

type MembersService struct {
	service.Base
	MembersDao            *membersDao.Members
	MembersDepartmentsDao *departmentsDao.MemberDepartments
}

func NewMembersService() *MembersService {
	return &MembersService{
		MembersDao:            membersDao.NewMembers(),
		MembersDepartmentsDao: departmentsDao.NewMemberDepartments(),
	}
}

// GetDetail 获取门店信息内容，门店的地址
func (s *MembersService) GetDetail(c app.Ctx, enterpriseId int64, mobile string) (memberReq.MemberDetailRes, error) {
	var department membersModel.MemberDepartment
	department, err := s.MembersDao.GetMemberInfoAndDepartmentsByMobile(c, enterpriseId, mobile)
	if err != nil {
		return memberReq.MemberDetailRes{}, err
	}
	storeDepartmentInfo, err := microWorkDao.NewStoreDepartment().GetInfo(c, microWorkModel.StoreDepartmentModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
		WechatDepartmentID: uint(department.DepartmentThirdId),
	})
	if err != nil {
		return memberReq.MemberDetailRes{MemberDepartment: department}, err
	}
	if helper.InArray(enterpriseId, []int64{3373738709328576, 3374939108972352}) { // 此处临时代码。判断企微为SW，SW的门店名称取store_department
		if storeDepartmentInfo.StoreName != "" {
			department.DepartmentName = storeDepartmentInfo.StoreName
		}
	}
	result := memberReq.MemberDetailRes{
		MemberDepartment: department,
		Address:          storeDepartmentInfo.Address,
		StoreCode:        storeDepartmentInfo.ThirdStoreCode,
		StoreName:        storeDepartmentInfo.StoreName,
	}
	return result, err
}

func (s *MembersService) GetInfoById(c app.Ctx, enterpriseId, id int64) (membersModel.MemberModel, error) {
	return s.MembersDao.GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel: models.BaseModel{
				ID: id,
			},
			EnterpriseID: enterpriseId,
		},
	})
}

func (s *MembersService) GetMemberByDepartmentId(c app.Ctx, enterpriseId int64, departments []int64) ([]membersModel.MemberDepartment, error) {
	return s.MembersDao.GetMemberByDepartmentId(c, enterpriseId, departments)
}

func (s *MembersService) GetInfoByMobile(c app.Ctx, enterpriseId int64, mobile string) (membersModel.MemberModel, error) {
	return s.MembersDao.GetMemberInfo(c, membersModel.MemberModel{
		Mobile: mobile,
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
	})
}
func (s *MembersService) GetThirdIdsInfoByUserId(c app.Ctx, enterpriseId int64, userId string) ([]int, error) {
	memberInfo, err := s.MembersDao.GetMemberInfo(c, membersModel.MemberModel{
		Userid: userId,
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
	})
	if err != nil {
		return nil, err
	}
	var thirdIds []int
	if memberInfo.ID > 0 && memberInfo.MemberDepartments != nil && len(memberInfo.MemberDepartments) > 0 {
		for _, relation := range memberInfo.MemberDepartments {
			if relation.DeletedAt == 0 && relation.Department.DeletedAt == 0 {
				thirdIds = append(thirdIds, relation.Department.ThirdID)
			}
		}
	}
	return thirdIds, nil
}

// GetStoreLists 获取门店列表
func (s *MembersService) GetStoreLists(c app.Ctx, request memberReq.MemberStoreReq) ([]memberReq.MemberStoreRes, error) {
	var result []memberReq.MemberStoreRes
	thirdIds := request.ThirdIds
	enterpriseId := request.EnterpriseId
	if request.IsSm == 0 {
		thirdIds = []int{}
		departments, err := s.MembersDao.GetMemberInfoAndAllDepartmentsByMobile(c, enterpriseId, request.Mobile)
		if err != nil {
			return nil, err
		}
		for _, v := range departments {
			var storeRes memberReq.MemberStoreRes
			storeRes.DepartmentId = v.DepartmentId
			storeRes.StoreId = int(v.DepartmentThirdId)
			storeRes.StoreName = v.DepartmentName
			result = append(result, storeRes)
		}
		return result, err
	} else {
		lists, err := departmentsDao.NewDepartments().GetDepartmentsInfoThirdIds(c, enterpriseId, thirdIds)
		if err != nil {
			return nil, err
		}
		for _, list := range lists {
			var storeRes memberReq.MemberStoreRes
			storeRes.DepartmentId = list.ID
			storeRes.StoreId = list.ThirdID
			storeRes.StoreName = list.Name
			result = append(result, storeRes)
		}
		return result, err
	}

}

func (s *MembersService) GetStaffLists(c app.Ctx, request memberReq.MemberStaffReq) ([]memberReq.MemberStaffRes, error) {
	var memberIds []int64
	var err error
	if request.UnionId != "" {
		memberIds, err = microScrmDao.NewUserRelation().GetMemberIdsByUnionId(c, request.EnterpriseId, request.UnionId)
		if err != nil {
			return []memberReq.MemberStaffRes{}, err
		}
	}
	info, err := s.MembersDao.GetMemberInfoAndDepartmentsByMobile(c, request.EnterpriseId, request.Mobile)
	if err != nil {
		return nil, err
	}
	lists, err := s.MembersDao.GetMemberByThirdIds(c, request.EnterpriseId, request.StoreIds, memberIds)
	if len(request.StoreIds) > 0 { // 店长所在门店不在管理的门店里，导购列表需要展示店长自己
		if !helper.InArray(int(info.DepartmentThirdId), request.StoreIds) {
			item := memberReq.MemberStaffRes{
				StoreId:       int(info.DepartmentThirdId),
				StoreName:     info.DepartmentName,
				MembersId:     info.MembersId,
				MembersAvatar: "",
				MembersName:   info.MembersName,
				MembersUserid: info.MembersUserid,
				Mobile:        info.Mobile,
			}
			lists = append(lists, item)
		}
	}
	return lists, err
}

// GetDetailByMemberId 获取门店信息内容，门店的地址
func (s *MembersService) GetDetailByMemberId(c app.Ctx, enterpriseId, memberId int64) (memberReq.MemberDetailRes, error) {
	var department membersModel.MemberDepartment
	department, err := s.MembersDao.GetMemberInfoAndDepartmentsByMemberId(c, enterpriseId, memberId)
	if err != nil {
		return memberReq.MemberDetailRes{}, err
	}
	storeDepartmentInfo, err := microWorkDao.NewStoreDepartment().GetInfo(c, microWorkModel.StoreDepartmentModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
		WechatDepartmentID: uint(department.DepartmentThirdId),
	})
	if err != nil {
		return memberReq.MemberDetailRes{}, err
	}
	department.DepartmentName = storeDepartmentInfo.StoreName
	result := memberReq.MemberDetailRes{
		MemberDepartment: department,
		Address:          storeDepartmentInfo.Address,
	}
	return result, err
}

func (s *MembersService) GetMemberAndDepartmentsByMemberId(c app.Ctx, enterpriseId int64, memberId int64) ([]membersModel.MemberDepartment, error) {
	return s.MembersDao.GetMemberAndDepartmentsByMemberId(c, enterpriseId, memberId)
}

func (s *MembersService) GetInfoByUserId(c app.Ctx, enterpriseId int64, userid string) (membersModel.MemberModel, error) {
	return s.MembersDao.GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
		Userid: userid,
	})
}
func (s *MembersService) GetInfoByOpenUserId(c app.Ctx, enterpriseId int64, openUserid string) (membersModel.MemberModel, error) {
	return s.MembersDao.GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
		OpenUserid: openUserid,
	})
}

func (s *MembersService) GetMemberInfoAndAllDepartmentsByMobile(c app.Ctx, enterpriseId int64, mobile string) ([]membersModel.MemberDepartment, error) {
	return s.MembersDao.GetMemberInfoAndAllDepartmentsByMobile(c, enterpriseId, mobile)
}

func (s *MembersService) GetMemberIdsByName(c app.Ctx, enterpriseId int64, name string) (memberIds []int64, err error) {

	return s.MembersDao.GetMemberIdsByName(c, enterpriseId, name)
}

func (s *MembersService) GetMemberIdsByThirdIds(c app.Ctx, enterpriseId int64, thirdIds []int) (memberIds []int64, err error) {
	return s.MembersDepartmentsDao.GetMemberIdsByThirdIds(c, enterpriseId, thirdIds)
}

func (s *MembersService) GetMemberIdsByUserid(c app.Ctx, enterpriseId int64, userid string) (memberIds []int64, err error) {
	return s.MembersDao.GetMemberIdsByUserid(c, enterpriseId, userid)

}

func (s *MembersService) GetThirdsByMemberIds(c app.Ctx, enterpriseId, memberId int64) (thirdIds []int, err error) {
	return s.MembersDepartmentsDao.GetThirdsByMemberIds(c, enterpriseId, memberId)
}

// SetMiniProgramLinkParamsByTask sw小程序埋点给小程序path增加参数(群发任务用）
func (s *MembersService) SetMiniProgramLinkParamsByTask(c app.Ctx, enterpriseId int64, userIds []string, rawPath string) (map[string]string, error) {
	var resMap = make(map[string]string)
	memberDepartmentList, err := s.MembersDao.GetMemberInfoAndDepartmentByUserIds(c, enterpriseId, userIds)
	var userIdThirdIdMap = make(map[string]int)
	var thirdIds []int
	for _, info := range memberDepartmentList {
		var thirdId int
		if len(info.MemberDepartments) > 0 && info.MemberDepartments[0].Department.ID > 0 {
			thirdId = info.MemberDepartments[0].Department.ThirdID
			thirdIds = append(thirdIds, thirdId)
		}
		userIdThirdIdMap[info.Userid] = thirdId
	}
	storeList, err := microWorkDao.NewStoreDepartment().GetAllByWechatDepartmentIds(c, enterpriseId, thirdIds)
	if err != nil {
		return resMap, err
	}
	var storeCodeMap = make(map[int]string)
	for _, storeInfo := range storeList {
		storeCodeMap[int(storeInfo.WechatDepartmentID)] = storeInfo.ThirdStoreCode
	}
	// 解析路径信息
	pathAndQuery := strings.Split(rawPath, "?")
	path := pathAndQuery[0]

	for _, userId := range userIds {
		storeCode := "unknown" // 兼容找不到门店编码的情况
		if _, ok := storeCodeMap[userIdThirdIdMap[userId]]; ok {
			storeCode = storeCodeMap[userIdThirdIdMap[userId]]
		}
		newQueryString := s.setMiniPathToSw(storeCode, userId, pathAndQuery)
		// 拼接新的路径和查询参数
		resMap[userId] = path + "?" + newQueryString
		fmt.Println("新路径:", resMap[userId])
	}
	return resMap, nil
}

// SetMiniProgramLinkParams sw小程序埋点给小程序path增加参数(单导购用）
func (s *MembersService) SetMiniProgramLinkParams(c app.Ctx, enterpriseId int64, memberId int64, pathMap map[int64]string) (map[int64]string, error) {
	var newPathMap = make(map[int64]string)
	memberInfo, err := s.MembersDao.GetMemberInfo(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"id":            memberId,
	})
	var thirdId int
	if len(memberInfo.MemberDepartments) > 0 {
		thirdId = memberInfo.MemberDepartments[0].Department.ThirdID
	} else {
		return newPathMap, errors.New("找不到该导购的门店")
	}

	storeInfo, err := microWorkDao.NewStoreDepartment().GetInfo(c, microWorkModel.StoreDepartmentModel{
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: enterpriseId,
		},
		WechatDepartmentID: uint(thirdId),
	})
	if err != nil {
		return newPathMap, err
	}

	for id, rawPath := range pathMap {
		// 解析路径信息
		pathAndQuery := strings.Split(rawPath, "?")
		path := pathAndQuery[0]
		newQueryString := s.setMiniPathToSw(storeInfo.ThirdStoreCode, memberInfo.Userid, pathAndQuery)
		// 拼接新的路径和查询参数
		newPathMap[id] = path + "?" + newQueryString
	}
	return newPathMap, err
}

// 拼接SW需要的小程序参数
func (s *MembersService) setMiniPathToSw(storeCode, userId string, pathAndQuery []string) string {
	valueStr := fmt.Sprintf("wc_%s_%s", storeCode, userId)
	// 组装新参数
	queryParams := make(url.Values)
	if len(pathAndQuery) > 1 {
		queryString := pathAndQuery[1]
		queryParams, _ = url.ParseQuery(queryString)
	}
	if _, ok := queryParams["utm_medium"]; ok {
		queryParams.Set("utm_medium", valueStr)
	} else {
		queryParams.Add("utm_medium", valueStr)
	}
	// 构建新的查询字符串
	return queryParams.Encode()
}

// GetStoreByDepartments 从导购的多个部门中找到真实门店id
func (s *MembersService) GetStoreByDepartments(memberDepartments []membersModel.MemberDepartmentsModel, rootDepartmentThirdId int) departmentsModels.DepartmentsModel {
	for _, department := range memberDepartments {
		if department.Department.ThirdParentID == rootDepartmentThirdId {
			// 如果存在真实门店则返回真实门店id
			return department.Department
		}
	}
	// 如果没有真实门店则返回首个门店id
	return memberDepartments[0].Department
}

// GetFirstOpen 获取首次进入某页面是否需要引导功能
func (s *MembersService) GetFirstOpen(c app.Ctx, req memberReq.FirstOpenReq) (int, error) {
	one, err := membersDao.NewFirstOpen().GetOne(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
		"mobile":        req.Mobile,
		"type":          req.Type,
	})
	if err != nil {
		return 0, err
	}
	if one.ID > 0 {
		return 0, nil
	}
	_ = membersDao.NewFirstOpen().Create(c, membersModel.FirstOpenModel{
		MicroBaseModel: models.CreateMicroBase(req.EnterpriseId, app.UID()),
		MarkType:       req.Type,
		Mobile:         req.Mobile,
	})

	return membersModel.IsFirstOpen, nil
}

func (s *MembersService) GetMemberListByStoreCodes(c ctx.Context, params memberReq.StoreCodesReq) (notificationResp.ImportDataResp, error) {
	list, err := microWorkDao.NewStoreDepartment().GetStoreInfoByStoreCode(c, params.EnterpriseId, params.StoreCodes)
	if err != nil {
		return notificationResp.ImportDataResp{}, err
	}
	var thirdIds []int64
	var storeCodes []string
	for _, store := range list {
		if store.WechatDepartmentID > 0 {
			thirdIds = append(thirdIds, int64(store.WechatDepartmentID))
			storeCodes = append(storeCodes, store.ThirdStoreCode)
		}
	}
	dpIds, err := departmentsDao.NewDepartments().GetDepartmentIdsByThirdIds(c, params.EnterpriseId, thirdIds)
	if err != nil {
		return notificationResp.ImportDataResp{}, err
	}
	memberData, err := membersDao.NewMemberDepartments().GetMembersByDepartmentIds(c, params.EnterpriseId, dpIds)
	if err != nil {
		return notificationResp.ImportDataResp{}, err
	}
	var memberList []notificationResp.CrowdDataResp
	for _, info := range memberData {
		if info.Member.ID > 0 {
			memberList = append(memberList, notificationResp.CrowdDataResp{
				MemberId:   info.MemberID,
				MemberName: info.Member.Name,
			})
		}
	}
	return notificationResp.ImportDataResp{MemberList: memberList, StoreCodes: storeCodes}, nil
}
