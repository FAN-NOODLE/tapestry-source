package microScrmModel

import "micro-scrm/internal/app/models"

type OrderAchievementModel struct {
	models.BaseModel
	ID               int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20)"`
	MemberID         int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);comment:导购id"`
	Userid           string `json:"userid" gorm:"column:userid;type:varchar(64);comment:工号"`
	StoreCode        string `json:"store_code" gorm:"column:store_code;type:varchar(255)"`
	OrderTime        int64  `json:"order_time" gorm:"column:order_time;type:int(11)"`
	OrderTotal       int64  `json:"order_total" gorm:"column:order_total;type:int(11)"`
	OrderPriceTotal  int64  `json:"order_price_total" gorm:"column:order_price_total;type:int(11)"`
	PayOrderTotal    int64  `json:"pay_order_total" gorm:"column:pay_order_total;type:int(11)"`
	RefundOrderTotal int64  `json:"refund_order_total" gorm:"column:refund_order_total;type:int(11)"`
	FileNumber       int64  `json:"file_number" gorm:"column:file_number;type:int(11)"`
}

// TableName returns the table name of the OrderAchievement model
func (OrderAchievementModel) TableName() string {
	return "order_achievement"
}
