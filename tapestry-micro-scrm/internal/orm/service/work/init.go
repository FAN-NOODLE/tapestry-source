package work

import "gorm.io/gorm"

// work数据库服务

type WorkDBService struct {
	*gorm.DB
	ServiceName string
}

func NewWorkDBService(db *gorm.DB) *WorkDBService {
	return &WorkDBService{
		DB:          db,
		ServiceName: "work数据库服务",
	}
}
