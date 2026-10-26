package onlineWarehouseDao

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	onlineWarehouseReq "micro-scrm/internal/app/http/request/online_warehouse"
	onlineWarehouseResp "micro-scrm/internal/app/http/response/online_warehouse"
	onlineWarehouseModel "micro-scrm/internal/app/models/micro_scrm/online_warehouse"
	"micro-scrm/internal/pkg/db/gormx"
	"strings"
)

type OrderDao struct {
	dao.EnterpriseRepo[onlineWarehouseModel.OnlineWarehouseOrderModel]
	model     *onlineWarehouseModel.OnlineWarehouseOrderModel
	tableName string
}

func NewOrderDao() *OrderDao {
	return &OrderDao{
		EnterpriseRepo: dao.EnterpriseRepo[onlineWarehouseModel.OnlineWarehouseOrderModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(onlineWarehouseModel.OnlineWarehouseOrderModel),
		tableName: onlineWarehouseModel.OnlineWarehouseOrderModel{}.TableName(),
	}
}

// GetPagination 分页获取订单数据
func (d *OrderDao) GetPagination(c app.Ctx, cond onlineWarehouseReq.ListReq, isExport bool) ([]onlineWarehouseModel.OnlineWarehouseOrderModel, int64, error) {
	var (
		total int64
		resp  []onlineWarehouseModel.OnlineWarehouseOrderModel
		err   error
	)
	query := d.SetDb(gormx.GetConnection(fmt.Sprintf("%s%d", app.Cfg().DB.DBPrefix, cond.EnterpriseId))).WithContext(c).Model(d.model).
		Scopes(func(db *gorm.DB) *gorm.DB {
			// 规则名称
			if cond.OrderSn != "" {
				if cond.IsH5 {
					db.Where("order_sn like ? or user_name like ?", strings.Trim(cond.OrderSn, " ")+"%", "%"+strings.Trim(cond.OrderSn, " ")+"%")
				} else {
					db.Where("order_sn like ?", strings.Trim(cond.OrderSn, " ")+"%")
				}
			}

			if cond.IsAll == true {
				if len(cond.LeaveMemberIds) > 0 {
					memberIds := cond.LeaveMemberIds
					memberIds = append(memberIds, cond.MemberId)
					db.Where("(third_id in ? ) or member_id in ?", cond.ThirdIds, memberIds)
				} else {
					db.Where("(third_id in ? ) or member_id =?", cond.ThirdIds, cond.MemberId)
				}
			} else {
				if len(cond.ThirdIds) > 0 {
					if len(cond.LeaveMemberIds) > 0 {
						db.Where("(third_id in ? ) or member_id in ?", cond.ThirdIds, cond.LeaveMemberIds)
					} else {
						db.Where("third_id in ?", cond.ThirdIds)
					}
				}

			}

			if len(cond.MemberIds) > 0 {
				memberIds := cond.MemberIds
				if len(cond.LeaveMemberIds) > 0 {
					memberIds = append(memberIds, cond.LeaveMemberIds...)
				}
				db.Where("member_id in ?", memberIds)

			}
			if len(cond.LeaveMemberIds) > 0 && len(cond.ThirdIds) == 0 && len(cond.MemberIds) == 0 {
				db.Where("member_id in ?", cond.LeaveMemberIds)
			}
			if cond.Status > 0 {
				db.Where("status = ?", cond.Status)
			}
			// 获得时间
			if cond.StartTime > 0 && cond.EndTime > 0 {
				db.Where("created_at between ? and ?", cond.StartTime, cond.EndTime)
			}
			if cond.TransferStatus > 0 {
				if cond.TransferStatus > onlineWarehouseModel.TransferNO {
					db.Where("transfer_file != ?", "")
				} else { // 待上传转货单
					db.Where("transfer_file = ?", "")
				}
			}
			if cond.OnlyPartRefund == 1 {
				db.Where("send_quantity<total_quantity")
			}
			if cond.MemberId > 0 && cond.IsAll == false {
				db.Where("member_id = ?", cond.MemberId)
			}
			return db
		})
	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size
	query.Order("created_at desc")
	if isExport {
		query.Preload("SkuLists")
	}
	err = query.Limit(cond.Size).Offset(offset).Find(&resp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return resp, total, err
}

// GetPaginationGroupBYStore 根据门店分组获取订单数据
func (d *OrderDao) GetPaginationGroupBYStore(c app.Ctx, cond onlineWarehouseReq.StoreListReq, thirdIds []int) ([]onlineWarehouseResp.StoreListResp, int64, error) {
	var (
		total int64
		resp  []onlineWarehouseResp.StoreListResp
		err   error
	)

	query := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).
		Scopes(func(db *gorm.DB) *gorm.DB {
			db.Where("status=? and transfer_file=?", onlineWarehouseModel.OrderStatusShipped, "")
			if len(thirdIds) > 0 {
				db.Where("third_id in ?", thirdIds)
			}
			// 获得时间
			if cond.ShippedTime > 0 {
				db.Where("shipped_at >= ?  and shipped_at < ?", cond.ShippedTime, cond.ShippedTime+86400)
			}
			return db
		})
	query.Group("CONCAT(`third_id`,UNIX_TIMESTAMP(DATE(FROM_UNIXTIME(shipped_at))))").Count(&total)

	offset := (cond.Page - 1) * cond.Size

	err = query.Select("third_id,count(1) as order_num,UNIX_TIMESTAMP(DATE(FROM_UNIXTIME(shipped_at))) as shipped_time,1 as transfer_status").Order("updated_at desc").Limit(cond.Size).Offset(offset).Scan(&resp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return resp, total, err
}

// GetStoreOrderPagination 门店id+发货时间获取订单列表数据
func (d *OrderDao) GetStoreOrderPagination(c app.Ctx, cond onlineWarehouseReq.StoreOrderListReq) ([]onlineWarehouseModel.OnlineWarehouseOrderModel, int64, error) {
	var (
		total int64
		resp  []onlineWarehouseModel.OnlineWarehouseOrderModel
		err   error
	)

	query := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).
		Scopes(func(db *gorm.DB) *gorm.DB {
			db.Where("shipped_at >= ?  and shipped_at < ?", cond.ShippedTime, cond.ShippedTime+86400).Where("third_id=?", cond.ThirdId)
			return db
		}).Where("transfer_file=?", "")
	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size

	err = query.Order("updated_at desc").Limit(cond.Size).Offset(offset).Find(&resp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return resp, total, err
}

// GetTipsNum 获取各状态统计素
func (d *OrderDao) GetTipsNum(c app.Ctx, cond onlineWarehouseReq.DashboardReq) ([]onlineWarehouseResp.TipsNums, error) {
	var result []onlineWarehouseResp.TipsNums
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).Select("status,count(1) as num").Scopes(func(db *gorm.DB) *gorm.DB {
		if cond.IsAll == true {
			if len(cond.LeaveMemberIds) > 0 {
				memberIds := cond.LeaveMemberIds
				memberIds = append(memberIds, cond.MemberId)
				db.Where("(third_id in ? ) or member_id in ?", cond.ThirdIds, memberIds)
			} else {
				db.Where("(third_id in ? ) or member_id =?", cond.ThirdIds, cond.MemberId)
			}
		} else {
			if len(cond.ThirdIds) > 0 {
				if len(cond.LeaveMemberIds) > 0 {
					db.Where("(third_id in ? ) or member_id in ?", cond.ThirdIds, cond.LeaveMemberIds)
				} else {
					db.Where("third_id in ?", cond.ThirdIds)
				}
			}

		}
		if cond.OrderSn != "" {
			if cond.IsH5 {
				db.Where("order_sn like ? or user_name like ?", strings.Trim(cond.OrderSn, " ")+"%", "%"+strings.Trim(cond.OrderSn, " ")+"%")
			} else {
				db.Where("order_sn like ?", strings.Trim(cond.OrderSn, " ")+"%")
			}
		}
		if len(cond.MemberIds) > 0 {
			memberIds := cond.MemberIds
			if len(cond.LeaveMemberIds) > 0 {
				memberIds = append(memberIds, cond.LeaveMemberIds...)
			}
			db.Where("member_id in ?", memberIds)
		}
		if len(cond.LeaveMemberIds) > 0 && len(cond.ThirdIds) == 0 && len(cond.MemberIds) == 0 {
			db.Where("member_id in?", cond.LeaveMemberIds)
		}
		if cond.StartTime > 0 && cond.EndTime > 0 {
			db.Where("created_at between ? and ?", cond.StartTime, cond.EndTime)
		}
		if cond.MemberId > 0 && cond.IsAll == false {
			db.Where("member_id =?", cond.MemberId)
		}
		return db
	}).Group("status").Scan(&result).Error
	return result, err
}

func (d *OrderDao) GetCountGroupThirdId(c ctx.Context, cond onlineWarehouseReq.DashboardReq) ([]onlineWarehouseResp.StoreRankResp, error) {
	var result []onlineWarehouseResp.StoreRankResp
	skuTableName := d.Dao.GetTable(onlineWarehouseModel.OnlineWarehouseOrderSkuModel{}.TableName(), cond.EnterpriseId)
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).Joins("left join "+skuTableName+" as b on online_warehouse_order.id=b.order_id").
		Select("third_id,sum(b.quantity) as sku_num,count(DISTINCT(online_warehouse_order.id)) as num").Scopes(func(db *gorm.DB) *gorm.DB {
		if cond.StartTime > 0 && cond.EndTime > 0 {
			db.Where("online_warehouse_order.created_at between ? and ?", cond.StartTime, cond.EndTime)
		}
		return db
	}).Where("b.status=?", onlineWarehouseModel.OrderStatusShipped).Group("third_id").Order("sku_num,num desc").Scan(&result).Error
	return result, err
}

