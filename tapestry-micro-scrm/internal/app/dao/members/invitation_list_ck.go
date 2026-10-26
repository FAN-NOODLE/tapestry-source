package membersDao

import (
	"fmt"
	"micro-scrm/internal/app"
	memberReq "micro-scrm/internal/app/http/request/member"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

func (d MemberInvitationList) InvitationCk(c app.Ctx, request memberReq.PerformanceListInvitationReq, membersID []int64, storeCodes []string) ([]memberReq.PerformanceListReserveOrderRes, int64, error) {
	var (
		total  int64
		result []memberReq.PerformanceListReserveOrderRes
		err    error
	)
	offset := (request.Page - 1) * request.Size

	startTime := request.FinancialTime[0]
	endTime := request.FinancialTime[1]

	// 设置库名和表名
	table := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", request.EnterpriseId, d.tableName)
	salesTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", request.EnterpriseId, microScrmModel.ReserveOrderSalesModel{}.TableName())
	// sqlstoreCodes := strings.Join(storeCodes, `", "`)
	// sqlstoreCodes = `"` + sqlstoreCodes + `"`
	// sqlSelect := fmt.Sprintf("reserve_order.member_id,count(DISTINCT IF(reserve_order.created_at BETWEEN %d AND %d,reserve_order.reserve_number,NULL)) as num,"+
	// 	"count(DISTINCT IF(reserve_order.is_pay=1 and reserve_order_sales.pay_time BETWEEN %d AND %d,reserve_order.reserve_number,NULL)) as success_num,"+
	// 	"sum(IF(reserve_order.is_pay=1 and reserve_order_sales.pay_time BETWEEN %d AND %d AND reserve_order_sales.store_code in (%s),reserve_order_sales.order_price,0)) as price_total",
	// 	startTime, endTime, startTime, endTime, startTime, endTime, sqlstoreCodes)
	query := gormx.GetCkInstance().WithContext(c).Table(table)

	if len(storeCodes) > 0 {
		query.Select("reserve_order.member_id,count(DISTINCT IF(reserve_order.created_at BETWEEN ? AND ? AND store_id in (?) ,reserve_order.reserve_number,NULL)) as num,"+
			"count(DISTINCT IF(reserve_order.status=1 and reserve_order_sales.pay_time BETWEEN ? AND ? ,reserve_order.reserve_number,NULL)) as success_num,"+
			"sum(IF(reserve_order.status=1 and reserve_order_sales.pay_time BETWEEN ? AND ? AND reserve_order_sales.store_code in (?),reserve_order_sales.order_price,0)) as price_total", startTime, endTime, request.ThirdId, startTime, endTime, startTime, endTime, storeCodes)
	} else {
		query.Select("reserve_order.member_id,count(DISTINCT IF(reserve_order.created_at BETWEEN ? AND ? ,reserve_order.reserve_number,NULL)) as num,"+
			"count(DISTINCT IF(reserve_order.status=1 and reserve_order_sales.pay_time BETWEEN ? AND ? ,reserve_order.reserve_number,NULL)) as success_num,"+
			"sum(IF(reserve_order.status=1 and reserve_order_sales.pay_time BETWEEN ? AND ? ,reserve_order_sales.order_price,0)) as price_total", startTime, endTime, startTime, endTime, startTime, endTime)
	}
	query.Joins("LEFT JOIN " + salesTable + " ON reserve_order.reserve_number=reserve_order_sales.reserve_number")
	query.Where("member_id>0")
	if len(request.ThirdId) > 0 { // 兼容处理导购换店，查询A门店的邀约单时，导购不在A门店的邀约单，也许展示
		query.Where("((reserve_order.created_at BETWEEN ? AND ? AND (reserve_order.store_id in (?) or reserve_order.member_id in (?)) ) or (reserve_order_sales.pay_time BETWEEN ? AND ? AND reserve_order_sales.store_code in (?))) ", startTime, endTime, request.ThirdId, membersID, startTime, endTime, storeCodes)
	}
	if len(membersID) > 0 && len(request.ThirdId) == 0 { // 该逻辑兼容明星榜单
		query.Where("((reserve_order.created_at BETWEEN ? AND ? ) or (reserve_order_sales.pay_time BETWEEN ? AND ? ))", startTime, endTime, startTime, endTime).Where("(reserve_order.member_id in (?) or reserve_order_sales.store_code in (?))", membersID, storeCodes)
	}

	query.Group("reserve_order.member_id").
		Order("price_total desc, num desc, success_num desc, member_id asc")

	// 如果是明星榜单就不传递导购id

	err = query.Offset(offset).Limit(request.Size).Find(&result).Error

	if len(membersID) > 0 {
		query := gormx.GetCkInstance().WithContext(c).Table(table).Joins("LEFT JOIN "+salesTable+" ON reserve_order.reserve_number=reserve_order_sales.reserve_number").
			Where("(reserve_order.created_at BETWEEN ? AND ? OR reserve_order_sales.pay_time BETWEEN ? AND ?) and (reserve_order.member_id in (?) or reserve_order_sales.store_code in (?))", startTime, endTime, startTime, endTime, membersID, storeCodes).
			Select("reserve_order.member_id").Group("reserve_order.member_id")
		err = gormx.GetCkInstance().WithContext(c).Table("(?) as a", query).Count(&total).Error
		// sqlMembersID := strings.Trim(strings.Replace(fmt.Sprint(membersID), " ", ",", -1), "[]")
		// sqlstoreCodes := strings.Join(storeCodes, `", "`)
		// sqlstoreCodes = `"` + sqlstoreCodes + `"`
		// sql := fmt.Sprintf("SELECT count(*) as total FROM (select reserve_order.member_id FROM %s  LEFT JOIN %s ON reserve_order.reserve_number=reserve_order_sales.reserve_number  WHERE (reserve_order.created_at BETWEEN %d AND %d OR "+
		// 	"reserve_order_sales.pay_time BETWEEN %d AND %d) and (reserve_order.member_id in (%s) or reserve_order_sales.store_code in (%s)) GROUP By reserve_order.member_id) a", table, salesTable, startTime, endTime, startTime, endTime, sqlMembersID, sqlstoreCodes)
		// gormx.GetCkInstance().WithContext(c).Raw(sql).Scan(&total)
	}

	return result, total, err
}

// InvitationMemberIdsCk 在规定时间内导购的
func (d MemberInvitationList) InvitationMemberIdsCk(c app.Ctx, request memberReq.PerformanceListInvitationReq, membersID []int64, storeCodes []string) ([]int64, error) {
	var membersId []int64
	startTime := request.FinancialTime[0]
	endTime := request.FinancialTime[1]

	table := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", request.EnterpriseId, d.tableName)
	salesTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", request.EnterpriseId, microScrmModel.ReserveOrderSalesModel{}.TableName())

	err := gormx.GetCkInstance().WithContext(c).Table(table).
		Joins("LEFT JOIN "+salesTable+" ON reserve_order.reserve_number=reserve_order_sales.reserve_number").
		Where("(reserve_order.created_at BETWEEN ? AND ? or reserve_order_sales.pay_time BETWEEN ? AND ? ) and member_id > 0", startTime, endTime, startTime, endTime).
		Where("(reserve_order.member_id in (?) or reserve_order_sales.store_code in (?))", membersID, storeCodes).
		Select("reserve_order.member_id").
		Group("reserve_order.member_id").Find(&membersId).Error
	return membersId, err
}
