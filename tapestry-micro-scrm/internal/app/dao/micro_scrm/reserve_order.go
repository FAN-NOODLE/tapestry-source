package microScrmDao

import (
	"context"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	customerReq "micro-scrm/internal/app/http/request/customer"
	reserveReq "micro-scrm/internal/app/http/request/open"
	planReq "micro-scrm/internal/app/http/request/plan"
	reserveRequest "micro-scrm/internal/app/http/request/reserve"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"time"
)

type ReservesOrder struct {
	dao       dao.Dao
	tableName string
}

func NewReservesOrder() *ReservesOrder {
	return &ReservesOrder{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.ReserveOrderModel{}.TableName(),
	}
}

func (d *ReservesOrder) getScrmTableCk(table string, enterpriseId int64) string {
	return fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, table)
}

// GetList
func (d *ReservesOrder) GetList(c app.Ctx, enterpriseId int64, page, size int, req reserveRequest.OrderListReq, unionIdsAndMemberId []customerReq.UserUnionIdAndMemberId) ([]microScrmModel.ReserveOrderModel, int64, error) {
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("reserve_time > 0 and deleted_at=0")
	var (
		total     int64
		orderList []microScrmModel.ReserveOrderModel
		err       error
	)

	if (req.ServiceId) > 0 {
		query.Where("reserve_id = ?", req.ServiceId)
	}
	if (req.Status) >= 0 {
		query.Where("status = ?", req.Status)
	}
	if req.IsSM == 1 { // 如果是店长
		query.Where("(member_id = ? or reserve_order.store_id in ?) ", req.MemberId, req.StoreIds)
	} else {
		if len(req.MemberIds) > 0 {
			query.Where("member_id in ?", req.MemberIds)
		}
		if len(req.StoreIds) > 0 {
			query.Where("store_id in ?", req.StoreIds)
		}
	}

	if len(req.Ids) > 0 {
		query.Where("id in ?", req.Ids)
	}
	if req.ReserveNumber != "" {
		query.Where("reserve_number = ?", req.ReserveNumber)
	}
	if len(req.ReserveNumbers) > 0 {
		query.Where("reserve_number in (?)", req.ReserveNumbers)
	}
	if req.ReserveStartTime > 0 && req.ReserveEndTime > 0 {
		query.Where("reserve_time >= ? AND reserve_time <= ?", req.ReserveStartTime, req.ReserveEndTime)

	}
	if req.StartTime > 0 && req.EndTime > 0 {
		query.Where("created_at >= ? AND created_at <= ?", req.StartTime, req.EndTime)

	}
	if req.IsStore >= 0 {
		query.Where("is_store = ?", req.IsStore)
	}
	if req.IsComment > 0 {
		query.Where("is_comment=?", req.IsComment)
	}
	if req.IsPay > 0 {
		query.Where("is_pay = ?", req.IsPay)
	}
	if req.ReserveOrderStartTime > 0 && req.ReserveOrderEndTime > 0 {
		query.Where("reserve_order.reserve_order_time >= ? AND reserve_order.reserve_order_time <= ?", req.ReserveOrderStartTime, req.ReserveOrderEndTime)
	}
	if req.Unionid != "" {
		query.Where("user_unionid = ?", req.Unionid)

	}
	if req.OrderOrigin > 0 {
		query.Where("order_origin = ?", req.OrderOrigin)
	}
	if len(unionIdsAndMemberId) > 0 {
		var userUnionidAndMemberId [][]interface{}
		for _, item := range unionIdsAndMemberId {
			userUnionidAndMemberId = append(userUnionidAndMemberId, []interface{}{
				item.Unionid, item.MemberId,
			})
		}
		query.Where("(user_unionid,member_id) in ?", userUnionidAndMemberId)
	}
	if req.Name != "" {
		query.Where("user_name like ?", "%"+req.Name+"%")
	}
	str := "id desc"
	if req.OrderSort == 1 {
		str = "FIELD(`status`,0,2,1,3,4) asc,reserve_time desc"
		err = query.Preload("ReserveService", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveServiceModel{}.TableName(), enterpriseId)).Unscoped()
		}).Preload("ReserveOrderSales", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveOrderSalesModel{}.TableName(), enterpriseId)).Preload("ReserveOrderGoods", func(db *gorm.DB) *gorm.DB {
				return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveOrderGoodsModel{}.TableName(), enterpriseId))
			}).Order("pay_time desc")
		}).Order(str).Offset((page - 1) * size).Limit(size).Find(&orderList).Error
	} else {
		err = query.Count(&total).Preload("ReserveService", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveServiceModel{}.TableName(), enterpriseId)).Unscoped()
		}).Preload("ReserveOrderSales", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveOrderSalesModel{}.TableName(), enterpriseId)).Preload("ReserveOrderGoods", func(db *gorm.DB) *gorm.DB {
				return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveOrderGoodsModel{}.TableName(), enterpriseId))
			}).Order("pay_time desc")
		}).Order(str).Offset((page - 1) * size).Limit(size).Find(&orderList).Error
	}

	return orderList, total, err
}

// GetListByCoach  Coach小程序端获取邀约
func (d *ReservesOrder) GetListByCoach(c app.Ctx, req reserveRequest.OrderListByCoachReq) ([]microScrmModel.ReserveOrderModel, int64, error) {
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).Where("reserve_time > 0 and deleted_at=0")
	var (
		total     int64
		orderList []microScrmModel.ReserveOrderModel
		err       error
	)
	if (req.Status) >= 0 {
		query.Where("status = ?", req.Status)
	}
	if len(req.MemberIds) > 0 {
		query.Where("member_id in ?", req.MemberIds)
	}
	if len(req.ThirdIds) > 0 {
		query.Where("store_id in ?", req.ThirdIds)
	}

	if req.ReserveStartTime > 0 && req.ReserveEndTime > 0 {
		query.Where("reserve_time >= ? AND reserve_time <= ?", req.ReserveStartTime, req.ReserveEndTime)
	}
	if len(req.UnionIds) > 0 {
		query.Where("user_unionid in ?", req.UnionIds)
	}
	str := "id desc"
	str = "FIELD(`status`,0,2,1,3,4) asc,reserve_time desc"
	err = query.Count(&total).Preload("ReserveService", func(db *gorm.DB) *gorm.DB {
		return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveServiceModel{}.TableName(), req.EnterpriseId)).Unscoped()
	}).Order(str).Offset((req.Page - 1) * req.Size).Limit(req.Size).Find(&orderList).Error

	return orderList, total, err
}

