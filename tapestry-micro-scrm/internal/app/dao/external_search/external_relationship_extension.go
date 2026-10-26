package externalSearchDao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	"micro-scrm/internal/pkg/db/gormx"
)

type ExternalRelationshipExtension struct {
	dao       dao.Dao
	tableName string
}

func NewExternalRelationshipExtension() *ExternalRelationshipExtension {
	return &ExternalRelationshipExtension{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: externalSearchModel.ExternalRelationshipExtensionModel{}.TableName(),
	}
}

func (d *ExternalRelationshipExtension) UpsertData(c app.Ctx, enterpriseId int64, data []externalSearchModel.ExternalRelationshipExtensionModel) error {
	if len(data) == 0 {
		return nil
	}
	return d.dao.WithContext(c).
		Scopes(d.dao.BranchDataBase(enterpriseId, d.tableName)).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "relation_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"unionid":       gorm.Expr("VALUES(unionid)"),
				"user_id":       gorm.Expr("VALUES(user_id)"),
				"interact_time": gorm.Expr("GREATEST(interact_time, VALUES(interact_time))"),
				"active_time":   gorm.Expr("GREATEST(active_time, VALUES(active_time))"),
				"updated_at":    gorm.Expr("GREATEST(updated_at, VALUES(updated_at))"),
			}),
		}).
		CreateInBatches(&data, extensionBatchSize).Error
}
