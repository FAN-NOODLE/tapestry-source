package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，大部分表都有的字段
type BaseModel struct {
	ID        int64 `gorm:"primaryKey;column:id;type:bigint(20)" json:"id"`
	CreatedAt int64 `gorm:"column:created_at;type:int(10);not null;default:0" json:"created_at"`
	UpdatedAt int64 `gorm:"column:updated_at;type:int(10);not null;default:0" json:"updated_at"`
	DeletedAt int64 `gorm:"column:deleted_at;type:int(10);not null;default:0" json:"deleted_at"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.CreatedAt = time.Now().Unix()
	b.UpdatedAt = time.Now().Unix()
	return
}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = time.Now().Unix()
	return
}

//func (b *BaseModel) AfterFind(tx *gorm.DB) (err error) {
//	b.UpdatedAt = 0
//	return
//}

type MicroScrmBaseModel struct {
	BaseModel
	EnterpriseID int64 `gorm:"column:enterprise_id;type:bigint(20);not null;index:enterprise_id_index;default:0" json:"enterprise_id"`
}

// BranchModel 分库表需要嵌入模型
type BranchModel struct {
	BaseModel
}

// DBPrefix 分库的数据库前缀
func (BranchModel) DBPrefix() string {
	return "micro_scrm_enterprise_"
}

type ITableName interface {
	DBPrefix() string
	TableName() string
}

// BranchTable 切库
func BranchTable(model ITableName, enterpriseID int64) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tableName := fmt.Sprintf("%s%d.%s", model.DBPrefix(), enterpriseID, model.TableName())
		return tx.Table(tableName)
	}
}
