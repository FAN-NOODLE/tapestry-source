package tagModels

import "micro-scrm/internal/app/models"

type ImpressionTagGroupModel struct {
	ThirdGroupID string `json:"third_group_id" gorm:"column:third_group_id;type:varchar(60);not null;index:group_id_index;comment:三方标签组id"`
	GroupName    string `json:"group_name" gorm:"column:group_name;type:varchar(255);not null;comment:标签组名称"`
	models.MicroBaseModel
}

// TableName 外部联系人印象标签组表
func (ImpressionTagGroupModel) TableName() string {
	return "impression_tag_group"
}
