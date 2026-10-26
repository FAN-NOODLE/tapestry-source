package plan_task

import (
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

const (
	FinishStatus = 1 // 子任务完成状态
)

// GetBaseInfo 任务的基本信息
func (p *PlanTaskService) GetBaseInfo(enterpriseID, planID int64) (result micro_scrm_enterprise.PlanInquireInfo, err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.PlanModel{}, enterpriseID)).
		Select("plans.name AS task_name,plans.start_time AS plan_start_time,plans.end_time AS plan_end_time").
		Where("plans.id = ? AND plans.deleted_at = 0 ", planID).
		Find(&result).Error
	return
}

// GetPlanBasicInfo 时间段内的基本信息
func (p *PlanTaskService) GetPlanBasicInfo(enterpriseID, planID, entTime int64) (result []micro_scrm_enterprise.PlanInquireInfo, err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).
		Select("id,dimension,status").
		Where("plan_id = ? AND deleted_at = 0", planID).
		Where("created_at <= ?", entTime).
		Find(&result).
		Error
	return
}

func (p *PlanTaskService) GetTaskCustomerInfo(enterpriseID, planID int64) (result []int64, err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).
		Select("id").
		Where("plan_id = ? ", planID).
		Where("dimension = ? ", 2).
		Where("deleted_at = 0").
		Find(&result).Error
	return
}

// GetTaskFinishCount 子任务完成数
func (p *PlanTaskService) GetTaskFinishCount(enterpriseID, endTime int64, taskIDList []int64) (result int64, err error) {
	err = p.Table("(?) AS T", p.Table("(?) AS M", p.Scopes(model.BranchTable(micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("IF(is_finished = 1, 1, 0) AS is_finished,task_id").
		Where("deleted_at = 0").
		Where("created_at <= ?", endTime)).
		Where("task_id IN (?)", taskIDList).
		Group("task_id").
		Select("task_id,MIN( is_finished ) AS is_finished")).
		Where("is_finished = 1").
		Count(&result).Error
	return
}

// GetTaskNumber 获取任务件数
func (p *PlanTaskService) GetTaskNumber(enterpriseID, endTime int64, taskIDList []int64) (result int64, err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("COUNT(1)").
		Where("created_at <= ?", endTime).
		Where("task_id IN (?) AND deleted_at = 0", taskIDList).
		Find(&result).Error
	return
}

// GetTaskFinishNumber 获取已完成任务件数
func (p *PlanTaskService) GetTaskFinishNumber(enterpriseID int64, taskIDList []int64) (result int64, err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("COUNT(1)").
		Where("task_id IN (?) AND deleted_at = 0", taskIDList).
		Where("is_finished = ?", 1).
		Find(&result).Error
	return
}

func (p *PlanTaskService) GetTaskGroupInfo(enterpriseID, planID, endTime int64, taskName string, pageSize, pageNum int,
	taskType int8) (result []int64, count int64, err error) {
	tx := p.Scopes(model.BranchTable(micro_scrm_enterprise.PlanTaskModel{}, enterpriseID))
	tx = tx.Select("id")
	if taskName != "" {
		tx = tx.Where("name LIKE ?", "%"+taskName+"%")
	}
	tx.Where("plan_id = ? AND deleted_at = 0 ", planID).
		Where("dimension = ?", taskType).
		Where("created_at <= ?", endTime).Count(&count)
	if pageNum != 0 && pageSize != 0 {
		tx = tx.Limit(pageSize).Offset(pageSize * (pageNum - 1))
	}
	tx = tx.Find(&result)
	err = tx.Error
	return
}

// CumulativeMatchGroup 获取累计匹配的社群
func (p *PlanTaskService) CumulativeMatchGroup(enterpriseID, endTime int64, groupIDList []int64) (
	result []micro_scrm_enterprise.MatchGroupIDList, err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("DISTINCT(group_chat_id) AS group_chat_id").
		Where("task_id IN (?)", groupIDList).
		Where("deleted_at = 0 ").
		Where("created_at <= ?", endTime).
		Find(&result).Error
	return
}

// CumulativeMatchLord 获取累计触达群主
func (p *PlanTaskService) CumulativeMatchLord(enterpriseID, endTime int64, groupList []int64) (result int64, err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("COUNT(DISTINCT(user_id))").
		Where("task_id IN (?)", groupList).
		Where("deleted_at = 0").
		Where("created_at <= ?", endTime).
		Find(&result).Error
	return
}

// GetReachGuideNumber 获取触达导购数
func (p *PlanTaskService) GetReachGuideNumber(enterpriseID, endTime int64, portraitIDList []int64) (result int64, err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("COUNT(DISTINCT(user_id))").
		Where("task_id IN (?)", portraitIDList).
		Where("deleted_at = 0").
		Where("created_at <= ?", endTime).
		Find(&result).Error
	return
}

// GetTaskName 获取子任务名称
func (p *PlanTaskService) GetTaskName(enterpriseID int64, taskIDList []int64) (result []micro_scrm_enterprise.TaskName, err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).
		Select("id,name").
		Where("id IN (?)", taskIDList).
		Find(&result).Error
	return
}

// GetTaskName2 获取子任务名称
func (p *PlanTaskService) GetTaskName2(enterpriseID int64, taskIDList []int64) (result []*micro_scrm_enterprise.CommunityDataDetailResult, err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).
		Select("id,name AS subtasks_number").
		Where("id IN (?)", taskIDList).
		Find(&result).Error
	return
}

