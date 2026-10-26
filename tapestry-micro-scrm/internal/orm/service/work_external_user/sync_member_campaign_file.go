package work_external_user

import (
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/service"
)

const (
	FILEPULLCOMPLETE = iota + 1
	FILESYNCCOMPLETE
)

func (w *WorkExternalUserService) CreateCTMemberFile(enterpriseID int64, fileModel micro_scrm_enterprise.SyncMemberCampaignFileModel) error {
	return w.Scopes(model.BranchTable(&micro_scrm_enterprise.SyncMemberCampaignFileModel{}, enterpriseID)).Create(&fileModel).Error
}
func (w *WorkExternalUserService) FindCTMemberFile(enterpriseID int64, isExtract int) (data micro_scrm_enterprise.SyncMemberCampaignFileModel, err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.SyncMemberCampaignFileModel{}, enterpriseID), service.ScopeCond_EQDelete).
		Where("is_extract=?", isExtract).First(&data).Error
	return
}

func (w *WorkExternalUserService) UpdateCTMemberSync(enterpriseID, Id int64) error {
	return w.Scopes(model.BranchTable(&micro_scrm_enterprise.SyncMemberCampaignFileModel{}, enterpriseID), service.ScopeCond_EQDelete).
		Where("id=?", Id).Update("is_extract", FILESYNCCOMPLETE).Error
}
