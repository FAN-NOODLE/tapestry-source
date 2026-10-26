package plan_task

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	planReq "micro-scrm/internal/app/http/request/plan"
	"strconv"
	"time"

	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm/group_chat"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/service"
	"micro-scrm/internal/pkg/db/gormx"
)

const (
	TASK_ALL = iota
	TASK_NOT_START
	TASK_GOING
	TASK_END
)

// GetPlanList 获取计划列表
func (g *PlanTaskService) GetPlanList(ctx ctx.Context, enterpriseID int64, searchType int, operationUserId []int64,
	createName string, page int,
	size int,
	param micro_scrm_enterprise.PlanModel,
	taskLabel int32,
) (
	data []micro_scrm_enterprise.PlanModel, total int64, err error) {

	if page < 0 {
		page = 1
	}

	if size < 0 {
		size = 10
	}

	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanModel{}, enterpriseID)).WithContext(ctx).
		Scopes(service.ScopeCond_EQDelete).Scopes(func(db *gorm.DB) *gorm.DB {
		if searchType == 0 && param.ID != 0 {
			db.Where("id = ?", param.ID)
		} else {
			if param.Name != "" {
				db.Where("name like ?", fmt.Sprintf("%s%s%s", "%", param.Name, "%"))
			}

		}
		if param.PlanTemplateID > 0 {
			db.Where("plan_template_id = ?", param.PlanTemplateID)
		}

		if taskLabel > 0 {
			db.Where("category = ?", taskLabel)
		}

		if param.Status >= 0 {
			db.Where("status = ?", param.Status)
		}

		if param.StartTime != 0 {
			db.Where("created_at >= ?", param.StartTime)
		}

		if param.EndTime != 0 {
			db.Where("created_at <= ?", param.EndTime)
		}

		if createName != "" {
			if len(operationUserId) != 0 {
				db.Where("operation_user_id in ?", operationUserId)
			} else {
				db.Where("operation_user_id in ?", []int{0})
			}
		}
		return db
	}).Count(&total).Order("updated_at desc").Limit(size).
		Offset((page - 1) * size).
		Find(&data).Error

	return
}

// GetTaskList 获取任务列表
func (g *PlanTaskService) GetTaskList(enterpriseID int64, page, size int, param planReq.TaskListRequest) (
	data []micro_scrm_enterprise.PlanTaskModel, total int64, err error) {

	if page < 0 {
		page = 1
	}

	if size < 0 {
		size = 10
	}

	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).
		Scopes(service.ScopeCond_EQDelete).Scopes(func(db *gorm.DB) *gorm.DB {
		if param.TaskId > 0 {
			db.Where("id=?", param.TaskId)
		}
		if param.TaskName != "" {
			db.Where("plan_name = ?", param.TaskName)
		}
		if param.TemplateId > 0 {
			db.Where("plan_template_id = ?", param.TemplateId)
		}
		if param.CreatedStartTime > 0 && param.CreatedEndTime > 0 {
			db.Where("created_at BETWEEN ? AND ?", param.CreatedStartTime, param.CreatedEndTime)
		}
		if param.Category >= 0 {
			db.Where("task_category=?", param.Category)
		}
		if param.GroupType >= 0 {
			db.Where("group_send_type=?", param.GroupType)
		}
		if param.TaskStarTime > 0 && param.TaskEndTime > 0 {
			db.Where("task_start_time BETWEEN ? AND ?", param.TaskStarTime, param.TaskEndTime)
		}
		if param.OriginType >= 0 {
			db.Where("origin_type=?", param.OriginType)
		}
		return db
	}).Where("plan_id=?", param.PlanId).
		Order("created_at desc").
		Count(&total).Limit(size).
		Offset((page - 1) * size).
		Find(&data).Error

	return
}

// GetPlanListByFront 小程序首页任务计划列表
func (g *PlanTaskService) GetPlanListByFront(eid, startTime, endTime int64, page, size int, userId, name string, status int) (data []ResponseFrontPlanList, err error) {

	db := g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, eid)).
		Select("plans.id,plans.`name`,plans.`remark`,plans.`status`,plans.`is_open`,plans.`start_time`,plans.`end_time`,task_allot_log.created_at ").
		Joins("left join "+getDBName(eid)+".plans on task_allot_log.plan_id=plans.id").
		Where("plans.deleted_at=0 and task_allot_log.deleted_at=0 and plans.is_open = 1 and task_allot_log.user_id=? AND task_allot_log.created_at>=? and task_allot_log.created_at<=?", userId, startTime, endTime).
		Order("FIELD(`status`, 1,0,2),end_time ASC")
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	if status > 0 {
		switch status {
		case TASK_NOT_START:
			db = db.Where("plans.status = ?", 0)
		case TASK_GOING:
			db = db.Where("plans.status = ?", 1)
		case TASK_END:
			db = db.Where("plans.status = ?", 2)
		}
	}

	err = db.Limit(size).Offset((page - 1) * size).Find(&data).Error
	return
}

