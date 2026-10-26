package fission

import "micro-scrm/internal/orm/model"

type FissionStyleConfigModel struct {
	EnterpriseID int64  `gorm:"column:enterprise_id;type:bigint(20);not null;index:enterprise_id_index;default:0"`
	RgbColors    string `gorm:"column:rgb_colors"`
	Styles       string `gorm:"column:styles"`
	IsCheck      int64  `gorm:"column:is_check"`
	model.MicroScrmBaseModel
}

func (FissionStyleConfigModel) TableName() string {
	return "fission_style_config"
}
