package externalSearchModel

type ExternalCustomerExtensionModel struct {
	Unionid                     string `json:"unionid" gorm:"primaryKey;column:unionid;type:varbinary(64);not null"`
	ValidReserveAt              int64  `json:"valid_reserve_at" gorm:"column:valid_reserve_at;type:int unsigned;not null;default:0"`
	SuccessReserveAt            int64  `json:"success_reserve_at" gorm:"column:success_reserve_at;type:int unsigned;not null;default:0"`
	TaskFinishAt                int64  `json:"task_finish_at" gorm:"column:task_finish_at;type:int unsigned;not null;default:0"`
	TaskCreatedAt               int64  `json:"task_created_at" gorm:"column:task_created_at;type:int unsigned;not null;default:0"`
	CouponStatus                int    `json:"coupon_status" gorm:"column:coupon_status;type:tinyint unsigned;not null;default:0"`
	CouponBy30daysStatus        int    `json:"coupon_by_30days_status" gorm:"column:coupon_by_30days_status;type:tinyint unsigned;not null;default:0"`
	BirthdayCouponStatus        int    `json:"birthday_coupon_status" gorm:"column:birthday_coupon_status;type:tinyint unsigned;not null;default:0"`
	RenewMembershipCouponStatus int    `json:"renew_membership_coupon_status" gorm:"column:renew_membership_coupon_status;type:tinyint unsigned;not null;default:0"`
	WelcomeCouponStatus         int    `json:"welcome_coupon_status" gorm:"column:welcome_coupon_status;type:tinyint unsigned;not null;default:0"`
	UpgradeCouponStatus         int    `json:"upgrade_coupon_status" gorm:"column:upgrade_coupon_status;type:tinyint unsigned;not null;default:0"`
	CouponUpdatedAt             int64  `json:"coupon_updated_at" gorm:"column:coupon_updated_at;type:int unsigned;not null;default:0"`
	CreatedAt                   int64  `json:"created_at" gorm:"column:created_at;type:int unsigned;not null;default:0"`
	UpdatedAt                   int64  `json:"updated_at" gorm:"column:updated_at;type:int unsigned;not null;default:0"`
}

func (ExternalCustomerExtensionModel) TableName() string {
	return "external_customer_extension"
}
