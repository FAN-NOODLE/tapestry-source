package work_external_user

import (
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/service"
)

func (w *WorkExternalUserService) CreateSyncExternalUserSftp(enterpriseID int64, data micro_scrm_enterprise.SyncExternalUserSftpModel) error {
	return w.Scopes(model.BranchTable(&micro_scrm_enterprise.SyncExternalUserSftpModel{}, enterpriseID)).Create(&data).Error
}
func (w *WorkExternalUserService) FindSyncExternalUserSftpLastTime(enterpriseID int64) (data micro_scrm_enterprise.SyncExternalUserSftpModel, err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.SyncExternalUserSftpModel{}, enterpriseID), service.ScopeCond_EQDelete).Last(&data).Error
	return
}
