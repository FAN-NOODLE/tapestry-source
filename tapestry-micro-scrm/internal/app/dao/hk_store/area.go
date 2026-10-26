package hkStoreDao

import (
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	"micro-scrm/internal/app/models"
	hkStoreModel "micro-scrm/internal/app/models/hk_store"
	"micro-scrm/internal/pkg/db/gormx"
	"strings"
)

const AreaTypeStore = 4

type Area struct {
	dao       dao.Dao
	model     *hkStoreModel.AreaModel
	tableName string
}

func NewArea() *Area {
	return &Area{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(hkStoreModel.AreaModel),
		tableName: hkStoreModel.AreaModel{}.TableName(),
	}
}

func (d *Area) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (hkStoreModel.AreaModel, error) {
	var result hkStoreModel.AreaModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Area) GetAll(c app.Ctx, enterpriseId int64) ([]hkStoreModel.AreaModel, error) {
	var result []hkStoreModel.AreaModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Find(&result).Error
	return result, err
}

func (d *Area) GetList(c app.Ctx, params dashboardReq.AreaTree) ([]hkStoreModel.AreaModel, error) {
	db := gormx.GetConnection(fmt.Sprintf("%s%d", app.Cfg().DB.DBPrefix, params.EnterpriseId))
	var result []hkStoreModel.AreaModel
	orm := db.Model(d.model).WithContext(c)
	if !params.HasChildren {
		orm = orm.Where("parent_id = ?", params.ParentId)
	}
	if !params.HasVirtual {
		orm = orm.Where("is_virtual = 0")
	}
	if params.IgnoreCloseStore {
		orm = orm.Where("type <> ? OR closed_at = 0", AreaTypeStore)
	}
	err := orm.Order("type asc, sort asc").Find(&result).Error
	return result, err
}

func (d *Area) GetListByAreaIds(c app.Ctx, enterpriseId int64, areaIds []int64) ([]hkStoreModel.AreaModel, error) {
	var result []hkStoreModel.AreaModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).
		Where("is_virtual = 0").
		Where("id in ?", areaIds).
		Where("closed_at = 0").
		Order("type asc, sort asc").
		Find(&result).Error
	return result, err
}

func (d *Area) GetThirdIdsByAreaIds(c app.Ctx, params dashboardReq.HkManagerAuthAreaIds) (thirdIds []int, err error) {
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, params.EnterpriseId)).Model(d.model).
		Where("type = ?", AreaTypeStore).
		Where("is_virtual = 0").
		Where("closed_at = 0")
	condition := ""
	var args []interface{}
	for idx, areaId := range params.AreaIds {
		if idx > 0 {
			condition += " OR "
		}
		condition += "(find_in_set(?, path) OR id = ?)"
		args = append(args, areaId, areaId)
	}
	if condition != "" {
		query = query.Where(condition, args...)
	}
	err = query.Pluck("third_id", &thirdIds).Error
	return
}

func (d *Area) UpdateStoreNode(c app.Ctx, enterpriseId int64, store hkStoreModel.StoreDepartmentModel) error {
	if store.WechatDepartmentID == 0 {
		return nil
	}
	data := map[string]interface{}{
		"name":       store.StoreCode,
		"closed_at":  store.ClosedAt,
		"updated_at": store.UpdatedAt,
	}
	if store.Area != "" {
		parent, err := d.GetOne(c, enterpriseId, "parent_id = 0 AND name = ?", strings.ToUpper(store.Area))
		if err != nil {
			return err
		}
		if parent.ID > 0 {
			data["parent_id"] = parent.ID
			data["path"] = fmt.Sprintf("%d", parent.ID)
		}
	}
	area, err := d.GetOne(c, enterpriseId, "type = ? AND third_id = ?", AreaTypeStore, store.WechatDepartmentID)
	if err != nil {
		return err
	}
	if area.ID == 0 {
		area = hkStoreModel.AreaModel{
			BaseModel: models.CreateBase(app.UID()),
			Type:      AreaTypeStore,
			Name:      store.StoreCode,
			ThirdID:   int(store.WechatDepartmentID),
			ClosedAt:  store.ClosedAt,
		}
		if parentID, ok := data["parent_id"].(int64); ok {
			area.ParentID = int(parentID)
		}
		if path, ok := data["path"].(string); ok {
			area.Path = path
		}
		return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&area).Error
	}
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).
		Where("id = ?", area.ID).
		UpdateColumns(data).Error
}
