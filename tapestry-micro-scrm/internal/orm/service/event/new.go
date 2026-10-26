package event

import "gorm.io/gorm"

type EventService struct {
	*gorm.DB
	ServiceName string
}

// 创建服务
func NewEventService(g *gorm.DB) *EventService {
	return &EventService{
		DB:          g,
		ServiceName: "enterprise",
	}
}
