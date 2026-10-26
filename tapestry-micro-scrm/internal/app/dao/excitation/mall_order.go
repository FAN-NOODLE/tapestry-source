package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type MallOrder struct {
	dao   dao.Dao
	model *excitationModels.ExcitationMallOrderModel
}

func NewMallOrder() *MallOrder {
	return &MallOrder{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationMallOrderModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *MallOrder) WithTx(tx *gorm.DB) *MallOrder {
	d.dao.DB = tx
	return d
}

func (d *MallOrder) GetOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationMallOrderModel, error) {
	var result excitationModels.ExcitationMallOrderModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MallOrder) GetMany(c app.Ctx, query interface{}) ([]excitationModels.ExcitationMallOrderModel, error) {
	var result []excitationModels.ExcitationMallOrderModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Find(&result).Error
	return result, err
}

func (d *MallOrder) CreateData(c app.Ctx, data excitationModels.ExcitationMallOrderModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *MallOrder) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}
