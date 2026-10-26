package onlineWarehouseReq

import "gw-work-scrm/internal/app/http/request"

type BaseReq struct {
	request.EnterpriseBaseReq
	MemberId int64 `json:"member_id"`
}

type CreateReq struct {
	BaseReq
	OrderSn        string    `json:"order_sn" binding:"required"`
	ThirdId        int       `json:"third_id"`
	ExternalUserId string    `json:"external_user_id" binding:"required"`
	UnionId        string    `json:"union_id"`
	UserName       string    `json:"user_name" binding:"required"`
	SkuList        []SkuList `json:"sku_list"`
}

type SkuList struct {
	UpcCode    string `json:"upc_code" binding:"required"`
	SkuName    string `json:"name" binding:"required"`
	ImageUrl   string `json:"image_url" binding:"required"`
	SkuInfo    string `json:"sku_info" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required"`
	Dept       string `json:"dept"`
	StyleColor string `json:"style_color"`
	SkuNumber  string `json:"sku_number"`
}

type ListReq struct {
	BaseReq
	request.PageReq
	OrderSn        string  `json:"order_sn"`
	ReceiverName   string  `json:"receiver_name"`
	StoreCode      string  `json:"store_code"`
	StoreName      string  `json:"store_name"`
	Status         int     `json:"status"`
	StartTime      int     `json:"start_time"`
	EndTime        int     `json:"end_time"`
	TransferStatus int     `json:"transfer_status" comment:"转货单状态，0全部，1待上传，2已上传"`
	MemberIds      []int64 `json:"member_ids"`
	ThirdIds       []int   `json:"third_ids"`
	LeaveMemberIds []int64 `json:"leave_member_ids" comment:"离职的导购id"`
	IsAll          bool    `json:"is_all" comment:"店长身份是否查询全部"`
	IsH5           bool    `json:"is_h5" comment:"是否是H5端"`
}

type InfoReq struct {
	BaseReq
	OrderId int64  `json:"order_id" comment:"网仓订单表id" binding:"required"`
	OrderSn string `json:"order_sn" comment:"销售单号"`
	Origin  int    `json:"origin" comment:"来源，1 H5端，0管理后台"`
}

type DashboardReq struct {
	BaseReq
	OrderSn        string  `json:"order_sn"`
	ThirdIds       []int   `json:"third_ids"`
	MemberIds      []int64 `json:"member_ids"`
	StartTime      int     `json:"start_time"`
	EndTime        int     `json:"end_time"`
	LeaveMemberIds []int64 `json:"leave_member_ids" comment:"离职的导购id"`
	IsAll          bool    `json:"is_all" comment:"店长身份是否查询全部"`
	IsH5           bool    `json:"is_h5" comment:"是否是H5端"`
}

type DcListReq struct {
	BaseReq
	UpcCode string `json:"upc_code" binding:"required"`
}

type InventoryListReq struct {
	BaseReq
	UpcCode      string   `json:"upc_code" binding:"required"`
	StoreCodeArr []string `json:"store_code_arr"` //前端传参门店标识
}

type ExportMailReq struct {
	ListReq
	Email string `json:"email" binding:"required"`
}

type SaveReceiveInfoReq struct {
	BaseReq
	OrderId         int64  `json:"order_id" binding:"required"`
	OrderSn         string `json:"order_sn"`
	ReceiverName    string `json:"receiver_name" binding:"required"`
	ReceiverPhone   string `json:"receiver_phone" binding:"required"`
	ReceiverAddress string `json:"receiver_address" binding:"required"`
	Province        string `json:"province" binding:"required"`
	City            string `json:"city" binding:"required"`
	County          string `json:"county" binding:"required"`
}

type LeaveSaReq struct {
	BaseReq
	ThirdIds []int `json:"third_ids"`
	CheckIds []int `json:"check_ids"`
}

type OrderInfoReq struct {
	BaseReq
	Id     int64 `json:"id" comment:"网仓订单表id" binding:"required"`
	Origin int   `json:"origin" comment:"来源，1 H5端，0管理后台"`
}

type CheckOrderReq struct {
	BaseReq
	OrderSn string `json:"order_sn" comment:"销售单号"`
	Origin  int    `json:"origin" comment:"来源，1 H5端，0管理后台"`
}
