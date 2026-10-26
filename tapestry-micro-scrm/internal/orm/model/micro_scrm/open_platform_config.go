package micro_scrm

import "micro-scrm/internal/orm/model"

type OpenPlatformConfigModel struct {
	model.MicroScrmBaseModel

	Name         string `json:"name" gorm:"column:name;type:varchar(32);not null;comment:接入方名称"`
	ClientID     string `json:"client_id" gorm:"column:client_id;type:varchar(32);not null;index:client_id_index;comment:对外提供的身份ID"`
	ClientSecret string `json:"client_secret" gorm:"column:client_secret;type:varchar(32);not null;comment:对外提供的身份密钥"`
	Salt         string `json:"salt" gorm:"column:salt;type:varchar(6);not null;comment:随机盐，6位"`
	Remark       string `json:"remark" gorm:"column:remark;type:varchar(100);not null;comment:接入方名称"`
}

func (OpenPlatformConfigModel) TableName() string {
	return "open_platform_config"
}
