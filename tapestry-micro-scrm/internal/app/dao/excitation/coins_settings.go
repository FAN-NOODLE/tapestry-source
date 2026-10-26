package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type CoinsSettings struct {
	dao   dao.Dao
	model *excitationModels.ExcitationCoinsSettingsModel
}

func NewCoinsSettings() *CoinsSettings {
	return &CoinsSettings{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationCoinsSettingsModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *CoinsSettings) WithTx(tx *gorm.DB) *CoinsSettings {
	d.dao.DB = tx
	return d
}

func (d *CoinsSettings) GetOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationCoinsSettingsModel, error) {
	var result excitationModels.ExcitationCoinsSettingsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *CoinsSettings) GetMany(c app.Ctx, query interface{}) ([]excitationModels.ExcitationCoinsSettingsModel, error) {
	var result []excitationModels.ExcitationCoinsSettingsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Find(&result).Error
	return result, err
}

func (d *CoinsSettings) CreateData(c app.Ctx, data excitationModels.ExcitationCoinsSettingsModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *CoinsSettings) CreateInBatchesData(c app.Ctx, data []excitationModels.ExcitationCoinsSettingsModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *CoinsSettings) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *CoinsSettings) DeleteData(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationCoinsSettingsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d *CoinsSettings) ForceDelete(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationCoinsSettingsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}
