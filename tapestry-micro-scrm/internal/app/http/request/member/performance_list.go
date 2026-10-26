package memberReq

type PerformanceListReq struct {
	Page          int     `json:"page" default:"1" comment:"页码"`
	Size          int     `json:"size" default:"10" comment:"分页数"`
	EnterpriseId  int64   `json:"enterprise_id" binding:"required" comment:"企业ID不能为空"`
	ThirdId       []int   `json:"third_id" binding:"required" comment:"门店id不能为空"`
	FinancialTime []int64 `json:"financial_time" binding:"required" comment:"时间不能为空"`
	ActiveType    int     `json:"active_type" binding:"required" comment:"类型不能为空"`
}

type PerformanceListRes struct {
	EnterpriseId    int64  `json:"enterprise_id"`
	Avatar          string `json:"avatar"`
	Name            string `json:"name"`
	StoreName       string `json:"store_name"`
	MemberID        int64  `json:"member_id"`
	ChatCnt         int    `json:"chat_cnt"`
	MessageCnt      int    `json:"message_cnt"`
	ReplyPercentage string `json:"reply_percentage"`
	NewContactCnt   int    `json:"new_contact_cnt"`
}

type PerformanceListActiveRes struct {
	MemberID        int64  `json:"member_id"`
	ChatCnt         int    `json:"chat_cnt"`
	MessageCnt      int    `json:"message_cnt"`
	ReplyPercentage string `json:"reply_percentage"`
	NewContactCnt   int    `json:"new_contact_cnt"`
}

// 邀约参数

type PerformanceListInvitationReq struct {
	Page          int     `json:"page" default:"1" comment:"页码"`
	Size          int     `json:"size" default:"10" comment:"分页数"`
	EnterpriseId  int64   `json:"enterprise_id" binding:"required" comment:"企业ID不能为空"`
	UserType      int64   `json:"user_type"`
	Mobile        string  `json:"mobile" binding:"required" comment:"导购手机号不能为空"`
	ThirdId       []int   `json:"third_id" comment:"门店id不能为空"`
	FinancialTime []int64 `json:"financial_time" binding:"required" comment:"时间不能为空"`
}

// 邀约榜

type PerformanceListInvitationRes struct {
	EnterpriseId   int64    `json:"enterprise_id"`
	Avatar         string   `json:"avatar"`
	Name           string   `json:"name"`
	StoreName      string   `json:"store_name"`
	MemberID       int64    `json:"member_id"`
	Num            int      `json:"num"`
	SuccessNum     int      `json:"success_num"`
	PriceTotal     int      `json:"price_total"`
	IsView         int      `json:"is_view"`
	ReserveNumbers []string `json:"reserve_numbers"`
}

// 邀约查询数据库返回的信息

type PerformanceListReserveOrderRes struct {
	MemberID   int64  `json:"member_id"`
	Num        int    `json:"num"`
	StoreCity  string `json:"store_city"`
	SuccessNum int    `json:"success_num"`
	PriceTotal int    `json:"price_total"`
}

// 业绩榜单

type PerformanceListAchievementReq struct {
	Page          int     `json:"page" default:"1" comment:"页码"`
	Size          int     `json:"size" default:"10" comment:"分页数"`
	EnterpriseId  int64   `json:"enterprise_id" binding:"required" comment:"企业ID不能为空"`
	ThirdId       []int   `json:"third_id"  comment:"门店id不能为空"`
	FinancialTime []int64 `json:"financial_time" binding:"required" comment:"时间不能为空"`
	Start         int     `json:"start"  comment:"榜单类型不能为空"`
}

// 业绩榜单

type PerformanceListAchievementRes struct {
	EnterpriseId     int64  `json:"enterprise_id"`
	Avatar           string `json:"avatar"`
	Name             string `json:"name"`
	StoreName        string `json:"store_name"`
	MemberID         int64  `json:"member_id"`
	OrderTotal       int64  `json:"order_total"`
	OrderPriceTotal  int64  `json:"order_price_total"`
	RefundOrderTotal int64  `json:"refund_order_total"`
	PayOrderTotal    int64  `json:"pay_order_total"`
}

// SwCustomersListReq SW顾客榜单
type SwCustomersListReq struct {
	Page          int     `json:"page" default:"1" comment:"页码"`
	Size          int     `json:"size" default:"10" comment:"分页数"`
	EnterpriseId  int64   `json:"enterprise_id" binding:"required" comment:"企业ID不能为空"`
	ThirdId       []int   `json:"third_id"  comment:"门店id不能为空"`
	Start         int     `json:"start"  comment:"榜单类型不能为空"`
	FinancialTime []int64 `json:"financial_time"  comment:"时间不能为空"`
}

// SwCustomersListReqRes SW顾客榜单返回
type SwCustomersListReqRes struct {
	EnterpriseId   int64  `json:"enterprise_id"`
	Avatar         string `json:"avatar"`
	Name           string `json:"name"`
	StoreName      string `json:"store_name"`
	MemberID       int64  `json:"member_id"`
	NewCustomers   int64  `json:"new_customers"`
	OldCustomers   int64  `json:"old_customers"`
	SleepCustomers int64  `json:"sleep_customers"`
	VipCustomers   int64  `json:"vip_customers"`
}

// SwCustomerBaseRes SW顾客榜单基本数据返回
type SwCustomerBaseRes struct {
	MemberID       int64 `json:"member_id"`
	NewCustomers   int64 `json:"new_customers"`
	OldCustomers   int64 `json:"old_customers"`
	SleepCustomers int64 `json:"sleep_customers"`
	VipCustomers   int64 `json:"vip_customers"`
}

// SwActiveListReqReq sw活跃榜单
type SwActiveListReqReq struct {
	Page          int     `json:"page" default:"1" comment:"页码"`
	Size          int     `json:"size" default:"10" comment:"分页数"`
	EnterpriseId  int64   `json:"enterprise_id" binding:"required" comment:"企业ID不能为空"`
	ThirdId       []int   `json:"third_id"  comment:"门店id不能为空"`
	FinancialTime []int64 `json:"financial_time" binding:"required" comment:"时间不能为空"`
	Start         int     `json:"start"  comment:"榜单类型不能为空"`
}

// SwActiveListReqRes sw活跃榜单
type SwActiveListReqRes struct {
	EnterpriseId    int64   `json:"enterprise_id"`
	Avatar          string  `json:"avatar"`
	Name            string  `json:"name"`
	StoreName       string  `json:"store_name"`
	MemberID        int64   `json:"member_id"`
	ChatCnt         int     `json:"chat_cnt"`
	MessageCnt      int     `json:"message_cnt"`
	ReplyPercentage string  `json:"reply_percentage"`
	NewContactCnt   int     `json:"new_contact_cnt"`
	CompletionRate  float64 `json:"completion_rate"`
}

type ExternalStatisticsRes struct {
	MemberID       int64 `json:"member_id"`
	NewCustomers   int64 `json:"new_customers"`
	OldCustomers   int64 `json:"old_customers"`
	SleepCustomers int64 `json:"sleep_customers"`
	VipCustomers   int64 `json:"vip_customers"`
}

type PerformanceListInvitationMemberInfo struct {
	MemberId  int64  `json:"member_id"`
	StoreCity string `json:"store_city"`
}
