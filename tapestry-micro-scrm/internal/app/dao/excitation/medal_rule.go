package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type MedalRule struct {
	dao   dao.Dao
	model *excitationModels.ExcitationMedalRuleModel
}

func NewMedalRule() *MedalRule {
	return &MedalRule{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationMedalRuleModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *MedalRule) WithTx(tx *gorm.DB) *MedalRule {
	d.dao.DB = tx
	return d
}

func (d *MedalRule) GetPagination(c app.Ctx, cond excitationReq.MedalRuleList) ([]excitationModels.ExcitationMedalRuleModel, int64, error) {
	var (
		total int64
		resp  []excitationModels.ExcitationMedalRuleModel
		err   error
	)

	query := d.dao.WithContext(c).Model(d.model).Where(map[string]interface{}{
		"enterprise_id": cond.EnterpriseId,
		"pid":           cond.Pid,
		"status":        excitationModels.MedalRuleStatusEnable,
	})
	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size
	err = query.Order("created_at DESC").Limit(cond.Size).Offset(offset).Find(&resp).Error

	return resp, total, err
}

func (d *MedalRule) GetOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationMedalRuleModel, error) {
	var result excitationModels.ExcitationMedalRuleModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MedalRule) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]excitationModels.ExcitationMedalRuleModel, error) {
	var result []excitationModels.ExcitationMedalRuleModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MedalRule) DeleteData(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationMedalRuleModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}