// GetExportList 数据导出
func (d *ReservesOrder) GetExportList(c app.Ctx, enterpriseId int64, page, size int, req reserveRequest.OrderListReq, isCoach bool) ([]reserveRequest.OrderListRes, int64, error) {
	var total int64
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId))
	if isCoach { // coach的查ck数据
		query = gormx.GetCkInstance().WithContext(c).Table(d.getScrmTableCk(d.tableName, enterpriseId))
		query.Scopes(func(db *gorm.DB) *gorm.DB {
			if req.ReserveOrderStartTime > 0 && req.ReserveOrderEndTime > 0 {
				db.Where("reserve_number in (?)", gormx.GetCkInstance().WithContext(c).Table(d.getScrmTableCk(microScrmModel.ReserveOrderSalesModel{}.TableName(), enterpriseId)).Select("reserve_number").Where("pay_time BETWEEN ? AND ? ", req.ReserveOrderStartTime, req.ReserveOrderEndTime))
			}
			return db
		})
	}
	if !isCoach {
		query.Joins("left join " + d.dao.GetTable(microScrmModel.ReserveServiceScoreModel{}.TableName(), enterpriseId) + " AS b ON reserve_order.id=b.reserve_order_id")
		if req.ReserveOrderStartTime > 0 && req.ReserveOrderEndTime > 0 {
			query.Where("reserve_order.reserve_order_time >= ? AND reserve_order.reserve_order_time <= ?", req.ReserveOrderStartTime, req.ReserveOrderEndTime)
		}
	}
	query.Where("reserve_order.reserve_time > 0 and reserve_order.deleted_at=0")
	var (
		orserList []reserveRequest.OrderListRes
		err       error
	)

	if (req.ServiceId) > 0 {
		query.Where("reserve_order.reserve_id = ?", req.ServiceId)
	}
	if (req.Status) >= 0 {
		query.Where("reserve_order.status = ?", req.Status)
	}
	if req.IsSM == 1 { // 如果是店长
		query.Where("reserve_order.member_id = ? or reserve_order.store_id in ? ", req.MemberId, req.StoreIds)
	} else {
		if len(req.MemberIds) > 0 {
			query.Where("reserve_order.member_id in ?", req.MemberIds)
		}
		if len(req.StoreIds) > 0 {
			query.Where("reserve_order.store_id in ?", req.StoreIds)
		}
	}
	if len(req.Ids) > 0 {
		query.Where("reserve_order.id in ?", req.Ids)
	}
	if req.ReserveNumber != "" {
		query.Where("reserve_order.reserve_number = ?", req.ReserveNumber)
	}
	if req.ReserveStartTime > 0 && req.ReserveEndTime > 0 {
		query.Where("reserve_order.reserve_time >= ? AND reserve_order.reserve_time <= ?", req.ReserveStartTime, req.ReserveEndTime)
	}
	if req.StartTime > 0 && req.EndTime > 0 {
		query.Where("reserve_order.created_at >= ? AND reserve_order.created_at <= ?", req.StartTime, req.EndTime)
	}
	if req.IsStore >= 0 {
		query.Where("reserve_order.is_store = ?", req.IsStore)
	}
	if req.IsComment > 0 {
		query.Where("reserve_order.is_comment=?", req.IsComment)
	}
	if req.IsPay > 0 {
		query.Where("reserve_order.is_pay = ?", req.IsPay)
	}

	if req.Unionid != "" {
		query.Where("reserve_order.user_unionid = ?", req.Unionid)

	}
	if req.OrderOrigin > 0 {
		query.Where("reserve_order.order_origin = ?", req.OrderOrigin)
	}
	if req.Name != "" {
		query.Where("reserve_order.user_name like ?", "%"+req.Name+"%")
	}
	str := "reserve_order.created_at desc"
	selectStr := "reserve_order.*"

	if !isCoach {
		selectStr = "reserve_order.*,b.service_score,b.store_score,b.remark"
	}
	if !isCoach {
		query.Count(&total).Preload("ReserveOrderSales", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveOrderSalesModel{}.TableName(), enterpriseId)).Preload("ReserveOrderGoods", func(db *gorm.DB) *gorm.DB {
				return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveOrderGoodsModel{}.TableName(), enterpriseId))
			})
		})
	} else {
		query.Count(&total).Preload("ReserveOrderSales", func(db *gorm.DB) *gorm.DB {
			return db.Table(d.getScrmTableCk(microScrmModel.ReserveOrderSalesModel{}.TableName(), enterpriseId)).Preload("ReserveOrderGoods", func(db *gorm.DB) *gorm.DB {
				return db.Table(d.getScrmTableCk(microScrmModel.ReserveOrderGoodsModel{}.TableName(), enterpriseId))
			}).Scopes(func(db *gorm.DB) *gorm.DB {
				if req.ReserveOrderStartTime > 0 && req.ReserveOrderEndTime > 0 {
					db.Where("pay_time >= ? AND pay_time <= ?", req.ReserveOrderStartTime, req.ReserveOrderEndTime)
				}
				return db
			})
		})
	}
	err = query.Order(str).
		Select(selectStr).
		Offset((page - 1) * size).
		Limit(size).
		Find(&orserList).
		Error

	return orserList, total, err
}

func (d *ReservesOrder) CreateOrder(c app.Ctx, enterpriseId int64, info microScrmModel.ReserveOrderModel) (microScrmModel.ReserveOrderModel, error) {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&info).Error
	return info, err
}

func (d *ReservesOrder) CancelOrder(c app.Ctx, enterpriseId int64, orderBase reserveRequest.OrderBase) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"id":             orderBase.ID,
		"reserve_number": orderBase.ReserveNumber,
		"status":         microScrmModel.RUNNING,
	}).Updates(map[string]interface{}{
		"status":     microScrmModel.CANCEL,
		"updated_at": time.Now().Unix(),
	}).Error
}

func (d *ReservesOrder) GetOrderInfo(c app.Ctx, enterpriseId int64, userOrderBase interface{}) (orderInfo microScrmModel.ReserveOrderModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Unscoped().Where(userOrderBase).Find(&orderInfo).Error
	return
}

func (d *ReservesOrder) GetOrderPreLoadDetail(c app.Ctx, enterpriseId int64, userOrderBase interface{}) (orderInfo microScrmModel.ReserveOrderModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Unscoped().Where(userOrderBase).Preload("ReserveService", func(db *gorm.DB) *gorm.DB {
		return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveServiceModel{}.TableName(), enterpriseId)).Unscoped()
	}).Preload("ReserveOrderSales", func(db *gorm.DB) *gorm.DB {
		return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveOrderSalesModel{}.TableName(), enterpriseId)).Preload("ReserveOrderGoods", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveOrderGoodsModel{}.TableName(), enterpriseId))
		}).Order("pay_time desc")
	}).Find(&orderInfo).Error
	return
}

func (d *ReservesOrder) OpenOrder(c app.Ctx, enterpriseId, id int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("id=?", id).Update("is_open", 1).Error
}

