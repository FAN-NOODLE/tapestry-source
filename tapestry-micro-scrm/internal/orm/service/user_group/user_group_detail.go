package user_group

import (
	"micro-scrm/internal/orm/model/micro_scrm"
	"time"
)

// 根据客户分组ID，删除客户分组详情
func (u *UserGroupService) DelCustomerGroupDetailByCustomerGroupID(enterpriseID, customerGroupID int64) error {
	return u.Model(&micro_scrm.UserGroupDetailModel{}).
		Where("enterprise_id = ?", enterpriseID).
		Where("user_group_id = ?", customerGroupID).
		Updates(map[string]interface{}{
			"updated_at": time.Now().Unix(),
			"deleted_at": time.Now().Unix(),
		}).Error
}


func (u *UserGroupService) UpdateDetailISExtract(enterpriseID, customerGroupID int64) error {
	return u.Model(&micro_scrm.UserGroupDetailModel{}).
		Where("enterprise_id = ?", enterpriseID).
		Where("user_group_id = ?", customerGroupID).
		Update("is_extract", 1).Error
}