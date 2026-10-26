package source

import (
	"gorm.io/gorm"
)

type SourceService struct {
	*gorm.DB
	ServiceName string
}

// NewPlanTaskService 创建服务
func NewSourceService(g *gorm.DB) *SourceService {
	return &SourceService{
		DB:          g,
		ServiceName: "source_service",
	}
}
