package dashboardSvc

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	dashboardDao "micro-scrm/internal/app/dao/dashboard"
	membersDao "micro-scrm/internal/app/dao/members"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/http/request"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	openReq "micro-scrm/internal/app/http/request/open"
	dashboardResp "micro-scrm/internal/app/http/response/dashboard"
	openResp "micro-scrm/internal/app/http/response/open"
	"micro-scrm/internal/app/models"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
	membersModel "micro-scrm/internal/app/models/members"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"strings"
	"time"
)

type AreaService struct {
	service.Base
	areaDao    *dashboardDao.Area
	managerDao *dashboardDao.Manager
}

func NewAreaService() *AreaService {
	return &AreaService{
		areaDao:    dashboardDao.NewArea(),
		managerDao: dashboardDao.NewManager(),
	}
}

// NotRelationStoreList  待关联门店列表
func (s *AreaService) NotRelationStoreList(c ctx.Context, enterpriseId int64) (res []microWorkModel.StoreDepartmentModel, err error) {
	allDepartmentList, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"deleted_at":    0,
	})
	// 获取dashboard中配置了的企微门店集合
	storeList, err := s.areaDao.GetMany(c, enterpriseId, map[string]interface{}{
		"type":       dashboardModel.AreaStore,
		"is_virtual": 0,
		"closed_at":  0,
	})
	var storeThirdIds []int
	for _, store := range storeList {
		storeThirdIds = append(storeThirdIds, store.ThirdID)
	}
	for _, department := range allDepartmentList {
		if department.WechatDepartmentID > 0 {
			if !helper.InArray(int(department.WechatDepartmentID), storeThirdIds) {
				res = append(res, department)
			}
		}
	}
	return
}

func (s *AreaService) GetAreaList(c ctx.Context, params dashboardReq.AreaTree) (list []dashboardResp.AreaList, err error) {
	areaList, err := s.areaDao.GetList(c, params)
	if err != nil {
		return nil, err
	}
	allDepartmentList, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"deleted_at":    0,
	})
	if err != nil {
		return nil, err
	}
	var departmentMap = make(map[int]microWorkModel.StoreDepartmentModel)
	for _, department := range allDepartmentList {
		if params.Language == "en" && department.InventoryName != "" {
			department.StoreName = department.InventoryName
		}
		if department.WechatDepartmentID > 0 {
			departmentMap[int(department.WechatDepartmentID)] = department
		}
	}
	for _, info := range areaList {
		list = append(list, s.recursionAreaInfo(info, departmentMap, ""))
	}
	return
}

// 递归整理门店信息
func (s *AreaService) recursionAreaInfo(info dashboardModel.DashboardAreaModel, departmentMap map[int]microWorkModel.StoreDepartmentModel, pathName string) dashboardResp.AreaList {
	var tmp = dashboardResp.AreaList{
		ID:        int(info.ID),
		Type:      info.Type,
		Name:      info.Name,
		ParentID:  info.ParentID,
		ThirdID:   info.ThirdID,
		PathName:  pathName,
		ClosedAt:  info.ClosedAt,
		UpdatedAt: info.UpdatedAt,
	}
	if val, ok := departmentMap[info.ThirdID]; ok {
		tmp.Name = fmt.Sprintf("%s_%s", val.StoreName, info.Name)
		tmp.StoreCode = val.ThirdStoreCode
		tmp.StoreName = val.StoreName
	}
	if len(info.Children) > 0 {
		for _, childInfo := range info.Children {
			tmp.Children = append(tmp.Children, s.recursionAreaInfo(childInfo, departmentMap, helper.RuleThree(pathName == "", info.Name, pathName+"/"+info.Name)))
		}
	}
	return tmp
}

func (s *AreaService) RemoveArea(c ctx.Context, params dashboardReq.AreaRemove) error {
	list, err := s.areaDao.GetList(c, dashboardReq.AreaTree{
		EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: params.EnterpriseId},
		ParentId:          params.Id,
		HasVirtual:        true,
	})
	if err != nil {
		return err
	}
	if len(list) > 0 {
		return errors.New("请先删除所有子区域")
	}
	return s.areaDao.DeleteData(c, params.EnterpriseId, map[string]interface{}{
		"id": params.Id,
	})
}