func (d *OrderDao) GetCountGroupMemberId(c ctx.Context, cond onlineWarehouseReq.DashboardReq) ([]onlineWarehouseResp.SaRankResp, error) {
	var result []onlineWarehouseResp.SaRankResp
	skuTableName := d.Dao.GetTable(onlineWarehouseModel.OnlineWarehouseOrderSkuModel{}.TableName(), cond.EnterpriseId)
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).Joins("left join "+skuTableName+" as b on online_warehouse_order.id=b.order_id").
		Select("third_id,member_id,sum(b.quantity) as sku_num,count(DISTINCT(online_warehouse_order.id)) as num").Scopes(func(db *gorm.DB) *gorm.DB {
		db.Where("third_id in ?", cond.ThirdIds)
		if cond.StartTime > 0 && cond.EndTime > 0 {
			db.Where("online_warehouse_order.created_at between ? and ?", cond.StartTime, cond.EndTime)
		}
		return db
	}).Where("b.status=?", onlineWarehouseModel.OrderStatusShipped).Group("member_id").Order("num desc").Limit(10).Scan(&result).Error
	return result, err
}

func (d *OrderDao) GetUpdateNumByMember(c ctx.Context, enterpriseid int64, startTime, endTime int64) ([]onlineWarehouseResp.MemberNumResp, error) {
	var result []onlineWarehouseResp.MemberNumResp
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, enterpriseid)).Model(d.model).Select("count(1) num,member_id,updated_at").Where("status= ? and updated_at between ? and ?", onlineWarehouseModel.OrderStatusWait, startTime, endTime).Group("member_id").Scan(&result).Error
	return result, err
}

