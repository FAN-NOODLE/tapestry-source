package microScrmDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type TaskConditionDataSaleDao struct {
	dao       dao.Dao
	tableName string
	model     *microScrmModel.TaskConditionDataSaleModel
}

func NewTaskConditionDataSaleDao() *TaskConditionDataSaleDao {
	return &TaskConditionDataSaleDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.TaskConditionDataSaleModel{}.TableName(),
		model:     new(microScrmModel.TaskConditionDataSaleModel),
	}
}

func (d *TaskConditionDataSaleDao) CreateData(c app.Ctx, enterpriseId int64, data microScrmModel.TaskConditionDataSaleModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Create(&data).Error
	return err
}

func (d *TaskConditionDataSaleDao) CreateGoodsData(c app.Ctx, enterpriseId int64, data microScrmModel.TaskConditionDataOrderGoodsModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.TaskConditionDataOrderGoodsModel{}.TableName(), enterpriseId)).Model(&microScrmModel.TaskConditionDataOrderGoodsModel{}).Create(&data).Error
	return err
}

func (d *TaskConditionDataSaleDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.TaskConditionDataSaleModel, error) {
	var result []microScrmModel.TaskConditionDataSaleModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

// GetSalesInfo 获取销售单数据
func (d *TaskConditionDataSaleDao) GetSalesInfo(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.TaskConditionDataSaleModel, error) {
	var data microScrmModel.TaskConditionDataSaleModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&data).Error
	return data, err
}

// UpdateSaleData 更新销售单数据
func (d *TaskConditionDataSaleDao) UpdateSaleData(c app.Ctx, enterpriseId, id int64, data microScrmModel.TaskConditionDataSaleModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(map[string]interface{}{
		"id": id,
	}).UpdateColumns(&data).Error
	return err
}

// GetOrderGoodsData 获取订单商品数据
func (d *TaskConditionDataSaleDao) GetOrderGoodsData(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.TaskConditionDataOrderGoodsModel, error) {
	var data microScrmModel.TaskConditionDataOrderGoodsModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.TaskConditionDataOrderGoodsModel{}.TableName(), enterpriseId)).Model(&microScrmModel.TaskConditionDataOrderGoodsModel{}).Where(query, args...).Find(&data).Error
	return data, err
}

// UpdateOrderGoodsData 更新销售商品数据
func (d *TaskConditionDataSaleDao) UpdateOrderGoodsData(c app.Ctx, enterpriseId, id int64, data microScrmModel.TaskConditionDataOrderGoodsModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.TaskConditionDataOrderGoodsModel{}.TableName(), enterpriseId)).Model(&microScrmModel.TaskConditionDataOrderGoodsModel{}).Where(map[string]interface{}{
		"id": id,
	}).UpdateColumns(&data).Error
	return err
}
