package micro_scrm

import (
	"gorm.io/gorm"
	"micro-scrm/internal/orm/model"
)

type WorkExternalTagGroupModel struct {
	model.MicroScrmBaseModel

	GroupID    string `gorm:"column:group_id;type:varchar(60);not null;index:group_id_index;comment:企业微信标签组id"`
	GroupName  string `gorm:"column:group_name;type:varchar(255);not null;comment:标签组名称"`
	Order      int    `gorm:"column:order;type:int(10);not null;default:0;comment:标签组次序值"`
	CreateTime int    `gorm:"column:create_time;type:int(10);not null;default:0;comment:标签组创建时间"`

	TagList []WorkExternalTagModel `gorm:"foreignKey:GroupID;references:ID"`
}

// TableName returns the table name of the struct_name model
func (WorkExternalTagGroupModel) TableName() string {
	return "work_external_tag_group"
}

// 查询后的值替换
func (w *WorkExternalTagGroupModel) AfterFind(tx *gorm.DB) (err error) {

	return
}