func (d *OrderDao) UpdateData(c ctx.Context, enterpriseId int64, query any, data any) error {
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Model(d.model).Where(query).Updates(data).Error
	return err
}

func (d *OrderDao) GetUpdateByOrderId(c ctx.Context, enterpriseid int64, orderId []int64) ([]onlineWarehouseResp.MemberNumResp, error) {
	var result []onlineWarehouseResp.MemberNumResp
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, enterpriseid)).Model(d.model).Select("count(1) num,member_id,updated_at").Where(map[string]interface{}{
		"id": orderId,
	}).Group("member_id").Scan(&result).Error
	return result, err
}

func (d *OrderDao) GetPendingNums(c ctx.Context, enterpriseId, startTime, endTime int64) (int64, error) {
	var total int64
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where("status=? and updated_at between ? and ?", onlineWarehouseModel.OrderStatusPending, startTime, endTime).Count(&total).Error
	return total, err
}

func (d *OrderDao) GetSendMail(c ctx.Context, enterpriseId int64, typeValue int32) ([]onlineWarehouseModel.OnlineWarehouseEmailSendModel, error) {
	var result []onlineWarehouseModel.OnlineWarehouseEmailSendModel
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(onlineWarehouseModel.OnlineWarehouseEmailSendModel{}.TableName(), enterpriseId)).Model(d.model).Where("type=?", typeValue).Find(&result).Error
	return result, err
}

