package gormx

import (
	"gorm.io/gorm"
	"micro-scrm/internal/pkg/db"
)

var gormCkConn *gorm.DB

// GetInstance 获取gorm连接
func GetInstance() *gorm.DB {
	return GetConnection("micro_scrm")
}

func GetCkInstance() *gorm.DB {
	return gormCkConn
}

func GetBossInstance() *gorm.DB {
	return GetConnection("micro_boss")
}

func GetWorkInstance() *gorm.DB {
	return GetConnection("micro_work")
}

func InitCKGorm(conf db.Config) (err error) {
	dbConfig := &Config{
		Dsn:                conf.DynamicDsn,
		MaxOpenConnections: conf.Pool.MaxOpenConns,
		MaxIdleConnections: conf.Pool.MaxIdleConns,
		MaxLifeTime:        int(conf.Pool.ConnMaxLifetime),
	}

	gormCkConn, err = connect(dbConfig, ClickhouseDriver)
	return err
}
