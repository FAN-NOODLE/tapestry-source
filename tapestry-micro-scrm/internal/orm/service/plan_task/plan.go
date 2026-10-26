package plan_task

import (
	"log"

	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

//@author :LiKun
//function : 新增计划
//start:2021-07-28
//update:-
//info : micro_scrm_enterprise.PlanModel
func (p *PlanTaskService) AddPlan(info micro_scrm_enterprise.PlanModel, enterpriseID int64) (planid int64, err error) {
	if err := p.DB.Scopes(model.BranchTable(&info, enterpriseID)).Create(&info).Error; err != nil {
		return 0, err
	}
	return info.ID, nil
}

//@author :LiKun
//function : 计划编辑
//start:2021-07-28
//update:-
//info : micro_scrm_enterprise.PlanModel
func (p *PlanTaskService) PlanEdit(info micro_scrm_enterprise.PlanModel, enterpriseID int64) (int64, error) {
	if err := p.Scopes(model.BranchTable(&info, enterpriseID)).Save(&info).Error; err != nil {
		return 0, err
	}
	return info.ID, nil
}

//@author :LiKun
//function : 任务计划删除
//start:2021-07-28
//update:-
//info : id int
func (p *PlanTaskService) DelPlan(id int64, deleteTime int64, enterpriseID int64) error {
	if err := p.DB.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanModel{}, enterpriseID)).
		Model(&micro_scrm_enterprise.PlanModel{}).Where("id = ?", id).Update("deleted_at", deleteTime).Error; err != nil {
		return err
	}
	return nil
}

//获取子任务id
func (p *PlanTaskService) GetPlanTaskId(planId, enterpriseID int64) ([]int64, error) {
	var planTaskId = make([]int64, 0)
	if err := p.DB.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).Model(&micro_scrm_enterprise.PlanTaskModel{}).Select("id").Where("plan_id = ?", planId).Pluck("id", &planTaskId).Error; err != nil {
		return nil, err
	}
	return planTaskId, nil
}

func (p *PlanTaskService) DelPlanTask(id []int64, deleteTime, enterpriseID int64) error {
	if err := p.DB.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).Where("id in (?)", id).Update("deleted_at", deleteTime).Error; err != nil {
		return err
	}
	return nil
}

//@author :LiKun
//function : 任务计划复制
//start:2021-07-28
//update:-
//info :*micro_scrm_enterprise.PlanModel

func (p *PlanTaskService) PlanCopy(info micro_scrm_enterprise.PlanModel, enterpriseID int64) (int64, error) {
	if err := p.DB.Scopes(model.BranchTable(&info, enterpriseID)).Create(&info).Error; err != nil {
		return 0, err
	}
	return info.ID, nil
}

func (p *PlanTaskService) PlanFound(id, enterpriseID int64) (info micro_scrm_enterprise.PlanModel, err error) {
	err = p.Scopes(model.BranchTable(&info, enterpriseID)).Where("id = ?", id).Find(&info).Error

	return
}

//获取对应子任务列表
func (p *PlanTaskService) PlanTaskFound(planId, enterpriseID int64) (result []micro_scrm_enterprise.PlanTaskModel, err error) {
	if err := p.DB.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).Where("plan_id = ? AND deleted_at = ?", planId, 0).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

//子任务复制
func (p *PlanTaskService) PlanTaskCopy(planTaskInfo []micro_scrm_enterprise.PlanTaskModel, enterpriseID int64) error {
	if err := p.DB.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskModel{}, enterpriseID)).CreateInBatches(planTaskInfo, len(planTaskInfo)).Error; err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (p *PlanTaskService) GetPlanTaskCondition(planId int64, enterpriseID int64) (result []micro_scrm_enterprise.PlanTaskConditionModel, err error) {
	if err = p.DB.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, enterpriseID)).
		Where("plan_id = ? AND deleted_at = ?", planId, 0).
		Find(&result).Error; err != nil {
		return
	}
	return
}
func (p *PlanTaskService) PlanTaskConditionCopy(planTaskConditionInfo []micro_scrm_enterprise.PlanTaskConditionModel, enterpriseID int64) error {
	if err := p.DB.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanTaskConditionModel{}, enterpriseID)).
		CreateInBatches(planTaskConditionInfo, len(planTaskConditionInfo)).
		Error; err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//更新计划状态
func (p *PlanTaskService) UpdatePlanStatus(planId, status, enterpriseID int64) error {
	if err := p.DB.Scopes(model.BranchTable(&micro_scrm_enterprise.PlanModel{}, enterpriseID)).
		Where("id = ?", planId).Update("is_open", status).Error; err != nil {
		return err
	}
	return nil
}
