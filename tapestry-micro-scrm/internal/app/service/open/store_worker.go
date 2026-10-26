package openSvc

import (
	"errors"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	workExternalDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	openDao "micro-scrm/internal/app/dao/open"
	openReq "micro-scrm/internal/app/http/request/open"
	openResp "micro-scrm/internal/app/http/response/open"
	"micro-scrm/internal/app/models"
	membersModel "micro-scrm/internal/app/models/members"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	openModel "micro-scrm/internal/app/models/open"
	"micro-scrm/internal/app/service"
	cacheSvc "micro-scrm/internal/app/service/cache"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/work"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type StoreWorkerService struct {
	service.Base
}

func NewStoreWorkerService() *StoreWorkerService {
	return &StoreWorkerService{}
}

// SaListByCustomer 获取顾客的导购
func (s StoreWorkerService) SaListByCustomer(c app.Ctx, params openReq.StoreWorkerSaListReq) ([]openResp.SaListByCustomerList, error) {
	var resp []openResp.SaListByCustomerList
	externalUserRelationList, err := workExternalDao.NewUserRelation().GetAllByUnionId(c, params.EnterpriseId, params.Unionid, []int{0})
	if err != nil {
		return nil, err
	}

	var memberIds []int64
	for _, externalUserRelation := range externalUserRelationList {
		memberIds = append(memberIds, externalUserRelation.MemberID)
	}

	memberList, err := membersDao.NewMembers().GetAllById(c, params.EnterpriseId, memberIds)
	if err != nil {
		return nil, err
	}

	for _, member := range memberList {
		for _, externalUserRelation := range externalUserRelationList {
			if member.ID == externalUserRelation.MemberID {
				resp = append(resp, openResp.SaListByCustomerList{
					SaCode:    member.Userid,
					Name:      member.Name,
					CreatedAt: int64(externalUserRelation.Createtime),
				})
			}
		}
	}

	return resp, nil
}

// GetSaList 获取顾客的导购
func (s StoreWorkerService) GetSaList(c app.Ctx, params openReq.SaListReq) ([]openResp.SaList, int64, error) {
	var resp []openResp.SaList
	employeeList, total, err := microWorkDao.NewStoreEmployee().GetPagination(c, params)
	if total == 0 || err != nil {
		return resp, total, nil
	}
	for _, employee := range employeeList {
		resp = append(resp, openResp.SaList{
			SaCode:    employee.EmployeeCode,
			StoreCode: employee.ThirdStoreCode,
			Name:      employee.Name,
		})
	}
	return resp, total, nil
}

// GetStoreList 获取门店列表
func (s StoreWorkerService) GetStoreList(c app.Ctx, params openReq.StoreListReq) ([]openResp.StoreList, int64, error) {
	var resp []openResp.StoreList
	departmentList, total, err := microWorkDao.NewStoreDepartment().GetPagination(c, params)
	if total == 0 || err != nil {
		return resp, total, nil
	}
	for _, department := range departmentList {
		resp = append(resp, openResp.StoreList{
			StoreName:          department.StoreName,
			StoreCode:          department.ThirdStoreCode,
			InventoryName:      department.InventoryName,
			Address:            department.Address,
			City:               department.City,
			Phone:              department.Phone,
			Email:              department.Email,
			StoreType:          department.StoreType,
			ClosedAt:           department.ClosedAt,
			WechatDepartmentID: department.WechatDepartmentID,
		})
	}

	return resp, total, err
}

