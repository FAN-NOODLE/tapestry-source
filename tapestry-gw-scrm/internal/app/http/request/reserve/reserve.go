package reserveReq

type EnterpriseBase struct {
	EnterpriseID int64 `json:"enterprise_id"`
}

type PageReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type OrderListReq struct {
	EnterpriseBase
	PageReq
	Total                 int64   `json:"total"`
	MemberIds             []int64 `json:"member_ids" comment:"服务导购的id"`
	StoreIds              []int64 `json:"store_ids" comment:"服务的门店id"`
	ServiceId             int64   `json:"service_id"`
	Status                int     `json:"status"`
	StartTime             int64   `json:"start_time"`
	EndTime               int64   `json:"end_time"`
	ReserveStartTime      int64   `json:"reserve_start_time"`
	ReserveEndTime        int64   `json:"reserve_end_time"`
	Ids                   []int64 `json:"ids" comment:"ids"`
	IsStore               int     `json:"is_store"`
	IsComment             int     `json:"is_comment"`
	IsPay                 int     `json:"is_pay"`
	OrderOrigin           int     `json:"order_origin"`
	Name                  string  `json:"name"`
	ReserveOrderStartTime int64   `json:"reserve_order_start_time" comment:"邀约成功交易时间"`
	ReserveOrderEndTime   int64   `json:"reserve_order_end_time" comment:"邀约成功交易时间"`
	ReserveNumber         string  `json:"reserve_number"`
}

type OrderDetailReq struct {
	EnterpriseBase
	ID             int64  `json:"id" binding:"required" comment:"ID"`
	EnterpriseName string `json:"enterprise_name" comment:"企业名称"`
}

type RoleRules struct {
	ID                 int64  `json:"id"` // is_dp为false时，id为member_id
	ThirdId            int    `json:"third_id"`
	Name               string `json:"name"` // is_dp为false时，name为member_name;is_dp为ture,name为store_name
	IsDP               bool   `json:"is_dp"`
	StoreType          int32  `json:"store_type"`
	MemSubscribeStatus int    `json:"mem_subscribe_status"`
	UserId             string `json:"user_id"`
}

type ReserveServiceCreateReq struct {
	EnterpriseBase
	ReserveName string      `json:"reserve_name" binding:"required" comment:"服务类型名称"`
	ReserveImg  string      `json:"reserve_img" comment:"服务类型的图片"`
	ReserveTime int32       `json:"reserve_time" binding:"required" comment:"服务时长"`
	EndTime     int32       `json:"end_time" comment:"截止时间"`
	ServiceType int32       `json:"service_type"`
	StartTime   int         `json:"start_time"`
	ServiceRole int32       `json:"service_role"`
	RoleRules   []RoleRules `json:"role_rules"`
	TimeSlot    []TimeSlot  `json:"time_slot"`
}

type ReserveServiceUpdateReq struct {
	EnterpriseBase
	ID          int64       `json:"id" binding:"required" comment:"ID"`
	ReserveName string      `json:"reserve_name" comment:"服务类型名称"`
	ReserveImg  string      `json:"reserve_img" comment:"服务类型的图片"`
	ReserveTime int32       `json:"reserve_time" comment:"服务时长"`
	EndTime     int32       `json:"end_time" binding:"required" comment:"截止时间"`
	ServiceType int32       `json:"service_type"`
	StartTime   int         `json:"start_time"`
	ServiceRole int32       `json:"service_role"`
	RoleRules   []RoleRules `json:"role_rules"`
	TimeSlot    []TimeSlot  `json:"time_slot"`
}

type ReserveServiceDeleteReq struct {
	EnterpriseBase
	ID int64 `json:"id" binding:"required" comment:"ID"`
}

