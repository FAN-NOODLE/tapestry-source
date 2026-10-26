package auth

import (
	"gorm.io/gorm"
)

type AuthService struct {
	*gorm.DB
	ServiceName string
}

func NewAuthService(g *gorm.DB) *AuthService {
	return &AuthService{
		DB:          g,
		ServiceName: "auth-service",
	}
}
