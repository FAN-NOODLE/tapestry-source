package source

type MiniProgramModel struct {
	ID            int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20)"`
	EnterpriseID  int64  `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20);not null;default:0;comment:企业id"`
	MiniProgramID int64  `json:"mini_program_id" gorm:"column:mini_program_id;type:bigint(20);not null;default:0;comment:mini_program 表主键"`
	Title         string `json:"title" gorm:"column:title;type:varchar(255);not null;comment:卡片标题"`
	Path          string `json:"path" gorm:"column:path;type:varchar(255);not null;comment:小程序页面路径"`
	CoverURL      string `json:"cover_url" gorm:"column:cover_url;type:varchar(255);not null;comment:封面图"`
	CreatedAt     int    `json:"created_at" gorm:"column:created_at;type:int(11);not null;default:0;comment:创建时间"`
	UpdatedAt     int    `json:"updated_at" gorm:"column:updated_at;type:int(11);not null;default:0;comment:更新时间"`
	DeletedAt     int    `json:"deleted_at" gorm:"column:deleted_at;type:int(11);not null;default:0;comment:删除时间"`
}

// TableName 素材小程序
func (MiniProgramModel) TableName() string {
	return "source_mini_program"
}