// GetPlanListByFront 小程序首页任务计划列表
func (g *PlanTaskService) GetPlanListByFrontBak(eid, startTime, endTime int64, page, size int, name string) (data []ResponseFrontPlanList, err error) {
	db := g.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanModel{}, eid)).
		Scopes(service.ScopeCond_EQDelete).
		Where("is_open = 1 AND ((start_time < ?  AND end_time > ?) or (start_time < ?  AND end_time>?) or (start_time > ?  AND end_time <= ?))", startTime, startTime, endTime, endTime, startTime, endTime).
		Order("FIELD(`status`, 1,0,2),end_time ASC")
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	err = db.Limit(size).Offset((page - 1) * size).Find(&data).Error
	return
}

// GetPlanListByFront 小程序首页任务计划列表
func (g *PlanTaskService) GetAllPlanIDS(eid, startTime, endTime int64) (ids []int64, err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanModel{}, eid)).
		Scopes(service.ScopeCond_EQDelete).
		Select("id").
		Where("is_open = 1 and status!=0 AND ((start_time < ?  AND end_time > ?) or (start_time < ?  AND end_time>?) or (start_time > ?  AND end_time <= ?))", startTime, startTime, endTime, endTime, startTime, endTime).
		Scan(&ids).Error
	return
}

// GetPlanListByFront 小程序首页任务计划列表
func (g *PlanTaskService) GetDayPlanListByFront(eid int64, page, size int, ids []int64, name string) (data []ResponseFrontPlanList, err error) {
	db := g.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanModel{}, eid)).
		Scopes(service.ScopeCond_EQDelete).
		Where("id in ? and is_open = 1", ids).
		Order("FIELD(`status`, 1,0,2) ASC,end_time desc")
	// Where("id in ? and is_open = 1 AND ((start_time < ?  AND end_time > ?) or (start_time < ?  AND end_time>?) or (start_time > ?  AND end_time <= ?))", ids, startTime, startTime, endTime, endTime, startTime, endTime)
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	err = db.Limit(size).Offset((page - 1) * size).Find(&data).Error
	return
}

// GetPlanListByFront 小程序首页任务计划列表
func (g *PlanTaskService) GetPlanListByFrontBAK(eid, startTime, endTime int64, page, size int, userId, name string) (data []ResponseFrontPlanList, err error) {
	db := g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, eid)).
		Select("plans.id,task_allot_log.is_finished,task_allot_log.created_at,plans.`name`,plans.remark,plans.type,plans.`STATUS`,plans.is_open,plans.start_time,plans.end_time").
		Joins("left join "+getDBName(eid)+".plans on plans.id = task_allot_log.plan_id").
		Where("task_allot_log.user_id=? and task_allot_log.created_at BETWEEN ? and ?", userId, startTime, endTime)
	// Where("plan.is_open = 1 AND ((start_time < ?  AND end_time > ?) or (start_time < ?  AND end_time>?) or (start_time > ?  AND end_time <= ?))", ids, startTime, startTime, endTime, endTime, startTime, endTime)
	if name != "" {
		db = db.Where("plans.`name` like ?", "%"+name+"%")
	}
	err = db.Limit(size).Offset((page - 1) * size).Find(&data).Error
	return
}

func (g *PlanTaskService) GetPlanInfoByFront(eid, pid int64) (data ResponseFrontPlanList, err error) {
	db := g.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanModel{}, eid)).
		Scopes(service.ScopeCond_EQDelete).
		Where("id=?", pid)
	err = db.Find(&data).Error
	return
}

func (g *PlanTaskService) GetTaskInfoByFront(eid, pid int64) (data micro_scrm_enterprise.TaskAllotLogModel, err error) {
	db := g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, eid)).
		Scopes(service.ScopeCond_EQDelete).
		Where("id=?", pid)
	err = db.Find(&data).Error
	return
}

