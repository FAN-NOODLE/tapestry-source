package marketing

import (
	"micro-scrm/internal/orm/model/micro_scrm/marketing"
)

// 营销策略执行的结果
func (m *MarketingService) GetStrategyResult(param MarketingResultQueryParam) (dataList []StrategyExecResult, err error) {
	cond := map[string]interface{}{
		"enterprise_id": param.EnterpriseID,
		"strategy_id":   param.StrategyID,
	}

	field := []string{
		"channel_type",
		"COUNT(`id`) AS total",
		"SUM(IF(`status` = 2, 1, 0)) AS num_success",
		"SUM(IF(`status` = 2, 0, 1)) AS num_fail",
	}

	err = m.Model(&marketing.MarketingStrategyRecordsModel{}).
		Select(field).
		Where(cond).
		Where("`created_at` BETWEEN ? AND ?", param.StartTime, param.EndTime).
		Group("channel_type").
		Find(&dataList).Error

	return
}