// GetStoreWorkList 获取导购信息列表
func (s StoreWorkerService) GetStoreWorkList(c app.Ctx, params openReq.StoreWorkListReq) ([]openResp.StoreWorkList, int64, error) {
	cacheRes, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, params.EnterpriseId, "store_name")
	if cacheRes.Name == "sw" {
		params.UserId = s.getUserIdBySwSaCode(c, params.EnterpriseId, params.UserId)
	}
	// 获取员工和部门集合
	memberList, total, err := membersDao.NewMembers().GetMemberAndDepartmentList(c, params)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	var resp []openResp.StoreWorkList
	var thirdIds []int

	for _, member := range memberList {
		tmpResp := openResp.StoreWorkList{
			MembersUserid:  member.Userid,
			MembersName:    member.Name,
			Mobile:         member.Mobile,
			Position:       member.Position,
			IsStoreManager: helper.RuleThree(helper.InArray(member.Position, membersDao.StoreManagerPosition), 1, 0),
			QrCode:         member.QrCode,
			Status:         member.Status,
		}
		if member.DeletedAt > 0 {
			tmpResp.Status = membersModel.MembersStatusLeave
		}
		var departments []openResp.StoreWorkListDepartments
		for _, department := range member.MemberDepartments {
			if department.DeletedAt == 0 {
				departments = append(departments, openResp.StoreWorkListDepartments{
					Name:    department.Department.Name,
					ThirdId: department.Department.ThirdID,
				})
				thirdIds = append(thirdIds, department.Department.ThirdID)
			}
		}
		tmpResp.Departments = departments
		if cacheRes.Name == "sw" && member.OpenUserid != "" {
			tmpResp.MembersUserid = member.OpenUserid
		}
		resp = append(resp, tmpResp)
	}
	// 查询门店编码
	thirdIdToStoreCode, err := s.GetStoreCodeByThirdIds(c, params.EnterpriseId, thirdIds)
	if err != nil {
		return nil, 0, err
	}
	for x, info := range resp {
		for y, department := range info.Departments {
			if _, ok := thirdIdToStoreCode[department.ThirdId]; ok {
				resp[x].Departments[y].StoreCode = thirdIdToStoreCode[department.ThirdId].ThirdStoreCode
				resp[x].Departments[y].StoreName = thirdIdToStoreCode[department.ThirdId].StoreName
			}
		}
	}
	return resp, total, nil
}

// GetCustomerStoreWorkList 获取客户添加的导购列表
func (s StoreWorkerService) GetCustomerStoreWorkList(c app.Ctx, params openReq.CustomerStoreWorkListReq) (interface{}, int64, error) {
	relationList, total, err := workExternalDao.NewUserRelation().GetCustomStoreWorkerByUnionId(c, params)
	if err != nil || total == 0 {
		return nil, 0, err
	}
	cacheRes, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(c, params.EnterpriseId, "store_name")
	var swSaCodeMap = make(map[string]string)
	if cacheRes.Name == "sw" {
		memberList, _ := membersDao.NewMembers().GetMany(c, "enterprise_id = ?", params.EnterpriseId)
		for _, model := range memberList {
			if model.OpenUserid != "" {
				swSaCodeMap[model.Userid] = model.OpenUserid
			}
		}
	}
	var resp []openResp.CustomerStoreWorkListResp
	var thirdIds []int
	for _, relation := range relationList {
		tmpResp := openResp.CustomerStoreWorkListResp{
			Id:           relation.ID,
			MemberUserId: relation.UserID,
			MemberName:   relation.Members.Name,
			Remark:       relation.Remark,
			CreateTime:   relation.Createtime,
			// IsKeySa:      relation.IsKeySa,
		}
		if cacheRes.Name == "sw" {
			if _, ok := swSaCodeMap[relation.UserID]; ok {
				tmpResp.MemberUserId = swSaCodeMap[relation.UserID]
			}
		}
		var departments []openResp.StoreWorkListDepartments
		for _, department := range relation.Members.MemberDepartments {
			departments = append(departments, openResp.StoreWorkListDepartments{
				Name:    department.Department.Name,
				ThirdId: department.Department.ThirdID,
			})
			thirdIds = append(thirdIds, department.Department.ThirdID)
		}
		tmpResp.Departments = departments
		resp = append(resp, tmpResp)
	}
	// 查询门店编码
	thirdIdToStoreCode, err := s.GetStoreCodeByThirdIds(c, params.EnterpriseId, thirdIds)
	if err != nil {
		return nil, 0, err
	}
	for x, info := range resp {
		for y, department := range info.Departments {
			if _, ok := thirdIdToStoreCode[department.ThirdId]; ok {
				resp[x].Departments[y].StoreCode = thirdIdToStoreCode[department.ThirdId].ThirdStoreCode
			}
		}
	}
	return resp, total, nil
}