// GetPlanTaskListByFront 小程序首页任务计划子任务列表
func (g *PlanTaskService) GetPlanTaskListByFrontBAK(eid, pid int64, page, size int) (data []ResponseFrontPlanTaskList, err error) {
	dbName := getDBName(eid)
	db := g.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, eid)).
		Select("plan_task_conditions.id, plan_tasks.`name`,plan_tasks.remark,plan_tasks.dimension,task_allot_log.is_finished as `status`,task_allot_log.customer_num").
		Joins("left join "+dbName+".plan_tasks on plan_tasks.id = plan_task_conditions.task_id").
		Joins("left join "+dbName+".task_allot_log on task_allot_log.condition_id = plan_task_conditions.id").
		Where(dbName+".plan_tasks.plan_id=? and "+dbName+".plan_tasks.deleted_at=0", pid)
	err = db.Limit(size).Offset((page - 1) * size).Find(&data).Error
	return
}

func (g *PlanTaskService) GetPlanTaskListByFront(eid, pid, startTime, endTime int64, page, size int, userId string, status int) (data []ResponseFrontPlanTaskList, err error) {
	dbName := getDBName(eid)
	// count(id) 修改为 count(DISTINCT work_external_user_id) 外部联系人有多个标签算一个人的统计调整
	db := g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, eid)).
		Select("task_allot_log.id,task_allot_log.customer_num,plan_tasks.`name`,plan_tasks.remark,plan_tasks.dimension,plan_tasks.condition_type,"+
			"task_allot_log.is_finished as `status`,task_allot_log.condition_id,task_allot_log.group_chat_id,"+
			"plans.start_time,plans.end_time,plan_tasks.task_type,plan_tasks.task_category,plan_tasks.group_send_type,task_allot_log.close_tips,task_allot_log.chat_num,plans.id as plan_id,plans.name as plan_name,plan_tasks.id as task_id,task_allot_log.finish_rate,plan_tasks.task_start_time,plan_tasks.task_end_time,plan_tasks.task_expire_type").
		Joins("left join "+dbName+".plan_tasks on plan_tasks.id = task_allot_log.task_id").
		Joins("left join "+dbName+".plans on plans.id = task_allot_log.plan_id").
		Where("task_allot_log.plan_id=? and task_allot_log.user_id=? and plan_tasks.deleted_at=0 and "+
			"task_allot_log.created_at BETWEEN ? and ?", pid, userId, startTime, endTime).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if status >= 0 {
				db.Where("task_allot_log.is_finished=?", status)
			}
			return db
		}).
		Group("task_allot_log.id").Order("status asc")
	err = db.Limit(size).Offset((page - 1) * size).Find(&data).Error
	return
}

func (g *PlanTaskService) GetLatestTaskByUserId(eid, planId int64, userId string) (data micro_scrm_enterprise.TaskAllotLogModel, err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, eid)).
		Scopes(service.ScopeCond_EQDelete).
		Where("user_id=? and plan_id=?", userId, planId).Order("created_at desc").First(&data).Error
	return
}

func (g *PlanTaskService) GetFinishedStatus(eid, logID int64, userId string, groupIDList []int64) (result []*micro_scrm_enterprise.AddGroupInfo) {
	g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, eid)).Select("is_finished,group_chat_id").
		Where("condition_id = ? AND user_id = ? AND group_chat_id IN (?) "+
			"and deleted_at = 0 ", logID, userId, groupIDList).
		Find(&result)
	return
}

// GetPlanTaskInfoFront 小程序任务计划子任务详情
func (g *PlanTaskService) GetPlanTaskInfoFront(eid, pid int64) (data []ResponseFrontPlanTaskList, err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskModel{}, eid)).
		Scopes(service.ScopeCond_EQDelete).
		Where("id=?", pid).Find(&data).Error
	return
}

func (g *PlanTaskService) GetConditionInfo(enterpriseID int64, idList []int64) (result []micro_scrm_enterprise.ConditionInfo, err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, enterpriseID)).
		Select("id,conditions,content").
		Scopes(service.ScopeCond_EQDelete).
		Where("id in (?)", idList).Find(&result).Error
	return
}

// GetTaskConInfoFrontBAK  小程序任务计划子任务条件详情
func (g *PlanTaskService) GetTaskConInfoFrontBAK(eid, id int64) (data ResponseFrontPlanTaskConList, err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, eid)).
		Scopes(service.ScopeCond_EQDelete).
		Where("id=?", id).Find(&data).Error
	return
}
func (g *PlanTaskService) GetTaskID(enterpriseID, logID int64) (result int64) {
	g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).Select("task_id").
		Where("deleted_at = 0 AND id = ?", logID).Find(&result)
	return

}

