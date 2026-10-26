package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type TaskConditionDataOrderGoodsModel struct {
	models.BaseModel
	OrderNo    string `json:"order_no" gorm:"column:order_no;type:varchar(255);not null;index:order_no_index;comment:关联订单id"`
	GoodsName  string `json:"goods_name" gorm:"column:goods_name;type:varchar(255);not null;comment:商品名称"`
	Style      string `json:"style" gorm:"column:style;type:varchar(255);not null;comment:款号"`
	Color      string `json:"color" gorm:"column:color;type:varchar(255);not null;comment:颜色"`
	GoodsPrice int    `json:"goods_price" gorm:"column:goods_price;type:int(10);not null;default:0;comment:商品价格"`
	PayTime    int    `json:"pay_time" gorm:"column:pay_time;type:int(10);not null;default:0;comment:交易时间"`
	GoodsURL   string `json:"goods_url" gorm:"column:goods_url;type:varchar(255);not null;comment:商品地址"`
}

// TableName returns the table name of the TaskConditionDataOrderGoods model
func (TaskConditionDataOrderGoodsModel) TableName() string {
	return "task_condition_data_order_goods"
}
