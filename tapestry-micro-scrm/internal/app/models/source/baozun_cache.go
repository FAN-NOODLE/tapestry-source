package source

import "micro-scrm/internal/app/models"

type BaozunCacheModel struct {
	models.MicroBaseModel
	Type            int    `json:"type" gorm:"column:type;not null;default:0;comment:类型(2 图片 5 视频)"`
	IsMain          int    `json:"is_main" gorm:"column:is_main;not null;default:0;comment:是否是主图 0  否 1是"`
	Name            string `json:"name" gorm:"column:name;not null;default:'';comment:图片名称"`
	URL             string `json:"url" gorm:"column:url;not null;default:'';comment:图片url"`
	ProductName     string `json:"product_name" gorm:"column:product_name;not null;default:'';comment:商品名称"`
	SkuItem         string `json:"sku_item" gorm:"column:sku_item;not null;default:'';comment:skuItem"`
	UpcCode         string `json:"upc_code" gorm:"column:upc_code;not null;default:'';comment:upcCode"`
	Color           string `json:"color" gorm:"column:color;not null;default:'';comment:color"`
	ResourceKey     string `json:"resource_key" gorm:"column:resource_key;not null;default:'';comment:资源key"`
	Size            int64  `json:"size" gorm:"column:size;not null"`
	Ext             string `json:"ext" gorm:"column:ext;not null"`
	MediaId         string `json:"media_id" gorm:"column:media_id;not null;default:'';comment:企业微信临时素材ID"`
	MediaExpireTime int64  `json:"media_expire_time" gorm:"column:media_expire_time;not null;default:0;comment:企业微信临时素材过期时间"`
	MediaStatus     int    `json:"media_status" gorm:"column:media_status;not null;default:0;comment:企业微信临时素材更新状态0=未更新 1=更新成功 2=更新失败"`
	MediaFailReason string `json:"media_fail_reason" gorm:"column:media_fail_reason;not null;default:'';comment:企业微信临时素材更新错误原因"`
	DownloadTimes   int    `json:"download_times" gorm:"column:download_times;not null;default:0;comment:下载次数"`
	ForwardingTimes int    `json:"forwarding_times" gorm:"column:forwarding_times;not null;default:0;comment:转发次数"`
	CopyTimes       int    `json:"copy_times" gorm:"column:copy_times;not null;default:0;comment:复制次数"`
}

// TableName returns the table name of the DashboardArea model
func (BaozunCacheModel) TableName() string {
	return "source_baozun_cache"
}

const BrandIsMainImage = 1
