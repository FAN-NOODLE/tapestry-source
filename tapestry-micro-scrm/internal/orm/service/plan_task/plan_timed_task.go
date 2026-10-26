package plan_task

import (
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"time"

	"micro-scrm/internal/orm/service"

	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/plan_task"

	"micro-scrm/internal/orm/model/micro_scrm"

	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

// 获取未开始计划列表
func (p *PlanTaskService) GetPlanStartTime(info micro_scrm_enterprise.PlanModel, enterpriseID int64) (result []micro_scrm_enterprise.PlanModel, err error) {
	err = p.Scopes(model.BranchTable(&info, enterpriseID)).Where("deleted_at = ? AND status = ? AND task_num >? AND start_time<=?", 0, 0, 0, time.Now().Unix()).Order("start_time ASC").Find(&result).Error
	return result, err

}

func (p *PlanTaskService) GetPerformPlan(info micro_scrm_enterprise.PlanModel, enterpriseID int64) (result []micro_scrm_enterprise.PlanModel, err error) {
	err = p.Scopes(model.BranchTable(&info, enterpriseID)).Where("deleted_at = ? AND status = ? AND end_time<=? ", 0, 1, time.Now().Unix()).Order("start_time ASC").Find(&result).Error
	return result, err

}

// 更新计划状态
func (p *PlanTaskService) UpdateStatus(enterpriseID, status int64, idList []int64, info micro_scrm_enterprise.PlanModel) error {
	if err := p.DB.Scopes(model.BranchTable(&info, enterpriseID)).Where("id in (?)", idList).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (p *PlanTaskService) UpdateTaskStatus(enterpriseID, status int64, planIdList []int64, info micro_scrm_enterprise.PlanTaskModel) error {
	if err := p.DB.Scopes(model.BranchTable(&info, enterpriseID)).Where("plan_id in (?)", planIdList).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

// GetPushData 获取推送标签类型的数据条件
func (p *PlanTaskService) GetPushData(enterpriseID int64, info micro_scrm_enterprise.PlanModel) (result []micro_scrm_enterprise.ResponCondition, err error) {
	err = p.Scopes(model.BranchTable(&info, enterpriseID)).
		Select("plans.is_push,plans.id AS plan_id,plans.name AS plan_name,pt.`name` AS task_name ,ptc.task_id,ptc.id AS condition_id,ptc.conditions,plans.start_time,plans.end_time,"+
			"ptc.content,ptc.content_type,pt.task_type,pt.send_type,pt.task_is_push,pt.is_specify_crowd,pt.plan_send_time,pt.origin_task_id,pt.task_start_time,pt.task_end_time,pt.task_expire_type").
		Joins("LEFT JOIN "+micro_scrm_enterprise.PlanTaskModel{}.AllTableName(enterpriseID)+" pt ON pt.plan_id = plans.id").
		Joins("LEFT JOIN "+micro_scrm_enterprise.PlanTaskConditionModel{}.AllTableName(enterpriseID)+" ptc ON ptc.task_id = pt.id").
		Where("plans.status = ? AND pt.status = ? AND pt.dimension=? AND plans.deleted_at = ? AND pt.deleted_at = ? AND ptc.deleted_at = ? AND plans.category=1 AND pt.task_is_push=0 ", 1, 1, microScrmModel.TaskDimensionTag, 0, 0, 0).Find(&result).Error
	return
}

// GetBirthdayPushData 获取生日类型的子任务
func (p *PlanTaskService) GetBirthdayPushData(enterpriseID int64, info micro_scrm_enterprise.PlanModel) (result []micro_scrm_enterprise.ResponseBirthday, err error) {
	err = p.Scopes(model.BranchTable(&info, enterpriseID)).
		Select("plans.id AS plan_id,plans.name AS plan_name,pt.`name` AS task_name ,ptc.task_id,ptc.id AS condition_id,ptc.conditions,plans.start_time,plans.end_time,"+
			"ptc.content,ptc.content_type,pt.is_specify_crowd,pt.birthday_pre_time").
		Joins("LEFT JOIN "+micro_scrm_enterprise.PlanTaskModel{}.AllTableName(enterpriseID)+" pt ON pt.plan_id = plans.id").
		Joins("LEFT JOIN "+micro_scrm_enterprise.PlanTaskConditionModel{}.AllTableName(enterpriseID)+" ptc ON ptc.task_id = pt.id").
		Where("plans.status = ? AND pt.status = ? AND pt.dimension=? AND pt.condition_type=? AND plans.deleted_at = ? AND pt.deleted_at = ? AND ptc.deleted_at = ? AND plans.category=1  ", 1, 1, 2, 8, 0, 0, 0).Find(&result).Error
	return
}

// GetTemplatePushData 获取导入客户模板的子任务
func (p *PlanTaskService) GetTemplatePushData(enterpriseID int64, info micro_scrm_enterprise.PlanModel) (result []micro_scrm_enterprise.ResponseTemplateCondition, err error) {
	err = p.Scopes(model.BranchTable(&info, enterpriseID)).
		Select("plans.is_push,plans.id AS plan_id,plans.name AS plan_name,pt.`name` AS task_name ,ptc.task_id,ptc.id AS condition_id,ptc.conditions,plans.start_time,plans.end_time,"+
			"ptc.content,ptc.content_type,pt.task_type,pt.send_type,pt.task_is_push,pt.is_specify_crowd,pt.consumer_file_id,pt.is_extract,pt.task_start_time,pt.task_end_time,pt.task_expire_type").
		Joins("LEFT JOIN "+micro_scrm_enterprise.PlanTaskModel{}.AllTableName(enterpriseID)+" pt ON pt.plan_id = plans.id").
		Joins("LEFT JOIN "+micro_scrm_enterprise.PlanTaskConditionModel{}.AllTableName(enterpriseID)+" ptc ON ptc.task_id = pt.id").
		Where("plans.status = ? AND pt.status = ? AND pt.dimension IN ? AND plans.deleted_at = ? AND pt.deleted_at = ? AND ptc.deleted_at = ? AND plans.category=1 ", 1, 1, []int{3, 4}, 0, 0, 0).Find(&result).Error
	return
}

// UpdatePushStatus 更新推送记录
func (p *PlanTaskService) UpdatePushStatus(enterpriseID, planId int64, info micro_scrm_enterprise.PlanModel) error {
	return p.Scopes(model.BranchTable(&info, enterpriseID)).Where("id = ?", planId).Update("is_push", 1).Error
}

func (p *PlanTaskService) GetTagInfo(tagName string, info micro_scrm.WorkExternalTagModel) (tag string, err error) {
	err = p.Model(&info).Select([]string{"tag_id"}).Where("name = ?", tagName).First(&tag).Error
	return
}

// 获取信息
func (p *PlanTaskService) GetCronData(enterpriseID int64, info micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel, tagName string, startTime, endTime int64, pagSize, pageNUm int) (result []*micro_scrm_enterprise.JobReceive, err error) {
	err = p.Scopes(model.BranchTable(&info, enterpriseID)).
		Select("DISTINCT campaign_tool_plan_task_cron_data.tag_name,wuce.unionid AS union_id,wuce.openid AS work_external_user_id,"+
			"m.userid AS user_id,m.mobile,tag.tag_id,weur.id AS relation_id,tag.id AS tag_table_id").
		Joins("INNER JOIN work_external_tag tag ON tag.name = campaign_tool_plan_task_cron_data.tag_name AND tag.deleted_at = 0").
		Joins("INNER JOIN "+micro_scrm_enterprise.WxUserCorpExternalModel{}.AllTableName(enterpriseID)+" wuce ON "+
			"wuce.union_code = campaign_tool_plan_task_cron_data.member_code AND wuce.deleted_at = 0").
		Joins("INNER JOIN work_external_user_relation weur ON weur.work_external_user_id = wuce.id AND weur.deleted_at = 0").
		Joins("INNER JOIN members m ON m.id = weur.member_id").
		Where("campaign_tool_plan_task_cron_data.tag_name = ? AND "+
			"m.enterprise_id = ?  AND m.deleted_at = 0 AND "+
			"campaign_tool_plan_task_cron_data.deleted_at = 0", tagName, enterpriseID).
		Limit(pagSize).Offset(pagSize * (pageNUm - 1)).
		Find(&result).Error

	return
}

func (p *PlanTaskService) CreateAllotLog(enterpriseID int64, info micro_scrm_enterprise.TaskAllotLogModel) (int64, error) {
	err := p.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).Create(&info).Error
	return info.ID, err
}

func (p *PlanTaskService) UpdateTaskPushStatus(enterpriseID int64, taskIds []int64) error {
	return p.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).Where("id IN (?)", taskIds).Update("task_is_push", 1).Error
}

func (p *PlanTaskService) GetUserTagInfo(relationID, tagID int64, info work_external_user.WorkExternalUserTagModel) (result int64, err error) {
	err = p.Model(&info).Scopes(service.ScopeCond_EQDelete).Select([]string{"id"}).Where("relation_id = ? AND external_tag_id = ?", relationID, tagID).
		First(&result).Error
	return
}

func (p *PlanTaskService) CreateUserTag(info work_external_user.WorkExternalUserTagModel) (err error) {
	err = p.Create(&info).Error
	return
}

func (p *PlanTaskService) UpdateAllotInfo(enterpriseID int64, planID []int64, info micro_scrm_enterprise.TaskAllotLogModel) (err error) {
	return p.Scopes(model.BranchTable(&info, enterpriseID)).Where("is_finished = ? AND plan_id in (?)", 0, planID).Update("is_finished", 2).Error
}

func (p *PlanTaskService) UpdateTagStatus(enterpriseID, externalTagID, relationID int64) (err error) {
	err = p.Model(&work_external_user.WorkExternalUserTagModel{}).Where("external_tag_id = ? AND relation_id = ? and enterprise_id=?", externalTagID, relationID, enterpriseID).Update("deleted_at", time.Now().Unix()).Error
	return
}

func (p *PlanTaskService) FindplanTaskTagData(enterpriseID, tagId int64, userId string) (data []*plan_task.TaskConditionDataLogModel, err error) {
	err = p.Scopes(model.BranchTable(&plan_task.TaskConditionDataLogModel{}, enterpriseID)).Where("tag_id=? AND user_id=?  AND deleted_at=?", tagId, userId, 0).Group("tag_id,relation_id,union_id").Find(&data).Error
	return
}

// UpdateplanTaskTagData  更新设置根据标签筛选的用户已经被使用了
func (p *PlanTaskService) UpdateplanTaskTagData(enterpriseID, tagId int64, userId string, RelationId int64) error {
	return p.Scopes(model.BranchTable(&plan_task.TaskConditionDataLogModel{}, enterpriseID)).Where("tag_id=? AND user_id=?  AND relation_id=?", tagId, userId, RelationId).Update("deleted_at", time.Now().Unix()).Error

}

func (p *PlanTaskService) FindtaskAllotlogByTaskIdAndDate(enterpriseID, taskId, startTime, endTime int64) (data micro_scrm_enterprise.TaskAllotLogModel, err error) {
	err = p.Scopes(model.BranchTable(&micro_scrm_enterprise.TaskAllotLogModel{}, enterpriseID)).Where("task_id=? AND created_at BETWEEN ? AND ? AND deleted_at=?", taskId, startTime, endTime, 0).First(&data).Error
	return
}
