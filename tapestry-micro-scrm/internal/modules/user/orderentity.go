package user

import (
	"encoding/json"
	"micro-scrm/internal/pkg/common"
)

const (
	OrderSyncTypePay    = 1 // 订单同步类型：支付订单
	OrderSyncTypeRefund = 2 // 订单同步类型：退款订单
)

// ShopOrderEntity 订单数据结构实例
type ShopOrderEntity struct {
	ID                  int64  `db:"id"`
	EnterpriseID        int64  `json:"enterpriseId" db:"enterprise_id"`            // EnterpriseID 企业id
	ShopUserID          int64  `json:"shop_user_id" db:"shop_user_id"`             // ShopUserID 会员ID
	UnionID             string `json:"unionId" db:"union_id"`                      // UnionID 用户union_id
	OpenID              string `json:"openId" db:"open_id"`                        // OpenID 用户open_id
	OrderNo             string `json:"orderNo" db:"order_no"`                      // OrderNo 订单号
	Status              int    `json:"status" db:"status"`                         // Status 订单状态：1 待付款，2 已完成，3 已取消, 4 订单确认中
	OrderType           int    `json:"order_type" db:"order_type"`                 // OrderType 订单类型: 1 线上订单，2 门店订单
	PaymentChannel      int    `json:"paymentChannel" db:"payment_channel"`        // PaymentChannel 支付渠道：1 微信小程序支付，2 美团支付
	TotalAmount         int64  `json:"totalAmount" db:"total_amount"`              // TotalAmount 原价 单位：分
	PreferentialPrice   int    `json:"preferential_price" db:"preferential_price"` // PreferentialPrice 优惠金额，单位：分
	DiscountFee         int    `json:"discountFee" db:"discount_fee"`              // DiscountFee 优惠总金额
	FinalAmount         int    `json:"finalAmount" db:"final_amount"`              // FinalAmount 实际付款价格，单位：分
	PayAt               int    `json:"payAt" db:"pay_at"`                          // PayAt 支付时间
	RefundAt            int    `json:"refundAt" db:"refund_at"`                    // RefundAt 退款时间
	OrderAt             int    `json:"createdAt" db:"order_at"`                    // OrderAt 下单时间
	Remarks             string `json:"remarks" db:"remarks"`                       // Remarks 卖家备注
	CashFee             int    `json:"cashFee" db:"cash_fee"`                      // CashFee 实际支付金额，单位：分
	CouponFee           int    `json:"couponFee" db:"coupon_fee"`                  // CouponFee 微信优惠，单位：分
	CouponCount         int    `json:"couponCount" db:"coupon_count"`              // CouponCount 微信优惠券数量
	OtherFee            int    `json:"otherFee" db:"other_fee"`                    // OtherFee 其他优惠，单位：分
	TotalRefundFee      int    `json:"totalRefundFee" db:"total_refund_fee"`       // TotalRefundFee 实际退款累计总额（包含微信优惠退回），单位：分
	ReadStatus          int    `json:"read_status" db:"read_status"`               // ReadStatus 已读状态 1 未读 2 已读
	LogisticsPrice      int    `json:"logisticsPrice" db:"logistics_price"`        // LogisticsPrice 物流邮费
	ThirdPaymentNo      string `json:"thirdPaymentNo" db:"third_payment_no"`       // ThirdPaymentNo 第三方支付流水号
	OrderSource         int    `json:"order_source" db:"order_source"`             // OrderSource 订单来源：1-crs商城
	SendAt              int    `json:"sendAt" db:"send_at"`                        // SendAt 发货时间
	SignAt              int    `json:"signAt" db:"sign_at"`                        // SignAt 签收时间
	CreatedAt           int    `json:"created_at" db:"created_at"`                 // CreatedAt 创建时间
	UpdatedAt           int    `json:"updated_at" db:"updated_at"`                 // UpdatedAt 更新时间
	DeletedAt           int    `json:"deleted_at" db:"deleted_at"`                 // DeletedAt 删除时间
	PreviousCashFee     int    // PreviousCashFee 更新数据之前的付款金额
	PreviousTotalRefund int    // PreviousTotalRefund 更新数据之前的退款金额
}

