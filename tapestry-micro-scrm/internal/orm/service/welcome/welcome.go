package welcome

import (
	"micro-scrm/internal/orm/model/micro_scrm/welcome"
	"micro-scrm/internal/orm/service"
	"micro-scrm/internal/pkg/helper"
)

func (w *WelcomeService) GetDataByTypeId(typeId int64, enterpriseId int64) (res []welcome.WelcomesExternalTypeDataModel, err error) {
	err = w.Scopes(service.ScopeCond_EQDelete).
		Where(map[string]interface{}{
			"welcomes_external_type_id": typeId,
			"enterprise_id":             enterpriseId,
		}).
		Order("code+0 desc").
		Preload("TypeInfo").
		Find(&res).Error
	return
}

func (w WelcomeService) GetTypeId(enterpriseId int64) (int64, error) {
	var res welcome.WelcomesExternalTypeModel
	err := w.Scopes(service.ScopeCond_EQDelete).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
		}).
		Order("updated_at desc").
		First(&res).Error
	return res.ID, err
}

func (w WelcomeService) GetAllType(enterpriseId int64) (res []welcome.WelcomesExternalTypeModel, err error) {
	query := w.Scopes(service.ScopeCond_EQDelete).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
		})
	if helper.CheckCoachTime() {
		query.Where("code != ?", welcome.MemberLevel)
	} else {
		query.Where("code != ?", welcome.MemberLevelFY26)
	}
	err = query.Preload("TypeDataList").
		Find(&res).Error
	return
}

func (w WelcomeService) DelWelcome(id int64) error {
	err := w.Scopes(service.ScopeCond_EQDelete).
		Delete(&welcome.DrainageWelcomesModel{}, id).Error
	return err
}

func (w WelcomeService) GetDetailById(id int64) (data welcome.DrainageWelcomesModel, err error) {
	err = w.Scopes(service.ScopeCond_EQDelete).
		Model(&welcome.DrainageWelcomesModel{}).
		Where("id = ?", id).
		Order("updated_at desc").
		Find(&data).Error
	return
}

func (w WelcomeService) GetWelcomeDetailById(id int64) (res welcome.DrainageWelcomesModel, err error) {
	err = w.Scopes(service.ScopeCond_EQDelete).
		Model(welcome.DrainageWelcomesModel{}).
		Where("id = ?", id).Find(&res).Error
	return
}

func (w WelcomeService) CreateWelcomeDetail(m welcome.DrainageWelcomesModel) (int64, error) {
	err := w.Create(&m).Error
	if err != nil {
		return 0, err
	}
	return m.ID, nil
}

func (w WelcomeService) EditWelcomeDetail(m welcome.DrainageWelcomesModel) (int64, error) {
	err := w.Model(&welcome.DrainageWelcomesModel{}).
		Scopes(service.ScopeCond_EQDelete).
		Where(map[string]interface{}{"id": m.ID}).UpdateColumns(m).Error
	if err != nil {
		return 0, err
	}
	return m.ID, nil
}

func (w WelcomeService) GetAllWelcome(enterpriseId int64, size int, page int, keyWords string) (data []welcome.DrainageWelcomesModel, total int64, err error) {
	if page < 0 {
		page = 1
	}
	if size < 0 {
		size = 10
	}
	s := w.Model(&welcome.DrainageWelcomesModel{}).
		Scopes(service.ScopeLeftCond_EQDelete).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"is_default":    welcome.ISDEFAULT_KD,
			"channel_id":    welcome.CHANNLEID_KD,
			"type":          welcome.TypeSelfHoo,
		})
	if keyWords != "" {
		s.Where("content like ?", "%"+keyWords+"%")
	}
	err = s.Order("created_at desc").Count(&total).Limit(size).Offset((page - 1) * size).Find(&data).Error
	return
}

func (w WelcomeService) GetTypeName(id int64) (string, error) {
	var res welcome.WelcomesExternalTypeDataModel
	err := w.Select("name").Scopes(service.ScopeCond_EQDelete).
		Model(&welcome.WelcomesExternalTypeDataModel{}).
		Where("id = ?", id).Find(&res).Error
	return res.Name, err
}
