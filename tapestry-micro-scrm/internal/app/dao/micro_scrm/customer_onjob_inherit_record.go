package microScrmDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type CustomerOnInheritRecordDao struct {
	dao       dao.Dao
	tableName string
}

func NewCustomerOnInheritRecord() *CustomerOnInheritRecordDao {
	return &CustomerOnInheritRecordDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.CustomerOnjobInheritRecordModel{}.TableName(),
	}
}

// GetLimitExternalIds 获取最近90天分配超过两次的external_id
func (d CustomerOnInheritRecordDao) GetLimitExternalIds(c ctx.Context, enterpriseId, startTime int64) ([]microScrmModel.OnInheritExternalIds, error) {
	var data []microScrmModel.OnInheritExternalIds
	err := d.dao.WithContext(c).Model(&microScrmModel.CustomerOnjobInheritRecordModel{}).Select("count(1) as num,external_id").Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"status": []int{
			microScrmModel.TransferSuccess, microScrmModel.TransferRefund, microScrmModel.TransferWait,
		},
	}).Where("created_at > ? ", startTime).Group("external_id").Having("num>=2").Find(&data).Error
	return data, err
}

// GetWaitExternalIds 获取等待客户接替的external_id
func (d CustomerOnInheritRecordDao) GetWaitExternalIds(c ctx.Context, enterpriseId, startTime int64) ([]int64, error) {
	var data []int64
	err := d.dao.WithContext(c).Model(&microScrmModel.CustomerOnjobInheritRecordModel{}).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"status":        microScrmModel.TransferWait,
	}).Where("created_at > ? ", startTime).Group("external_id").Pluck("external_id", &data).Error
	return data, err
}

func (d CustomerOnInheritRecordDao) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]microScrmModel.CustomerOnjobInheritRecordModel, error) {
	var result []microScrmModel.CustomerOnjobInheritRecordModel
	err := d.dao.WithContext(c).Model(&microScrmModel.CustomerOnjobInheritRecordModel{}).Where(query, args...).Find(&result).Error
	return result, err
}
