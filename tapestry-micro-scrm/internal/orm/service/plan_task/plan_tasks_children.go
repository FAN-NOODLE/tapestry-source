package plan_task

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	"micro-scrm/internal/pkg/helper"
	"time"

	"micro-scrm/internal/orm/model/micro_scrm/group_chat"

	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/service"
)

func (p *PlanTaskService) PlanTasksAdd(enterpriseID, planId int64, taskNum int, planTask micro_scrm_enterprise.PlanTaskModel, planTaskConditionModels []micro_scrm_enterprise.PlanTaskConditionModel, planTaskSpecifyCrowds []micro_scrm_enterprise.PlanTaskSpecifyCrowdModel) (err error) {
	tx := p.Begin()
	if err = tx.Scopes(model.BranchTable(planTask, enterpriseID)).Create(&planTask).Error; err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, enterpriseID)).Create(&planTaskConditionModels).Error; err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanModel{}, enterpriseID)).Where("id=?", planId).Update("task_num", taskNum).Error; err != nil {
		tx.Rollback()
		return
	}
	if len(planTaskSpecifyCrowds) > 0 {
		if err = tx.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskSpecifyCrowdModel{}, enterpriseID)).Create(&planTaskSpecifyCrowds).Error; err != nil {
			tx.Rollback()
			return
		}
	}
	err = tx.Commit().Error
	return

}
func (p *PlanTaskService) PlanTasksUpdate(enterpriseID int64, planTask micro_scrm_enterprise.PlanTaskModel) error {
	return p.Scopes(model.BranchTable(planTask, enterpriseID)).Updates(&planTask).Error
}
func (p *PlanTaskService) PlanTasksConditionsBatchAdd(enterpriseID int64, planTaskConditionModels []micro_scrm_enterprise.PlanTaskConditionModel) error {
	return p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, enterpriseID)).Create(&planTaskConditionModels).Error
}
func (p *PlanTaskService) PlanTasksConditionsUpdateById(enterpriseID int64, planTaskConditionModels micro_scrm_enterprise.PlanTaskConditionModel) error {
	return p.Scopes(model.BranchTable(planTaskConditionModels, enterpriseID)).Updates(&planTaskConditionModels).Error

}
func (p *PlanTaskService) PlanTasksConditionsBatchUpdate(enterpriseID int64, deletPlanTaskConditions []int64, addplanTaskConditionModels, updateplanTaskConditionModels []micro_scrm_enterprise.PlanTaskConditionModel) (err error) {
	tx := p.Begin()
	if len(addplanTaskConditionModels) > 0 {
		if err = tx.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, enterpriseID)).Create(&addplanTaskConditionModels).Error; err != nil {
			tx.Rollback()
			return
		}
	}
	if len(deletPlanTaskConditions) > 0 {
		if err = tx.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, enterpriseID)).Where("id IN (?)", deletPlanTaskConditions).Update("deleted_at", time.Now().Unix()).Error; err != nil {
			tx.Rollback()
			return
		}
	}
	if len(updateplanTaskConditionModels) > 0 {
		for _, value := range updateplanTaskConditionModels {
			if err = tx.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, enterpriseID)).Updates(&value).Error; err != nil {
				tx.Rollback()
				return
			}
		}
	}
	err = tx.Commit().Error
	return
}
func (p *PlanTaskService) FindplanTaskInfoById(enterpriseID, taskId int64) *micro_scrm_enterprise.PlanTaskModel {
	result := micro_scrm_enterprise.PlanTaskModel{}
	p.Scopes(model.BranchTable(&result, enterpriseID)).Scopes(service.ScopeCond_EQDelete).Where("id", taskId).Find(&result)
	return &result
}

func (p *PlanTaskService) FindPlanTasksConditionsById(enterpriseID, Id int64) *micro_scrm_enterprise.PlanTaskConditionModel {
	result := micro_scrm_enterprise.PlanTaskConditionModel{}
	p.Scopes(model.BranchTable(&result, enterpriseID)).Scopes(service.ScopeCond_EQDelete).Where("id", Id).Find(&result)
	return &result
}
func (p *PlanTaskService) PlanTasksConditionsById(enterpriseID, taskId int64) []micro_scrm_enterprise.PlanTaskConditionModel {
	planTaskConditions := []micro_scrm_enterprise.PlanTaskConditionModel{}
	p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, enterpriseID)).Scopes(service.ScopeCond_EQDelete).Where("task_id", taskId).Order("id").Find(&planTaskConditions)
	return planTaskConditions
}

func (p *PlanTaskService) PlanTasksDelet(enterpriseID, taskID, DeletedAt, planId int64, taskNum int) (err error) {
	tx := p.Begin()
	if err = tx.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).Where("id=?", taskID).Update("deleted_at", DeletedAt).Error; err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, enterpriseID)).Where("task_id=?", taskID).Update("deleted_at", DeletedAt).Error; err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanModel{}, enterpriseID)).Where("id=?", planId).Update("task_num", taskNum).Error; err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	return
}

