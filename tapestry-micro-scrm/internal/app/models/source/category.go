package source

type CategoryModel struct {
	ID           uint64 `json:"id" gorm:"primaryKey;column:id;type:bigint(20) unsigned"`
	EnterpriseID uint64 `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20) unsigned;not null;default:0;comment:企业ID"`
	Name         string `json:"name" gorm:"column:name;type:varchar(255);not null;index:name_idx;comment:名称"`
	CreatedAt    uint   `json:"created_at" gorm:"column:created_at;type:int(10) unsigned;not null;default:0;comment:创建时间"`
	UpdatedAt    uint   `json:"updated_at" gorm:"column:updated_at;type:int(10) unsigned;not null;default:0;comment:修改时间"`
	DeletedAt    uint   `json:"deleted_at" gorm:"column:deleted_at;type:int(10) unsigned;not null;default:0;comment:删除时间"`
}

// TableName returns the table name of the SourceCategory model
func (CategoryModel) TableName() string {
	return "source_category"
}
