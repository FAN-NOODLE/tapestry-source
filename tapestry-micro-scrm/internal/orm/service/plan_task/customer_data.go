package plan_task

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/modules/plan_task/delivery/http_entity"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"strings"
)

/*
查询子任务的用户画像&社群任务的任务总数
*/
func (w *PlanTaskService) TaskCount(enterpriseID, taskId int64) (data int64, err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Where("task_id = ? ", taskId).
		Where("deleted_at  = ? ", micro_scrm_enterprise.ISNOTDELETE).
		Count(&data).Error
	return
}

/*
查询子任务的用户画像&社群任务的完成的任务的总数
*/
func (w *PlanTaskService) TaskComplateCount(enterpriseID, taskId, startTime, endTime int64) (data int64, err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Where("task_id = ? ", taskId).
		Where("end_time >= ? ", startTime).
		Where("end_time <= ? ", endTime).
		Where("is_finished  = ? ", micro_scrm_enterprise.FINISHED).
		Where("deleted_at  = ? ", micro_scrm_enterprise.ISNOTDELETE).
		Count(&data).Error
	return
}

/*
查询子任务用户画像的触达导购的人数
*/
func (w *PlanTaskService) TouchBaCount(enterpriseID, taskId int64) (data int64, err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Where("task_id = ? ", taskId).
		Where("deleted_at  = ? ", micro_scrm_enterprise.ISNOTDELETE).
		Distinct("user_id").
		Count(&data).Error
	return
}

/*
查询子任务的用户画像匹配客户数
*/
func (w *PlanTaskService) CustomerCount(enterpriseID, taskId, planId int64) (data int64, err error) {
	enterpriseIdStr := fmt.Sprintf("micro_scrm_enterprise_%d", enterpriseID)
	tableName := fmt.Sprintf("%s%d.%s_%d", app.Cfg().DB.DBPrefix, enterpriseID, micro_scrm_enterprise.TaskConditionDataModel{}.TableName(), planId)
	sql := "SELECT count(union_id) as count" +
		" FROM " + enterpriseIdStr + ".task_allot_log AS allot LEFT JOIN " +
		" " + tableName + " AS cond " +
		" ON cond.task_allot_log_id = allot.id" +
		" where allot.task_id = ?" +
		" and allot.deleted_at = ?" +
		" and cond.deleted_at = ?"
	err = w.Raw(sql, taskId, micro_scrm_enterprise.ISNOTDELETE,
		micro_scrm_enterprise.ISNOTDELETE).
		Scan(&data).Error
	return
}

/*
查询子任务的用户画像匹配客户数(去除重复的union_id)
*/
func (w *PlanTaskService) CustomerDistinctCount(enterpriseID, taskId, planId int64) (data int64, err error) {
	enterpriseIdStr := fmt.Sprintf("micro_scrm_enterprise_%d", enterpriseID)
	tableName := fmt.Sprintf("%s%d.%s_%d", app.Cfg().DB.DBPrefix, enterpriseID, micro_scrm_enterprise.TaskConditionDataModel{}.TableName(), planId)
	sql := "SELECT count(DISTINCT (CONCAT(union_id,customer_id))) as count" +
		" FROM " + enterpriseIdStr + ".task_allot_log AS allot LEFT JOIN " +
		" " + tableName + " AS cond " +
		" ON cond.task_allot_log_id = allot.id" +
		" where allot.task_id = ?" +
		" and allot.deleted_at = ?" +
		" and cond.deleted_at = ?"
	err = w.Raw(sql, taskId, micro_scrm_enterprise.ISNOTDELETE,
		micro_scrm_enterprise.ISNOTDELETE).
		Scan(&data).Error
	return
}

