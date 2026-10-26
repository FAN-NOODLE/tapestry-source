package source

import (
	"micro-scrm/internal/app/models"
)

type ImageModel struct {
	models.MicroBaseModel

	Name            string      `gorm:"column:name;type:varchar(255);not null;comment:图片名称" json:"name"`
	URL             string      `gorm:"column:url;type:varchar(255);not null;comment:图片url" json:"url"`
	Size            int         `gorm:"column:size;type:int(11);not null" json:"size"`
	Ext             string      `gorm:"column:ext;type:varchar(255);not null" json:"ext"`
	MediaID         string      `gorm:"column:media_id;type:varchar(128);not null;comment:企业微信临时素材ID" json:"media_id"`
	MediaExpireTime int         `gorm:"column:media_expire_time;type:int(10);not null;default:0;comment:企业微信临时素材过期时间" json:"media_expire_time"`
	MediaStatus     int32       `gorm:"column:media_status;type:tinyint(1);not null;default:0;comment:企业微信临时素材更新状态0=未更新 1=更新成功 2=更新失败" json:"media_status"`
	MediaFailReason string      `gorm:"column:media_fail_reason;type:varchar(255);not null;comment:企业微信临时素材更新错误原因" json:"media_fail_reason"`
	Source          SourceModel `json:"source" gorm:"foreignKey:ID;references:ID" json:"source,omitempty"`
}

// TableName 素材库图片扩展表
func (ImageModel) TableName() string {
	return "source_image"
}

type ImageCategoryRelationDetail struct {
	ImageModel
	SourceId      int64  `json:"source_id"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	CategoryIDs   string `json:"category_ids"`
	CategoryNames string `json:"category_names"`
	GroupId       int64  `json:"group_id"`
	GroupName     string `json:"group_name"`
	MergeId       int64  `json:"merge_id"`
}

type EndSourceIds struct {
	ID int64 `json:"id"`
}
