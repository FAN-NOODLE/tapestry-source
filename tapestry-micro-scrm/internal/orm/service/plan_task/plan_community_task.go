package plan_task

import (
	"fmt"
	"time"

	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm/group_chat"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/plan_task"
)

const (
	PALNNOSTART = iota
	PALNSTART
)

func (p *PlanTaskService) GetPlanCommunityTask(enterpriseID, planTime int64, conditionType int) (list []micro_scrm_enterprise.ResponCondition, err error) {
	planTasksTable := "micro_scrm_enterprise_" + fmt.Sprintf("%d", enterpriseID) + ".plan_tasks"
	planTaskConditionsTable := "micro_scrm_enterprise_" + fmt.Sprintf("%d", enterpriseID) + ".plan_task_conditions"
	err = p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanModel{}, enterpriseID)).
		Select("plans.id plan_id,plans.name plan_name,a.id task_id,a.name task_name,a.is_specify_crowd,a.task_is_push,b.id condition_id,plans.start_time,plans.end_time,a.condition_type,b.conditions,plans.status plan_status,a.status task_status").
		Joins("INNER JOIN "+planTasksTable+" a ON a.plan_id = plans.id AND a.dimension = 1").
		Joins("INNER JOIN "+planTaskConditionsTable+" b ON b.task_id = a.id").
		Where("plans.deleted_at = ? AND a.deleted_at = ? AND b.deleted_at = ? AND ? BETWEEN a.start_time AND a.end_time  AND a.condition_type = ? AND plans.status IN ?", 0, 0, 0, planTime, conditionType, []int{0, 1}).Group("a.id").
		Scan(&list).Error

	return
}

/**
*获取群人数
 */
func (p *PlanTaskService) WorkGroupChatTotalNumber(enterpriseID, chatID int64, specifyCrowdUser []string) (result group_chat.GroupChatPoeple, err error) {
	query := p.Table("work_group_chat").Select("work_group_chat.id,work_group_chat.owner_member_id,c.userid user_id,member_num number").
		Joins("INNER JOIN members c ON work_group_chat.owner_member_id=c.id  AND c.enterprise_id=? AND c.deleted_at=0", enterpriseID).
		Where("work_group_chat.id=? AND work_group_chat.enterprise_id=? AND work_group_chat.status=? AND work_group_chat.dissolve_time=?", chatID, enterpriseID, 1, 0)
	if len(specifyCrowdUser) > 0 {
		query.Where("c.userid IN (?)", specifyCrowdUser)
	}
	err = query.Find(&result).Error
	return
}

/**
*获取群人数
 */
func (p *PlanTaskService) WorkGroupChatTotalNumberByenterpriseID(enterpriseID int64) (result []group_chat.GroupChatPoeple, err error) {
	query := p.Table("work_group_chat").Select("work_group_chat.id,work_group_chat.owner_member_id,c.userid user_id,member_num number").
		Joins("INNER JOIN members c ON work_group_chat.owner_member_id=c.id  AND c.enterprise_id=? AND c.deleted_at=0", enterpriseID).
		Where("work_group_chat.enterprise_id=? AND work_group_chat.status=? AND work_group_chat.dissolve_time=?", enterpriseID, 1, 0)
	err = query.Find(&result).Error
	return
}

/**
*获取群新增人数
 */
func (p *PlanTaskService) WorkGroupChatAddNumber(enterpriseID, chatID int64, startTime, endTime int64, specifyCrowdUser []string) (result group_chat.GroupChatPoeple, err error) {

	query := p.Table("work_group_chat_member").Select("b.id,b.owner_member_id,c.userid user_id,SUM(1) number ").
		Joins("INNER JOIN work_group_chat b ON work_group_chat_member.group_chat_id = b.id  AND work_group_chat_member.deleted_at=0").
		Joins("INNER JOIN members c ON b.owner_member_id=c.id  AND c.enterprise_id=? AND c.deleted_at=0", enterpriseID).
		Where("b.id=? AND b.enterprise_id = ? AND b.deleted_at = ? AND b.status =? AND b.dissolve_time =? AND work_group_chat_member.status=?  AND work_group_chat_member.join_time BETWEEN ? AND ?  AND work_group_chat_member.join_time>work_group_chat_member.leave_time", chatID, enterpriseID, 0, 1, 0, 1, startTime, endTime)
	if len(specifyCrowdUser) > 0 {
		query.Where("c.userid IN (?)", specifyCrowdUser)
	}
	err = query.Group("b.id,b.owner_member_id,c.userid").First(&result).Error
	return
}

/**
*获取退群人数
 */
func (p *PlanTaskService) WorkGroupChatRefundNumber(enterpriseID, chatID int64, startTime, endTime int64, specifyCrowdUser []string) (result group_chat.GroupChatPoeple, err error) {

	query := p.Table("work_group_chat_member").Select("b.id,b.owner_member_id,c.userid user_id,SUM(1) number ").
		Joins("INNER JOIN work_group_chat b ON work_group_chat_member.group_chat_id = b.id  AND work_group_chat_member.deleted_at=0").
		Joins("INNER JOIN members c ON b.owner_member_id=c.id AND c.enterprise_id=? AND c.deleted_at=0", enterpriseID).
		Where("b.id=? AND b.enterprise_id = ? AND b.deleted_at = ? AND b.status =? AND b.dissolve_time =? AND work_group_chat_member.status=?  AND work_group_chat_member.leave_time BETWEEN ? AND ?  AND work_group_chat_member.leave_time>work_group_chat_member.join_time", chatID, enterpriseID, 0, 1, 0, 2, startTime, endTime)
	if len(specifyCrowdUser) > 0 {
		query.Where("c.userid IN (?)", specifyCrowdUser)
	}
	err = query.Group("b.id,b.owner_member_id,c.userid").Find(&result).Error

	return
}

