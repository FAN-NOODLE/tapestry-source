package openReq

const (
	HasCoupon = iota + 1
)

type CouponReq struct {
	EnterpriseId int64 `json:"enterprise_id"`
	UpdateCouponParams
}

type UpdateCouponParams struct {
	UnionId                  string `json:"union_id" comment:"union_id，plan_send_type=2时，unionid必传"`
	HasCoupon                int    `json:"has_coupon" comment:"待使用券， 1有 2无"`
	HasCouponBy30days        int    `json:"has_coupon_by_30days" comment:"30天内临期券， 1有 2无"`
	HasBirthdayCoupon        int    `json:"has_birthday_coupon" comment:"生日券， 1有 2无"`
	HasRenewMembershipCoupon int    `json:"has_renew_membership_coupon" comment:"续会券， 1有 2无"`
	HasWelcomeCoupon         int    `json:"has_welcome_coupon" comment:"欢迎券， 1有 2无"`
	HasUpgradeCoupon         int    `json:"has_upgrade_coupon" comment:"升级券， 1有 2无"`
	UpdatedAt                int    `json:"updated_at" comment:"更新时间"`
}

type BatchUpdateCouponReq struct {
	EnterpriseId int64                `json:"enterprise_id"`
	DataList     []UpdateCouponParams `json:"data_list"`
}

type UpdateActiveParams struct {
	UnionId      string `json:"union_id" comment:"union_id，plan_send_type=2时，unionid必传"`
	UserId       string `json:"user_id" comment:"员工编码"`
	ActiveTime   int64  `json:"active_time" comment:"最近一次活跃会话的时间（10位unix时间戳）"`
	InteractTime int64  `json:"interact_time" comment:"最近一次触达时间（10位unix时间戳）"`
	UpdatedAt    int64  `json:"updated_at" comment:"更新时间"`
}

type BatchUpdateActiveReq struct {
	EnterpriseId int64                `json:"enterprise_id"`
	DataList     []UpdateActiveParams `json:"data_list"`
}
