package workModels

import "micro-scrm/internal/app/models"

type WorkExternalTagGroupModel struct {
	models.MicroBaseModel
	GroupID    string `json:"group_id" gorm:"column:group_id;type:varchar(60);not null;index:group_id_index;comment:企业微信标签组id"`
	GroupName  string `json:"group_name" gorm:"column:group_name;type:varchar(255);not null;comment:标签组名称"`
	Order      int    `json:"order" gorm:"column:order;type:int(10);not null;default:0;comment:标签组次序值"`
	IsThirdTag int32  `json:"is_third_tag" gorm:"column:is_third_tag;type:tinyint(1);not null;default:0;comment:是否三方标签 0否 1是"`
	Type       int32  `json:"type" gorm:"column:type;type:tinyint(1);not null;default:1;comment:标签组类型 1-基础标签组 2-第三方标签组"`
	CreateTime int    `json:"create_time" gorm:"column:create_time;type:int(10);not null;default:0;comment:标签组创建时间"`
}

// TableName returns the table name of the WorkExternalTagGroup model
func (WorkExternalTagGroupModel) TableName() string {
	return "work_external_tag_group"
}
