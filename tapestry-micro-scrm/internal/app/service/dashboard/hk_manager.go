package dashboardSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	hkStoreDao "micro-scrm/internal/app/dao/hk_store"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/http/request"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	dashboardResp "micro-scrm/internal/app/http/response/dashboard"
	"micro-scrm/internal/app/models"
	hkStoreModel "micro-scrm/internal/app/models/hk_store"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"sort"
	"strconv"
	"strings"
	"time"
)

type HkManagerService struct {
	service.Base
	managerDao *hkStoreDao.StoreManager
	storeDao   *hkStoreDao.StoreDepartment
	areaDao    *hkStoreDao.Area
}

func NewHkManagerService() *HkManagerService {
	return &HkManagerService{
		managerDao: hkStoreDao.NewStoreManager(),
		storeDao:   hkStoreDao.NewStoreDepartment(),
		areaDao:    hkStoreDao.NewArea(),
	}
}

func (s *HkManagerService) StoreList(c ctx.Context, params dashboardReq.HkManagerStoreList) ([]dashboardResp.HkStoreInfo, error) {
	storeList, err := s.storeDao.GetStoreList(c, params)
	if err != nil {
		return nil, err
	}
	return s.formatStoreList(storeList), nil
}

func (s *HkManagerService) SaveManagerAuth(c ctx.Context, params dashboardReq.HkManagerAuthSave) error {
	params.AreaIds = helper.ArrayUniq(params.AreaIds)
	areaList, err := s.areaDao.GetListByAreaIds(c, params.EnterpriseId, params.AreaIds)
	if err != nil {
		return err
	}
	if len(areaList) != len(params.AreaIds) {
		return errorx.NewErrorMsg(1001, "您选择的港澳区域不存在")
	}
	areaPathMap := make(map[int64][]int64)
	for _, area := range areaList {
		for _, areaIdStr := range strings.Split(area.Path, ",") {
			areaId := int64(helper.StringToInt(areaIdStr))
			if areaId > 0 {
				areaPathMap[area.ID] = append(areaPathMap[area.ID], areaId)
			}
		}
	}
	sort.Slice(areaList, func(i, j int) bool {
		return areaList[i].Type < areaList[j].Type
	})
	var areaIds []int64
	var areaIdsStr string
	for _, area := range areaList {
		if len(helper.ArrayIntersect(areaPathMap[area.ID], areaIds)) > 0 {
			continue
		}
		areaIds = append(areaIds, area.ID)
		areaIdsStr += strconv.FormatInt(area.ID, 10) + ","
	}
	memberInfo, _ := s.managerDao.GetOne(c, params.EnterpriseId, "member_id = ?", params.MemberId)
	if params.ID > 0 {
		info, _ := s.managerDao.GetOne(c, params.EnterpriseId, "id = ?", params.ID)
		if info.ID == 0 {
			return errorx.NewErrorMsg(1001, "您要更新的职位信息不存在")
		}
		if memberInfo.ID > 0 && memberInfo.ID != info.ID {
			return errorx.NewErrorMsg(1001, "您选择的导购已经存在职位")
		}
		return s.managerDao.UpdateData(c, params.EnterpriseId, map[string]interface{}{
			"id": info.ID,
		}, map[string]interface{}{
			"member_id":      params.MemberId,
			"area_ids":       strings.TrimRight(areaIdsStr, ","),
			"position":       params.Position,
			"position_level": params.PositionLevel,
			"updated_at":     time.Now().Unix(),
		})
	}
	if memberInfo.ID > 0 {
		return errorx.NewErrorMsg(1001, "您选择的导购已经存在职位")
	}
	unixTime := time.Now().Unix()
	return s.managerDao.CreateData(c, params.EnterpriseId, hkStoreModel.StoreManagerModel{
		BaseModel: models.BaseModel{
			CreatedAt: unixTime,
			UpdatedAt: unixTime,
		},
		MemberId:      params.MemberId,
		AreaIds:       strings.TrimRight(areaIdsStr, ","),
		Position:      params.Position,
		PositionLevel: params.PositionLevel,
	})
}

