package source

import (
	"micro-scrm/internal/app/models"
)

type LinkModel struct {
	models.BaseModel
	EnterpriseID int64 `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20);not null;default:0"`

	Title       string `json:"title" gorm:"column:title;type:varchar(255);not null;comment:标题"`
	Description string `json:"description" gorm:"column:description;type:varchar(255);not null;comment:链接描述"`
	CoverURL    string `json:"cover_url" gorm:"column:cover_url;type:varchar(1000);not null;comment:封面图URL"`
	LinkType    int32  `json:"link_type" gorm:"column:link_type;type:tinyint(3);not null;default:0;comment:类型(1 URL 2 自定义)"`
	LinkURL     string `json:"link_url" gorm:"column:link_url;type:varchar(1000);not null;comment:链接URL"`
	Author      string `json:"author" gorm:"column:author;type:varchar(255);not null;comment:作者"`
	Content     string `json:"content" gorm:"column:content;type:text;not null;comment:正文"`

	Source SourceModel `json:"source" gorm:"foreignKey:ID;references:ID"`
}

// TableName returns the table name of the SourceFile model
func (LinkModel) TableName() string {
	return "source_link"
}

type LinkCategoryRelationDetail struct {
	LinkModel
	SourceId      int64  `json:"source_id"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	CategoryIDs   string `json:"category_ids"`
	CategoryNames string `json:"category_names"`
	GroupId       int64  `json:"group_id"`
	GroupName     string `json:"group_name"`
}