func (d *ReservesOrder) CheckOrderInfo(c app.Ctx, enterpriseId int64, userOrderBase reserveRequest.UserOrderBase) (orderInfo microScrmModel.ReserveOrderModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"id":             userOrderBase.ID,
		"reserve_number": userOrderBase.ReserveNumber,
		"user_unionid":   userOrderBase.UnionId,
	}).Find(&orderInfo).Error
	return
}

func (d *ReservesOrder) UserCommitOrder(c app.Ctx, enterpriseId int64, UserOrderBase reserveRequest.UserOrderBase) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"id":             UserOrderBase.ID,
		"reserve_number": UserOrderBase.ReserveNumber,
		"user_unionid":   UserOrderBase.UnionId,
	}).Update("is_commit", 1).Error
}

func (d *ReservesOrder) UserEditOrder(c app.Ctx, enterpriseId int64, request reserveRequest.UserEditOrder) error {
	if request.Status == 1 { // 用户确认订单
		return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
			"id":             request.ID,
			"reserve_number": request.ReserveNumber,
			"user_unionid":   request.UnionId,
		}).Updates(map[string]interface{}{
			"reserve_id":       request.ReserveId,
			"reserve_time":     request.ReserveTime,
			"reserve_end_time": request.ReserveEndTime,
			"user_remark":      request.UserRemark,
			"is_commit":        1,
		}).Error
	}
	if request.Status == 2 { // 用户取消订单
		return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
			"id":             request.ID,
			"reserve_number": request.ReserveNumber,
		}).Updates(map[string]interface{}{
			"status":      microScrmModel.CANCEL,
			"updated_at":  time.Now().Unix(),
			"user_remark": request.UserRemark,
		}).Error
	}
	return nil

}

func (d *ReservesOrder) UserConfirmOrder(c app.Ctx, enterpriseId int64, request reserveRequest.UserConfirmOrder, isConfirm int) error {
	if isConfirm == 1 {
		return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
			Where(map[string]interface{}{
				"id": request.Id,
			}).Updates(map[string]interface{}{
			"reserve_time":     request.ReserveTime,
			"reserve_end_time": request.ReserveEndTime,
			"invite_num":       1,
			"is_commit":        1,
		}).Error
	} else {
		return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
			Where(map[string]interface{}{
				"id": request.Id,
			}).Updates(map[string]interface{}{
			"is_open": 1,
		}).Error
	}

}

func (d *ReservesOrder) Update(c app.Ctx, enterpriseId int64, cond interface{}, data interface{}) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(cond).UpdateColumns(data).Error
}

// GetReserveTimeList 时间段内的邀约单数
func (d *ReservesOrder) GetReserveTimeList(c app.Ctx, enterpriseId, memberId, startTime, endTime int64) (result []reserveRequest.OrderTimeInfo, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Select("reserve_time,count(1) as num").
		Where(map[string]interface{}{
			"member_id": memberId,
			"status":    []int{microScrmModel.RUNNING, microScrmModel.WAITFINISH},
		}).Where("reserve_time BETWEEN ? and ? ", startTime, endTime).Group("reserve_time").Find(&result).Error
	return
}

// GetReserveOrderList 获取邀约单列表
func (d *ReservesOrder) GetReserveOrderList(c app.Ctx, enterpriseId int64, page, size int, req reserveRequest.ReserveOrderListReq, isCoach bool) (reserveRequest.ReserveSum, []microScrmModel.ReserveOrderModel, int64, error) {
	query := gormx.GetCkInstance().WithContext(c).Table(d.getScrmTableCk(d.tableName, enterpriseId)).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(req.MemberIds) > 0 {
			db.Where("member_id in ?", req.MemberIds)
		}
		if len(req.StoreIds) > 0 {
			db.Where("store_id in ?", req.StoreIds)
		}
		if len(req.Ids) > 0 && req.Ids[0] > 0 {
			db.Where("id in ?", req.Ids)
		}
		if req.StartTime > 0 && req.EndTime > 0 {
			db.Where("created_at >= ? AND created_at <= ?", req.StartTime, req.EndTime)
		}
		// if req.ReserveStartTime > 0 && req.ReserveEndTime > 0 {
		// 	db.Where("reserve_order_time >= ? AND reserve_order_time <= ?", req.ReserveStartTime, req.ReserveEndTime)
		// }
		if req.CommentStartTime > 0 && req.CommentEndTime > 0 {
			db.Where("comment_time >= ? AND comment_time <= ?", req.CommentStartTime, req.CommentEndTime)
		}
		if req.OrderOrigin > 0 {
			db.Where("order_origin = ?", req.OrderOrigin)
		}
		if req.ReserveNumber != "" {
			db.Where("reserve_number = ?", req.ReserveNumber)
		}
		if req.Status >= 0 {
			db.Where("status = ?", req.Status)
		}
		if req.IsComment > 0 {
			db.Where("is_comment = ?", req.IsComment)
		}
		if req.IsPay > 0 {
			db.Where("is_pay = ?", req.IsPay)
		}
		if req.IsStore >= 0 {
			db.Where("is_store = ?", req.IsStore)
		}
		return db
	}).Where("reserve_time>0 and deleted_at=0").Scopes(func(db *gorm.DB) *gorm.DB {
		if req.ReserveOrderStartTime > 0 && req.ReserveOrderEndTime > 0 {
			db.Where("reserve_number in (?)",
				gormx.GetCkInstance().WithContext(c).Table(d.getScrmTableCk(microScrmModel.ReserveOrderSalesModel{}.TableName(), enterpriseId)).Select("reserve_number").Where("pay_time BETWEEN ? AND ? ", req.ReserveOrderStartTime, req.ReserveOrderEndTime))
		}
		return db
	})
	var (
		total            int64
		reserveOrderList []microScrmModel.ReserveOrderModel
		err              error
	)

	// var querySum = query
	var nums reserveRequest.ReserveSum
	// querySum.Select("sum(invite_num) as invite_num,sum(is_open) as open_num,sum(is_commit) as commit_num").Find(&nums)

	var inviteSuccessNum, paySuccessNum, successPrice int64
	inviteSuccessNum, paySuccessNum, successPrice, nums = d.buildReserveOrderListCountSql(c, req, isCoach)

	// querySum.Where("is_pay = ?", 1).Count(&inviteSuccessNum)
	// querySum.Where("is_store = ?", 1).Count(&paySuccessNum)

	nums.InviteSuccessNum = inviteSuccessNum
	nums.PaySuccessNum = paySuccessNum
	nums.SuccessPrice = successPrice
	err = query.Count(&total).
		Select("*").
		Order("created_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&reserveOrderList).
		Error
	if nums.OpenNum == 0 {
		nums.OpenPer = "0%"
	} else {
		nums.OpenPer = helper.MathRoundedPer(float64(nums.OpenNum), float64(total))
	}

	if nums.CommitNum == 0 {
		nums.CommitPer = "0%"
	} else {
		nums.CommitPer = helper.MathRoundedPer(float64(nums.CommitNum), float64(total))
	}
	return nums, reserveOrderList, total, err
}

