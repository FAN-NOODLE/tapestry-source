package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type MedalTaskRecord struct {
	dao   dao.Dao
	model *excitationModels.ExcitationMedalTaskRecordModel
}

func NewMedalTaskRecord() *MedalTaskRecord {
	return &MedalTaskRecord{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationMedalTaskRecordModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *MedalTaskRecord) WithTx(tx *gorm.DB) *MedalTaskRecord {
	d.dao.DB = tx
	return d
}

func (d *MedalTaskRecord) GetOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationMedalTaskRecordModel, error) {
	var result excitationModels.ExcitationMedalTaskRecordModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MedalTaskRecord) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]excitationModels.ExcitationMedalTaskRecordModel, error) {
	var result []excitationModels.ExcitationMedalTaskRecordModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MedalTaskRecord) CreateData(c app.Ctx, data excitationModels.ExcitationMedalTaskRecordModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *MedalTaskRecord) SaveData(c app.Ctx, data excitationModels.ExcitationMedalTaskRecordModel) error {
	err := d.dao.WithContext(c).Model(d.model).Save(&data).Error
	return err
}

func (d *MedalTaskRecord) UpdateData(c app.Ctx, query interface{}, data interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *MedalTaskRecord) DeleteData(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationMedalTaskRecordModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d *MedalTaskRecord) ForceDelete(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationMedalTaskRecordModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}

func (d *MedalTaskRecord) UpdateStatus(c app.Ctx, enterpriseId, id int64, status int64) error {
	return d.UpdateData(c, map[string]any{
		"enterprise_id": enterpriseId,
		"medal_id":      id,
	}, map[string]any{"status": status})
}
