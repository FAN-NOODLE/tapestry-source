package user_group

import (
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/service"
	"time"
)

// 根据客户分组ID，删除客户分组
func (u *UserGroupService) DelCustomerGroupByID(enterpriseID, customerGroupID int64) error {
	return u.Model(&micro_scrm.UserGroupModel{}).
		Where("enterprise_id = ?", enterpriseID).
		Where("id = ?", customerGroupID).
		Updates(map[string]interface{}{
			"updated_at": time.Now().Unix(),
			"deleted_at": time.Now().Unix(),
		}).Error
}

// 根据标签ID，获取包含该标签的所有客户分组的标签
func (u *UserGroupService) GetUserGroupIncludeTagInfoByTagID(enterpriseID int64, tagIDList []int64) []entity.UserGroupDetailTagInfo {
	// 获取标签所在的客户分组
	userGroupIDList, _ := u.GetUserMaybeCustomerGroup(enterpriseID, tagIDList)

	// 根据客户分组ID，获取客户分组的
	userGroupDetail := make([]entity.UserGroupDetailTagInfo, 0)
	u.Table("`user_group_detail` AS UGD").
		Select("UG.`name` AS user_group_name, UG.`id` AS user_group_id, TG.`id` AS tag_group_id, TG.`group_name` AS tag_group_name, T.`id` AS tag_id, T.`name` AS tag_name, UGD.`include_type`").
		Joins("LEFT JOIN `user_group` UG ON UG.`id` = UGD.`user_group_id`").
		Joins("LEFT JOIN `work_external_tag_group` TG ON TG.`id` = UGD.`tag_group_id`").
		Joins("LEFT JOIN `work_external_tag` T ON T.`id` = UGD.`tag_id`").
		Where(map[string]interface{}{
			"UGD.enterprise_id": enterpriseID,
			"UGD.user_group_id": userGroupIDList,
		}).Scan(&userGroupDetail)

	return userGroupDetail
}

// 根据客户分组ID，获取客户分组详情
func (u *UserGroupService) GetCustomerGroupDetail(enterpriseID, customerGroupID int64) (data []micro_scrm.UserGroupDetailModel, err error) {
	where := map[string]interface{}{
		"enterprise_id": enterpriseID,
		"user_group_id": customerGroupID,
	}

	err = u.Model(&data).
		Scopes(service.ScopeCond_EQDelete).
		Where(where).
		Preload("UserGroup").Preload("TagGroup").Preload("Tag").
		Find(&data).
		Error

	return
}

// 根据标签ID，获取用户可能存在的客户分组ID
func (u *UserGroupService) GetUserMaybeCustomerGroup(enterpriseID int64, tagIDList []int64) (customerGroupIDList []int64, err error) {
	err = u.Model(&micro_scrm.UserGroupDetailModel{}).
		Scopes(service.ScopeCond_EQDelete).
		Distinct("user_group_id").
		Where(map[string]interface{}{
			"enterprise_id": enterpriseID,
			"tag_id":        tagIDList,
			"include_type":  micro_scrm.UGD_IncludeType_Include,
		}).Pluck("user_group_id", &customerGroupIDList).Error

	return
}



