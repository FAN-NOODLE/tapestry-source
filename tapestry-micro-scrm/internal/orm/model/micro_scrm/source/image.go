package source

import (
	"micro-scrm/internal/orm/model"
)

type ImageModel struct {
	model.MicroScrmBaseModel

	Name            string `gorm:"column:name;type:varchar(255);not null;comment:图片名称"`
	URL             string `gorm:"column:url;type:varchar(255);not null;comment:图片url"`
	Size            int    `gorm:"column:size;type:int(11);not null"`
	Ext             string `gorm:"column:ext;type:varchar(255);not null"`
	MediaID         string `gorm:"column:media_id;type:varchar(128);not null;comment:企业微信临时素材ID"`
	MediaExpireTime int    `gorm:"column:media_expire_time;type:int(10);not null;default:0;comment:企业微信临时素材过期时间"`
	MediaStatus     int32  `gorm:"column:media_status;type:tinyint(1);not null;default:0;comment:企业微信临时素材更新状态0=未更新 1=更新成功 2=更新失败"`
	MediaFailReason string `gorm:"column:media_fail_reason;type:varchar(255);not null;comment:企业微信临时素材更新错误原因"`
}

// TableName 素材库图片扩展表
func (ImageModel) TableName() string {
	return "source_image"
}
