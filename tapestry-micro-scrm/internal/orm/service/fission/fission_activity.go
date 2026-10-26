package fission

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm/fission"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

// 获取活动的详情
func (f *FissionService) GetFissionDetail(enterpriseID, activityID int64) (data fission.FissionActivityModel, err error) {
	cond := map[string]interface{}{
		"enterprise_id": enterpriseID,
		"id":            activityID,
	}

	err = f.Where(cond).Preload(clause.Associations).Find(&data).Error

	return
}

// 获取参与裂变活动人员的unionID,但不是品牌VIP的客户 （关注了公众号/小程序后才有）
func (f *FissionService) GetFissionAssistantUserOneFieldList(enterpriseID, activityID, fissionUserID int64) (dataList []UpdateVipFlag, err error) {
	err = f.Table(fmt.Sprintf("%s AS %s", fission.FissionAssistantRecordsModel{}.TableName(), "A")).
		Select("A.`id` AS id, B.`unionid` AS union_id").
		Joins(fmt.Sprintf("INNER JOIN %s%d.%s AS B ON B.`id` = A.`assistant_user_id`", model.BranchModel{}.DBPrefix(), enterpriseID, micro_scrm_enterprise.WxUserCorpExternalModel{}.TableName())).
		Joins(fmt.Sprintf("INNER JOIN %s AS FU ON FU.`id` = A.`fission_user_id`", fission.FissionUserModel{}.TableName())).
		Where("A.`enterprise_id` = ?", enterpriseID).
		Where("A.`activity_id` = ?", activityID).
		Where("FU.`user_id` = ?", fissionUserID).
		Where("A.is_vip = ?", 0).
		Scan(&dataList).Error

	return
}

// 更新裂变客户is_vip字段
func (f *FissionService) UpdateFissionIsVipField(enterpriseID int64, idList []int64) error {
	return f.Model(&fission.FissionAssistantRecordsModel{}).
		Where("enterprise_id = ?", enterpriseID).
		Where("id IN ?", idList).
		Update("is_vip", fission.FAR_IsVip_Yes).Error
}

func (f *FissionService) GetFissionInfoOne(cond fission.FissionActivityModel) (data fission.FissionActivityModel, err error) {
	err = f.Where(&cond).First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return data, nil
	}

	return
}
