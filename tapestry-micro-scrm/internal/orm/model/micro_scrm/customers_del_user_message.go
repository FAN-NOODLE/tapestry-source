package micro_scrm

type CustomersDelUserMessageModel struct {
	ID           int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20)"`
	EnterpriseID int64  `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20);not null;uniqueIndex:index_user_id;default:0"`
	DpID         int64  `json:"dp_id" gorm:"column:dp_id;type:bigint(20);not null;comment:部门Id"`
	SendType     uint32 `json:"send_type" gorm:"column:send_type;type:tinyint(3) unsigned;not null;default:0;comment:类别 0 选择所有 1客户删除员工 2员工删除客户"`
	UserID       string `json:"user_id" gorm:"column:user_id;type:varchar(50);not null;comment:被通知员工id"`
	CreatedAt    int    `json:"created_at" gorm:"column:created_at;type:int(10);not null;default:0"`
	UpdatedAt    int    `json:"updated_at" gorm:"column:updated_at;type:int(10);not null;default:0"`
	DeletedAt    int    `json:"deleted_at" gorm:"column:deleted_at;type:int(10);not null;uniqueIndex:index_user_id;default:0"`
}

// TableName 删除客户员工消息配置表
func (CustomersDelUserMessageModel) TableName() string {
	return "customers_del_user_message"
}
