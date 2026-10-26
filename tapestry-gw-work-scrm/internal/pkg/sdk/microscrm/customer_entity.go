package microScrm

import "gw-work-scrm/internal/pkg/sdk/work"

// CustomerRequest 客户画像请求参数
type CustomerRequest struct {
	UserID int64 `json:"user_id" validate:"required" required:"客户ID不能为空"`
	Page   int   `json:"page"`
	Size   int   `json:"size"`
	Status int   `json:"status"`
	RequestBase
}

type CustomerDetailRequest struct {
	RequestBase
	ExternalUserID string `json:"external_userid" validate:"required" required:"外部联系人userid不能为空"`
	Mobile         string `json:"mobile"`
	MemberID       int64  `json:"member_id"`
	MemberUserID   string `json:"member_userid"`
	MemberType     int    `json:"member_type"  comment:"当前顾问类型"`
	ThirdIds       []int  `json:"third_ids" binding:"required_if=MemberType 4" required_if:"third_ids不能为空" comment:"当前SM管理的门店"`
}

type CustomerCouponsRequest struct {
	ConsumerID string `json:"consumerID"`
	TransType  string `json:"transType"`
	RequestBase
}

type CustomerDLOrderListRequest struct {
	UnionID string `json:"union_id" validate:"required" required:"外部联系人union_id不能为空"`
}

type CustomerDLPointsDetailRequest struct {
	UnionID        string `json:"union_id"`
	PointTypeGroup string `json:"point_type_group"`
	PageNo         string `json:"page_no"`
	Limit          string `json:"limit"`
}

type CustomerCrmPointsRequest struct {
	UnionID        string `json:"union_id"`
	PointTypeGroup string `json:"point_type_group"`
	MemberID       int64  `json:"member_id"`
	Page           int    `json:"page"`
	Size           int    `json:"size"`
	StartTime      string `json:"startTime"`
	EndTime        string `json:"endTime"`
	RequestBase
}

type CustomerCrmCouponsRequest struct {
	UnionID     string `json:"union_id"`
	ChannelType string `json:"channel_type"`
	Status      int    `json:"status"`
	Source      int    `json:"source"`
	Page        int    `json:"page"`
	Size        int    `json:"size"`
	RequestBase
}

type CustomerDLOrderGoodsRequest struct {
	CustomerDLOrderListRequest
	TradeNo string `json:"trade_no"`
}

type CustomerCrmOrderRequest struct {
	UnionID       string `json:"union_id" validate:"required" required:"外部联系人 union_id 不能为空"`
	StartDateTime int    `json:"startDateTime" required:"startDateTime 不能为空"`
	EndDateTime   int    `json:"endDateTime" required:"endDateTime 不能为空"`
}

// GroupChatResponse 返回结构体
type GroupChatResponse struct {
	ResponseBase
	Data GroupChatList `json:"data"`
}

// GroupChatList 列表结构
type GroupChatList struct {
	ListBase
	ListData []BelongToGroupChat `json:"list"`
}

// BelongToGroupChat 群聊
type BelongToGroupChat struct {
	GroupChatId   int64                 `json:"group_chat_id"`
	GroupChatName string                `json:"group_chat_name"`
	MemberId      int64                 `json:"member_id"`
	MemberName    string                `json:"member_name"`
	UserId        string                `json:"user_id"`
	JoinTime      int64                 `json:"join_time"`
	MemberIsJoin  int                   `json:"member_is_join"`
	GroupMember   []WorkGroupChatMember `json:"group_member"`
}

