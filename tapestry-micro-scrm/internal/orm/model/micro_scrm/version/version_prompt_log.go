package version

import "micro-scrm/internal/orm/model"

type VersionPromptLogModel struct {
	model.MicroScrmBaseModel

	AdminID         int64 `json:"admin_id" gorm:"column:admin_id;type:bigint(20);not null;index:admin_id_index;default:0;comment:账号ID"`
	VersionRecordID int64 `json:"version_record_id" gorm:"column:version_record_id;type:bigint(11);not null;index:version_record_id_index;default:0;comment:version_record.id"`
}

func (VersionPromptLogModel) TableName() string {
	return "version_prompt_log"
}