func (p *PlanTaskService) PlanTasksConditionsDeleteByTaskId(enterpriseID, taskId int64) error {

	return p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, enterpriseID)).Where("task_id=? AND deleted_at=?", taskId, 0).Update("deleted_at", time.Now().Unix()).Error
}
func (p *PlanTaskService) PlanTasksConditionsDeleteById(enterpriseID, id int64) error {

	return p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, enterpriseID)).Where("id=? AND deleted_at=?", id, 0).Update("deleted_at", time.Now().Unix()).Error
}

func (p *PlanTaskService) FindTaskConditionsById(enterpriseID, id int64) (data micro_scrm_enterprise.PlanTaskConditionModel, err error) {
	err = p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, enterpriseID)).Where("id=?", id, 0).Find(&data).Error
	return
}
func (p *PlanTaskService) PlanTasksConditionsCompletion(enterpriseID, id, departmentId, planStartTime, planEndTime, completionStartTime, completionEndTime int64, baName string, isFinished, page, size int) ([]micro_scrm_enterprise.PlanTasksConditionsAllotLog, int64) {
	var (
		user  []micro_scrm_enterprise.PlanTasksConditionsAllotLog
		total = int64(0)
	)
	query := p.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).Select("task_allot_log.user_id,b.name,d.id department_id,d.name department_name,"+
		"task_allot_log.start_time,task_allot_log.end_time,task_allot_log.is_finished,task_allot_log.group_chat_id AS chat_id,task_allot_log.tag_id").
		Joins("INNER JOIN  members b ON task_allot_log.user_id = b.userid  AND b.deleted_at=0").
		Joins("LEFT JOIN member_departments c ON b.id = c.member_id AND c.deleted_at=0").
		Joins("LEFT JOIN departments d ON c.department_id = d.id  AND d.deleted_at=0").
		Where("task_allot_log.condition_id=?", id).
		Where("task_allot_log.deleted_at=?", 0).
		Where("b.enterprise_id=?", enterpriseID)
	if departmentId != 0 {
		query = query.Where("d.id=?", departmentId)
	}
	if planStartTime != 0 && planEndTime != 0 {
		query = query.Where("task_allot_log.start_time BETWEEN ? AND ?", planStartTime, planEndTime)
	}
	if completionStartTime != 0 && completionEndTime != 0 {
		query = query.Where("task_allot_log.end_time BETWEEN ? AND ?", completionStartTime, completionEndTime)
	}
	if baName != "" {
		query = query.Where("b.name LIKE ?", fmt.Sprintf("%s%s%s", "%", baName, "%"))
	}
	if isFinished == 0 {
		query = query.Where("task_allot_log.is_finished=?", 0)
	} else if isFinished == 1 {
		query = query.Where("task_allot_log.is_finished=?", 1)
	} else if isFinished == 2 {
		query = query.Where("task_allot_log.is_finished=?", 2)
	}
	query.Count(&total).Order("task_allot_log.created_at DESC").Limit(size).Offset((page - 1) * size).Scan(&user)

	return user, total
}

func (p *PlanTaskService) PlanTasksConditionsCompletioncustomerNumv1(enterpriseID, Id, planId int64, userIds []string) ([]micro_scrm_enterprise.TaskAllotcustomerNum, error) {
	result := []micro_scrm_enterprise.TaskAllotcustomerNum{}
	params := []interface{}{Id}
	userWHere := ""
	for k, userId := range userIds {
		if k == 0 {
			userWHere = "?"
		} else {
			userWHere += ",?"
		}
		params = append(params, userId)
	}
	task_condition_data_table := fmt.Sprintf("%s%d.%s_%d", app.Cfg().DB.DBPrefix, enterpriseID, micro_scrm_enterprise.TaskConditionDataModel{}.TableName(), planId)
	sql := " SELECT user_id,union_id FROM " + task_condition_data_table + " WHERE condition_id =? AND user_id IN (" + userWHere + ") AND union_id!='' GROUP BY user_id,union_id "
	tablesql := "SELECT t.user_id,sum(1) customer_num FROM (" + sql + ") t  GROUP BY t.user_id"
	err := p.Raw(tablesql, params...).Scan(&result).Error
	return result, err
}