// WorkGroupChatMember 群成员
type WorkGroupChatMember struct {
	ID           int64  `json:"id" desc:"主键"`
	GroupChatID  int64  `json:"group_chat_id"`
	Type         int    `json:"type"`
	UserID       string `json:"userid"`
	JoinTime     int64  `json:"join_time"`
	JoinScene    int    `json:"join_scene"`
	LeaveTime    int64  `json:"leave_time"`
	Status       int    `json:"status"`
	MemberID     int64  `json:"member_id"`
	UnionID      string `json:"unionid"`
	EnterpriseID int64  `json:"enterprise_id"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
	DeletedAt    int64  `json:"deleted_at"`
}

// RemarkRequest 客户备注信息
type RemarkInfo struct {
	UserID       int64  `json:"user_id" db:"user_id" desc:"客户userID" validate:"required" required:"客户ID不能为空"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Age          int    `json:"age" db:"age" desc:"年龄"`
	Birthday     int64  `json:"birthday" db:"birthday" desc:"生日时间戳"`
	BirthdayDate int    `json:"birthday_date" db:"birthday_date" desc:"生日格式:20200908"`
	Gender       int    `json:"gender" db:"gender" desc:"性别：0-未知 1-男性 2-女性"`
	Mobile       string `json:"mobile" db:"mobile" desc:"电话"`
	Email        string `json:"email" db:"email" desc:"邮箱"`
	Qq           string `json:"qq" db:"qq" desc:"QQ号"`
	Describe     string `json:"describe" db:"describe" desc:"描述信息"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"添加时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// ReturnUID 更新客户备注信息返回数据
type ReturnUID struct {
	UserID int64 `json:"user_id" db:"user_id"`
}

// RequestOrderList 订单列表请求数据(侧边栏)
type RequestOrderList struct {
	UserID       int64  `json:"user_id" db:"user_id" validate:"required" required:"客户ID必填"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	OrderNo      int    `json:"order_no" db:"order_no"`
	Status       int    `json:"status" db:"status"`
	ShowCancel   int    `json:"show_cancel" usage:"是否展示未付款的订单 0=不展示（默认） 1=展示"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
	SortKey      string `json:"sort_key"`
	SortType     string `json:"sort_type"`
}

// ResponseOrderList 订单列表
type ResponseOrderList []struct {
	ID             int64                  `json:"id" db:"id"`
	OrderNo        string                 `json:"order_no" db:"order_no"`
	Status         int                    `json:"status" db:"status" usage:"订单状态"`
	FinalAmount    int64                  `json:"final_amount" db:"final_amount" usage:"含券支付金额"`
	CashFee        int64                  `json:"cash_fee" db:"cash_fee" usage:"实际支付金额"`
	TotalRefundFee int64                  `json:"total_refund_fee" db:"total_refund_fee" usage:"退款金额"`
	SellPrice      int                    `json:"sell_price" db:"sell_price"`
	GoodsPayPrice  int                    `json:"goods_pay_price" db:"goods_pay_price"`
	PayAt          int64                  `json:"pay_at" db:"pay_at"`
	RefundAt       int64                  `json:"refund_at" db:"refund_at"`
	GoodsList      GoodsList              `json:"goods_list" usage:"商品列表"`
	CouponList     CouponList             `json:"coupon_list" usage:"优惠券列表"`
	UserInfo       *OrderListCustomerInfo `json:"user_info"`
}

// OrderListCustomerInfo 订单列表用户信息
type OrderListCustomerInfo struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// GoodsList 商品列表
type GoodsList []struct {
	ID               int64                    `json:"id" db:"id"`
	ShopOrderID      int64                    `json:"shop_order_id" db:"shop_order_id"`
	SpuName          string                   `json:"spu_name" db:"spu_name"`
	SkuPicture       string                   `json:"sku_picture" db:"sku_picture"`
	Quantity         int                      `json:"quantity" db:"quantity"`
	SellPrice        int                      `json:"sell_price" db:"sell_price"`
	GoodsPayPrice    int                      `json:"goods_pay_price" db:"goods_pay_price"`
	SkuId            int64                    `json:"sku_id" db:"sku_id"`
	SpuId            int64                    `json:"spu_id" db:"spu_id"`
	SkuAttribute     string                   `json:"sku_attribute" db:"sku_attribute"`
	SkuAttributeList []map[string]interface{} `json:"sku_attribute_list"`
}

// CouponList 优惠券列表
type CouponList []struct {
	ShopOrderID int64  `json:"shop_order_id"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	StockId     string `json:"stock_id"`
}

// RequestOrderDetail 订单详情请求数据
type RequestOrderDetail struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
	OrderNo      int64 `json:"order_no" db:"order_no"`
}

