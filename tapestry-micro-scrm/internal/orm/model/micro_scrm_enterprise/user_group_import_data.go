package micro_scrm_enterprise

import "micro-scrm/internal/orm/model"

type UserGroupImportDataModel struct {
	model.BranchModel

	FileName   string `json:"file_name" gorm:"column:file_name;type:varchar(255);not null"`
	FileID     int64    `json:"file_id" gorm:"column:file_id;type:int(11);not null;default:0;comment:user_group_file表的主键id"`
	UnionID    string `json:"union_id" gorm:"column:union_id;type:varchar(255);not null;comment:导入的unionid"`
	CreateTime string `json:"create_time" gorm:"column:create_time;type:varchar(255);not null;comment:推送时间"`
}

// TableName returns the table name of the UserGroupImportData model
func (UserGroupImportDataModel) TableName() string {
	return "user_group_import_data"
}
