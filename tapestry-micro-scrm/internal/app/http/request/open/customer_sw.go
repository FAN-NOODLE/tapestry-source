package openReq

type UpdateCouponSwParams struct {
	UnionId              string `json:"union_id" comment:"union_id，plan_send_type=2时，unionid必传"`
	HasCoupon            int    `json:"has_coupon" comment:"待使用券， 1有 2无"`
	HasCouponBy30days    int    `json:"has_coupon_by_30days" comment:"30天内临期券， 1有 2无"`
	HasNewMemberCoupon   int    `json:"has_new_member_coupon" comment:"新人礼券， 1有 2无"`
	HasBirthdayCoupon    int    `json:"has_birthday_coupon" comment:"生日礼券， 1有 2无"`
	HasRebateCoupon      int    `json:"has_rebate_coupon" comment:"欢迎礼券， 1有 2无"`
	HasPromotionalCoupon int    `json:"has_promotional_coupon" comment:"活动礼券， 1有 2无"`
	UpdatedAt            int    `json:"updated_at" comment:"更新时间"`
}

type BatchUpdateCouponSwReq struct {
	EnterpriseId int64                  `json:"enterprise_id"`
	DataList     []UpdateCouponSwParams `json:"data_list"`
}
