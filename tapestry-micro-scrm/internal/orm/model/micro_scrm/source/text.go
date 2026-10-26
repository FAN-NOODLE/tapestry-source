package source

import (
	"micro-scrm/internal/orm/model"
)

type TextModel struct {
	model.MicroScrmBaseModel

	Content string `gorm:"column:content;type:text;not null;comment:内容"`
}

// TableName 文本素材表
func (TextModel) TableName() string {
	return "source_text"
}
