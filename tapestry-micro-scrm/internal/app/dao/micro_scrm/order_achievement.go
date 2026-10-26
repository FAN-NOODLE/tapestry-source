package microScrmDao

import (
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	memberReq "micro-scrm/internal/app/http/request/member"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type OrderAchievement struct {
	dao       dao.Dao
	tableName string
}

func NewOrderAchievement() *OrderAchievement {
	return &OrderAchievement{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.OrderAchievementModel{}.TableName(),
	}
}

func (d OrderAchievement) Achievement(c app.Ctx, request memberReq.PerformanceListAchievementReq, membersID []int64) ([]memberReq.PerformanceListAchievementRes, int64, error) {
	var (
		total  int64
		result []memberReq.PerformanceListAchievementRes
		err    error
	)
	offset := (request.Page - 1) * request.Size

	startTime := request.FinancialTime[0]
	endTime := request.FinancialTime[1]
	sqlSelect := fmt.Sprintf("order_achievement.member_id,sum(order_total) as order_total," +
		"sum(order_price_total) as order_price_total," +
		"sum(pay_order_total) as pay_order_total," +
		"sum(refund_order_total) as refund_order_total")
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).
		Select(sqlSelect).
		Where("order_time BETWEEN ? AND ?", startTime, endTime).
		Group("member_id").
		Order("order_price_total desc, order_total desc, member_id asc")
	//如果是明星榜单就不传递导购id
	if len(membersID) > 0 {
		query.Where("member_id in (?)", membersID)
		query.Count(&total)
	}
	err = query.Offset(offset).Limit(request.Size).Find(&result).Error

	return result, total, err
}

// InsertAchievement 业绩信息写入
func (d OrderAchievement) InsertAchievement(c app.Ctx, enterpriseId int64, data microScrmModel.OrderAchievementModel) bool {

	var result microScrmModel.OrderAchievementModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("id").
		Where("userid = ? and store_code = ? and order_time = ? ", data.Userid, data.StoreCode, data.OrderTime).
		Find(&result).Error
	if err != nil {
		return false
	}
	if result.ID > 0 {
		return false
	}

	data.ID = app.UID()
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
	if err != nil {
		return false
	}
	return true
}

// InsertAchievementMemberIds 在规定时间内导购的id
func (d OrderAchievement) InsertAchievementMemberIds(c app.Ctx, request memberReq.PerformanceListAchievementReq, membersID []int64) ([]int64, error) {
	var membersId []int64
	startTime := request.FinancialTime[0]
	endTime := request.FinancialTime[1]

	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).
		Where("order_time BETWEEN ? AND ?", startTime, endTime).
		Where("member_id in (?)", membersID).
		Group("member_id").Pluck("member_id", &membersId).Error
	return membersId, err
}
