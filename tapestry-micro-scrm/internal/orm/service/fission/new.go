package fission

import (
	"gorm.io/gorm"
)

type FissionService struct {
	*gorm.DB
	ServiceName string
}

// 创建服务
func NewFissionService(g *gorm.DB) *FissionService {
	return &FissionService{
		DB:          g,
		ServiceName: "fission",
	}
}
