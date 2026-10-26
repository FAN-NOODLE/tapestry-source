package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type ReserveOrderSalesModel struct {
	models.BaseModel
	ReserveNumber     string                   `json:"reserve_number" gorm:"column:reserve_number;type:varchar(255);not null;index:reserve_number_index;comment:邀约单号"`
	MemberUserID      string                   `json:"member_user_id" gorm:"column:member_user_id;type:varchar(255);not null;index:member_user_id_index;comment:导购的user_id"`
	UserUnionid       string                   `json:"user_unionid" gorm:"column:user_unionid;type:varchar(255);not null;index:unionid_index;comment:用户的unionid"`
	StoreCode         string                   `json:"store_code" gorm:"column:store_code;type:varchar(255);not null;comment:交易门店code"`
	RelationOrderID   string                   `json:"relation_order_id" gorm:"column:relation_order_id;type:varchar(255);not null;comment:关联的订单id"`
	OriginOrderID     string                   `json:"origin_order_id" gorm:"column:origin_order_id;type:varchar(255);not null;comment:原订单id"`
	OrderPrice        int                      `json:"order_price" gorm:"column:order_price;type:int(10);not null;default:0;comment:关联订单金额"`
	PayTime           int                      `json:"pay_time" gorm:"column:pay_time;type:int(10);not null;default:0;comment:交易时间"`
	OrderType         int32                    `json:"order_type" gorm:"column:order_type;type:tinyint(1);not null;default:0;comment:1匹配成功,2取消匹配,3邀约过期"`
	NewCustomer       int32                    `json:"new_customer" gorm:"column:new_customer;type:tinyint(1);not null;default:0;comment:新客，老客"`
	ReserveOrderGoods []ReserveOrderGoodsModel `json:"reserve_order_goods" gorm:"foreignKey:RelationOrderID,UpdatedAt;references:RelationOrderID,UpdatedAt"`
	Order             ReserveOrderModel        `json:"-" gorm:"foreignKey:ReserveNumber;references:ReserveNumber"`
}

// TableName returns the table name of the ReserveOrderSales model
func (ReserveOrderSalesModel) TableName() string {
	return "reserve_order_sales"
}
