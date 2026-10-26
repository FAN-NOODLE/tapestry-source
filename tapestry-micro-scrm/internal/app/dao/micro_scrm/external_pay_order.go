package microScrmDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	openReq "micro-scrm/internal/app/http/request/open"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type ExternalPayOrderDao struct {
	dao       dao.Dao
	tableName string
}

func NewExternalPayOrderDao() *ExternalPayOrderDao {
	return &ExternalPayOrderDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.ExternalPayOrderModel{}.TableName(),
	}
}

func (d *ExternalPayOrderDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.ExternalPayOrderModel, error) {
	var result []microScrmModel.ExternalPayOrderModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *ExternalPayOrderDao) GetOrderMemberUserId(c app.Ctx, OrderNo string, EnterpriseID int64) (openReq.SendWaitPaymentOrderRes, error) {
	var result openReq.SendWaitPaymentOrderRes
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, EnterpriseID)).
		Where("out_trade_no", OrderNo).
		First(&result).Error

	return result, err
}
