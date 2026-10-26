package microScrmDao

import (
	"context"
	"gorm.io/gorm/clause"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type UnionidMapDao struct {
	dao       dao.Dao
	tableName string
}

func NewUnionidMap() *UnionidMapDao {
	return &UnionidMapDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.UnionidMapModel{}.TableName(),
	}
}

func (d *UnionidMapDao) CreateInBatches(c app.Ctx, enterpriseId int64, data []microScrmModel.UnionidMapModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(&data, len(data)).Error
	return err
}

// GetBatch  中添加
func (d *UnionidMapDao) GetBatch(c context.Context, enterpriseId int64, lastID int64, limit int) ([]microScrmModel.UnionidMapModel, error) {
	var rows = make([]microScrmModel.UnionidMapModel, 0)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("id > ?", lastID).Order("id asc").Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}
func (d *UnionidMapDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.UnionidMapModel, error) {
	var result []microScrmModel.UnionidMapModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}
