package work_external_user

import "micro-scrm/internal/orm/model"

type WorkExternalUserTagModel struct {
	model.MicroScrmBaseModel

	RelationID    int64 `gorm:"column:relation_id;type:bigint(20);not null;index:relation_id_index;default:0;comment:work_external_user_relation表id"` // work_external_user_relation表id
	ExternalTagID int64 `gorm:"column:external_tag_id;type:bigint(20);not null;default:0;comment:标签表标签id"`                                            // 标签表标签id
}

// TableName returns the table name of the struct_name model
func (WorkExternalUserTagModel) TableName() string {
	return "work_external_user_tag"
}
