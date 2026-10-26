package micro_scrm_enterprise

import (
	"micro-scrm/internal/orm/model"
)

type BaMonthlyKanbanModel struct {
	model.BranchModel

	Month                                  string  `gorm:"column:month;NOT NULL" json:"month"`                                                                                   // 月份
	BaName                                 string  `gorm:"column:ba_name;NOT NULL" json:"ba_name"`                                                                               // BA姓名
	BaMarscode                             string  `gorm:"column:ba_marscode" json:"ba_marscode"`                                                                                // BA marscode
	Mobile                                 string  `gorm:"column:mobile" json:"mobile"`                                                                                          // 手机
	CounterName                            string  `gorm:"column:counter_name" json:"counter_name"`                                                                              // 柜台名称
	SalesType                              int     `gorm:"column:sales_type;default:1;NOT NULL" json:"sales_type"`                                                               // 业绩统计类型 1BA、2 店经理&柜长、3 分区主管、4 大区主管
	Sales                                  float64 `gorm:"column:sales;default:0.00;NOT NULL" json:"sales"`                                                                      // 月度销售额
	MonthlySales                           float64 `gorm:"column:monthly_sales;default:0.00;NOT NULL" json:"monthly_sales"`                                                      // 月度目标销售额
	FinishRate                             int     `gorm:"column:finish_rate;default:0" json:"finish_rate"`                                                                      // 月度销售额完成率(销售额 / 完成率 *100%)
	GrowthRate                             int     `gorm:"column:growth_rate;default:0" json:"growth_rate"`                                                                      // 同比增长率
	DifferenceSales                        float64 `gorm:"column:difference_sales;default:0.00" json:"difference_sales"`                                                         // 差额(目标金额 - 销售额)
	OrderNum                               int     `gorm:"column:order_num;default:0;NOT NULL" json:"order_num"`                                                                 // 订单数据量
	OrderNumPrice                          float64 `gorm:"column:order_num_price;default:0.00" json:"order_num_price"`                                                           // 客单价
	CustomerNum                            int     `gorm:"column:customer_num;default:0;NOT NULL" json:"customer_num"`                                                           // 客人数据量
	NewCustomerNum                         int     `gorm:"column:new_customer_num;default:0;NOT NULL" json:"new_customer_num"`                                                   // 新客人数据量
	NewCustomerSales                       float64 `gorm:"column:new_customer_sales;default:0.00;NOT NULL" json:"new_customer_sales"`                                            // 新客人月度销售额
	NewCustomerRate                        int     `gorm:"column:new_customer_rate;default:0;NOT NULL" json:"new_customer_rate"`                                                 // 新客户销售额占总销售额(新二购销售额 / 总销售额)
	VisitCustomerNum                       int     `gorm:"column:visit_customer_num;default:0;NOT NULL" json:"visit_customer_num"`                                               // 游客人数据量
	VisitCustomerSales                     float64 `gorm:"column:visit_customer_sales;default:0;NOT NULL" json:"visit_customer_sales"`                                           // 游客月度销售额
	VisitCustomerSalesRate                 float64 `gorm:"column:visit_customer_sales_rate;default:0" json:"visit_customer_sales_rate"`                                          // 游客销售额占总销售额占比
	NewSecondCustomerNum                   int     `gorm:"column:new_second_customer_num;default:0;NOT NULL" json:"new_second_customer_num"`                                     // 新二购客人数据量
	NewSecondCustomerYNum                  int     `gorm:"column:new_second_customer_y_num;default:0;NOT NULL" json:"new_second_customer_y_num"`                                 // 新二购客人数据量
	NewSecondCustomerSales                 float64 `gorm:"column:new_second_customer_sales;default:0.00;NOT NULL" json:"new_second_customer_sales"`                              // 新二购月度销售额
	NewSecondCustomerSalesRate             int     `gorm:"column:new_second_customer_sales_rate;default:0" json:"new_second_customer_sales_rate"`                                // 新二购销售额占总销售额(新二购销售额 / 总销售额)
	NewSecondAndThreeScan                  int     `gorm:"column:new_second_and_three_scan;default:0;NOT NULL" json:"new_second_and_three_scan"`                                 // 新二购且一码三观客人数据量
	MonthlySecondBuyAndThreeScan           int     `gorm:"column:monthly_second_buy_and_three_scan;default:0;NOT NULL" json:"monthly_second_buy_and_three_scan"`                 // 新二购且一码三观客人指标数据量
	NewSecondAndThreeScanRate              int     `gorm:"column:new_second_and_three_scan_rate;default:0" json:"new_second_and_three_scan_rate"`                                // 新二购一码三关注率( 新二购且一码三关注的完成人数 / 新二购人数 * 100%)
	NewSecondAndThreeScanWageProgress      float64 `gorm:"column:new_second_and_three_scan_wage_progress;default:0.00" json:"new_second_and_three_scan_wage_progress"`           // 新二购工资完成进度(新二购 且 一码三关注 的完成 人数 / 新二购 且 一码三关注 的指标 人数 * 100%)
	OldCustomerSales                       float64 `gorm:"column:old_customer_sales;default:0.00;NOT NULL" json:"old_customer_sales"`                                            // 老客月度销售额
	OldCustomerNum                         int     `gorm:"column:old_customer_num;default:0;NOT NULL" json:"old_customer_num"`                                                   // 老客人数据量
	OldCustomerSalesRate                   int     `gorm:"column:old_customer_sales_rate;default:0" json:"old_customer_sales_rate"`                                              // 老客销售额占总销售额
	OldAndThreeScan                        int     `gorm:"column:old_and_three_scan;default:0;NOT NULL" json:"old_and_three_scan"`                                               // 老客人一码三观数据量
	MonthlyOldCustomerAndThreeScan         int     `gorm:"column:monthly_old_customer_and_three_scan;default:0;NOT NULL" json:"monthly_old_customer_and_three_scan"`             // 老客且一码三关注的指标人数
	OldXCustomerNum                        int     `gorm:"column:old_x_customer_num;default:0" json:"old_x_customer_num"`                                                        // 未知老客数
	MonthlyOldCustomerAndThreeScanRate     int     `gorm:"column:monthly_old_customer_and_three_scan_rate;default:0" json:"monthly_old_customer_and_three_scan_rate"`            // 老客一码三关注率(老客且一码三关注的完成人数 / 老客人数 * 100%)
	MonthlyOldCustomerAndThreeWageProgress float64 `gorm:"column:monthly_old_customer_and_three_wage_progress;default:0.00" json:"monthly_old_customer_and_three_wage_progress"` // 老客工资完成进度(老客且 一码三关注 的完成 人数 / 老客 且 一码三关注 的指标 人数 * 100%)
	NewCustomerAndThreeScan                int     `gorm:"column:new_customer_and_three_scan;default:0;NOT NULL" json:"new_customer_and_three_scan"`                             // 新客且一码三观客人数据量
	MonthlyNewCustomerAndThreeScan         int     `gorm:"column:monthly_new_customer_and_three_scan;default:0;NOT NULL" json:"monthly_new_customer_and_three_scan"`             // 新客且一码三关注的目标人数
	NewCustomerAndThreeScanRate            int     `gorm:"column:new_customer_and_three_scan_rate;default:0" json:"new_customer_and_three_scan_rate"`                            // 新客一码三关注率(新客且一码三关注的完成人数 / 新客人数 * 100%)
}

