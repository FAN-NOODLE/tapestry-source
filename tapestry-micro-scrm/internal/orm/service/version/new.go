package version

import (
	"gorm.io/gorm"
)

type VersionService struct {
	*gorm.DB
	ServiceName string
}

// NewVersionService 创建服务
func NewVersionService(g *gorm.DB) *VersionService {
	return &VersionService{
		DB:          g,
		ServiceName: "version",
	}
}
