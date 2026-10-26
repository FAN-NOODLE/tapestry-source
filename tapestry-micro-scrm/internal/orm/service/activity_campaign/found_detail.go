package activity_campaign

import (
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

func (a *ActivityCampaignService) GetSignInList(enterpriseID, activityID, startTime, endTime int64, pageSize, pageNum int, typeInt int) (result []*micro_scrm_enterprise.SignInList, count int64, err error) {
	var user = micro_scrm_enterprise.OperationActivityUsers{}.AllTableName(enterpriseID)
	tx := a.Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivitySignInLog{}, enterpriseID)).
		Select("operation_activity_sign_in_log.union_id,count(operation_activity_sign_in_log.union_id) AS count,"+
			"MAX(operation_activity_sign_in_log.created_at) AS last_sign_time,user.receiver_name AS name,user.receiver_phone AS phone").
		Joins("LEFT JOIN "+user+"user ON user.id = operation_activity_sign_in_log.activity_user_id AND user.deleted_at = 0").
		Where("operation_activity_sign_in_log.activity_id= ? "+
			"AND operation_activity_sign_in_log.deleted_at = 0", activityID)
	if startTime != 0 && endTime != 0 {
		tx = tx.Where("operation_activity_sign_in_log.created_at <= ? AND operation_activity_sign_in_log.created_at >= ?", endTime, startTime)
	}
	if typeInt == 1 {
		tx = tx.Group("operation_activity_sign_in_log.union_id").Order("MAX(operation_activity_sign_in_log.created_at) DESC")
		err = tx.Find(&result).Error
		return
	}
	tx = tx.Group("operation_activity_sign_in_log.union_id").Order("MAX(operation_activity_sign_in_log.created_at) DESC").Count(&count)
	tx = tx.Limit(pageSize).Offset(pageSize * (pageNum - 1)).Find(&result)
	err = tx.Error
	return
}

func (a *ActivityCampaignService) UserInfo(enterpriseID, activityID int64, unionID string) (result micro_scrm_enterprise.UserInfo, err error) {
	err = a.Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivityUserModel{}, enterpriseID)).
		Select("wx_name AS name,wx_phone AS phone,receiver_name AS shr_name,receiver_phone AS shr_phone,"+
			"receiver_address AS shr_address").
		Where("activity_id = ? AND union_id = ? AND deleted_at = 0 ", activityID, unionID).
		Find(&result).Error
	return
}

func (a *ActivityCampaignService) GetAvatarInfo(enterpriseID int64, unionID string) (result micro_scrm_enterprise.AvatarInfo, err error) {
	err = a.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, enterpriseID)).
		Select("avatar").
		Where("unionid = ? AND deleted_at = 0", unionID).Find(&result).Error
	return
}

func (a *ActivityCampaignService) GetSignDetail(enterpriseID, activityID, startTime, endTime int64, unionID string, pageSize, pageNum int) (result []micro_scrm_enterprise.SignInInfoDetail, err error) {
	tx := a.Scopes(model.BranchTable(&micro_scrm_enterprise.OperationActivitySignInLog{}, enterpriseID)).
		Select("created_at AS sign_time,sign_in_content AS sign_content,sign_in_pic AS sign_pic").
		Where("union_id = ? AND activity_id= ? AND deleted_at = 0", unionID, activityID)
	if startTime != 0 && endTime != 0 {
		tx = tx.Where("created_at <= ? AND created_at >= ?")
	}
	tx = tx.Limit(pageSize).Offset(pageSize * (pageNum - 1)).Find(&result)
	err = tx.Error
	return
}
