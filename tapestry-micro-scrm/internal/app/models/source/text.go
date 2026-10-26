package source

import (
	"micro-scrm/internal/app/models"
)

type TextModel struct {
	models.MicroBaseModel

	Content string `json:"content" gorm:"column:content;type:text;not null;comment:正文"`

	Source SourceModel `json:"source" gorm:"foreignKey:ID;references:ID"`
}

// TableName returns the table name of the SourceFile model
func (TextModel) TableName() string {
	return "source_text"
}

type TextCategoryRelationDetail struct {
	TextModel
	SourceId      int64  `json:"source_id"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	CategoryIDs   string `json:"category_ids"`
	CategoryNames string `json:"category_names"`
	GroupId       int64  `json:"group_id"`
	GroupName     string `json:"group_name"`
}
