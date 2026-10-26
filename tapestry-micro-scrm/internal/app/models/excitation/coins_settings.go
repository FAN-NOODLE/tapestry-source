package excitationModels

import "micro-scrm/internal/app/models"

const (
	CoinsExpireClose = 0 // 金币过期关闭
	CoinsExpireOpen  = 1 // 金币过期开启

)

type ExcitationCoinsSettingsModel struct {
	models.MicroBaseModel
	Ratio        float64 `json:"ratio" gorm:"column:ratio;type:decimal(10,6);not null;default:0.00000000;comment:积分金币兑换比率"`
	RatioRemark  string  `json:"ratio_remark" gorm:"column:ratio_remark;type:varchar(50);not null;comment:积分金币兑换比率数值"`
	MallIsOpen   int32   `json:"mall_is_open" gorm:"column:mall_is_open;type:tinyint(1);not null;default:0;comment:是否开启商城 0 否 1是"`
	ExpireAt     int64   `json:"expire_at" gorm:"column:expire_at;type:int(10);not null;default:0;comment:过期时间"`
	ExpireIsOpen int32   `json:"expire_is_open" gorm:"column:expire_is_open;type:tinyint(1);not null;default:0;comment:是否开启金币过期 0否 1是"`
}

// TableName returns the table name of the ExcitationCoinsSettings model
func (ExcitationCoinsSettingsModel) TableName() string {
	return "excitation_coins_settings"
}
