package source

import (
	"micro-scrm/internal/app/models"
)

type MiniProgramModel struct {
	models.BaseModel

	MiniProgramID int64  `json:"mini_program_id" gorm:"column:mini_program_id;type:bigint(20);not null;default:0;comment:mini_program 表主键"`
	Title         string `json:"title" gorm:"column:title;type:varchar(255);not null;comment:卡片标题"`
	Path          string `json:"path" gorm:"column:path;type:varchar(255);not null;comment:小程序页面路径"`
	CoverURL      string `json:"cover_url" gorm:"column:cover_url;type:varchar(255);not null;comment:封面图"`

	Source SourceModel `json:"source" gorm:"foreignKey:ID;references:ID"`
}

// TableName returns the table name of the SourceFile model
func (MiniProgramModel) TableName() string {
	return "source_mini_program"
}

type MiniProgramCategoryRelationDetail struct {
	MiniProgramModel
	SourceId      int64  `json:"source_id"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	CategoryIDs   string `json:"category_ids"`
	CategoryNames string `json:"category_names"`
	GroupId       int64  `json:"group_id"`
	GroupName     string `json:"group_name"`
}
