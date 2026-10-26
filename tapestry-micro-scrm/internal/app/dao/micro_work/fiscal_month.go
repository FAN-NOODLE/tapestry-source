package microWorkDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationModels "micro-scrm/internal/app/models/excitation"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"time"
)

type FiscalMonth struct {
	dao       dao.Dao
	tableName string
}

func NewFiscalMonth() *FiscalMonth {
	return &FiscalMonth{
		dao:       dao.Dao{DB: gormx.GetWorkInstance()},
		tableName: microWorkModel.FiscalMonthModel{}.TableName(),
	}
}

func (d *FiscalMonth) GetByYear(c app.Ctx, enterpriseId int64, year int64) (result []microWorkModel.FiscalMonthModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where(map[string]any{
		"enterprise_id": enterpriseId,
		"fiscal_year":   year,
	}).Find(&result).Error
	return
}

func (d *FiscalMonth) GetDataByYearAndTime(c app.Ctx, enterpriseId int64, time int64) (result microWorkModel.FiscalMonthModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ?", enterpriseId).
		Where("start_time < ? AND end_time > ?", time, time).
		Find(&result).Error
	return
}

func (d *FiscalMonth) GetAllMonth(c app.Ctx, enterpriseId int64) (result []microWorkModel.FiscalMonthModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ?", enterpriseId).
		Find(&result).Error
	return
}

// GetFiscalMonth 根据时间戳获取所在的财月
func (d *FiscalMonth) GetFiscalMonth(c app.Ctx, enterpriseId int64, fiscalTime int64) (info microWorkModel.FiscalMonthModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id=?", enterpriseId).
		Where("start_time < ? and end_time > ?", fiscalTime, fiscalTime).
		First(&info).Error
	return
}

// GetFiscalMonthByEndTime 根据end_time获取所在的财月
func (d *FiscalMonth) GetFiscalMonthByEndTime(c app.Ctx, enterpriseId int64, fiscalTime int64) (info microWorkModel.FiscalMonthModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id=?", enterpriseId).Where(" end_time=?", fiscalTime).Find(&info).Error
	return
}

// GetFiscalYear 根据时间戳获取所在财年的开始时间和结束时间
func (d *FiscalMonth) GetFiscalYear(c app.Ctx, enterpriseId int64, fiscalTime int64) (startTime int64, endTime int64, err error) {
	month, err := d.GetFiscalMonth(c, enterpriseId, fiscalTime)
	var result []microWorkModel.FiscalMonthModel
	err = d.dao.WithContext(c).Table(d.tableName).Where(map[string]any{
		"enterprise_id": enterpriseId,
		"fiscal_year":   month.Year,
	}).Where("Period in ?", []int64{1, 12}).Find(&result).Error
	for _, v := range result {
		if v.Period == 1 {
			startTime = v.StartTime
		}
		if v.Period == 12 {
			endTime = v.EndTime
		}
	}
	return
}

// GetTimeByAstrictType 根据积分上限类型获取开始结束时间
func (d *FiscalMonth) GetTimeByAstrictType(c ctx.Context, enterpriseId int64, astrictType int64, nowTime time.Time) (startTime int64, endTime int64, err error) {
	// 如果该维度有积分上限限制
	switch astrictType {
	//当日上限
	case excitationModels.Astrict_Day:
		//这几个改为当日的上限（2023.4.6）,
		startTime = helper.GetZeroTimeByTime(nowTime).Unix()
		endTime = startTime + 86400
	case excitationModels.Astrict_Week:
		endTime = helper.GetLastDateOfWeek(helper.GetZeroTimeByTime(nowTime)).Unix()
		startTime = endTime - 86400*7
	case excitationModels.Astrict_Month:
		// 获取统计范围内财月的时间范围
		month, err := d.GetFiscalMonth(c, enterpriseId, nowTime.Unix())
		if err != nil {
			return 0, 0, err
		}
		startTime = month.StartTime
		endTime = month.EndTime
	case excitationModels.Astrict_Year:
		// 获取统计范围内财月的时间范围
		startTime, endTime, _ = d.GetFiscalYear(c, enterpriseId, nowTime.Unix())
	}
	return
}

// GetFiscalMonthByPeriod 根据时间戳获取所在的财月
func (d *FiscalMonth) GetFiscalMonthByPeriod(c app.Ctx, enterpriseId, period, fiscalYear int64) (info microWorkModel.FiscalMonthModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id=?", enterpriseId).
		Where("period = ?", period).
		Where("fiscal_year = ?", fiscalYear).
		First(&info).Error
	return
}