// GetTodayOrderList 当天的邀约单列表
func (d *ReservesOrder) GetTodayOrderList(c app.Ctx, enterpriseId, startTime, endTime int64) (result []reserveRequest.PushMemberInfo, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("member_user_id,member_id,count(1) as num").
		Where("reserve_time BETWEEN ? AND ? ", startTime, endTime).
		Where("status=?", microScrmModel.RUNNING).Group("member_user_id").Find(&result).Error
	return
}

// GetPreTodayOrderList 前一天的待完成邀约单
func (d *ReservesOrder) GetPreTodayOrderList(c app.Ctx, enterpriseId, startTime, endTime int64, status int) (result []reserveRequest.PushMemberInfo, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("member_user_id,member_id,count(1) as num").
		Where("reserve_time BETWEEN ? AND ? ", startTime, endTime).
		Where("status=?", status).Group("member_user_id").Find(&result).Error
	return
}

// GetRankBaseList 邀约榜
func (d *ReservesOrder) GetRankBaseList(c app.Ctx, enterpriseId int64, request reserveRequest.OrderRankList) (result []reserveRequest.OrderRankRes, err error) {
	offset := (request.Page - 1) * request.Size
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("reserve_order.member_id,reserve_order.store_id,count(1) as num,"+
			"sum(IF(reserve_order.is_commit,1,0)) as success_num,"+
			"m.name as member_name,m.avatar").
		Joins("LEFT JOIN members as m ON m.id=reserve_order.member_id").
		Where("reserve_order.created_at BETWEEN ? AND ? ", request.StartTime, request.EndTime).
		// Where("reserve_order.status=?", microScrmModel.RUNNING).//统计的为所有的邀约单，不是进行中的
		Group("member_id,store_id").
		Order("num desc, member_id desc").
		Offset(offset).Limit(request.Size).
		Find(&result).Error
	return

}

// GetRankPrice 获取排行的成单金额
func (d *ReservesOrder) GetRankPrice(c app.Ctx, enterpriseId int64, memberIds []int64, startTime, endTime int64) (result []reserveRequest.OrderPriceRank, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("member_id,sum(reserve_order_price) as price").
		Where("member_id", memberIds).
		Where("updated_at BETWEEN ? AND ? ", startTime, endTime).
		Group("member_id").
		Find(&result).Error
	return
}

func (d *ReservesOrder) Increment(c app.Ctx, enterpriseId int64, cond interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where(cond).
		UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s + ?", column), num)).Error
	return err
}

// GetInfoByThirdOrderNo 判断coach的三方订单是否被核销过
func (d *ReservesOrder) GetInfoByThirdOrderNo(c app.Ctx, enterpriseId int64, unionId, orderNo string) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"user_unionid":     unionId,
		"reserve_order_sn": orderNo,
	}).Count(&total).Error
	return
}

// UpdatePreTodayOrders 更新前一天的订单状态
func (d *ReservesOrder) UpdatePreTodayOrders(c app.Ctx, enterpriseId int64, startTime, endTime int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("reserve_time BETWEEN ? AND ? ", startTime, endTime).
		Where("status=?", microScrmModel.RUNNING).
		Update("status", microScrmModel.WAITFINISH).Error
}

// UserSignIn 用户确认订单（SW的功能）
func (d *ReservesOrder) UserSignIn(c app.Ctx, enterpriseId int64, request reserveRequest.UserSignIn) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"id":             request.ID,
		"reserve_number": request.ReserveNumber,
		"user_unionid":   request.UnionId,
	}).Updates(map[string]interface{}{
		"is_sign": 1,
	}).Error

}

// GetReserveOrderBaseList 获取基本的预约单列表
func (d *ReservesOrder) GetReserveOrderBaseList(c app.Ctx, enterpriseId int64, page, size int, req reserveRequest.ReserveOrderListReq, status int) ([]microScrmModel.ReserveOrderModel, int64, error) {

	var (
		total            int64
		reserveOrderList []microScrmModel.ReserveOrderModel
	)

	if (req.ReserveStartTime == 0 || req.ReserveEndTime == 0) && (req.StartTime == 0 || req.EndTime == 0) {
		return reserveOrderList, 0, nil
	}
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Scopes(func(db *gorm.DB) *gorm.DB {
		if req.ReserveStartTime > 0 && req.ReserveEndTime > 0 {
			db.Where("reserve_time >= ? AND reserve_time <= ?", req.ReserveStartTime, req.ReserveEndTime)
		}
		if req.StartTime > 0 && req.EndTime > 0 {
			db.Where("created_at BETWEEN ? AND ?", req.ReserveStartTime, req.ReserveEndTime)
		}
		if status >= 0 {
			db.Where("status=?", status)
		}
		return db
	})
	err := query.Count(&total).
		Select("*").
		Order("created_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&reserveOrderList).
		Error
	return reserveOrderList, total, err
}

// ChangeStatus  导购确认订单是否是到店购买,设置新客老客
func (d *ReservesOrder) ChangeStatus(c app.Ctx, enterpriseId int64, request reserveRequest.ChangeStatus) error {
	mapStruct := make(map[string]interface{})
	mapStruct["updated_at"] = time.Now().Unix()
	if request.IsStore > 0 {
		mapStruct["is_store"] = request.IsStore
	}
	if request.CustomerStatus > 0 {
		mapStruct["customer_status"] = request.CustomerStatus
	}
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"id":             request.ID,
		"reserve_number": request.ReserveNumber,
	}).Updates(mapStruct).Error
}

// ConfirmStore 导购确认用户是否到店
func (d *ReservesOrder) ConfirmStore(c ctx.Context, enterpriseId int64, request reserveRequest.OrderAssignReq) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"id":             request.ID,
		"reserve_number": request.ReserveNumber,
	}).Updates(map[string]interface{}{
		"is_confirm": 1,
		"status":     microScrmModel.WAITFINISH,
		"updated_at": time.Now().Unix(),
	}).Error
}

// UserCommentOrder 用户评论订单
func (d *ReservesOrder) UserCommentOrder(c app.Ctx, enterpriseId int64, request reserveRequest.UserOrderBase) error {
	var commentTime = int(time.Now().Unix())
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"id":             request.ID,
		"reserve_number": request.ReserveNumber,
	}).Updates(map[string]interface{}{
		"is_comment":   microScrmModel.COMMENT_YES,
		"comment_time": commentTime,
	}).Error
}

