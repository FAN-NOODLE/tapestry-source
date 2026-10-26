package enterprise

import "micro-scrm/internal/orm/model"

type EnterpriseConfigModel struct {
	model.MicroScrmBaseModel

	Hash  string `gorm:"column:hash;type:varchar(64);not null;comment:hash"`
	Key   string `gorm:"column:key;type:varchar(100);uniqueIndex:idx_enterprise;comment:配置项"`
	Value string `gorm:"column:value;type:varchar(1024);not null;comment:配置值"`
	IsUse uint32 `gorm:"column:is_use;type:tinyint(1) unsigned;not null;default:1"`
}

func (EnterpriseConfigModel) TableName() string {
	return "enterprise_config"
}