func (d *OrderDao) GetAllMemberId(c ctx.Context, enterpriseId int64, thirdIds []int) ([]int64, error) {
	var result []int64
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(map[string]interface{}{
		"third_id": thirdIds,
	}).Group("member_id").Pluck("member_id", &result).Error
	return result, err
}

func (d *OrderDao) GetAllThirdId(c ctx.Context, enterpriseId int64, thirdIds []int) ([]int, error) {
	var result []int
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(map[string]interface{}{
		"third_id": thirdIds,
	}).Group("third_id").Pluck("third_id", &result).Error
	return result, err
}

// GetAllTransferInfoByThird 获取门店的未上传的转货单
func (d *OrderDao) GetAllTransferInfoByThird(c app.Ctx, enterpriseId int64, thirdId int) ([]onlineWarehouseModel.OnlineWarehouseOrderModel, error) {
	var resp []onlineWarehouseModel.OnlineWarehouseOrderModel
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).
		Where("third_id=? and status=?", thirdId, onlineWarehouseModel.OrderStatusShipped).
		Where("transfer_file=?", "").Group("id").Find(&resp).Error
	return resp, err
}

func (d *OrderDao) GetLastOrderInfo(c app.Ctx, enterpriseId int64, orderSn string) (onlineWarehouseModel.OnlineWarehouseOrderModel, error) {
	var resp onlineWarehouseModel.OnlineWarehouseOrderModel
	err := d.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(map[string]interface{}{
		"order_sn": orderSn,
	}).Where("status<?", onlineWarehouseModel.OrderStatusCancel).Order("id desc").First(&resp).Error
	return resp, err
}

func (d *OrderDao) GetOrderIds(c ctx.Context, cond onlineWarehouseReq.DashboardReq) ([]int64, error) {
	var orderIds []int64
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).Scopes(func(db *gorm.DB) *gorm.DB {
		if cond.IsAll == true {
			if len(cond.LeaveMemberIds) > 0 {
				memberIds := cond.LeaveMemberIds
				memberIds = append(memberIds, cond.MemberId)
				db.Where("(third_id in ? ) or member_id in ?", cond.ThirdIds, memberIds)
			} else {
				db.Where("(third_id in ? ) or member_id =?", cond.ThirdIds, cond.MemberId)
			}
		} else {
			if len(cond.ThirdIds) > 0 {
				if len(cond.LeaveMemberIds) > 0 {
					db.Where("(third_id in ? ) or member_id in ?", cond.ThirdIds, cond.LeaveMemberIds)
				} else {
					db.Where("third_id in ?", cond.ThirdIds)
				}
			}

		}
		if cond.OrderSn != "" {
			if cond.IsH5 {
				db.Where("order_sn = ? or user_name like ?", strings.Trim(cond.OrderSn, " "), "%"+strings.Trim(cond.OrderSn, " ")+"%")
			} else {
				db.Where("order_sn = ?", strings.Trim(cond.OrderSn, " "))
			}
		}
		if len(cond.MemberIds) > 0 {
			memberIds := cond.MemberIds
			if len(cond.LeaveMemberIds) > 0 {
				memberIds = append(memberIds, cond.LeaveMemberIds...)
			}
			db.Where("member_id in ?", memberIds)
		}
		if len(cond.LeaveMemberIds) > 0 && len(cond.ThirdIds) == 0 && len(cond.MemberIds) == 0 {
			db.Where("member_id in?", cond.LeaveMemberIds)
		}
		if cond.StartTime > 0 && cond.EndTime > 0 {
			db.Where("created_at between ? and ?", cond.StartTime, cond.EndTime)
		}
		if cond.MemberId > 0 && cond.IsAll == false {
			db.Where("member_id =?", cond.MemberId)
		}
		return db
	}).Pluck("id", &orderIds).Error
	return orderIds, err
}
