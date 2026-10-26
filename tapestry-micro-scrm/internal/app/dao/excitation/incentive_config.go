package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type IncentiveConfig struct {
	dao   dao.Dao
	model *excitationModels.ExcitationIncentiveConfigModel
}

func NewIncentiveConfig() *IncentiveConfig {
	return &IncentiveConfig{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationIncentiveConfigModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *IncentiveConfig) WithTx(tx *gorm.DB) *IncentiveConfig {
	d.dao.DB = tx
	return d
}

func (d *IncentiveConfig) GetMany(c app.Ctx, cond excitationReq.IncentiveConfigListReq) ([]excitationModels.ExcitationIncentiveConfigModel, int64, error) {
	var (
		total int64
		resp  []excitationModels.ExcitationIncentiveConfigModel
		err   error
	)

	where := map[string]interface{}{
		"enterprise_id": cond.EnterpriseId,
	}
	query := d.dao.WithContext(c).Model(d.model).
		Where(where)
	query.Count(&total)

	err = query.Preload("IncentiveRelation", where).
		Order("created_at ASC").Find(&resp).Error

	return resp, total, err
}

func (d *IncentiveConfig) CreateInBatchesData(c app.Ctx, data []excitationModels.ExcitationIncentiveConfigModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}
