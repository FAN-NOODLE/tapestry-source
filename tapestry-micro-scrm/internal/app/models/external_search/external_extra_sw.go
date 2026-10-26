package externalSearchModel

import (
	"micro-scrm/internal/app/models"
)

type ExternalExtraSwModel struct {
	models.BaseModel
	UserUnionid             string   `json:"user_unionid" gorm:"column:user_unionid;type:varchar(255);not null;index:unionid_index;comment:预约服务的用户unionid"`
	ValidReserveAt          int      `json:"valid_reserve_at" gorm:"column:valid_reserve_at;type:int(10);not null;default:0;comment:最新的进行中的邀约的时间"`
	SuccessReserveAt        int      `json:"success_reserve_at" gorm:"column:success_reserve_at;type:int(10);not null;default:0;comment:最新的邀约成功时间"`
	TaskFinishAt            int      `json:"task_finish_at" gorm:"column:task_finish_at;type:int(10);not null;default:0;comment:最新的任务触达时间"`
	TaskCreatedAt           int64    `json:"task_created_at" gorm:"column:task_created_at;type:int(10);not null;default:0;comment:任务加入时间"`
	InteractTime            int      `json:"interact_time" gorm:"column:interact_time;type:int(10);not null;default:0;comment:最近互动时间"`
	CouponUpdatedAt         int      `json:"coupon_updated_at" gorm:"column:coupon_updated_at;type:int(10);not null;default:0;comment:优惠券同步时间"`
	CouponStatus            int      `json:"coupon_status" gorm:"column:coupon_status;type:int(10);not null;default:0;comment:未使用优惠券状态 0 无 1 有"`
	CouponBy30daysStatus    int      `json:"coupon_by_30days_status" gorm:"column:coupon_by_30days_status;type:int(10);not null;default:0;comment:30天临期券状态 0 无 1 有"`
	NewMemberCouponStatus   int      `json:"new_member_coupon_status" gorm:"column:new_member_coupon_status;type:int(10);not null;default:0;comment:新人礼券状态 0 无 1 有"`
	BirthdayCouponStatus    int      `json:"birthday_coupon_status" gorm:"column:birthday_coupon_status;type:int(10);not null;default:0;comment:生日礼券状态 0 无 1 有"`
	RebateCouponStatus      int      `json:"rebate_coupon_status" gorm:"column:rebate_coupon_status;type:int(10);not null;default:0;comment:欢迎礼券状态 0 无 1 有"`
	PromotionalCouponStatus int      `json:"promotional_coupon_status" gorm:"column:promotional_coupon_status;type:int(10);not null;default:0;comment:活动礼券状态 0 无 1 有"`
	PriceTag                string   `json:"price_tag" gorm:"column:price_tag;type:varchar(255);not null"`
	AgeTag                  string   `json:"age_tag" gorm:"column:age_tag;type:varchar(255);not null"`
	R12PurchaseAmount       int      `json:"r12_purchase_amount" gorm:"column:r12_purchase_amount;type:int(10);not null;default:0"`
	CategoryPreference      []string `json:"category_preference" gorm:"column:category_preference;type:Array(String);"`
	PurchaseDate            []int64  `json:"purchase_date" gorm:"column:purchase_date;type:Array(Int64);"`
}

// TableName returns the table name of the WorkExternalUserRelation model
func (ExternalExtraSwModel) TableName() string {
	return "external_extra"
}
