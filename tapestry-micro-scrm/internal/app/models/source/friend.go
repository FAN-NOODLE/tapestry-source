package source

import (
	"micro-scrm/internal/app/models"
)

type FriendModel struct {
	models.MicroBaseModel

	TextContent    string `json:"text_content" gorm:"column:text_content;type:text;not null;comment:文本内容"`
	Name           string `json:"name" gorm:"column:name;type:varchar(100);not null;comment:资源名称"`
	VideoCoverURL  string `json:"video_cover_url" gorm:"column:video_cover_url;type:varchar(255);not null;comment:视频封面地址"`
	OriginalURL    string `json:"original_url" gorm:"column:original_url;type:varchar(255);not null;comment:原资源地址"`
	CompressionURL string `json:"compression_url" gorm:"column:compression_url;type:varchar(255);not null;comment:压缩资源地址"`
	Ext            string `json:"ext" gorm:"column:ext;type:varchar(20);not null;comment:文件后缀"`
	Size           int    `json:"size" gorm:"column:size;type:int(10);not null;default:0;comment:文件大小"`
	Type           int32  `json:"type" gorm:"column:type;type:tinyint(1);not null;default:0;comment:类型 1:图片 2:视频"`

	Source SourceModel `json:"source" gorm:"foreignKey:ID;references:ID"`
}

// TableName 素材朋友圈子表
func (FriendModel) TableName() string {
	return "source_friend"
}

type FriendDetail struct {
	FriendModel
	SourceId      int64  `json:"source_id"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	CategoryIDs   string `json:"category_ids"`
	CategoryNames string `json:"category_names"`
	GroupId       int64  `json:"group_id"`
	GroupName     string `json:"group_name"`
	MergeId       int64  `json:"merge_id"`
}

const (
	FriendTypeText = 0 // 文本
	FriendTypeImg  = 1
)
