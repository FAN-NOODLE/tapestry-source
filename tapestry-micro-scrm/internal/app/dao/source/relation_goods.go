package source

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/source"
)

type RelationGoodsDao struct {
	dao.Repo[source.RelationGoodsModel]
	model     *source.RelationGoodsModel
	tableName string
}

func NewRelationGoodsDao() *RelationGoodsDao {
	return &RelationGoodsDao{
		Repo: dao.Repo[source.RelationGoodsModel]{
			DB: app.DB(),
		},
		model:     new(source.RelationGoodsModel),
		tableName: source.RelationGoodsModel{}.TableName(),
	}
}
