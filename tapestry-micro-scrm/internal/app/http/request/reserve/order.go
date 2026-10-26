package reserveReq

import (
	"micro-scrm/internal/app/http/request"
	customerSearchReq "micro-scrm/internal/app/http/request/customer_search"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
)

type EnterpriseBase struct {
	EnterpriseId int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
}

type PageReq struct {
	Page int `json:"page" binding:"min=1" comment:"页码"`
	Size int `json:"size" binding:"min=1,max=99999999" comment:"分页数"`
}

type OrderBase struct {
	ID            int64  `json:"id" binding:"required" comment:"ID"`
	ReserveNumber string `json:"reserve_number" binding:"required" comment:"预约单编号"`
}

type OrderBaseReq struct {
	OrderBase
	EnterpriseBase
}

type OrderInfo struct {
	ReserveId         int64  `json:"reserve_id" binding:"required"`
	ReserveTime       int    `json:"reserve_time" binding:"required"`
	StoreId           int64  `json:"store_id" binding:"required"`
	StoreName         string `json:"store_name" binding:"required"`
	StoreAddress      string `json:"store_address"`
	ReserveType       int    `json:"reserve_type"`
	Status            int    `json:"status"`
	UserName          string `json:"user_name" binding:"required"`
	UserUnionid       string `json:"user_unionid" binding:"required"`
	UserExternalId    string `json:"user_external_id" binding:"required"`
	Mobile            string `json:"mobile"`
	ReserveRemark     string `json:"reserve_remark"`
	UserRemark        string `json:"user_remark"`
	IsOpen            int    `json:"is_open"`
	IsCommit          int    `json:"is_commit"`
	ReserveOrderSn    string `json:"reserve_order_sn"`
	ReserveOrderPrice int64  `json:"reserve_order_price"`
	ReserveEndTime    int    `json:"reserve_end_time"`
}

type OrderInfoReq struct {
	EnterpriseBase
	OrderInfo
}

type OrderInfoRes struct {
	microScrmModel.ReserveOrderModel
	ServeName    string `json:"serve_name"`
	ServeImg     string `json:"serve_img"`
	ServeTime    int32  `json:"serve_time"`
	ServiceScore int32  `json:"service_score"`
	StoreScore   int32  `json:"store_score"`
	Remark       string `json:"remark"`
}

type MemberOrderInfoRes struct { // 导购端的预约单信息
	microScrmModel.ReserveOrderModel
	ServeName         string `json:"serve_name"`
	ServeImg          string `json:"serve_img"`
	ServeTime         int32  `json:"serve_time"`
	Mobile            string `json:"mobile"`
	UcDepartmentId    string `json:"uc_department_id"`
	ServiceScore      int32  `json:"service_score"`
	StoreScore        int32  `json:"store_score"`
	Remark            string `json:"remark"`
	IsOwn             int    `json:"is_own" comment:"判断导购是否是自己的邀约单"`
	CustomerStatusStr string `json:"customer_status_str"`
	RelationOrderNum  int    `json:"relation_order_num"`
}

type OrderListInfo struct {
	OrderBase
	OrderInfo
}

type UserOrderBase struct {
	OrderBase
	EnterpriseBase
	UnionId string `json:"union_id"`
}

type UserEditOrder struct {
	UserOrderBase
	ReserveId      int64  `json:"reserve_id" binding:"required"`
	ReserveTime    int    `json:"reserve_time" binding:"required"`
	ReserveEndTime int    `json:"reserve_end_time"`
	UserRemark     string `json:"user_remark"`
	Status         int    `json:"status"`
}

type UserConfirmOrder struct {
	Id             int64 `json:"id"`
	ReserveTime    int   `json:"reserve_time"`
	ReserveEndTime int64 `json:"reserve_end_time"`
}

type MemberEditOrder struct {
	ID int64 `json:"id" binding:"required" comment:"ID"`
	EnterpriseBase
	ReserveId      int64  `json:"reserve_id" binding:"required"`
	Mobile         string `json:"mobile"`
	ReserveTime    int    `json:"reserve_time" binding:"required"`
	ReserveEndTime int    `json:"reserve_end_time"`
	ReserveRemark  string `json:"reserve_remark"`
	StoreAddress   string `json:"store_address"`
}

