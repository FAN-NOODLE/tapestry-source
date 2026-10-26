package source

import (
	"micro-scrm/internal/orm/model"
)

type SourceModel struct {
	model.MicroScrmBaseModel

	GroupID int64 `gorm:"column:group_id;type:bigint(20);not null;index:group_id_index;default:0;comment:分组ID 0 未分组"`
	Type    int32 `gorm:"column:type;type:tinyint(3);not null;default:0;comment:类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)"`
	AdminID int64 `gorm:"column:admin_id;type:bigint(20);not null;index:admin_id_index;default:0;comment:操作人ID"`
}

// TableName 素材主表
func (SourceModel) TableName() string {
	return "source"
}
