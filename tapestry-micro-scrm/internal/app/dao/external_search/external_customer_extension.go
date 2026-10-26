package externalSearchDao

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	"micro-scrm/internal/pkg/db/gormx"
)

const extensionBatchSize = 1000

var customerCouponColumns = map[string]struct{}{
	"coupon_status":                  {},
	"coupon_by_30days_status":        {},
	"birthday_coupon_status":         {},
	"renew_membership_coupon_status": {},
	"welcome_coupon_status":          {},
	"upgrade_coupon_status":          {},
}

type ExternalCustomerExtension struct {
	dao       dao.Dao
	tableName string
}

func NewExternalCustomerExtension() *ExternalCustomerExtension {
	return &ExternalCustomerExtension{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: externalSearchModel.ExternalCustomerExtensionModel{}.TableName(),
	}
}

func (d *ExternalCustomerExtension) UpsertReserveData(c app.Ctx, enterpriseId int64, data []externalSearchModel.ExternalCustomerExtensionModel) error {
	if len(data) == 0 {
		return nil
	}
	return d.upsert(c, enterpriseId, data, map[string]interface{}{
		"valid_reserve_at":   gorm.Expr("GREATEST(valid_reserve_at, VALUES(valid_reserve_at))"),
		"success_reserve_at": gorm.Expr("GREATEST(success_reserve_at, VALUES(success_reserve_at))"),
		"updated_at":         gorm.Expr("GREATEST(updated_at, VALUES(updated_at))"),
	})
}

func (d *ExternalCustomerExtension) UpsertTaskData(c app.Ctx, enterpriseId int64, data []externalSearchModel.ExternalCustomerExtensionModel) error {
	if len(data) == 0 {
		return nil
	}
	return d.upsert(c, enterpriseId, data, map[string]interface{}{
		"task_finish_at":  gorm.Expr("GREATEST(task_finish_at, VALUES(task_finish_at))"),
		"task_created_at": gorm.Expr("GREATEST(task_created_at, VALUES(task_created_at))"),
		"updated_at":      gorm.Expr("GREATEST(updated_at, VALUES(updated_at))"),
	})
}

func (d *ExternalCustomerExtension) UpsertCouponData(c app.Ctx, enterpriseId int64, data []externalSearchModel.ExternalCustomerExtensionModel, columns []string) error {
	if len(data) == 0 || len(columns) == 0 {
		return nil
	}
	assignments := map[string]interface{}{
		"coupon_updated_at": gorm.Expr("GREATEST(coupon_updated_at, VALUES(coupon_updated_at))"),
		"updated_at":        gorm.Expr("GREATEST(updated_at, VALUES(updated_at))"),
	}
	for _, column := range columns {
		if _, ok := customerCouponColumns[column]; !ok {
			return fmt.Errorf("unsupported coupon column: %s", column)
		}
		assignments[column] = gorm.Expr(fmt.Sprintf("IF(VALUES(coupon_updated_at) >= coupon_updated_at, VALUES(%s), %s)", column, column))
	}
	return d.upsert(c, enterpriseId, data, assignments)
}

// UpsertBackfillData 将 CK 中已有的客户扩展字段幂等回填到企业 MySQL。
// 时间字段只前进；优惠券状态仅在来源更新时间不旧于现有数据时覆盖，避免回填覆盖实时影子写入。
func (d *ExternalCustomerExtension) UpsertBackfillData(c app.Ctx, enterpriseId int64, data []externalSearchModel.ExternalCustomerExtensionModel) error {
	if len(data) == 0 {
		return nil
	}
	updates := clause.Set{
		{Column: clause.Column{Name: "valid_reserve_at"}, Value: gorm.Expr("GREATEST(valid_reserve_at, VALUES(valid_reserve_at))")},
		{Column: clause.Column{Name: "success_reserve_at"}, Value: gorm.Expr("GREATEST(success_reserve_at, VALUES(success_reserve_at))")},
		{Column: clause.Column{Name: "task_finish_at"}, Value: gorm.Expr("GREATEST(task_finish_at, VALUES(task_finish_at))")},
		{Column: clause.Column{Name: "task_created_at"}, Value: gorm.Expr("GREATEST(task_created_at, VALUES(task_created_at))")},
	}
	for _, column := range []string{
		"coupon_status",
		"coupon_by_30days_status",
		"birthday_coupon_status",
		"renew_membership_coupon_status",
		"welcome_coupon_status",
		"upgrade_coupon_status",
	} {
		updates = append(updates, clause.Assignment{
			Column: clause.Column{Name: column},
			Value:  gorm.Expr(fmt.Sprintf("IF(VALUES(coupon_updated_at) >= coupon_updated_at, VALUES(%s), %s)", column, column)),
		})
	}
	// coupon_updated_at 必须最后更新，确保上面的状态表达式与回填前的时间比较。
	updates = append(updates,
		clause.Assignment{Column: clause.Column{Name: "coupon_updated_at"}, Value: gorm.Expr("GREATEST(coupon_updated_at, VALUES(coupon_updated_at))")},
		clause.Assignment{Column: clause.Column{Name: "updated_at"}, Value: gorm.Expr("GREATEST(updated_at, VALUES(updated_at))")},
	)
	return d.dao.WithContext(c).
		Scopes(d.dao.BranchDataBase(enterpriseId, d.tableName)).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "unionid"}},
			DoUpdates: updates,
		}).
		CreateInBatches(&data, extensionBatchSize).Error
}

func (d *ExternalCustomerExtension) upsert(c app.Ctx, enterpriseId int64, data []externalSearchModel.ExternalCustomerExtensionModel, assignments map[string]interface{}) error {
	return d.dao.WithContext(c).
		Scopes(d.dao.BranchDataBase(enterpriseId, d.tableName)).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "unionid"}},
			DoUpdates: clause.Assignments(assignments),
		}).
		CreateInBatches(&data, extensionBatchSize).Error
}
