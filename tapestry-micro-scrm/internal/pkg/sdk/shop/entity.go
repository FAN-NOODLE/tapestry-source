package shop

// WxCouponListByIdsRequest 根据id查询微信优惠券信息请求参数
type WxCouponListByIdsRequest struct {
	CouponIds []int64 `json:"coupon_ids"`
}

// WxCouponListByIdsResponse 根据id查询微信优惠券信息返回数据格式
type WxCouponListByIdsResponse struct {
	ResponseBase
	Data []WxCouponInfo `json:"data"`
}

// WxCouponInfo 微信优惠券基本信息
type WxCouponInfo struct {
	ID             int64  `json:"id"`
	EnterpriseID   int64  `json:"enterprise_id"`
	CouponStockID  string `json:"coupon_stock_id"`
	CouponName     string `json:"coupon_name"`
	CouponValue    int    `json:"coupon_value"`
	CouponMininumn int    `json:"coupon_mininumn"`
}

// LiveListRequest live列表请求参数
type LiveListRequest struct {
	EnterpriseID int64   `json:"enterprise_id"`
	Ids          []int64 `json:"ids"`
	Page         int     `json:"page"`
	Size         int     `json:"size"`
}

// LiveListResponse live列表返回格式
type LiveListResponse struct {
	ResponseBase
	Data struct {
		List  []LiveInfo `json:"list"`
		Total int        `json:"total"`
		Page  int        `json:"page"`
		Size  int        `json:"size"`
	} `json:"data"`
}

// LiveInfo live基本信息
type LiveInfo struct {
	ID        int64  `json:"id"`
	LiveTitle string `json:"live_title"`
}

// ActivityPageListRequest live列表请求参数
type ActivityPageListRequest struct {
	EnterpriseID int64   `json:"enterprise_id"`
	Ids          []int64 `json:"ids"`
	Page         int     `json:"page"`
	Size         int     `json:"size"`
}

// ActivityPageListResponse live列表返回格式
type ActivityPageListResponse struct {
	ResponseBase
	Data struct {
		List  []ActivityPageInfo `json:"list"`
		Total int                `json:"total"`
		Page  int                `json:"page"`
		Size  int                `json:"size"`
	} `json:"data"`
}

// ActivityPageInfo live基本信息
type ActivityPageInfo struct {
	ID           int64  `json:"id"`
	TemplateName string `json:"template_name"`
}
