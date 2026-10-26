package thirdCrmSftpDao

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type ThirdCrmSftpFile struct {
	dao       dao.Dao
	tableName string
}

func NewThirdCrmSftpFile() *ThirdCrmSftpFile {
	return &ThirdCrmSftpFile{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: crmSftpFileModels.CrmSftpDownloadedFileModel{}.TableName(),
	}
}

func (d *ThirdCrmSftpFile) GetExistsFile(c app.Ctx, enterpriseID int64, fileType []int) (result []crmSftpFileModels.CrmSftpDownloadedFileModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ? AND file_type in ? AND deleted_at = 0 ", enterpriseID, fileType).
		Find(&result).Error
	return
}

func (d *ThirdCrmSftpFile) GetNotExtractFile(c app.Ctx, enterpriseID int64, fileType []int) (result []crmSftpFileModels.CrmSftpDownloadedFileModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ? AND file_type in ? AND is_extract = 0 AND deleted_at = 0 ", enterpriseID, fileType).
		Find(&result).Error
	return
}

func fileIdentityConditions(enterpriseID int64, fileType int, fullPath, fileName string) map[string]interface{} {
	return map[string]interface{}{
		"enterprise_id":  enterpriseID,
		"file_type":      fileType,
		"file_full_path": fullPath,
		"file_name":      fileName,
		"deleted_at":     0,
	}
}

// GetFileByIdentity 按企业、业务类型、完整路径和文件名获取唯一文件处理记录。
func (d *ThirdCrmSftpFile) GetFileByIdentity(c app.Ctx, enterpriseID int64, fileType int, fullPath, fileName string) (result crmSftpFileModels.CrmSftpDownloadedFileModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where(fileIdentityConditions(enterpriseID, fileType, fullPath, fileName)).
		First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return result, gorm.ErrRecordNotFound
	}
	return result, err
}

func (d *ThirdCrmSftpFile) SaveFile(c app.Ctx, data []crmSftpFileModels.CrmSftpDownloadedFileModel) error {
	return d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
}

func (d *ThirdCrmSftpFile) UpdateFileExTract(c app.Ctx, enterpriseID int64, id int64, status int32) error {
	return d.dao.WithContext(c).Table(d.tableName).Where("id=? AND enterprise_id=?", id, enterpriseID).Updates(map[string]interface{}{
		"is_extract": status,
		"updated_at": time.Now().Unix(),
	}).Error
}

func (d *ThirdCrmSftpFile) CreateData(c app.Ctx, data crmSftpFileModels.CrmSftpDownloadedFileModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
	return err
}
