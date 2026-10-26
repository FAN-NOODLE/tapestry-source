package brand_config

import "micro-scrm/internal/orm/model"

type BrandConfigCouponCenterModel struct {
	model.MicroScrmBaseModel

	BrandCode string `json:"brand_code" gorm:"column:brand_code;type:varchar(20);not null;comment:品牌简称"`        // 品牌简称
	BrandName string `json:"brand_name" gorm:"column:brand_name;type:varchar(50);not null;comment:品牌名称"`        // 品牌名称
	UUID      string `json:"uuid" gorm:"column:uuid;type:varchar(50);not null;comment:coupon_center需要的品牌的uuid"` // coupon_center需要的品牌的uuid
}

func (BrandConfigCouponCenterModel) TableName() string {
	return "brand_config_coupon_center"
}