func (BaMonthlyKanbanModel) TableName() string {
	return "ba_monthly_kanban"
}

type BAMonthlyKanBanResult struct {
	Name string
	// 销售数据
	Sales        float32 // 销售额
	TargetSales  float32 // 销售目标
	FinishRate   int8    // 销售完成率
	SurplusSales float32 // 销售差额
	SurplusRate  int8    // 销售差额率
	OrderNum     int     // 订单数
	OrderPrice   float32 // 客单价

	CustomerNum          int     // 客人数量
	NewCustomerNum       int     // 新客数量
	NewCustomerSales     float32 // 新客月度销售额
	NewCustomerSalesRate int8    // 新客销售占比

	OldCustomerNum       int     // 老客数量
	OldCustomerSales     float32 // 老客月度销售额
	OldCustomerSalesRate int8    // 老客销售占比

	NewSecondCustomerNum       int     // 新二购数量
	NewSecondCustomerSales     float32 // 新二购月度销售额
	NewSecondCustomerSalesRate int8    // 新二购销售占比

	VisitCustomerNum       int     // 老客数量
	VisitCustomerSales     float32 // 老客月度销售额
	VisitCustomerSalesRate int8    // 老客销售占比

	// 运营数据
	MonthlySecondBuyAndThreeScanNum int // 新二购指标人数
	NewSecondAndThreeScanFinishRate int // 新二购指标完成率
	NewSecondAndThreeScanNum        int // 新二购一码三关注人数
	NewSecondAndThreeScanRate       int // 新二购一码三关注率

	MonthlyOldCustomerAndThreeScanNum int // 老客指标人数
	OldAndThreeScanFinishRate         int // 老客一码三关注人数
	OldAndThreeScanNum                int // 老顾客一码三关注人数
	OldAndThreeScanRate               int // 老顾客一码三关注率

	MonthlyNewCustomerAndThreeScanNum int // 新顾客指标人数
	NewCustomerAndThreeScanFinishRate int // 新顾客一码三关注人数
	NewCustomerAndThreeScanNum        int // 新顾客一码三关注人数
	NewCustomerAndThreeScanRate       int // 新顾客一码三关注率
}
