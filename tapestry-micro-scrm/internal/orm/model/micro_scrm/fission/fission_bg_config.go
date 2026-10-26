package fission

import "micro-scrm/internal/orm/model"

type FissionBgConfigModel struct {
	EnterpriseID int64  `gorm:"column:enterprise_id;type:bigint(20);not null;index:enterprise_id_index;default:0"`
	BgImg        string `gorm:"column:bg_img"`
	ShareImg     string `gorm:"column:share_img"`
	model.MicroScrmBaseModel
}

func (FissionBgConfigModel) TableName() string {
	return "fission_bg_config"
}
