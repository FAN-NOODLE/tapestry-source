package excitationModels

import "micro-scrm/internal/app/models"

type ExcitationCoinsSettingsEditRecordModel struct {
	models.MicroBaseModel
	AdminID     int64  `json:"admin_id" gorm:"column:admin_id;type:bigint(20);not null;default:0"`
	AdminName   string `json:"admin_name" gorm:"column:admin_name;type:varchar(50);not null;comment:操作人标识"`
	RatioRemark string `json:"ratio_remark" gorm:"column:ratio_remark;type:varchar(50);not null;comment:修改内容"`
	MallIsOpen  int32  `json:"mall_is_open" gorm:"column:mall_is_open;type:tinyint(1);not null;default:0;comment:是否开启商城 0 否 1是"`
	ExpireAt    int64  `json:"expire_at" gorm:"column:expire_at;type:int(10);not null;default:0;comment:过期时间"`
}

// TableName returns the table name of the ExcitationCoinsSettingsEditRecord model
func (ExcitationCoinsSettingsEditRecordModel) TableName() string {
	return "excitation_coins_settings_edit_record"
}
