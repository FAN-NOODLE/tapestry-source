package customer

import "gw-scrm/pkg/micro-scrm-sdk/microScrm"

// UpdateMessageRequest 更新message请求
type UpdateMessageRequest struct {
	ID           int64  `json:"id" db:"id" desc:"ID"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	Status       int    `json:"status" db:"status" desc:"开启开关" validate:"required" required:"提醒设置开关必传"`
	SendType     int    `json:"send_type" db:"send_type" desc:"类别" validate:"required" required:"type不能为空"`
	StatTxt      string `json:"stat_txt" db:"stat_txt" desc:"设置详情" validate:"required" required:"设置不能为空"`
	QyStatTxt    string `json:"qy_stat_txt" db:"qy_stat_txt" desc:"设置详情" validate:"required" required:"设置不能为空"`
}

// MessageListRequest message list请求
type MessageListRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	SendType     int   `json:"send_type" db:"send_type" desc:"类别" validate:"required" required:"Type不能为空"`
}

type RepeatListRequest struct {
	EnterpriseID int64        `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	Page         int          `json:"page" validate:"min=1" min:"页面参数不能小于1"`
	Size         int          `json:"size" validate:"min=1" min:"分页参数不能小于1"`
	Search       ParamsSearch `json:"search"`
}

type ParamsSearch struct {
	CustomerName string `json:"customer_name"`
	DepartmentID int64  `json:"department_id"`
	MemberID     int64  `json:"member_id"`
}

type ExternalRelationByUserIdRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	UserWorkID   int64 `json:"user_work_id" db:"user_work_id" validate:"required" required:"UserWorkID不能为空"`
	UserID       int64 `json:"user_id" db:"user_id"`
	DepartmentID int64 `json:"department_id" db:"department_id"`
	MemberID     int64 `json:"member_id" db:"member_id"`
	Page         int   `json:"page" validate:"min=1" min:"页面参数不能小于1"`
	Size         int   `json:"size" validate:"min=1" min:"分页参数不能小于1"`
}

// LossCustomersListRequest 流失客户列表接口请求参数结构
type LossCustomersListRequest struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业id" validate:"required" required:"企业ID不能为空"`
	Type         int    `json:"type" desc:"流失客户类型" validate:"required" required:"type不能为空"`
	CustomerName string `json:"customer_name"`
	MemberID     int    `json:"member_id"`
	DepartmentID int    `json:"department_id"`
	TimeSlotType int    `json:"time_slot_type"`
	TimeStart    int    `json:"time_start"`
	TimeEnd      int    `json:"time_end"`
	Page         int    `json:"page" validate:"min=1" min:"页面参数不能小于1"`
	Size         int    `json:"size" validate:"min=1" min:"分页参数不能小于1"`
	ChannelId    int64  `json:"channel_id"`
}

