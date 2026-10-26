package dashboardSvc

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	hkStoreDao "micro-scrm/internal/app/dao/hk_store"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	dashboardResp "micro-scrm/internal/app/http/response/dashboard"
	hkStoreModel "micro-scrm/internal/app/models/hk_store"
	"micro-scrm/internal/app/service"
)

type HkAreaService struct {
	service.Base
	areaDao  *hkStoreDao.Area
	storeDao *hkStoreDao.StoreDepartment
}

func NewHkAreaService() *HkAreaService {
	return &HkAreaService{
		areaDao:  hkStoreDao.NewArea(),
		storeDao: hkStoreDao.NewStoreDepartment(),
	}
}

func (s *HkAreaService) GetAreaList(c ctx.Context, params dashboardReq.AreaTree) ([]dashboardResp.AreaList, error) {
	areaList, err := s.areaDao.GetList(c, params)
	if err != nil {
		return nil, err
	}
	if params.HasChildren {
		areaList = buildHkAreaTree(areaList, params.ParentId)
	}
	storeList, err := s.storeDao.GetMany(c, map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"deleted_at":    0,
	})
	if err != nil {
		return nil, err
	}
	storeMap := make(map[int]hkStoreModel.StoreDepartmentModel)
	for _, store := range storeList {
		if store.WechatDepartmentID > 0 {
			storeMap[int(store.WechatDepartmentID)] = store
		}
	}
	result := make([]dashboardResp.AreaList, 0, len(areaList))
	for _, area := range areaList {
		result = append(result, s.formatArea(area, storeMap, ""))
	}
	return result, nil
}

func (s *HkAreaService) formatArea(area hkStoreModel.AreaModel, storeMap map[int]hkStoreModel.StoreDepartmentModel, pathName string) dashboardResp.AreaList {
	result := dashboardResp.AreaList{
		ID:        int(area.ID),
		Type:      area.Type,
		Name:      area.Name,
		ParentID:  area.ParentID,
		ThirdID:   area.ThirdID,
		PathName:  pathName,
		ClosedAt:  area.ClosedAt,
		UpdatedAt: area.UpdatedAt,
	}
	if store, ok := storeMap[area.ThirdID]; ok {
		result.Name = fmt.Sprintf("%s_%s", store.StoreName, area.Name)
		result.StoreCode = store.StoreCode
		result.StoreName = store.StoreName
	}
	for _, child := range area.Children {
		childPath := area.Name
		if pathName != "" {
			childPath = pathName + "/" + area.Name
		}
		result.Children = append(result.Children, s.formatArea(child, storeMap, childPath))
	}
	return result
}

func buildHkAreaTree(areaList []hkStoreModel.AreaModel, parentId int) []hkStoreModel.AreaModel {
	childrenMap := make(map[int][]hkStoreModel.AreaModel)
	for _, area := range areaList {
		childrenMap[area.ParentID] = append(childrenMap[area.ParentID], area)
	}
	var attachChildren func(area hkStoreModel.AreaModel) hkStoreModel.AreaModel
	attachChildren = func(area hkStoreModel.AreaModel) hkStoreModel.AreaModel {
		for _, child := range childrenMap[int(area.ID)] {
			area.Children = append(area.Children, attachChildren(child))
		}
		return area
	}
	result := make([]hkStoreModel.AreaModel, 0)
	for _, area := range childrenMap[parentId] {
		result = append(result, attachChildren(area))
	}
	return result
}