func (s *AreaService) SaveData(c ctx.Context, data dashboardReq.AreaSave) error {
	unixTime := time.Now().Unix()
	if data.Id > 0 {
		info, _ := s.areaDao.GetOne(c, data.EnterpriseId, "id = ?", data.Id)
		if info.ID == 0 {
			return errors.New("id参数错误")
		}
		var updateData = make(map[string]interface{})
		if data.ParentId > 0 {
			if info.Type < dashboardModel.AreaStore {
				return errors.New("只能移动门店")
			}
			parentAreaInfo, _ := s.areaDao.GetOne(c, data.EnterpriseId, "id = ?", data.ParentId)
			path := fmt.Sprintf("%s,%d", parentAreaInfo.Path, data.ParentId)
			updateData = map[string]interface{}{
				"parent_id":  data.ParentId,
				"updated_at": unixTime,
				"path":       path,
			}
		} else if data.Name != "" {
			if data.Name != "" {
				exist, _ := s.areaDao.GetOne(c, data.EnterpriseId, " id!=? AND parent_id = ? AND name = ? ", data.Id, info.ParentID, data.Name)
				if exist.ID > 0 {
					return errors.New("同级别下有重名区域")
				}
			}
			updateData = map[string]interface{}{
				"name":       data.Name,
				"updated_at": unixTime,
			}

		}
		return s.areaDao.UpdateData(c, data.EnterpriseId, map[string]interface{}{
			"id": data.Id,
		}, updateData)
	}

	var parentInfo dashboardModel.DashboardAreaModel
	var brother dashboardModel.DashboardAreaModel
	if data.ParentId > 0 {
		parentInfo, _ = s.areaDao.GetOne(c, data.EnterpriseId, "id = ?", data.ParentId)
		if parentInfo.ID == 0 {
			return errors.New("该父节点不存在")
		}
	} else if data.BrotherId > 0 {
		// 增加同级别节点,将大于基本节点的其他兄弟节点的order都加1
		brother, _ = s.areaDao.GetOne(c, data.EnterpriseId, "id = ? ", data.BrotherId)
		if brother.ID == 0 {
			return errors.New("该节点不存在")
		}
		if brother.ParentID > 0 {
			parentInfo, _ = s.areaDao.GetOne(c, data.EnterpriseId, "id = ?", brother.ParentID)
		}
	}
	if data.Name != "" {
		exist, _ := s.areaDao.GetOne(c, data.EnterpriseId, "parent_id = ? AND name = ?", parentInfo.ID, data.Name)
		if exist.ID > 0 {
			return errors.New("同级别下有重名区域")
		}
	}

	if data.ThirdId > 0 {
		exist, _ := s.areaDao.GetOne(c, data.EnterpriseId, "third_id = ?", data.ThirdId)
		if exist.ID > 0 {
			return errors.New("该门店已存在")
		}
		storeCodeMap, _ := microWorkDao.NewStoreDepartment().GetStoreCodeMap(c, data.EnterpriseId)
		if storeCode, ok := storeCodeMap[data.ThirdId]; ok {
			data.Name = storeCode
		}
	}
	addData := dashboardModel.DashboardAreaModel{
		Name:     data.Name,
		ParentID: data.ParentId,
		ThirdID:  data.ThirdId,
		Sort:     s.areaDao.GetMaxOrderByParentId(c, data.EnterpriseId, parentInfo.ID) + 1,
		BaseModel: models.BaseModel{
			CreatedAt: unixTime,
			UpdatedAt: unixTime,
		},
	}
	if parentInfo.ID > 0 {
		addData.Path = strings.Trim(parentInfo.Path+","+strconv.FormatInt(parentInfo.ID, 10), ",")
	}
	switch parentInfo.Type {
	case 0:
		addData.Type = dashboardModel.AreaSource
	case dashboardModel.AreaSource:
		addData.Type = dashboardModel.AreaRegion
	case dashboardModel.AreaRegion:
		addData.Type = dashboardModel.AreaArea
	case dashboardModel.AreaArea:
		addData.Type = dashboardModel.AreaStore
		if data.ThirdId == 0 {
			return errors.New("门店thirdId参数错误")
		}
	case dashboardModel.AreaStore:
		return errors.New("parentId参数错误")
	default:
		return errors.New("type参数错误")
	}
	if data.ParentId == 0 && data.BrotherId > 0 {
		addData.ParentID = brother.ParentID
		sign := ">"
		addData.Sort = brother.Sort + 1
		if data.IsAbove {
			sign = ">="
			addData.Sort = brother.Sort
		}
		sqlConn := fmt.Sprintf("parent_id = %d AND `sort` %s %d", parentInfo.ID, sign, brother.Sort)
		err := s.areaDao.Increment(c, data.EnterpriseId, sqlConn, "sort", 1)
		if err != nil {
			return err
		}
	}
	return s.areaDao.CreateData(c, data.EnterpriseId, addData)
}

