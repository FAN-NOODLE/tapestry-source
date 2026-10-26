package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"micro-scrm/internal/app"
)

type EnterpriseRepo[T schema.Tabler] struct {
	Dao
	model T
}

// QueryClone 克隆一个带原条件的新实例,WithContext必须放在Session后，否则不会带上之前的条件
func (d *EnterpriseRepo[T]) QueryClone(c context.Context, query *gorm.DB) *gorm.DB {
	return query.Session(&gorm.Session{NewDB: true}).WithContext(c)
}

// SetDb 切换db连接,用于事物
func (d *EnterpriseRepo[T]) SetDb(db *gorm.DB) *EnterpriseRepo[T] {
	d.DB = db
	return d
}

func (d *EnterpriseRepo[T]) GetOne(c app.Ctx, enterpriseId int64, query any, args ...any) (T, error) {
	var result T
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *EnterpriseRepo[T]) GetOneUnscoped(c app.Ctx, enterpriseId int64, query any, args ...any) (T, error) {
	var result T
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Model(d.model).Unscoped().Where(query, args...).Find(&result).Error
	return result, err
}

func (d *EnterpriseRepo[T]) GetMany(c app.Ctx, enterpriseId int64, query any, args ...any) ([]T, error) {
	var result []T
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *EnterpriseRepo[T]) CreateData(c app.Ctx, enterpriseId int64, data T) error {
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Model(d.model).Create(&data).Error
	return err
}

func (d *EnterpriseRepo[T]) CreateInBatchesData(c app.Ctx, enterpriseId int64, data []T) error {
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *EnterpriseRepo[T]) SaveData(c app.Ctx, enterpriseId int64, data T) error {
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Model(d.model).Save(&data).Error
	return err
}

// UpdateData 该方法不会更新update_at字段
func (d *EnterpriseRepo[T]) UpdateData(c app.Ctx, enterpriseId int64, query any, data map[string]any) error {
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *EnterpriseRepo[T]) DeleteData(c app.Ctx, enterpriseId int64, query any, args ...any) error {
	var model T
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Model(d.model).Where(query, args...).Delete(&model).Error
	return err
}

func (d *EnterpriseRepo[T]) ForceDelete(c app.Ctx, enterpriseId int64, query any, args ...any) error {
	var model T
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).Model(d.model).Where(query, args...).Unscoped().Delete(&model).Error
	return err
}

func (d *EnterpriseRepo[T]) Increment(c app.Ctx, enterpriseId int64, query any, column string, num int64) error {
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).
		Model(d.model).
		Where(query).
		UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s + ?", column), num)).Error
	return err
}

func (d *EnterpriseRepo[T]) Decrement(c app.Ctx, enterpriseId int64, query any, column string, num int64) error {
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).
		Model(d.model).
		Where(query).
		UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s - ?", column), num)).Error
	return err
}