func (g *PlanTaskService) GetTaskConInfoFront(eid, id int64) (data []*ResponseFrontPlanTaskConList, err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, eid)).
		Select("task_allot_log.id,task_allot_log.is_finished,plan_task_conditions.content,plan_task_conditions.conditions,plan_task_conditions.content_type,"+
			"plan_tasks.dimension,plan_tasks.condition_type,task_allot_log.group_chat_id,plans.start_time,plans.end_time,plan_tasks.remark,task_allot_log.condition_id").
		Joins("left join "+getDBName(eid)+".plan_task_conditions on plan_task_conditions.id = task_allot_log.condition_id").
		Joins("left join "+micro_scrm_enterprise.PlanTaskModel{}.AllTableName(eid)+" on plan_tasks.id = plan_task_conditions.task_id").
		Where("task_allot_log.condition_id IN (SELECT condition_id from "+micro_scrm_enterprise.TaskAllotLogModel{}.AllTableName(eid)+
			" where id = ? AND deleted_at = 0)", id).Find(&data).Error
	return
}
func (g *PlanTaskService) GetChatName(enterpriseID int64) (result []micro_scrm_enterprise.GroupChatInfo) {
	g.Model(&group_chat.WorkGroupChatModel{}).Select("id,name,chat_id").Where("deleted_at = 0 AND "+
		"enterprise_id = ?", enterpriseID).Find(&result)
	return
}

func (g *PlanTaskService) GetChatId(enterpriseID, logID int64) (result micro_scrm_enterprise.Receives) {
	g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("task_allot_log.id,task_allot_log.plan_id,task_allot_log.created_at,task_allot_log.group_chat_id,task_allot_log.user_id,task_allot_log.condition_id,"+
			"task_allot_log.is_finished,plans.start_time,task_allot_log.task_id,plans.end_time,"+
			"plan_task_conditions.conditions,plan_task_conditions.content,plan_tasks.remark,plan_tasks.dimension,"+
			"plan_tasks.condition_type,plan_task_conditions.content_type,task_allot_log.finish_rate,plan_tasks.task_start_time,plan_tasks.task_end_time,plan_tasks.task_expire_type").
		Joins("left join "+micro_scrm_enterprise.PlanModel{}.AllTableName(enterpriseID)+" on plans.id = task_allot_log.plan_id").
		Joins("left join "+micro_scrm_enterprise.PlanTaskModel{}.AllTableName(enterpriseID)+" on plan_tasks.id = task_allot_log.task_id").
		Joins("left join"+micro_scrm_enterprise.PlanTaskConditionModel{}.AllTableName(enterpriseID)+" on "+
			"plan_task_conditions.id = task_allot_log.condition_id").
		Where("task_allot_log.id = ?", logID).Find(&result)
	return
}
func (g *PlanTaskService) GetConditionIdInfo(enterpriseID, logID int64) (result micro_scrm_enterprise.Receives2) {
	g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("task_allot_log.id,task_allot_log.user_id,plans.start_time,plans.end_time,task_allot_log.condition_id,task_allot_log.is_finished,plans.status,task_allot_log.is_forward,task_allot_log.task_id,plans.name as plan_name,plans.id as plan_id").
		Joins("left join "+micro_scrm_enterprise.PlanModel{}.AllTableName(enterpriseID)+" on plans.id = task_allot_log.plan_id").
		Where("task_allot_log.id = ?", logID).Find(&result)
	return
}

func (g *PlanTaskService) GetGroupInfo(enterpriseID, conditionID int64, userID string) (result []*micro_scrm_enterprise.AddGroupInfo) {
	g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("task_allot_log.id,task_allot_log.group_chat_id,task_allot_log.is_finished").
		Where("task_allot_log.condition_id = ? AND user_id = ? AND deleted_at = 0", conditionID, userID).
		Group("task_allot_log.group_chat_id").
		Find(&result)
	return
}

func (g *PlanTaskService) GetGroupInfoByTaskAllotLogId(enterpriseID, id int64, userID string) (result []*micro_scrm_enterprise.AddGroupInfo) {
	g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Select("task_allot_log.id,task_allot_log.group_chat_id,task_allot_log.is_finished").
		Where("task_allot_log.id = ? AND user_id = ? AND deleted_at = 0", id, userID).
		Find(&result)
	return
}