type UserScoreReq struct {
	UserOrderBase
	ServiceScore int32  `json:"service_score" binding:"gte=0,lte=5" gte:"评分必须大于等于0" lte:"评分不能大于5"`
	StoreScore   int32  `json:"store_score" binding:"gte=0,lte=5" gte:"评分必须大于等于0" lte:"评分不能大于5"`
	FlowScore    int32  `json:"flow_score"`
	Remark       string `json:"remark"`
}

type OrderListReq struct {
	EnterpriseBase
	PageReq
	MemberIds             []int64  `json:"member_ids" comment:"服务导购的id"`
	StoreIds              []int64  `json:"store_ids" comment:"服务的门店id"`
	ServiceId             int64    `json:"service_id"`
	Status                int      `json:"status"`
	StartTime             int64    `json:"start_time"`
	EndTime               int64    `json:"end_time"`
	ReserveStartTime      int64    `json:"reserve_start_time"`
	ReserveEndTime        int64    `json:"reserve_end_time"`
	Ids                   []int64  `json:"ids" comment:"ids"`
	IsStore               int      `json:"is_store"`
	IsComment             int      `json:"is_comment"`
	IsPay                 int      `json:"is_pay"`
	Unionid               string   `json:"unionid"`
	IsSM                  int      `json:"is_sm" comment:"是否是店长"`
	OrderOrigin           int      `json:"order_origin"`
	Name                  string   `json:"name"`
	OrderSort             int      `json:"order_sort"`
	MemberId              int64    `json:"member_id" comment:"店长自己的MemberId"`
	ReserveOrderStartTime int64    `json:"reserve_order_start_time" comment:"邀约成功交易时间"`
	ReserveOrderEndTime   int64    `json:"reserve_order_end_time" comment:"邀约成功交易时间"`
	ReserveNumber         string   `json:"reserve_number" comment:"邀约单号"`
	ReserveNumbers        []string `json:"reserve_numbers" comment:"支付订单的邀约单号"`
	Search                request.SearchList
}

type OrderListByCoachReq struct {
	EnterpriseBase
	PageReq
	Status           int                          `json:"status"`
	ReserveStartTime int64                        `json:"reserve_start_time"`
	ReserveEndTime   int64                        `json:"reserve_end_time"`
	MemberId         int64                        `json:"member_id" comment:"自己的MemberId"`
	MemberIds        []int64                      `json:"member_ids" comment:"手动选择的导购id"`
	ThirdIds         []int                        `json:"third_ids" comment:"选择或者SM管理的门店"`
	UnionIds         []string                     `json:"union_ids"`
	IsSearch         bool                         `json:"is_search"`
	Search           customerSearchReq.SearchList `json:"search"`
}

type OrderDetailReq struct {
	EnterpriseBase
	ID             int64  `json:"id" binding:"required" comment:"ID"`
	EnterpriseName string `json:"enterprise_name" binding:"required" comment:"企业名称"`
}

type OrderAssignReq struct {
	OrderBaseReq
	Mobile string `json:"mobile" binding:"required" required:"导购手机号"`
}

type OrderTimeListReq struct {
	EnterpriseBase
	MemberId int64  `json:"member_id"`
	Mobile   string `json:"mobile"`
}

type OrderTimeInfo struct {
	ReserveTime int   `json:"reserve_time"`
	Num         int64 `json:"num"`
}

type ReserveOrderListReq struct {
	EnterpriseBase
	PageReq
	MemberIds             []int64 `json:"member_ids" comment:"服务导购的id"`
	StoreIds              []int64 `json:"store_ids" comment:"服务的门店id"`
	StartTime             int64   `json:"start_time" comment:"创建时间"`
	EndTime               int64   `json:"end_time" comment:"创建时间"`
	Ids                   []int64 `json:"ids" comment:"ids"`
	ReserveStartTime      int64   `json:"reserve_start_time" comment:"邀约的到店时间"`
	ReserveEndTime        int64   `json:"reserve_end_time" comment:"邀约的到店时间"`
	ReserveOrderStartTime int64   `json:"reserve_order_start_time" comment:"订单时间"`
	ReserveOrderEndTime   int64   `json:"reserve_order_end_time" comment:"订单时间"`
	CommentStartTime      int64   `json:"comment_start_time" comment:"评论时间"`
	CommentEndTime        int64   `json:"comment_end_time"  comment:"评论时间"`
	OrderOrigin           int     `json:"order_origin"`
	Status                int     `json:"status"`
	IsComment             int     `json:"is_comment"`
	IsPay                 int     `json:"is_pay"`
	IsStore               int     `json:"is_store"`
	ReserveNumber         string  `json:"reserve_number"`
	ServiceId             int64   `json:"service_id"`
}

