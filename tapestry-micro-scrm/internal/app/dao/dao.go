package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
)

type Dao struct {
	*gorm.DB
}

// BranchTable 切分库表
func (d *Dao) BranchTable(table string, enterpriseID int64) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tableName := fmt.Sprintf("%s%d.%s", app.Cfg().DB.DBPrefix, enterpriseID, table)
		return tx.Table(tableName)
	}
}

func (d *Dao) GetTable(table string, enterpriseID int64) string {
	return fmt.Sprintf("%s%d.%s", app.Cfg().DB.DBPrefix, enterpriseID, table)
}

func (d *Dao) BranchDataBase(enterpriseID int64, table string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tableName := fmt.Sprintf("%s%d.%s", app.Cfg().DB.DBPrefix, enterpriseID, table)
		return tx.Table(tableName)
	}
}

// QueryClone 克隆一个带原条件的新实例,WithContext必须放在Session后，否则不会带上之前的条件
func (d *Dao) QueryClone(c context.Context, query *gorm.DB) *gorm.DB {
	return query.Session(&gorm.Session{NewDB: true}).WithContext(c)
}
