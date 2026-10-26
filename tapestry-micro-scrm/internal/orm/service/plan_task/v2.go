package plan_task

import (
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/plan_task"
)

func (db *PlanTaskService) GetCTDownloadedFileList(enterpriseID int64, fileType uint8) (fileList []string, err error) {
	err = db.Model(&plan_task.PlanTaskDownloadedFileModel{}).
		Scopes(model.BranchTable(&plan_task.PlanTaskDownloadedFileModel{}, enterpriseID)).
		Where(map[string]interface{}{
			"file_type": fileType,
		}).Pluck("file_name", &fileList).Error

	return
}