type ResponseFootprintList []struct {
	ID        int64  `json:"id" db:"id"`
	Event     string `json:"event" db:"event"`
	EventDesc string `json:"event_desc" db:"event_desc"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
}

type ResponseMemberList []struct {
	MemberId               int64                 `json:"member_id" db:"member_id"`
	MemberName             string                `json:"member_name" db:"member_name"`
	DepartmentId           int64                 `json:"department_id" db:"department_id"`
	DepartmentName         string                `json:"department_name" db:"department_name"`
	CreateTime             int64                 `json:"create_time" db:"create_time" usage:"关系建立时间"`
	State                  string                `json:"state" json:"state" usage:"来源"`
	BelongToGroupChatCount int64                 `json:"belong_to_group_chat_count" usage:"当前客户占员工群总数的个数"`
	BelongToGroupChatList  BelongToGroupChatList `json:"belong_to_group_chat_list" usage:"当前客户占员工群列表"`
}

type BelongToGroupChat struct {
	GroupChatId   int64  `json:"group_chat_id" db:"group_chat_id"`
	GroupChatName string `json:"group_chat_name" db:"group_chat_name"`
	MemberId      int64  `json:"member_id" db:"member_id"`
	MemberName    string `json:"member_name" db:"member_name"`
	UserId        int64  `json:"user_id" db:"user_id"`
	JoinTime      int64  `json:"join_time" db:"join_time"`
}

type BelongToGroupChatList []BelongToGroupChat

type ResponseOrderList []struct {
	ID             int64      `json:"id" db:"id"`
	OrderNo        string     `json:"order_no" db:"order_no"`
	FinalAmount    int64      `json:"final_amount" db:"final_amount" usage:"含券支付金额"`
	CashFee        int64      `json:"cash_fee" db:"cash_fee" usage:"实际支付金额"`
	TotalRefundFee int64      `json:"total_refund_fee" db:"total_refund_fee" usage:"退款金额"`
	PayAt          int64      `json:"pay_at" db:"pay_at"`
	RefundAt       int64      `json:"refund_at" db:"refund_at"`
	GoodsList      GoodsList  `json:"goods_list" usage:"商品列表"`
	CouponList     CouponList `json:"coupon_list" usage:"优惠券列表"`
}

type GoodsList []struct {
	ID          int64  `json:"id" db:"id"`
	ShopOrderId int64  `json:"shop_order_id" db:"shop_order_id"`
	SpuName     string `json:"spu_name" db:"spu_name"`
}

type CouponList []struct {
	Name     string `json:"name"`
	CouponNo string `json:"coupon_no"`
}

// ResponseAddressList 收货人地址列表
type ResponseAddressList []Address
type Address struct {
	ID               int64  `json:"id" db:"id"`
	ConsigneeName    string `json:"consignee_name" db:"consignee_name" usage:"收货人姓名"`
	ConsigneeMobile  string `json:"consignee_mobile" db:"consignee_mobile" usage:"收货人手机号"`
	ConsigneeAddress string `json:"consignee_address" usage:"收货人地址"`
	Province         string `json:"province" json:"province"`
	City             string `json:"city" db:"city"`
	Area             string `json:"area" db:"area"`
	Address          string `json:"address" db:"address"`
}

// ResponseCustomerDetail 客户详情
type ResponseCustomerDetail struct {
	UserID       int64                          `json:"user_id" db:"user_id"`
	UserName     string                         `json:"user_name" db:"user_name"`
	Nickname     string                         `json:"nick_name" db:"nickname"`
	Avatar       string                         `json:"avatar" db:"avatar"`
	Gender       int                            `json:"gender" db:"gender"`
	Country      string                         `json:"country" db:"country"`
	Province     string                         `json:"province" db:"province"`
	City         string                         `json:"city" db:"city"`
	Remark       string                         `json:"remark" db:"remark"`
	RemarkList   []string                       `json:"remark_list"`
	TagList      []string                       `json:"tag_list"`
	GroupList    []string                       `json:"group_list"`
	AddressList  *microScrm.ResponseAddressList `json:"address_list"`
	IsWorkUser   bool                           `json:"is_work_user"`
	OfficialList []*DetailOfficialAccount       `json:"official_list"`
	ShowOfficial bool                           `json:"show_official"`
	LevelName    string                         `json:"level_name"`
	Birthday     string                         `json:"birthday" db:"birthday"`
	UnionCode    string                         `json:"union_code" db:"union_code"`
	Point        string                         `json:"point" db:"point"`
	Grade        string                         `json:"grade" db:"grade"`
	StoreCode    string                         `json:"store_code" db:"store_code"`
	BACode       string                         `json:"ba_code" db:"ba_code"`
	Mobile       string                         `json:"mobile" db:"mobile"`
}

// DetailOfficialAccount 用户公众号信息
type DetailOfficialAccount struct {
	UserID        int64  `json:"user_id" db:"user_id"`
	Name          string `json:"name" db:"name"`
	SubscribeTime int64  `json:"subscribe_time" db:"subscribe_time"`
}

// RequestCustomer 客户详情参数
type RequestCustomer struct {
	UserID   int64  `json:"user_id" db:"user_id" validate:"required" required:"客户ID必填"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	SortKey  string `json:"sort_key"`
	SortType string `json:"sort_type"`
}

// RequestCustomerGroupChat 客户群请求
type RequestCustomerGroupChat struct {
	RequestCustomer
	Status       int   `json:"status"`
	EnterpriseID int64 `json:"enterprise_id"`
}