// GetOrderMonthLists 获取每个月的预约单列表
func (d *ReservesOrder) GetOrderMonthLists(c ctx.Context, request planReq.MonthListsReq, memberId int64) ([]microScrmModel.ReserveOrderModel, error) {
	var result []microScrmModel.ReserveOrderModel
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).Where("member_id = ? and deleted_at=0", memberId).
		Where("reserve_time BETWEEN ? AND ?", request.StartTime, request.EndTime)
	err := query.Order("status asc").Find(&result).Error
	return result, err
}

// GetOrderInfoByServiceId 根据邀约服务ID查询邀约单
func (d *ReservesOrder) GetOrderInfoByServiceId(c app.Ctx, enterpriseId, ReserveId int64) (orderInfo microScrmModel.ReserveOrderModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("reserve_id=?", ReserveId).Find(&orderInfo).Error
	return
}

// GetPlanTipsCount 获取到点预约进行中的订单数
func (d *ReservesOrder) GetPlanTipsCount(c app.Ctx, enterpriseId, startTime, endTime, memberId int64) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("member_id = ? and deleted_at=0 and status=0", memberId).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if startTime > 0 && endTime > 0 {
				db.Where("reserve_time BETWEEN ? AND ?", startTime, endTime)
			}
			return db
		}).Count(&total).Error
	return
}

// GetSearchUnionid 获取要筛选的unionid
func (d *ReservesOrder) GetSearchUnionid(c ctx.Context, request reserveRequest.OrderListReq, memberIds []int64) (data []string, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).Scopes(func(db *gorm.DB) *gorm.DB {
		db.Where("member_id in ?", memberIds)
		if len(request.Search.CustomerType) > 0 {
			var customerType []int
			// vip客
			if helper.InArray(microScrmModel.VipTypeChoice, request.Search.CustomerType) {
				db.Where("is_vip = ?", microScrmModel.VipType)

				item := []int{microScrmModel.VipTypeChoice}
				customerType = helper.DelItems(request.Search.CustomerType, item)
			}
			if len(customerType) > 0 {
				db.Where("customer_status in (?)", customerType)
			}
		}
		return db
	}).Where("reserve_time BETWEEN ? AND ? AND deleted_at=0 ", request.ReserveStartTime, request.ReserveEndTime).Group("user_unionid").Pluck("user_unionid", &data).Error
	return
}

// GetReserveTipsCount 获取客户的所有邀约单总数
func (d *ReservesOrder) GetReserveTipsCount(c app.Ctx, enterpriseId int64, userExternalId string) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("user_external_id = ? ", userExternalId).
		Where("reserve_time > 0  and deleted_at=0").
		Count(&total).Error
	return
}

func (d *ReservesOrder) buildReserveOrderListCountSql(c ctx.Context, req reserveRequest.ReserveOrderListReq, isCoach bool) (inviteSuccessNum, paySuccessNum, successPrice int64, reserveSum reserveRequest.ReserveSum) {
	// query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).Scopes(func(db *gorm.DB) *gorm.DB {
	isNotAllSearch := false
	query := gormx.GetCkInstance().WithContext(c).Table(d.getScrmTableCk(d.tableName, req.EnterpriseId)).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(req.MemberIds) > 0 {
			isNotAllSearch = true
			db.Where("member_id in ?", req.MemberIds)
		}
		if len(req.StoreIds) > 0 {
			isNotAllSearch = true
			db.Where("store_id in ?", req.StoreIds)
		}
		if len(req.Ids) > 0 && req.Ids[0] > 0 {
			isNotAllSearch = true
			db.Where("id in ?", req.Ids)
		}
		if req.StartTime > 0 && req.EndTime > 0 {
			isNotAllSearch = true
			db.Where("created_at >= ? AND created_at <= ?", req.StartTime, req.EndTime)
		}
		if req.ReserveStartTime > 0 && req.ReserveEndTime > 0 {
			isNotAllSearch = true
			db.Where("reserve_time >= ? AND reserve_time <= ?", req.ReserveStartTime, req.ReserveEndTime)
		}
		if req.CommentStartTime > 0 && req.CommentEndTime > 0 {
			isNotAllSearch = true
			db.Where("comment_time >= ? AND comment_time <= ?", req.CommentStartTime, req.CommentEndTime)
		}
		if req.ReserveNumber != "" {
			isNotAllSearch = true
			db.Where("reserve_number = ?", req.ReserveNumber)
		}
		if req.OrderOrigin > 0 {
			isNotAllSearch = true
			db.Where("order_origin = ?", req.OrderOrigin)
		}
		if req.Status >= 0 {
			isNotAllSearch = true
			db.Where("status = ?", req.Status)
		}
		if req.IsComment > 0 {
			isNotAllSearch = true
			db.Where("is_comment = ?", req.IsComment)
		}
		if req.IsPay > 0 {
			isNotAllSearch = true
			db.Where("is_pay = ?", req.IsPay)
		}
		if req.IsStore >= 0 {
			isNotAllSearch = true
			db.Where("is_store = ?", req.IsStore)
		}
		if req.ServiceId > 0 {
			isNotAllSearch = true
			db.Where("reserve_id=?", req.ServiceId)
		}
		return db
	}).Where("reserve_time>0 and deleted_at=0")
	subQuery := query.Session(&gorm.Session{NewDB: true}).WithContext(c).Select("reserve_number")
	query.Scopes(func(db *gorm.DB) *gorm.DB {
		if req.ReserveOrderStartTime > 0 && req.ReserveOrderEndTime > 0 {
			db.Where("reserve_number in (?)", gormx.GetCkInstance().WithContext(c).Table(d.getScrmTableCk(microScrmModel.ReserveOrderSalesModel{}.TableName(), req.EnterpriseId)).Select("reserve_number").Where("pay_time BETWEEN ? AND ? ", req.ReserveOrderStartTime, req.ReserveOrderEndTime))
		}
		return db
	})
	query.Select("sum(invite_num) as invite_num,sum(is_open) as open_num,sum(is_commit) as commit_num").Find(&reserveSum)
	query.Where("status = ?", microScrmModel.FINISH).Count(&inviteSuccessNum)
	query.Where("is_store in ?", []int{microScrmModel.STORE_NORMAL, microScrmModel.STORE_YES}).Count(&paySuccessNum)
	if !isCoach { // 如果企业不是coach则，successPrice计算逻辑调整
		query.Where("status = ?", microScrmModel.FINISH).Select("COALESCE(sum(reserve_order_price),0) as success_price").Find(&successPrice)
	} else {
		salesQuery := gormx.GetCkInstance().WithContext(c).Table(d.getScrmTableCk(microScrmModel.ReserveOrderSalesModel{}.TableName(), req.EnterpriseId)).Select("COALESCE(sum(order_price),0) as success_price").Scopes(func(db *gorm.DB) *gorm.DB {
			if req.ReserveOrderStartTime > 0 && req.ReserveOrderEndTime > 0 {
				db.Where("pay_time BETWEEN ? AND ? ", req.ReserveOrderStartTime, req.ReserveOrderEndTime)
			}
			if isNotAllSearch {
				db.Where("reserve_number in (?)", subQuery)
			}
			return db
		})
		salesQuery.Find(&successPrice)
	}
	return
}

