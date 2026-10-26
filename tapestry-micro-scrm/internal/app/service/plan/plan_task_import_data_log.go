package planService

import (
	"micro-scrm/internal/app"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	planModels "micro-scrm/internal/app/models/micro_scrm"
)

type planTaskImportDataLog struct {
	planTaskImportDataLogDao *planDao.PlanTaskImportDataLogDao
}

func NewPlanTaskImportDataLog() *planTaskImportDataLog {
	return &planTaskImportDataLog{
		planTaskImportDataLogDao: planDao.NewPlanTaskImportDataLogDao(),
	}
}

func (p planTaskImportDataLog) GetImportDataLogCount(c app.Ctx, enterpriseID, fileID int64) (total int64, err error) {
	return p.planTaskImportDataLogDao.GetImportDataLogCount(c, enterpriseID, fileID)
}

func (p planTaskImportDataLog) GetImportDataLogList(c app.Ctx, enterpriseID, fileID int64, page int, size int) (data []planModels.PlanTaskImportDataLogModel, err error) {
	return p.planTaskImportDataLogDao.GetImportDataLogList(c, enterpriseID, fileID, page, size)
}
