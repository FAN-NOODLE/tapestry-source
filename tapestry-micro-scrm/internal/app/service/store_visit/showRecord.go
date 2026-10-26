package storeVisitSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	storeVisitDao "micro-scrm/internal/app/dao/store_visit"
	storeVisitReq "micro-scrm/internal/app/http/request/store_visit"
	storeVisitResp "micro-scrm/internal/app/http/response/store_visit"
	"micro-scrm/internal/app/models"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
	"micro-scrm/internal/app/service"
	"time"
)

type ShowRecordSvc struct {
	service.Base
	visitRecordDao *storeVisitDao.ShowRecordDao
	visitMediaDao  *storeVisitDao.MediaDao
}

func NewShowRecordSvc() *ShowRecordSvc {
	return &ShowRecordSvc{
		visitRecordDao: storeVisitDao.NewShowRecordDao(),
		visitMediaDao:  storeVisitDao.NewMediaDao(),
	}
}

func (s *ShowRecordSvc) SaveShowRecordContent(c ctx.Context, req storeVisitReq.SaveShowRecordContentReq) (id int64, err error) {
	if req.ID > 0 {
		err = s.visitRecordDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
			"id":        req.ID,
			"member_id": req.MemberId,
		}, map[string]interface{}{
			"content":    req.Content,
			"type_name":  req.TypeName,
			"updated_at": time.Now().Unix(),
		})
	} else {
		req.ID = app.UID()
		err = s.visitRecordDao.CreateData(c, req.EnterpriseId, storeVisitModel.StoreVisitShowRecordModel{
			BaseModel: models.CreateBase(req.ID),
			VisitID:   req.VisitID,
			TypeName:  req.TypeName,
			MemberID:  req.MemberId,
			Content:   req.Content,
		})
	}
	return req.ID, err
}

func (s *ShowRecordSvc) SaveShowRecordMedia(c ctx.Context, req storeVisitReq.SaveShowRecordMediaReq) (info storeVisitResp.ShowRecordMediaInfo, err error) {
	mediaID := app.UID()
	if req.ID == 0 {
		req.ID = app.UID()
		err = s.visitRecordDao.CreateData(c, req.EnterpriseId, storeVisitModel.StoreVisitShowRecordModel{
			BaseModel: models.CreateBase(req.ID),
			VisitID:   req.VisitID,
			TypeName:  req.TypeName,
			MemberID:  req.MemberId,
			Content:   "",
		})
	}
	err = s.visitMediaDao.CreateData(c, req.EnterpriseId, storeVisitModel.StoreVisitMediaModel{
		BaseModel:     models.CreateBase(mediaID),
		VisitID:       req.VisitID,
		MemberID:      req.MemberId,
		RelID:         req.ID,
		RelType:       storeVisitModel.ShowRecord,
		MediaType:     req.MediaType,
		URL:           req.Url,
		VideoImageUrl: req.VideoImageUrl,
		VideoSecond:   req.VideoSeconds,
	})
	info.ID = req.ID
	info.MediaID = mediaID
	return info, err
}

func (s *ShowRecordSvc) DelShowRecord(c ctx.Context, req storeVisitReq.DelShowRecordReq) (err error) {
	if req.ID == 0 {
		return nil
	}
	_ = s.visitRecordDao.DeleteData(c, req.EnterpriseId, map[string]interface{}{
		"id":        req.ID,
		"member_id": req.MemberId,
	})
	_ = s.visitMediaDao.DeleteData(c, req.EnterpriseId, map[string]interface{}{
		"visit_id":  req.ID,
		"member_id": req.MemberId,
	})

	return nil
}

func (s *ShowRecordSvc) DelShowRecordMedia(c ctx.Context, req storeVisitReq.DelShowRecordMediaReq) (err error) {
	if req.MediaID == 0 {
		return nil
	}
	return s.visitMediaDao.DeleteData(c, req.EnterpriseId, map[string]interface{}{
		"id":        req.MediaID,
		"member_id": req.MemberId,
	})
}

func (s *ShowRecordSvc) ShowRecordList(c ctx.Context, req storeVisitReq.VisitInfoReq) ([]storeVisitModel.StoreVisitShowRecordModel, error) {
	return s.visitRecordDao.GetListPreloadMedia(c, req.EnterpriseId, "visit_id = ?", req.ID)
}
