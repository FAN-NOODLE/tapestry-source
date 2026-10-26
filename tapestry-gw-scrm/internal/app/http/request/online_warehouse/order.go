package onlineWarehouseReq

import "gw-scrm/internal/app/http/request"

type BaseReq struct {
	request.EnterpriseBaseReq
	AdminID int64 `json:"admin_id"`
}

type OrderReq struct {
	BaseReq
	request.PageReq
	OrderSn        string `json:"order_sn"`
	ReceiverName   string `json:"receiver_name"`
	StoreCode      string `json:"store_code"`
	StoreName      string `json:"store_name"`
	Status         int    `json:"status"`
	StartTime      int    `json:"start_time"`
	EndTime        int    `json:"end_time"`
	TransferStatus int    `json:"transfer_status"`
	ExportType     string `json:"export_type"`
}

type TransferFileInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type SkuList struct {
	UpcCode         string `json:"upc_code" binding:"required"`
	SkuName         string `json:"name" binding:"required"`
	ImageUrl        string `json:"image_url"`
	SkuInfo         string `json:"sku_info" binding:"required"`
	Quantity        int    `json:"quantity" comment:"商品总数"`
	RefusalQuantity int    `json:"refusal_quantity" comment:"退单数"`
}

type SubmitOrderReq struct {
	InfoReq
	RefusalInfo    []SkuList        `json:"refusal_info" comment:"退单信息"`
	TrackingNumber string           `json:"tracking_number" comment:"运单号"`
	TransferFile   TransferFileInfo `json:"transfer_file" comment:"转货单信息"`
}

type InfoReq struct {
	BaseReq
	OrderId int64  `json:"order_id" binding:"required"`
	OrderSn string `json:"order_sn" comment:"销售单号"`
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

type ListResp struct {
	OrderSn         string           `json:"order_sn"`
	ThirdId         string           `json:"third_id"`
	MemberId        int64            `json:"member_id"`
	MemberName      string           `json:"member_name"`
	StoreCode       string           `json:"store_code"`
	StoreName       string           `json:"store_name"`
	ReceiverName    string           `json:"receiver_name"`
	TrackingNumber  string           `json:"tracking_number" `
	ShippedAt       int              `json:"shipped_at"`
	ReceiverPhone   string           `json:"receiver_phone"`
	ReceiverAddress string           `json:"receiver_address"`
	TransferFile    TransferFileInfo `json:"transfer_file"`
	Status          int              `json:"status"`
	TransferStatus  int              `json:"transfer_status"`
	CreatedAt       int64            `json:"created_at"`
	UpdatedAt       int64            `json:"updated_at"`
	TotalQuantity   int              `json:"total_quantity"`
}

type ExportListsResp struct {
	Total int64
	List  []ExportListRes
}

type ExportListRes struct {
	ListResp
	SkuList
}

type ImportOrderReq struct {
	BaseReq
	FileUrl string `json:"file_url" binding:"required"`
}

type BatchTransferReq struct {
	BaseReq
	ThirdId      int              `json:"third_id" binding:"required"`
	ShippedTime  int              `json:"shipped_time" comment:"发货时间"  binding:"required"`
	TransferFile TransferFileInfo `json:"transfer_file" comment:"转货单信息" binding:"required"`
}

type TransferExportReq struct {
	BaseReq
	ThirdId int `json:"third_id"`
}

type TransferImportReq struct {
	BaseReq
	ThirdId      int              `json:"third_id"`
	FileUrl      string           `json:"file_url" binding:"required"`
	TransferFile TransferFileInfo `json:"transfer_file" binding:"required"`
}
