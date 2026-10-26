package reserveReq

type OrderInfo struct {
	ReserveId         int64  `json:"reserve_id" binding:"required"`
	ReserveTime       int    `json:"reserve_time" binding:"required"`
	StoreId           int64  `json:"store_id" binding:"required"`
	StoreName         string `json:"store_name" binding:"required"`
	StoreAddress      string `json:"store_address" binding:"required"`
	ReserveType       int    `json:"reserve_type"`
	Status            int    `json:"status"`
	UserName          string `json:"user_name" binding:"required"`
	UserUnionid       string `json:"user_unionid" binding:"required"`
	UserExternalId    string `json:"user_external_id" binding:"required"`
	MemberId          int64  `json:"member_id" binding:"required"`
	MemberName        string `json:"member_name" binding:"required"`
	MemberUserId      string `json:"member_user_id" binding:"required"`
	ReserveRemark     string `json:"reserve_remark"`
	UserRemark        string `json:"user_remark"`
	IsOpen            int    `json:"is_open"`
	IsCommit          int    `json:"is_commit"`
	ReserveOrderSn    string `json:"reserve_order_sn"`
	ReserveOrderPrice int64  `json:"reserve_order_price"`
}