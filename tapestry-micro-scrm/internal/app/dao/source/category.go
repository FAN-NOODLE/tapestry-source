package source

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/source"
	"micro-scrm/internal/pkg/db/gormx"
)

type CategoryDao struct {
	dao       dao.Dao
	tableName string
}

func NewCategoryDao() *CategoryDao {
	return &CategoryDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: source.CategoryModel{}.TableName(),
	}
}

func (d *CategoryDao) GetList(c app.Ctx, name string, page, size int, enterpriseId int64, ids []int64) (error, []source.CategoryModel, int64) {
	query := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id", enterpriseId)
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	query.Scopes(func(q *gorm.DB) *gorm.DB {
		if len(ids) > 0 {
			return query.Where("id in ?", ids)
		}
		return query
	})
	var (
		total        int64
		categoryList []source.CategoryModel
		err          error
	)

	err = query.Count(&total).Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&categoryList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return err, categoryList, total
}

func (d *CategoryDao) GetOneByName(c app.Ctx, name string, enterpriseId int64) (error, source.CategoryModel) {
	var result source.CategoryModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"name":          name,
		"enterprise_id": enterpriseId,
	}).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return err, result
}

func (d *CategoryDao) GetOneById(c app.Ctx, id, enterpriseId int64) (error, source.CategoryModel) {
	var result source.CategoryModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"id":            id,
		"enterprise_id": enterpriseId,
	}).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return err, result
}

func (d *CategoryDao) Create(c app.Ctx, data source.CategoryModel) error {
	return d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
}

func (d *CategoryDao) Update(c app.Ctx, query, column map[string]interface{}) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).Updates(&column).Error
}

func (d *CategoryDao) Delete(c app.Ctx, query map[string]interface{}) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(query).Delete(source.CategoryModel{}).Error
}