// GetOrderListsByStatus 查询昨天coach的邀约单数据给到安客诚(acxiom),推送进行中，邀约成功的，不是当天已取消的订单
// reserve_status 状态码为安客诚给到的订单匹配状态，1为匹配成功,2为取消匹配,3为邀约过期
func (d *ReservesOrder) GetOrderListsByStatus(c app.Ctx, enterpriseId, startTime, endTime int64) (result []microScrmModel.ReserveOrderModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("created_at BETWEEN ? AND ? AND reserve_time>0  ", startTime, endTime).Where("reserve_order_status=?", 0).Scopes(func(db *gorm.DB) *gorm.DB {
		db.Where("status in ? or (status=3 and updated_at>created_at+8600)", []int{microScrmModel.RUNNING, microScrmModel.WAITFINISH})
		return db
	}).Find(&result).Error
	return
}

// GetSuccessfulOrders 门店邀约数量(Successful_Appointment),当天日期-6的邀约成功记录
func (d *ReservesOrder) GetSuccessfulOrders(c app.Ctx, enterpriseId int64, startTime, endTime int64) (result []microScrmModel.ReserveOrderModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("reserve_order_time BETWEEN ? AND ?  ", startTime, endTime).Where("status=?", microScrmModel.FINISH).Find(&result).Error
	return
}

// GetStoreAppointmentLists Monthly的Store成功邀约文件，每个财月结束后上传(STORE_APPOINTMENTQTY)
func (d *ReservesOrder) GetStoreAppointmentLists(c ctx.Context, enterpriseId int64, startTime, endTime int64) (result []reserveRequest.AcxiomStoreAppointmentQTYInfo, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("created_at BETWEEN ? AND ? AND deleted_at=0  ", startTime, endTime).
		Select("`store_id`,`store_name` as `store_code`,count(1) as `appointment_qty`").
		Group("store_id").Scan(&result).Error
	return
}

// GetUserReserveOrderLists 获取用户的邀约单列表
func (d *ReservesOrder) GetUserReserveOrderLists(c app.Ctx, req reserveReq.UserReserveListReq) ([]microScrmModel.ReserveOrderModel, int64, error) {

	var (
		total            int64
		reserveOrderList []microScrmModel.ReserveOrderModel
	)
	enterpriseId := req.EnterpriseId
	page := req.Page
	size := req.Size

	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Scopes(func(db *gorm.DB) *gorm.DB {
		if req.ReserveTime > 0 {
			reserveStartTime := req.ReserveTime
			reserveEndTime := reserveStartTime + 86400
			if req.Month == 1 { // 查询指定月的时间
				year, MonthTime, _ := time.Unix(reserveStartTime, 10).Date()
				thisMonth := time.Date(year, MonthTime, 1, 0, 0, 0, 0, time.Local)
				reserveEndTimestr := thisMonth.AddDate(0, 1, -1).Format("2006-01-02") + " 23:59:59"
				reserveEndTime = helper.GetDateToTimes(reserveEndTimestr)
			}
			db.Where("reserve_time >= ? AND reserve_time <= ?", reserveStartTime, reserveEndTime)
		}
		if req.Status >= 0 {
			db.Where("status =?", req.Status)
		}
		return db
	}).Where("user_unionid=?", req.UnionId)
	err := query.Count(&total).
		Select("*").
		Order("created_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&reserveOrderList).
		Error
	return reserveOrderList, total, err
}

func (d *ReservesOrder) GetStoreReportDataLists(c ctx.Context, enterpriseId int64, startTime, endTime int64) (result []reserveRequest.AcxiomStoreReportDataInfo, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("created_at BETWEEN ? AND ? ", startTime, endTime).Where("reserve_time>0 and deleted_at=0").
		Select("`member_id`,`member_user_id`,`store_id`,`store_name` as `store_code`,count(1) as `appointment_qty`").
		Group("CONCAT(`store_id`,`member_id`)").Scan(&result).Error
	return
}

// GetAverage  获取邀约平均值相关统计
func (d *ReservesOrder) GetAverage(c app.Ctx, params reserveRequest.AchievementParams) (res int64, err error) {
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, params.EnterpriseId)).
		Where("member_id in ?", params.MemberIds).
		Where("reserve_time > 0")
	err = orm.Scopes(func(db *gorm.DB) *gorm.DB {
		switch params.Type {
		case excitationModels.RuleCategoryId2:
			db.Where("updated_at BETWEEN ? AND ?  ", params.StartTime, params.EndTime).
				Select("round(count(IF(is_confirm=?,true,null))/count(1) * 100) as res", microScrmModel.CONFIRM_YES)
		case excitationModels.RuleCategoryId3:
			db.Where("created_at BETWEEN ? AND ?  ", params.StartTime, params.EndTime).
				Select("round(count(IF(status=?,true,null))/count(1) * 100) as res", microScrmModel.FINISH)
		case excitationModels.RuleCategoryId5:
			db.Where("is_pay=?", microScrmModel.PAY_SUCCESS).
				Where("updated_at BETWEEN ? AND ?  ", params.StartTime, params.EndTime).
				Select("round(sum(reserve_order_price)) as res")

		}
		return db
	}).Pluck("res", &res).Error
	if err != nil {
		return 0, err
	}
	return
}

// GetAllMemberAchievement  获取所有导购邀约相关统计
func (d *ReservesOrder) GetAllMemberAchievement(c app.Ctx, params reserveRequest.AchievementParams) (map[int64]float64, error) {
	var resMap = make(map[int64]float64)
	var allRes []reserveRequest.AllRes
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, params.EnterpriseId)).
		Where("reserve_time > 0").
		Group("member_id")
	err := orm.Scopes(func(db *gorm.DB) *gorm.DB {
		switch params.Type {
		case excitationModels.RuleCategoryId2:
			db.Where("updated_at BETWEEN ? AND ?  ", params.StartTime, params.EndTime).
				Select("member_id,round(count(IF(is_confirm=?,true,null))/count(1) * 100) as all_total", microScrmModel.CONFIRM_YES)
		case excitationModels.RuleCategoryId3:
			db.Where("created_at BETWEEN ? AND ?  ", params.StartTime, params.EndTime).
				Select("member_id,round(count(IF(status=?,true,null))/count(1) * 100) as all_total", microScrmModel.FINISH)
		case excitationModels.RuleCategoryId5:
			db.Where("is_pay=?", microScrmModel.PAY_SUCCESS).
				Where("updated_at BETWEEN ? AND ?  ", params.StartTime, params.EndTime).
				Select("member_id,sum(reserve_order_price) as all_total")
		}
		return db
	}).Scan(&allRes).Error
	if err != nil {
		return nil, err
	}
	for _, tmp := range allRes {
		resMap[tmp.MemberId] = tmp.AllTotal
	}
	return resMap, nil
}

