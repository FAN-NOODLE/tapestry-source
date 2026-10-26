package source

import (
	"micro-scrm/internal/orm/model"
)

type GroupModel struct {
	model.MicroScrmBaseModel
	IgnoreIndex int    `json:"ignore_index" gorm:"column:ignore_index;type:tinyint(1);not null"`
	GroupName   string `gorm:"column:group_name;type:varchar(255);not null" json:"group_name"`
	ParentId    int64  `gorm:"column:parent_id;NOT NULL" json:"parent_id"` // 父级id
	Orders      int    `gorm:"column:orders;NOT NULL" json:"orders"`       // 排序字段

}

// TableName 素材分组表
func (GroupModel) TableName() string {
	return "source_group"
}
