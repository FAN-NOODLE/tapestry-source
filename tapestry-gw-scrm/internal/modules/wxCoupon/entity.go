package wxcoupon

type CommonRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" validate:"required" required:"企业ID不能为空"`
	Page         int64 `json:"page" db:"page"`
	Size         int64 `json:"size" db:"size"`
}

type AvailableWxCouponRequest struct {
	CommonRequest
	CouponName string `json:"coupon_name" db:"coupon_name"`
	StockID    string `json:"stock_id" db:"coupon_stock_id"`
	NowTime    int64
}
