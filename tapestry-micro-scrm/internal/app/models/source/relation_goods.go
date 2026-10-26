package source

import "micro-scrm/internal/app/models"

type RelationGoodsModel struct {
	models.MicroBaseModel
	SourceID int64  `gorm:"column:source_id;type:bigint(20);not null;default:0;comment:素材id" json:"source_id"`
	Type     int8   `gorm:"column:type;type:tinyint(3);not null;default:0;comment:类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件 7 朋友圈组合)" json:"type"`
	UPC      string `gorm:"column:upc;type:varchar(60);not null;default:'';comment:upc" json:"upc"`
	SKU      string `gorm:"column:sku;type:varchar(60);not null;default:'';comment:sku" json:"sku"`
	Name     string `gorm:"column:name;type:varchar(60);not null;default:'';comment:商品名称" json:"name"`
	Url      string `gorm:"column:url;type:varchar(255);not null;default:'';comment:图片地址" json:"url"`
	Color    string `gorm:"column:color;type:varchar(60);not null;default:'';comment:颜色" json:"color"`
	SkuCode  string `gorm:"column:sku_code;type:varchar(64);not null;default:'';comment:款 色" json:"sku_code"`
	AddType  int8   `gorm:"column:add_type;type:tinyint(1);not null;default:1;comment:添加方式 1 手动添加 2 AI识别" json:"add_type"`
	Score    string `gorm:"column:score;type:varchar(32);not null;default:'';comment:置信度（AI识别）" json:"score"`
}

func (RelationGoodsModel) TableName() string {
	return "source_relation_goods"
}

// IdWithType  商品检索结果：source_id + 素材类型
type IdWithType struct {
	SourceID int64 `gorm:"column:source_id"`
	Type     int8  `gorm:"column:type"`
}