func (g *PlanTaskService) GetConditionDataInfo(enterpriseID, startTime, endTime int64, groupChatID int64) (result []micro_scrm_enterprise.DetailInfo) {
	g.Model(group_chat.WorkGroupChatMemberModel{}).
		Select("wxce.name,work_group_chat_member.userid AS work_external_user_id,wxce.sex AS gender,wxce.avatar").
		Joins("left join "+micro_scrm_enterprise.WxUserCorpExternalModel{}.AllTableName(enterpriseID)+
			"wxce ON wxce.unionid = work_group_chat_member.unionid AND wxce.deleted_at = 0").
		Where("group_chat_id = ? AND leave_time BETWEEN ? and ? AND "+
			"leave_time > join_time AND work_group_chat_member.deleted_at =0", groupChatID, startTime, endTime).
		Group("work_group_chat_member.userid").
		Find(&result)
	return
}

func (g *PlanTaskService) GetConditionDataInfo2(enterpriseID, startTime, endTime int64, groupChatID []int64) (result []micro_scrm_enterprise.DetailInfo2) {
	g.Model(group_chat.WorkGroupChatMemberModel{}).
		Select("wxce.name,work_group_chat_member.userid AS work_external_user_id,wxce.sex AS gender,wxce.avatar,group_chat_id").
		Joins("left join "+micro_scrm_enterprise.WxUserCorpExternalModel{}.AllTableName(enterpriseID)+
			"wxce ON wxce.unionid = work_group_chat_member.unionid AND wxce.deleted_at = 0 ").
		Where("group_chat_id IN (?) AND leave_time BETWEEN ? and ? AND "+
			"leave_time > join_time AND work_group_chat_member.deleted_at =0", groupChatID, startTime, endTime).
		Group("work_group_chat_member.group_chat_id,work_group_chat_member.userid").
		Find(&result)
	return
}

func (g *PlanTaskService) GetStatus(enterpriseID, startTime, endTime int64, conditionID []int64) (result []micro_scrm_enterprise.Receives3) {
	g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).Select("id,is_finished,condition_id").
		Where("start_time >= ? AND end_time<= ?  AND condition_id IN (?)", startTime, endTime, conditionID).
		Find(&result)
	return
}

func (g *PlanTaskService) GetAddDataInfo(enterpriseID, startTime, endTime int64, groupChatID []int64) (result []micro_scrm_enterprise.DetailInfo) {
	g.Model(group_chat.WorkGroupChatMemberModel{}).
		Select("wxce.name,work_group_chat_member.userid AS work_external_user_id,wxce.sex AS gender "+
			"work_group_chat.name,work_group_chat.id").
		Joins("left join "+micro_scrm_enterprise.WxUserCorpExternalModel{}.AllTableName(enterpriseID)+
			"wxce ON wxce.unionid = work_group_chat_member.unionid").
		Joins("left join work_group_chat ON work_group_chat.id = work_group_chat_member.group_chat_id ").
		Where("group_chat_id IN (?) AND join_time BETWEEN ? and ? AND "+
			"join_time > join_time AND wxce.deleted_at = 0 "+
			"AND work_group_chat_member.deleted_at =0", groupChatID, startTime, endTime).
		Group("work_group_chat_member.unionid").Find(&result)
	return
}

func (g *PlanTaskService) GetNumCount(enterpriseID int64, idList int64) (number int64) {
	g.Model(&group_chat.WorkGroupChatModel{}).Select("member_num").Where("enterprise_id = ? "+
		"AND id = ? AND status = ? AND deleted_at = 0", enterpriseID, idList, 1).Find(&number)
	return
}

func (g *PlanTaskService) GetLeaveNumber(startTime, endTime, enterpriseID, chatID int64) (result []micro_scrm_enterprise.CountInfo) {
	g.Model(&group_chat.WorkGroupChatMemberModel{}).Select("unionid").Where("enterprise_id = ? AND leave_time <= ? AND leave_time >= ?"+
		" AND group_chat_id = ? AND deleted_at = 0", enterpriseID, endTime, startTime, chatID).Group("unionid").Find(&result)
	return
}

func (g *PlanTaskService) GetTaskConPlanInfoFront(eid, id int64) (data ResponseFrontPlanTaskConList, err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, eid)).
		Select("task_allot_log.id,task_allot_log.is_finished,plan_task_conditions.content,plan_task_conditions.conditions,plan_task_conditions.content_type").
		Joins("left join "+getDBName(eid)+".plan_task_conditions on plan_task_conditions.id = task_allot_log.condition_id").
		Where("task_allot_log.id=?", id).Find(&data).Error
	return
}

