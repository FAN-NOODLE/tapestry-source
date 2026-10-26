package membersDao

import (
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	memberReq "micro-scrm/internal/app/http/request/member"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type MemberInvitationList struct {
	dao       dao.Dao
	tableName string
}

func NewMemberInvitationList() *MemberInvitationList {
	return &MemberInvitationList{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.ReserveOrderModel{}.TableName(),
	}
}

// 查询条件待确认 按照到店时间统计

func (d MemberInvitationList) Invitation(c app.Ctx, request memberReq.PerformanceListInvitationReq, membersID []int64, storeCodes []string) ([]memberReq.PerformanceListReserveOrderRes, int64, error) {
	var (
		total  int64
		result []memberReq.PerformanceListReserveOrderRes
		err    error
	)
	offset := (request.Page - 1) * request.Size

	startTime := request.FinancialTime[0]
	endTime := request.FinancialTime[1]

	sqlSelect := fmt.Sprintf("reserve_order.member_id,count(IF(reserve_order.created_at BETWEEN %d AND %d,reserve_order.reserve_number,NULL)) as num,"+
		"count(IF(reserve_order.is_pay=1 and reserve_order_sales.pay_time BETWEEN %d AND %d,reserve_order.reserve_number,NULL)) as success_num,"+
		"sum(IF(reserve_order.is_pay=1 and reserve_order_sales.pay_time BETWEEN %d AND %d,reserve_order_sales.order_price,0)) as price_total",
		startTime, endTime, startTime, endTime, startTime, endTime)
	salesTable := d.dao.GetTable(microScrmModel.ReserveOrderSalesModel{}.TableName(), request.EnterpriseId)
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).
		Joins("left join "+salesTable+" ON reserve_order.reserve_number=reserve_order_sales.reserve_number").
		Select(sqlSelect).
		Where("(reserve_order.created_at BETWEEN ? AND ? or reserve_order_sales.pay_time BETWEEN ? AND ? ) and member_id > 0 ", startTime, endTime, startTime, endTime).
		Group("member_id").
		Order("price_total desc, num desc, success_num desc, member_id asc")

	// 如果是明星榜单就不传递导购id
	if len(membersID) > 0 {
		query.Where("(reserve_order.member_id in (?) or reserve_order_sales.store_code in (?))", membersID, storeCodes)
	}
	err = query.Offset(offset).Limit(request.Size).Find(&result).Error

	if len(membersID) > 0 {
		query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).Joins("LEFT JOIN "+salesTable+" ON reserve_order.reserve_number=reserve_order_sales.reserve_number").
			Where("(reserve_order.created_at BETWEEN ? AND ? OR reserve_order_sales.pay_time BETWEEN ? AND ?) and (reserve_order.member_id in (?) or reserve_order_sales.store_code in (?))", startTime, endTime, startTime, endTime, membersID, storeCodes).
			Select("reserve_order.member_id").Group("reserve_order.member_id")
		err = d.dao.WithContext(c).Table("(?) as a", query).Count(&total).Error
		// sqlMembersID := strings.Trim(strings.Replace(fmt.Sprint(membersID), " ", ",", -1), "[]")
		// sqlstoreCodes := strings.Join(storeCodes, `", "`)
		// sqlstoreCodes = `"` + sqlstoreCodes + `"`
		// sql := fmt.Sprintf("SELECT count(*) as total FROM (select member_id FROM micro_scrm_enterprise_%d.reserve_order left join micro_scrm_enterprise_%d.reserve_order_sales ON reserve_order.reserve_number=reserve_order_sales.reserve_number WHERE (reserve_order.created_at BETWEEN %d AND %d OR "+
		// 	"reserve_order_sales.pay_time BETWEEN %d AND %d) and (member_id in (%s) or reserve_order_sales.store_code in (%s)) GROUP By member_id) a", request.EnterpriseId, request.EnterpriseId, startTime, endTime, startTime, endTime, sqlMembersID, sqlstoreCodes)
		// d.dao.WithContext(c).Raw(sql).Scan(&total)
	}

	return result, total, err
}

// InvitationMemberIds 在规定时间内导购的
func (d MemberInvitationList) InvitationMemberIds(c app.Ctx, request memberReq.PerformanceListInvitationReq, membersID []int64, storeCodes []string) ([]int64, error) {
	var membersId []int64
	startTime := request.FinancialTime[0]
	endTime := request.FinancialTime[1]
	salesTable := d.dao.GetTable(microScrmModel.ReserveOrderSalesModel{}.TableName(), request.EnterpriseId)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).
		Joins("LEFT JOIN "+salesTable+" ON reserve_order.reserve_number=reserve_order_sales.reserve_number").
		Where("(reserve_order.created_at BETWEEN ? AND ? or reserve_order_sales.pay_time BETWEEN ? AND ? ) and member_id > 0", startTime, endTime, startTime, endTime).
		Where("(member_id in (?) or reserve_order_sales.store_code in (?))", membersID, storeCodes).
		Select("member_id").
		Group("member_id").Find(&membersId).Error
	return membersId, err
}
