package reserve

type EnterpriseBase struct {
	EnterpriseId int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
}

type OrderInfo struct {
	ReserveId         int64  `json:"reserve_id" binding:"required"`
	ReserveTime       int    `json:"reserve_time" binding:"required"`
	StoreId           int64  `json:"store_id" binding:"required"`
	StoreName         string `json:"store_name" binding:"required"`
	StoreAddress      string `json:"store_address"`
	ReserveType       int    `json:"reserve_type"`
	Status            int    `json:"status"`
	UserName          string `json:"user_name" binding:"required"`
	UserUnionid       string `json:"user_unionid" binding:"required"`
	UserExternalId    string `json:"user_external_id" binding:"required"`
	Mobile            string `json:"mobile"`
	ReserveRemark     string `json:"reserve_remark"`
	UserRemark        string `json:"user_remark"`
	IsOpen            int    `json:"is_open"`
	IsCommit          int    `json:"is_commit"`
	ReserveOrderSn    string `json:"reserve_order_sn"`
	ReserveOrderPrice int64  `json:"reserve_order_price"`
	ReserveEndTime    int    `json:"reserve_end_time"`
}

type OrderInfoReq struct {
	EnterpriseBase
	OrderInfo
}

type OrderListReq struct {
	EnterpriseBase
	Unionid   string  `json:"unionid"`
	Page      int     `json:"page"`
	Size      int     `json:"size"`
	Mobile    string  `json:"mobile"`
	MemberIds []int64 `json:"member_ids"`
	OrderSort int     `json:"order_sort"`
	Status    int     `json:"status"`
	StoreIds  []int   `json:"store_ids"`
	IsSM      int     `json:"is_sm"`
	MemberId  int64   `json:"member_id"`
}

type OrderTimeListReq struct {
	EnterpriseBase
	MemberId int64  `json:"member_id"`
	Mobile   string `json:"mobile"`
}

type ReserveServiceListByMemberReq struct {
	EnterpriseBase
	Mobile string `json:"mobile"`
	Status int    `json:"status"`
}

type OrderBase struct {
	ID            int64  `json:"id" binding:"required" comment:"ID"`
	ReserveNumber string `json:"reserve_number" binding:"required" comment:"预约单编号"`
}

type OrderBaseReq struct {
	OrderBase
	EnterpriseBase
}

type MemberDetailReq struct {
	EnterpriseBase
	Mobile string `json:"mobile" comment:"手机号"`
}

type UserOrderBase struct {
	OrderBase
	EnterpriseBase
	UnionId string `json:"union_id"`
}

type ChangeStatus struct {
	UserOrderBase
	IsStore        int    `json:"is_store"`
	CustomerStatus int    `json:"customer_status"`
	Mobile         string `json:"mobile"`
}

type CheckReserveTime struct {
	EnterpriseBase
	ReserveTime int64 `json:"reserve_time"`
	MemberId    int64 `json:"member_id"`
	StoreId     int64 `json:"store_id" comment:"门店id"`
}

type OrderAssignReq struct {
	OrderBaseReq
	Mobile string `json:"mobile" binding:"required" required:"导购手机号"`
}

type CustomerReserveListsReq struct {
	EnterpriseBase
	ExternalUserID string `json:"external_userid" validate:"required" required:"无效的外部联系人"`
	Page           int    `json:"page"`
	Size           int    `json:"size"`
}

type OrderShowReq struct {
	OrderBaseReq
	MemberId int64 `json:"member_id"`
}