// GetCorrespondCustomerNumber 获取符合条件的客户
func (p *PlanTaskService) GetCorrespondCustomerNumber(enterpriseID, logID, planId int64, pageSize, pageNum int, customerName string) (result []micro_scrm_enterprise.TaskConditionDataCustomer, count int64, err error) {
	var wxUserCropExternalTable = micro_scrm_enterprise.WxUserCorpExternalModel{}.AllTableName(enterpriseID)
	tableName := fmt.Sprintf("%s%d.%s_%d as a", app.Cfg().DB.DBPrefix, enterpriseID, micro_scrm_enterprise.TaskConditionDataModel{}.TableName(), planId)
	tx := p.Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Table(tableName)
	}).Select("a.tag_name,wxUser.name,wxUser.sex,wxUser.avatar,a.ww_tag_id").
		Joins("LEFT JOIN "+wxUserCropExternalTable+" AS wxUser ON wxUser.unionid = a.union_id AND wxUser.deleted_at = 0 and a.union_id!=''").
		Where("a.task_allot_log_id = ?", logID).
		Where("a.deleted_at = 0")
	if customerName != "" {
		tx = tx.Where("wxUser.name LIKE ?", "%"+customerName+"%")
	}
	tx = tx.Count(&count)
	tx = tx.Limit(pageSize).Offset(pageSize * (pageNum - 1)).Find(&result)
	err = tx.Error
	return
}

func (p *PlanTaskService) GetGroupChatName(enterpriseID int64) (result []group_chat.WorkGroupChatModel) {
	p.Model(&group_chat.WorkGroupChatModel{}).Where("deleted_at = ? AND enterprise_id = ?", 0, enterpriseID).
		Group("name").Find(&result)
	return
}

func (p *PlanTaskService) GetLeaveCount() (result []micro_scrm_enterprise.GetLeaveInfo) {
	p.Model(&group_chat.WorkGroupChatMemberModel{}).Select("group_chat_id,leave_time").
		Where("status = ?", 2).
		Find(&result)
	return
}

func (p *PlanTaskService) GetLeaveTime(enterpriseID int64, userID string) (result int64) {
	p.Model(&group_chat.WorkGroupChatMemberModel{}).Select("leave_time").
		Where("enterprise_id = ? AND userid = ?", enterpriseID, userID).Find(&result)
	return
}
func (p *PlanTaskService) PlanTasksIsSpecifyCrowd(enterpriseID, taskId int64, isSpecifyCrowd int, planTaskSpecifyCrowds []micro_scrm_enterprise.PlanTaskSpecifyCrowdModel) error {
	tx := p.Begin()
	if err := tx.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskSpecifyCrowdModel{}, enterpriseID)).Where("task_id=?", taskId).Update("deleted_at", time.Now().Unix()).Error; err != nil {
		tx.Rollback()
		return err
	}
	if isSpecifyCrowd == 2 && len(planTaskSpecifyCrowds) > 0 {
		if err := tx.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskSpecifyCrowdModel{}, enterpriseID)).Create(&planTaskSpecifyCrowds).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	err := tx.Commit().Error
	return err
}
func (p *PlanTaskService) GetListByTaskId(enterpriseID, taskId int64) (data []micro_scrm_enterprise.PlanTaskSpecifyCrowdData, err error) {
	err = p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskSpecifyCrowdModel{}, enterpriseID)).Select("plan_id,task_id,user_id,user_name,is_ba,third_id,type").
		Where("task_id=? AND deleted_at=?", taskId, 0).Find(&data).Error

	return
}
func (p *PlanTaskService) PlanTasksIsSpecifyCrowdByTaskId(c ctx.Context, enterpriseID, taskId int64) (userIds []string, err error) {
	var list []micro_scrm_enterprise.PlanTaskSpecifyCrowdData
	err = p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskSpecifyCrowdModel{}, enterpriseID)).Select("plan_id,task_id,user_id,user_name,is_ba,third_id,type").
		Where("task_id=? AND deleted_at=?", taskId, 0).Find(&list).Error
	if err != nil {
		return
	}
	var thirdIds = make([]int, 0)
	for _, v := range list {
		if v.Type == micro_scrm_enterprise.TaskCrowdTypeSa && v.UserId != "" {
			userIds = append(userIds, v.UserId)
		} else if v.Type == micro_scrm_enterprise.TaskCrowdTypeDept && v.ThirdId > 0 {
			thirdIds = append(thirdIds, v.ThirdId)
		}
	}
	if len(thirdIds) > 0 {
		memberList, _ := departmentsDao.NewDepartments().GetMemberDataByThirdIds(c, enterpriseID, thirdIds)
		if len(memberList) > 0 {
			for _, v := range memberList {
				userIds = append(userIds, v.Userid)
			}
		}
	}
	userIds = helper.ArrayUniq(userIds)
	return
}

func (p *PlanTaskService) PlanTasksIsSpecifyCrowdUpdateDelete(enterpriseID, taskId int64) error {
	return p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskSpecifyCrowdModel{}, enterpriseID)).Where("task_id=?", taskId).Update("deleted_at", time.Now().Unix()).Error
}

func (p *PlanTaskService) PlanTasksIsSpecifyCrowdUpdateCreate(enterpriseID int64, planTaskSpecifyCrowds []micro_scrm_enterprise.PlanTaskSpecifyCrowdModel) error {

	return p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskSpecifyCrowdModel{}, enterpriseID)).Create(&planTaskSpecifyCrowds).Error
}
