package dashboardDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	dashboardResp "micro-scrm/internal/app/http/response/dashboard"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"strings"
)

type Manager struct {
	dao       dao.Dao
	model     *dashboardModel.DashboardManagerModel
	tableName string
}

func NewManager() *Manager {
	return &Manager{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(dashboardModel.DashboardManagerModel),
		tableName: dashboardModel.DashboardManagerModel{}.TableName(),
	}
}

func (d *Manager) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (dashboardModel.DashboardManagerModel, error) {
	var result dashboardModel.DashboardManagerModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Manager) GetAll(c app.Ctx, enterpriseId int64) ([]dashboardModel.DashboardManagerModel, error) {
	var result []dashboardModel.DashboardManagerModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Find(&result).Error
	return result, err
}

func (d *Manager) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]dashboardModel.DashboardManagerModel, error) {
	var result []dashboardModel.DashboardManagerModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Manager) CreateData(c app.Ctx, enterpriseId int64, data dashboardModel.DashboardManagerModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
	return err
}

func (d *Manager) CreateInBatchesData(c app.Ctx, enterpriseId int64, data []dashboardModel.DashboardManagerModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *Manager) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).UpdateColumns(data).Error
	return err
}

func (d *Manager) DeleteData(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) error {
	var model dashboardModel.DashboardManagerModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Delete(&model).Error
	return err
}

func (d *Manager) GetManagerAuthAreaIds(c app.Ctx, params dashboardReq.ManagerAuthAreaIds) (areaIds []int64, err error) {
	var tmpString string
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, params.EnterpriseId)).Model(d.model).
		Where("member_id=?", params.MemberId).
		Pluck("area_ids", &tmpString).Error
	tmpSlice := strings.Split(tmpString, ",")
	for _, val := range tmpSlice {
		areaIds = append(areaIds, int64(helper.StringToInt(val)))
	}
	return
}

func (d *Manager) GetAllOrderByPosition(c app.Ctx, enterpriseId int64) ([]dashboardModel.DashboardManagerModel, error) {
	var result []dashboardModel.DashboardManagerModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Order("position_level desc").Find(&result).Error
	return result, err
}

func (d *Manager) GetPagination(c app.Ctx, cond dashboardReq.ManagerAuthList) (list []dashboardResp.ManagerList, total int64, err error) {
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).
		Select("dashboard_manager.id,dashboard_manager.member_id,dashboard_manager.position,dashboard_manager.area_ids,members.name,members.userid").
		Joins("LEFT JOIN micro_scrm.members ON dashboard_manager.member_id = members.id").
		Where("dashboard_manager.deleted_at = 0 AND members.deleted_at = 0").
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
	err = query.Order("dashboard_manager.created_at DESC").Limit(cond.Size).Offset(offset).Find(&list).Error

	return list, total, err
}
