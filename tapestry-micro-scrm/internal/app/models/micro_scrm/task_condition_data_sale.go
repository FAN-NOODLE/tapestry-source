package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type TaskConditionDataSaleModel struct {
	models.BaseModel
	PlanID        int64  `json:"plan_id" gorm:"column:plan_id;type:bigint(20);not null;index:plan_id_index,plan_task_id_index;default:0;comment:主任务id"`
	ConditionID   int    `json:"condition_id" gorm:"column:condition_id;type:int(10);not null;comment:task_condition_data 分表的id"`
	TaskID        int64  `json:"task_id" gorm:"column:task_id;type:bigint(20);not null;index:plan_task_id_index;default:0;comment:子任务id"`
	UnionID       string `json:"union_id" gorm:"column:union_id;type:varchar(255);not null;comment:顾客的unionid"`
	UserID        string `json:"user_id" gorm:"column:user_id;type:varchar(255);not null;index:user_id_index;comment:员工的工号"`
	StoreCode     string `json:"store_code" gorm:"column:store_code;type:varchar(255);not null;comment:交易门店"`
	OrderNo       string `json:"order_no" gorm:"column:order_no;type:varchar(255);not null;comment:匹配的订单id"`
	OriginOrderNo string `json:"origin_order_no" gorm:"column:origin_order_no;type:varchar(255);not null;comment:原订单ID"`
	OrderPrice    int    `json:"order_price" gorm:"column:order_price;type:int(10);not null;default:0;comment:订单金额"`
	OrderPayTime  int    `json:"order_pay_time" gorm:"column:order_pay_time;type:int(10);not null;default:0;comment:订单交易时间"`
	OrderStatus   int32  `json:"order_status" gorm:"column:order_status;type:tinyint(1);not null;default:0;comment:订单匹配类型，1匹配成功订单，2取消匹配订单"`
}

// TableName returns the table name of the TaskConditionDataSale model
func (TaskConditionDataSaleModel) TableName() string {
	return "task_condition_data_sale"
}