func (p *PlanTaskService) ExportPlanTasksCompletion(enterpriseID int64, param http_entity.ExportTasksRequst, searchUserIds []string) (result []*micro_scrm_enterprise.ExportTaskDataResp, err error) {
	enterpriseStr := fmt.Sprintf("micro_scrm_enterprise_%v", enterpriseID)

	query := p.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("task_allot_log.id,task_allot_log.user_id,task_allot_log.is_finished,"+
			" p.name AS plan_name,pt.name AS task_name,task_allot_log.created_at AS start_time,"+
			" task_allot_log.end_time AS completed_time,"+
			" task_allot_log.customer_num,p.category,task_allot_log.finish_rate,task_allot_log.contacted_count,task_allot_log.moment_look_count").
		Joins("LEFT JOIN "+enterpriseStr+".plan_tasks pt ON task_allot_log.task_id = pt.id  AND pt.deleted_at=0").
		Joins("LEFT JOIN "+enterpriseStr+".plans p ON task_allot_log.plan_id = p.id AND p.deleted_at=0").
		Where("task_allot_log.condition_id=?", param.Id).
		Where("task_allot_log.deleted_at=?", 0)
	if param.DepartmentId != 0 {
		query = query.Where("d.id=?", param.DepartmentId)
	}
	if param.PlanStartTime != 0 && param.PlanEndTime != 0 {
		query = query.Where("task_allot_log.created_at BETWEEN ? AND ?",
			param.PlanStartTime, param.PlanEndTime)
	}
	if param.CompletionStartTime != 0 && param.CompletionEndTime != 0 {
		query = query.Where("task_allot_log.end_time BETWEEN ? AND ?",
			param.CompletionStartTime, param.CompletionEndTime)
	}
	if len(searchUserIds) > 0 {
		query = query.Where("task_allot_log.user_id IN (?)", searchUserIds)
	}
	if param.IsFinished == 0 {
		query = query.Where("task_allot_log.is_finished=?", 0)
	} else if param.IsFinished == 1 {
		query = query.Where("task_allot_log.is_finished=?", 1)
	} else if param.IsFinished == 2 {
		query = query.Where("(task_allot_log.is_finished=? or task_allot_log.is_finished=?)", 0, 2)
	}
	err = query.Group("task_allot_log.id").
		Order("task_allot_log.id DESC").
		Find(&result).Error
	return
}

/*
查询子任务中社群任务的触达的群主人数
*/
func (w *PlanTaskService) TouchMasterCount(enterpriseID, taskId, startTime, endTime int64) (data int64, err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Where("task_id = ? ", taskId).
		Where("created_at  >= ? ", startTime).
		Where("created_at  <= ? ", endTime).
		Where("deleted_at  = ? ", micro_scrm_enterprise.ISNOTDELETE).
		Distinct("user_id").
		Count(&data).Error
	return
}

/*
查询子任务的社群任务匹配社群数
*/
func (w *PlanTaskService) GroupChatCount(enterpriseID, planId, taskId, startTime,
	endTime int64) (data int64, err error) {
	err = w.Scopes(func(db *gorm.DB) *gorm.DB {
		tableName := fmt.Sprintf("%s%d.%s_%d", app.Cfg().DB.DBPrefix, enterpriseID, micro_scrm_enterprise.TaskConditionDataModel{}.TableName(), planId)
		return db.Table(tableName)
	}).Where("task_id = ? ", taskId).
		Where("created_at  >= ? ", startTime).
		Where("created_at  <= ? ", endTime).
		Where("deleted_at  = ? ", micro_scrm_enterprise.ISNOTDELETE).
		Distinct("group_chat_id").
		Count(&data).Error
	return
}

/*
查询子任务中的社群任务中时间段内新创建的任务数
*/
func (w *PlanTaskService) GroupNewTaskCount(enterpriseID, taskId, startTime,
	endTime int64) (data int64, err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Where("task_id = ? ", taskId).
		Where("created_at >= ? ", startTime).
		Where("created_at <= ? ", endTime).
		Where("deleted_at = ? ", micro_scrm_enterprise.ISNOTDELETE).
		Count(&data).Error
	return
}

