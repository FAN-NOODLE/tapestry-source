package dashboardSvc

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	hkStoreDao "micro-scrm/internal/app/dao/hk_store"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	dashboardResp "micro-scrm/internal/app/http/response/dashboard"
	"micro-scrm/internal/app/models"
	hkStoreModel "micro-scrm/internal/app/models/hk_store"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk/boss"
	"strconv"
	"time"
)

const hkStoreRootDepartmentConfigKey = "hkmo_store_root_dept_id"
const hkStoreLastSyncTimeKey = "%d:HKMO_STORE_LAST_SYNC_TIME"

type HkStoreService struct {
	service.Base
	storeDao *hkStoreDao.StoreDepartment
	areaDao  *hkStoreDao.Area
}

func NewHkStoreService() *HkStoreService {
	return &HkStoreService{
		storeDao: hkStoreDao.NewStoreDepartment(),
		areaDao:  hkStoreDao.NewArea(),
	}
}

func (s *HkStoreService) Sync(c ctx.Context, params dashboardReq.HkStoreSync) (dashboardResp.HkStoreSyncResult, error) {
	var result dashboardResp.HkStoreSyncResult
	configData, err := boss.EnterpriseConfig(c, map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"key":           hkStoreRootDepartmentConfigKey,
	})
	if err != nil {
		return result, err
	}
	if len(configData.Data) == 0 || configData.Data[0].Value <= 0 {
		return result, errors.New("未获取到香港门店父级部门配置")
	}
	departmentList, err := departmentsDao.NewDepartments().GetListByThirdParentId(
		c,
		params.EnterpriseId,
		int(configData.Data[0].Value),
	)
	if err != nil {
		return result, err
	}
	departmentIdMap := make(map[int]struct{}, len(departmentList))
	for _, department := range departmentList {
		departmentIdMap[department.ThirdID] = struct{}{}
		store, err := s.storeDao.GetOne(c, map[string]interface{}{
			"enterprise_id":        params.EnterpriseId,
			"wechat_department_id": department.ThirdID,
		})
		if err != nil {
			return result, err
		}
		if store.ID == 0 {
			store = hkStoreModel.StoreDepartmentModel{
				MicroBaseModel:     models.CreateMicroBase(params.EnterpriseId, app.UID()),
				StoreCode:          department.Name,
				WechatDepartmentID: uint(department.ThirdID),
				StoreName:          department.Name,
			}
			err = s.storeDao.CreateData(c, store)
			if err != nil {
				return result, err
			}
			if err = s.areaDao.UpdateStoreNode(c, params.EnterpriseId, store); err != nil {
				return result, err
			}
			result.CreatedCount++
			continue
		}
		if err = s.areaDao.UpdateStoreNode(c, params.EnterpriseId, store); err != nil {
			return result, err
		}
		result.SkippedCount++
	}
	activeStoreList, err := s.storeDao.GetMany(c, map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"deleted_at":    0,
		"closed_at":     0,
	})
	if err != nil {
		return result, err
	}
	for _, store := range activeStoreList {
		if _, ok := departmentIdMap[int(store.WechatDepartmentID)]; ok {
			continue
		}
		updateData := map[string]interface{}{
			"closed_at":  time.Now().Unix(),
			"updated_at": time.Now().Unix(),
		}
		err = s.storeDao.UpdateData(c, map[string]interface{}{
			"id":            store.ID,
			"enterprise_id": params.EnterpriseId,
		}, updateData)
		if err != nil {
			return result, err
		}
		store.ClosedAt = updateData["closed_at"].(int64)
		store.UpdatedAt = updateData["updated_at"].(int64)
		if err = s.areaDao.UpdateStoreNode(c, params.EnterpriseId, store); err != nil {
			return result, err
		}
		result.ClosedCount++
	}
	lastSyncTime := time.Now().Unix()
	if err = redis.Client().Set(c, fmt.Sprintf(hkStoreLastSyncTimeKey, params.EnterpriseId), lastSyncTime, 0).Err(); err != nil {
		return result, err
	}
	return result, nil
}

func (s *HkStoreService) LastSyncTime(c ctx.Context, params dashboardReq.HkStoreLastSyncTime) (dashboardResp.HkStoreLastSyncTime, error) {
	value, err := redis.Client().Get(c, fmt.Sprintf(hkStoreLastSyncTimeKey, params.EnterpriseId)).Result()
	if err != nil {
		if err == redis.Nil {
			return dashboardResp.HkStoreLastSyncTime{}, nil
		}
		return dashboardResp.HkStoreLastSyncTime{}, err
	}
	lastSyncTime, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return dashboardResp.HkStoreLastSyncTime{}, err
	}
	return dashboardResp.HkStoreLastSyncTime{LastSyncTime: lastSyncTime}, nil
}

