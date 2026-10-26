package source

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	sourceModels "micro-scrm/internal/app/models/source"
	"micro-scrm/internal/pkg/db/gormx"
)

type FileDao struct {
	dao       dao.Dao
	tableName string
}

func NewFileDao() *FileDao {
	return &FileDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: sourceModels.FileModel{}.TableName(),
	}
}

func (d *FileDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.FileModel, error) {
	var result []sourceModels.FileModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *FileDao) GetManyPreload(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.FileModel, error) {
	var result []sourceModels.FileModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Preload("Source").Where(query, args...).Find(&result).Error
	return result, err
}

func (d *FileDao) Update(c app.Ctx, query, column map[string]interface{}) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).Updates(column).Error
}

func (d *FileDao) GetManyLimit(c app.Ctx, enterpriseId int64, limit int, query interface{}, args ...interface{}) ([]sourceModels.ImageModel, error) {
	var result []sourceModels.ImageModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Where(query, args...).Order("id asc").Limit(limit).Find(&result).Error
	return result, err
}