type ReserveServiceListReq struct {
	EnterpriseBase
	PageReq
	ReserveName   string `json:"reserve_name" comment:"服务类型名称"`
	ReserveStatus int    `json:"reserve_status" comment:"是否有效, -1全部，0为有效，1为失效"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
}

type ServiceScoreListReq struct {
	EnterpriseBase
	PageReq
	Total            int64   `json:"total"`
	MemberIds        []int64 `json:"member_ids" comment:"服务导购的id"`
	StoreIds         []int64 `json:"store_ids" comment:"服务的门店id"`
	ReserveStartTime int     `json:"reserve_start_time" comment:"预约开始时间"`
	ReserveEndTime   int     `json:"reserve_end_time" comment:"预约结束时间"`
	ScoreStartTime   int     `json:"score_start_time" comment:"评分开始时间"`
	ScoreEndTime     int     `json:"score_end_time" comment:"评分结束时间"`
	Ids              []int64 `json:"ids" comment:"ids"`
}

type ReserveOrderListReq struct {
	EnterpriseBase
	PageReq
	Total                 int64   `json:"total"`
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

type ResultReserveOrderList struct {
	Total int64 `json:"total"`
	List  []ResultReserveOrderData
}
type ResultReserveOrderData struct {
	EnterpriseBase
	ReserveNumber       string `json:"reserve_number" comment:"预约单的编号"`
	ReserveID           int64  `json:"reserve_id" comment:"预约服务的项目id"`
	ReserveTime         int    `json:"reserve_time" comment:"预约服务的时间"`
	StoreID             int64  `json:"store_id" comment:"门店id(部门id)"`
	StoreName           string `json:"store_name" comment:"预约的门店名称"`
	StoreAddress        string `json:"store_address" comment:"预约的门店地址"`
	ReserveType         int32  `json:"reserve_type" comment:"预约订单是否是分配单,0:默认预约单，1:分配的预约单"`
	Status              int32  `json:"status" comment:"预约单的状态"`
	UserName            string `json:"user_name" comment:"预约服务的用户名称"`
	UserUnionid         string `json:"user_unionid" comment:"预约服务的用户unionid"`
	UserExternalID      string `json:"user_external_id" comment:"预约服务的用户external_id"`
	MemberID            int64  `json:"member_id" comment:"预约单服务的导购id"`
	MemberName          string `json:"member_name" comment:"预约单服务的导购名称"`
	ReserveRemark       string `json:"reserve_remark" comment:"预约服务的备注信息(导购)"`
	UserRemark          string `json:"user_remark" comment:"预约服务的备注信息(用户)"`
	IsOpen              int32  `json:"is_open" comment:"用户是否打开预约单 0否，1是"`
	IsCommit            int32  `json:"is_commit" comment:"用户是否确认预约单 0否，1是"`
	InviteNum           int    `json:"invite_num" comment:"邀约次数"`
	MemberUserID        string `json:"member_user_id" comment:"预约单服务的导购user_id"`
	ReserveOrderSn      string `json:"reserve_order_sn" comment:"coach核销单号"`
	ReserveOrderPrice   int    `json:"reserve_order_price" comment:"coach核销导购的订单金额（元）"`
	CreatedAt           int64  `json:"created_at"`
	ReserveName         string `json:"reserve_name"`
	ReserveTimeStr      string `json:"reserve_time_str"`
	CreateTimeStr       string `json:"create_time_str"`
	IsOpenStr           string `json:"is_open_str"`
	IsCommitStr         string `json:"is_commit_str"`
	StatusStr           string `json:"status_str"`
	ReserveTypeStr      string `json:"reserve_type_str"`
	IsSignInStr         string `json:"is_sign_in_str"`
	IsCommentStr        string `json:"is_comment_str"`
	IsStoreStr          string `json:"is_store_str"`
	IsStore             int    `json:"is_store"`
	IsPay               int    `json:"is_pay"`
	IsPayStr            string `json:"is_pay_str"`
	CommentTimeStr      string `json:"comment_time_str"`
	ReserveOrderTimeStr string `json:"reserve_order_time_str"`
	OrderOrigin         int32  `json:"order_origin"`
	OrderOriginStr      string `json:"order_origin_str"`
	CustomerStatus      int32  `json:"customer_status"`
}

type ResultOrderList struct {
	Total int64 `json:"total"`
	List  []ResultOrderData
}

type ResultOrderData struct {
	EnterpriseBase
	ReserveNumber       string `json:"reserve_number" comment:"预约单的编号"`
	ReserveID           int64  `json:"reserve_id" comment:"预约服务的项目id"`
	ReserveTime         int    `json:"reserve_time" comment:"预约服务的时间"`
	StoreID             int64  `json:"store_id" comment:"门店id(部门id)"`
	StoreName           string `json:"store_name" comment:"预约的门店名称"`
	StoreAddress        string `json:"store_address" comment:"预约的门店地址"`
	ReserveType         int32  `json:"reserve_type" comment:"预约订单是否是分配单,0:默认预约单，1:分配的预约单"`
	Status              int32  `json:"status" comment:"预约单的状态"`
	UserName            string `json:"user_name" comment:"预约服务的用户名称"`
	UserUnionid         string `json:"user_unionid" comment:"预约服务的用户unionid"`
	UserExternalID      string `json:"user_external_id" comment:"预约服务的用户external_id"`
	MemberID            int64  `json:"member_id" comment:"预约单服务的导购id"`
	MemberName          string `json:"member_name" comment:"预约单服务的导购名称"`
	ReserveRemark       string `json:"reserve_remark" comment:"预约服务的备注信息(导购)"`
	UserRemark          string `json:"user_remark" comment:"预约服务的备注信息(用户)"`
	IsOpen              int32  `json:"is_open" comment:"用户是否打开预约单 0否，1是"`
	IsCommit            int32  `json:"is_commit" comment:"用户是否确认预约单 0否，1是"`
	InviteNum           int    `json:"invite_num" comment:"邀约次数"`
	MemberUserID        string `json:"member_user_id" comment:"预约单服务的导购user_id"`
	ReserveOrderSn      string `json:"reserve_order_sn" comment:"coach核销单号"`
	ReserveOrderPrice   int    `json:"reserve_order_price" comment:"coach核销导购的订单金额（元）"`
	ReserveName         string `json:"reserve_name"`
	ReserveTimeStr      string `json:"reserve_time_str"`
	CreateTimeStr       string `json:"create_time_str"`
	IsOpenStr           string `json:"is_open_str"`
	IsCommitStr         string `json:"is_commit_str"`
	StatusStr           string `json:"status_str"`
	ReserveTypeStr      string `json:"reserve_type_str"`
	IsSignInStr         string `json:"is_sign_in_str"`
	ReserveOrderTimeStr string `json:"reserve_order_time_str"`
	IsCommentStr        string `json:"is_comment_str"`
	IsStoreStr          string `json:"is_store_str"`
	IsStore             int    `json:"is_store"`
	IsPayStr            string `json:"is_pay_str"`
	IsPay               int    `json:"is_pay"`
	ServiceScore        int32  `json:"service_score"`
	StoreScore          int32  `json:"store_score"`
	Remark              string `json:"remark"`
	CommentTimeStr      string `json:"comment_time_str"`
	OrderOrigin         int32  `json:"order_origin"`
	OrderOriginStr      string `json:"order_origin_str"`
	CustomerStatus      int32  `json:"customer_status"`
	CustomerStatusStr   string `json:"customer_status_str"`
	VipType             int    `json:"vip_type"`
}

type ResultScoreList struct {
	Total int64 `json:"total"`
	List  []ResultScoreData
}
type ResultScoreData struct {
	EnterpriseBase
	ReserveNumber  string `json:"reserve_number" comment:"预约单编号"`
	ReserveOrderID int64  `json:"reserve_order_id" comment:"预约单表id"`
	UserName       string `json:"user_name" comment:"评论的用户昵称"`
	UserUnionid    string `json:"user_unionid" comment:"评论的用户unionid"`
	UserExternalID string `json:"user_external_id" comment:"评论的用户external_id"`
	MemberID       int64  `json:"member_id" comment:"服务的导购id"`
	MemberName     string `json:"member_name" comment:"服务的导购昵称"`
	MemberUserID   string `json:"member_user_id" comment:"服务导购的user_id"`
	StoreID        int64  `json:"store_id" comment:"服务的门店id"`
	StoreName      string `json:"store_name" comment:"服务的门店名称"`
	ServiceScore   int32  `json:"service_score" comment:"默认服务分数为5分（KS端店员服务,sw端活动安排）"`
	StoreScore     int32  `json:"store_score" comment:"门店的评分，默认5分（KS端店铺环境,sw端整体服务）"`
	FlowScore      int32  `json:"flow_score" comment:"默认5分'流程的评分（ks端服务流程）"`
	Remark         string `json:"remark" comment:"用户的评论信息"`
	CreatedAt      int64  `json:"created_at"`
	CreateTimeStr  string `json:"create_time_str"`
	ReserveTime    int    `json:"reserve_time" comment:"预约时间(预约单的到店时间)"`
	ReserveTimeStr string `json:"reserve_time_str" comment:"预约时间(预约单的到店时间)"`
}

type TimeSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
