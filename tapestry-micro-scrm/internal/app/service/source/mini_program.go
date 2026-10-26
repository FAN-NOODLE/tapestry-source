package source

import (
	"micro-scrm/internal/app"
	sourceDao "micro-scrm/internal/app/dao/source"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	"micro-scrm/internal/app/models/source"
	"time"
)

type MiniProgramService struct {
	sourceDao      *sourceDao.Dao
	miniProgramDao *sourceDao.MiniProgramDao
}

func NewMiniProgramService() MiniProgramService {
	return MiniProgramService{
		sourceDao:      sourceDao.NewSourceDao(),
		miniProgramDao: sourceDao.NewMiniProgramDao(),
	}
}

func (s MiniProgramService) Detail(c app.Ctx, request sourceRequest.DetailRequest) ([]source.MiniProgramCategoryRelationDetail, error) {
	return s.sourceDao.GetSourceMiniProgramDetail(c, request.SourceID, request.EnterpriseID)
}

func (s MiniProgramService) Update(c app.Ctx, request sourceRequest.MiniProgramUpdateRequest) error {
	// 更新小程序相关信息
	query := map[string]interface{}{
		"id":            request.ID,
		"enterprise_id": request.EnterpriseID,
	}
	column := map[string]interface{}{
		"mini_program_id": request.MiniProgramID,
		"title":           request.Title,
		"path":            request.Path,
		"cover_url":       request.CoverUrl,
		"updated_at":      time.Now().Unix(),
	}
	if err := s.miniProgramDao.Update(c, query, column); err != nil {
		return err
	}

	request.SourceType = source.TypeMiniProgram
	err := NewService().Update(c, request.Request)
	return err
}
