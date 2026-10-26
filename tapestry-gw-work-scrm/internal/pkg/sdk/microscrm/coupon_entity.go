package microScrm

// UserCouponListRequest 用户优惠券列表查询请求
type WxCouponRequest struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	MemberUserId string `json:"member_user_id" db:"member_user_id" desc:"企业员工Id"`
	SearchText   string `json:"search_text" db:"search_text"`
	Page         int64  `json:"page" db:"page"`
	Size         int64  `json:"size" db:"size"`
}

// MemberGroupPageListResponse 员工的客户列表响应
type WxCouponListResponse struct {
	Page  int64      `json:"page"`
	Size  int64      `json:"size"`
	Total int64      `json:"total"`
	List  []WxCoupon `json:"list" desc:"列表数据"`
}

// WxCouponResponse 微信优惠券返回结构体
type WxCoupon struct {
	ID                int64  `json:"id" db:"id"`
	EnterpriseID      int64  `json:"enterprise_id" db:"enterprise_id"`
	CouponID          int64  `json:"coupon_id" db:"coupon_id"`
	CouponStockID     string `json:"coupon_stock_id" db:"coupon_stock_id"`
	CouponName        string `json:"coupon_name" db:"coupon_name"`
	CouponValue       int64  `json:"coupon_value" db:"coupon_value"`
	CouponMininumn    int64  `json:"coupon_mininumn" db:"coupon_mininumn"`
	CouponStatus      int64  `json:"coupon_status" db:"coupon_status"`
	CouponStockStatus int64  `json:"coupon_stock_status" db:"coupon_stock_status"`
	CouponTotal       int64  `json:"coupon_total" db:"coupon_total"`
	UseDayLimit       int64  `json:"use_day_limit" db:"use_day_limit"`
	MaxQuota          int64  `json:"max_quota" db:"max_quota"`
	IsSendNum         int64  `json:"is_send_num" db:"is_send_num"`
	CouponBudget      int64  `json:"coupon_budget" db:"coupon_budget"`
	Scope             string `json:"scope" db:"scope"`
	FavorType         string `json:"favor_type" db:"favor_type"`
	StartTime         int64  `json:"start_time" db:"start_time"`
	EndTime           int64  `json:"end_time" db:"end_time"`
	IsOverlay         int64  `json:"is_overlay" db:"is_overlay"`
	OverlayCoupons    string `json:"overlay_coupons" db:"overlay_coupons"`
	Channel           int64  `json:"channel" db:"channel"`
	IsAll             int64  `json:"is_all" db:"is_all"`
	Range             string `json:"range" db:"range"`
	RangeHide         string `json:"range_hide" db:"range_hide"`
	TagGoods          string `json:"tag_goods" db:"tag_goods"`
	TagName           string `json:"tag_name" db:"tag_name"`
	TagRuleType       int64  `json:"tag_rule_type" db:"tag_rule_type"`
	CreatedAt         int64  `json:"created_at" db:"created_at"`
	UpdatedAt         int64  `json:"updated_at" db:"updated_at"`
	DeletedAt         int64  `json:"deleted_at" db:"deleted_at"`
	ShareAppId        string `json:"share_app_id" db:"share_app_id"`
	SharePath         string `json:"share_path" db:"share_path"`
	DefaultPicture    string `json:"default_picture" db:"default_picture"`
}

// WithGoodsRequest 优惠券部分商品
type WithGoodsRequest struct {
	EnterpriseID int64 `json:"enterprise_id"`
	CouponId     int64 `json:"coupon_id"`
	Page         int   `json:"page"`
	Size         int   `json:"size"`
}

// WithGoodsRequest 优惠券部分商品
type WithGoodsResponse struct {
	List  []WithGoodsSpuItem `json:"list"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
	Total int64              `json:"total"`
}

type WithGoodsSpuItem struct {
	SpuID   int64  `json:"spu_id"`
	SpuName string `json:"spu_name"`
	Image   string `json:"image"`
}
