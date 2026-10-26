package microScrmModel

import "micro-scrm/internal/app/models"

type WorkExternalUserTagModel struct {
	models.MicroBaseModel
	RelationID    int64 `json:"relation_id" gorm:"column:relation_id;type:bigint(20);not null;index:relation_id_index;default:0;comment:work_external_user_relation表id"`
	ExternalTagID int64 `json:"external_tag_id" gorm:"column:external_tag_id;type:bigint(20);not null;index:external_tag_id_index;default:0;comment:标签表标签id"`
}

// TableName returns the table name of the WorkExternalUserTag model
func (WorkExternalUserTagModel) TableName() string {
	return "work_external_user_tag"
}
