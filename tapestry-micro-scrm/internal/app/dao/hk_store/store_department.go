package hkStoreDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	openReq "micro-scrm/internal/app/http/request/open"
	hkStoreModel "micro-scrm/internal/app/models/hk_store"
	"micro-scrm/internal/pkg/db/gormx"
)

type StoreDepartment struct {
	dao       dao.Dao
	tableName string
}

func NewStoreDepartment() *StoreDepartment {
	return &StoreDepartment{
		dao:       dao.Dao{DB: gormx.GetWorkInstance()},
		tableName: hkStoreModel.StoreDepartmentModel{}.TableName(),
	}
}

func (d *StoreDepartment) GetPagination(c app.Ctx, cond openReq.HkStoreListReq) ([]hkStoreModel.StoreDepartmentModel, int64, error) {
	var (
		total int64
		resp  []hkStoreModel.StoreDepartmentModel
	)
	query := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id = ? and deleted_at = 0", cond.EnterpriseId)
	if len(cond.StoreCodes) > 0 {
		query = query.Where("store_code in ?", cond.StoreCodes)
	}
	if cond.StoreName != "" {
		query = query.Where("store_name like ?", "%"+cond.StoreName+"%")
	}
	if cond.Area != "" {
		query = query.Where("area = ?", cond.Area)
	}
	if cond.WechatDepartmentID > 0 {
		query = query.Where("wechat_department_id = ?", cond.WechatDepartmentID)
	}
	if cond.ExcludeClosed {
		query = query.Where("closed_at = 0")
	}
	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size
	err := query.Order("created_at DESC").Limit(cond.Size).Offset(offset).Find(&resp).Error
	return resp, total, err
}

func (d *StoreDepartment) GetStoreList(c app.Ctx, cond dashboardReq.HkManagerStoreList) ([]hkStoreModel.StoreDepartmentModel, error) {
	var resp []hkStoreModel.StoreDepartmentModel
	query := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id = ? and deleted_at = 0", cond.EnterpriseId)
	if cond.StoreName != "" {
		query = query.Where("store_name like ?", "%"+cond.StoreName+"%")
	}
	if cond.StoreCode != "" {
		query = query.Where("store_code like ?", "%"+cond.StoreCode+"%")
	}
	if cond.Area != "" {
		query = query.Where("area = ?", cond.Area)
	}
	if cond.WechatDepartmentID > 0 {
		query = query.Where("wechat_department_id = ?", cond.WechatDepartmentID)
	}
	if cond.ExcludeClosed {
		query = query.Where("closed_at = 0")
	}
	err := query.Order("created_at DESC").Find(&resp).Error
	return resp, err
}

func (d *StoreDepartment) GetManagePagination(c app.Ctx, cond dashboardReq.HkStoreList) ([]hkStoreModel.StoreDepartmentModel, int64, error) {
	var (
		total int64
		resp  []hkStoreModel.StoreDepartmentModel
	)
	query := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ? and deleted_at = 0", cond.EnterpriseId)
	if cond.SearchStr != "" {
		searchStr := "%" + cond.SearchStr + "%"
		query = query.Where("(store_name like ? OR store_code like ?)", searchStr, searchStr)
	}
	if cond.Area != "" {
		query = query.Where("area = ?", cond.Area)
	}
	if cond.WechatDepartmentID > 0 {
		query = query.Where("wechat_department_id = ?", cond.WechatDepartmentID)
	}
	if cond.StoreType >= 0 {
		query = query.Where("store_type = ?", cond.StoreType)
	}
	switch cond.ClosedStatus {
	case 1:
		query = query.Where("closed_at = 0")
	case 2:
		query = query.Where("closed_at > 0")
	}
	if err := query.Count(&total).Error; err != nil {
		return resp, total, err
	}
	offset := (cond.Page - 1) * cond.Size
	err := query.Order("created_at DESC").Limit(cond.Size).Offset(offset).Find(&resp).Error
	return resp, total, err
}

func (d *StoreDepartment) GetOne(c app.Ctx, query interface{}, args ...interface{}) (hkStoreModel.StoreDepartmentModel, error) {
	var result hkStoreModel.StoreDepartmentModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *StoreDepartment) CreateData(c app.Ctx, data hkStoreModel.StoreDepartmentModel) error {
	return d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
}

func (d *StoreDepartment) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).Updates(data).Error
}

func (d *StoreDepartment) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]hkStoreModel.StoreDepartmentModel, error) {
	var result []hkStoreModel.StoreDepartmentModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *StoreDepartment) GetAllWechatDepartmentIds(c app.Ctx, enterpriseId int64) (departmentIds []int, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ? and deleted_at = 0", enterpriseId).
		Pluck("wechat_department_id", &departmentIds).Error
	return
}

func (d *StoreDepartment) GetWechatDepartmentIdsByWorkListReq(c app.Ctx, cond openReq.HkStoreWorkListReq) (departmentIds []int, err error) {
	query := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ? and deleted_at = 0", cond.EnterpriseId)
	if len(cond.StoreCodes) > 0 {
		query = query.Where("store_code in ?", cond.StoreCodes)
	}
	if cond.Area != "" {
		query = query.Where("area = ?", cond.Area)
	}
	if cond.WechatDepartmentID > 0 {
		query = query.Where("wechat_department_id = ?", cond.WechatDepartmentID)
	}
	err = query.Pluck("wechat_department_id", &departmentIds).Error
	return
}

func (d *StoreDepartment) GetMapByWechatDepartmentIds(c app.Ctx, enterpriseId int64, departmentIds []int) (map[int]hkStoreModel.StoreDepartmentModel, error) {
	storeList, err := d.GetMany(c, map[string]interface{}{
		"enterprise_id":        enterpriseId,
		"wechat_department_id": departmentIds,
		"deleted_at":           0,
	})
	if err != nil {
		return nil, err
	}
	result := make(map[int]hkStoreModel.StoreDepartmentModel)
	for _, store := range storeList {
		result[int(store.WechatDepartmentID)] = store
	}
	return result, nil
}

func (d *StoreDepartment) GetActiveMapByWechatDepartmentIds(c app.Ctx, enterpriseId int64, departmentIds []int) (map[int]hkStoreModel.StoreDepartmentModel, error) {
	storeList, err := d.GetMany(c, map[string]interface{}{
		"enterprise_id":        enterpriseId,
		"wechat_department_id": departmentIds,
		"deleted_at":           0,
		"closed_at":            0,
	})
	if err != nil {
		return nil, err
	}
	result := make(map[int]hkStoreModel.StoreDepartmentModel)
	for _, store := range storeList {
		result[int(store.WechatDepartmentID)] = store
	}
	return result, nil
}