type ReserveSum struct {
	InviteNum        int    `json:"invite_num"`         // 邀约数
	OpenNum          int    `json:"open_num"`           // 打开数
	CommitNum        int    `json:"commit_num"`         // 确认数
	OpenPer          string `json:"open_per"`           // 打开占比
	CommitPer        string `json:"commit_per"`         // 确认占比
	InviteSuccessNum int64  `json:"invite_success_num"` // 邀约成功数
	PaySuccessNum    int64  `json:"pay_success_num"`    // 到店购买数
	SuccessPrice     int64  `json:"success_price"`      // 导购成单金额
}

type PushMemberInfo struct {
	MemberUserID string `json:"member_user_id"`
	MemberId     int64  `json:"member_id"`
	Num          int64  `json:"num"`
}

type OrderListRes struct {
	microScrmModel.ReserveOrderModel
	ServiceScore int32  `json:"service_score"`
	StoreScore   int32  `json:"store_score"`
	Remark       string `json:"remark"`
}

type OrderExportListResult struct {
	OrderListRes
	ReserveName         string  `json:"reserve_name"`
	ReserveTimeStr      string  `json:"reserve_time_str"`
	CreateTimeStr       string  `json:"create_time_str"`
	IsOpenStr           string  `json:"is_open_str"`
	IsCommitStr         string  `json:"is_commit_str"`
	StatusStr           string  `json:"status_str"`
	ReserveTypeStr      string  `json:"reserve_type_str"`
	IsSignInStr         string  `json:"is_sign_in_str"`
	IsCommentStr        string  `json:"is_comment_str"`
	IsStoreStr          string  `json:"is_store_str"`
	IsPayStr            string  `json:"is_pay_str"`
	CommentTimeStr      string  `json:"comment_time_str"`
	ReserveOrderTimeStr string  `json:"reserve_order_time_str"`
	OrderOriginStr      string  `json:"order_origin_str"`
	CustomerStatusStr   string  `json:"customer_status_str"`
	ServiceScoreStr     string  `json:"service_score_str"`
	StoreScoreStr       string  `json:"store_score_str"`
	OrderPrice          float64 `json:"order_price"`
	RelationOrderNumStr string  `json:"relation_order_num_str"`
	RelationOrderNum    int     `json:"relation_order_num"`
}

type OrderListResult struct {
	microScrmModel.ReserveOrderModel
	CreateTimeStr     string `json:"create_time_str"`
	ReserveName       string `json:"reserve_name"`
	ReserveTimeStr    string `json:"reserve_time_str"`
	StatusStr         string `json:"status_str"`
	ReserveTypeStr    string `json:"reserve_type_str"`
	CustomerStatusStr string `json:"customer_status_str"`
	UnionCode         string `json:"union_code"`
	RelationOrderNum  int    `json:"relation_order_num"`
}

type OrderDetailRes struct {
	microScrmModel.ReserveOrderModel
	LogList          []LogList      `json:"log_list"`
	LogListArr       [][]LogListArr `json:"log_list_arr"`
	Avatar           string         `json:"avatar"`
	ReserveName      string         `json:"reserve_name"`
	ReserveTimeStr   string         `json:"reserve_time_str"`
	CreateTimeStr    string         `json:"create_time_str"`
	IsOpenStr        string         `json:"is_open_str"`
	IsCommitStr      string         `json:"is_commit_str"`
	StatusStr        string         `json:"status_str"`
	IsCommentStr     string         `json:"is_comment_str"` // 是否评价
	IsPayStr         string         `json:"is_pay_str"`     // 是否成单
	ServiceScore     int32          `json:"service_score"`
	StoreScore       int32          `json:"store_score"`
	Remark           string         `json:"remark"`
	IsSignInStr      string         `json:"is_sign_in_str"`
	ServiceScoreStr  string         `json:"service_score_str"`
	StoreScoreStr    string         `json:"store_score_str"`
	RelationOrderNum int            `json:"relation_order_num"`
}

