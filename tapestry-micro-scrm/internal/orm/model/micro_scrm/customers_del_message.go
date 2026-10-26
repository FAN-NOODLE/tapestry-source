package micro_scrm

type CustomersDelMessageModel struct {
	ID           int64  `gorm:"primaryKey;column:id;type:bigint(20)"`
	EnterpriseID int64  `gorm:"column:enterprise_id;type:bigint(20);not null;default:0"`
	Status       uint32 `gorm:"column:status;type:tinyint(3) unsigned;not null;default:1;comment:开启开关"`
	SendType     uint32 `gorm:"column:send_type;type:tinyint(3) unsigned;not null;default:0;comment:类别 1客户删除员工 2员工删除客户"`
	StatTxt      string `gorm:"column:stat_txt;type:text;not null;comment:设置员工字段json"`
	QyStatTxt    string `gorm:"column:qy_stat_txt;type:text;not null;comment:设置员工字段json 企业微信相关"`
	CreatedAt    int    `gorm:"column:created_at;type:int(10);not null;default:0"`
	UpdatedAt    int    `gorm:"column:updated_at;type:int(10);not null;default:0"`
	DeletedAt    int    `gorm:"column:deleted_at;type:int(10);not null;default:0"`
}

// TableName 删除客户消息配置表
func (CustomersDelMessageModel) TableName() string {
	return "customers_del_message"
}
