package storeVisitSvc

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	storeVisitDao "micro-scrm/internal/app/dao/store_visit"
	storeVisitReq "micro-scrm/internal/app/http/request/store_visit"
	storeVisitResp "micro-scrm/internal/app/http/response/store_visit"
	"micro-scrm/internal/app/models"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"time"
)

type ConfigSvc struct {
	service.Base
	visitConfigDao *storeVisitDao.ConfigDao
}

func NewConfigSvc() *ConfigSvc {
	return &ConfigSvc{
		visitConfigDao: storeVisitDao.NewConfigDao(),
	}
}
func (s *ConfigSvc) GetConfig(c ctx.Context, req storeVisitReq.VisitConfigReq) (storeVisitResp.ConfigInfoResp, error) {
	configModel, err := s.visitConfigDao.GetOne(c, req.EnterpriseId, "`key` = ?", req.Key)
	if err != nil {
		return storeVisitResp.ConfigInfoResp{}, err
	}
	memberConfig, _ := s.visitConfigDao.GetOne(c, req.EnterpriseId, "`key` = ?", helper.Int64ToString(req.MemberId))
	return storeVisitResp.ConfigInfoResp{
		Key:          configModel.Key,
		Value:        configModel.Value,
		MemberConfig: memberConfig.Value,
	}, nil
}

func (s *ConfigSvc) SaveLanguageConfig(c ctx.Context, req storeVisitReq.BaseReq) {
	configModel, _ := s.visitConfigDao.GetOne(c, req.EnterpriseId, "`key` = ?", helper.Int64ToString(req.MemberId))
	if configModel.ID > 0 {
		_ = s.visitConfigDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("id = %d", configModel.ID), map[string]any{
			"value":      req.Language,
			"updated_at": time.Now().Unix(),
		})
	} else {
		_ = s.visitConfigDao.CreateData(c, req.EnterpriseId, storeVisitModel.StoreVisitConfigModel{
			BaseModel: models.CreateBase(app.UID()),
			Key:       helper.Int64ToString(req.MemberId),
			Value:     req.Language,
		})
	}
}
