package coupon

import "micro-scrm/internal/pkg/common"

// WxCouponIsAllYes 微信优惠券是全场商品
var WxCouponIsAllYes = 1

// 活动状态
const (
	DraftStatus      = 0
	PendingStatus    = 1
	NormalStatus     = 2
	DownStatus       = 3
	FinishStatus     = 4
	StopStatus       = 5
	TimingUpStatus   = 6
	TimingDownStatus = 7
)

type DDCRequestBase struct {
	Method   string `json:"method"`
	Platform string `json:"platform"`
}

// DDCRequestBody DDC请求数据结构
type DDCRequestBody struct {
	DDCRequestBase
	Data interface{} `json:"data"`
}

// SyncCouponRequest 同步优惠券接口请求参数
type SyncCouponRequest struct {
	DDCRequestBase
	CouponInfo BasicInfoEntity `json:"data"`
}

// SyncCouponUserRequest 同步用户优惠券接口请求参数
type SyncCouponUserRequest struct {
	DDCRequestBase
	DataList DataCouponUserList `json:"data"`
}

type DataCouponUserList struct {
	CouponUserInfoList []BasicCouponUserInfoEntity `json:"list"`
}

// BasicInfoEntity 微信优惠券详情结构体
type BasicInfoEntity struct {
	ID                int64  `json:"id" db:"id"`
	EnterpriseId      int64  `json:"enterprise_id" db:"enterprise_id"`
	CouponId          int64  `json:"coupon_id" db:"coupon_id"`
	CouponStockId     string `json:"coupon_stock_id" db:"coupon_stock_id"`
	CouponName        string `json:"coupon_name" db:"coupon_name"`
	CouponValue       int64  `json:"coupon_value" db:"coupon_value"`
	CouponMininumn    int64  `json:"coupon_mininumn" db:"coupon_mininumn"`
	CouponStatus      int8   `json:"coupon_status" db:"coupon_status"`
	CouponStockStatus int8   `json:"coupon_stock_status" db:"coupon_stock_status"`
	CouponTotal       int64  `json:"coupon_total" db:"coupon_total"`
	UseDayLimit       int64  `json:"use_day_limit" db:"use_day_limit"`
	MaxQuota          int64  `json:"max_quota" db:"max_quota"`
	IsSendNum         int64  `json:"is_send_num" db:"is_send_num"`
	IsAll             int64  `json:"is_all" db:"is_all"`
	CouponBudget      int64  `json:"coupon_budget" db:"coupon_budget"`
	Scope             string `json:"scope" db:"scope"`
	FavorType         string `json:"favor_type" db:"favor_type"`
	StartTime         int64  `json:"start_time" db:"start_time"`
	EndTime           int64  `json:"end_time" db:"end_time"`
	IsOverlay         int64  `json:"is_overlay" db:"is_overlay"`
	OverlayCoupons    string `json:"overlay_coupons" db:"overlay_coupons"`
	Channel           int64  `json:"channel" db:"channel"`
	Range             string `json:"range" db:"range"`
	RangeHide         string `json:"range_hide" db:"range_hide"`
	TagGoods          string `json:"tag_goods" db:"tag_goods"`
	TagName           string `json:"tag_name" db:"tag_name"`
	TagRuleType       int64  `json:"tag_rule_type" db:"tag_rule_type"`
	CreatedAt         int64  `json:"created_at" db:"created_at"`
	UpdatedAt         int64  `json:"updated_at" db:"updated_at"`
	DeletedAt         int64  `json:"deleted_at" db:"deleted_at"`
}

// BasicCouponUserInfoEntity 用户微信优惠券详情结构体
type BasicCouponUserInfoEntity struct {
	ID            int64  `json:"id" db:"id"`
	EnterpriseId  int64  `json:"enterprise_id" db:"enterprise_id"`
	CouponId      int64  `json:"coupon_id" db:"coupon_id"`
	OpenId        string `json:"open_id" db:"open_id"`
	CouponUserId  int64  `json:"coupon_user_id" db:"coupon_user_id"`
	UnionId       string `json:"union_id" db:"union_id"`
	UserId        int64  `db:"user_id" `
	OrderNo       string `json:"order_no" db:"order_no"`
	CouponStockId string `json:"coupon_stock_id" db:"coupon_stock_id"`
	Effective     int64  `json:"effective" db:"effective"`
	StartTime     int64  `json:"start_time" db:"start_time"`
	EndTime       int64  `json:"end_time" db:"end_time"`
	Status        int8   `json:"status" db:"status"`
	SourceId      int64  `json:"source_id" db:"source_id"`
	Source        int64  `json:"source" db:"source"`
	ActionId      string `json:"action_id" db:"action_id"`
	ActionType    int64  `json:"action_type" db:"action_type"`
	UseTime       int64  `json:"use_time" db:"use_time"`
	CreatedAt     int64  `json:"created_at" db:"created_at"`
	UpdatedAt     int64  `json:"updated_at" db:"updated_at"`
	DeletedAt     int64  `json:"deleted_at" db:"deleted_at"`
	WxCouponID    string `json:"wx_coupon_id"`
}

// 总数返回值
type CountResponse struct {
	Count int64 `json:"count" db:"count"`
}

