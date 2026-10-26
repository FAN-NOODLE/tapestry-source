package dashboardSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	dashboardDao "micro-scrm/internal/app/dao/dashboard"
	membersDao "micro-scrm/internal/app/dao/members"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/http/request"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	dashboardResp "micro-scrm/internal/app/http/response/dashboard"
	"micro-scrm/internal/app/models"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ManagerService struct {
	service.Base
	managerDao *dashboardDao.Manager
}

func NewMemberService() *ManagerService {
	return &ManagerService{
		managerDao: dashboardDao.NewManager(),
	}
}

// SaveManagerAuth  更新管理者dashboard权限
func (s *ManagerService) SaveManagerAuth(c ctx.Context, params dashboardReq.ManagerAuthSave) error {
	//获取节点列表
	areaList, err := dashboardDao.NewArea().GetListByAreaIds(c, params.EnterpriseId, params.AreaIds)
	if err != nil {
		return err
	}
	var areaPathMap = make(map[int64][]int64) //节点id与节点path的map关系
	for _, area := range areaList {
		tmpIds := strings.Split(area.Path, ",")
		var intAreaIds []int64
		for _, areaIdStr := range tmpIds {
			intAreaIds = append(intAreaIds, int64(helper.StringToInt(areaIdStr)))
		}
		areaPathMap[area.ID] = intAreaIds
	}
	// 按照节点type类型由高到低排序
	sort.Slice(areaList, func(i, j int) bool {
		return areaList[i].Type < areaList[j].Type
	})
	var areaIds []int64 // 需要保存的areaIds
	var areaIdsStr string
	for _, area := range areaList {
		if pathIds, ok := areaPathMap[area.ID]; ok {
			// 如果新增节点的path中已经存在需要保存的areaId ，则说明该节点的父级以及被保存了（忽略该节点）
			if len(helper.ArrayIntersect(pathIds, areaIds)) > 0 {
				continue
			}
			areaIds = append(areaIds, area.ID)
			areaIdsStr += strconv.FormatInt(area.ID, 10) + ","
		}
	}
	// 获取入参中的memberId职位
	memberInfo, _ := s.managerDao.GetOne(c, params.EnterpriseId, "member_id = ?", params.MemberId)
	if params.ID > 0 {
		info, _ := s.managerDao.GetOne(c, params.EnterpriseId, "id = ?", params.ID)
		if info.ID == 0 {
			return errorx.NewErrorMsg(1001, "您要更新的职位信息不存在")
		}
		if memberInfo.ID > 0 && memberInfo.ID != info.ID {
			// 更新职位时，如果传的member_id已经存在职位且职位id不是当前编辑的职位id
			return errorx.NewErrorMsg(1001, "您选择的导购已经存在职位")
		}
		return s.managerDao.UpdateData(c, params.EnterpriseId, map[string]interface{}{
			"id": info.ID,
		}, map[string]interface{}{
			"member_id":  params.MemberId,
			"area_ids":   strings.TrimRight(areaIdsStr, ","),
			"position":   params.Position,
			"updated_at": time.Now().Unix(),
		})
	} else {
		if memberInfo.ID > 0 {
			return errorx.NewErrorMsg(1001, "您选择的导购已经存在职位")
		}
		unixTime := time.Now().Unix()
		return s.managerDao.CreateData(c, params.EnterpriseId, dashboardModel.DashboardManagerModel{
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

}

// ManagerAuthDetail  管理者权限详情
func (s *ManagerService) ManagerAuthDetail(c ctx.Context, params dashboardReq.ManagerAuthDetail) (detail dashboardResp.ManagerAuthDetail, err error) {
	info, err := s.managerDao.GetOne(c, params.EnterpriseId, "id = ?", params.ID)
	if err != nil {
		return detail, err
	}
	memberInfo, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{"enterprise_id": params.EnterpriseId, "id": info.MemberId})
	detail.MemberInfo = memberInfo
	detail.ID = info.ID
	detail.Position = info.Position
	tmpIds := strings.Split(info.AreaIds, ",")
	for _, areaIdStr := range tmpIds {
		detail.AreaIds = append(detail.AreaIds, int64(helper.StringToInt(areaIdStr)))
	}
	thirdIds, _ := dashboardDao.NewArea().GetThirdIdsByAreaIds(c, dashboardReq.ManagerAuthAreaIds{
		EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: params.EnterpriseId},
		AreaIds:           detail.AreaIds,
	})
	detail.StoreNum = len(thirdIds)
	return
}

// ManagerAuthList  管理者权限列表
func (s *ManagerService) ManagerAuthList(c ctx.Context, params dashboardReq.ManagerAuthList) ([]dashboardResp.ManagerList, int64, error) {
	list, total, err := s.managerDao.GetPagination(c, params)
	if err != nil {
		return nil, 0, err
	}
	allList, err := dashboardDao.NewArea().GetAll(c, params.EnterpriseId)
	if err != nil {
		return nil, 0, err
	}
	var storeList []dashboardModel.DashboardAreaModel
	var areaMap = make(map[int64]dashboardModel.DashboardAreaModel)
	for _, area := range allList {
		if area.Type == dashboardModel.AreaStore {
			storeList = append(storeList, area)
		}
		areaMap[area.ID] = area
	}
	storeNameMap, _ := microWorkDao.NewStoreDepartment().GetStoreNameMap(c, params.EnterpriseId)
	for i, manager := range list {
		areaIds := strings.Split(manager.AreaIds, ",")
		for _, tmpAreaId := range areaIds {
			//处理管理的节点名
			if area, ok := areaMap[int64(helper.StringToInt(tmpAreaId))]; ok {
				if area.Type == dashboardModel.AreaStore {
					if storeName, storeOk := storeNameMap[area.ThirdID]; storeOk {
						manager.PathNameArr = append(manager.PathNameArr, storeName)
					}
				} else {
					manager.PathNameArr = append(manager.PathNameArr, area.Name)
				}
			}

		}
		// 如果门店的path和导购的管理权限有交集 说明这个门店就是管理范围内的
		for _, store := range storeList {
			pathIds := strings.Split(store.Path, ",")
			if len(helper.ArrayIntersect(append(pathIds, helper.Int64ToString(store.ID)), areaIds)) > 0 {
				manager.StoreNum++
			}
		}
		list[i] = manager
	}
	return list, total, nil
}

// GetAllManagerMemberIds 获取所有已配置管理者
func (s *ManagerService) GetAllManagerMemberIds(c ctx.Context, enterpriseId int64) (memberIds []int64, err error) {

	list, err := s.managerDao.GetAll(c, enterpriseId)
	if err != nil {
		return memberIds, err
	}
	for _, manager := range list {
		memberIds = append(memberIds, manager.MemberId)
	}
	return
}

// DeleteManagerAuth  获取所有已配置管理者
func (s *ManagerService) DeleteManagerAuth(c ctx.Context, params dashboardReq.ManagerRemove) error {
	return s.managerDao.DeleteData(c, params.EnterpriseId, "id=?", params.Id)

}
