package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type MedalStyle struct {
	dao   dao.Dao
	model *excitationModels.ExcitationMedalStyleModel
}

func NewMedalStyle() *MedalStyle {
	return &MedalStyle{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationMedalStyleModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *MedalStyle) WithTx(tx *gorm.DB) *MedalStyle {
	d.dao.DB = tx
	return d
}

func (d *MedalStyle) GetOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationMedalStyleModel, error) {
	var result excitationModels.ExcitationMedalStyleModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MedalStyle) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]excitationModels.ExcitationMedalStyleModel, error) {
	var result []excitationModels.ExcitationMedalStyleModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MedalStyle) GetManyUnscoped(c app.Ctx, query interface{}, args ...interface{}) ([]excitationModels.ExcitationMedalStyleModel, error) {
	var result []excitationModels.ExcitationMedalStyleModel
	err := d.dao.WithContext(c).Model(d.model).Unscoped().Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MedalStyle) CreateData(c app.Ctx, data excitationModels.ExcitationMedalStyleModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *MedalStyle) CreateInBatchesData(c app.Ctx, data []excitationModels.ExcitationMedalStyleModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *MedalStyle) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *MedalStyle) DeleteData(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationMedalStyleModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d *MedalStyle) ForceDelete(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationMedalStyleModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}

func (d *MedalStyle) GetOneByNumber(c app.Ctx, enterpriseId int64, number string) (excitationModels.ExcitationMedalStyleModel, error) {
	var result excitationModels.ExcitationMedalStyleModel
	err := d.dao.WithContext(c).Model(d.model).Where("enterpriseId = ? and number = ?", enterpriseId, number).Find(&result).Error
	return result, err
}

func (d *MedalStyle) GetManyByMedalNumber(c app.Ctx, enterpriseId int64, medalNumber []string) ([]excitationModels.ExcitationMedalStyleModel, error) {
	var result []excitationModels.ExcitationMedalStyleModel
	err := d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ? and medal_number in ?", enterpriseId, medalNumber).Find(&result).Error
	return result, err
}
