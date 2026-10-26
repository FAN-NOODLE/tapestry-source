package welcome

import (
	"gorm.io/gorm"
)

type WelcomeService struct {
	*gorm.DB
	ServiceName string
}

// NewWelcomeService 创建服务
func NewWelcomeService(g *gorm.DB) *WelcomeService {
	return &WelcomeService{
		DB:          g,
		ServiceName: "欢迎语模块",
	}
}
