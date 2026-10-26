package marketing

import (
	"gorm.io/gorm"
)

type MarketingService struct {
	*gorm.DB
	ServiceName string
}

// 创建服务
func NewMarketingService(g *gorm.DB) *MarketingService {
	return &MarketingService{
		DB:          g,
		ServiceName: "marketing",
	}
}
