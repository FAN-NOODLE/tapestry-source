package marketing

import (
	"micro-scrm/internal/orm/model/micro_scrm/marketing"
	"micro-scrm/internal/orm/service"
)

// @title 				IsEditByUserGroupID
// @description 		根据客户分组ID，判断该客户分组是否可以进行编辑、删除
// @param userGroupID 	int64 	客户分组ID
// @return isExist 		bool 	是否可以编辑、删除；true 可以，否则不可以
// @return err 			error 	错误
func (m *MarketingService) IsEditByUserGroupID(userGroupID int64) (bool, error) {
	var count int64
	if err := m.Model(&marketing.MarketingStrategiesModel{}).
		Scopes(service.ScopeCond_EQDelete).
		Where("customer_group_id = ?", userGroupID).
		Where("strategy_status IN ?", []int{StrategyStatusNotStart, StrategyStatusRunning, StrategyStatusSuspend}).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count <= 0, nil
}

// 获取营销策略信息详情
func (m *MarketingService) GetMarketingStrategyDetail(param MarketingExportParam) (data marketing.MarketingStrategiesModel, err error) {
	err = m.Scopes(service.ScopeCond_EQDelete).
		Where("enterprise_id = ?", param.EnterpriseID).
		Where("id = ?", param.ID).
		Preload("MarketingStrategyRecords", "created_at BETWEEN ? AND ?", param.StartTime, param.EndTime).
		Preload("MarketingStrategyTriggerCount").
		First(&data).Error

	return
}
