package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type IntegralRuleDetail struct {
	dao   dao.Dao
	model *excitationModels.ExcitationIntegralRuleDetailModel
}

// WithTx 切换db连接,用于事物
func (d *IntegralRuleDetail) WithTx(tx *gorm.DB) *IntegralRuleDetail {
	d.dao.DB = tx
	return d
}

func NewIntegralRuleDetail() *IntegralRuleDetail {
	return &IntegralRuleDetail{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationIntegralRuleDetailModel),
	}
}

func (d *IntegralRuleDetail) GetDetail(c app.Ctx, ruleId []int64) ([]excitationModels.ExcitationIntegralRuleDetailModel, error) {
	var (
		resp []excitationModels.ExcitationIntegralRuleDetailModel
		err  error
	)

	err = d.dao.WithContext(c).Model(d.model).Where("rule_id in ?", ruleId).Group("rule_category_id").Find(&resp).Error
	return resp, err
}

func (d *IntegralRuleDetail) CreateData(c app.Ctx, data []excitationModels.ExcitationIntegralRuleDetailModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *IntegralRuleDetail) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]excitationModels.ExcitationIntegralRuleDetailModel, error) {
	var result []excitationModels.ExcitationIntegralRuleDetailModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *IntegralRuleDetail) GetRule(c app.Ctx, query interface{}, args ...interface{}) ([]excitationModels.ExcitationIntegralRuleDetailModel, error) {
	var result []excitationModels.ExcitationIntegralRuleDetailModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Preload("Rule").Find(&result).Error
	return result, err
}

func (d *IntegralRuleDetail) DeleteData(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationIntegralRuleDetailModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d *IntegralRuleDetail) GetOnCategoryList(c app.Ctx, query interface{}, args ...interface{}) ([]excitationModels.ExcitationIntegralRuleDetailModel, error) {
	var result []excitationModels.ExcitationIntegralRuleDetailModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).
		Group("rule_category_id").
		Find(&result).Error
	return result, err
}

func (d *IntegralRuleDetail) GetOne(c app.Ctx, query interface{}) (excitationModels.ExcitationIntegralRuleDetailModel, error) {
	var result excitationModels.ExcitationIntegralRuleDetailModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Find(&result).Error
	return result, err
}