// ShopOrderCouponEntity 订单优惠券结构实例
type ShopOrderCouponEntity struct {
	ID           int64  `db:"id"`
	EnterpriseID int64  `db:"enterprise_id"`                 // EnterpriseID 企业id
	ShopOrderID  int64  `db:"shop_order_id"`                 // ShopOrderID shop_order表id
	Type         int    `db:"type"`                          // Type 优惠券类型：1-微信优惠券
	StockID      string `db:"stock_id" json:"couponStockId"` // StockID 优惠券批次号
	Name         string `db:"name" json:"couponName"`        // Name 优惠券名称
	Value        int    `db:"value" json:"couponValue"`      // Value 优惠金额（分）
	Minimum      int    `db:"minimum" json:"couponMininumn"` // Minimum 使用限制金额（分）
	CreatedAt    int    `db:"created_at"`
	UpdatedAt    int    `db:"updated_at"`
	DeletedAt    int    `db:"deleted_at"`
}

// ShopOrderDeliveryEntity 商城订单发货信息结构实例
type ShopOrderDeliveryEntity struct {
	ID           int64  `db:"id"`
	EnterpriseID int64  `db:"enterprise_id"`                  // EnterpriseID 企业id
	ShopOrderID  int64  `db:"shop_order_id"`                  // ShopOrderID shop_order表id
	Name         string `json:"deliveryName" db:"name"`       // Name 收货人姓名
	Mobile       string `json:"deliveryMobile" db:"mobile"`   // Mobile 收货人手机号码
	Address      string `json:"deliveryAddress" db:"address"` // Address 收货人地址
	Province     string `json:"province" db:"province"`       // Province 省
	City         string `json:"city" db:"city"`               // City 市
	Area         string `json:"area" db:"area"`               // Area 区
	CreatedAt    int    `db:"created_at"`                     // CreatedAt 创建时间
	UpdatedAt    int    `db:"updated_at"`                     // UpdatedAt 更新时间
	DeletedAt    int    `db:"deleted_at"`                     // DeletedAt 删除时间
}

// ShopOrderGoodEntity 商城订单商品信息实体
type ShopOrderGoodEntity struct {
	ID             int64  `db:"id"`
	EnterpriseID   int64  `db:"enterprise_id"` // EnterpriseID 企业id
	ShopOrderID    int64  `db:"shop_order_id"` // ShopOrderID shop_order表id
	OrderNO        string `db:"order_no"`
	SkuNo          string `json:"skuNo" db:"sku_no"`                   // SkuNo sku编号
	SpuID          int64  `json:"spuId" db:"spu_id"`                   // SpuID SPU ID
	SkuID          int64  `json:"skuId" db:"sku_id"`                   // SkuID SKU ID
	SpuName        string `json:"spuName" db:"spu_name"`               // SpuName SPU商品名称
	Type           int    `json:"type" db:"type"`                      // Type 订单商品类型：1普通商品，2拼团，3秒杀，5超低购，6助力，7砍价，11新人购商品，12限时购商品
	GoodsType      int    `json:"goodsType" db:"goods_type"`           // GoodsType 商品类型：1普通，2赠品，3奖品
	CategoryName   string `json:"categoryName" db:"category_name"`     // CategoryName 分类名称
	SkuPicture     string `json:"skuPicture" db:"sku_picture"`         // SkuPicture SKU主图URL
	SkuAttribute   string `json:"skuAttribute" db:"sku_attribute"`     // SkuAttribute SKU商品属性值集合
	Quantity       int    `json:"quantity" db:"quantity"`              // Quantity 购买数量
	OriginalPrice  int    `json:"originalPrice" db:"original_price"`   // OriginalPrice 商品原价，单位：分
	SellPrice      int    `json:"sellPrice" db:"sell_price"`           // SellPrice 销售单价，单位：分
	SingleDiscount int    `json:"singleDiscount" db:"single_discount"` // SingleDiscount 单品优惠，单位：分
	SubPrice       int    `json:"subPrice" db:"sub_price"`             // SubPrice 总价，单位：分
	GoodsPayPrice  int    `json:"goodsPayPrice" db:"goods_pay_price"`  // GoodsPayPrice 订单商品分摊实付金额
	CreatedAt      int    `json:"created_at" db:"created_at"`          // CreatedAt 创建时间
	UpdatedAt      int    `json:"updated_at" db:"updated_at"`          // UpdatedAt 更新时间
	DeletedAt      int    `json:"deleted_at" db:"deleted_at"`          // DeletedAt 删除时间
}

