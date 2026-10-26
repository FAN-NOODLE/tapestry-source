package activity_campaign

import (
	"gorm.io/gorm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/service"
)

func (a *ActivityCampaignService) ActivitieInfo(enterPriseID, activityID int64) (result micro_scrm_enterprise.OperationActivityModel, err error) {
	err = a.Scopes(model.BranchTable(&result, enterPriseID), service.ScopeCond_EQDelete).
		Where("id=?", activityID).Find(&result).Error

	return
}

func (a *ActivityCampaignService) GetUserSignInDate(enterPriseID, activitieId int64, unionId string) (result []string) {
	a.Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivitySignInLog{}, enterPriseID), service.ScopeCond_EQDelete).
		Where("activity_id=? AND union_id=?", activitieId, unionId).Select("FROM_UNIXTIME(created_at,'%Y-%m-%d') sign_in_date").Order("id").Find(&result)
	return
}
func (s *ActivityCampaignService) InsertactivityUser(interpriseId int64, user micro_scrm_enterprise.UserModel) (int64, error) {
	err := s.Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityUserModel{}, interpriseId)).Create(&user).Error
	return user.ID, err
}
func (a *ActivityCampaignService) ActivitieAddNum(enterPriseID, activityID int64) error {
	return a.Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityModel{}, enterPriseID), service.ScopeCond_EQDelete).
		Where("id=?", activityID).Update("nums", gorm.Expr("nums + 1")).Error

}
func (a *ActivityCampaignService) Findactivitysigninlog(enterPriseID, activitieId, startTiem, endTime int64, unionId string) (result micro_scrm_enterprise.OperationActivitySignInLog, err error) {
	err = a.Scopes(model.BranchTable(&result, enterPriseID), service.ScopeCond_EQDelete).
		Where("activity_id=? AND union_id=? AND created_at BETWEEN ? AND ?", activitieId, unionId, startTiem, endTime).Find(&result).Error
	return
}

func (a *ActivityCampaignService) Createactivitysigninlog(enterPriseID int64, result micro_scrm_enterprise.OperationActivitySignInLog) error {
	return a.Scopes(model.BranchTable(&result, enterPriseID)).Create(&result).Error
}

func (a *ActivityCampaignService) FindUserSignup(enterPriseID, activityId int64, unionId string) (data micro_scrm_enterprise.OperationActivityUserModel, err error) {
	err = a.Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityUserModel{}, enterPriseID), service.ScopeCond_EQDelete).
		Where("activity_id=? AND union_id = ? AND status=?", activityId, unionId, 2).Find(&data).Error
	return
}

func (a *ActivityCampaignService) FindActivityUserinfo(enterPriseID, activityId int64, unionId string) (data micro_scrm_enterprise.OperationActivityUserModel, err error) {
	err = a.Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityUserModel{}, enterPriseID), service.ScopeCond_EQDelete).
		Where("activity_id=? AND union_id = ?", activityId, unionId).Find(&data).Error
	return
}