func (d *ReservesOrder) DeleteOrder(c ctx.Context, enterpriseId int64, orderBase reserveRequest.OrderBase) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"id":             orderBase.ID,
		"reserve_number": orderBase.ReserveNumber,
	}).Delete(&microScrmModel.ReserveOrderModel{}).Error
}

func (d *ReservesOrder) GetTopList(c app.Ctx, enterpriseId int64, inviteTop int, memberIds []int64, startTime, endTime int64) ([]models.TopList, error) {
	var result []models.TopList
	var lastTopTotal int64
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("SUM(reserve_order_price) AS total").
		Where("member_id in ?", memberIds).
		Where("status = 1 AND created_at BETWEEN ? AND ?", startTime, endTime).
		Group("member_id").
		Order("total desc")
	// 获取top以及和top最后一位数量相同的其他导购
	newOrm := orm.Session(&gorm.Session{NewDB: true}).WithContext(c)
	err := orm.Offset(inviteTop-1).Limit(1).Pluck("total", &lastTopTotal).Error
	err = newOrm.Select("member_id,SUM(reserve_order_price) AS total").
		Having("total > 0 AND total >= ?", lastTopTotal).Scan(&result).Error
	return result, err
}

// GetReserveTimeCount 获取到店邀约时间段的邀约数和门店的邀约数
func (d *ReservesOrder) GetReserveTimeCount(c ctx.Context, request reserveRequest.CheckReserveTime) (int64, int64, error) {
	var memberCount, storeCount int64
	var err error
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).Where(map[string]interface{}{
		"member_id":    request.MemberId,
		"reserve_time": request.ReserveTime,
		"status":       []int{microScrmModel.RUNNING, microScrmModel.WAITFINISH},
	}).Count(&memberCount).Error
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).Where(map[string]interface{}{"store_id": request.StoreId, "reserve_time": request.ReserveTime, "status": []int{microScrmModel.RUNNING, microScrmModel.WAITFINISH}}).Count(&storeCount).Error
	return memberCount, storeCount, err
}

func (d *ReservesOrder) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.ReserveOrderModel, error) {
	var result []microScrmModel.ReserveOrderModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Preload("ReserveService", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveServiceModel{}.TableName(), enterpriseId))
		}).
		Where(query, args...).Find(&result).Error
	return result, err
}

// GetReserveTimeOrderLists sw推送cdp的邀约单数据列表
func (d *ReservesOrder) GetReserveTimeOrderLists(c ctx.Context, enterpriseId int64, startTime int64, endTime int64) (result []microScrmModel.ReserveOrderModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("reserve_time BETWEEN ? AND ?  ", startTime, endTime).Where("reserve_order_status=?", 0).Scopes(func(db *gorm.DB) *gorm.DB {
		db.Where("status in ? or (status=3 and updated_at>created_at+8600)", []int{microScrmModel.RUNNING, microScrmModel.WAITFINISH})
		return db
	}).Find(&result).Error
	return
}

// GetCustomerList 获取用户的所有邀约单,只查询了邀约单列表，这里不需要记录查询邀约单总数
func (d *ReservesOrder) GetCustomerList(c app.Ctx, enterpriseId int64, page, size int, request reserveRequest.CustomerReserveListsReq) ([]microScrmModel.ReserveOrderModel, int64, error) {
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("reserve_time > 0 and deleted_at=0").Where("user_external_id=?", request.ExternalUserID)

	var (
		total     int64
		orderList []microScrmModel.ReserveOrderModel
		err       error
	)
	str := "FIELD(`status`,0,2,1,3,4) asc,reserve_time desc"
	err = query.Preload("ReserveService", func(db *gorm.DB) *gorm.DB {
		return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveServiceModel{}.TableName(), enterpriseId)).Unscoped()
	}).Preload("ReserveOrderSales", func(db *gorm.DB) *gorm.DB {
		return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveOrderSalesModel{}.TableName(), enterpriseId)).Preload("ReserveOrderGoods", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveOrderGoodsModel{}.TableName(), enterpriseId))
		})
	}).Order(str).Offset((page - 1) * size).Limit(size).Find(&orderList).Error

	return orderList, total, err
}

func (d *ReservesOrder) UpdateOrdersOverTime(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where(query, args...).
		Update("status", microScrmModel.OVERTIME).Error
}

