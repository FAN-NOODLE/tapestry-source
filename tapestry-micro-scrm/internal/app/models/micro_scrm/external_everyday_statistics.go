package microScrmModel

import "micro-scrm/internal/app/models"

type ExternalEverydayStatisticsModel struct {
	models.BaseModel
	MemberID            int64  `json:"member_id" gorm:"column:member_id;type:bigint(10) unsigned;not null;index:member_id_index;default:0;comment:导购member表id"`
	UnixTime            int    `json:"unix_time" gorm:"column:unix_time;type:int(10);not null;index:unix_time_index;default:0;comment:时间戳"`
	Date                string `json:"date" gorm:"column:date;type:varchar(64);not null;comment:日期"`
	NewCustomer         int    `json:"new_customer" gorm:"column:new_customer;type:int(10);not null;default:0;comment:新客"`
	NewCustomerChange   int    `json:"new_customer_change" gorm:"column:new_customer_change;type:int(10);not null;default:0;comment:新客变化数量"`
	OldCustomer         int    `json:"old_customer" gorm:"column:old_customer;type:int(10);not null;default:0;comment:老客"`
	OldCustomerChange   int    `json:"old_customer_change" gorm:"column:old_customer_change;type:int(10);not null;default:0;comment:老客变化数量"`
	SleepCustomer       int    `json:"sleep_customer" gorm:"column:sleep_customer;type:int(10);not null;default:0;comment:沉睡客"`
	SleepCustomerChange int    `json:"sleep_customer_change" gorm:"column:sleep_customer_change;type:int(10);not null;default:0;comment:沉睡客变化数量"`
	VipCustomer         int    `json:"vip_customer" gorm:"column:vip_customer;type:int(10);not null;default:0;comment:vip客"`
	VipCustomerChange   int    `json:"vip_customer_change" gorm:"column:vip_customer_change;type:int(10);not null;default:0;comment:vip客变化数量"`
}

// TableName returns the table name of the ExternalEverydayStatistics model
func (ExternalEverydayStatisticsModel) TableName() string {
	return "external_everyday_statistics"
}
