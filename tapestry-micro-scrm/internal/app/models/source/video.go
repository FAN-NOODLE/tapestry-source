package source

import (
	"micro-scrm/internal/app/models"
)

type VideoModel struct {
	models.MicroBaseModel

	Name            string `gorm:"column:name;type:varchar(255);not null;comment:视频名称"`
	CoverURL        string `gorm:"column:cover_url;type:varchar(255);not null;comment:视频封面"`
	URL             string `gorm:"column:url;type:varchar(255);not null"`
	Status          int32  `gorm:"column:status;type:tinyint(4);not null;default:0;comment:封面图"`
	FileID          string `gorm:"column:file_id;type:varchar(200);not null;default:0;comment:第三方 ID"`
	Ext             string `gorm:"column:ext;type:varchar(255);not null;comment:文件后缀"`
	MediaID         string `gorm:"column:media_id;type:varchar(128);not null;comment:企业微信临时素材ID"`
	MediaExpireTime int    `gorm:"column:media_expire_time;type:int(10);not null;default:0;comment:企业微信临时素材过期时间"`
	MediaStatus     int32  `gorm:"column:media_status;type:tinyint(1);not null;default:0;comment:企业微信临时素材更新状态0=未更新 1=更新成功 2=更新失败"`
	MediaFailReason string `gorm:"column:media_fail_reason;type:varchar(255);not null;comment:企业微信临时素材更新错误原因"`
	Size            string `gorm:"column:size;type:varchar(20);not null;default:0;comment:文件大小"`

	Source SourceModel `json:"source" gorm:"foreignKey:ID;references:ID"`
}

type VideoCategoryRelationDetail struct {
	VideoModel

	SourceId      int64  `json:"source_id"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	CategoryIDs   string `json:"category_ids"`
	CategoryNames string `json:"category_names"`
	GroupId       int64  `json:"group_id"`
	GroupName     string `json:"group_name"`
}

// TableName 素材视频子表
func (VideoModel) TableName() string {
	return "source_video"
}
