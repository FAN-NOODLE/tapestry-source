package brand_config

import (
	"micro-scrm/internal/orm/model/micro_scrm"
)

// 获取品牌SFTP配置

func (brandconf *BrandConfigService) GetSFTPConfig(enterpriseID []int64, sftpType int) (result []micro_scrm.BrandConfigSftpModel, err error) {

	err = brandconf.Model(&micro_scrm.BrandConfigSftpModel{}).
		Where("enterprise_id IN (?) AND sftp_type =? AND deleted_at=?", enterpriseID, sftpType, 0).Find(&result).Error
	return
}