// GetManagerAuthAreaIds  获取管理者组织架构权限节点
func (s *AreaService) GetManagerAuthAreaIds(c ctx.Context, params dashboardReq.ManagerAuthAreaIds) (areaIds []int64, err error) {
	areaIds, err = dashboardDao.NewManager().GetManagerAuthAreaIds(c, params)
	return
}

// GetThirdIdsByAreaList  根据节点集合获取管理范围内的门店thirdId集合
func (s *AreaService) GetThirdIdsByAreaList(c ctx.Context, params dashboardReq.ManagerAuthAreaIds) (areaList []dashboardResp.AreaInfo, err error) {
	// 需要返回的节点的级别
	var returnTypeSign int
	var returnAreaIdSign int64
	var relationAreaIds []int64
	oneParent := true
	allAreaList, err := s.areaDao.GetNormalAreaList(c, params.EnterpriseId)
	if params.TargetAreaId == 0 {
		// 非下钻则直接获取参数中的节点
		relationAreaIds = params.AreaIds
	}
	for _, area := range allAreaList {
		if helper.InArray(area.ID, params.AreaIds) {
			// 如果存在下钻目标节点则只获取与下钻目标节点相关的节点
			tmpPathSlice := s.getPathSlice(area, false) // 组装节点的path
			// 获取最高等级的节点(列表是根据级别由高到低）
			if returnTypeSign == 0 {
				returnTypeSign = area.Type
				returnAreaIdSign = area.ID
			} else {
				if !helper.InArray(returnAreaIdSign, tmpPathSlice) {
					// 如果某个节点的path中不存在目前的顶级节点id，说明有该顶级节点之外的节点，这种情况就需要展示顶级节点的层级
					oneParent = false
				}
			}
			if params.TargetAreaId > 0 {
				if helper.InArray(params.TargetAreaId, tmpPathSlice) {
					// 如果节点与入参的节点有交集 则说明该节点以及的path与需要的数据有关
					relationAreaIds = append(relationAreaIds, append(tmpPathSlice, area.ID)...)
				}
			} else {
				relationAreaIds = append(relationAreaIds, tmpPathSlice...)
			}
		}
		if params.TargetAreaId > 0 && params.TargetAreaId == area.ID {
			returnTypeSign = area.Type + 1 // 目标层级+1
		}
	}
	if oneParent && returnTypeSign < dashboardModel.AreaStore {
		// 如果只有一个顶级节点则只获取下级节点
		returnTypeSign++
	}
	var highestAreaMap = make(map[int64]dashboardResp.AreaInfo)
	for _, area := range allAreaList {
		if area.Type == dashboardModel.AreaStore && area.ThirdID <= 0 {
			continue
		}
		if area.Type >= returnTypeSign {
			if helper.InArray(area.ID, relationAreaIds) {
				if area.Type == returnTypeSign {
					// 如果节点的path与检索的节点有交集则说明是我们需要找的节点
					// 如果层级是最高级则是我们要找的顶级节点
					highestAreaMap[area.ID] = dashboardResp.AreaInfo{
						ID:   area.ID,
						Type: area.Type,
						Name: area.Name,
					}
				}
				if area.ThirdID > 0 {
					// 组装节点的path
					tmpPathSlice := s.getPathSlice(area, true)
					// 如果是门店 则将门店的thirdId写入顶级节点的数据中
					for _, pathAreaId := range tmpPathSlice {
						if areaInfo, ok := highestAreaMap[pathAreaId]; ok {
							highestAreaMap[pathAreaId] = dashboardResp.AreaInfo{
								ID:               areaInfo.ID,
								Type:             areaInfo.Type,
								Name:             areaInfo.Name,
								ChildrenThirdIds: append(highestAreaMap[pathAreaId].ChildrenThirdIds, area.ThirdID),
							}
						}
					}
				}
			}

		}
	}
	for _, info := range highestAreaMap {
		areaList = append(areaList, info)
	}
	return
}

