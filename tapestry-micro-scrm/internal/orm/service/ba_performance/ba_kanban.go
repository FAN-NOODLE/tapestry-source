package ba_performance

import (
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/service"
)

func (ba *BaPerformanceService) GetBaSaleDataByMobile(enterpriseID int64, mobiles []string, month string) (data micro_scrm_enterprise.BAMonthlyKanBanResult, err error) {
	// mainTable := fmt.Sprintf("%s%d.%s", micro_scrm_enterprise.BaPerformanceIndicatorModel{}.DBPrefix(), enterpriseID, micro_scrm_enterprise.BaMonthlyKanbanModel{}.TableName())
	// joinTable := fmt.Sprintf("%s%d.%s", micro_scrm_enterprise.BaPerformanceIndicatorModel{}.DBPrefix(), enterpriseID, micro_scrm_enterprise.BaPerformanceIndicatorModel{}.TableName())
	err = ba.Scopes(model.BranchTable(&micro_scrm_enterprise.BaMonthlyKanbanModel{}, enterpriseID)).
		Select(
			"month",
			// ----- 销售数据 -----
			"SUM(sales) AS sales",                // 销售额
			"SUM(monthly_sales) AS target_sales", // 销售目标
			"ROUND(ROUND(SUM(sales)/SUM(monthly_sales), 2)*100, 0) AS finish_rate",                          // 销售完成率
			"(SUM(monthly_sales)  - SUM(sales)) AS surplus_sales",                                           // 销售差额
			"ROUND(ROUND((SUM(sales)  - SUM(monthly_sales))/SUM(monthly_sales), 2)*100, 0) AS surplus_rate", // 销售完成率
			"ROUND(SUM(sales)/COUNT(order_num)) AS order_price",                                             // 客单价
			"SUM(order_num) AS order_num",                                                                   // 订单数

			"SUM(customer_num) AS customer_num", // 客人数
			// ----- 新客 -----
			"SUM(new_customer_num) AS new_customer_num",                                                             // 新客数
			"SUM(new_customer_sales) AS new_customer_sales",                                                         // 新客销售额
			"ROUND(ROUND((SUM(new_customer_sales)  - SUM(sales))/SUM(sales), 2)*100, 0) AS new_customer_sales_rate", // 新客销售占比
			// ----- 老客 -----
			"SUM(old_customer_num)",   // 老客数
			"SUM(old_customer_sales)", // 老客销售额
			"ROUND(ROUND((SUM(old_customer_sales)  - SUM(sales))/SUM(sales), 2)*100, 0) AS old_customer_sales_rate", // 老客销售占比
			// ----- 新二购 -----
			"SUM(new_second_customer_num) AS new_second_customer_num",                                                             // 新二购客数
			"SUM(new_second_customer_sales) AS new_second_customer_sales",                                                         // 新二购客销售额
			"ROUND(ROUND((SUM(new_second_customer_sales)  - SUM(sales))/SUM(sales), 2)*100, 0) AS new_second_customer_sales_rate", // 新二购客销售占比
			// ----- 游客 -----
			"SUM(visit_customer_num) AS visit_customer_num",                                                             // 游客数
			"SUM(visit_customer_sales) AS visit_customer_sales",                                                         // 游客销售额
			"ROUND(ROUND((SUM(visit_customer_sales)  - SUM(sales))/SUM(sales), 2)*100, 0) AS visit_customer_sales_rate", // 游客销售占比

			// ----- 运营数据 ----
			// ----- 新二购且一码三关注 -----
			"SUM(monthly_second_buy_and_three_scan) AS  monthly_second_buy_and_three_scan_num",                                      // 新二购指标人数
			"SUM(new_second_and_three_scan) AS new_second_and_three_scan_num",                                                       // 一码三关注人数
			"ROUND(ROUND(SUM(new_second_and_three_scan)/SUM(new_second_customer_num), 2)*100, 0) AS new_second_and_three_scan_rate", // 一码三关注率
			// ----- 老客且一码三关注 -----
			"SUM(monthly_old_customer_and_three_scan) AS monthly_old_customer_and_three_scan_num",              // 老客指标人数
			"SUM(old_and_three_scan) AS old_and_three_scan_num",                                                // 一码三关注人数
			"ROUND(ROUND(SUM(old_and_three_scan)/SUM(old_customer_num), 2)*100, 0) AS old_and_three_scan_rate", // 一码三关注率
			// ----- 新客且一码三关注 -----
			"SUM(monthly_new_customer_and_three_scan) AS  monthly_new_customer_and_three_scan_num",                               // 新客指标人数
			"SUM(new_customer_and_three_scan) AS new_customer_and_three_scan_num",                                                // 一码三关注人数
			"ROUND(ROUND(SUM(new_customer_and_three_scan)/SUM(new_customer_num), 2)*100, 0) AS new_customer_and_three_scan_rate", // 一码三关注率
		).
		Where("mobile IN ?", mobiles).
		Where("month = ?", month).
		Where("deleted_at=?", 0).
		Group("month").
		Take(&data).
		Error

	return
}

