package excitationModels

import "micro-scrm/internal/app/models"

const (
	OrderStatusNormal = 1
	OrderStatusCancel = 2
)

type ExcitationMallOrderModel struct {
	models.MicroBaseModel
	UserID      string  `json:"user_id" gorm:"column:user_id;type:varchar(50);not null;comment:导购编码"`
	MemberID    int64   `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;comment:员工id"`
	OrderNum    string  `json:"order_num" gorm:"column:order_num;type:varchar(100);not null;index:order_num_index;comment:订单号"`
	OrderCoins  float64 `json:"order_coins" gorm:"column:order_coins;type:float;not null;default:0;comment:金币金额"`
	OrderStatus int     `json:"order_status" gorm:"column:order_status;type:tinyint(1);not null;default:0;comment:订单状态1 下单成功 2退单成功"`
}

// TableName returns the table name of the ExcitationMallOrder model
func (ExcitationMallOrderModel) TableName() string {
	return "excitation_mall_order"
}
