package fission

import (
	"micro-scrm/internal/orm/model"
)

type FissionMpUserModel struct {
	model.MicroScrmBaseModel

	OpenID    string `gorm:"column:open_id;type:varchar(100);not null;uniqueIndex:unq_open_id_enterprise"`
	UnionID   string `gorm:"column:union_id;type:varchar(100);not null"`
	NickName  string `gorm:"column:nick_name;type:varchar(255);not null;comment:昵称"`
	Sex       int32  `gorm:"column:sex;type:tinyint(1) unsigned;not null;default:0;comment:值为1时是男性，值为2时是女性，值为0时是未知"`
	Province  string `gorm:"column:province;type:varchar(100);not null;comment:省份"`
	City      string `gorm:"column:city;type:varchar(255);not null;comment:城市"`
	Country   string `gorm:"column:country;type:varchar(255);not null;comment:国家"`
	ImgURL    string `gorm:"column:img_url;type:varchar(255);not null;comment:头像"`
	Privilege string `gorm:"column:privilege;type:varchar(255);not null;comment:用户特权信息"`
}

func (FissionMpUserModel) TableName() string {
	return "fission_mp_user"
}
