package tagsSvc

import (
	"micro-scrm/internal/app"
	tagsDao "micro-scrm/internal/app/dao/tags"
	tagModels "micro-scrm/internal/app/models/tag"
)

type ImpressionTagLog struct {
	impressionTagLogDao *tagsDao.ImpressionOperateLog
}

func NewImpressionTagLog() *ImpressionTagLog {
	return &ImpressionTagLog{
		impressionTagLogDao: tagsDao.NewTagsOperateLog(),
	}
}

func (s ImpressionTagLog) GetTagLog(c app.Ctx, enterpriseId int64) ([]tagModels.ImpressionTagOperateLogModel, error) {
	return s.impressionTagLogDao.GetLogList(c, enterpriseId)
}

func (s ImpressionTagLog) DeleteList(c app.Ctx, enterpriseId int64, createdAt int64) error {
	return s.impressionTagLogDao.DeleteList(c, enterpriseId, createdAt)
}