func (s *HkManagerService) ManagerAuthDetail(c ctx.Context, params dashboardReq.HkManagerAuthDetail) (detail dashboardResp.HkManagerAuthDetail, err error) {
	info, err := s.managerDao.GetOne(c, params.EnterpriseId, "id = ?", params.ID)
	if err != nil {
		return detail, err
	}
	memberInfo, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{"enterprise_id": params.EnterpriseId, "id": info.MemberId})
	if err != nil {
		return detail, err
	}
	for _, areaIdStr := range strings.Split(info.AreaIds, ",") {
		areaId := int64(helper.StringToInt(areaIdStr))
		if areaId > 0 {
			detail.AreaIds = append(detail.AreaIds, areaId)
		}
	}
	thirdIds, err := s.areaDao.GetThirdIdsByAreaIds(c, dashboardReq.HkManagerAuthAreaIds{
		EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: params.EnterpriseId},
		AreaIds:           detail.AreaIds,
	})
	if err != nil {
		return detail, err
	}
	storeMap, err := s.storeDao.GetActiveMapByWechatDepartmentIds(c, params.EnterpriseId, thirdIds)
	if err != nil {
		return detail, err
	}
	detail.ID = info.ID
	detail.Position = info.Position
	detail.StoreList = make([]dashboardResp.HkStoreInfo, 0)
	detail.MemberInfo = dashboardResp.HkManagerMemberInfo{
		ID:     memberInfo.ID,
		Name:   memberInfo.Name,
		Userid: memberInfo.Userid,
		Mobile: memberInfo.Mobile,
	}
	detail.StoreNum = len(storeMap)
	for _, thirdId := range thirdIds {
		if store, ok := storeMap[thirdId]; ok {
			detail.StoreList = append(detail.StoreList, s.formatStoreInfo(store))
		}
	}
	return
}

func (s *HkManagerService) ManagerAuthList(c ctx.Context, params dashboardReq.HkManagerAuthList) ([]dashboardResp.HkManagerList, int64, error) {
	list, total, err := s.managerDao.GetPagination(c, params)
	if err != nil {
		return nil, 0, err
	}
	allAreaList, err := s.areaDao.GetAll(c, params.EnterpriseId)
	if err != nil {
		return nil, 0, err
	}
	var storeAreaList []hkStoreModel.AreaModel
	areaMap := make(map[int64]hkStoreModel.AreaModel)
	for _, area := range allAreaList {
		if area.Type == hkStoreDao.AreaTypeStore {
			storeAreaList = append(storeAreaList, area)
		}
		areaMap[area.ID] = area
	}
	storeList, err := s.storeDao.GetMany(c, map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"deleted_at":    0,
		"closed_at":     0,
	})
	if err != nil {
		return nil, 0, err
	}
	storeMap := make(map[int]string)
	for _, store := range storeList {
		storeMap[int(store.WechatDepartmentID)] = store.StoreName
	}
	for i, manager := range list {
		areaIds := strings.Split(manager.AreaIds, ",")
		manager.StoreNameArr = make([]string, 0)
		for _, areaIdStr := range areaIds {
			areaId := int64(helper.StringToInt(areaIdStr))
			area, ok := areaMap[areaId]
			if !ok {
				continue
			}
			if area.Type == hkStoreDao.AreaTypeStore {
				if storeName, ok := storeMap[area.ThirdID]; ok {
					manager.PathNameArr = append(manager.PathNameArr, storeName)
					manager.StoreNameArr = append(manager.StoreNameArr, storeName)
				}
				continue
			}
			manager.PathNameArr = append(manager.PathNameArr, area.Name)
		}
		for _, storeArea := range storeAreaList {
			pathIds := strings.Split(storeArea.Path, ",")
			if len(helper.ArrayIntersect(append(pathIds, helper.Int64ToString(storeArea.ID)), areaIds)) > 0 {
				manager.StoreNum++
			}
		}
		list[i] = manager
	}
	return list, total, nil
}

func (s *HkManagerService) GetAllManagerMemberIds(c ctx.Context, enterpriseId int64) (memberIds []int64, err error) {
	list, err := s.managerDao.GetAll(c, enterpriseId)
	if err != nil {
		return memberIds, err
	}
	for _, manager := range list {
		memberIds = append(memberIds, manager.MemberId)
	}
	return
}

func (s *HkManagerService) DeleteManagerAuth(c ctx.Context, params dashboardReq.HkManagerRemove) error {
	return s.managerDao.DeleteData(c, params.EnterpriseId, "id=?", params.Id)
}

func (s *HkManagerService) formatStoreList(storeList []hkStoreModel.StoreDepartmentModel) []dashboardResp.HkStoreInfo {
	var resp []dashboardResp.HkStoreInfo
	for _, store := range storeList {
		resp = append(resp, s.formatStoreInfo(store))
	}
	return resp
}

func (s *HkManagerService) formatStoreInfo(store hkStoreModel.StoreDepartmentModel) dashboardResp.HkStoreInfo {
	return dashboardResp.HkStoreInfo{
		ID:                 store.ID,
		StoreCode:          store.StoreCode,
		StoreName:          store.StoreName,
		Area:               store.Area,
		StoreType:          store.StoreType,
		InventoryName:      store.InventoryName,
		WechatDepartmentID: store.WechatDepartmentID,
		ClosedAt:           store.ClosedAt,
	}
}
