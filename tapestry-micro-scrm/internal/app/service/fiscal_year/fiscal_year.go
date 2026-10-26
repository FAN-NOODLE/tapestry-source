package fiscalYearSvc

import (
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	excitationModels "micro-scrm/internal/app/models/excitation"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
)

type FiscalYear struct {
	service.Base

	dao *microWorkDao.FiscalMonth
}

func NewFiscalYear() *FiscalYear {
	return &FiscalYear{
		dao: microWorkDao.NewFiscalMonth(),
	}
}

// JobIsRun
// timeType 1月,2季度,3年
// month 财年的月份
// startTime 财年的开始时间
func (s *FiscalYear) JobIsRun(timeType int32, formatResult map[int64]microWorkModel.FiscalMonthModel, fiscalMonth microWorkModel.FiscalMonthModel, nowTime int64) (int64, int64, bool) {
	if timeType == 0 {
		return 0, 0, false
	}
	if fiscalMonth.EndTime-nowTime >= 60*60*24 {
		return 0, 0, false
	}
	// 季度开始结束
	var startTime, endTime int64
	// 判断 月,季度,年
	switch timeType {
	case excitationModels.MedalRuleTimeTypeMonth:
		startTime = fiscalMonth.StartTime
		endTime = fiscalMonth.EndTime
		// 取当前时间是否为月初第一天
	case excitationModels.MedalRuleTimeTypeQuarter:
		// 取当前时间是否为季度第一天
		if !helper.InArray(fiscalMonth.Period, []int64{3, 6, 9, 12}) {
			return 0, 0, false
		}
		endTime = fiscalMonth.EndTime
		startTime = formatResult[fiscalMonth.Period-2].StartTime
	case excitationModels.MedalRuleTimeTypeYear:
		// 取当前时间是否为新的一年第一天
		if fiscalMonth.Period != 12 {
			return 0, 0, false
		}
		endTime = fiscalMonth.EndTime
		startTime = formatResult[1].StartTime

	}

	return startTime, endTime, true
}
