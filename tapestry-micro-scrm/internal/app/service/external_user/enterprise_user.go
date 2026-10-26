package externalUserService

import (
	"micro-scrm/internal/app"
	externalUserDao "micro-scrm/internal/app/dao/micro_scrm"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
)

type EnterpriseUser struct {
	wxUserCorpExternalDao *externalUserDao.WxUserCorpExternal
}

func NewEnterpriseUser() *EnterpriseUser {
	return &EnterpriseUser{
		wxUserCorpExternalDao: externalUserDao.NewWxUserCorpExternal(),
	}
}

func (d EnterpriseUser) GetWxUserCorpExternalCount(c app.Ctx, enterpriseID int64) (total int64, err error) {
	return d.wxUserCorpExternalDao.GetCount(c, enterpriseID)
}

func (d EnterpriseUser) GetWxUserCorpExternalList(c app.Ctx, enterpriseID int64, Page, Size int) ([]microScrmModel.WxUserCorpExternalModel, error) {
	return d.wxUserCorpExternalDao.GetList(c, enterpriseID, Page, Size)
}
