package source

import (
	"micro-scrm/internal/app/models"
)

type CategoryRelationModel struct {
	models.MicroBaseModel
	CategoryID uint64 `json:"category_id" gorm:"column:category_id;type:bigint(20) unsigned;not null;index:category_id_idx;default:0;comment:关联类别表id"`
	SourceID   uint64 `json:"source_id" gorm:"column:source_id;type:bigint(20) unsigned;not null;index:source_id_idx;default:0;comment:关联资源表id"`
	Type       int32  `json:"type" gorm:"column:type;type:tinyint(3);not null;default:0;comment:类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)"`
}

// TableName returns the table name of the SourceCategoryRelation model
func (CategoryRelationModel) TableName() string {
	return "source_category_relation"
}