// GetTaskConInfoFront 小程序任务计划子任务多条件
func (g *PlanTaskService) GetTaskConListFront(eid, taskId int64) (data []ResponseFrontPlanTaskConList, err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, eid)).
		Scopes(service.ScopeCond_EQDelete).
		Where("task_id=?", taskId).Find(&data).Error
	return
}

func (g *PlanTaskService) UpdateGroupStatus(enterpriseID, isFinish, id int64) (err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).Where(" id = ? ", id).
		Updates(map[string]interface{}{"is_finished": 1, "end_time": time.Now().Unix()}).Error
	return
}

// SetStatusByFrontBAK 完成状态设置
func (g *PlanTaskService) SetStatusByFrontBAK(eid, conId int64, userId string) (data []ResponseFrontTagUserList, err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, eid)).
		Scopes(service.ScopeCond_EQDelete).
		Where("condition_id=? and user_id=?", conId, userId).Update("is_finished", 1).Error
	return
}

func (g *PlanTaskService) SetStatusByFront(eid, dimension, id, conditionID, groupChatID int64, userID string) (err error) {
	if dimension == 1 {
		var result []micro_scrm_enterprise.TaskConditionDataModel

		one, err := g.FindShequnTaskPushByConditionId(eid, conditionID)
		if err != nil {
			return err
		}
		tableName := fmt.Sprintf("%s%d.%s_%d", app.Cfg().DB.DBPrefix, eid, micro_scrm_enterprise.TaskConditionDataModel{}.TableName(), one.PlanId)
		err = g.Scopes(func(db *gorm.DB) *gorm.DB {
			db.Table(tableName)
			return db
		}).Scopes(service.ScopeCond_EQDelete).Where("task_allot_log_id = ? and is_finished=0", id).Scan(&result).Error
		if len(result) == 0 {
			return err
		}
		err = g.Scopes(func(db *gorm.DB) *gorm.DB {
			db.Table(tableName)
			return db
		}).Scopes(service.ScopeCond_EQDelete).Where(map[string]interface{}{
			"task_allot_log_id": id,
			"group_chat_id":     groupChatID,
		}).Updates(map[string]interface{}{"is_finished": 1, "finished_at": time.Now().Unix()}).Error
		if (len(result) == 1 && result[0].GroupChatId == groupChatID) || groupChatID == 0 { // 如果是最后一个社群任务完成，task_allot_log应变为完成,或者点击任务完成按钮变完成
			err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, eid)).
				Scopes(service.ScopeCond_EQDelete).Where("id = ?", id).Updates(map[string]interface{}{"is_finished": 1, "end_time": time.Now().Unix()}).Error
		}
	} else {
		tx := g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, eid)).
			Scopes(service.ScopeCond_EQDelete).Where("id = ?", id).Updates(map[string]interface{}{"is_finished": 1, "end_time": time.Now().Unix()})
		err = tx.Error
	}
	return
}

// getPlanIdsByBA 获取BA指定日的所有计划
func (g *PlanTaskService) GetPlanIdsByBA(eid int64, userId string, starTime, endTime int64) (data []int64, err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, eid)).
		Scopes(service.ScopeCond_EQDeleteAndGroup).
		Select("plan_id").
		Where("user_id=? and created_at BETWEEN ? and ?", userId, starTime, endTime).
		Scan(&data).Error
	return
}

// getPlanIdsByBA 获取BA所有计划
func (g *PlanTaskService) GetPlanIdsByBABAK(eid int64, userId string) (data []int64, err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, eid)).
		Scopes(service.ScopeCond_EQDelete).
		Select("plan_id").
		Where("user_id=?", userId).Scan(&data).Error
	return
}

// GetTaskIds 计划区间内包含的所有已完成的子任务完成时间(日)
func (g *PlanTaskService) GetFinishDaysBAK(eid, start_time, end_time int64, planid []int64, userId string) []int64 {
	taskIds := make([]int64, 0)
	db := gormx.GetInstance()
	sql := fmt.Sprintf("SELECT id FROM %s.task_allot_log WHERE user_id ='%s' AND is_finished = 1 AND plan_id in ? AND created_at BETWEEN %d AND %d GROUP BY id HAVING SUM(is_finished)=COUNT(id)", getDBName(eid), userId, start_time, end_time)
	db.Raw(sql, planid).Scan(&taskIds)
	dateTime := make([]int64, 0)
	db.Raw("select created_at from "+getDBName(eid)+".task_allot_log where id in ? group by created_at", taskIds).Scan(&dateTime)
	return dateTime
}

