package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type Achievement struct {
	dao   dao.Dao
	model *excitationModels.ExcitationAchievementModel
}

func NewAchievement() *Achievement {
	return &Achievement{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationAchievementModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *Achievement) WithTx(tx *gorm.DB) *Achievement {
	d.dao.DB = tx
	return d
}

func (d *Achievement) GetOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationAchievementModel, error) {
	var result excitationModels.ExcitationAchievementModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Achievement) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]excitationModels.ExcitationAchievementModel, error) {
	var result []excitationModels.ExcitationAchievementModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Achievement) CreateData(c app.Ctx, data excitationModels.ExcitationAchievementModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *Achievement) CreateInBatchesData(c app.Ctx, data []excitationModels.ExcitationAchievementModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *Achievement) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *Achievement) DeleteData(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationAchievementModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}
