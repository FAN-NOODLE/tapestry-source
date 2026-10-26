package onlineWarehouseReq

import (
	"micro-scrm/internal/app/http/request"
	onlineWarehouseModel "micro-scrm/internal/app/models/micro_scrm/online_warehouse"
)

type BaseReq struct {
	AdminManager
	Mobile string `json:"mobile"`
}

type CreateReq struct {
	BaseReq
	OrderSn        string    `json:"order_sn"`
	ThirdId        int       `json:"third_id"`
	ExternalUserId string    `json:"external_user_id"`
	UnionId        string    `json:"union_id"`
	UserName       string    `json:"user_name"`
	SkuList        []SkuList `json:"sku_list"`
}

type SaveReceiveInfoReq struct {
	BaseReq
	OrderId         int64  `json:"order_id"`
	OrderSn         string `json:"order_sn"`
	ReceiverName    string `json:"receiver_name" binding:"required"`
	ReceiverPhone   string `json:"receiver_phone" binding:"required"`
	ReceiverAddress string `json:"receiver_address" binding:"required"`
	Province        string `json:"province" binding:"required"`
	City            string `json:"city" binding:"required"`
	County          string `json:"county" binding:"required"`
}

type SkuList struct {
	UpcCode    string `json:"upc_code"`
	SkuName    string `json:"name"`
	ImageUrl   string `json:"image_url"`
	SkuInfo    string `json:"sku_info" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required"`
	Dept       string `json:"dept"`
	StyleColor string `json:"style_color"`
	SkuNumber  string `json:"sku_number"`
}

type AdminManager struct {
	request.EnterpriseBaseReq
	MemberId int64 `json:"member_id"`
	AdminID  int64 `json:"admin_id"`
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
	OnlyPartRefund int     `json:"only_part_refund" comment:"仅展示部分发货的已发货的订单，1展示。0不展示"`
	IsAll          bool    `json:"is_all" comment:"店长身份是否查询全部"`
	IsH5           bool    `json:"is_h5" comment:"是否是H5端"`
	ExportType     string  `json:"export_type" comment:"导出类型"`
}

type SubmitOrderReq struct {
	InfoReq
	TrackingNumber string                                `json:"tracking_number" comment:"运单号"`
	RefusalInfo    []SkuList                             `json:"refusal_info" comment:"退单信息"`
	TransferFile   onlineWarehouseModel.TransferFileInfo `json:"transfer_file" comment:"转货单信息"`
}

type InfoReq struct {
	BaseReq
	OrderId int64  `json:"order_id" binding:"required"`
	OrderSn string `json:"order_sn" comment:"销售单号"`
	Origin  int    `json:"origin" comment:"h5端，值为1，管理后台端为0"`
}

type StoreListReq struct {
	BaseReq
	request.PageReq
	StoreCode   string `json:"store_code"`
	StoreName   string `json:"store_name"`
	ShippedTime int    `json:"shipped_time" comment:"发货时间"`
}

type StoreOrderListReq struct {
	BaseReq
	request.PageReq
	ThirdId     int `json:"third_id" binding:"required"`
	ShippedTime int `json:"shipped_time" comment:"发货时间"`
}

type ImportOrderReq struct {
	BaseReq
	FileUrl string `json:"file_url" binding:"required"`
}

type ImportOrderData struct {
	Id              int64  `json:"id"`
	OrderSn         string `json:"order_sn"`
	Sku             string `json:"sku"`
	UpcCode         string `json:"upc_code"`
	IsRefusal       int    `json:"is_refusal"`
	RefusalQuantity int    `json:"refusal_quantity"`
	TrackingNumber  string `json:"tracking_number"`
	TransferFileUrl string `json:"transfer_file_url"`
	Quantity        int    `json:"quantity"`
	ReceiverName    string `json:"receiver_name"`
	ReceiverPhone   string `json:"receiver_phone"`
	ReceiverAddress string `json:"receiver_address"`
	StoreCode       string `json:"store_code"`
	StoreName       string `json:"store_name"`
	OrderStatus     string `json:"order_status"`
	MemberName      string `json:"member_name"`
	CreateTime      string `json:"create_time"`
	UpdateTime      string `json:"update_time"`
	FailReason      string `json:"fail_reason"`
	Province        string `json:"province"`
	City            string `json:"city"`
	County          string `json:"county"`
}

type BatchTransferReq struct {
	BaseReq
	ThirdId      int                                   `json:"third_id" binding:"required"`
	ShippedTime  int                                   `json:"shipped_time" comment:"发货时间" binding:"required"`
	TransferFile onlineWarehouseModel.TransferFileInfo `json:"transfer_file" comment:"转货单信息" binding:"required"`
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

type ExportMailReq struct {
	ListReq
	Email string `json:"email" binding:"required"`
}

type LeaveSaReq struct {
	BaseReq
	ThirdIds []int `json:"third_ids"`
	CheckIds []int `json:"check_ids"`
}

type TransferExportReq struct {
	BaseReq
	ThirdId int `json:"third_id" binding:"required"`
}

type TransferImportReq struct {
	TransferExportReq
	FileUrl      string                                `json:"file_url" binding:"required"`
	TransferFile onlineWarehouseModel.TransferFileInfo `json:"transfer_file" binding:"required"`
}

type OrderInfoReq struct {
	BaseReq
	Id     int64 `json:"id"`
	Origin int   `json:"origin" comment:"h5端，值为1，管理后台端为0"`
}

type TransferImportData struct {
	OrderId    int64  `json:"order_id"`
	OrderSn    string `json:"order_sn"`
	FailReason string `json:"fail_reason"`
}

type CheckOrderSnReq struct {
	BaseReq
	OrderSn string `json:"order_sn" comment:"销售单号"`
}
