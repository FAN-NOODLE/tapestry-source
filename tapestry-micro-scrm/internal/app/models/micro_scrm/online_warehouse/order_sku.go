package onlineWearhouseModel

import "micro-scrm/internal/app/models"

type OnlineWarehouseOrderSkuModel struct {
	OrderId       int64  `json:"order_id" gorm:"column:order_id;type:bigint(20);not null;comment:销售单号表主键"`
	OrderSn       string `json:"order_sn" gorm:"column:order_sn;type:varchar(64);not null;index:order_sn_index;comment:销售单号表主键"`
	UpcCode       string `json:"upc_code" gorm:"column:upc_code;type:varchar(64);not null;comment:商品upc"`
	SkuName       string `json:"sku_name" gorm:"column:sku_name;type:varchar(64);not null;comment:商品名称"`
	SkuInfo       string `json:"sku_info" gorm:"column:sku_info;type:varchar(64);not null;comment:商品skuItem 颜色 型号"`
	ImageUrl      string `json:"image_url" gorm:"column:image_url;type:varchar(255);not null;comment:商品图片地址"`
	Quantity      int    `json:"quantity" gorm:"column:quantity;type:int(10);not null;default:0;comment:商品数量"`
	Status        int32  `json:"status" gorm:"column:status;type:tinyint(1);not null;default:0;comment:状态 0待发货 1 已发货 2 已拒单"`
	RefusalReason string `json:"refusal_reason" gorm:"column:refusal_reason;type:varchar(255);not null;comment:拒单原因"`
	Dept          string `json:"dept" gorm:"column:dept;type:varchar(32);not null;comment:ross返回的dept"`
	StyleColor    string `json:"style_color" gorm:"column:style_color;type:varchar(32);not null;comment:ross返回的  StyleCode+color"`
	SkuNumber     string `json:"sku_number" gorm:"column:sku_number;type:varchar(64);not null;comment:ross返回的 skuNumber"`
	models.BaseModel
}

// TableName returns the table name of the OnlineWarehouseOrderSku model
func (OnlineWarehouseOrderSkuModel) TableName() string {
	return "online_warehouse_order_sku"
}
