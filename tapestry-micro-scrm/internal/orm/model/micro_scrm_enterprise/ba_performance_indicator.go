package micro_scrm_enterprise

import (
	"micro-scrm/internal/orm/model"
)

type BaPerformanceIndicatorModel struct {
	model.BranchModel

	Regional                       string `json:"regional" gorm:"column:regional;type:varchar(50);not null;comment:区域（大区）"`                                                                       // 区域（大区）
	Area                           string `json:"area" gorm:"column:area;type:varchar(50);not null;comment:区域"`                                                                                   // 区域
	City                           string `json:"city" gorm:"column:city;type:varchar(50);not null;comment:城市"`                                                                                   // 城市
	CounterCode                    string `json:"counter_code" gorm:"column:counter_code;type:varchar(50);not null;comment:柜台code"`                                                               // 柜台code
	CounterName                    string `json:"counter_name" gorm:"column:counter_name;type:varchar(50);not null;comment:柜台名称"`                                                                 // 柜台名称
	BaName                         string `json:"ba_name" gorm:"column:ba_name;type:varchar(50);not null;comment:BA姓名"`                                                                           // BA姓名
	Mobile                         string `json:"mobile" gorm:"column:mobile;type:varchar(20);comment:手机"`                                                                                        // 手机
	MonthlySales                   int64  `json:"monthly_sales" gorm:"column:monthly_sales;type:bigint(20);not null;default:0;comment:月度销售目标（元）"`                                                 // 月度销售目标（元）
	MonthlySecondBuyAndThreeScan   int    `json:"monthly_second_buy_and_three_scan" gorm:"column:monthly_second_buy_and_three_scan;type:int(10);not null;default:0;comment:新二购且一码三关注月度目标（人）"`     // 新二购且一码三关注月度目标（人）
	MonthlyOldCustomerAndThreeScan int    `json:"monthly_old_customer_and_three_scan" gorm:"column:monthly_old_customer_and_three_scan;type:int(10);not null;default:0;comment:老客购且一码三关注月度目标（人）"` // 老客购且一码三关注月度目标（人）
	MonthlyNewCustomerAndThreeScan int    `json:"monthly_new_customer_and_three_scan" gorm:"column:monthly_new_customer_and_three_scan;type:int(10);not null;default:0;comment:新客购且一码三关注月度目标（人）"` // 新客购且一码三关注月度目标（人）
	OperationUserID                int64  `json:"operation_user_id" gorm:"column:operation_user_id;type:bigint(11);not null;default:0;comment:成员表id"`                                             // 成员表id
	Month                          string `json:"month" gorm:"column:month;type:char(6);not null;comment:导入月份：yyyymm"`                                                                            // 导入月份：yyyymm
}

func (BaPerformanceIndicatorModel) TableName() string {
	return "ba_performance_indicator"
}
