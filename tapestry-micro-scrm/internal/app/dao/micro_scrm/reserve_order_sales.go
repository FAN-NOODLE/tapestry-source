package microScrmDao

import (
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"

	"gorm.io/gorm"
)

type ReservesOrderSales struct {
	dao       dao.Dao
	tableName string
}

func NewReservesOrderSales() *ReservesOrderSales {
	return &ReservesOrderSales{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.ReserveOrderSalesModel{}.TableName(),
	}
}
func (d *ReservesOrderSales) getScrmTableCk(table string, enterpriseId int64) string {
	return fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, table)
}

// GetReserveSalesCount 获取邀约单的订单数
func (d *ReservesOrderSales) GetReserveSalesCount(c app.Ctx, enterpriseId int64, reserveNumbers []string, startTime, endTime int64) ([]reserveRequest.ReserveSalesNum, error) {
	var result []reserveRequest.ReserveSalesNum
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Select("count(1) as num,reserve_number,sum(order_price) as reserve_order_price").Where(map[string]interface{}{
		"reserve_number": reserveNumbers,
	}).Scopes(func(db *gorm.DB) *gorm.DB {
		if startTime > 0 && endTime > 0 {
			db.Where("pay_time between ? and ? ", startTime, endTime)
		}
		return db
	}).Group("reserve_number").Find(&result).Error
	return result, err
}

// GetReserveSaleGoods 获取邀约单的商品订单列表
func (d *ReservesOrderSales) GetReserveSaleGoods(c app.Ctx, enterpriseId int64, req reserveRequest.ReserveOrderListReq) ([]reserveRequest.ExportOrderRecord, int64, error) {
	page := req.Page
	size := req.Size
	var result []reserveRequest.ExportOrderRecord
	var total int64
	reserveSales := d.getScrmTableCk(d.tableName, enterpriseId)
	reserveGoodsSales := d.getScrmTableCk(microScrmModel.ReserveOrderGoodsModel{}.TableName(), enterpriseId)
	selectStr := `reserve_order_goods.id,reserve_order.reserve_number as reserve_number,reserve_order.reserve_time,store_name,status,user_name,
member_id,member_name,is_open,is_commit,invite_num,is_store,is_pay,order_origin,customer_status,reserve_order_goods.relation_order_id as relation_order_id,goods_name,
style,color,goods_price,reserve_order_goods.pay_time as pay_time_int,sales.store_code as store_code`
	query := gormx.GetCkInstance().WithContext(c).Table(d.getScrmTableCk(microScrmModel.ReserveOrderModel{}.TableName(), enterpriseId)).
		Joins(" LEFT JOIN " + reserveSales + " sales ON sales.reserve_number = reserve_order.reserve_number").
		Joins(" LEFT JOIN " + reserveGoodsSales + " reserve_order_goods ON sales.relation_order_id = reserve_order_goods.relation_order_id").Scopes(func(db *gorm.DB) *gorm.DB {
		if len(req.MemberIds) > 0 {
			db.Where("reserve_order.member_id in ?", req.MemberIds)
		}
		if len(req.StoreIds) > 0 {
			db.Where("reserve_order.store_id in ?", req.StoreIds)
		}
		if len(req.Ids) > 0 && req.Ids[0] > 0 {
			db.Where("reserve_order.id in ?", req.Ids)
		}
		if req.StartTime > 0 && req.EndTime > 0 {
			db.Where("reserve_order.created_at >= ? AND reserve_order.created_at <= ?", req.StartTime, req.EndTime)
		}
		if req.ReserveStartTime > 0 && req.ReserveEndTime > 0 {
			db.Where("reserve_order.reserve_time >= ? AND reserve_order.reserve_time <= ?", req.ReserveStartTime, req.ReserveEndTime)
		}
		if req.CommentStartTime > 0 && req.CommentEndTime > 0 {
			db.Where("reserve_order.comment_time >= ? AND reserve_order.comment_time <= ?", req.CommentStartTime, req.CommentEndTime)
		}
		if req.ReserveNumber != "" {
			db.Where("reserve_order.reserve_number = ?", req.ReserveNumber)
		}
		if req.OrderOrigin > 0 {
			db.Where("reserve_order.order_origin = ?", req.OrderOrigin)
		}
		if req.Status >= 0 {
			db.Where("reserve_order.status = ?", req.Status)
		}
		if req.IsComment > 0 {
			db.Where("reserve_order.is_comment = ?", req.IsComment)
		}
		if req.IsPay > 0 {
			db.Where("reserve_order.is_pay = ?", req.IsPay)
		}
		if req.IsStore >= 0 {
			db.Where("reserve_order.is_store = ?", req.IsStore)
		}
		if req.ServiceId > 0 {
			db.Where("reserve_order.reserve_id=?", req.ServiceId)
		}
		db.Where("reserve_order.reserve_time>0 and reserve_order.deleted_at=0")
		if req.ReserveOrderStartTime > 0 && req.ReserveOrderEndTime > 0 {
			db.Where("sales.pay_time >= ? AND sales.pay_time <= ?", req.ReserveOrderStartTime, req.ReserveOrderEndTime)
		}
		return db
	})
	err := query.Session(&gorm.Session{NewDB: true}).WithContext(c).Count(&total).Error
	err = query.Select(selectStr).Offset((page - 1) * size).Limit(size).Find(&result).Error
	return result, total, err
}

