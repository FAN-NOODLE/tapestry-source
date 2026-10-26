package group_chat

import (
	"gorm.io/gorm"
)

type GroupChatService struct {
	*gorm.DB
	ServiceName string
}

// NewGroupChatService 创建服务
func NewGroupChatService(g *gorm.DB) *GroupChatService {
	return &GroupChatService{
		DB:          g,
		ServiceName: "group-chat",
	}
}
