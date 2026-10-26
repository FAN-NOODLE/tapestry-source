package dashboardDao

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	openReq "micro-scrm/internal/app/http/request/open"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"time"
)

type Area struct {
	dao       dao.Dao
	model     *dashboardModel.DashboardAreaModel
	tableName string
}

func NewArea() *Area {
	return &Area{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(dashboardModel.DashboardAreaModel),
		tableName: dashboardModel.DashboardAreaModel{}.TableName(),
	}
}

func (d *Area) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (dashboardModel.DashboardAreaModel, error) {
	var result dashboardModel.DashboardAreaModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Area) GetAll(c app.Ctx, enterpriseId int64) ([]dashboardModel.DashboardAreaModel, error) {
	var result []dashboardModel.DashboardAreaModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Find(&result).Error
	return result, err
}

func (d *Area) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]dashboardModel.DashboardAreaModel, error) {
	var result []dashboardModel.DashboardAreaModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Area) CreateData(c app.Ctx, enterpriseId int64, data dashboardModel.DashboardAreaModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
	return err
}

func (d *Area) CreateInBatchesData(c app.Ctx, enterpriseId int64, data []dashboardModel.DashboardAreaModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *Area) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *Area) DeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model dashboardModel.DashboardAreaModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).Delete(&model).Error
	return err
}

func (d *Area) Increment(c app.Ctx, enterpriseId int64, query interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s + ?", column), num)).Error
	return err
}

func (d *Area) GetList(c app.Ctx, params dashboardReq.AreaTree) ([]dashboardModel.DashboardAreaModel, error) {
	db := gormx.GetConnection(fmt.Sprintf("%s%d", app.Cfg().DB.DBPrefix, params.EnterpriseId))
	var result []dashboardModel.DashboardAreaModel
	orm := db.Model(d.model).WithContext(c).Where("parent_id = ?", params.ParentId)
	if !params.HasVirtual {
		orm = orm.Where("is_virtual = 0")
	}
	if params.HasChildren {
		orm = orm.Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("type asc,sort asc").Preload("Children", func(db *gorm.DB) *gorm.DB {
				return db.Order("type asc,sort asc").Preload("Children", func(db *gorm.DB) *gorm.DB {
					if params.IgnoreCloseStore {
						return db.Where("closed_at = 0")
					}
					return db
				})
			})
		})
	}
	err := orm.Order("type asc,sort asc").Find(&result).Error
	return result, err
}

func (d *Area) StoreCloseOrOpen(c app.Ctx, params dashboardReq.StoreCloseOrOpen) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, params.EnterpriseId)).Model(d.model).
		Where("third_id = ?", params.ThirdId).UpdateColumns(map[string]any{
		"closed_at":  helper.RuleThree(params.IsClosed, time.Now().Unix(), 0),
		"updated_at": time.Now().Unix(),
	}).Error
	return err
}

func (d *Area) GetThirdIdsByAreaIds(c app.Ctx, params dashboardReq.ManagerAuthAreaIds) (thirdIds []int, err error) {
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, params.EnterpriseId)).Model(d.model).
		Where("type=?", dashboardModel.AuthStore).
		Where("is_virtual = 0").
		Where("closed_at=0")
	// 查询管理范围内所有门店
	condition := ""
	var param []interface{}
	for idx, areaId := range params.AreaIds {
		if idx > 0 {
			condition += " OR "
		}
		condition += "(find_in_set(?, path) OR id=?)"
		param = append(param, areaId, areaId)
	}

	// 添加条件到查询中
	if condition != "" {
		query = query.Where(condition, param...)
	}
	// tmpSql := "("
	// for i, areaId := range params.AreaIds {
	// 	if i > 0 {
	// 		tmpSql += " or "
	// 	}
	// 	tmpSql += "find_in_set(" + strconv.Itoa(int(areaId)) + ",path)"
	// }
	// tmpSql += ")"
	err = query.Pluck("third_id", &thirdIds).Error
	return
}

func (d *Area) GetNormalAreaList(c app.Ctx, enterpriseId int64) ([]dashboardModel.DashboardAreaModel, error) {
	var result []dashboardModel.DashboardAreaModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).
		Where("is_virtual = 0").
		Where("closed_at=0").
		Order("type ASC").
		Find(&result).Error
	return result, err
}

func (d *Area) AreaListForOpen(c app.Ctx, params openReq.AreaTreeReq) ([]dashboardModel.DashboardAreaModel, error) {
	db := gormx.GetConnection(fmt.Sprintf("%s%d", app.Cfg().DB.DBPrefix, params.EnterpriseId))
	var result []dashboardModel.DashboardAreaModel
	orm := db.Model(d.model).WithContext(c).
		Where("parent_id = ?", params.ParentId).
		Where("is_virtual = 0")
	if params.HasChildren {
		orm = orm.Preload("Children.Children.Children", func(db *gorm.DB) *gorm.DB {
			if params.ExcludeClosed {
				return db.Where("closed_at = 0")
			}
			return db
		}) // 向下查3层
	}
	err := orm.Find(&result).Error
	return result, err
}

func (d *Area) GetAreaListByTypeAsc(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]dashboardModel.DashboardAreaModel, error) {
	var result []dashboardModel.DashboardAreaModel
	err := d.dao.WithContext(c).
		Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).
		Where(query, args...).
		Order("type asc").
		Find(&result).Error
	return result, err
}

func (d *Area) GetMaxOrderByParentId(c app.Ctx, enterpriseId int64, parentId int64) int64 {
	var result int64
	_ = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).
		Select("max(`sort`) as max_sort").
		Where("parent_id = ?", parentId).Pluck("max_sort", &result).Error
	return result
}

func (d *Area) GetListByAreaIds(c app.Ctx, enterpriseId int64, areaIds []int64) ([]dashboardModel.DashboardAreaModel, error) {
	var result []dashboardModel.DashboardAreaModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).
		Where("is_virtual = 0").
		Where("id in ?", areaIds).
		Where("closed_at=0").
		Order("type asc, sort asc").
		Find(&result).Error
	return result, err
}

func (d *Area) GetAllStore(c app.Ctx, enterpriseId int64) ([]dashboardModel.DashboardAreaModel, error) {
	var result []dashboardModel.DashboardAreaModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).
		Where("type = ?", dashboardModel.AreaStore).Find(&result).Error
	return result, err
}

func (d *Area) DeleteOldData(c app.Ctx, enterpriseId int64, unixTime int64) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).
		Where("updated_at < ?", unixTime).Delete(&d.model).Error
	return err
}
