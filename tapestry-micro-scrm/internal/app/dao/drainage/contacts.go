package drainageContactDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	drainageContactGroupReq "micro-scrm/internal/app/http/request/drainage_contact"
	"micro-scrm/internal/app/models/drainage"
	"micro-scrm/internal/pkg/db/gormx"
)

type DrainageContactsDao struct {
	dao       dao.Dao
	tableName string
	model     *drainageContactModel.DrainageContactsModel
}

func NewDrainageContactsDao() *DrainageContactsDao {
	return &DrainageContactsDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: drainageContactModel.DrainageContactsModel{}.TableName(),
		model:     new(drainageContactModel.DrainageContactsModel),
	}
}

func (d *DrainageContactsDao) GetOne(c app.Ctx, query interface{}, args ...interface{}) (drainageContactModel.DrainageContactsModel, error) {
	var result drainageContactModel.DrainageContactsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *DrainageContactsDao) GetFirstByGroupId(c app.Ctx, enterpriseId, groupId int64) (drainageContactModel.DrainageContactsModel, error) {
	var result drainageContactModel.DrainageContactsModel
	err := d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ? AND contact_group_id = ? AND deleted_at = 0", enterpriseId, groupId).
		Order("id asc").
		Limit(1).
		Find(&result).Error
	return result, err
}

func (d *DrainageContactsDao) GetContactListByGroupId(c app.Ctx, enterpriseId, groupId int64) ([]drainageContactModel.DrainageContactsModel, error) {
	var result []drainageContactModel.DrainageContactsModel
	err := d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ? AND contact_group_id = ?", enterpriseId, groupId).
		Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at = 0").Order("id asc")
		}).
		Order("id asc").Find(&result).Error
	return result, err
}

func (d *DrainageContactsDao) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]drainageContactModel.DrainageContactsModel, error) {
	var result []drainageContactModel.DrainageContactsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *DrainageContactsDao) CreateData(c app.Ctx, data drainageContactModel.DrainageContactsModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *DrainageContactsDao) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *DrainageContactsDao) CreateInBatchesData(c app.Ctx, data []drainageContactModel.DrainageContactsModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *DrainageContactsDao) GetPage(c ctx.Context, req drainageContactGroupReq.DetailListReq) ([]drainageContactModel.DrainageContactsModel, int64, error) {
	var (
		result []drainageContactModel.DrainageContactsModel
		total  int64
	)
	offset := (req.Page - 1) * req.Size

	err := d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ?", req.EnterpriseId).
		Where("contact_group_id = ?", req.ID).
		Count(&total).
		Order("is_contact_group_single asc,id asc").
		Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at = 0").Order("id asc")
		}).
		Offset(offset).Limit(req.Size).
		Find(&result).Error
	return result, total, err
}

func (d *DrainageContactsDao) GetManyPreloadThirdMembers(c app.Ctx, query interface{}, args ...interface{}) ([]drainageContactModel.DrainageContactsModel, error) {
	var result []drainageContactModel.DrainageContactsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).
		Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at = 0")
		}).
		Find(&result).Error
	return result, err
}
