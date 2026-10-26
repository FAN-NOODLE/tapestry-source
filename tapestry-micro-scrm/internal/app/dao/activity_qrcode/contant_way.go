package activityQRCodeDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	activityQRCodeModel "micro-scrm/internal/app/models/activity_qrcode"
)

type ContactWayDao struct {
	dao.EnterpriseRepo[activityQRCodeModel.ActivityQRCodeContactWayModel]
	model     *activityQRCodeModel.ActivityQRCodeContactWayModel
	tableName string
}

func NewContactWayDao() *ContactWayDao {
	return &ContactWayDao{
		EnterpriseRepo: dao.EnterpriseRepo[activityQRCodeModel.ActivityQRCodeContactWayModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(activityQRCodeModel.ActivityQRCodeContactWayModel),
		tableName: activityQRCodeModel.ActivityQRCodeContactWayModel{}.TableName(),
	}
}
