package activityQRCodeDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	activityQRCodeModel "micro-scrm/internal/app/models/activity_qrcode"
)

type MemberScopeDao struct {
	dao.EnterpriseRepo[activityQRCodeModel.ActivityQRCodeMemberScopeModel]
	model     *activityQRCodeModel.ActivityQRCodeMemberScopeModel
	tableName string
}

func NewMemberScopeDao() *MemberScopeDao {
	return &MemberScopeDao{
		EnterpriseRepo: dao.EnterpriseRepo[activityQRCodeModel.ActivityQRCodeMemberScopeModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(activityQRCodeModel.ActivityQRCodeMemberScopeModel),
		tableName: activityQRCodeModel.ActivityQRCodeMemberScopeModel{}.TableName(),
	}
}
