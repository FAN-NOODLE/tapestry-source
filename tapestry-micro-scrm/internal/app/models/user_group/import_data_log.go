package userGroupModel

import "micro-scrm/internal/app/models"

type UserGroupImportDataLogModel struct {
	models.BaseModel

	UserGroupID int64  `json:"user_group_id" gorm:"column:user_group_id;type:bigint(20);not null;default:0"`
	UnionID     string `json:"union_id" gorm:"column:union_id;type:varchar(255);not null"`
	ImportRes   string `json:"import_res" gorm:"column:import_res;type:varchar(255);not null"`
}

// TableName returns the table name of the UserGroupImportDataLog model
func (UserGroupImportDataLogModel) TableName() string {
	return "user_group_import_data_log"
}