func (p *PlanTaskService) ExportGroupTasksListData(enterpriseID int64,
	param http_entity.ExportTasksRequst) []micro_scrm_enterprise.ExportGroupChatDataResp {
	var user []micro_scrm_enterprise.ExportGroupChatDataResp
	enterpriseStr := fmt.Sprintf("micro_scrm_enterprise_%v", enterpriseID)
	query := p.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("task_allot_log.id,gc.name group_chat_name,task_allot_log.user_id,b.name user_name,"+
			" task_allot_log.is_finished,task_allot_log.created_at start_time,"+
			" task_allot_log.end_time completed_time,pt.name task_name,"+
			" d.NAME department_name,p.name plan_name").
		Joins("INNER JOIN  members b "+
			" ON task_allot_log.user_id = b.userid  AND b.deleted_at=0 and b.enterprise_id = ?", enterpriseID).
		Joins("LEFT JOIN member_departments c "+
			" ON b.id = c.member_id AND c.deleted_at=0").
		Joins("LEFT JOIN departments d "+
			" ON c.department_id = d.id  AND d.deleted_at=0").
		Joins("LEFT JOIN "+enterpriseStr+".plan_tasks pt "+
			" ON task_allot_log.task_id = pt.id AND pt.deleted_at=0").
		Joins("LEFT JOIN "+enterpriseStr+".plans p "+
			" ON task_allot_log.plan_id = p.id AND p.deleted_at=0").
		Joins("LEFT JOIN work_group_chat gc "+
			" ON task_allot_log.group_chat_id = gc.id AND gc.deleted_at=0").
		Where("task_allot_log.condition_id=?", param.Id).
		Where("task_allot_log.deleted_at=?", 0)
	if param.DepartmentId != 0 {
		query = query.Where("d.id=?", param.DepartmentId)
	}
	if param.PlanStartTime != 0 && param.PlanEndTime != 0 {
		query = query.Where("task_allot_log.created_at BETWEEN ? AND ?",
			param.PlanStartTime, param.PlanEndTime)
	}
	if param.CompletionStartTime != 0 && param.CompletionEndTime != 0 {
		query = query.Where("task_allot_log.end_time BETWEEN ? AND ?",
			param.CompletionStartTime, param.CompletionEndTime)
	}
	if param.BaName != "" {
		query = query.Where("b.name LIKE ?",
			fmt.Sprintf("%s%s%s", "%", param.BaName, "%"))
	}
	if param.IsFinished == 0 {
		query = query.Where("task_allot_log.is_finished=?", 0)
	} else if param.IsFinished == 1 {
		query = query.Where("task_allot_log.is_finished=?", 1)
	} else if param.IsFinished == 2 {
		query = query.Where("(task_allot_log.is_finished=? or task_allot_log.is_finished=?)", 0, 2)
	}
	query.Group("task_allot_log.id").Order("task_allot_log.id DESC").Scan(&user)

	return user
}

