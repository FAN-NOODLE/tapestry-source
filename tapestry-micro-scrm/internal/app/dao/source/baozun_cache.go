package source

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	sourceModels "micro-scrm/internal/app/models/source"
	"micro-scrm/internal/pkg/db/gormx"
)

type BaozunCacheDao struct {
	dao       dao.Dao
	model     *sourceModels.BaozunCacheModel
	tableName string
}

func NewBaozunCacheDao() *BaozunCacheDao {
	return &BaozunCacheDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(sourceModels.BaozunCacheModel),
		tableName: sourceModels.BaozunCacheModel{}.TableName(),
	}
}

func (d *BaozunCacheDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (sourceModels.BaozunCacheModel, error) {
	var result sourceModels.BaozunCacheModel
	err := d.dao.WithContext(c).Model(d.model).Where("enterprise_id=?", enterpriseId).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *BaozunCacheDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]sourceModels.BaozunCacheModel, error) {
	var result []sourceModels.BaozunCacheModel
	err := d.dao.WithContext(c).Model(d.model).Where("enterprise_id=?", enterpriseId).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *BaozunCacheDao) CreateData(c app.Ctx, data sourceModels.BaozunCacheModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *BaozunCacheDao) CreateInBatchesData(c app.Ctx, data []sourceModels.BaozunCacheModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *BaozunCacheDao) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *BaozunCacheDao) Increment(c app.Ctx, query interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s + ?", column), num)).Error
	return err
}
