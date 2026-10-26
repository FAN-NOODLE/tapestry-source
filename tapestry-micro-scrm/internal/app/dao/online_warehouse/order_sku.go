package onlineWarehouseDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	onlineWarehouseReq "micro-scrm/internal/app/http/request/online_warehouse"
	onlineWarehouseResp "micro-scrm/internal/app/http/response/online_warehouse"
	onlineWarehouseModel "micro-scrm/internal/app/models/micro_scrm/online_warehouse"
)

type OrderSkuDao struct {
	dao.EnterpriseRepo[onlineWarehouseModel.OnlineWarehouseOrderSkuModel]
	model     *onlineWarehouseModel.OnlineWarehouseOrderSkuModel
	tableName string
}

func NewOrderSkuDao() *OrderSkuDao {
	return &OrderSkuDao{
		EnterpriseRepo: dao.EnterpriseRepo[onlineWarehouseModel.OnlineWarehouseOrderSkuModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(onlineWarehouseModel.OnlineWarehouseOrderSkuModel),
		tableName: onlineWarehouseModel.OnlineWarehouseOrderSkuModel{}.TableName(),
	}
}

// UpdateData 重写该方法更新update_at字段
func (d *OrderSkuDao) UpdateData(c ctx.Context, enterpriseId int64, query any, data any) error {
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Model(d.model).Where(query).Updates(data).Error
	return err
}

func (d *OrderSkuDao) GetDashboardNum(c ctx.Context, enterpriseId int64, orderIds []int64) ([]onlineWarehouseResp.TipsNums, error) {
	var result []onlineWarehouseResp.TipsNums
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Model(d.model).Select("status,sum(quantity) as num").Where(map[string]interface{}{
		"order_id": orderIds,
	}).Group("status").Scan(&result).Error
	return result, err
}

func (d *OrderSkuDao) GoodsRank(c ctx.Context, cond onlineWarehouseReq.DashboardReq) ([]onlineWarehouseResp.GoodsRankResp, error) {
	result := make([]onlineWarehouseResp.GoodsRankResp, 0)
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).
		Select("sku_name,style_color as sku_info,max(image_url) as image_url,dept,sum(quantity) as num").Scopes(func(db *gorm.DB) *gorm.DB {
		if cond.StartTime > 0 && cond.EndTime > 0 {
			db.Where("created_at between ? and ?", cond.StartTime, cond.EndTime)
		}
		return db
	}).Where("status=?", onlineWarehouseModel.OrderStatusShipped).Group("style_color").Order("num desc,image_url desc").Limit(100).Scan(&result).Error
	return result, err
}
