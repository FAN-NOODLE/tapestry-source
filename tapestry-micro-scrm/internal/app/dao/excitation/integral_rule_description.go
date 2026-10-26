package excitationDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type IntegralRuleDescription struct {
	dao   dao.Dao
	model *excitationModels.ExcitationIntegralRuleDescriptionModel
}

func NewIntegralRuleDescription() *IntegralRuleDescription {
	return &IntegralRuleDescription{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationIntegralRuleDescriptionModel),
	}
}

func (d *IntegralRuleDescription) CreateOUpdate(c app.Ctx, data excitationModels.ExcitationIntegralRuleDescriptionModel) error {
	var (
		count int64
		err   error
	)
	d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ?", data.EnterpriseID).Count(&count)

	if count > 0 {
		err = d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ?", data.EnterpriseID).UpdateColumns(&data).Error
	} else {
		err = d.dao.WithContext(c).Model(d.model).Create(&data).Error
	}

	return err
}

func (d *IntegralRuleDescription) Detail(c app.Ctx, params excitationReq.DescriptionReq) (excitationResp.DescriptionResp, error) {
	var result excitationResp.DescriptionResp
	err := d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ?", params.EnterpriseId).Find(&result).Error

	return result, err
}
