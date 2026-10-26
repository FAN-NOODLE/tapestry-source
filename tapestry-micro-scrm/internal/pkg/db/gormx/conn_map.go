package gormx

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
	"micro-scrm/internal/pkg/db"
	"strings"
	"sync"
)

var databasesConn sync.Map

func InitConnection(conf db.Config) error {
	open, err := sqlx.Open(conf.Driver, fmt.Sprintf(conf.DynamicDsn, "micro_scrm"))
	if err != nil {
		return err
	}
	defer func() {
		_ = open.Close()
	}()

	var databases []string
	err = open.Select(&databases, "SHOW DATABASES")
	if err != nil || len(databases) == 0 {
		return err
	}

	for _, database := range databases {
		if !strings.Contains(database, "micro_") {
			continue
		}
		microDatabase := &Config{
			Dsn:                fmt.Sprintf(conf.DynamicDsn, database),
			MaxOpenConnections: conf.Pool.MaxOpenConns,
			MaxIdleConnections: conf.Pool.MaxIdleConns,
			MaxLifeTime:        int(conf.Pool.ConnMaxLifetime),
		}
		microConn, err := connect(microDatabase, MysqlDriver)
		if err != nil {
			panic(err)
		}
		databasesConn.Store(database, microConn)
	}

	return err
}

// GetConnection 获取数据库连接
func GetConnection(name string) *gorm.DB {
	// 判断map中是否存在数据库链接对象key
	conn, ok := databasesConn.Load(name)
	if !ok {
		return nil
	}

	return conn.(*gorm.DB)
}
