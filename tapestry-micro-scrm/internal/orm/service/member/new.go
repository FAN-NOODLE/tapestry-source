package member

import (
	"gorm.io/gorm"
)

type MemberService struct {
	*gorm.DB
	ServiceName string
}

// 创建服务
func NewMemberService(g *gorm.DB) *MemberService {
	return &MemberService{
		DB:          g,
		ServiceName: "member",
	}
}
