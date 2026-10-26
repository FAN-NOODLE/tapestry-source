package source

import (
	"micro-scrm/internal/app/models"
)

type CollectRelationModel struct {
	models.MicroBaseModel
	MemberId    int64  `json:"member_id" gorm:"column:member_id;type:bigint(20) unsigned;not null;default:0;comment:members表id"`
	SourceID    int64  `json:"source_id" gorm:"column:source_id;type:bigint(20) unsigned;not null;default:0;comment:素材库表id"`
	Type        int32  `json:"type" gorm:"column:type;type:tinyint(3);not null;default:0;comment:类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)"`
	IsCollect   int32  `json:"is_collect" gorm:"column:is_collect;type:tinyint(3);not null;default:0;comment:是否收藏，0不收藏，1收藏"`
	ProductCode string `json:"product_code" gorm:"column:product_code;type:varchar(255);not null;default:'';comment:官网素材库的商品code"`
	MergeId     int64  `json:"merge_id" gorm:"column:merge_id;type:bigint(20) unsigned;not null;default:0;comment:朋友圈素材合并id"`
}

// TableName returns the table name of the SourceCategoryRelation model
func (CollectRelationModel) TableName() string {
	return "source_collect_relation"
}

const (
	SourceCollectYes = 1 // 收藏
	SourceCollectNo  = 0 // 未收藏
)

type ImgMergeCollect struct {
	Num     int   `json:"num"`
	MergeId int64 `json:"merge_id"`
}
