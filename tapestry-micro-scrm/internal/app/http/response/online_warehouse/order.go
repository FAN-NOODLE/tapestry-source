package onlineWarehouseResp

import (
	onlineWarehouseModel "micro-scrm/internal/app/models/micro_scrm/online_warehouse"
)

type ListResp struct {
	Id              int64                                 `json:"id"`
	OrderSn         string                                `json:"order_sn"`
	ThirdId         int                                   `json:"third_id"`
	MemberId        int64                                 `json:"member_id"`
	MemberName      string                                `json:"member_name"`
	StoreCode       string                                `json:"store_code"`
	StoreName       string                                `json:"store_name"`
	ReceiverName    string                                `json:"receiver_name"`
	TrackingNumber  string                                `json:"tracking_number" `
	ShippedAt       int                                   `json:"shipped_at"`
	ReceiverPhone   string                                `json:"receiver_phone"`
	ReceiverAddress string                                `json:"receiver_address"`
	TransferFile    onlineWarehouseModel.TransferFileInfo `json:"transfer_file"`
	Status          int                                   `json:"status"`
	TransferStatus  int                                   `json:"transfer_status"`
	CreatedAt       int64                                 `json:"created_at"`
	UpdatedAt       int64                                 `json:"updated_at"`
	TotalQuantity   int                                   `json:"total_quantity"`
	SendQuantity    int                                   `json:"send_quantity"`
	UserName        string                                `json:"user_name"`
	ExternalUserId  string                                `json:"external_user_id"`
	Province        string                                `json:"province"`
	City            string                                `json:"city"`
	County          string                                `json:"county"`
}

type SkuList struct {
	UpcCode  string `json:"upc_code" binding:"required"`
	SkuName  string `json:"name" binding:"required"`
	SkuInfo  string `json:"sku_info" binding:"required"`
	ImageUrl string `json:"image_url" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
	Status   int    `json:"status"`
}

type InfoResp struct {
	Id              int64                                 `json:"id"`
	OrderSn         string                                `json:"order_sn"`
	MemberId        int64                                 `json:"member_id"`
	MemberName      string                                `json:"member_name"`
	StoreCode       string                                `json:"store_code"`
	StoreName       string                                `json:"store_name"`
	ReceiverName    string                                `json:"receiver_name"`
	ReceiverPhone   string                                `json:"receiver_phone"`
	ReceiverAddress string                                `json:"receiver_address"`
	ShippedAt       int                                   `json:"shipped_at"`
	CreatedAt       int64                                 `json:"created_at"`
	UpdatedAt       int64                                 `json:"updated_at"`
	TrackingNumber  string                                `json:"tracking_number"`
	TransferFile    onlineWarehouseModel.TransferFileInfo `json:"transfer_file"`
	SkuList         []InfoSku                             `json:"sku_list"`
	AllNum          int                                   `json:"all_num"`
	RefusalNum      int                                   `json:"refusal_num"`
	Status          int32                                 `json:"status"`
	UserName        string                                `json:"user_name"`
	ExternalUserId  string                                `json:"external_user_id"`
	SharePosterUrl  string                                `json:"share_poster_url"`
	OperateTime     int64                                 `json:"operate_time"`
}

type InfoSku struct {
	SkuList
	RefusalQuantity int `json:"refusal_quantity"`
}

type StoreListResp struct {
	StoreCode      string `json:"store_code"`
	StoreName      string `json:"store_name"`
	ThirdId        int    `json:"third_id"`
	OrderNum       int    `json:"order_num"`
	ShippedTime    int    `json:"shipped_time"`
	TransferStatus int    `json:"transfer_status"`
}

type StoreOrderListResp struct {
	OrderSn       string `json:"order_sn"`
	MemberId      int64  `json:"member_id"`
	MemberName    string `json:"member_name"`
	ShippedAt     int    `json:"shipped_at"`
	CreatedAt     int64  `json:"created_at"`
	TotalQuantity int    `json:"total_quantity"`
}

type SkuExportList struct {
	UpcCode         string `json:"upc_code"`
	SkuName         string `json:"name"`
	SkuInfo         string `json:"sku_info"`
	ImageUrl        string `json:"image_url"`
	Quantity        int    `json:"quantity" comment:"商品总数"`
	RefusalQuantity int    `json:"refusal_quantity" comment:"退单数"`
	SkuStatus       int    `json:"sku_status"`
	Dept            string `json:"dept"`
	StyleColor      string `json:"style_color"`
	NewFormatSku    string `json:"new_format_sku"`
}
type ExportListResp struct {
	ListResp
	SkuExportList
}

type ImportOrderResp struct {
	UploadNum    int    `json:"upload_num"`
	Success      int    `json:"success"`
	DuplicateNum int    `json:"duplicate_num"`
	FailedNum    int    `json:"failed_num"`
	FailUrl      string `json:"fail_url"`
}

type TipsNums struct {
	Status int `json:"status"`
	Num    int `json:"num"`
}
type TipsStatusNums struct {
	AllNum     int `json:"all_num"`
	WaitNum    int `json:"wait_num"`
	PendingNum int `json:"pending_num"`
	SendNum    int `json:"send_num"`
	RefuseNum  int `json:"refuse_num"`
	CancelNum  int `json:"cancel_num"`
}

type StoreRankResp struct {
	ThirdId   int    `json:"third_id"`
	SkuNum    int    `json:"sku_num"`
	Num       int    `json:"num"`
	StoreName string `json:"store_name"`
	StoreCode string `json:"store_code"`
}

type SaRankResp struct {
	StoreRankResp
	MemberId   int64  `json:"member_id"`
	MemberName string `json:"member_name"`
	Avatar     string `json:"avatar"`
	SkuNum     int    `json:"sku_num"`
	Num        int    `json:"num"`
	SendNum    int    `json:"send_num"`
}

type MemberNumResp struct {
	MemberId int64  `json:"member_id"`
	Num      int64  `json:"num"`
	UserId   string `json:"user_id"`
	UpdateAt int64  `json:"update_at"`
}

type CustomerOrderInfoResp struct {
	OrderStatus int      `json:"order_status" comment:"订单状态;0正常，1成功，2已作废"`
	OrderSn     string   `json:"order_sn" comment:"订单号"`
	OrderInfo   InfoResp `json:"order_info"`
}

type TransferImportResp struct {
	UploadNum    int    `json:"upload_num"`
	Success      int    `json:"success"`
	DuplicateNum int    `json:"duplicate_num"`
	FailedNum    int    `json:"failed_num"`
	FailUrl      string `json:"fail_url"`
}

type DashboardResp struct {
	Order TipsStatusNums `json:"order"`
	Goods TipsStatusNums `json:"goods"`
}

type GoodsRankResp struct {
	SkuName  string `json:"sku_name"`
	SkuInfo  string `json:"sku_info"`
	Dept     string `json:"dept"`
	ImageUrl string `json:"image_url"`
	Num      int    `json:"num"`
}
