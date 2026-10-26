package structure

import (
	"gorm.io/gorm"
)

type StructureService struct {
	*gorm.DB
	ServiceName string
}

// NewStructureService 创建服务
func NewStructureService(g *gorm.DB) *StructureService {
	return &StructureService{
		DB:          g,
		ServiceName: "structure",
	}
}