// GetMemberAuthAreaIds 获取导购组织架构权限节点
func (s *AreaService) getPathSlice(area dashboardModel.DashboardAreaModel, hasSelf bool) []int64 {
	// 组装节点的path
	var tmpPathSlice []int64
	if hasSelf {
		tmpPathSlice = append(tmpPathSlice, area.ID)
	}
	for _, tmp := range strings.Split(area.Path, ",") {
		tmpPathSlice = append(tmpPathSlice, int64(helper.StringToInt(tmp)))
	}
	return tmpPathSlice
}

func (s *AreaService) AreaListForOpen(c ctx.Context, params openReq.AreaTreeReq) (list []openResp.AreaList, err error) {
	managerList, err := s.managerDao.GetAllOrderByPosition(c, params.EnterpriseId)
	if err != nil {
		return nil, err
	}
	var memberIds []int64
	var areaManagerMap = make(map[int64][]dashboardModel.DashboardManagerModel)
	for _, manager := range managerList {
		memberIds = append(memberIds, manager.MemberId)
		tmpSlice := strings.Split(manager.AreaIds, ",")
		for _, areaIdStr := range tmpSlice {
			areaId := int64(helper.StringToInt(areaIdStr))
			areaManagerMap[areaId] = append(areaManagerMap[areaId], manager)
		}
	}
	memberList, err := membersDao.NewMembers().GetAllById(c, params.EnterpriseId, memberIds)
	var memberMap = make(map[int64]membersModel.MemberModel)
	for _, member := range memberList {
		memberMap[member.ID] = member
	}
	allDepartmentList, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"deleted_at":    0,
	})
	if err != nil {
		return nil, err
	}
	var departmentMap = make(map[int]microWorkModel.StoreDepartmentModel)
	for _, department := range allDepartmentList {
		if department.WechatDepartmentID > 0 {
			departmentMap[int(department.WechatDepartmentID)] = department
		}
	}
	storeNameMap, _ := microWorkDao.NewStoreDepartment().GetStoreNameMap(c, params.EnterpriseId)
	areaList, err := s.areaDao.AreaListForOpen(c, params)
	if err != nil {
		return nil, err
	}
	for _, info := range areaList {
		list = append(list, s.recursionAreaInfoForOpen(info, departmentMap, memberMap, storeNameMap, areaManagerMap))
	}
	return
}

