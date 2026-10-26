package version

import "micro-scrm/internal/orm/model"

type VersionRecordModel struct {
	model.BaseModel

	VersionNumber string `json:"version_number" gorm:"column:version_number;type:varchar(20);not null;comment:版本号"`
	Theme         string `json:"theme" gorm:"column:theme;type:varchar(50);not null;comment:版本主题"`
	VersionInfo   string `json:"version_info" gorm:"column:version_info;type:text;comment:版本记，json的形式记录"`
}

func (VersionRecordModel) TableName() string {
	return "version_record"
}
