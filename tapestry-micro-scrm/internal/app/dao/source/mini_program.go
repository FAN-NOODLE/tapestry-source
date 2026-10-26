package source

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	sourceModels "micro-scrm/internal/app/models/source"
	"micro-scrm/internal/pkg/db/gormx"
)

type MiniProgramDao struct {
	dao       dao.Dao
	tableName string
}

func NewMiniProgramDao() *MiniProgramDao {
	return &MiniProgramDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: sourceModels.MiniProgramModel{}.TableName(),
	}
}

func (d *MiniProgramDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.MiniProgramModel, error) {
	var result []sourceModels.MiniProgramModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MiniProgramDao) GetManyPreload(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.MiniProgramModel, error) {
	var result []sourceModels.MiniProgramModel
	err := d.dao.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("enterprise_id = ?", enterpriseId)
	}).Preload("Source").Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MiniProgramDao) Update(c app.Ctx, query, column map[string]interface{}) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).Updates(column).Error
}
