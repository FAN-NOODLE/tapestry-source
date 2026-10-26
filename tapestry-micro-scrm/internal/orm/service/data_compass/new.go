package data_compass

import "gorm.io/gorm"

type DataCompassService struct {
	*gorm.DB
	ServiceName string
}

// 创建服务
func NewDataCompassService(g *gorm.DB) *DataCompassService {
	return &DataCompassService{
		DB:          g,
		ServiceName: "data_compass",
	}
}
