package user_event

import "gorm.io/gorm"

type UserEventService struct {
	*gorm.DB
	ServiceName string
}

func NewUserEventService(g *gorm.DB) *UserEventService {
	return &UserEventService{
		DB:          g,
		ServiceName: "user_event",
	}
}