func (s StoreWorkerService) GenQrCode(c app.Ctx, req openReq.GenQrCodeReq) (openResp.GenQrCodeResp, error) {
	cacheRes, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, req.EnterpriseId, "store_name")
	if cacheRes.Name == "sw" {
		req.UserId = s.getUserIdBySwSaCode(c, req.EnterpriseId, req.UserId)
	}
	var resp openResp.GenQrCodeResp
	// 判断人员是否存在
	memberInfo, err := membersDao.NewMembers().GetMemberInfo(ctx.Wrap(c),
		map[string]interface{}{
			"enterprise_id": req.EnterpriseId,
			"userid":        req.UserId,
		})
	if err != nil || memberInfo.ID == 0 {
		return resp, errors.New("获取员工信息异常")
	}
	resp = openResp.GenQrCodeResp{
		StoreWorkList: openResp.StoreWorkList{
			MembersUserid: memberInfo.Userid,
			MembersName:   memberInfo.Name,
			Mobile:        memberInfo.Mobile,
			Position:      memberInfo.Position,
			Status:        memberInfo.Status,
		},
	}
	//查询缓存
	wayModel, _ := openDao.NewOpenQrcodeContactWayDao().GetOne(c, req.EnterpriseId, "user_id = ? AND state = ?", req.UserId, req.State)
	if wayModel.ID > 0 {
		resp.StoreWorkList.QrCode = wayModel.QrCode
	} else {
		request := work.ContactCreateRequest{
			CodeCame:     "三方商城二维码",
			Type:         1, // 1-单人, 2-多人
			AutoFriend:   1,
			UserIds:      []string{req.UserId},
			EnterpriseId: req.EnterpriseId,
			State:        req.State,
		}

		workImpl := work.NewQyWorkImpl()
		contact, err := workImpl.SetEnterprise(req.EnterpriseId).GetWxEexternalContact(ctx.Wrap(c), request)
		if err != nil {
			return resp, errors.New("生成二维码异常")
		}
		resp.StoreWorkList.QrCode = contact.QrCode
		//入库
		_ = openDao.NewOpenQrcodeContactWayDao().CreateData(c, req.EnterpriseId, openModel.OpenQrcodeContactWayModel{
			BaseModel: models.CreateBase(app.UID()),
			State:     req.State,
			UserId:    req.UserId,
			ConfigId:  contact.ConfigId,
			QrCode:    contact.QrCode,
		})
	}
	var departments []openResp.StoreWorkListDepartments
	var thirdIds []int
	for _, memberDepartment := range memberInfo.MemberDepartments {
		departments = append(departments, openResp.StoreWorkListDepartments{
			Name:    memberDepartment.Department.Name,
			ThirdId: memberDepartment.Department.ThirdID,
		})
		thirdIds = append(thirdIds, memberDepartment.Department.ThirdID)
	}
	resp.StoreWorkList.Departments = departments

	// 查询门店编码
	thirdIdToStoreCode, err := s.GetStoreCodeByThirdIds(c, req.EnterpriseId, thirdIds)
	if err != nil {
		return openResp.GenQrCodeResp{}, err
	}
	for key, department := range resp.StoreWorkList.Departments {
		if _, ok := thirdIdToStoreCode[department.ThirdId]; ok {
			resp.StoreWorkList.Departments[key].StoreCode = thirdIdToStoreCode[department.ThirdId].ThirdStoreCode
		}
	}
	if cacheRes.Name == "sw" && memberInfo.OpenUserid != "" {
		resp.StoreWorkList.MembersUserid = memberInfo.OpenUserid
	}
	return resp, nil
}

func (s StoreWorkerService) GetStoreCodeByThirdIds(c app.Ctx, enterpriseId int64, thirdIds []int) (map[int]microWorkModel.StoreDepartmentModel, error) {
	// 查询门店编码
	storeList, err := microWorkDao.NewStoreDepartment().GetAllByWechatDepartmentIds(c, enterpriseId, thirdIds)
	if err != nil {
		return nil, err
	}
	var thirdIdToStoreCode = make(map[int]microWorkModel.StoreDepartmentModel)
	for _, store := range storeList {
		thirdIdToStoreCode[int(store.WechatDepartmentID)] = store
	}
	return thirdIdToStoreCode, nil
}

func (s StoreWorkerService) getUserIdBySwSaCode(c app.Ctx, enterpriseId int64, saCode string) string {
	if saCode != "" {
		memberInfo, _ := membersDao.NewMembers().GetOne(c, map[string]interface{}{
			"enterprise_id": enterpriseId,
			"open_userid":   saCode,
		})
		if memberInfo.Userid != "" {
			return memberInfo.Userid
		}
	}
	return saCode
}