// GetTaskName3 获取子任务名称
func (p *PlanTaskService) GetTaskName3(enterpriseID int64, taskIDList []int64) (result []*micro_scrm_enterprise.CustomerPortraitDetailResult, err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).
		Select("id,name AS subtasks_number,dimension").
		Where("id IN (?)", taskIDList).
		Find(&result).Error
	return
}

// GetGroupNumber 获取群数,群主数,任务总数
func (p *PlanTaskService) GetGroupNumber(enterpriseID, endTime int64, taskIDList []int64) (result []micro_scrm_enterprise.TaskGroupNumberCount, err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("task_id,COUNT(DISTINCT(group_chat_id)) AS group_number,"+
			"COUNT(DISTINCT(user_id)) AS lord_number,COUNT(task_id) AS task_number").
		Where("task_id IN (?)", taskIDList).
		Where("created_at <= ?", endTime).
		Where("deleted_at = 0").
		Group("task_id").
		Find(&result).Error

	return
}

// // GetGuideNumber 获取导购数
// func (p *PlanTaskService) GetGuideNumber(enterpriseID int64, taskIDList []int64) (result []micro_scrm_enterprise.TaskNumberCount, err error) {
//	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
//		Select("task_id,COUNT(DISTINCT(user_id)) AS number").
//		Where("task_id IN (?)", taskIDList).
//		Where("deleted_at = 0").
//		Group("task_id").
//		Find(&result).Error
//	return
// }

// // GetTaskCountNumber 获取任务总数
// func (p *PlanTaskService) GetTaskCountNumber(enterpriseID int64, taskIDList []int64) (result []micro_scrm_enterprise.TaskNumberCount, err error) {
//	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
//		Select("task_id,COUNT(DISTINCT(task_id)) AS number ").
//		Where("task_id IN (?)", taskIDList).
//		Where("deleted_at = 0").
//		Group("task_id").
//		Find(&result).Error
//	return
// }

// GetTaskCountFinishNumber 获取完成任务数
func (p *PlanTaskService) GetTaskCountFinishNumber(enterpriseID, endTime int64, taskIDList []int64) (result []micro_scrm_enterprise.TaskNumberCount, err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("task_id,COUNT(task_id) AS number").
		Where("task_id IN (?)", taskIDList).
		Where("deleted_at = 0").
		Where("created_at <= ?", endTime).
		Where("is_finished = ?", FinishStatus).
		Group("task_id").
		Find(&result).Error
	return
}

// // GetDistinctCustomNumber 获取客户数(去重)
// func (p *PlanTaskService) GetDistinctCustomNumber(enterpriseID int64, taskIDList []int64) (result []micro_scrm_enterprise.TaskNumberCount, err error) {
//	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.TaskConditionDataModel{}, enterpriseID)).
//		Select("task_id,COUNT(DISTINCT(union_id)) AS number").
//		Where("task_id IN (?)", taskIDList).
//		Where("deleted_at = 0").
//		Group("task_id").
//		Find(&result).Error
//	return
// }
