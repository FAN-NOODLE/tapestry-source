package customerDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/customer"
	"micro-scrm/internal/pkg/db/gormx"
)

type SowTemplateNotice struct {
	dao   dao.Dao
	model *customerModels.SowTemplateNoticeModel
}

func NewSowTemplateNotice() *SowTemplateNotice {
	return &SowTemplateNotice{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(customerModels.SowTemplateNoticeModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *SowTemplateNotice) WithTx(tx *gorm.DB) *SowTemplateNotice {
	d.dao.DB = tx
	return d
}

func (d *SowTemplateNotice) GetOne(c app.Ctx, query interface{}, args ...interface{}) (customerModels.SowTemplateNoticeModel, error) {
	var result customerModels.SowTemplateNoticeModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *SowTemplateNotice) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]customerModels.SowTemplateNoticeModel, error) {
	var result []customerModels.SowTemplateNoticeModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *SowTemplateNotice) GetList(c app.Ctx, page, size int, cond interface{}, option map[string]interface{}) ([]customerModels.SowTemplateNoticeModel, int64, error) {
	var result []customerModels.SowTemplateNoticeModel
	var total int64
	offset := (page - 1) * size
	query := d.dao.WithContext(c).Model(d.model).Where(cond)
	if _, ok := option["start_time"]; ok {
		query.Where("remind_time BETWEEN ? AND ?", option["start_time"], option["end_time"])
	}
	err := query.Count(&total).Offset(offset).Limit(size).
		Order("status asc").
		Order("remind_time desc").
		Find(&result).Error

	return result, total, err
}

func (d *SowTemplateNotice) CreateData(c app.Ctx, data customerModels.SowTemplateNoticeModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *SowTemplateNotice) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).Updates(data).Error
	return err
}

func (d *SowTemplateNotice) DeleteData(c app.Ctx, query interface{}, args ...interface{}) error {
	var model customerModels.SowTemplateNoticeModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Delete(&model).Error
	return err
}

func (d *SowTemplateNotice) UnscopedDeleteData(c app.Ctx, query interface{}) error {
	var model customerModels.SowTemplateNoticeModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}

func (d *SowTemplateNotice) CreateInBatchesData(c app.Ctx, data []customerModels.SowTemplateNoticeModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}
