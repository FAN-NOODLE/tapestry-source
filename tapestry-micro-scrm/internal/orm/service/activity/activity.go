package activity

import (
	"gorm.io/gorm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/service"
)

func (s ActivityService) GetTypeList(interpriseId int64, page int, size int) (int64, []micro_scrm_enterprise.OperationActivityModel, error) {
	var data []micro_scrm_enterprise.OperationActivityModel
	if page < 0 {
		page = 1
	}

	if size < 0 {
		size = 10
	}
	var total int64
	err := s.Scopes(service.ScopeCond_EQDelete).
		Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityModel{}, interpriseId)).
		Count(&total).
		Order("updated_at desc").
		Limit(size).
		Offset((page - 1) * size).
		Find(&data).Error
	return total, data, err
}

func (s ActivityService) GetUserName(interpriseId int64, unionid string) (micro_scrm_enterprise.WxUserCorpExternalModel, error) {
	var user micro_scrm_enterprise.WxUserCorpExternalModel
	err := s.Scopes(service.ScopeCond_EQDelete).
		Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, interpriseId)).
		Where("unionid = ?", unionid).
		Find(&user).Error
	return user, err
}

func (s ActivityService) GetUserPhone(interpriseId int64, euid int64) (micro_scrm_enterprise.EnterpriseUserModel, error) {
	var user micro_scrm_enterprise.EnterpriseUserModel
	err := s.Scopes(model.BranchTable(&micro_scrm_enterprise.EnterpriseUserModel{}, interpriseId)).
		Where("id = ?", euid).
		Find(&user).Error
	return user, err
}

func (s ActivityService) InsertUser(interpriseId int64, user micro_scrm_enterprise.UserModel) error {
	err := s.Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityUserModel{}, interpriseId)).Create(&user).Error
	return err
}

func (s ActivityService) UpdateActivityNums(interpriseId int64, id int64) error {
	err := s.Scopes(service.ScopeCond_EQDelete).
		Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityModel{}, interpriseId)).
		Where("id = ?", id).Update("nums", gorm.Expr("nums + ?", 1)).Error
	return err
}

func (s ActivityService) GetActivityId(interpriseId int64) (int64, error) {
	var data micro_scrm_enterprise.OperationActivityModel
	err := s.Scopes(service.ScopeCond_EQDelete).
		Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityModel{}, interpriseId)).
		Where("type = ?", 1).Order("created_at desc").Find(&data).Error
	return data.ID, err
}

func (s ActivityService) GetInfoByUnionid(interpriseId int64, unionId string, activityId int64) (data micro_scrm_enterprise.OperationActivityUserModel, err error) {
	err = s.Scopes(service.ScopeCond_EQDelete).
		Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityUserModel{}, interpriseId)).
		Where("union_id = ? and activity_id = ?", unionId, activityId).Find(&data).Error
	return
}

func (s ActivityService) ExportUserList(interpriseId int64, ActivityId int64, startTime int64, endTime int64, status int64) ([]micro_scrm_enterprise.OperationActivityUserModel, error) {
	var data []micro_scrm_enterprise.OperationActivityUserModel
	scope := s.Scopes(service.ScopeCond_EQDelete).
		Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityUserModel{}, interpriseId))
	if startTime > 0 && endTime > 0 && endTime > startTime {
		scope.Where("created_at > ? and created_at < ?", startTime, endTime)
	}
	if status > 0 {
		scope.Where("status = ?", status)
	}

	err := scope.Where(map[string]interface{}{
		"activity_id": ActivityId,
	}).Find(&data).Error
	return data, err
}

func (s ActivityService) GetUserList(interpriseId int64, ActivityId int64, startTime int64, endTime int64, status int64, page int, size int) (int64, []micro_scrm_enterprise.OperationActivityUserListModel, error) {
	var data []micro_scrm_enterprise.OperationActivityUserListModel
	var total int64
	if page < 0 {
		page = 1
	}

	if size < 0 {
		size = 10
	}
	scope := s.Scopes(service.ScopeCond_EQDelete).
		Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityUserModel{}, interpriseId))
	if startTime > 0 && endTime > 0 && endTime > startTime {
		scope.Where("created_at > ? and created_at < ?", startTime, endTime)
	}
	if status > 0 {
		scope.Where("status = ?", status)
	}
	err := scope.Where(map[string]interface{}{
		"activity_id": ActivityId,
	}).
		Count(&total).
		Order("created_at desc").
		Limit(size).
		Offset((page - 1) * size).
		Find(&data).Error
	return total, data, err
}

func (s ActivityService) GetUserInfo(enterpriseId int64, activityUserId int64) (micro_scrm_enterprise.OperationActivityUserModel, error) {
	var data micro_scrm_enterprise.OperationActivityUserModel
	err := s.Scopes(service.ScopeCond_EQDelete).
		Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityUserModel{}, enterpriseId)).
		Where("id = ?", activityUserId).Find(&data).Error
	return data, err
}

func (s ActivityService) VerifyUser(enterpriseId int64, activityUserId int64) error {
	err := s.Scopes(service.ScopeCond_EQDelete).
		Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityUserModel{}, enterpriseId)).
		Where("id = ?", activityUserId).Update("status", 2).Error
	return err
}
