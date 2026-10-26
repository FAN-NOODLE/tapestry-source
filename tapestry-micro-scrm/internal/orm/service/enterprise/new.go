package enterprise

import (
	"gorm.io/gorm"
)

type EnterpriseService struct {
	*gorm.DB
	ServiceName string
}

// 创建服务
func NewEnterpriseService(g *gorm.DB) *EnterpriseService {
	return &EnterpriseService{
		DB:          g,
		ServiceName: "enterprise",
	}
}
