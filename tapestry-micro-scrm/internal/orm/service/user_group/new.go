package user_group

import "gorm.io/gorm"

// 客户分组服务

type UserGroupService struct {
	*gorm.DB
	ServiceName string
}

// 创建服务
func NewUserGroupService(g *gorm.DB) *UserGroupService {
	return &UserGroupService{
		DB:          g,
		ServiceName: "user_group",
	}
}