type StatusList struct {
	Status    int32  `json:"status"`
	StatusStr string `json:"status_str"`
}

type OrderRankList struct {
	EnterpriseBase
	StartTime int64 `json:"start_time" binding:"required" comment:"开始时间"`
	EndTime   int64 `json:"end_time" binding:"required" comment:"结束时间"`
	PageReq
}

type OrderRankRes struct {
	MemberID       int64   `json:"member_id"`
	Num            int64   `json:"num"`
	SuccessNum     int64   `json:"success_num"`
	Price          float64 `json:"price"`
	MemberName     string  `json:"member_name"`
	Avatar         string  `json:"avatar"`
	Index          int     `json:"index"`
	DepartmentName string  `json:"department_name"`
	StoreId        int64   `json:"store_id"`
}

type OrderPriceRank struct {
	MemberID int64   `json:"member_id"`
	Price    float64 `json:"price"`
}

type WriteOffReq struct {
	OrderBaseReq
	ReserveOrderSn    string `json:"reserve_order_sn"`
	ReserveOrderTime  int    `json:"reserve_order_time"`
	ReserveOrderPrice int    `json:"reserve_order_price"`
}

type ScoreInfoRes struct {
	microScrmModel.ReserveOrderModel
	IsScore      int32  `json:"is_score"`
	Remark       string `json:"remark"`
	ServiceScore int32  `json:"service_score"`
	StoreScore   int32  `json:"store_score"`
	FlowScore    int32  `json:"flow_score"`
	ServeName    string `json:"serve_name"`
}

type UserSignIn struct {
	UserOrderBase
	IsSign int `json:"is_sign"`
}

type OrderDetailResBaseStr struct {
	microScrmModel.ReserveOrderModel
	ReserveTimeStr      string  `json:"reserve_time_str"`
	CreateTimeStr       string  `json:"create_time_str"`
	IsOpenStr           string  `json:"is_open_str"`
	IsCommitStr         string  `json:"is_commit_str"`
	StatusStr           string  `json:"status_str"`
	ReserveTypeStr      string  `json:"reserve_type_str"`
	IsSignInStr         string  `json:"is_sign_in_str"`
	IsCommentStr        string  `json:"is_comment_str"`
	IsStoreStr          string  `json:"is_store_str"`
	IsPayStr            string  `json:"is_pay_str"`
	CommentTimeStr      string  `json:"comment_time_str"`
	ReserveOrderTimeStr string  `json:"reserve_order_time_str"`
	OrderOriginStr      string  `json:"order_origin_str"`
	CustomerStatusStr   string  `json:"customer_status_str"`
	OrderPrice          float64 `json:"order_price"`
	RelationOrderNum    int     `json:"relation_order_num"`
	RelationOrderNumStr string  `json:"relation_order_num_str"`
}

type ChangeStatus struct {
	UserOrderBase
	IsStore        int    `json:"is_store"`
	CustomerStatus int    `json:"customer_status"`
	Mobile         string `json:"mobile"`
}
type ConfirmStore struct {
	UserOrderBase
}

type AcxiomSftpInfo struct {
	ReserveNumber      string  `json:"reserve_number"`
	ReserveOrderTime   int     `json:"reserve_order_time"`
	ReserveOrderPrice  float64 `json:"reserve_order_price"`
	ReserveOrderSn     int     `json:"reserve_order_sn"`
	ReserveOrderStatus int32   `json:"reserve_order_status"`
}

type AcxiomStoreAppointmentQTYInfo struct {
	StoreId        int64  `json:"store_id"`
	StoreCode      string `json:"store_code"`
	AppointmentQTY int64  `json:"appointment_qty" gorm:"appointment_qty"`
}

type AcxiomStoreReportDataInfo struct {
	MemberId       int64  `json:"member_id"`
	MemberUserId   string `json:"member_user_id"`
	StoreId        int64  `json:"store_id"`
	StoreCode      string `json:"store_code"`
	AppointmentQTY int64  `json:"appointment_qty" gorm:"appointment_qty"`
}

