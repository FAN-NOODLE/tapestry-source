package work_external_user

import "gorm.io/gorm"

type WorkExternalUserService struct {
	*gorm.DB
	ServiceName string
}

// 创建服务
func NewWorkExternalUserService(g *gorm.DB) *WorkExternalUserService {
	return &WorkExternalUserService{
		DB:          g,
		ServiceName: "work_external_user",
	}
}