func (g *PlanTaskService) GetFinishDays(eid, start_time, end_time int64, userId string) []int64 {
	var finishDays []ResponseFinishDays
	db := gormx.GetInstance()
	// db.Raw("SELECT user_id,FROM_UNIXTIME( created_at, '%Y-%m-%d' ),COUNT( id ) AS log_count,SUM(IF( is_finished = 1, 1, 0 )) AS is_finish_num FROM micro_scrm_enterprise_"+strconv.FormatInt(eid, 10)+".task_allot_log WHERE  user_id = ? AND created_at BETWEEN ? AND ? AND deleted_at = 0 GROUP BY user_id,FROM_UNIXTIME( created_at, '%Y-%m-%d' )", userId, start_time, end_time).Scan(&finishDays)
	db.Raw("SELECT t.user_id,FROM_UNIXTIME( t.created_at, '%Y-%m-%d' )as created_at,COUNT( t.id ) AS log_count,SUM( IF ( t.is_finished = 1, 1, 0 ) ) AS is_finish_num FROM micro_scrm_enterprise_"+strconv.FormatInt(eid, 10)+".task_allot_log as t LEft JOIN micro_scrm_enterprise_"+strconv.FormatInt(eid, 10)+".plans as p on t.plan_id=p.id WHERE p.is_open=1 and user_id = ? AND t.created_at BETWEEN ? AND ? AND  t.deleted_at = 0 GROUP BY t.user_id,FROM_UNIXTIME(t.created_at,'%Y-%m-%d' )", userId, start_time, end_time).Scan(&finishDays)
	dateTime := make([]int64, 0)
	for _, v := range finishDays {
		if v.IsFinishNum == v.LogCount {
			tm, _ := time.Parse("2006-01-02", v.CreatedAt)
			dateTime = append(dateTime, tm.Unix())
		}
	}
	return dateTime
}

func (g *PlanTaskService) GetPlanTimeFinishDays(eid, pid, start_time, end_time int64, userId string) []int64 {
	var finishDays []ResponseFinishDays
	db := gormx.GetInstance()
	// db.Raw("SELECT user_id,FROM_UNIXTIME( created_at, '%Y-%m-%d' ),COUNT( id ) AS log_count,SUM(IF( is_finished = 1, 1, 0 )) AS is_finish_num FROM micro_scrm_enterprise_"+strconv.FormatInt(eid, 10)+".task_allot_log WHERE  user_id = ? AND created_at BETWEEN ? AND ? AND deleted_at = 0 GROUP BY user_id,FROM_UNIXTIME( created_at, '%Y-%m-%d' )", userId, start_time, end_time).Scan(&finishDays)
	db.Raw("SELECT t.user_id,FROM_UNIXTIME( t.created_at, '%Y-%m-%d' )as created_at,COUNT( t.id ) AS log_count,SUM( IF ( t.is_finished = 1, 1, 0 ) ) AS is_finish_num FROM micro_scrm_enterprise_"+strconv.FormatInt(eid, 10)+".task_allot_log as t LEft JOIN micro_scrm_enterprise_"+strconv.FormatInt(eid, 10)+".plans as p on t.plan_id=p.id WHERE p.is_open=1 and t.plan_id=? and user_id = ? AND t.created_at BETWEEN ? AND ? AND  t.deleted_at = 0 GROUP BY t.user_id,FROM_UNIXTIME(t.created_at,'%Y-%m-%d' )", pid, userId, start_time, end_time).Scan(&finishDays)
	dateTime := make([]int64, 0)
	for _, v := range finishDays {
		if v.IsFinishNum == v.LogCount {
			tm, _ := time.Parse("2006-01-02", v.CreatedAt)
			dateTime = append(dateTime, tm.Unix())
		}
	}
	return dateTime
}

