package user_group

import (
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

// GetFileNameDetection 验证客户分组导入模板是否已存在
func (u *UserGroupService) GetFileNameDetection(enterpriseID int64, fileName string) (result int64, err error) {
	err = u.Scopes(model.BranchTable(micro_scrm_enterprise.UserGroupFileModel{}, enterpriseID)).
		Select("COUNT(1)").
		Where("file_name = ? AND deleted_at = 0 AND is_extract != 0", fileName).
		Find(&result).Error
	return
}

// CreateUploadInfo 导入客户的模板文件记录保存
func (u *UserGroupService) CreateUploadInfo(enterpriseID int64, info micro_scrm_enterprise.UserGroupFileModel) (ID int64, err error) {
	err = u.Scopes(model.BranchTable(micro_scrm_enterprise.UserGroupFileModel{}, enterpriseID)).
		Create(&info).Error
	ID = info.ID
	return
}

func (u *UserGroupService) CreateCSVInfo(enterpriseID int64, data []micro_scrm_enterprise.UserGroupImportDataModel) (err error) {
	err = u.Scopes(model.BranchTable(micro_scrm_enterprise.UserGroupImportDataModel{}, enterpriseID)).
		Create(&data).Error
	return
}

func (u *UserGroupService) UpdateISExtract(enterpriseID, ID int64) (err error) {
	err = u.Scopes(model.BranchTable(micro_scrm_enterprise.UserGroupFileModel{}, enterpriseID)).
		Where("id = ?", ID).Update("is_extract", 1).Error
	return
}

func (u *UserGroupService) GetImportFile(enterpriseID, ID int64) (data micro_scrm_enterprise.UserGroupFileModel, err error) {
	err = u.Scopes(model.BranchTable(micro_scrm_enterprise.UserGroupFileModel{}, enterpriseID)).
		Where("id = ?", ID).First(&data).Error
	return
}

func (u *UserGroupService) GetImportData(enterpriseID, FileID int64, Page, PageSize int) (data []micro_scrm_enterprise.UserGroupImportDataModel, err error) {
	err = u.Scopes(model.BranchTable(micro_scrm_enterprise.UserGroupImportDataModel{}, enterpriseID)).
		Where("file_id = ?", FileID).Offset(Page).Limit(PageSize).Find(&data).Error
	return
}

func (u *UserGroupService) GetImportDataCount(enterpriseID, FileID int64) (total int64, err error) {
	err = u.Scopes(model.BranchTable(micro_scrm_enterprise.UserGroupImportDataModel{}, enterpriseID)).
		Where("file_id = ?", FileID).Count(&total).Error
	return
}

func (u *UserGroupService) CreateExportLog(enterpriseID int64, data []micro_scrm_enterprise.UserGroupImportDataLogModel) (err error) {
	err = u.Scopes(model.BranchTable(micro_scrm_enterprise.UserGroupImportDataLogModel{}, enterpriseID)).
		Create(&data).Error
	return
}

func (u *UserGroupService) UpdateExportLogStatus(enterpriseID, GroupID int64, unionId string) error {
	return u.Scopes(model.BranchTable(micro_scrm_enterprise.UserGroupImportDataLogModel{}, enterpriseID)).
		Where(map[string]interface{}{
			"user_group_id": GroupID,
			"union_id":      unionId}).Update("import_res", "成功").Error
}

func (u *UserGroupService) GetUserGroupLogCount(enterpriseID, GroupID int64) (total int64, err error) {
	err = u.Scopes(model.BranchTable(micro_scrm_enterprise.UserGroupImportDataLogModel{}, enterpriseID)).
		Where("user_group_id = ?", GroupID).Count(&total).Error
	return
}

func (u *UserGroupService) GetUserGroupLogList(enterpriseID, GroupID int64, Page, PageSize int) (data []entity.ExportDataInfo, err error) {
	offset := (Page - 1) * PageSize
	err = u.Scopes(model.BranchTable(micro_scrm_enterprise.UserGroupImportDataLogModel{}, enterpriseID)).
		Select("union_id", "import_res").
		Where("user_group_id = ?", GroupID).Offset(offset).Limit(PageSize).Find(&data).Error
	return
}