// UserCouponListRequest 用户优惠券列表查询请求
type UserCouponListRequest struct {
	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" validate:"required" required:"企业ID不能为空"`
	UserID       int64   `json:"user_id" db:"user_id"`
	UserIds      []int64 `json:"user_ids"`
	UserIdCnt    int
	UnionID      string `json:"union_id" db:"union_id"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
}

// UserCouponListRequest 用户优惠券列表查询请求
type WxCouponRequest struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" validate:"required" required:"企业ID不能为空"`
	MemberUserId string `json:"member_user_id" db:"member_user_id" desc:"员工ID"`
	SearchText   string `json:"search_text" db:"search_text"`
	Page         int64  `json:"page" db:"page"`
	Size         int64  `json:"size" db:"size"`
}

// MemberGroupPageListResponse 员工的客户列表响应
type WxCouponListResponse struct {
	common.PageType
	List []WxCoupon `json:"list" desc:"列表数据"`
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

// UserCouponListResponse 用户优惠券列表返回结果
type UserCouponListResponse struct {
	UserCouponID   int64   `json:"user_coupon_id" db:"user_coupon_id"`
	EnterpriseID   int64   `json:"enterprise_id" db:"enterprise_id"`
	UnionID        string  `json:"union_id" db:"union_id"`
	UserID         int64   `json:"user_id" db:"user_id"`
	Effective      int64   `json:"effective" db:"effective"`
	StartTime      int64   `json:"start_time" db:"start_time"`
	EndTime        int64   `json:"end_time" db:"end_time"`
	Status         string  `json:"status" db:"status"`
	CouponID       int64   `json:"coupon_id" db:"coupon_id"`
	CouponName     string  `json:"coupon_name" db:"coupon_name"`
	CouponValue    int64   `json:"coupon_value" db:"coupon_value"`
	CouponStockID  string  `json:"coupon_stock_id" db:"coupon_stock_id"`
	MaxQuota       int64   `json:"max_quota" db:"max_quota"`
	OrderNo        string  `json:"order_no" db:"order_no" desc:"子订单编号"`
	IsAll          int     `json:"is_all" db:"is_all"`
	Range          string  `json:"range" db:"range"`
	RangeGoods     []int64 `json:"range_goods"`
	UseTime        int64   `json:"use_time" db:"use_time" desc:"优惠券使用时间"`
	CouponMininumn int64   `json:"coupon_mininumn" db:"coupon_mininumn"`
}

type BatchCouponUserEntity struct {
	CouponUserIds []interface{}
	EnterpriseId  int64 `json:"enterprise_id" db:"enterprise_id"`
}

// WithGoodsRequest 优惠券部分商品
type WithGoodsRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" validate:"required" required:"企业ID不能为空"`
	CouponId     int64 `json:"coupon_id" db:"coupon_id"`
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

type CommonRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" validate:"required" required:"企业ID不能为空"`
	Page         int64 `json:"page" db:"page"`
	Size         int64 `json:"size" db:"size"`
}

type WxCouponAvailableRequest struct {
	CommonRequest
	CouponName string `json:"coupon_name" db:"coupon_name"`
	StockID    string `json:"stock_id" db:"coupon_stock_id"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
	NowTime    int64
}

// WxCouponAvailableListResponse 微信可用的优惠券列表响应体
type WxCouponAvailableListResponse struct {
	common.PageType
	List []*WxCouponAvailable `json:"list" desc:"列表数据"`
}

type WxCouponAvailable struct {
	WxCoupon
	IsUsed    int64  `json:"is_used" desc:"是否被使用"`
	Condition string `json:"condition" desc:"条件"`
}

type MarketingStrategiesCouponLog struct {
	Id                  int64  `json:"id" db:"id"`
	EnterpriseId        int64  `json:"enterprise_id" db:"enterprise_id"`
	ShopUserId          int64  `json:"shop_user_id" db:"shop_user_id"`
	Openid              string `json:"openid" db:"openid"`
	UnionId             string `json:"union_id" db:"union_id"`
	MarketingStrategyId int64  `json:"marketing_strategy_id" db:"marketing_strategy_id"`
	ScrmWxCouponId      int64  `json:"scrm_wx_coupon_id" db:"scrm_wx_coupon_id"`
	WxStockId           string `json:"wx_stock_id" db:"wx_stock_id"`
	WxCouponId          int64  `json:"wx_coupon_id" db:"wx_coupon_id"`
	Type                int    `json:"type" db:"type"`
	CreatedAt           int64  `json:"created_at" db:"created_at"`
	UpdatedAt           int64  `json:"updated_at" db:"updated_at"`
	DeletedAt           int64  `json:"deleted_at" db:"deleted_at"`
}

type StrategyCouponLogRes struct {
	OpenID    string `json:"openid" db:"openid"`
	WxStockId string `json:"wx_stock_id" db:"wx_stock_id"`
}

type StrategyCouponLogResp struct {
	ShopUserId          int64  `json:"shop_user_id" db:"shop_user_id"`
	Openid              string `json:"openid" db:"openid"`
	UnionId             string `json:"union_id" db:"union_id"`
	MarketingStrategyId int64  `json:"marketing_strategy_id" db:"marketing_strategy_id"`
	ScrmWxCouponId      int64  `json:"scrm_wx_coupon_id" db:"scrm_wx_coupon_id"`
	WxStockId           string `json:"wx_stock_id" db:"wx_stock_id"`
}
