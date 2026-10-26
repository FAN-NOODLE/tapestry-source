package microScrmDao

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type EnterpriseUser struct {
	dao       dao.Dao
	tableName string
}

func NewEnterpriseUser() *EnterpriseUser {
	return &EnterpriseUser{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.EnterpriseUserModel{}.TableName(),
	}
}

func (d *EnterpriseUser) GetInfo(c app.Ctx, param microScrmModel.EnterpriseUserModel, enterpriseId int64) (microScrmModel.EnterpriseUserModel, error) {
	var result microScrmModel.EnterpriseUserModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(&param).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return result, err
}

// GetBatch 按id 批量查询
func (d *EnterpriseUser) GetBatch(c context.Context, enterpriseId int64, lastID int64, limit int) ([]microScrmModel.EnterpriseUserModel, error) {
	var rows = make([]microScrmModel.EnterpriseUserModel, 0)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("id > ?", lastID).Order("id asc").Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *EnterpriseUser) UpdateById(c app.Ctx, data microScrmModel.EnterpriseUserModel, enterpriseId, id int64) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("id", id).UpdateColumns(data).Error
	return err
}

func (d *EnterpriseUser) DeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model []microScrmModel.EnterpriseUserModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).Unscoped().Delete(&model).Error
	return err
}

func (d *EnterpriseUser) CreateInBatchesData(c app.Ctx, enterpriseId int64, data []microScrmModel.EnterpriseUserModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *EnterpriseUser) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).UpdateColumns(data).Error
	return err
}

func (d *EnterpriseUser) BatchUpdateByUnionIds(c context.Context, enterpriseId int64, mappings []microScrmModel.UnionidMapModel, nowUnix int64) (int64, error) {
	if len(mappings) == 0 {
		return 0, nil
	}

	// 获取当前批次的最小和最大ID
	minID := mappings[0].ID
	maxID := mappings[len(mappings)-1].ID

	query := `
        UPDATE ` + fmt.Sprintf("%s%d.%s", app.Cfg().DB.DBPrefix, enterpriseId, d.tableName) + ` AS t1
        INNER JOIN ` + fmt.Sprintf("%s%d.%s", app.Cfg().DB.DBPrefix, enterpriseId, "unionid_map") + ` AS t2 ON t1.wx_union = CONCAT('{"default":"', t2.old_unionid, '"}')
        SET t1.wx_union = CONCAT('{"default":"', t2.new_unionid, '"}'),
            t1.updated_at = ?
        AND t2.id >= ? AND t2.id <= ?`

	db := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Exec(query, nowUnix, minID, maxID)
	return db.RowsAffected, db.Error
}
