package plan_task

import (
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/plan_task"
	"time"
)

func (p *PlanTaskService) CreateCSVInfo(enterpriseID int64, info []micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel) (err error) {
	err = p.Scopes(model.BranchTable(micro_scrm_enterprise.CampaignToolPlanTaskCronDataModel{}, enterpriseID)).
		Create(&info).Error
	return
}

func (p *PlanTaskService) CreateUploadInfo(enterpriseID int64, info plan_task.PlanTaskDownloadedFileModel) (ID int64, err error) {
	err = p.Scopes(model.BranchTable(plan_task.PlanTaskDownloadedFileModel{}, enterpriseID)).
		Create(&info).Error
	ID = info.ID
	return
}

func (p *PlanTaskService) GetFileInfo(enterpriseID int64, fileName string) (result int64, err error) {
	err = p.Scopes(model.BranchTable(plan_task.PlanTaskDownloadedFileModel{}, enterpriseID)).
		Select("COUNT(1)").
		Where("file_name = ? AND deleted_at = 0 AND is_extract != 0", fileName).
		Find(&result).Error
	return
}
func (p *PlanTaskService) UpdatedFileInfo(enterpriseID int64, fileID int64) (err error) {
	err = p.Scopes(model.BranchTable(plan_task.PlanTaskDownloadedFileModel{}, enterpriseID)).
		Where("id = ?", fileID).Update("deleted_at", time.Now().Unix()).Error
	return
}

func (p *PlanTaskService) UpdateISExtract(enterpriseID, ID int64) (err error) {
	err = p.Scopes(model.BranchTable(plan_task.PlanTaskDownloadedFileModel{}, enterpriseID)).
		Where("id = ?", ID).Update("is_extract", 1).Error
	return
}
