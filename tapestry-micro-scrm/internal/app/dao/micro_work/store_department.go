package microWorkDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	openReq "micro-scrm/internal/app/http/request/open"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/pkg/db/gormx"
)

const (
	// STORE_TYPE_OUTLETS 奥莱
	STORE_TYPE_OUTLETS = 1
	// STORE_TYPE_RETAIL 零售
	STORE_TYPE_RETAIL = 2
)

type StoreDepartment struct {
	dao       dao.Dao
	tableName string
}

func NewStoreDepartment() *StoreDepartment {
	return &StoreDepartment{
		dao:       dao.Dao{DB: gormx.GetWorkInstance()},
		tableName: microWorkModel.StoreDepartmentModel{}.TableName(),
	}
}

func (d *StoreDepartment) GetInfo(c app.Ctx, cond interface{}) (microWorkModel.StoreDepartmentModel, error) {
	var result microWorkModel.StoreDepartmentModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).Find(&result).Error
	return result, err
}

func (d *StoreDepartment) GetPagination(c app.Ctx, cond openReq.StoreListReq) ([]microWorkModel.StoreDepartmentModel, int64, error) {
	var (
		total int64
		resp  []microWorkModel.StoreDepartmentModel
		err   error
	)

	query := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id = ? and deleted_at = 0", cond.EnterpriseId)
	if len(cond.StoreCodes) > 0 {
		query.Where("third_store_code in ?", cond.StoreCodes)
	}
	if cond.StoreName != "" {
		query.Where(" store_name like ? ", "%"+cond.StoreName+"%")
	}
	if cond.StoreType > 0 {
		query.Where(" store_type=?", cond.StoreType)
	}
	if cond.WechatDepartmentID > 0 {
		query.Where(" wechat_department_id=?", cond.WechatDepartmentID)
	}
	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size
	err = query.Order("created_at DESC").Limit(cond.Size).Offset(offset).Find(&resp).Error

	return resp, total, err
}

func (d *StoreDepartment) GetAllByWechatDepartmentIds(c app.Ctx, enterpriseId int64, departmentIds []int) ([]microWorkModel.StoreDepartmentModel, error) {
	var result []microWorkModel.StoreDepartmentModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ? and wechat_department_id in (?)", enterpriseId, departmentIds).
		Find(&result).Error
	return result, err
}

// GetStoreName 根据wechatDepartmentId 获取门店名称
func (d *StoreDepartment) GetStoreName(c app.Ctx, enterpriseId int64, wechatDepartmentId []int) ([]microWorkModel.StoreDepartmentModel, error) {
	var result []microWorkModel.StoreDepartmentModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id":        enterpriseId,
		"wechat_department_id": wechatDepartmentId,
	}).Find(&result).Error
	return result, err
}

// GetDpIdsByType 获取奥莱或零售的企微部门id集合
func (d *StoreDepartment) GetDpIdsByType(c app.Ctx, enterpriseId int64, storeType int64) (dpIds []string, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"store_type":    storeType,
	}).Pluck("wechat_department_id", &dpIds).Error
	return
}

func (d *StoreDepartment) GetStoreInfoByStoreCode(c app.Ctx, enterpriseId int64, storeCode []string) ([]microWorkModel.StoreDepartmentModel, error) {
	var result []microWorkModel.StoreDepartmentModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id":    enterpriseId,
		"third_store_code": storeCode,
	}).Find(&result).Error
	return result, err
}

func (d *StoreDepartment) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]microWorkModel.StoreDepartmentModel, error) {
	var result []microWorkModel.StoreDepartmentModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *StoreDepartment) GetAll(c app.Ctx, enterpriseId int64) ([]microWorkModel.StoreDepartmentModel, error) {
	var result []microWorkModel.StoreDepartmentModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ? ", enterpriseId).
		Where("wechat_department_id > 0 ").
		Find(&result).Error
	return result, err
}

func (d *StoreDepartment) GetStoreNameMap(c app.Ctx, enterpriseId int64) (map[int]string, error) {
	storeNameMap := make(map[int]string, 0)
	many, err := d.GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
	})
	if err != nil {
		return storeNameMap, err
	}
	for _, v := range many {
		storeNameMap[int(v.WechatDepartmentID)] = v.StoreName
	}
	return storeNameMap, err
}

func (d *StoreDepartment) GetStoreCodeMap(c app.Ctx, enterpriseId int64) (map[int]string, error) {
	storeCodeMap := make(map[int]string, 0)
	many, err := d.GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
	})
	if err != nil {
		return storeCodeMap, err
	}
	for _, v := range many {
		storeCodeMap[int(v.WechatDepartmentID)] = v.ThirdStoreCode
	}
	return storeCodeMap, err
}

func (d *StoreDepartment) GetThirdIdMap(c app.Ctx, enterpriseId int64) (map[string]int, error) {
	thirdIdMap := make(map[string]int, 0)
	many, err := d.GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
	})
	if err != nil {
		return thirdIdMap, err
	}
	for _, v := range many {
		if v.WechatDepartmentID > 0 {
			thirdIdMap[v.ThirdStoreCode] = int(v.WechatDepartmentID)
		}
	}
	return thirdIdMap, err
}

func (d *StoreDepartment) SearchStore(c app.Ctx, cond dashboardReq.SearchStoreReq) ([]microWorkModel.StoreDepartmentModel, error) {
	var (
		resp []microWorkModel.StoreDepartmentModel
		err  error
	)
	query := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id = ? and deleted_at = 0", cond.EnterpriseId)
	if cond.StoreStr != "" {
		query.Where("(store_name like ? OR third_store_code like ? OR inventory_name like ? )", "%"+cond.StoreStr+"%", "%"+cond.StoreStr+"%", "%"+cond.StoreStr+"%")
	}
	if cond.IgnoreCloseStore {
		query.Where("closed_at = 0")
	}
	err = query.Order("created_at DESC").Limit(100).Find(&resp).Error

	return resp, err
}

// GetThirdIdsByStoreType 获取奥莱或零售的企微部门id集合
func (d *StoreDepartment) GetThirdIdsByStoreType(c app.Ctx, enterpriseId int64, storeType int) (dpIds []int, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"store_type":    storeType,
	}).Pluck("wechat_department_id", &dpIds).Error
	return
}
