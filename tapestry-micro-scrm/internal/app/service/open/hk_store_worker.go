package openSvc

import (
	"micro-scrm/internal/app"
	hkStoreDao "micro-scrm/internal/app/dao/hk_store"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/http/request"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	openReq "micro-scrm/internal/app/http/request/open"
	openResp "micro-scrm/internal/app/http/response/open"
	hkStoreModel "micro-scrm/internal/app/models/hk_store"
	membersModel "micro-scrm/internal/app/models/members"
	"micro-scrm/internal/app/service"
	cacheSvc "micro-scrm/internal/app/service/cache"
	"micro-scrm/internal/pkg/helper"
)

type HkStoreWorkerService struct {
	service.Base
	storeDao   *hkStoreDao.StoreDepartment
	managerDao *hkStoreDao.StoreManager
}

func NewHkStoreWorkerService() *HkStoreWorkerService {
	return &HkStoreWorkerService{
		storeDao:   hkStoreDao.NewStoreDepartment(),
		managerDao: hkStoreDao.NewStoreManager(),
	}
}

func (s HkStoreWorkerService) GetStoreList(c app.Ctx, params openReq.HkStoreListReq) ([]openResp.HkStoreList, int64, error) {
	var resp []openResp.HkStoreList
	storeList, total, err := s.storeDao.GetPagination(c, params)
	if total == 0 || err != nil {
		return resp, total, err
	}
	for _, store := range storeList {
		resp = append(resp, openResp.HkStoreList{
			StoreName:          store.StoreName,
			StoreCode:          store.StoreCode,
			Area:               store.Area,
			InventoryName:      store.InventoryName,
			Address:            store.Address,
			City:               store.City,
			Phone:              store.Phone,
			Email:              store.Email,
			ClosedAt:           store.ClosedAt,
			WechatDepartmentID: store.WechatDepartmentID,
		})
	}
	return resp, total, nil
}

func (s HkStoreWorkerService) GetStoreWorkList(c app.Ctx, params openReq.HkStoreWorkListReq) ([]openResp.HkStoreWorkList, int64, error) {
	cacheRes, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(c, params.EnterpriseId, "store_name")
	if cacheRes.Name == "sw" {
		params.UserId = s.getUserIdBySwSaCode(c, params.EnterpriseId, params.UserId)
	}
	departmentIds, err := s.storeDao.GetWechatDepartmentIdsByWorkListReq(c, params)
	if err != nil || len(departmentIds) == 0 {
		return nil, 0, err
	}
	storeMap, err := s.storeDao.GetMapByWechatDepartmentIds(c, params.EnterpriseId, departmentIds)
	if err != nil {
		return nil, 0, err
	}
	memberList, total, err := membersDao.NewMembers().GetHkStoreMemberAndDepartmentList(c, params, departmentIds)
	if err != nil || total == 0 {
		return nil, total, err
	}

	var resp []openResp.HkStoreWorkList
	for _, member := range memberList {
		tmpResp := openResp.HkStoreWorkList{
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
		for _, department := range member.MemberDepartments {
			store, ok := storeMap[department.Department.ThirdID]
			if !ok {
				continue
			}
			tmpResp.Departments = append(tmpResp.Departments, openResp.HkStoreWorkListDepartment{
				Name:      department.Department.Name,
				ThirdId:   department.Department.ThirdID,
				StoreCode: store.StoreCode,
				StoreName: store.StoreName,
			})
		}
		if len(tmpResp.Departments) == 0 {
			continue
		}
		if cacheRes.Name == "sw" && member.OpenUserid != "" {
			tmpResp.MembersUserid = member.OpenUserid
		}
		resp = append(resp, tmpResp)
	}
	return resp, total, nil
}

func (s HkStoreWorkerService) AreaManager(c app.Ctx, params openReq.HkAreaManagerReq) (openResp.HkManagerInfo, error) {
	memberInfo, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"userid":        params.UserId,
	})
	if err != nil {
		return openResp.HkManagerInfo{}, err
	}
	info := openResp.HkManagerInfo{
		UserId:   memberInfo.Userid,
		UserName: memberInfo.Name,
	}
	manager, err := s.managerDao.GetOne(c, params.EnterpriseId, map[string]interface{}{
		"member_id": memberInfo.ID,
	})
	if err != nil {
		return info, err
	}
	info.Position = manager.Position
	var areaIds []int64
	for _, areaIdStr := range helper.StringToIntSlice(manager.AreaIds, ",") {
		areaIds = append(areaIds, int64(areaIdStr))
	}
	departmentIds, err := hkStoreDao.NewArea().GetThirdIdsByAreaIds(c, openReqToHkManagerAuthAreaIds(params.EnterpriseId, areaIds))
	if err != nil {
		return info, err
	}
	storeMap, err := s.storeDao.GetMapByWechatDepartmentIds(c, params.EnterpriseId, departmentIds)
	if err != nil {
		return info, err
	}
	for _, departmentId := range departmentIds {
		if store, ok := storeMap[departmentId]; ok {
			info.StoreList = append(info.StoreList, s.formatStoreInfo(store))
		}
	}
	return info, nil
}

func (s HkStoreWorkerService) formatStoreInfo(store hkStoreModel.StoreDepartmentModel) openResp.HkStoreInfo {
	return openResp.HkStoreInfo{
		ThirdId:   int(store.WechatDepartmentID),
		StoreCode: store.StoreCode,
		StoreName: store.StoreName,
		Area:      store.Area,
		EnName:    store.InventoryName,
	}
}

func openReqToHkManagerAuthAreaIds(enterpriseId int64, areaIds []int64) dashboardReq.HkManagerAuthAreaIds {
	return dashboardReq.HkManagerAuthAreaIds{
		EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: enterpriseId},
		AreaIds:           areaIds,
	}
}

func (s HkStoreWorkerService) getUserIdBySwSaCode(c app.Ctx, enterpriseId int64, saCode string) string {
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
