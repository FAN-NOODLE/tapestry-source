package brand_config

import (
	"gorm.io/gorm"
)

type BrandConfigService struct {
	*gorm.DB
	ServiceName string
}

// NewBrandConfigService 创建服务
func NewBrandConfigService(g *gorm.DB) *BrandConfigService {
	return &BrandConfigService{
		DB:          g,
		ServiceName: "品牌配置信息",
	}
}
