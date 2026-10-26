package tagModels

import "micro-scrm/internal/app/models"

type ImpressionTagRelationModel struct {
	ExternalUserID int64  `json:"external_user_id" gorm:"column:external_user_id;type:bigint(20);not null;index:external_user_id_index;default:0;comment:外部联系人ID"`
	TagID          int64  `json:"tag_id" gorm:"column:tag_id;type:bigint(20);not null;index:tag_id_index;default:0;comment:外部联系人印象标签表ID"`
	UserID         string `json:"user_id" gorm:"column:user_id;type:varchar(50);not null;comment:企微用户Id"`
	models.MicroBaseModel
}

// TableName 外部联系人印象标签关系表
func (ImpressionTagRelationModel) TableName() string {
	return "impression_tag_relation"
}
