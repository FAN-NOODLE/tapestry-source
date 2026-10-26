package source

import (
	"micro-scrm/internal/app"
	sourceDao "micro-scrm/internal/app/dao/source"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	"micro-scrm/internal/app/models/source"
	"time"
)

type FileService struct {
	sourceDao *sourceDao.Dao
	fileDao   *sourceDao.FileDao
}

func NewFileService() FileService {
	return FileService{
		sourceDao: sourceDao.NewSourceDao(),
		fileDao:   sourceDao.NewFileDao(),
	}
}

func (s FileService) Detail(c app.Ctx, request sourceRequest.DetailRequest) ([]source.FileCategoryRelationDetail, error) {
	return s.sourceDao.GetSourceFileDetail(c, request.SourceID, request.EnterpriseID)
}

func (s FileService) Update(c app.Ctx, request sourceRequest.FileUpdateRequest) error {
	// 更新文件名
	query := map[string]interface{}{
		"id":            request.ID,
		"enterprise_id": request.EnterpriseID,
	}
	column := map[string]interface{}{
		"name":       request.Name,
		"updated_at": time.Now().Unix(),
	}
	if err := s.fileDao.Update(c, query, column); err != nil {
		return err
	}

	request.SourceType = source.TypeFile
	err := NewService().Update(c, request.Request)
	return err
}
