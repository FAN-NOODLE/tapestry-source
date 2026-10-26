package hkStoreDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	dashboardResp "micro-scrm/internal/app/http/response/dashboard"
	hkStoreModel "micro-scrm/internal/app/models/hk_store"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"strings"
)

type StoreManager struct {
	dao       dao.Dao
	model     *hkStoreModel.StoreManagerModel
	tableName string
}

func NewStoreManager() *StoreManager {
	return &StoreManager{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(hkStoreModel.StoreManagerModel),
		tableName: hkStoreModel.StoreManagerModel{}.TableName(),
	}
}

func (d *StoreManager) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (hkStoreModel.StoreManagerModel, error) {
	var result hkStoreModel.StoreManagerModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *StoreManager) GetAll(c app.Ctx, enterpriseId int64) ([]hkStoreModel.StoreManagerModel, error) {
	var result []hkStoreModel.StoreManagerModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Find(&result).Error
	return result, err
}

func (d *StoreManager) CreateData(c app.Ctx, enterpriseId int64, data hkStoreModel.StoreManagerModel) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
}

func (d *StoreManager) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).UpdateColumns(data).Error
}

func (d *StoreManager) DeleteData(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) error {
	var model hkStoreModel.StoreManagerModel
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Delete(&model).Error
}

func (d *StoreManager) GetPagination(c app.Ctx, cond dashboardReq.HkManagerAuthList) (list []dashboardResp.HkManagerList, total int64, err error) {
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).
		Select("hkmo_store_manager.id,hkmo_store_manager.member_id,hkmo_store_manager.position,hkmo_store_manager.area_ids,members.name,members.userid").
		Joins("LEFT JOIN micro_scrm.members ON hkmo_store_manager.member_id = members.id").
		Where("hkmo_store_manager.deleted_at = 0 AND members.deleted_at = 0").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if cond.UserId != "" {
				db.Where("members.userid = ?", cond.UserId)
			} else if cond.UserName != "" {
				db.Where("name like ?", "%"+cond.UserName+"%")
			}
			return db
		})
	query.Count(&total)
	offset := (cond.Page - 1) * cond.Size
	err = query.Order("hkmo_store_manager.created_at DESC").Limit(cond.Size).Offset(offset).Find(&list).Error
	return
}

func (d *StoreManager) GetManagerAuthAreaIds(c app.Ctx, params dashboardReq.HkManagerAuthAreaIds) (areaIds []int64, err error) {
	var areaIdsStr string
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, params.EnterpriseId)).Model(d.model).
		Where("member_id = ?", params.MemberId).
		Pluck("area_ids", &areaIdsStr).Error
	for _, areaIdStr := range strings.Split(areaIdsStr, ",") {
		areaId := int64(helper.StringToInt(areaIdStr))
		if areaId > 0 {
			areaIds = append(areaIds, areaId)
		}
	}
	return
}
