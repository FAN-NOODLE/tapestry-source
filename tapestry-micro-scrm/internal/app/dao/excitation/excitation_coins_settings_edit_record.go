package excitationDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type CoinsSettingsEditRecord struct {
	dao   dao.Dao
	model *excitationModels.ExcitationCoinsSettingsEditRecordModel
}

func NewCoinsSettingsEditRecord() *CoinsSettingsEditRecord {
	return &CoinsSettingsEditRecord{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationCoinsSettingsEditRecordModel),
	}
}

func (d *CoinsSettingsEditRecord) GetMany(c app.Ctx, query interface{}) ([]excitationModels.ExcitationCoinsSettingsEditRecordModel, error) {
	var result []excitationModels.ExcitationCoinsSettingsEditRecordModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Order("id desc").Find(&result).Error
	return result, err
}

func (d *CoinsSettingsEditRecord) CreateData(c app.Ctx, data excitationModels.ExcitationCoinsSettingsEditRecordModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}