func (ba *BaPerformanceService) GetBaSaleDataByMobileList(enterpriseID int64, mobiles []string, month string, param service.PaginateParam) (data []micro_scrm_enterprise.BAMonthlyKanBanResult, err error) {
	// mainTable := fmt.Sprintf("%s%d.%s", micro_scrm_enterprise.BaPerformanceIndicatorModel{}.DBPrefix(), enterpriseID, micro_scrm_enterprise.BaMonthlyKanbanModel{}.TableName())
	// joinTable := fmt.Sprintf("%s%d.%s", micro_scrm_enterprise.BaPerformanceIndicatorModel{}.DBPrefix(), enterpriseID, micro_scrm_enterprise.BaPerformanceIndicatorModel{}.TableName())
	err = ba.Scopes(model.BranchTable(&micro_scrm_enterprise.BaMonthlyKanbanModel{}, enterpriseID)).
		Select(
			"ba_name AS name",
			// ----- 销售数据 -----
			"SUM(sales) AS sales",                // 销售额
			"SUM(monthly_sales) AS target_sales", // 销售目标
			"ROUND(ROUND(SUM(sales)/SUM(monthly_sales), 2)*100, 0) AS finish_rate",                          // 销售完成率
			"(SUM(monthly_sales)  - SUM(sales)) AS surplus_sales",                                           // 销售差额
			"ROUND(ROUND((SUM(sales)  - SUM(monthly_sales))/SUM(monthly_sales), 2)*100, 0) AS surplus_rate", // 销售完成率
			"ROUND(SUM(sales)/COUNT(order_num)) AS order_price",                                             // 客单价
			"SUM(order_num) AS order_num",                                                                   // 订单数

			"SUM(customer_num) AS customer_num", // 客人数
			// ----- 新客 -----
			"SUM(new_customer_num) AS new_customer_num",                                                             // 新客数
			"SUM(new_customer_sales) AS new_customer_sales",                                                         // 新客销售额
			"ROUND(ROUND((SUM(new_customer_sales)  - SUM(sales))/SUM(sales), 2)*100, 0) AS new_customer_sales_rate", // 新客销售占比
			// ----- 老客 -----
			"SUM(old_customer_num)",   // 老客数
			"SUM(old_customer_sales)", // 老客销售额
			"ROUND(ROUND((SUM(old_customer_sales)  - SUM(sales))/SUM(sales), 2)*100, 0) AS old_customer_sales_rate", // 老客销售占比
			// ----- 新二购 -----
			"SUM(new_second_customer_num) AS new_second_customer_num",                                                             // 新二购客数
			"SUM(new_second_customer_sales) AS new_second_customer_sales",                                                         // 新二购客销售额
			"ROUND(ROUND((SUM(new_second_customer_sales)  - SUM(sales))/SUM(sales), 2)*100, 0) AS new_second_customer_sales_rate", // 新二购客销售占比
			// ----- 游客 -----
			"SUM(visit_customer_num) AS visit_customer_num",                                                             // 游客数
			"SUM(visit_customer_sales) AS visit_customer_sales",                                                         // 游客销售额
			"ROUND(ROUND((SUM(visit_customer_sales)  - SUM(sales))/SUM(sales), 2)*100, 0) AS visit_customer_sales_rate", // 游客销售占比

			// ----- 运营数据 ----
			// ----- 新二购且一码三关注 -----
			"SUM(monthly_second_buy_and_three_scan) AS  monthly_second_buy_and_three_scan_num",                                                       // 新二购指标人数
			"ROUND(ROUND(SUM(new_second_and_three_scan)/SUM(monthly_second_buy_and_three_scan), 2)*100, 0) AS new_second_and_three_scan_finish_rate", // 完成指标率
			"SUM(new_second_and_three_scan) AS new_second_and_three_scan_num",                                                                        // 一码三关注人数
			"ROUND(ROUND(SUM(new_second_and_three_scan)/SUM(new_second_customer_num), 2)*100, 0) AS new_second_and_three_scan_rate",                  // 一码三关注率
			// ----- 老客且一码三关注 -----
			"SUM(monthly_old_customer_and_three_scan) AS monthly_old_customer_and_three_scan",                                            // 老客指标人数
			"ROUND(ROUND(SUM(old_and_three_scan)/SUM(monthly_old_customer_and_three_scan), 2)*100, 0) AS old_and_three_scan_finish_rate", // 完成指标率
			"SUM(old_and_three_scan) AS old_and_three_scan_num",                                                                          // 一码三关注人数
			"ROUND(ROUND(SUM(old_and_three_scan)/SUM(old_customer_num), 2)*100, 0) AS old_and_three_scan_rate",                           // 一码三关注率
			// ----- 新客且一码三关注 -----
			"SUM(monthly_new_customer_and_three_scan) AS  monthly_new_customer_and_three_scan_num",                                                         // 新客指标人数
			"ROUND(ROUND(SUM(new_customer_and_three_scan)/SUM(monthly_new_customer_and_three_scan), 2)*100, 0) AS new_customer_and_three_scan_finish_rate", // 完成指标率
			"SUM(new_customer_and_three_scan) AS new_customer_and_three_scan_num",                                                                          // 一码三关注人数
			"ROUND(ROUND(SUM(new_customer_and_three_scan)/SUM(new_customer_num), 2)*100, 0) AS new_customer_and_three_scan_rate",                           // 一码三关注率
		).
		Where("mobile IN ?", mobiles).
		Where("month = ?", month).
		Where("deleted_at=?", 0).
		Group("month, mobile, ba_name").
		Find(&data).
		Error

	return
}

