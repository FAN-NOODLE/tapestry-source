package transfer_customer

import (
	"fmt"
	"time"

	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/transfer_customer"
	ormService "micro-scrm/internal/orm/service"
)

// TasksSummaryData 汇总任务的数据信息
func (service *TransferCustomerService) TasksSummaryData(enterpriseID, taskID int64) (data transfer_customer.TransferCustomerDataModel_TasksReportTitle, err error) {
	err = service.Table("(?) AS M", service.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, enterpriseID)).
		Select("`external_unionid` AS union_id", "IF(`transfer_customer_status` = 1, 1, 0) AS success_flag",
			"IF(`transfer_customer_status` = 3 || `transfer_customer_status` = 4, 1, 0) AS fail_flag",
			"`takeover_userid` AS ba_user_id").
		Where(map[string]interface{}{
			"task_id": taskID,
		}).Group("external_unionid")).
		Select("COUNT(1) AS customer_total_num", "SUM(success_flag) AS customer_success_num",
			"SUM(fail_flag) AS customer_fail_num", "COUNT(DISTINCT ba_user_id) AS finish_ba_num").
		Scan(&data).Error
	return
}

// TasksTransferCustomerList 根据任务ID，查询继承客户列表
func (service *TransferCustomerService) TasksTransferCustomerList(enterpriseID, taskID int64, param ormService.PaginateParam) (data []transfer_customer.TransferCustomerDataModel_TasksReportListData, err error) {
	brandDBPrefix := fmt.Sprintf("%s%d", model.BranchModel{}.DBPrefix(), enterpriseID)
	mainTable := fmt.Sprintf("%s.%s", brandDBPrefix, transfer_customer.TransferCustomerDataModel{}.TableName())
	wuceTable := fmt.Sprintf("%s.%s", brandDBPrefix, micro_scrm_enterprise.WxUserCorpExternalModel{}.TableName())
	mTable := fmt.Sprintf("%s.%s", brandDBPrefix, micro_scrm_enterprise.MembersModel{}.TableName())
	dTable := micro_scrm.DepartmentsModel{}.TableName()

	fields := []string{
		fmt.Sprintf("%s.name", wuceTable),                                        // 客户名称
		fmt.Sprintf("COUNT(%s.id) AS friend_num", mainTable),                     // 好友数
		fmt.Sprintf("%s.name AS current_admin_ba_name", mTable),                  // 当前管理BA
		fmt.Sprintf("%s.name AS ba_department_name", dTable),                     // 所属部门
		fmt.Sprintf("%s.transfer_customer_status AS transfer_status", mainTable), // 继承状态
		fmt.Sprintf("CASE %s.transfer_customer_status "+
			"WHEN 0 THEN '未开始' "+
			"WHEN 1 THEN '已完成' "+
			"WHEN 2 THEN '等待接替' "+
			"WHEN 5 THEN '无接替记录' "+
			"WHEN 6 THEN '接替中' "+
			"ELSE '继承失败' "+
			"END AS transfer_status_msg", mainTable), // 继承状态文本
		fmt.Sprintf("CASE %s.transfer_customer_status "+
			"WHEN 3 THEN '客户拒绝' "+
			"WHEN 4 THEN '接替成员客户达到上限' "+
			"ELSE '/' "+
			"END AS fail_msg", mainTable), // 失败明文
	}
	err = service.Table(mainTable).Select(fields).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s", wuceTable, wuceTable, "unionid", mainTable, "external_unionid")).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s", mTable, mTable, "userid", mainTable, "takeover_userid")).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s", dTable, dTable, "third_id", mTable, "main_department")).
		Group(fmt.Sprintf("%s.%s", mainTable, "external_unionid")).
		Where(fmt.Sprintf("%s.task_id = ?", mainTable), taskID).
		Scopes(ormService.ScopeCond_Paginate(param)).
		Scan(&data).Error

	return
}

// TasksCheck 检测任务是否存在
func (service *TransferCustomerService) TasksCheck(enterpriseID int64, taskName string) (count int64, err error) {
	err = service.Scopes(model.BranchTable(&transfer_customer.TransferCustomerTasksModel{}, enterpriseID)).
		Where("name = ?", taskName).
		Count(&count).Error
	return
}

// TasksCreate 新增任务
func (service *TransferCustomerService) TasksCreate(enterpriseID int64, data transfer_customer.TransferCustomerTasksModel) error {
	return service.Scopes(model.BranchTable(&data, enterpriseID)).
		Create(&data).Error
}

// TasksUpdate 新增任务
func (service *TransferCustomerService) TasksUpdate(enterpriseID int64, data transfer_customer.TransferCustomerTasksModel) error {
	return service.Scopes(model.BranchTable(&data, enterpriseID)).
		Updates(data).Error
}

// TasksGet 获取任务详情
func (service *TransferCustomerService) TasksGet(enterpriseID int64, cond transfer_customer.TransferCustomerTasksModel) (result transfer_customer.TransferCustomerTasksModel, err error) {
	err = service.Scopes(model.BranchTable(&result, enterpriseID)).
		Scopes(ormService.ScopeCond_EQDelete).
		Where(&cond).
		First(&result).Error
	return
}

