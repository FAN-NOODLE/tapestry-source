package microScrmModel

const (
	ONLINE_SALES_VOLUME = iota + 1
	ALL_SALES_VOLUME
)

type MemberSalesVolumeModel struct {
	ID           int64   `json:"id" gorm:"primaryKey;column:id;type:int(10) auto_increment"`
	MemberID     int64   `json:"member_id" gorm:"column:member_id;type:bigint(10) unsigned;not null;index:member_id_start_time_deleted_at_index;default:0;comment:导购member表id"`
	DepartmentID int64   `json:"department_id" gorm:"column:department_id;type:int(10) unsigned;not null;default:0;comment:departments表主键"`
	StartTime    int64   `json:"start_time" gorm:"column:start_time;type:int(10) unsigned;not null;index:member_id_start_time_deleted_at_index;default:0;comment:0点时间戳"`
	Date         string  `json:"date" gorm:"column:date;type:varchar(32);not null;comment:日期"`
	StoreCode    string  `json:"store_code" gorm:"column:store_code;type:varchar(32);not null;comment:门店编码"`
	HomeStore    float64 `json:"home_store" gorm:"column:home_store;type:decimal(9,2);not null;default:0.00"`
	LlEshop      float64 `json:"ll_eshop" gorm:"column:ll_eshop;type:decimal(9,2);not null;default:0.00"`
	LlLivestream float64 `json:"ll_livestream" gorm:"column:ll_livestream;type:decimal(9,2);not null;default:0.00"`
	PpLivestream float64 `json:"pp_livestream" gorm:"column:pp_livestream;type:decimal(9,2);not null;default:0.00"`
	PpPost       float64 `json:"pp_post" gorm:"column:pp_post;type:decimal(9,2);not null;default:0.00"`
	RemoteSales  float64 `json:"remote_sales" gorm:"column:remote_sales;type:decimal(9,2);not null;default:0.00"`
	CloudStore   float64 `json:"cloud_store" gorm:"column:cloud_store;type:decimal(9,2);not null;default:0.00"`
	O2o          float64 `json:"o2o" gorm:"column:o2o;type:decimal(9,2);not null;default:0.00"`
	TryAtHome    float64 `json:"try_at_home" gorm:"column:try_at_home;type:decimal(9,2) unsigned;not null;default:0.00"`
	LinkToEc     float64 `json:"link_to_ec" gorm:"column:link_to_ec;type:decimal(9,2);not null;default:0.00"`
	CreatedAt    int64   `json:"created_at" gorm:"column:created_at;type:int(10) unsigned;not null;default:0;comment:创建时间"`
	UpdatedAt    int64   `json:"updated_at" gorm:"column:updated_at;type:int(10) unsigned;not null;default:0;comment:修改时间"`
	DeletedAt    int64   `json:"deleted_at" gorm:"column:deleted_at;type:int(10) unsigned;not null;index:member_id_start_time_deleted_at_index;default:0;comment:删除时间"`
}

// TableName returns the table name of the MemberSalesVolume model
func (MemberSalesVolumeModel) TableName() string {
	return "member_sales_volume"
}

type CpsRecord struct {
	Data     string `json:"data"`
	SaId     string `json:"saId"`
	StoreId  string `json:"store_id"`
	LinkToEc string `json:"link_to_ec"`
}