// ResponseOrderDetail 订单详情返回结果
type ResponseOrderDetail struct {
	ID             int64          `json:"id" db:"id"`
	OrderNo        int64          `json:"order_no" db:"order_no"`
	Status         int            `json:"status" db:"status" usage:"订单状态"`
	PaymentChannel int            `json:"payment_channel" db:"payment_channel" usage:"支付渠道"`
	PayAt          int64          `json:"pay_at" db:"pay_at"`
	ThirdPaymentNo string         `json:"third_payment_no" db:"third_payment_no" usage:"第三方支付流水号"`
	GoodsPayPrice  int            `json:"goods_pay_price" db:"goods_pay_price"`
	LogisticsPrice int            `json:"logistics_price" db:"logistics_price" usage:"物流邮费"`
	DiscountFee    int            `json:"discount_fee" db:"discount_fee" usage:"优惠总金额"`
	CouponFee      int            `json:"coupon_fee" db:"coupon_fee" usage:"微信优惠金额"`
	OtherFee       int            `json:"other_fee" db:"other_fee" usage:"微信优惠金额"`
	TotalAmount    int            `json:"total_amount" db:"total_amount" usage:"原价"`
	CashFee        int            `json:"cash_fee" db:"cash_fee" usage:"实际支付金额"`
	TotalRefundFee int64          `json:"total_refund_fee" db:"total_refund_fee" usage:"退款金额"`
	OrderAt        int64          `json:"order_at" db:"order_at"`
	GoodsList      []GoodsEntity  `json:"goods_list" usage:"商品列表"`
	DeliveryInfo   DeliveryEntity `json:"delivery_info"`
}

// GoodsEntity 商品实体
type GoodsEntity struct {
	ID               int64                    `json:"id" db:"id"`
	ShopOrderID      int64                    `json:"shop_order_id" db:"shop_order_id"`
	SpuName          string                   `json:"spu_name" db:"spu_name"`
	SkuPicture       string                   `json:"sku_picture" db:"sku_picture"`
	Quantity         int                      `json:"quantity" db:"quantity"`
	SellPrice        int                      `json:"sell_price" db:"sell_price"`
	GoodsPayPrice    int                      `json:"goods_pay_price" db:"goods_pay_price"`
	GoodsType        int                      `json:"goods_type" db:"goods_type"`
	SkuId            int64                    `json:"sku_id" db:"sku_id"`
	SpuId            int64                    `json:"spu_id" db:"spu_id"`
	SkuAttribute     string                   `json:"sku_attribute" db:"sku_attribute"`
	SkuAttributeList []map[string]interface{} `json:"sku_attribute_list"`
}

// RequestDelivery 收货人信息请求参数
type RequestDelivery struct {
	ShopOrderID int64 `json:"shop_order_id" db:"shop_order_id"`
}

// DeliveryEntity 收货人信息返回参数
type DeliveryEntity struct {
	Name     string `json:"name" db:"name"`
	Mobile   string `json:"mobile" db:"mobile"`
	Province string `json:"province" db:"province"`
	City     string `json:"city" db:"city"`
	Area     string `json:"area" db:"area"`
	Address  string `json:"address" db:"address"`
}

