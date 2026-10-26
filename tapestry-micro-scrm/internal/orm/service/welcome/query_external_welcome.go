package welcome

import (
	"micro-scrm/internal/orm/model/micro_scrm/welcome"
	ormService "micro-scrm/internal/orm/service"
)

func (service *WelcomeService) GetWelcomesExternalTypeDataAll(param welcome.WelcomesExternalTypeDataModel) ([]welcome.WelcomesExternalTypeDataModel, error) {
	var result []welcome.WelcomesExternalTypeDataModel
	err := service.Scopes(ormService.ScopeCond_EQDelete).Where(&param).Order("updated_at desc").Find(&result).Error
	return result, err
}

// QueryExternalWelcome 查询外部自定义欢迎语
func (service *WelcomeService) QueryExternalWelcome(param welcome.DrainageWelcomesModel) (data welcome.DrainageWelcomesModel, err error) {
	err = service.Scopes(ormService.ScopeCond_EQDelete).Where(&param).Order("updated_at desc").First(&data).Error

	return
}

// QueryWelcomeExternalData 查询外部查询条件数据
func (service *WelcomeService) QueryWelcomeExternalData(param welcome.WelcomesExternalTypeDataModel) (data welcome.WelcomesExternalTypeDataModel, err error) {
	err = service.Scopes(ormService.ScopeCond_EQDelete).
		Where(&param).
		Where("welcomes_external_type_id > 0").
		First(&data).Error

	return
}
