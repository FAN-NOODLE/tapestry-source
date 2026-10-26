package excitationModels

import "micro-scrm/internal/app/models"

type ExcitationCoinsExpireRecordModel struct {
	models.MicroBaseModel
	ExpiredAt int64 `gorm:"column:expired_at;type:int(10);not null;default:0" json:"expired_at"`
}

// TableName returns the table name of the ExcitationCoinsExpireRecord model
func (ExcitationCoinsExpireRecordModel) TableName() string {
	return "excitation_coins_expire_record"
}
