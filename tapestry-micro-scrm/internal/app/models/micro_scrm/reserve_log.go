package microScrmModel

import "micro-scrm/internal/app/models"

type ReserveLogModel struct {
	models.BaseModel
	ReserveOrderID int64  `json:"reserve_order_id" gorm:"column:reserve_order_id;type:bigint(20);not null;default:0;comment:预约单id"`
	LogInfo        string `json:"log_info" gorm:"column:log_info;type:text;comment:操作日志"`
	OrderInfo      string `json:"order_info" gorm:"column:order_info;type:text;comment:订单信息"`
	LogType        int32  `json:"log_type" gorm:"column:log_type;type:tinyint(1);not null;default:1;comment:日志类型"`
}

// TableName returns the table name of the ReserveLog model
func (ReserveLogModel) TableName() string {
	return "reserve_log"
}
