package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type CoinsExpireRecord struct {
	dao   dao.Dao
	model *excitationModels.ExcitationCoinsExpireRecordModel
}

func NewCoinsExpireRecord() *CoinsExpireRecord {
	return &CoinsExpireRecord{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationCoinsExpireRecordModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *CoinsExpireRecord) WithTx(tx *gorm.DB) *CoinsExpireRecord {
	d.dao.DB = tx
	return d
}

func (d *CoinsExpireRecord) GetMany(c app.Ctx, query interface{}) ([]excitationModels.ExcitationCoinsExpireRecordModel, error) {
	var result []excitationModels.ExcitationCoinsExpireRecordModel
	err := d.dao.WithContext(c).Model(d.model).Order("id desc").Where(query).Find(&result).Error
	return result, err
}

func (d *CoinsExpireRecord) CreateData(c app.Ctx, data excitationModels.ExcitationCoinsExpireRecordModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}
