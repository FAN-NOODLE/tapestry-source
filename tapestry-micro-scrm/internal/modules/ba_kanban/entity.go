package ba_kanban

import "micro-scrm/internal/orm/model/micro_scrm_enterprise"

type BaKanBanDataResponse struct {
	SaleData            `json:"sale_data"`
	GrowthData          `json:"growth_data"`
	CustomerNumData     `json:"customer_num_data"`
	OneQrThreeFocusData `json:"one_qr_three_focus_data"`
}

type SaleData struct {
	SaleAmount   float32 `json:"sales"`
	FinishRate   int8    `json:"finish_rate"`
	TargetAmount float32 `json:"target_amount"`
	SurplusSales float32 `json:"surplus_sales"` // 销售差额
	SurplusRate  int8    `json:"surplus_rate"`
	OrderNum     int     `json:"order_num"`
	OrderPrice   float32 `json:"order_price"`
}

type GrowthData struct {
	CurrentSaleAmount float32 `json:"current_sale_amount"`
	LastSaleAmount    float32 `json:"last_sale_amount"`
	GrowthRate        int8    `json:"growth_rate"`
}

// CustomerNumData 客户数据
type CustomerNumData struct {
	TotalCustomer         int     `json:"total_customer"`
	TotalSaleAmount       float32 `json:"total_sale_amount"`
	NewCustomerData       `json:"new_customer_data"`
	OldCustomerData       `json:"old_customer_data"`
	NewSecondCustomerData `json:"new_second_customer_data"`
	VisitCustomerData     `json:"visit_customer_data"`
}

type CustomerData struct {
	Num        int     `json:"num"`
	SaleAmount float32 `json:"sale_amount"`
	Rate       int8    `json:"rate"`
}

// NewCustomerData 新客户数据
type NewCustomerData struct {
	CustomerData
}

// OldCustomerData 老客户数据
type OldCustomerData struct {
	CustomerData
}

// NewSecondCustomerData 新二购客户数据
type NewSecondCustomerData struct {
	CustomerData
}

// VisitCustomerData 游客数据
type VisitCustomerData struct {
	CustomerData
}

type OneQrThreeFocusData struct {
	NewSecondOneQrThreeFocusData   `json:"new_second_one_qr_three_focus_data"`
	OldCustomerOneQrThreeFocusData `json:"old_customer_one_qr_three_focus_data"`
	NewCustomerOneQrThreeFocusData `json:"new_customer_one_qr_three_focus_data"`
}

type OneQrThreeFocusBaseData struct {
	Num        int `json:"num"`
	QrNum      int `json:"qr_num"`
	FinishRate int `json:"finish_rate"`
	Rate       int `json:"rate"`
}

type NewSecondOneQrThreeFocusData struct {
	OneQrThreeFocusBaseData
}

type OldCustomerOneQrThreeFocusData struct {
	OneQrThreeFocusBaseData
}

type NewCustomerOneQrThreeFocusData struct {
	OneQrThreeFocusBaseData
}

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

type BaPerformanceQueryReq struct {
	EnterpriseId int64    `json:"enterprise_id"`
	UserIds      []string `json:"user_ids"`
	Name         string   `json:"name"`
	IsPage       int      `json:"is_page"`
	Page         int      `json:"page"`
	Size         int      `json:"size"`
}

type BaPerformanceSavaReq struct {
	EnterpriseId int64                                               `json:"enterprise_id"`
	ModelArray   []micro_scrm_enterprise.BaPerformanceIndicatorModel `json:"model_array"`
}
type BaSourceDeptRequst struct {
	EnterpriseId int64  `json:"enterprise_id" binding:"required" desc:"公司编号"`
	Mobile       string `json:"mobile" binding:"required" desc:"当前BA的电话号码"`
}
