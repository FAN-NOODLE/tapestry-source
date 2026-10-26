package tagsDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	tagModels "micro-scrm/internal/app/models/work"
	"micro-scrm/internal/pkg/db/gormx"
)

type WorkExternalTag struct {
	dao   dao.Dao
	model *tagModels.WorkExternalTagModel
}

func NewWorkExternalTag() *WorkExternalTag {
	return &WorkExternalTag{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(tagModels.WorkExternalTagModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *WorkExternalTag) WithTx(tx *gorm.DB) *WorkExternalTag {
	d.dao.DB = tx
	return d
}

func (d *WorkExternalTag) GetOne(c app.Ctx, query interface{}, args ...interface{}) (tagModels.WorkExternalTagModel, error) {
	var result tagModels.WorkExternalTagModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *WorkExternalTag) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]tagModels.WorkExternalTagModel, error) {
	var result []tagModels.WorkExternalTagModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *WorkExternalTag) CreateData(c app.Ctx, data tagModels.WorkExternalTagModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *WorkExternalTag) CreateInBatchesData(c app.Ctx, data []tagModels.WorkExternalTagModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *WorkExternalTag) SaveData(c app.Ctx, data tagModels.WorkExternalTagModel) error {
	err := d.dao.WithContext(c).Model(d.model).Save(&data).Error
	return err
}

func (d *WorkExternalTag) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *WorkExternalTag) DeleteData(c app.Ctx, query interface{}) error {
	var model tagModels.WorkExternalTagModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d *WorkExternalTag) ForceDelete(c app.Ctx, query interface{}) error {
	var model tagModels.WorkExternalTagModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}

func (d *WorkExternalTag) GetOneById(c app.Ctx, enterpriseId, id int64) (tagModels.WorkExternalTagModel, error) {
	return d.GetOne(c, "enterprise_id = ? and id = ?", enterpriseId, id)
}
