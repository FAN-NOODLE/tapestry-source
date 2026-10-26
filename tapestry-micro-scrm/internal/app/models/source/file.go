package source

import "micro-scrm/internal/app/models"

type FileModel struct {
	models.BaseModel

	Name            string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:文件名称"`
	URL             string `json:"url" gorm:"column:url;type:varchar(255);not null;comment:文件url"`
	AttachName      string `json:"attach_name" gorm:"column:attach_name;type:varchar(255);not null;comment:第三方文件名称"`
	AttachID        int64  `json:"attach_id" gorm:"column:attach_id;type:bigint(20);not null;default:0;comment:第三方file ID"`
	FileType        int32  `json:"file_type" gorm:"column:file_type;type:tinyint(3);not null;default:0;comment:文件类型 1 word 2 excel 3 ppt 4 pdf 6 其他"`
	Ext             string `json:"ext" gorm:"column:ext;type:varchar(255);not null;comment:文件后缀"`
	MediaID         string `json:"media_id" gorm:"column:media_id;type:varchar(128);not null;comment:企业微信临时素材ID"`
	MediaExpireTime int    `json:"media_expire_time" gorm:"column:media_expire_time;type:int(10);not null;default:0;comment:企业微信临时素材过期时间"`
	MediaStatus     int32  `json:"media_status" gorm:"column:media_status;type:tinyint(1);not null;default:0;comment:企业微信临时素材更新状态0=未更新 1=更新成功 2=更新失败"`
	MediaFailReason string `json:"media_fail_reason" gorm:"column:media_fail_reason;type:varchar(255);not null;comment:企业微信临时素材更新错误原因"`
	Size            int    `json:"size" gorm:"column:size;type:int(11);not null;default:0;comment:文件大小"`

	Source SourceModel `json:"source" gorm:"foreignKey:ID;references:ID"`
}

// TableName returns the table name of the SourceFile model
func (FileModel) TableName() string {
	return "source_file"
}

type FileCategoryRelationDetail struct {
	FileModel
	SourceId      int64  `json:"source_id"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	CategoryIDs   string `json:"category_ids"`
	CategoryNames string `json:"category_names"`
	GroupId       int64  `json:"group_id"`
	GroupName     string `json:"group_name"`
}
