package enterprise

import (
	"gorm.io/gorm"
)

type EnterpriseBossService struct {
	*gorm.DB
	ServiceName string
}

// NewEnterpriseBossService 创建服务
func NewEnterpriseBossService(g *gorm.DB) *EnterpriseBossService {
	return &EnterpriseBossService{
		DB:          g,
		ServiceName: "enterprise",
	}
}
