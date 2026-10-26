package notificationDao

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	notificationReq "micro-scrm/internal/app/http/request/notification"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"strings"
)

type CrowdDataDao struct {
	dao       dao.Dao
	model     *microScrmModel.NotificationCrowdDataModel
	tableName string
}

func NewCrowdDataDao() *CrowdDataDao {
	return &CrowdDataDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(microScrmModel.NotificationCrowdDataModel),
		tableName: microScrmModel.NotificationCrowdDataModel{}.TableName(),
	}
}

func (d CrowdDataDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.NotificationCrowdDataModel, error) {
	var result microScrmModel.NotificationCrowdDataModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Preload("DataList", func(db *gorm.DB) *gorm.DB {
		return db.Scopes(d.dao.BranchTable(microScrmModel.NotificationCrowdDataModel{}.TableName(), enterpriseId))
	}).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d CrowdDataDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.NotificationCrowdDataModel, error) {
	var result []microScrmModel.NotificationCrowdDataModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d CrowdDataDao) GetPagination(c app.Ctx, cond notificationReq.CrowdListReq) ([]microScrmModel.NotificationCrowdDataModel, int64, error) {
	var (
		total int64
		resp  []microScrmModel.NotificationCrowdDataModel
		err   error
	)

	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).
		Scopes(func(db *gorm.DB) *gorm.DB {
			// 规则名称
			if cond.Name != "" {
				db.Where("name like ?", "%"+strings.Trim(cond.Name, " ")+"%")
			}
			// 获得时间
			if cond.StartTime > 0 && cond.EndTime > 0 {
				db.Where("created_at between ? and ?", cond.StartTime, cond.EndTime)
			}
			return db
		})
	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size

	err = query.Order("created_at DESC").Order("id DESC").Limit(cond.Size).Offset(offset).Find(&resp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return resp, total, err
}

func (d CrowdDataDao) CreateData(c app.Ctx, enterpriseId int64, data microScrmModel.NotificationCrowdDataModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Create(&data).Error
	return err
}

func (d CrowdDataDao) CreateInBatchesData(c app.Ctx, enterpriseId int64, data []microScrmModel.NotificationCrowdDataModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d CrowdDataDao) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d CrowdDataDao) DeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model microScrmModel.NotificationCrowdDataModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).Delete(&model).Error
	return err
}
func (d CrowdDataDao) UnscopedDeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model microScrmModel.NotificationCrowdDataModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}
