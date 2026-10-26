package storeVisitModel

import (
	"micro-scrm/internal/app/models"
)

type StoreVisitMediaModel struct {
	models.BaseModel
	VisitID       int64  `json:"visit_id" gorm:"column:visit_id;type:bigint(20);not null;index:visit_id_index;comment:巡店表id"`
	MemberID      int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:member_id_index;comment:导购id"`
	RelID         int64  `json:"rel_id" gorm:"column:rel_id;type:bigint(20);not null;comment:相关业务id"`
	RelType       int    `json:"rel_type" gorm:"column:rel_type;type:tinyint(1);not null;comment:1 表现反馈 2 留言"`
	MediaType     int    `json:"media_type" gorm:"column:media_type;type:tinyint(1);not null;comment:1图片2视频"`
	URL           string `json:"url" gorm:"column:url;type:varchar(256);default:'';comment:图片视频相关url"`
	VideoImageUrl string `json:"video_image_url" gorm:"column:video_image_url;type:varchar(256);default:'';comment:视频图片url"`
	VideoSecond   int    `json:"video_second" gorm:"column:video_second;type:int(10);default:0;comment:视频时长（秒）"`
}

func (StoreVisitMediaModel) TableName() string {
	return "store_visit_media"
}

const ShowRecord = 1
const Comments = 2

const Image = 1
const Video = 2
