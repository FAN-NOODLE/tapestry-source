package onlineWearhouseModel

import "micro-scrm/internal/app/models"

type OnlineWarehouseOrderModel struct {
	OrderSn         string                         `json:"order_sn" gorm:"column:order_sn;type:varchar(64);not null;comment:销售单号"`
	ThirdId         int                            `json:"third_id" gorm:"column:third_id;type:int(1);not null;default:0;comment:企微部门id"`
	MemberID        int64                          `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:member_id_index;comment:创建人"`
	TotalQuantity   int                            `json:"total_quantity" gorm:"column:total_quantity;type:int(10);not null;comment:总商品数量"`
	Status          int32                          `json:"status" gorm:"column:status;type:tinyint(1);not null;default:0;comment:状态 1待客户填写 2待发货 3 已发货 4 已拒单 5 已作废"`
	AdminId         int64                          `json:"admin_id" gorm:"column:admin_id;type:bigint(20);not null;default:0;comment:操作人id"`
	TrackingNumber  string                         `json:"tracking_number" gorm:"column:tracking_number;type:varchar(64);not null;comment:运单号"`
	ShippedAt       int                            `json:"shipped_at" gorm:"column:shipped_at;type:int(10);not null;default:0;comment:发件时间"`
	ReceiverName    string                         `json:"receiver_name" gorm:"column:receiver_name;type:varchar(64);not null;comment:收件人"`
	ReceiverPhone   string                         `json:"receiver_phone" gorm:"column:receiver_phone;type:varchar(64);not null;comment:收件电话"`
	ReceiverAddress string                         `json:"receiver_address" gorm:"column:receiver_address;type:varchar(255);not null;comment:收件地址"`
	TransferFile    string                         `json:"transfer_file" gorm:"column:transfer_file;type:varchar(255);not null;comment:转货单文件json"`
	SendQuantity    int                            `json:"send_quantity" gorm:"column:send_quantity;type:int(10);not null;comment:发货数量"`
	ExternalUserId  string                         `json:"external_user_id" gorm:"column:external_user_id;type:varchar(255);not null;comment:外部联系人id"`
	UnionId         string                         `json:"unionid" gorm:"column:unionid;type:varchar(255);not null;comment:unionid"`
	UserName        string                         `json:"user_name" gorm:"column:user_name;type:varchar(255);not null;comment:下单用户昵称"`
	IsExport        int                            `json:"is_export" gorm:"column:is_export;type:tinyint(1);not null;default:0;comment:是否导出，导出的订单，待发货的不能更新为已废弃"`
	SharePosterUrl  string                         `json:"share_poster_url" gorm:"column:share_poster_url;type:varchar(255);not null;comment:海报二维码url"`
	IsPoster        int                            `json:"is_poster" gorm:"column:is_poster;type:tinyint(1);not null;default:0;comment:是否生成二维码"`
	OperateTime     int64                          `json:"operate_time" gorm:"column:operate_time;type:int(10);not null;default:0;comment:操作时间"`
	Province        string                         `json:"province" gorm:"column:province;type:varchar(64);not null;comment:省"`
	City            string                         `json:"city" gorm:"column:city;type:varchar(64);not null;comment:市"`
	County          string                         `json:"county" gorm:"column:county;type:varchar(64);not null;comment:区"`
	SkuLists        []OnlineWarehouseOrderSkuModel `json:"sku_lists" gorm:"foreignKey:OrderId;references:ID"`
	models.BaseModel
}

// TableName returns the table name of the OnlineWarehouseOrder model
func (OnlineWarehouseOrderModel) TableName() string {
	return "online_warehouse_order"
}

const (
	OrderStatusWait    = iota + 1 // 待用户填写
	OrderStatusPending            // 待发货
	OrderStatusShipped            // 已发货
	OrderStatusRefusal            // 已拒单
	OrderStatusCancel             // 已作废
)

const (
	TransferNO = iota + 1
	TransferYES
)

const (
	OrderExportYes = iota + 1 //待发货订单已导出
)

type TransferFileInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Url  string `json:"url"`
}
