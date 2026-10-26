package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type IntegralRule struct {
	dao   dao.Dao
	model *excitationModels.ExcitationIntegralRuleModel
}

// WithTx 切换db连接,用于事物

func (d *IntegralRule) WithTx(tx *gorm.DB) *IntegralRule {
	d.dao.DB = tx
	return d
}

func NewIntegralRule() *IntegralRule {
	return &IntegralRule{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationIntegralRuleModel),
	}
}

func (d *IntegralRule) GetList(c app.Ctx, cond excitationReq.RuleList) ([]excitationModels.ExcitationIntegralRuleModel, int64, error) {
	var (
		total int64
		resp  []excitationModels.ExcitationIntegralRuleModel
		err   error
	)

	query := d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ?", cond.EnterpriseId).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if cond.Type > 0 {
				db.Where("type = ?", cond.Type)
			}
			if cond.RuleName != "" {
				db.Where("rule_name like ?", "%"+cond.RuleName+"%")
			}

			if cond.StartTime > 0 && cond.EndTime > 0 {
				db.Where("created_at between ? and ?", cond.StartTime, cond.EndTime)
			}
			return db
		})

	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size
	query.Preload("IncentiveConfig", "enterprise_id=?", cond.EnterpriseId).Preload("IncentiveRelation", "enterprise_id=?", cond.EnterpriseId)
	err = query.Order("created_at DESC").Limit(cond.Size).Offset(offset).Find(&resp).Error

	return resp, total, err
}

func (d *IntegralRule) CreateData(c app.Ctx, data excitationModels.ExcitationIntegralRuleModel) (error, int64) {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err, data.ID
}

func (d *IntegralRule) GetOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationIntegralRuleModel, error) {
	var result excitationModels.ExcitationIntegralRuleModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Preload("Details").Find(&result).Error
	return result, err
}

func (d *IntegralRule) DeleteData(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationIntegralRuleModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d *IntegralRule) UpdateData(c app.Ctx, query interface{}, data excitationModels.ExcitationIntegralRuleModel) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *IntegralRule) GetMany(c app.Ctx, cond excitationReq.RuleMany) ([]excitationModels.ExcitationIntegralRuleModel, int64, error) {
	var (
		total int64
		resp  []excitationModels.ExcitationIntegralRuleModel
		err   error
	)

	query := d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ?", cond.EnterpriseId).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if len(cond.Types) > 0 {
				db.Where("type in (?)", cond.Types)
			}
			return db
		})

	err = query.Order("created_at DESC").Find(&resp).Error

	return resp, total, err
}
