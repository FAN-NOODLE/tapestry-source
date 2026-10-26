package tagsDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	tagModels "micro-scrm/internal/app/models/work"
	"micro-scrm/internal/pkg/db/gormx"
)

type WorkExternalTagGroup struct {
	dao   dao.Dao
	model *tagModels.WorkExternalTagGroupModel
}

func NewWorkExternalTagGroup() *WorkExternalTagGroup {
	return &WorkExternalTagGroup{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(tagModels.WorkExternalTagGroupModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *WorkExternalTagGroup) WithTx(tx *gorm.DB) *WorkExternalTagGroup {
	d.dao.DB = tx
	return d
}

func (d *WorkExternalTagGroup) GetOne(c app.Ctx, query interface{}, args ...interface{}) (tagModels.WorkExternalTagGroupModel, error) {
	var result tagModels.WorkExternalTagGroupModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *WorkExternalTagGroup) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]tagModels.WorkExternalTagGroupModel, error) {
	var result []tagModels.WorkExternalTagGroupModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *WorkExternalTagGroup) CreateData(c app.Ctx, data tagModels.WorkExternalTagGroupModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *WorkExternalTagGroup) CreateInBatchesData(c app.Ctx, data []tagModels.WorkExternalTagGroupModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *WorkExternalTagGroup) SaveData(c app.Ctx, data tagModels.WorkExternalTagGroupModel) error {
	err := d.dao.WithContext(c).Model(d.model).Save(&data).Error
	return err
}

func (d *WorkExternalTagGroup) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *WorkExternalTagGroup) DeleteData(c app.Ctx, query interface{}) error {
	var model tagModels.WorkExternalTagGroupModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d *WorkExternalTagGroup) ForceDelete(c app.Ctx, query interface{}) error {
	var model tagModels.WorkExternalTagGroupModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}
