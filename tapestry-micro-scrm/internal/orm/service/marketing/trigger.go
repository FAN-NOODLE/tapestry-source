package marketing

import (
	"micro-scrm/internal/orm/model/micro_scrm/marketing"
)

// RecordTriggerCount 记录策略触发的人数
func (m *MarketingService) RecordTriggerCount(model marketing.MarketingStrategyTriggerCountModel) error {
	result := m.Create(&model)
	return result.Error
}

func (m *MarketingService) UpdateRecordTriggerCount(strategyId, count int64) error {
	err := m.Model(&marketing.MarketingStrategyTriggerCountModel{}).
		Where("strategy_id=?", strategyId).Update("count", count).Error
	return err
}