// SetExternalTags 设置外部联系人企业标签请求参数
type SetExternalTags struct {
	EnterpriseID int64    `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	RelationID   int64    `json:"relation_id" db:"relation_id" desc:"work_external_user_relation表id"`
	AddTag       []string `json:"add_tag" db:"external_tag_id" desc:"新增的企微标签IDS"`
	RemoveTag    []string `json:"remove_tag" db:"external_tag_id" desc:"删除的企微标签IDS"`
	CreatedAt    int64    `json:"created_at" db:"created_at"`
	UpdatedAt    int64    `json:"updated_at" db:"updated_at"`
	DeletedAt    int64    `json:"deleted_at" db:"deleted_at"`
}

// PortraitResponse 客户画像详情返回结果
type PortraitResponse struct {
	UserID                int64                  `json:"user_id" db:"user_id"`
	UnionID               string                 `json:"union_id" db:"union_id"`
	ExternalUserID        string                 `json:"external_userid" db:"external_userid"`
	RelationID            int64                  `json:"relation_id" db:"relation_id"`
	NickName              string                 `json:"nick_name" db:"nick_name"`
	Avatar                string                 `json:"avatar" db:"avatar"`
	CreateTime            int64                  `json:"createtime" db:"createtime"`
	WorkExternalUserID    int64                  `json:"work_external_user_id" db:"work_external_user_id"`
	BindGuideName         string                 `json:"bind_guide_name" db:"bind_guide_name"`
	ChannelName           string                 `json:"channel_name"`
	GroupName             []string               `json:"group_name"`
	GroupNumber           int64                  `json:"group_number"`
	TagNames              []string               `json:"tag_names"`
	BindGuideMemberID     int64                  `json:"bind_guide_member_id" db:"bind_member_id"`
	BindGuideMemberUserid string                 `json:"bind_guide_member_userid"`
	Age                   int                    `json:"age" db:"age"`
	QQ                    string                 `json:"qq" db:"qq"`
	Email                 string                 `json:"email" db:"email"`
	Describe              string                 `json:"describe" db:"describe"`
	Gender                int                    `json:"gender" db:"gender"`
	GenderIndex           int                    `json:"gender_index"`
	UserOrder             UserOrderStatistics    `json:"user_order"`
	CustomerMember        CustomerMemberRelation `json:"customer_member"`
	EvalidPoints          int64                  `json:"evalid_points"`
	ExpiredPoints         int64                  `json:"expired_points"`
	OrderList             []work.Orders          `json:"order_list"`
	Grade                 string                 `json:"grade"`
	NextGrade             string                 `json:"nextGrade"`
	GradeCode             string                 `json:"gradeCode"`
	NextGradeCode         string                 `json:"nextGradeCode"`
	Birthday              int64                  `json:"birthday" db:"birthday" desc:"生日时间戳"`
	Name                  string                 `json:"name"`
	Mobile                string                 `json:"mobile" db:"mobile"`                                     // 电话
	BirthdayStr           string                 `json:"birthday_str"`                                           // 生日
	Location              string                 `json:"location"`                                               // 地址
	CustomerCode          string                 `json:"customer_code" db:"customer_code"`                       // 会员卡号
	CustomerType          string                 `json:"customer_type" db:"customer_type"`                       // 会员等级
	Status                string                 `json:"status"`                                                 // 会员状态
	ExpiryDate            string                 `json:"expiry_date"`                                            // 会员到期时间
	UpgradeDateTime       string                 `json:"upgrade_date_time" db:"upgrade_date_time"`               // 会员最近升级时间
	UpgradeNeedPoints     int                    `json:"upgrade_need_points" db:"upgrade_need_points"`           // 会员升级所需积分
	Last12mPurchaseAmount float32                `json:"last_12m_purchase_amount" db:"last_12m_purchase_amount"` // 过去12个月消费总额
	Last36mPurchaseAmount float32                `json:"last_36m_purchase_amount" db:"last_36m_purchase_amount"` // 过去36个月消费总额

	// UpGradeInfo UpGradeInfo //升级信息
}

type UpGradeInfo struct {
	UpGradeTxt string `json:"up_grade_txt"`
	Percentage int64  `json:"percentage"`
}

type CustomerMemberRelation struct {
	CustomerMemberList []GroupUsersAddEmployees `json:"customer_member_list"`
	MemberCount        int64                    `json:"member_count"`
}

// UserOrderStatistics 订单统计的
type UserOrderStatistics struct {
	NumberOfOrders   int `json:"number_of_orders" desc:"订单成交数"`
	NumberOfRefunded int `json:"number_of_refunded" desc:"退款订单数"`
	OrdersAmount     int `json:"orders_amount" desc:"实收金额"`
	RefundAmount     int `json:"refund_amount" desc:"退款金额"`
}

// CustomerFootPrintRequest 客户画像请求参数
type CustomerFootPrintRequest struct {
	UserID   int64 `json:"user_id" validate:"required" required:"客户ID不能为空"`
	MemberID int64 `json:"member_id"`
	Page     int   `json:"page"`
	Size     int   `json:"size"`
	RequestBase
}

// PortraitDetailResponse 客户画像详情返回结果
type PortraitDetailResponse struct {
	UserID                int64                    `json:"user_id" db:"user_id"`
	UnionID               string                   `json:"union_id" db:"union_id"`
	ExternalUserID        string                   `json:"external_userid" db:"external_userid"`
	RelationID            int64                    `json:"relation_id" db:"relation_id"`
	NickName              string                   `json:"nick_name" db:"nick_name"`
	Remark                string                   `json:"remark" db:"remark"`
	Avatar                string                   `json:"avatar" db:"avatar"`
	CreateTime            int64                    `json:"createtime" db:"createtime"`
	WorkExternalUserID    int64                    `json:"work_external_user_id" db:"work_external_user_id"`
	GroupName             []string                 `json:"group_name"`
	GroupNumber           int64                    `json:"group_number"`
	Tags                  []string                 `json:"tags"`
	UserOrder             UserOrderStatistics      `json:"user_order"`
	CustomerMember        CustomerMemberRelation   `json:"customer_member"`
	OrderList             []work.Orders            `json:"order_list"`
	Birthday              int64                    `json:"birthday" db:"birthday" desc:"生日时间戳"`
	Sex                   int                      `json:"sex"`
	Name                  string                   `json:"name"`
	Mobile                string                   `json:"mobile" db:"mobile"`                                     // 电话
	BirthdayStr           string                   `json:"birthday_str"`                                           // 生日
	Location              string                   `json:"location"`                                               // 地址
	CustomerCode          string                   `json:"customer_code" db:"customer_code"`                       // 会员卡号
	Membership            string                   `json:"membership" db:"membership"`                             // 主会员卡号
	CustomerType          string                   `json:"customer_type" db:"customer_type"`                       // 会员等级
	Status                string                   `json:"status"`                                                 // 会员状态
	ExpiryDate            string                   `json:"expiry_date"`                                            // 会员到期时间
	UpgradeDateTime       string                   `json:"upgrade_date_time" db:"upgrade_date_time"`               // 会员最近升级时间
	UpgradeNeedPoints     float32                  `json:"upgrade_need_points" db:"upgrade_need_points"`           // 会员升级所需积分
	Last12mPurchaseAmount float32                  `json:"last_12m_purchase_amount" db:"last_12m_purchase_amount"` // 过去12个月消费总额
	Last36mPurchaseAmount float32                  `json:"last_36m_purchase_amount" db:"last_36m_purchase_amount"` // 过去36个月消费总额
	TotalValidPoints      float32                  `json:"total_valid_points" db:"total_valid_points"`             // 积分
	ClientType            string                   `json:"client_type" db:"client_type"`                           // CTL 类型
	Address               string                   `json:"address" db:"address"`
	UserTags              []UserTagsInfo           `json:"user_tags" db:"user_tags"`                       // 客户标签
	UserImpressionTags    []UserImpressionTagsInfo `json:"user_impression_tags" db:"user_impression_tags"` // 客户印象标签
	UserInfoAdd           UserInfoAdd              `json:"user_info_add"`
}

type UserTagsRequest struct {
	EnterpriseID       int64  `json:"enterprise_id"`
	MemberID           int64  `json:"member_id"`
	WorkExternalUserID int64  `json:"work_external_user_id" db:"work_external_user_id"`
	Userid             string `json:"userid"`
}

type UserTagsRes struct {
	ResponseBase
	Data []UserTagsInfo `json:"data"`
}
type UserTagsInfo struct {
	RelationID    int64      `json:"relation_id"`
	ExternalTagID int64      `json:"external_tag_id"`
	Type          int        `json:"type"`
	GroupID       int64      `json:"group_id"`
	GroupName     string     `json:"group_name"`
	TagsList      []TagsList `json:"tags_list"`
}

type TagsList struct {
	TagID               string `json:"tag_id"`
	Name                string `json:"name"`
	ThirdTagDescription string `json:"third_tag_description"`
}

type UserImpressionTagsRes struct {
	ResponseBase
	Data []UserImpressionTagsInfo `json:"data"`
}
type UserImpressionTagsInfo struct {
	GroupID            int64                `json:"group_id"`
	GroupName          string               `json:"group_name"`
	ImpressionTagsList []ImpressionTagsList `json:"impression_tags_list"`
}
type ImpressionTagsList struct {
	TagID  int64  `json:"tag_id"`
	Name   string `json:"name"`
	IsHave int    `json:"is_have"`
}
type UserInfoAdd struct {
	Address  string `json:"address"`
	ShoeSize string `json:"shoe_size"`
	Month    int    `json:"month"`
	Day      int    `json:"day"`
}
type UserInfoAddRes struct {
	ResponseBase
	Data UserInfoAdd `json:"data"`
}

type UpdateRemarkRequest struct {
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	MemberID       int64  `json:"member_id" `
	ExternalUserID string `json:"external_userid"  validate:"required"  db:"external_userid"`
	Mobile         string `json:"mobile"`
	Remark         string `json:"remark"  validate:"required" `
}

// CrmUserDetail h5 客户画像详情返回结果
type CrmUserDetail struct {
	UserID                int64                    `json:"user_id"`
	UnionID               string                   `json:"union_id"`
	ExternalUserID        string                   `json:"external_userid"`
	WorkExternalUserID    int64                    `json:"work_external_user_id"`
	NickName              string                   `json:"nick_name"`
	Name                  string                   `json:"name"`
	Avatar                string                   `json:"avatar"`
	Mobile                string                   `json:"mobile"`
	Gender                int                      `json:"gender"`
	RelationID            int64                    `json:"relation_id"`
	Remark                string                   `json:"remark"`
	IsRelation            int                      `json:"is_relation"`
	BirthdayStr           string                   `json:"birthday_str"`             // 生日
	CustomerCode          string                   `json:"customer_code"`            // 会员卡号
	CustomerCard          string                   `json:"customer_card"`            // 会员卡等类型
	CustomerCardStr       string                   `json:"customer_card_str"`        // 会员卡等类型
	Membership            string                   `json:"membership"`               // 主会员卡号
	ExpiryDate            string                   `json:"expiry_date"`              // 会员到期时间
	UpgradeDateTime       string                   `json:"upgrade_date_time"`        // 会员最近升级时间
	UpgradeNeedPoints     float32                  `json:"upgrade_need_points"`      // 会员升级所需积分
	Last12mPurchaseAmount float32                  `json:"last_12m_purchase_amount"` // 过去12个月消费总额
	Last36mPurchaseAmount float32                  `json:"last_36m_purchase_amount"` // 过去36个月消费总额
	TotalValidPoints      float32                  `json:"total_valid_points"`       // 积分
	ClientType            string                   `json:"client_type"`              // CTL 类型
	Address               string                   `json:"address"`
	AmountTotal           int64                    `json:"amount_total" `
	LastPurchaseDate      int                      `json:"last_purchase_date"`
	LastPurchaseDateStr   string                   `json:"last_purchase_date_str"`
	LastPurchaseStoreCode string                   `json:"last_purchase_store_code"`
	LastPurchaseStoreName string                   `json:"last_purchase_store_name"`
	UpgradeNeedAmount     int                      `json:"upgrade_need_amount"`
	CustomerType          int                      `json:"customer_type"`
	CrmTag                []string                 `json:"crm_tag" comment:"CRM标签"`
	UserTags              []UserTagsInfo           `json:"user_tags"`            // 客户标签
	UserImpressionTags    []UserImpressionTagsInfo `json:"user_impression_tags"` // 客户印象标签
	VipStatus             CrmVipStatus             `json:"VipStatus"`
	UserInfoAdd           UserInfoAdd              `json:"user_info_add"`
}

type CrmVipStatus struct {
	IssueDate    string `json:"issueDate"`
	IssueShop    string `json:"issueShop"`
	PtsToUpgrade string `json:"ptsToUpgrade"`
	ExpiryDate   string `json:"expiryDate"`
}

type CrmUserDetailRes struct {
	ResponseBase
	Data CrmUserDetail `json:"data"`
}

type UserDetailTipsNum struct {
	ReserveNum int `json:"reserve_num"`
	SowTaskNum int `json:"sow_task_num"`
}

type UserDetailTipsNumRes struct {
	ResponseBase
	Data UserDetailTipsNum `json:"data"`
}

type UserSaListData struct {
	MemberId       int64  `json:"member_id"`
	Unionid        string `json:"unionid"`
	Userid         string `json:"userid"`
	ThumbAvatar    string `json:"thumb_avatar"`
	QrCode         string `json:"qr_code"`
	IsKeySa        int    `json:"is_key_sa" db:"is_key_sa"` // 是否专属导购：0否，1是
	SaUserId       string `json:"sa_user_id" db:"sa_user_id"`
	SaUserName     string `json:"sa_user_name" db:"sa_user_name"`
	ThirdId        int    `json:"third_id"`
	DepartmentName string `json:"department_name"`
	CreatedAt      int    `json:"created_at"`
}

type UserSaListDataRes struct {
	ResponseBase
	Data []UserSaListData `json:"data"`
}
