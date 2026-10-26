package activity

import (
	"gorm.io/gorm"
)

type ActivityService struct {
	*gorm.DB
	ServiceName string
}

// NewActivityService 创建服务
func NewActivityService(g *gorm.DB) *ActivityService {
	return &ActivityService{
		DB:          g,
		ServiceName: "activity_service",
	}
}
