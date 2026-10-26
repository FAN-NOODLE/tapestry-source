package planService

import (
	"micro-scrm/internal/app"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	planModels "micro-scrm/internal/app/models/micro_scrm"
)

type planTaskImportData struct {
	planTaskImportDataDao *planDao.PlanTaskImportDataDao
}

func NewPlanTaskImportData() *planTaskImportData {
	return &planTaskImportData{
		planTaskImportDataDao: planDao.NewPlanTaskImportDataDao(),
	}
}

func (p planTaskImportData) GetImportDataCount(c app.Ctx, enterpriseID, fileID int64) (total int64, err error) {
	return p.planTaskImportDataDao.GetImportDataCount(c, enterpriseID, fileID)
}

func (p planTaskImportData) GetImportDataList(c app.Ctx, enterpriseID, fileID int64, page int, size int) (data []planModels.PlanTaskImportDataModel, err error) {
	return p.planTaskImportDataDao.GetImportDataList(c, enterpriseID, fileID, page, size)
}
