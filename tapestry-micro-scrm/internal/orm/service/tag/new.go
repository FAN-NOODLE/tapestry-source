package tag

import (
	"gorm.io/gorm"
)

type TagService struct {
	*gorm.DB
	ServiceName string
}

// 创建服务
func NewTagService(g *gorm.DB) *TagService {
	return &TagService{
		DB:          g,
		ServiceName: "tag",
	}
}
