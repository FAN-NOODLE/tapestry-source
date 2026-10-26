package microWorkModel

import "micro-scrm/internal/app/models"

type FiscalMonthModel struct {
	models.MicroBaseModel
	ID           uint64 `json:"id" gorm:"primaryKey;column:id;type:bigint(20) unsigned;comment:主键ID"`
	EnterpriseID uint64 `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20) unsigned;not null;default:0;comment:企业ID"`
	Year         string `json:"year" gorm:"column:year;type:varchar(50);not null;comment:财年标识"`
	Period       int64  `json:"period" gorm:"column:period;type:int(10) unsigned;not null;default:0;comment:财月自增月份周期1-12"`
	StartTime    int64  `json:"start_time" gorm:"column:start_time;type:int(10) unsigned;not null;default:0;comment:开始时间"`
	EndTime      int64  `json:"end_time" gorm:"column:end_time;type:int(10) unsigned;not null;default:0;comment:开始时间"`
	FiscalYear   int64  `json:"fiscal_year" gorm:"column:fiscal_year;type:int(10) unsigned;not null;default:0;comment:所属的财年"`
}

// TableName returns the table name of the FiscalMonth model
func (FiscalMonthModel) TableName() string {
	return "fiscal_month"
}
