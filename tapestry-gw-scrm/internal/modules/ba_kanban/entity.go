package ba_kanban

// SalesDataDetailResponse 销售数据明细 response
type SalesDataDetailResponse struct {
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
	List     []SalesDataDetailList `json:"list"`
}

type SalesDataDetailList struct {
	Name string `json:"name"`
	// 销售数据
	Sales        float32 `json:"sales"`         // 销售额
	TargetSales  float32 `json:"target_sales"`  // 销售目标
	FinishRate   int8    `json:"finish_rate"`   // 销售完成率
	SurplusSales float32 `json:"surplus_sales"` // 销售差额
	OrderNum     int     `json:"order_num"`     // 订单数
	OrderPrice   float32 `json:"order_price"`   // 客单价

	CustomerNum          int     `json:"customer_num"`            // 客人数量
	NewCustomerNum       int     `json:"new_customer_num"`        // 新客数量
	NewCustomerSales     float32 `json:"new_customer_sales"`      // 新客月度销售额
	NewCustomerSalesRate int8    `json:"new_customer_sales_rate"` // 新客销售占比

	OldCustomerNum       int     `json:"old_customer_num"`        // 老客数量
	OldCustomerSales     float32 `json:"old_customer_sales"`      // 老客月度销售额
	OldCustomerSalesRate int8    `json:"old_customer_sales_rate"` // 老客销售占比

	NewSecondCustomerNum       int     `json:"new_second_customer_num"`        // 新二购数量
	NewSecondCustomerSales     float32 `json:"new_second_customer_sales"`      // 新二购月度销售额
	NewSecondCustomerSalesRate int8    `json:"new_second_customer_sales_rate"` // 新二购销售占比

	VisitCustomerNum       int     `json:"visit_customer_num"`        // 老客数量
	VisitCustomerSales     float32 `json:"visit_customer_sales"`      // 老客月度销售额
	VisitCustomerSalesRate int8    `json:"visit_customer_sales_rate"` // 老客销售占比
}

// OperationDataDetailResponse 运营数据明细 response
type OperationDataDetailResponse struct {
	Total    int64                     `json:"total"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"page_size"`
	List     []OperationDataDetailList `json:"list"`
}

type OperationDataDetailList struct {
	Name string `json:"name"`
	// 运营数据
	MonthlySecondBuyAndThreeScanNum int `json:"monthly_second_buy_and_three_scan_num"` // 新二购指标人数
	NewSecondAndThreeScanNum        int `json:"new_second_and_three_scan_num"`         // 新二购一码三关注人数
	NewSecondAndThreeScanRate       int `json:"new_second_and_three_scan_rate"`        // 新二购一码三关注率

	MonthlyOldCustomerAndThreeScanNum int `json:"monthly_old_customer_and_three_scan_num"` // 老客指标人数
	OldAndThreeScanNum                int `json:"old_and_three_scan_num"`                  // 老顾客一码三关注人数
	OldAndThreeScanRate               int `json:"old_and_three_scan_rate"`                 // 老顾客一码三关注率

	MonthlyNewCustomerAndThreeScanNum int `json:"monthly_new_customer_and_three_scan_num"` // 新顾客指标人数
	NewCustomerAndThreeScanNum        int `json:"new_customer_and_three_scan_num"`         // 新顾客一码三关注人数
	NewCustomerAndThreeScanRate       int `json:"new_customer_and_three_scan_rate"`        // 新顾客一码三关注率
}

type BaSourceDeptRequst struct {
	EnterpriseId int64  `json:"enterprise_id" binding:"required" desc:"公司编号"`
	Mobile       string `json:"mobile" binding:"required" desc:"当前BA的电话号码"`
}

