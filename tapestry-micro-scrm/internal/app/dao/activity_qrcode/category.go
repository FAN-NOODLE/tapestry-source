package activityQRCodeDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	activityQRCodeModel "micro-scrm/internal/app/models/activity_qrcode"
)

type CategoryDao struct {
	dao.EnterpriseRepo[activityQRCodeModel.ActivityQRCodeCategoryModel]
	model     *activityQRCodeModel.ActivityQRCodeCategoryModel
	tableName string
}

func NewCategoryDao() *CategoryDao {
	return &CategoryDao{
		EnterpriseRepo: dao.EnterpriseRepo[activityQRCodeModel.ActivityQRCodeCategoryModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(activityQRCodeModel.ActivityQRCodeCategoryModel),
		tableName: activityQRCodeModel.ActivityQRCodeCategoryModel{}.TableName(),
	}
}
