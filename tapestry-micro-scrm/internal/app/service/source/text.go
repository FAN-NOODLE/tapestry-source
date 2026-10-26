package source

import (
	"micro-scrm/internal/app"
	sourceDao "micro-scrm/internal/app/dao/source"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	"micro-scrm/internal/app/models/source"
	"time"
)

type TextService struct {
	sourceDao *sourceDao.Dao
	textDao   *sourceDao.TextDao
}

func NewTextService() TextService {
	return TextService{
		sourceDao: sourceDao.NewSourceDao(),
		textDao:   sourceDao.NewTextDao(),
	}
}

func (s TextService) Detail(c app.Ctx, request sourceRequest.DetailRequest) ([]source.TextCategoryRelationDetail, error) {

	return s.sourceDao.GetSourceTextDetail(c, request.SourceID, request.EnterpriseID)

}

func (s TextService) Update(c app.Ctx, request sourceRequest.TextUpdateRequest) error {
	// 更新文本内容
	query := map[string]interface{}{
		"id":            request.ID,
		"enterprise_id": request.EnterpriseID,
	}
	column := map[string]interface{}{
		"content":    request.Content,
		"updated_at": time.Now().Unix(),
	}
	if err := s.textDao.Update(c, query, column); err != nil {
		return err
	}

	request.SourceType = source.TypeText
	return NewService().Update(c, request.Request)
}