// CreateSaleData 添加邀约销售数据
func (d *ReservesOrderSales) CreateSaleData(c app.Ctx, enterpriseId int64, data microScrmModel.ReserveOrderSalesModel) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(&microScrmModel.ReserveOrderSalesModel{}).Create(&data).Error
}

// CreateOrderGoodsData 添加订单商品数据
func (d *ReservesOrderSales) CreateOrderGoodsData(c app.Ctx, enterpriseId int64, data microScrmModel.ReserveOrderGoodsModel) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.ReserveOrderGoodsModel{}.TableName(), enterpriseId)).Model(&microScrmModel.ReserveOrderGoodsModel{}).Create(&data).Error
}

// GetReserveOrderSales 获取邀约销售订单数据
func (d *ReservesOrderSales) GetReserveOrderSales(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.ReserveOrderSalesModel, error) {
	var result []microScrmModel.ReserveOrderSalesModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(&microScrmModel.ReserveOrderSalesModel{}).Where(query, args...).Find(&result).Error
	return result, err
}

// GetReserveNumberForSales 获取订单商品数据的销售单号
func (d *ReservesOrderSales) GetReserveNumberForSales(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]string, error) {
	result := make([]string, 0)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(&microScrmModel.ReserveOrderSalesModel{}).Where(query, args...).Group("reserve_number").Pluck("reserve_number", &result).Error
	return result, err
}

// GetSalesInfo 获取销售单数据
func (d *ReservesOrderSales) GetSalesInfo(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.ReserveOrderSalesModel, error) {
	var data microScrmModel.ReserveOrderSalesModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(&microScrmModel.ReserveOrderSalesModel{}).Where(query, args...).Find(&data).Error
	return data, err
}

// GetOrderGoodsData 获取订单商品数据
func (d *ReservesOrderSales) GetOrderGoodsData(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.ReserveOrderGoodsModel, error) {
	var data []microScrmModel.ReserveOrderGoodsModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.ReserveOrderGoodsModel{}.TableName(), enterpriseId)).Model(&microScrmModel.ReserveOrderGoodsModel{}).Where(query, args...).Find(&data).Error
	return data, err
}

// UpdateSaleData 更新销售单数据
func (d *ReservesOrderSales) UpdateSaleData(c app.Ctx, enterpriseId, id int64, data microScrmModel.ReserveOrderSalesModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(&microScrmModel.ReserveOrderSalesModel{}).Where(map[string]interface{}{
		"id": id,
	}).UpdateColumns(&data).Error
	return err
}

// GetSalesStatusList 获取邀约状态的列表数据
func (d *ReservesOrderSales) GetSalesStatusList(c app.Ctx, enterpriseId, updateTime int64, status int) ([]reserveRequest.ReserveSalesStatusInfo, error) {
	var data []reserveRequest.ReserveSalesStatusInfo
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(&microScrmModel.ReserveOrderSalesModel{}).Select("reserve_number,SUM(order_price) as order_price,max(pay_time) as last_pay_time,new_customer").Where(map[string]interface{}{
		"updated_at": updateTime,
		"order_type": status,
	}).Group("reserve_number").Find(&data).Error
	return data, err
}

// UpdateOrderGoodsData 更新销售商品数据
func (d *ReservesOrderSales) UpdateOrderGoodsData(c app.Ctx, enterpriseId, id int64, data interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.ReserveOrderGoodsModel{}.TableName(), enterpriseId)).Model(&microScrmModel.ReserveOrderGoodsModel{}).Where(map[string]interface{}{
		"id": id,
	}).UpdateColumns(&data).Error
	return err
}

// GetSalesStatusListByReserveNumber 获取邀约单的数据
func (d *ReservesOrderSales) GetSalesStatusListByReserveNumber(c app.Ctx, enterpriseId int64, reserveNumber []string) ([]reserveRequest.ReserveSalesStatusInfo, error) {
	var data []reserveRequest.ReserveSalesStatusInfo
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(&microScrmModel.ReserveOrderSalesModel{}).Select("reserve_number,SUM(order_price) as order_price,max(pay_time) as last_pay_time,new_customer").Where(map[string]interface{}{
		"reserve_number": reserveNumber,
	}).Group("reserve_number").Find(&data).Error
	return data, err
}