// ShopOrderLogisticEntity 商城订单物流信息实体
type ShopOrderLogisticEntity struct {
	ID               int64  `db:"id"`
	EnterpriseID     int64  `db:"enterprise_id"`
	ShopOrderID      int64  `db:"shop_order_id"`     // ShopOrderID shop_order表id
	LogisticIdentity string `db:"logistic_identity"` // LogisticIdentity 物流公司标识
	LogisticNo       string `db:"logistic_no"`       // LogisticNo 物流单号
	CreatedAt        int    `db:"created_at"`        // CreatedAt 创建时间
	UpdatedAt        int    `db:"updated_at"`        // UpdatedAt 更新时间
	DeletedAt        int    `db:"deleted_at"`        // DeletedAt 删除时间
}

// ShopOrderStatusTrace 订单状态追踪实体
type ShopOrderStatusTraceEntity struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `db:"enterprise_id"`
	ShopOrderID  int64  `db:"shop_order_id"`               // ShopOrderID shop_order表id
	Event        string `json:"event" db:"event"`          // Event 订单变化事件
	EventDesc    string `json:"eventDesc" db:"event_desc"` // EventDesc 订单变化事件描述
	From         int    `json:"from" db:"from"`            // From 变化之前状态
	To           int    `json:"to" db:"to"`                // To 变化之后状态
	ChangeAt     int    `json:"createdAt" db:"change_at"`  // ChangeAt 状态变化时间
	CreatedAt    int    `db:"created_at"`                  // CreatedAt 创建时间
	UpdatedAt    int    `db:"updated_at"`                  // UpdatedAt 更新时间
	DeletedAt    int    `db:"deleted_at"`                  // DeletedAt 删除时间
}

type OrderSyncRequest struct {
	Method    common.DdcEvent `json:"method"`
	Platform  string          `json:"platform"`
	OrderData OrderSyncData   `json:"data"`
}

// OrderSyncData 订单同步请求参数结构
type OrderSyncData struct {
	Info struct {
		ShopOrderEntity
		ShopOrderDeliveryEntity
	} `json:"orderInfo"`
	Goods []ShopOrderGoodEntity `json:"orderGoods"`
	User  struct {
		EnterpriseID int64  `json:"enterpriseId"`
		UnionID      string `json:"unionId"`
		AppID        string `json:"appId"`
		OpenId       string `json:"openId"`
		Avatar       string `json:"avatar"`
		Nickname     string `json:"nickName"`
		Mobile       string `json:"mobile"`
		Sex          int    `json:"sex"`
		Country      string `json:"country"`
		Province     string `json:"province"`
		City         string `json:"city"`
		Birthday     string `json:"birthday"`
	} `json:"orderUser"`
	StatusTrace json.RawMessage `json:"orderStatusTrace"`
	Logistic    []struct {
		LogisticIdentity string `json:"logisticIdentity"`
		LogisticNo       string `json:"logisticNo"`
	} `json:"orderLogistic"`
	Coupons []ShopOrderCouponEntity `json:"wxCoupons"`
	PayInfo json.RawMessage         `json:"payInfo"`
}

func (data *OrderSyncData) GetPayInfo() *PayInfo {
	var info PayInfo
	json.Unmarshal(data.PayInfo, &info)
	return &info
}

func (data *OrderSyncData) GetStatusTrace() *ShopOrderStatusTraceEntity {
	var info ShopOrderStatusTraceEntity
	json.Unmarshal(data.StatusTrace, &info)
	return &info
}

type PayInfo struct {
	ThirdPaymentNo string `json:"thirdPaymentNo"`
}

// OrderStatisticsEntity 订单统计
type OrderStatisticsEntity struct {
	ID                     int64 `db:"id"`
	UserID                 int64 `db:"user_id"`
	EnterpriseID           int64 `db:"enterprise_id"`
	PayOrder               int   `db:"pay_order"`
	TotalFee               int   `db:"total_fee"`
	RefundOrder            int   `db:"refund_order"`
	RefundFee              int   `db:"refund_fee"`
	PerCustomerTransaction int   `db:"per_customer_transaction"`
	LatestMarketTime       int   `db:"latest_market_time"`
	MarketCount            int   `db:"market_count"`
	CreatedAt              int   `db:"created_at"`
	UpdatedAt              int   `db:"updated_at"`
	DeletedAt              int   `db:"deleted_at"`
}

// OrderCompleted 成交订单统计
type OrderCompleted struct {
	CompletedOrders       int `json:"completed_orders" db:"completed_orders"`
	CompletedOrdersAmount int `json:"completed_orders_amount" db:"completed_orders_amount"`
}

// RefundOrder 退款订单统计
type RefundOrder struct {
	RefundedOrderNumber int `json:"refunded_order_number" db:"refunded_order_number"`
	RefundAmount        int `json:"refund_amount" db:"refund_amount"`
}
