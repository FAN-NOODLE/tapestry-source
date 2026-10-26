package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type IncentiveRelationRule struct {
	dao   dao.Dao
	model *excitationModels.ExcitationIncentiveRelationRuleModel
}

func NewIncentiveRelationRule() *IncentiveRelationRule {
	return &IncentiveRelationRule{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationIncentiveRelationRuleModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *IncentiveRelationRule) WithTx(tx *gorm.DB) *IncentiveRelationRule {
	d.dao.DB = tx
	return d
}

func (d *IncentiveRelationRule) CreateInBatchesData(c app.Ctx, data []excitationModels.ExcitationIncentiveRelationRuleModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *IncentiveRelationRule) ForceDelete(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationIncentiveRelationRuleModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}

func (d *IncentiveRelationRule) DeleteData(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationIncentiveRelationRuleModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d *IncentiveRelationRule) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]excitationModels.ExcitationIncentiveRelationRuleModel, error) {
	var result []excitationModels.ExcitationIncentiveRelationRuleModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Where("id > 0 ").Find(&result).Error
	return result, err
}
