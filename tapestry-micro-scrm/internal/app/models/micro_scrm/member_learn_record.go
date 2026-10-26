package microScrmModel

type MemberLearnRecordModel struct {
	ID           int64  `json:"id" gorm:"primaryKey;column:id;type:int(10) auto_increment"`
	MemberID     int64  `json:"member_id" gorm:"column:member_id;type:bigint(10) unsigned;not null;index:member_id_deleted_at_start_time_index;default:0;comment:导购member表id"`
	DepartmentID int64  `json:"department_id" gorm:"column:department_id;type:int(10) unsigned;not null;default:0;comment:departments表主键"`
	StartTime    int64  `json:"start_time" gorm:"column:start_time;type:int(10) unsigned;not null;index:member_id_deleted_at_start_time_index;default:0;comment:时间戳"`
	StoreCode    string `json:"store_code" gorm:"column:store_code;type:varchar(32);not null;comment:门店编码"`
	TitleRemark  string `json:"title_remark" gorm:"column:title_remark;type:varchar(64);not null;comment:培训类目名称"`
	Points       int64  `json:"points" gorm:"column:points;type:int(10);not null;default:0;comment:积分"`
	CreatedAt    int64  `json:"created_at" gorm:"column:created_at;type:int(10) unsigned;not null;default:0;comment:创建时间"`
	UpdatedAt    int64  `json:"updated_at" gorm:"column:updated_at;type:int(10) unsigned;not null;default:0;comment:修改时间"`
	DeletedAt    int64  `json:"deleted_at" gorm:"column:deleted_at;type:int(10) unsigned;not null;default:0;comment:删除时间"`
}

// TableName returns the table name of the MemberLearnRecord model
func (MemberLearnRecordModel) TableName() string {
	return "member_learn_record"
}
