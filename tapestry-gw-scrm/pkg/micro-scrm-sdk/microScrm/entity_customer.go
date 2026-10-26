package microScrm

// RequestCustomer 客户请求
type RequestCustomer struct {
	UserID       int64  `json:"user_id" db:"user_id" binding:"required" required:"客户ID必填"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必填"`
	Page         int    `json:"page" binding:"min=1" min:"页面参数不能小于1"`
	Size         int    `json:"size" binding:"min=1" min:"分页参数不能小于1"`
	SortKey      string `json:"sort_key"`
	SortType     string `json:"sort_type"`
	ChannelType  int8   `json:"channel_type"`
}

// RequestCustomerGroup 客户群请求
type RequestCustomerGroup struct {
	RequestCustomer
	Status int `json:"status" db:"status" desc:"群状态，1=入群 2=退群"`
}

// ResponseFootprintList 足迹
type ResponseFootprintList []struct {
	ID        int64  `json:"id" db:"id"`
	EventTxt  string `json:"event_txt" db:"event_txt"`
	EventDesc string `json:"event_desc" db:"event_desc"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
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

// ResponseAddressList 收货人地址列表
type ResponseAddressList []struct {
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
	UserID       int64                    `json:"user_id" db:"user_id"`
	UserName     string                   `json:"user_name" db:"name"`
	Nickname     string                   `json:"nick_name" db:"nick"`
	Avatar       string                   `json:"avatar" db:"avatar"`
	Gender       int                      `json:"gender" db:"gender"`
	Country      string                   `json:"country" db:"country"`
	Province     string                   `json:"province" db:"province"`
	City         string                   `json:"city" db:"city"`
	Remark       string                   `json:"remark" db:"remark"`
	RemarkList   []string                 `json:"remark_list"`
	TagList      []string                 `json:"tag_list"`
	GroupList    []string                 `json:"group_list"`
	IsWorkUser   bool                     `json:"is_work_user"`
	OfficialList []*DetailOfficialAccount `json:"official_list"`
	ShowOfficial bool                     `json:"show_official"`
	LevelName    string                   `json:"level_name"`
	Birthday     string                   `json:"birthday" db:"birthday"`
	UnionCode    string                   `json:"union_code" db:"union_code"`
	Point        string                   `json:"point" db:"point"`
	Grade        string                   `json:"grade" db:"grade"`
	StoreCode    string                   `json:"store_code" db:"store_code"`
	BACode       string                   `json:"ba_code" db:"ba_code"`
	Mobile       string                   `json:"mobile" db:"mobile"`
}

// DetailOfficialAccount 用户公众号信息
type DetailOfficialAccount struct {
	UserID        int64  `json:"user_id" db:"user_id"`
	Name          string `json:"name" db:"name"`
	SubscribeTime int64  `json:"subscribe_time" db:"subscribe_time"`
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

type ResponseMemberList []struct {
	MemberId               int64                 `json:"member_id" db:"member_id"`
	MemberName             string                `json:"member_name" db:"member_name"`
	DepartmentList         MemberDepartsList     `json:"department_list"`
	CreateTime             int64                 `json:"create_time" db:"create_time" usage:"关系建立时间"`
	State                  string                `json:"state" json:"state" usage:"来源"`
	BelongToGroupChatCount int64                 `json:"belong_to_group_chat_count" usage:"当前客户占员工群总数的个数"`
	BelongToGroupChatList  BelongToGroupChatList `json:"belong_to_group_chat_list" usage:"当前客户占员工群列表"`
	RemarkList             []string              `json:"remark_list" usage:"备注"`
}

type MemberDeparts struct {
	MemberID       int64  `json:"member_id" db:"member_id"`
	DepartmentID   int64  `json:"department_id" db:"department_id"`
	DepartmentName string `json:"department_name" db:"department_name"`
}

type MemberDepartsList []MemberDeparts

type BelongToGroupChat struct {
	MemberId      int64  `json:"member_id" db:"member_id"`
	MemberName    string `json:"member_name" db:"member_name"`
	GroupChatId   int64  `json:"group_chat_id" db:"group_chat_id"`
	GroupChatName string `json:"group_chat_name" db:"group_chat_name"`
	UserId        string `json:"user_id" db:"user_id"`
	JoinTime      int64  `json:"join_time" db:"join_time"`
}

type BelongToGroupChatList []BelongToGroupChat

// RemarkRequest 客户备注信息
type RemarkInfo struct {
	UserID       int64  `json:"user_id" db:"user_id" desc:"客户userID" validate:"required" required:"客户ID不能为空"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" validate:"required" required:"企业ID不能为空"`
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
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必填"`
	OrderNo      int    `json:"order_no" db:"order_no"`
	Status       int    `json:"status" db:"status"`
	ShowCancel   int    `json:"show_cancel" usage:"是否展示未付款的订单 0=不展示（默认） 1=展示"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
	SortKey      string `json:"sort_key"`
	SortType     string `json:"sort_type"`
}

// RequestOrderDetail 订单详情请求数据
type RequestOrderDetail struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必填"`
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

// ReturnID 设置客户企业标签返回
type ReturnID struct {
	ID int64 `json:"id"`
}