func (p *PlanTaskService) CustomerTaskListData(enterpriseID, id, planStartTime, planEndTime, completionStartTime, completionEndTime int64, searchUserIds []string, isFinished, page, size int) (user []micro_scrm_enterprise.PlanTasksConditionsAllotLog, total int64) {
	tx := p.Scopes(model.BranchTable(micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("id,user_id,start_time,end_time,is_finished,group_chat_id AS chat_id,tag_id,created_at,customer_num,finish_rate,contacted_count,moment_look_count").
		Where("deleted_at = 0").
		Where("condition_id = ?", id)
	if planStartTime != 0 && planEndTime != 0 {
		tx = tx.Where("ask_allot_log.created_at BETWEEN ? AND ?", planStartTime, planEndTime)
	}
	if completionStartTime != 0 && completionEndTime != 0 {
		tx = tx.Where("end_time BETWEEN ? AND ?", completionStartTime, completionEndTime)
	}
	if len(searchUserIds) > 0 {
		tx = tx.Where("user_id IN (?)", searchUserIds)
	}
	if isFinished == 0 {
		tx = tx.Where("is_finished=0")
	} else if isFinished == microScrmModel.PLAN_FINISHED {
		tx = tx.Where("is_finished=1")
	} else if isFinished == microScrmModel.PLAN_NOFINISH {
		tx = tx.Where("is_finished IN (?)", []int{0, 2})
	}
	tx = tx.Count(&total)
	tx = tx.Order("id DESC")
	tx = tx.Limit(size).Offset((page - 1) * size).Find(&user)
	return
}

func (p *PlanTaskService) GroupTasksListData(enterpriseID int64,
	param http_entity.TaskListRequst) (user []micro_scrm_enterprise.ExportGroupChatDataResp, total int64) {
	// enterpriseStr := fmt.Sprintf("micro_scrm_enterprise_%v", enterpriseID)
	enterprise := fmt.Sprintf("%v", enterpriseID)
	sqlTotal := " SELECT count(1) FROM ( "
	sqlStr := "select task_allot_log.id,gc.name group_chat_name,task_allot_log.user_id,b.name user_name," +
		"task_allot_log.is_finished,task_allot_log.created_at start_time," +
		"task_allot_log.end_time completed_time, pt.name task_name,d.NAME department_name," +
		"p.name plan_name " +
		" from micro_scrm_enterprise_" + enterprise + ".task_allot_log" +
		" INNER JOIN  members b  ON task_allot_log.user_id = b.userid  AND b.deleted_at=0 " +
		" LEFT JOIN member_departments c  ON b.id = c.member_id AND c.deleted_at=0" +
		" LEFT JOIN departments d ON c.department_id = d.id  AND d.deleted_at=0 and d.enterprise_id= ?" +
		" LEFT JOIN micro_scrm_enterprise_" + enterprise + ".plan_tasks pt  ON task_allot_log.task_id = pt.id AND pt.deleted_at=0" +
		" LEFT JOIN micro_scrm_enterprise_" + enterprise + ".plans p  ON task_allot_log.plan_id = p.id AND p.deleted_at=0" +
		" LEFT JOIN work_group_chat gc ON task_allot_log.group_chat_id = gc.id AND gc.deleted_at=0" +
		" where task_allot_log.condition_id=?" +
		" and task_allot_log.deleted_at=0" +
		" "

	// 修改为查询部门下所有部门
	if len(param.DepartmentIds) > 0 {

		// 切片转字符串
		departmentIds := strings.Replace(strings.Trim(fmt.Sprint(param.DepartmentIds), "[]"), " ", ",", -1)
		sqlStr += fmt.Sprintf(" and d.id  in  (%v)", departmentIds)

	}
	// if param.DepartmentId != 0 {
	//	sqlStr += fmt.Sprintf(" and d.id=%v", param.DepartmentId)
	// }
	if param.PlanStartTime != 0 && param.PlanEndTime != 0 {
		sqlStr += fmt.Sprintf(" and task_allot_log.created_at BETWEEN %v AND %v", param.PlanStartTime, param.PlanEndTime)
	}
	if param.CompletionStartTime != 0 && param.CompletionEndTime != 0 {
		sqlStr += fmt.Sprintf(" and task_allot_log.end_time BETWEEN %v AND %v", param.CompletionStartTime, param.CompletionEndTime)
	}
	if param.BaName != "" {
		sqlStr += fmt.Sprintf(" and gc.name LIKE '%s%s%s'", "%", param.BaName, "%")
	}
	if param.IsFinished == 0 {
		sqlStr += fmt.Sprintf(" and task_allot_log.is_finished=%v", 0)
		// query = query.Where(, 0)
	} else if param.IsFinished == 1 {
		sqlStr += fmt.Sprintf(" and task_allot_log.is_finished=%v", 1)
	} else if param.IsFinished == 2 {
		sqlStr += fmt.Sprintf(" and (task_allot_log.is_finished=%v or task_allot_log.is_finished=%v)", 0, 2)
	}
	sqlTotal += sqlStr
	sqlTotal += " GROUP BY task_allot_log.id) AS T"
	sqlStr += " GROUP BY task_allot_log.id"
	sqlStr += " order by task_allot_log.id DESC  limit ? offset " + fmt.Sprintf("%v", (param.Page-1)*param.Size)
	p.Raw(sqlTotal, enterpriseID, param.Id).Scan(&total)
	p.Raw(sqlStr, enterpriseID, param.Id, param.Size).Scan(&user)
	return
}

func (p *PlanTaskService) ExportPlanTasksCompletionV2(enterpriseID, id, planStartTime, planEndTime, completionStartTime, completionEndTime int64, searchUserIds []string, isFinished int) (user []micro_scrm_enterprise.PlanTasksConditionsAllotLog, total int64) {
	tx := p.Scopes(model.BranchTable(micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("id,user_id,start_time,end_time,is_finished,group_chat_id AS chat_id,tag_id,created_at,customer_num,finish_rate,contacted_count,moment_look_count").
		Where("deleted_at = 0").
		Where("condition_id = ?", id)
	if planStartTime != 0 && planEndTime != 0 {
		tx = tx.Where("ask_allot_log.created_at BETWEEN ? AND ?", planStartTime, planEndTime)
	}
	if completionStartTime != 0 && completionEndTime != 0 {
		tx = tx.Where("end_time BETWEEN ? AND ?", completionStartTime, completionEndTime)
	}
	if len(searchUserIds) > 0 {
		tx = tx.Where("user_id IN (?)", searchUserIds)
	}
	if isFinished == 0 {
		tx = tx.Where("is_finished=0")
	} else if isFinished == microScrmModel.PLAN_FINISHED {
		tx = tx.Where("is_finished=1")
	} else if isFinished == microScrmModel.PLAN_NOFINISH {
		tx = tx.Where("is_finished IN (?)", []int{0, 2})
	}
	tx = tx.Count(&total)
	tx = tx.Order("id DESC")
	tx = tx.Find(&user)
	return
}