// GetReserveOrderListNew 获取受邀统计单列表--下载中心
func (d *ReservesOrder) GetReserveOrderListNew(c app.Ctx, enterpriseId int64, page, size int, req reserveRequest.ReserveOrderListReq) ([]microScrmModel.ReserveOrderModel, int64, error) {
	query := gormx.GetCkInstance().Table(d.getScrmTableCk(d.tableName, enterpriseId))
	query = query.WithContext(c).Scopes(func(db *gorm.DB) *gorm.DB {
		if len(req.MemberIds) > 0 {
			db.Where("member_id in ?", req.MemberIds)
		}
		if len(req.StoreIds) > 0 {
			db.Where("store_id in ?", req.StoreIds)
		}
		if len(req.Ids) > 0 && req.Ids[0] > 0 {
			db.Where("id in ?", req.Ids)
		}
		if req.StartTime > 0 && req.EndTime > 0 {
			db.Where("created_at >= ? AND created_at <= ?", req.StartTime, req.EndTime)
		}
		if req.CommentStartTime > 0 && req.CommentEndTime > 0 {
			db.Where("comment_time >= ? AND comment_time <= ?", req.CommentStartTime, req.CommentEndTime)
		}
		if req.ReserveNumber != "" {
			db.Where("reserve_number = ?", req.ReserveNumber)
		}
		if req.OrderOrigin > 0 {
			db.Where("order_origin = ?", req.OrderOrigin)
		}
		if req.Status >= 0 {
			db.Where("status = ?", req.Status)
		}
		if req.IsComment > 0 {
			db.Where("is_comment = ?", req.IsComment)
		}
		if req.IsPay > 0 {
			db.Where("is_pay = ?", req.IsPay)
		}
		if req.IsStore >= 0 {
			db.Where("is_store = ?", req.IsStore)
		}
		return db
	}).Where("reserve_time>0 and deleted_at=0").Scopes(func(db *gorm.DB) *gorm.DB {
		if req.ReserveOrderStartTime > 0 && req.ReserveOrderEndTime > 0 {
			db.Where("reserve_number in (?)",
				gormx.GetCkInstance().WithContext(c).Table(d.getScrmTableCk(microScrmModel.ReserveOrderSalesModel{}.TableName(), enterpriseId)).Select("reserve_number").Where("pay_time BETWEEN ? AND ? ", req.ReserveOrderStartTime, req.ReserveOrderEndTime))
		}
		return db
	})
	var (
		total            int64
		reserveOrderList []microScrmModel.ReserveOrderModel
		err              error
	)

	err = query.Count(&total).Preload("ReserveOrderSales", func(db *gorm.DB) *gorm.DB {
		return db.Table(d.getScrmTableCk(microScrmModel.ReserveOrderSalesModel{}.TableName(), enterpriseId)).Preload("ReserveOrderGoods", func(db *gorm.DB) *gorm.DB {
			return db.Table(d.getScrmTableCk(microScrmModel.ReserveOrderGoodsModel{}.TableName(), enterpriseId))
		}).Scopes(func(db *gorm.DB) *gorm.DB {
			if req.ReserveOrderStartTime > 0 && req.ReserveOrderEndTime > 0 {
				db.Where("pay_time >= ? AND pay_time <= ?", req.ReserveOrderStartTime, req.ReserveOrderEndTime)
			}
			return db
		})
	}).Order("created_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&reserveOrderList).
		Error
	return reserveOrderList, total, err
}

// GetReserveByDays  获取1天内的邀约
func (d *ReservesOrder) GetReserveByDays(c app.Ctx, enterpriseId int64, lastId int64, size int) (list []microScrmModel.ReserveOrderModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("id,member_id,user_unionid,created_at,reserve_order_time").
		Where("updated_at > ?", time.Now().Unix()-86400).
		Where("reserve_time > 0").                       // 有到店时间的才算是邀约，没有到店时间说明是还未接受的邀请函
		Where("status <= ?", microScrmModel.WAITFINISH). // 进行中 已完成 都是进行中
		Where("deleted_at = 0").
		Where("id > ?", lastId).
		Order("id asc").
		Limit(size).
		Find(&list).Error
	return
}

// GetListCk 邀约单列表查询通过CK查询
func (d *ReservesOrder) GetListCk(c app.Ctx, enterpriseId int64, page, size int, req reserveRequest.OrderListReq) ([]microScrmModel.ReserveOrderModel, int64, error) {
	query := gormx.GetCkInstance().WithContext(c).Table(d.getScrmTableCk(d.tableName, enterpriseId))
	var (
		total     int64
		orderList []microScrmModel.ReserveOrderModel
		err       error
	)
	query.Where("reserve_order.reserve_time > 0 and reserve_order.deleted_at=0")
	if (req.ServiceId) > 0 {
		query.Where("reserve_id = ?", req.ServiceId)
	}
	if (req.Status) >= 0 {
		query.Where("status = ?", req.Status)
	}

	if len(req.MemberIds) > 0 {
		query.Where("member_id in ?", req.MemberIds)
	}
	if len(req.StoreIds) > 0 {
		query.Where("store_id in ?", req.StoreIds)
	}

	if len(req.Ids) > 0 {
		query.Where("id in ?", req.Ids)
	}
	if req.ReserveNumber != "" {
		query.Where("reserve_number = ?", req.ReserveNumber)
	}
	if req.ReserveStartTime > 0 && req.ReserveEndTime > 0 {
		query.Where("reserve_time >= ? AND reserve_time <= ?", req.ReserveStartTime, req.ReserveEndTime)

	}
	if req.StartTime > 0 && req.EndTime > 0 {
		query.Where("created_at >= ? AND created_at <= ?", req.StartTime, req.EndTime)

	}
	if req.IsStore >= 0 {
		query.Where("is_store = ?", req.IsStore)
	}
	if req.IsComment > 0 {
		query.Where("is_comment=?", req.IsComment)
	}
	if req.IsPay > 0 {
		query.Where("is_pay = ?", req.IsPay)
	}

	if req.Unionid != "" {
		query.Where("user_unionid = ?", req.Unionid)

	}
	if req.OrderOrigin > 0 {
		query.Where("order_origin = ?", req.OrderOrigin)
	}
	if req.Name != "" {
		query.Where("user_name like ?", "%"+req.Name+"%")
	}
	query.Scopes(func(db *gorm.DB) *gorm.DB {
		if req.ReserveOrderStartTime > 0 && req.ReserveOrderEndTime > 0 {
			db.Where("reserve_number in (?)", gormx.GetCkInstance().WithContext(c).Table(d.getScrmTableCk(microScrmModel.ReserveOrderSalesModel{}.TableName(), enterpriseId)).Select("reserve_number").Where("pay_time BETWEEN ? AND ? ", req.ReserveOrderStartTime, req.ReserveOrderEndTime))
		}
		return db
	})
	str := "created_at desc"
	err = query.Count(&total).Order(str).
		Offset((page - 1) * size).
		Limit(size).
		Scan(&orderList).
		Error

	return orderList, total, err
}

func (d *ReservesOrder) GetManyById(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.ReserveOrderModel, error) {
	var result []microScrmModel.ReserveOrderModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Preload("ReserveService", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveServiceModel{}.TableName(), enterpriseId))
		}).
		Where(query, args...).Order("id asc").Limit(1000).Find(&result).Error
	return result, err
}

// GetManyByIdPagination 通过id分页获取数据
func (d *ReservesOrder) GetManyByIdPagination(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.ReserveOrderModel, error) {
	var result []microScrmModel.ReserveOrderModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Preload("ReserveService", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.dao.BranchTable(microScrmModel.ReserveServiceModel{}.TableName(), enterpriseId))
		}).Where(query, args...).Order("id asc").Limit(1000).Find(&result).Error
	return result, err
}

func (d *ReservesOrder) BatchUpdateByUnionIds(c context.Context, enterpriseId int64, mappings []microScrmModel.UnionidMapModel, nowUnix int64) (int64, error) {
	if len(mappings) == 0 {
		return 0, nil
	}

	// 获取当前批次的最小和最大ID
	minID := mappings[0].ID
	maxID := mappings[len(mappings)-1].ID

	query := `
       UPDATE ` + fmt.Sprintf("%s%d.%s", app.Cfg().DB.DBPrefix, enterpriseId, d.tableName) + ` AS t1
        INNER JOIN ` + fmt.Sprintf("%s%d.%s", app.Cfg().DB.DBPrefix, enterpriseId, "unionid_map") + ` AS t2 ON t1.user_unionid  = t2.old_unionid
        SET t1.user_unionid = t2.new_unionid, t1.updated_at = ?
        WHERE t2.id >= ? AND t2.id <= ?`

	db := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Exec(query, nowUnix, minID, maxID)
	return db.RowsAffected, db.Error
}
