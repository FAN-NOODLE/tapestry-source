package source

import (
	"micro-scrm/internal/orm/model"
)

type LinkModel struct {
	model.MicroScrmBaseModel
	Title       string `gorm:"column:title;type:varchar(255);not null;comment:标题"`
	Description string `gorm:"column:description;type:varchar(255);not null;comment:链接描述"`
	CoverURL    string `gorm:"column:cover_url;type:varchar(255);not null;comment:封面图URL"`
	LinkType    int32  `gorm:"column:link_type;type:tinyint(3);not null;default:0;comment:类型(1 URL 2 自定义)"`
	LinkURL     string `gorm:"column:link_url;type:varchar(255);not null;comment:链接URL"`
	Author      string `gorm:"column:author;type:varchar(255);not null;comment:作者"`
	Content     string `gorm:"column:content;type:text;not null;comment:正文"`
}

// TableName 链接素材表
func (LinkModel) TableName() string {
	return "source_link"
}
