package externalSearchModel

import (
	"micro-scrm/internal/app/models"
)

const (
	NoCoupon = iota
	HasCoupon
)

const (
	WeChat = iota + 1
)

type ExternalExtraModel struct {
	models.BaseModel
	UserUnionid                 string   `json:"user_unionid" gorm:"column:user_unionid;type:varchar(255);not null;index:unionid_index;comment:预约服务的用户unionid"`
	ValidReserveAt              int      `json:"valid_reserve_at" gorm:"column:valid_reserve_at;type:int(10);not null;default:0;comment:最新的进行中的邀约的时间"`
	SuccessReserveAt            int      `json:"success_reserve_at" gorm:"column:success_reserve_at;type:int(10);not null;default:0;comment:最新的邀约成功时间"`
	TaskFinishAt                int      `json:"task_finish_at" gorm:"column:task_finish_at;type:int(10);not null;default:0;comment:最新的任务触达时间"`
	TaskCreatedAt               int64    `json:"task_created_at" gorm:"column:task_created_at;type:int(10);not null;default:0;comment:任务加入时间"`
	CouponUpdatedAt             int      `json:"coupon_updated_at" gorm:"column:coupon_updated_at;type:int(10);not null;default:0;comment:优惠券同步时间"`
	CouponStatus                int      `json:"coupon_status" gorm:"column:coupon_status;type:int(10);not null;default:0;comment:未使用优惠券状态 0 无 1 有"`
	CouponBy30daysStatus        int      `json:"coupon_by_30days_status" gorm:"column:coupon_by_30days_status;type:int(10);not null;default:0;comment:30天临期券状态 0 无 1 有"`
	BirthdayCouponStatus        int      `json:"birthday_coupon_status" gorm:"column:birthday_coupon_status;type:int(10);not null;default:0;comment:生日券状态 0 无 1 有"`
	RenewMembershipCouponStatus int      `json:"renew_membership_coupon_status" gorm:"column:renew_membership_coupon_status;type:int(10);not null;default:0;comment:续会券状态 0 无 1 有"`
	WelcomeCouponStatus         int      `json:"welcome_coupon_status" gorm:"column:welcome_coupon_status;type:int(10);not null;default:0;comment:欢迎券状态 0 无 1 有"`
	UpgradeCouponStatus         int      `json:"upgrade_coupon_status" gorm:"column:upgrade_coupon_status;type:int(10);not null;default:0;comment:升级券状态 0 无 1 有"`
	Openid                      string   `json:"openid" gorm:"column:openid;type:varchar(255);"`
	Nick                        string   `json:"nick" gorm:"column:nick;type:varchar(255);"`
	Name                        string   `json:"name" gorm:"column:name;type:varchar(255);"`
	Position                    string   `json:"position" gorm:"column:position;type:varchar(255);"`
	Avatar                      string   `json:"avatar" gorm:"column:avatar;type:varchar(255);"`
	Type                        int      `json:"type" gorm:"column:type;type:int(10);"`
	Sex                         int      `json:"sex" gorm:"column:sex;type:int(10);"`
	Birthday                    string   `json:"birthday" gorm:"column:birthday;type:varchar(255);"`
	UnionCode                   string   `json:"union_code" gorm:"column:union_code;type:varchar(255);"`
	Point                       string   `json:"point" gorm:"column:point;type:varchar(255);"`
	Grade                       string   `json:"grade" gorm:"column:grade;type:varchar(255);"`
	BaCode                      string   `json:"ba_code" gorm:"column:ba_code;type:varchar(255);"`
	BirthYear                   string   `json:"birth_year" gorm:"column:birth_year;type:varchar(255);"`
	BirthMonth                  string   `json:"birth_month" gorm:"column:birth_month;type:varchar(255);"`
	BirthDay                    string   `json:"birth_day" gorm:"column:birth_day;type:varchar(255);"`
	Mobile                      string   `json:"mobile" gorm:"column:mobile;type:varchar(255);"`
	CardLevel                   string   `json:"card_level" gorm:"column:card_level;type:varchar(255);"`
	CltType                     string   `json:"clt_type" gorm:"column:clt_type;type:varchar(255);"`
	AmountTotal                 int      `json:"amount_total" gorm:"column:amount_total;type:int(10);"`
	HistoryAmountTotal          int      `json:"history_amount_total" gorm:"column:history_amount_total;type:int(10);"`
	CustomerCode                string   `json:"customer_code" gorm:"column:customer_code;type:varchar(255);"`
	CustomerId                  string   `json:"customer_id" gorm:"column:customer_id;type:varchar(255);"`
	LastPurchaseDate            int64    `json:"last_purchase_date" gorm:"column:last_purchase_date;type:int(10);"`
	LastPurchaseStoreCode       string   `json:"last_purchase_store_code" gorm:"column:last_purchase_store_code;type:varchar(255);"`
	UpgradeNeedAmount           int      `json:"upgrade_need_amount" gorm:"column:upgrade_need_amount;type:int(10);"`
	CustomerType                int      `json:"customer_type" gorm:"column:customer_type;type:int(10);"`
	CustomerTypeUpdatedAt       int      `json:"customer_type_updated_at" gorm:"column:customer_type_updated_at;type:int(10);"`
	VipType                     int      `json:"vip_type" gorm:"column:vip_type;type:int(10);"`
	CustomerDeletedAt           int      `json:"customer_deleted_at" gorm:"column:customer_deleted_at;type:int(10);"`
	AgeTag                      string   `json:"age_tag" gorm:"column:age_tag;type:varchar(255);"`
	CategoryPreference          []string `json:"category_preference" gorm:"column:category_preference;type:Array(String);"`
	PurchasedDate               []int64  `json:"purchased_date" gorm:"column:purchased_date;type:Array(Int64);"`
	ExternalImpressionTagIds    []int64  `json:"external_impression_tag_ids" gorm:"column:external_impression_tag_ids;type:Array(Int64);"`
	HighWilling                 string   `json:"high_willing" gorm:"column:high_willing;type:varchar(255);"`
	PriceSensitivity            string   `json:"price_sensitivity" gorm:"column:price_sensitivity;type:varchar(255);"`
}

// TableName returns the table name of the WorkExternalUserRelation model
func (ExternalExtraModel) TableName() string {
	return "external_extra"
}

type MaterializeMainId struct {
	RelationId         int64  `json:"relation_id"`
	WorkExternalUserId int64  `json:"work_external_user_id"`
	UnionId            string `json:"union_id"`
}

type MaterializeMainData struct {
	RelationIds         []int64  `json:"relation_ids"`
	WorkExternalUserIds []int64  `json:"work_external_user_ids"`
	UnionIds            []string `json:"union_ids"`
}
