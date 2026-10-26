package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type IncentiveRelation struct {
	dao   dao.Dao
	model *excitationModels.ExcitationIncentiveRelationModel
}

func NewIncentiveRelation() *IncentiveRelation {
	return &IncentiveRelation{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationIncentiveRelationModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *IncentiveRelation) WithTx(tx *gorm.DB) *IncentiveRelation {
	d.dao.DB = tx
	return d
}

func (d *IncentiveRelation) Create(c app.Ctx, data excitationModels.ExcitationIncentiveRelationModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *IncentiveRelation) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *IncentiveRelation) ForceDelete(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationIncentiveRelationModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}

func (d *IncentiveRelation) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]excitationModels.ExcitationIncentiveRelationModel, error) {
	var result []excitationModels.ExcitationIncentiveRelationModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *IncentiveRelation) GetOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationIncentiveRelationModel, error) {
	var result excitationModels.ExcitationIncentiveRelationModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).First(&result).Error
	return result, err
}

func (d *IncentiveRelation) GetIds(c app.Ctx, query interface{}, args ...interface{}) []int64 {
	var result []int64
	d.dao.WithContext(c).Model(d.model).Where(query, args...).Pluck("id", &result)
	return result
}

func (d *IncentiveRelation) GetDetailList(c app.Ctx, enterpriseId int64, ruleConditions []int, query interface{}, args ...interface{}) ([]excitationModels.ExcitationIncentiveRelationModel, error) {
	var (
		resp []excitationModels.ExcitationIncentiveRelationModel
		err  error
	)

	err = d.dao.WithContext(c).Model(d.model).
		Where(query, args...).
		Preload("IncentiveRelationRules", "enterprise_id=? and rule_condition in (?) ", enterpriseId, ruleConditions).
		Order("incentive_config_type ASC").Find(&resp).Error

	return resp, err
}
