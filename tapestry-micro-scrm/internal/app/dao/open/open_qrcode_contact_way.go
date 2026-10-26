package openDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	openModel "micro-scrm/internal/app/models/open"
)

type OpenQrcodeContactWayDao struct {
	dao.EnterpriseRepo[openModel.OpenQrcodeContactWayModel]
	model     *openModel.OpenQrcodeContactWayModel
	tableName string
}

func NewOpenQrcodeContactWayDao() *OpenQrcodeContactWayDao {
	return &OpenQrcodeContactWayDao{
		EnterpriseRepo: dao.EnterpriseRepo[openModel.OpenQrcodeContactWayModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(openModel.OpenQrcodeContactWayModel),
		tableName: openModel.OpenQrcodeContactWayModel{}.TableName(),
	}
}
