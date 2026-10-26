package plan_task

import (
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

// InsertCampaignToolData 新增Campaign-Tool推送的会员数据
func (w *PlanTaskService) InsertCampaignToolData(enterpriseID int64, newData []micro_scrm_enterprise.CampaignToolCronDataModel) error {
	return w.Scopes(model.BranchTable(&micro_scrm_enterprise.CampaignToolCronDataModel{}, enterpriseID)).Create(&newData).Error
}

// InsertCampaignToolPlanTaskData 新增Campaign-Tool推送的标签数据
func (w *PlanTaskService) InsertCampaignToolPlanTaskData(enterpriseID int64, newData []micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel) error {
	return w.Scopes(model.BranchTable(&micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel{}, enterpriseID)).Create(&newData).Error
}

// SelectCampaignToolPlanTaskDataUnionId  查询出数据campaign-tool中的union_id
func (w *PlanTaskService) SelectCampaignToolPlanTaskDataUnionId(enterpriseID int64, fileId int) (data []micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel, err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel{}, enterpriseID)).Where("file_id=?", fileId).Find(&data).Error
	return data, err
}