/**
*批量添加task_allot_log
 */
func (p *PlanTaskService) BatchCreatetaskAllotLog(enterpriseID int64, taskAllotLogs []micro_scrm_enterprise.TaskAllotLogModel) error {
	return p.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).Create(&taskAllotLogs).Error
}

func (p *PlanTaskService) CreatetaskAllotLog(enterpriseID int64, taskAllotLog micro_scrm_enterprise.TaskAllotLogModel) (int64, error) {

	err := p.Scopes(model.BranchTable(&taskAllotLog, enterpriseID)).Create(&taskAllotLog).Error
	return taskAllotLog.ID, err
}

func (p *PlanTaskService) FindworExternalUser(enterpriseID int64, userID string) (result map[string]string, err error) {
	err = p.Table("work_external_user a").Joins("INNER JOIN user b ON a.user_id=b.id").
		Select("a.union_id,a.external_userid,a.name,b.mobile").
		Where("a.enterprise_id=? AND a.external_userid=? AND a.deleted_at=? AND b.deleted_at=?", enterpriseID, userID, 0, 0).
		Scan(&result).Error
	return
}
func (p *PlanTaskService) UpdateTaskPush(enterpriseID, taskId int64) (err error) {
	err = p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).Where("id=?", taskId).
		Update("task_is_push", 1).Error
	return
}
func (p *PlanTaskService) UpdateplanTaskStatusByID(enterpriseID, planId int64) (err error) {
	err = p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).Where("plan_id=? AND  status=? AND deleted_at=?", planId, PALNNOSTART, PALNNOSTART).
		Update("status", PALNSTART).Error
	return
}
func (p *PlanTaskService) UpdateplanStatusByID(enterpriseID, planId int64) (err error) {
	err = p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanModel{}, enterpriseID)).Where("id=?", planId).
		Update("status", 1).Error
	return
}
func (p *PlanTaskService) UpdateplanTaskStatusByIDs(enterpriseID int64, planIds []int64) (err error) {
	err = p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).Where("plan_id IN (?) AND  status=? AND deleted_at=?", planIds, PALNNOSTART, PALNNOSTART).
		Update("status", PALNSTART).Error
	return
}
func (p *PlanTaskService) UpdateplanStatusByIDs(enterpriseID int64, planIds []int64) (err error) {
	err = p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanModel{}, enterpriseID)).Where("id IN (?) AND  status=?", planIds, PALNNOSTART).
		Update("status", 1).Error
	return
}
func (p *PlanTaskService) FindGroupmenenberByuserId(enterpriseID, groupChatID int64, userId string) (memberlog group_chat.WorkGroupChatMemberLogInfo, err error) {
	err = p.Table("work_group_chat_member").Select("work_group_chat_member.*,a.name chat_name,c.name user_name").
		Joins("INNER JOIN work_group_chat a ON work_group_chat_member.group_chat_id=a.id AND a.deleted_at=0").
		Joins("LEFT JOIN members c ON work_group_chat_member.userid=c.userid AND c.enterprise_id=?", enterpriseID).
		Where("work_group_chat_member.group_chat_id=? AND work_group_chat_member.userid=? AND work_group_chat_member.deleted_at=?", groupChatID, userId, 0).First(&memberlog).Error
	return
}
func (p *PlanTaskService) CreatememberLos(memberLos []group_chat.WorkGroupChatMemberLog) (err error) {
	err = p.Create(&memberLos).Error
	return
}

func (p *PlanTaskService) FindShequnTaskPushIS(enterpriseID int64, taskAllotLog micro_scrm_enterprise.TaskAllotLogModel) (data micro_scrm_enterprise.TaskAllotLogModel, err error) {
	// startTime, endTime := helper.PeriodOfTimestamp()
	err = p.Scopes(model.BranchTable(&taskAllotLog, enterpriseID)).
		Where("task_id =? AND condition_id=? AND user_id=?  AND deleted_at=? ", taskAllotLog.TaskId, taskAllotLog.ConditionId, taskAllotLog.UserId, 0).
		First(&data).Error
	return
}

func (p *PlanTaskService) UpdatetaskAllotLogTime(enterpriseID, id int64) (err error) {
	err = p.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Where("id =? ", id).Update("updated_at", time.Now().Unix()).Error
	return
}

func (p *PlanTaskService) UPdatetaskDownloadedFile(enterpriseID, id int64) error {
	err := p.Scopes(model.BranchTable(&plan_task.PlanTaskDownloadedFileModel{}, enterpriseID)).
		Where("id=? ", id).Update("is_extract", 3).Error
	return err
}

func (p *PlanTaskService) FindShequnTaskPushByConditionId(enterpriseID, conditionId int64) (data micro_scrm_enterprise.TaskAllotLogModel, err error) {
	err = p.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).
		Where("condition_id=?  AND deleted_at=? ", conditionId, 0).
		First(&data).Error
	return
}