type AchievementParams struct {
	EnterpriseId int64   `json:"enterprise_id"`
	MemberIds    []int64 `json:"member_ids"`
	Type         int64   `json:"type"`
	StartTime    int64   `json:"start_time"`
	EndTime      int64   `json:"end_time"`
}

type AllRes struct {
	MemberId int64   `json:"member_id"`
	AllTotal float64 `json:"all_total"`
}

type CheckReserveTime struct {
	EnterpriseBase
	ReserveTime int64 `json:"reserve_time"`
	MemberId    int64 `json:"member_id"`
	StoreId     int64 `json:"store_id"`
}

type CustomerReserveListsReq struct {
	EnterpriseBase
	ExternalUserID string `json:"external_userid"`
	Page           int    `json:"page"`
	Size           int    `json:"size"`
}

type OrderShowReq struct {
	OrderBaseReq
	MemberId int64 `json:"member_id"`
}

type ReserveSalesNum struct {
	Num               int    `json:"num"`
	ReserveNumber     string `json:"reserve_number"`
	ReserveOrderPrice int    `json:"reserve_order_price"`
}

type ExportOrderRecord struct {
	ID                int64  `json:"id"`
	ReserveNumber     string `json:"reserve_number"`
	ReserveTime       int    `json:"reserve_time"`
	StoreName         string `json:"store_name"`
	Status            int32  `json:"status"`
	UserName          string `json:"user_name"`
	MemberID          int64  `json:"member_id"`
	MemberName        string `json:"member_name"`
	IsOpen            int32  `json:"is_open"`
	IsCommit          int32  `json:"is_commit"`
	InviteNum         int    `json:"invite_num"`
	IsStore           int32  `json:"is_store"`
	IsPay             int32  `json:"is_pay"`
	OrderOrigin       int32  `json:"order_origin"`
	CustomerStatus    int32  `json:"customer_status"`
	ReserveTimeStr    string `json:"reserve_time_str"`
	IsOpenStr         string `json:"is_open_str"`
	IsCommitStr       string `json:"is_commit_str"`
	StatusStr         string `json:"status_str"`
	IsStoreStr        string `json:"is_store_str"`
	IsPayStr          string `json:"is_pay_str"`
	OrderOriginStr    string `json:"order_origin_str"`
	CustomerStatusStr string `json:"customer_status_str"`
	ExportOrderSales
}
type ExportOrderSales struct {
	RelationOrderID string  `json:"relation_order_id"`
	GoodsName       string  `json:"goods_name"`
	Style           string  `json:"style"`
	Color           string  `json:"color"`
	OrderPrice      float64 `json:"order_price"`
	GoodsPrice      int     `json:"goods_price"`
	PayTimeInt      int     `json:"pay_time_int"`
	PayTime         string  `json:"pay_time"`
	StoreCode       string  `json:"store_code"`
}

type ReserveMatchOrder struct {
	OrderNo           string  `json:"order_no"`
	SaCode            string  `json:"sa_code"`
	UnionId           string  `json:"union_id"`
	PurchaseStoreCode string  `json:"purchase_store_code"`
	PurchaseCode      string  `json:"purchase_code"`
	OriPurchaseCode   string  `json:"ori_purchase_code"`
	PurchaseAmount    float64 `json:"purchase_amount"`
	PurchaseDate      string  `json:"purchase_date"`
	Status            int64   `json:"status"`
	NewCustomer       int64   `json:"new_customer"`
	CreateDate        string  `json:"create_date"`
}

type ReserveMatchOrderProduct struct {
	PurchaseCode string  `json:"purchase_code"`
	SkuName      string  `json:"sku_name"`
	Style        string  `json:"style"`
	Color        string  `json:"color"`
	Amount       float64 `json:"amount"`
	PurchaseDate string  `json:"purchase_date"`
	CreateDate   string  `json:"create_date"`
}

type ReserveSalesStatusInfo struct {
	ReserveNumber string `json:"reserve_number"`
	OrderPrice    int64  `json:"order_price"`
	LastPayTime   int64  `json:"last_pay_time"`
	NewCustomer   int64  `json:"new_customer"`
}
