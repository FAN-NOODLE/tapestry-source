package crmSftpFileModels

import "micro-scrm/internal/app/models"

type BrandConfigSftpModel struct {
	models.MicroBaseModel
	BrandCode  string `json:"brand_code" gorm:"column:brand_code;type:varchar(20);not null;comment:品牌简称"`
	BrandName  string `json:"brand_name" gorm:"column:brand_name;type:varchar(50);not null;comment:品牌名称"`
	SftpType   int32  `json:"sftp_type" gorm:"column:sftp_type;type:tinyint(1);default:0;comment:配置类型:1 计划任务 2 ba业绩看板 3 社群营销数据 4 安客诚 crm"`
	SftpConfig string `json:"sftp_config" gorm:"column:sftp_config;type:text;comment:SFTP配置"`
	SftpPath   string `json:"sftp_path" gorm:"column:sftp_path;type:text;not null"`
}

// TableName returns the table name of the BrandConfigSftp model
func (BrandConfigSftpModel) TableName() string {
	return "brand_config_sftp"
}
