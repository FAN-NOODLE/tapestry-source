package microScrmDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	workExternalModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type RepeatCustomer struct {
	dao       dao.Dao
	tableName string
}

func NewRepeatCustomer() *RepeatCustomer {
	return &RepeatCustomer{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: workExternalModel.WxUserCorpExternalRepeatCustomerModel{}.TableName(),
	}
}

func (d RepeatCustomer) GetRepeatCustomerInfo(c app.Ctx, enterpriseId int64, page int, size int) (data []workExternalModel.WxUserCorpExternalRepeatCustomerModel, err error) {
	offset := page * size
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Offset(offset).Limit(size).Find(&data).Error
	return
}

func (d RepeatCustomer) SyncRepeatCustomerTotal(c app.Ctx, enterpriseId int64) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Count(&total).Error
	return
}

func (d RepeatCustomer) DeleteRepeatCustomerByExternalIds(c app.Ctx, enterpriseId int64, ids []int64) error {
	var model workExternalModel.WxUserCorpExternalRepeatCustomerModel
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("work_external_user_id in ?", ids).Delete(&model).Error
}