// 递归整理门店信息
func (s *AreaService) recursionAreaInfoForOpen(info dashboardModel.DashboardAreaModel, departmentMap map[int]microWorkModel.StoreDepartmentModel, memberMap map[int64]membersModel.MemberModel, storeNameMap map[int]string, areaManagerMap map[int64][]dashboardModel.DashboardManagerModel) openResp.AreaList {
	var tmp = openResp.AreaList{
		ID:          int(info.ID),
		LMSID:       info.LMSID,
		Type:        info.Type,
		Name:        info.Name,
		ParentID:    info.ParentID,
		ThirdID:     info.ThirdID,
		ManagerList: make([]openResp.Manager, 0),
		ClosedAt:    info.ClosedAt,
		UpdatedAt:   info.UpdatedAt,
	}
	if tmp.Type == dashboardModel.AreaSource {
		tmp.CustomGroup = info.FiscalYear
	}
	if tmp.Type == dashboardModel.AreaStore {
		storeCode := info.Name
		if v, ok := storeNameMap[tmp.ThirdID]; ok {
			tmp.Name = fmt.Sprintf("%s_%s", v, storeCode)
		}
	}
	if dept, ok := departmentMap[info.ThirdID]; ok {
		tmp.StoreCode = dept.ThirdStoreCode
		tmp.Name = dept.StoreName
		tmp.EnName = dept.InventoryName
	}
	var tmpManager openResp.Manager
	if managerList, ok := areaManagerMap[info.ID]; ok {
		for _, manager := range managerList {
			tmpManager.Position = manager.Position
			if member, ok := memberMap[manager.MemberId]; ok {
				tmpManager.UserId = member.Userid
				tmpManager.UserName = member.Name
			}
			tmp.ManagerList = append(tmp.ManagerList, tmpManager)
		}
	}
	if len(info.Children) > 0 {
		for _, childInfo := range info.Children {
			tmp.Children = append(tmp.Children, s.recursionAreaInfoForOpen(childInfo, departmentMap, memberMap, storeNameMap, areaManagerMap))
		}
	}
	return tmp
}

func (s *AreaService) AreaManager(c ctx.Context, params openReq.AreaManagerReq) (openResp.ManagerInfo, error) {
	memberInfo, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"userid":        params.UserId,
	})
	if err != nil {
		return openResp.ManagerInfo{}, err
	}
	info := openResp.ManagerInfo{
		UserId:    memberInfo.Userid,
		UserName:  memberInfo.Name,
		Position:  "",
		StoreList: nil,
	}
	manager, err := dashboardDao.NewManager().GetOne(c, params.EnterpriseId, map[string]interface{}{
		"member_id": memberInfo.ID,
	})
	tmpSlice := strings.Split(manager.AreaIds, ",")
	var areaIds []int64
	for _, val := range tmpSlice {
		areaIds = append(areaIds, int64(helper.StringToInt(val)))
	}

	info.Position = manager.Position

	thirdIds, err := s.areaDao.GetThirdIdsByAreaIds(c, dashboardReq.ManagerAuthAreaIds{
		EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: params.EnterpriseId},
		AreaIds:           areaIds,
	})
	if err != nil {
		return info, err
	}
	departmentList, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
		"enterprise_id":        params.EnterpriseId,
		"wechat_department_id": thirdIds,
		"deleted_at":           0,
	})
	if err != nil {
		return info, err
	}
	for _, store := range departmentList {
		info.StoreList = append(info.StoreList, openResp.StoreInfo{
			ThirdId:   int(store.WechatDepartmentID),
			StoreCode: store.ThirdStoreCode,
			StoreName: store.StoreName,
			EnName:    store.InventoryName,
		})
	}

	return info, nil
}

// StoreCloseOrOpen   门店闭店或者重新开店
func (s *AreaService) StoreCloseOrOpen(c ctx.Context, params dashboardReq.StoreCloseOrOpen) error {
	return dashboardDao.NewArea().StoreCloseOrOpen(c, params)
}

func (s *AreaService) SearchStore(c ctx.Context, params dashboardReq.SearchStoreReq) (list []microWorkModel.StoreDepartmentModel, err error) {
	searchStore, err := microWorkDao.NewStoreDepartment().SearchStore(c, params)
	if err != nil {
		return nil, err
	}
	for i, model := range searchStore {
		if params.Language == "en" && model.InventoryName != "" {
			searchStore[i].StoreName = model.InventoryName
		}
	}
	return searchStore, nil
}

func (s *AreaService) AreaList(c ctx.Context, enterpriseId int64) ([]dashboardResp.AreaInfo, error) {
	result := make([]dashboardResp.AreaInfo, 0)
	many, err := s.areaDao.GetMany(c, enterpriseId, map[string]interface{}{
		"type":       dashboardModel.AreaArea,
		"is_virtual": 0,
		"closed_at":  0,
	})
	if err != nil {
		return result, nil
	}
	for _, v := range many {
		result = append(result, dashboardResp.AreaInfo{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	return result, err
}
