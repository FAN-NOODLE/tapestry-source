package departments

import (
	"gorm.io/gorm"
)

type DepartmentsService struct {
	*gorm.DB
	ServiceName string
}

// NewDepartmentsService 创建服务
func NewDepartmentsService(g *gorm.DB) *DepartmentsService {
	return &DepartmentsService{
		DB:          g,
		ServiceName: "departments",
	}
}
