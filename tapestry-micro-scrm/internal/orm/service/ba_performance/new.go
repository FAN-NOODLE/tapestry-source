package ba_performance

import (
	"gorm.io/gorm"
)

type BaPerformanceService struct {
	*gorm.DB
	ServiceName string
}

// 创建服务
func NewBaPerformanceService(g *gorm.DB) *BaPerformanceService {
	return &BaPerformanceService{
		DB:          g,
		ServiceName: "ba_performance",
	}
}