// TasksGetList 获取品牌任务集合
func (service *TransferCustomerService) TasksGetList(enterpriseID, startTime, endTime int64, status uint8, paginate ormService.PaginateParam) (result []transfer_customer.TransferCustomerTasksModel, err error) {
	m := service.Scopes(ormService.ScopeCond_EQDelete).
		Scopes(model.BranchTable(&transfer_customer.TransferCustomerTasksModel{}, enterpriseID))
	if status != transfer_customer.TCT_STATUS_ZERO {
		m.Where("status = ?", status)
	}
	if startTime != 0 {
		m.Where("created_at >= ?", startTime)
	}
	if endTime != 0 {
		m.Where("created_at <= ?", endTime)
	}

	err = m.Scopes(ormService.ScopeCond_Paginate(paginate)).
		Order("created_at DESC").
		Find(&result).Error

	return
}

// TasksGetListCount 获取品牌任务的总数
func (service *TransferCustomerService) TasksGetListCount(enterpriseID, startTime, endTime int64, status uint8) (count int64, err error) {
	m := service.Scopes(ormService.ScopeCond_EQDelete).
		Scopes(model.BranchTable(&transfer_customer.TransferCustomerTasksModel{}, enterpriseID))
	if status != transfer_customer.TCT_STATUS_ZERO {
		m.Where("status = ?", status)
	}
	if startTime != 0 {
		m.Where("start_time >= ?", startTime)
	}
	if endTime != 0 {
		m.Where("end_time <= ?", endTime)
	}
	err = m.Count(&count).Error
	return
}

// TasksGetListCount 获取数据已经拉完正在执行的计划
func (service *TransferCustomerService) TasksPullDataStatus(enterpriseID int64, paginate ormService.PaginateParam) (result []transfer_customer.TransferCustomerTasksModel, err error) {
	err = service.Scopes(model.BranchTable(&transfer_customer.TransferCustomerTasksModel{}, enterpriseID),
		ormService.ScopeCond_EQDelete, ormService.ScopeCond_Paginate(paginate)).
		Where("status = ? AND pull_data_status=?", 2, 1).Find(&result).Error
	return
}
func (service *TransferCustomerService) FindNeedTransferrResultDate(enterpriseID, taskId int64, paginate ormService.PaginateParam) (result []*transfer_customer.TransferCustomerDataModel, err error) {
	err = service.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, enterpriseID), ormService.ScopeCond_EQDelete,
		ormService.ScopeCond_Paginate(paginate)).Where(" task_id=? AND transfer_customer_status=?", taskId, 2).Order("id").Find(&result).Error
	return
}
func (service *TransferCustomerService) FindNeedTransferrResultDateCount(enterpriseID, taskId int64) (count int64, err error) {
	err = service.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, enterpriseID), ormService.ScopeCond_EQDelete).
		Where(" task_id=? AND transfer_customer_status=?", taskId, 2).Count(&count).Error
	return
}
func (service *TransferCustomerService) UpdatetransfercustomerStatus(enterpriseID int64, p *transfer_customer.TransferCustomerDataModel, transfercustomerlog transfer_customer.TransferCustomerLogModel) error {
	tx := service.Begin()
	err := tx.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, enterpriseID), ormService.ScopeCond_EQDelete).
		Where("id=?", p.ID).Update("transfer_customer_status", p.TransferCustomerStatus).Error
	if err != nil {
		return err
	}

	err = tx.Scopes(model.BranchTable(&transfercustomerlog, enterpriseID), ormService.ScopeCond_EQDelete).
		Create(&transfercustomerlog).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	return err
}

func (service *TransferCustomerService) Deletetransfercustomer(enterpriseID int64, p *transfer_customer.TransferCustomerDataModel) error {
	tx := service.Begin()
	err := tx.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, enterpriseID), ormService.ScopeCond_EQDelete).
		Where("id=?", p.ID).Update("deleted_at", time.Now().Unix()).Error
	if err != nil {
		return err
	}
	err = tx.Scopes(model.BranchTable(&transfer_customer.TransferCustomerLogModel{}, enterpriseID), ormService.ScopeCond_EQDelete).
		Where("handover_userid=? AND takeover_userid=? AND external_userid=? AND task_id=?",
			p.HandoverUserid, p.TakeoverUserid, p.ExternalUserid, p.ID).Update("deleted_at", time.Now().Unix()).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	return err
}
func (service *TransferCustomerService) UpdateTaskStatus(enterpriseID, taskId int64) error {
	updateMap := map[string]interface{}{
		"status":   3,
		"end_time": time.Now().Unix(),
	}
	return service.Scopes(model.BranchTable(&transfer_customer.TransferCustomerTasksModel{}, enterpriseID)).
		Where("id", taskId).Updates(updateMap).Error
}