// BaDataCount 计算需要统计的BA的数量
func (ba *BaPerformanceService) BaDataCount(enterpriseID int64, mobiles []string, month string) (count int64, err error) {
	err = ba.Scopes(model.BranchTable(&micro_scrm_enterprise.BaMonthlyKanbanModel{}, enterpriseID)).
		Scopes(service.ScopeCond_EQDelete).
		Where("mobile IN ?", mobiles).
		Where("month = ?", month).
		Where("deleted_at=?", 0).
		Count(&count).Error

	return
}

// GetBATotalSales 获取BA某月的销售总额
func (ba *BaPerformanceService) GetBATotalSales(enterpriseID int64, mobiles []string, month string) (sales float32, err error) {
	err = ba.Scopes(model.BranchTable(&micro_scrm_enterprise.BaMonthlyKanbanModel{}, enterpriseID)).
		Select("SUM(sales) AS sales").
		Where("mobile IN ?", mobiles).
		Where("month = ?", month).
		Where("deleted_at=?", 0).
		Scan(&sales).
		Error

	return
}

func (ba *BaPerformanceService) BaSourceMaxDept(enterpriseID int64, mobile string) (baLevel int, err error) {
	err = ba.Scopes(model.BranchTable(&micro_scrm_enterprise.BaKanbanDictionaryModel{}, enterpriseID)).
		Select("IFNULL(MAX(ba_level),0) ba_level").
		Where("ba_mobile=?", mobile, 0).
		Scan(&baLevel).
		Error

	return
}
func (ba *BaPerformanceService) BaKanbaninfo(enterpriseID int64, baLevel int, mobile string) (info []micro_scrm_enterprise.BaKanbanDictionaryModel, err error) {
	query := ba.Scopes(model.BranchTable(&micro_scrm_enterprise.BaKanbanDictionaryModel{}, enterpriseID))
	if baLevel == 1 {
		query.Where("ba_level=? AND ba_mobile=? AND deleted_at=? ", baLevel, mobile, 0)
	} else if baLevel == 2 {
		query.Where("ba_level=? AND manager_mobile=? AND deleted_at=? ", 1, mobile, 0)
	} else if baLevel == 3 {
		query.Where(" ba_level=? AND partition_leader_mobile=? AND deleted_at=? ", 2, mobile, 0)
	} else if baLevel == 4 {
		query.Where("ba_level=? AND regional_leader_mobie=? AND deleted_at=? ", 3, mobile, 0)
	}
	err = query.Scan(&info).Error
	return
}
func (ba *BaPerformanceService) FindBAkanBanDateBymobie(enterpriseID int64, baLevel int, mobile, month string) (info micro_scrm_enterprise.BaMonthlyKanbanModel, err error) {
	err = ba.Scopes(model.BranchTable(&info, enterpriseID)).
		Where("sales_type=? AND mobile=?  AND month=? AND deleted_at=? ", baLevel, mobile, month, 0).
		First(&info).Error
	return
}
func (ba *BaPerformanceService) BaKanbanBaMobiles(enterpriseID int64, baLevel int, mobile string) (mobiles []string, err error) {
	query := ba.Scopes(model.BranchTable(&micro_scrm_enterprise.BaKanbanDictionaryModel{}, enterpriseID)).
		Select("ba_mobile")
	if baLevel == 1 {
		query.Where("ba_mobile=? AND deleted_at=? ", mobile, 0)
	} else if baLevel == 2 {
		query.Where("manager_mobile=? AND deleted_at=? ", mobile, 0)
	} else if baLevel == 3 {
		query.Where(" partition_leader_mobile=? AND deleted_at=? ", mobile, 0)
	} else if baLevel == 4 {
		query.Where(" regional_leader_mobie=? AND deleted_at=? ", mobile, 0)
	}
	err = query.Scan(&mobiles).Error
	return
}

func (ba *BaPerformanceService) FindBAkanBanDateBymobies(enterpriseID int64, month string, mobiles []string, page, pageSize int) (info []micro_scrm_enterprise.BaMonthlyKanbanModel, total int64) {
	query := ba.Scopes(model.BranchTable(&micro_scrm_enterprise.BaMonthlyKanbanModel{}, enterpriseID)).
		Where(" mobile IN ?  AND month=? AND sales_type=? AND  deleted_at=? ", mobiles, month, 1, 0).
		Count(&total).Order("sales DESC,ba_name")
	if page != 0 && pageSize != 0 {
		query.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	query.Scan(&info)
	return
}
