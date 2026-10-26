package activityQRCodeModel

import (
	"micro-scrm/internal/app/models"
)

type ActivityQRCodeMainModel struct {
	models.BaseModel
	Name           string `json:"name" gorm:"column:name;type:varchar(255);not null;default:'';comment:活动名称"`
	AliasName      string `json:"alias_name" gorm:"column:alias_name;type:varchar(255);not null;default:'';comment:对外名称"`
	BannerUrl      string `json:"banner_url" gorm:"column:banner_url;type:varchar(255);not null;default:'';comment:banner图url"`
	TimeType       int    `json:"time_type" gorm:"column:time_type;type:tinyint(1);not null;default:0;comment:时间类型，0 长期 1 指定时间"`
	StartTime      int    `json:"start_time" gorm:"column:start_time;type:int(10);not null;default:0;comment:开始时间"`
	EndTime        int    `json:"end_time" gorm:"column:end_time;type:int(10);not null;default:0;comment:结束时间"`
	CategoryIds    string `json:"category_ids" gorm:"column:category_ids;type:varchar(255);not null;default:'';comment:投放渠道id集合"`
	AddNameTag     int    `json:"add_name_tag" gorm:"column:add_name_tag;type:tinyint(1);default:0;comment:是否添加活动名称标签"`
	AddCategoryTag int    `json:"add_category_tag" gorm:"column:add_category_tag;type:tinyint(1);not null;default:0;comment:是否添加投放渠道标签"`
	TagIds         string `json:"tag_ids" gorm:"column:tag_ids;type:varchar(255);not null;default:'';comment:标签id集合，英文逗号分割"`
	NameTagId      int64  `json:"name_tag_id" gorm:"column:name_tag_id;type:int(10);not null;default:0;comment:活动名称标签id"`
	SkipVerify     int    `json:"skip_verify" gorm:"column:skip_verify;type:tinyint(1);not null;default:1;comment:是否无需验证"`
}

func (ActivityQRCodeMainModel) TableName() string {
	return "activity_qrcode_main"
}

const TimeTypeLongTerm = 0