func (g *PlanTaskService) GetTaskFinishDays(eid, start_time, end_time int64, userId string, pid int64) []int64 {
	var finishDays []ResponseFinishDays
	db := gormx.GetInstance()
	// db.Raw("SELECT user_id,FROM_UNIXTIME( created_at, '%Y-%m-%d' ),COUNT( id ) AS log_count,SUM(IF( is_finished = 1, 1, 0 )) AS is_finish_num FROM micro_scrm_enterprise_"+strconv.FormatInt(eid, 10)+".task_allot_log WHERE  user_id = ? AND created_at BETWEEN ? AND ? AND deleted_at = 0 GROUP BY user_id,FROM_UNIXTIME( created_at, '%Y-%m-%d' )", userId, start_time, end_time).Scan(&finishDays)
	db.Raw("SELECT t.user_id,FROM_UNIXTIME( t.created_at, '%Y-%m-%d' )as created_at,COUNT( t.id ) AS log_count,SUM( IF ( t.is_finished = 1, 1, 0 ) ) AS is_finish_num FROM micro_scrm_enterprise_"+strconv.FormatInt(eid, 10)+".task_allot_log as t LEft JOIN micro_scrm_enterprise_"+strconv.FormatInt(eid, 10)+".plans as p on t.plan_id=p.id WHERE t.pland_id=? and p.is_open=1 and t.user_id = ? AND t.created_at BETWEEN ? AND ? AND  t.deleted_at = 0 GROUP BY t.user_id,FROM_UNIXTIME(t.created_at,'%Y-%m-%d' )", pid, userId, start_time, end_time).Scan(&finishDays)
	dateTime := make([]int64, 0)
	for _, v := range finishDays {
		if v.IsFinishNum == v.LogCount {
			tm, _ := time.Parse("2006-01-02", v.CreatedAt)
			dateTime = append(dateTime, tm.Unix())
		}
	}
	return dateTime
}

// IsFinishByPlanTaskBAK 某个计划内的所有子任务是否全部完成
func (g *PlanTaskService) IsFinishByPlanTaskBAK(eid, start_time, end_time int64, planid int64, userId string) int {
	count := 0
	db := gormx.GetInstance()
	sql := fmt.Sprintf("SELECT count(id) FROM %s.task_allot_log WHERE user_id ='%s' AND plan_id=? and is_finished=1 AND created_at BETWEEN %d AND %d GROUP BY id HAVING SUM(is_finished)=COUNT(id)", getDBName(eid), userId, start_time, end_time)
	db.Raw(sql, planid).Scan(&count)
	return count
}

// IsFinishByPlanTask 某个计划内的所有子任务是否全部完成
func (g *PlanTaskService) IsFinishByPlanTask(eid, pid, start_time, end_time int64, userId string) int {
	var finishDays []ResponseFinishDays
	db := gormx.GetInstance()
	db.Raw("SELECT t.user_id,FROM_UNIXTIME( t.created_at, '%Y-%m-%d' )as created_at,COUNT( t.id ) AS log_count,SUM( IF ( t.is_finished = 1, 1, 0 ) ) AS is_finish_num FROM micro_scrm_enterprise_"+strconv.FormatInt(eid, 10)+".task_allot_log as t LEft JOIN micro_scrm_enterprise_"+strconv.FormatInt(eid, 10)+".plans as p on t.plan_id=p.id WHERE p.is_open=1 and t.plan_id=? and user_id = ? AND t.created_at BETWEEN ? AND ? AND  t.deleted_at = 0 GROUP BY t.user_id,FROM_UNIXTIME(t.created_at,'%Y-%m-%d' )", pid, userId, start_time, end_time).Scan(&finishDays)
	dateTime := make([]int64, 0)
	for _, v := range finishDays {
		if v.IsFinishNum == v.LogCount {
			tm, _ := time.Parse("2006-01-02", v.CreatedAt)
			dateTime = append(dateTime, tm.Unix())
		}
	}
	return len(dateTime)
}

func getDBName(eid int64) string {
	return fmt.Sprintf("micro_scrm_enterprise_%d", eid)
}

func (g *PlanTaskService) GetPlanTaskAllDays(eid, pid, start_time, end_time int64, userId string) []int64 {
	var allDays []int64
	var newAllDays []int64
	db := gormx.GetInstance()
	db.Raw("select created_at from micro_scrm_enterprise_"+strconv.FormatInt(eid, 10)+".task_allot_log where deleted_at=0 and user_id=? and plan_id=? and created_at BETWEEN ? and ?", userId, pid, start_time, end_time).Scan(&allDays)
	allCreateAt := make(map[int64]bool)
	for _, v := range allDays {
		if allCreateAt[v] {
			continue
		}
		allCreateAt[v] = true
		newAllDays = append(newAllDays, v)
	}
	return newAllDays
}
