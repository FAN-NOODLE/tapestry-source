package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type ReserveOrderGoodsModel struct {
	models.BaseModel
	RelationOrderID string                 `json:"relation_order_id" gorm:"column:relation_order_id;type:varchar(255);not null;index:relation_order_index;comment:关联订单id"`
	GoodsName       string                 `json:"goods_name" gorm:"column:goods_name;type:varchar(255);not null;comment:商品名称"`
	Style           string                 `json:"style" gorm:"column:style;type:varchar(255);not null;comment:款号"`
	Color           string                 `json:"color" gorm:"column:color;type:varchar(255);not null;comment:颜色"`
	GoodsPrice      int                    `json:"goods_price" gorm:"column:goods_price;type:int(10);not null;default:0;comment:商品价格"`
	PayTime         int                    `json:"pay_time" gorm:"column:pay_time;type:int(10);not null;default:0;comment:交易时间"`
	GoodsUrl        string                 `json:"goods_url" gorm:"column:goods_url;type:varchar(255);not null;comment:商品地址"`
	ReserveSales    ReserveOrderSalesModel `json:"-" gorm:"foreignKey:RelationOrderID,UpdatedAt;references:RelationOrderID,UpdatedAt"`
}

// TableName returns the table name of the ReserveOrderGoods model
func (ReserveOrderGoodsModel) TableName() string {
	return "reserve_order_goods"
}