type BaKanBanDataRequest struct {
	Month              int64  `json:"month" binding:"required" desc:"当前月份"`
	Default_mobile     string `json:"default_mobile" binding:"required" desc:"默认查询电话"`
	Dept_leader_mobile string `json:"dept_leader_mobile" desc:"部门电话"`
	Ba_mobile          string `json:"ba_mobile" desc:"BA电话"`
	EnterpriseId       int64  `json:"enterprise_id" binding:"required" desc:"公司编号"`
}
type BaKanBanpOerationDataRequest struct {
	Month              int64  `json:"month" binding:"required" desc:"当前月份"`
	Default_mobile     string `json:"default_mobile" binding:"required" desc:"默认查询电话"`
	Dept_leader_mobile string `json:"dept_leader_mobile" desc:"部门电话"`
	Ba_mobile          string `json:"ba_mobile" desc:"BA电话"`
	EnterpriseId       int64  `json:"enterprise_id" binding:"required" desc:"公司编号"`
	PageSize           int    `json:"page_size"`
	Page               int    `json:"page"`
}
type OperationDataDetailResponseV1 struct {
	Total    int64                       `json:"total"`
	Page     int                         `json:"page"`
	PageSize int                         `json:"page_size"`
	List     []OperationDataDetailListV1 `json:"list"`
}
type OperationDataDetailListV1 struct {
	Month                                  string  `json:"month"`                                        // 月份
	BaName                                 string  `json:"ba_name"`                                      // BA姓名
	BaMarscode                             string  `json:"ba_marscode"`                                  // BA marscode
	Mobile                                 string  `json:"mobile"`                                       // 手机
	CounterName                            string  `json:"counter_name"`                                 // 柜台名称
	SalesType                              int     `json:"sales_type"`                                   // 业绩统计类型 1BA、2 店经理&柜长、3 分区主管、4 大区主管
	Sales                                  float64 `json:"sales"`                                        // 月度销售额
	MonthlySales                           float64 `json:"monthly_sales"`                                // 月度目标销售额
	FinishRate                             int     `json:"finish_rate"`                                  // 月度销售额完成率(销售额 / 完成率 *100%)
	GrowthRate                             int     `json:"growth_rate"`                                  // 同比增长率
	DifferenceSales                        float64 `json:"difference_sales"`                             // 差额(目标金额 - 销售额)
	OrderNum                               int     `json:"order_num"`                                    // 订单数据量
	OrderNumPrice                          float64 `json:"order_num_price"`                              // 客单价
	CustomerNum                            int     `json:"customer_num"`                                 // 客人数据量
	NewCustomerNum                         int     `json:"new_customer_num"`                             // 新客人数据量
	NewCustomerSales                       float64 `json:"new_customer_sales"`                           // 新客人月度销售额
	NewCustomerRate                        int     `json:"new_customer_rate"`                            // 新客户销售额占总销售额(新二购销售额 / 总销售额)
	VisitCustomerNum                       int     `json:"visit_customer_num"`                           // 游客人数据量
	VisitCustomerSales                     int64   `json:"visit_customer_sales"`                         // 游客月度销售额
	NewSecondCustomerNum                   int     `json:"new_second_customer_num"`                      // 新二购客人数据量
	NewSecondCustomerSales                 float64 `json:"new_second_customer_sales"`                    // 新二购月度销售额
	NewSecondCustomerSalesRate             int     `json:"new_second_customer_sales_rate"`               // 新二购销售额占总销售额(新二购销售额 / 总销售额)
	NewSecondAndThreeScan                  int     `json:"new_second_and_three_scan"`                    // 新二购且一码三观客人数据量
	MonthlySecondBuyAndThreeScan           int     `json:"monthly_second_buy_and_three_scan"`            // 新二购且一码三观客人指标数据量
	NewSecondAndThreeScanRate              int     `json:"new_second_and_three_scan_rate"`               // 新二购一码三关注率( 新二购且一码三关注的完成人数 / 新二购人数 * 100%)
	NewSecondAndThreeScanWageProgress      float64 `json:"new_second_and_three_scan_wage_progress"`      // 新二购工资完成进度(新二购 且 一码三关注 的完成 人数 / 新二购 且 一码三关注 的指标 人数 * 100%)
	OldCustomerSales                       float64 `json:"old_customer_sales"`                           // 老客月度销售额
	OldCustomerNum                         int     `json:"old_customer_num"`                             // 老客人数据量
	OldCustomerSalesRate                   int     `json:"old_customer_sales_rate"`                      // 老客销售额占总销售额
	OldAndThreeScan                        int     `json:"old_and_three_scan"`                           // 老客人一码三观数据量
	MonthlyOldCustomerAndThreeScan         int     `json:"monthly_old_customer_and_three_scan"`          // 老客且一码三关注的指标人数
	MonthlyOldCustomerAndThreeScanRate     int     `json:"monthly_old_customer_and_three_scan_rate"`     // 老客一码三关注率(老客且一码三关注的完成人数 / 老客人数 * 100%)
	MonthlyOldCustomerAndThreeWageProgress float64 `json:"monthly_old_customer_and_three_wage_progress"` // 老客工资完成进度(老客且 一码三关注 的完成 人数 / 老客 且 一码三关注 的指标 人数 * 100%)
	NewCustomerAndThreeScan                int     `json:"new_customer_and_three_scan"`                  // 新客且一码三观客人数据量
	MonthlyNewCustomerAndThreeScan         int     `json:"monthly_new_customer_and_three_scan"`          // 新客且一码三关注的目标人数
	NewCustomerAndThreeScanRate            int     `json:"new_customer_and_three_scan_rate"`             // 新客一码三关注率(新客且一码三关注的完成人数 / 新客人数 * 100%)
}