func (s *HkStoreService) List(c ctx.Context, params dashboardReq.HkStoreList) ([]dashboardResp.HkStoreDetail, int64, error) {
	if params.StoreType < -1 || params.StoreType > 2 {
		return nil, 0, errors.New("门店类型参数错误")
	}
	storeList, total, err := s.storeDao.GetManagePagination(c, params)
	if err != nil {
		return nil, 0, err
	}
	var result []dashboardResp.HkStoreDetail
	for _, store := range storeList {
		result = append(result, s.formatStore(store))
	}
	return result, total, nil
}

func (s *HkStoreService) Detail(c ctx.Context, params dashboardReq.HkStoreDetail) (dashboardResp.HkStoreDetail, error) {
	store, err := s.storeDao.GetOne(c, map[string]interface{}{
		"id":            params.ID,
		"enterprise_id": params.EnterpriseId,
		"deleted_at":    0,
	})
	if err != nil {
		return dashboardResp.HkStoreDetail{}, err
	}
	if store.ID == 0 {
		return dashboardResp.HkStoreDetail{}, errors.New("香港门店不存在")
	}
	return s.formatStore(store), nil
}

func (s *HkStoreService) Save(c ctx.Context, params dashboardReq.HkStoreSave) error {
	if params.StoreType < 0 || params.StoreType > 2 {
		return errors.New("门店类型参数错误")
	}
	store, err := s.storeDao.GetOne(c, map[string]interface{}{
		"id":            params.ID,
		"enterprise_id": params.EnterpriseId,
		"deleted_at":    0,
	})
	if err != nil {
		return err
	}
	if store.ID == 0 {
		return errors.New("香港门店不存在")
	}
	updateData := map[string]interface{}{
		"store_name":     params.StoreName,
		"area":           params.Area,
		"store_type":     params.StoreType,
		"inventory_name": params.InventoryName,
		"address":        params.Address,
		"city":           params.City,
		"phone":          params.Phone,
		"email":          params.Email,
		"updated_at":     time.Now().Unix(),
	}
	if err = s.storeDao.UpdateData(c, map[string]interface{}{
		"id":            params.ID,
		"enterprise_id": params.EnterpriseId,
	}, updateData); err != nil {
		return err
	}
	store.StoreName = params.StoreName
	store.Area = params.Area
	store.StoreType = int32(params.StoreType)
	store.InventoryName = params.InventoryName
	store.Address = params.Address
	store.City = params.City
	store.Phone = params.Phone
	store.Email = params.Email
	store.UpdatedAt = updateData["updated_at"].(int64)
	return s.areaDao.UpdateStoreNode(c, params.EnterpriseId, store)
}

func (s *HkStoreService) Close(c ctx.Context, params dashboardReq.HkStoreClose) error {
	store, err := s.storeDao.GetOne(c, map[string]interface{}{
		"id":            params.ID,
		"enterprise_id": params.EnterpriseId,
		"deleted_at":    0,
	})
	if err != nil {
		return err
	}
	if store.ID == 0 {
		return errors.New("香港门店不存在")
	}
	if store.ClosedAt > 0 {
		return nil
	}
	updateData := map[string]interface{}{
		"closed_at":  time.Now().Unix(),
		"updated_at": time.Now().Unix(),
	}
	if err = s.storeDao.UpdateData(c, map[string]interface{}{
		"id":            params.ID,
		"enterprise_id": params.EnterpriseId,
	}, updateData); err != nil {
		return err
	}
	store.ClosedAt = updateData["closed_at"].(int64)
	store.UpdatedAt = updateData["updated_at"].(int64)
	return s.areaDao.UpdateStoreNode(c, params.EnterpriseId, store)
}

func (s *HkStoreService) formatStore(store hkStoreModel.StoreDepartmentModel) dashboardResp.HkStoreDetail {
	return dashboardResp.HkStoreDetail{
		ID:                 store.ID,
		StoreCode:          store.StoreCode,
		WechatDepartmentID: store.WechatDepartmentID,
		StoreName:          store.StoreName,
		Area:               store.Area,
		StoreType:          store.StoreType,
		InventoryName:      store.InventoryName,
		Address:            store.Address,
		City:               store.City,
		Phone:              store.Phone,
		Email:              store.Email,
		ClosedAt:           store.ClosedAt,
		CreatedAt:          store.CreatedAt,
		UpdatedAt:          store.UpdatedAt,
	}
}
