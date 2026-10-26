package microScrmDao

import (
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	memberReq "micro-scrm/internal/app/http/request/member"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type ExternalEverydayStatisticsDao struct {
	dao       dao.Dao
	model     microScrmModel.ExternalEverydayStatisticsModel
	tableName string
}

func NewExternalEverydayStatisticsDao() *ExternalEverydayStatisticsDao {
	return &ExternalEverydayStatisticsDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.ExternalEverydayStatisticsModel{}.TableName(),
	}
}

func (d *ExternalEverydayStatisticsDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.ExternalEverydayStatisticsModel, error) {
	var result microScrmModel.ExternalEverydayStatisticsModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *ExternalEverydayStatisticsDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.ExternalEverydayStatisticsModel, error) {
	var result []microScrmModel.ExternalEverydayStatisticsModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}
func (d *ExternalEverydayStatisticsDao) GetLastData(c app.Ctx, enterpriseId int64) ([]microScrmModel.ExternalEverydayStatisticsModel, error) {
	var result []microScrmModel.ExternalEverydayStatisticsModel
	RawSql := fmt.Sprintf("SELECT ees.* FROM micro_scrm_enterprise_%d.external_everyday_statistics ees INNER JOIN "+
		"(SELECT member_id, MAX(id) AS max_id FROM micro_scrm_enterprise_%d.external_everyday_statistics GROUP BY member_id) ees2"+
		" ON ees.member_id = ees2.member_id AND ees.id = ees2.max_id", enterpriseId, enterpriseId)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Raw(RawSql).Scan(&result).Error
	return result, err
}

func (d *ExternalEverydayStatisticsDao) CreateData(c app.Ctx, enterpriseId int64, data microScrmModel.ExternalEverydayStatisticsModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
	return err
}

func (d *ExternalEverydayStatisticsDao) CreateInBatchesData(c app.Ctx, enterpriseId int64, data []microScrmModel.ExternalEverydayStatisticsModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *ExternalEverydayStatisticsDao) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).UpdateColumns(data).Error
	return err
}

func (d *ExternalEverydayStatisticsDao) DeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model microScrmModel.ExternalEverydayStatisticsModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).Delete(&model).Error
	return err
}

func (d *ExternalEverydayStatisticsDao) GetSwCustomers(c app.Ctx, request memberReq.SwCustomersListReq, memberIds []int64) ([]memberReq.SwCustomerBaseRes, int64, error) {
	startTime := request.FinancialTime[0]
	endTime := request.FinancialTime[1]
	offset := (request.Page - 1) * request.Size
	limit := request.Size
	var result []memberReq.SwCustomerBaseRes
	var total int64
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).
		Select("member_id,SUM(`new_customer_change`) as new_customers,SUM(`old_customer_change`) as old_customers,SUM(`sleep_customer_change`) as sleep_customers,SUM(`vip_customer_change`) as vip_customers").
		Where("unix_time BETWEEN ? AND ?", startTime, endTime).Where(map[string]interface{}{
		"member_id": memberIds,
	}).Group("member_id")
	query.Count(&total)
	query = query.Order("new_customers desc,old_customers desc,vip_customers desc").Order("member_id asc")
	err := query.Offset(offset).Limit(limit).Scan(&result).Error
	return result, total, err
}
