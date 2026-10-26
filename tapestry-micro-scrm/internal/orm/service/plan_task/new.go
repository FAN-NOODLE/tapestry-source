package plan_task

import (
	"gorm.io/gorm"
)

type PlanTaskService struct {
	*gorm.DB
	ServiceName string
}

// NewPlanTaskService 创建服务
func NewPlanTaskService(g *gorm.DB) *PlanTaskService {
	return &PlanTaskService{
		DB:          g,
		ServiceName: "plan_task_service",
	}
}
