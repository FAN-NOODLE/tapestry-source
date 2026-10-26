package userGroupModel

import (
	"micro-scrm/internal/app/models"
)

type UserGroupFileModel struct {
	models.BaseModel
	FileName  string `json:"file_name" gorm:"column:file_name;type:varchar(255);not null"`
	FileType  uint32 `json:"file_type" gorm:"column:file_type;type:tinyint(1) unsigned;not null;default:0"`
	FileURL   string `json:"file_url" gorm:"column:file_url;type:varchar(255);not null;comment:文件的url地址"`
	IsExtract uint32 `json:"is_extract" gorm:"column:is_extract;type:tinyint(1) unsigned;not null;default:0;comment:数据文件是否已经提取完成：1: 文件数据提取完成；2、预处理数据完成；3、数据打标签完成"`
}

// TableName returns the table name of the UserGroupFile model
func (UserGroupFileModel) TableName() string {
	return "user_group_file"
}

