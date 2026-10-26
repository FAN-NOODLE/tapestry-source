package drainageContactDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	drainageContactModel "micro-scrm/internal/app/models/drainage"
)

type DrainageContactGroupWelcomesDao struct {
	dao.EnterpriseRepo[drainageContactModel.DrainageContactsWelcomesModel]
	tableName string
	model     *drainageContactModel.DrainageContactsWelcomesModel
}

func NewDrainageContactGroupWelcomesDao() *DrainageContactGroupWelcomesDao {
	return &DrainageContactGroupWelcomesDao{
		EnterpriseRepo: dao.EnterpriseRepo[drainageContactModel.DrainageContactsWelcomesModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		tableName: drainageContactModel.DrainageContactsWelcomesModel{}.TableName(),
		model:     new(drainageContactModel.DrainageContactsWelcomesModel),
	}
}
