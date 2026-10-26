package microScrmModel

import "micro-scrm/internal/app/models"

type ExternalPayOrderModel struct {
	models.BaseModel
	TransactionID  string `json:"transaction_id" gorm:"column:transaction_id;type:varchar(255);not null;comment:交易单号"`
	TradeState     int32  `json:"trade_state" gorm:"column:trade_state;type:tinyint(4);default:0;comment:交易状态。1：已完成 3：已完成有退款"`
	OrderState     int32  `json:"order_state" gorm:"column:order_state;type:tinyint(4);default:1;comment:订单状态。1：待发货 2：已发货3：已完成"`
	PayTime        uint   `json:"pay_time" gorm:"column:pay_time;type:int(11) unsigned;index:pay_time_index;default:0;comment:交易时间"`
	OutTradeNo     string `json:"out_trade_no" gorm:"column:out_trade_no;type:varchar(255);index:out_trade_no_index;comment:商户单号"`
	ExternalUserid string `json:"external_userid" gorm:"column:external_userid;type:varchar(255);comment:付款人的userid"`
	ExternalID     int64  `json:"external_id" gorm:"column:external_id;type:bigint(20);default:0;comment:外部联系人用户id"`
	TotalFee       uint   `json:"total_fee" gorm:"column:total_fee;type:int(11) unsigned;default:0;comment:收款总金额，单位为分"`
	PayeeUserid    string `json:"payee_userid" gorm:"column:payee_userid;type:varchar(255);comment:收款人企业内账号userid"`
	MemberID       int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);index:member_id_sales_order_no_index;default:0;comment:导购用户id"`
	PaymentType    int32  `json:"payment_type" gorm:"column:payment_type;type:tinyint(4);default:-1;comment:收款方式。0：在聊天中收款 1：收款码收款 2：在直播间收款 3：用产品图册收款"`
	MchID          string `json:"mch_id" gorm:"column:mch_id;type:varchar(50);comment:收款商户号id"`
	Remark         string `json:"remark" gorm:"column:remark;type:varchar(255);comment:收款备注"`
	CommodityList  string `json:"commodity_list" gorm:"column:commodity_list;type:varchar(255);comment:商品信息详情列表,商品信息结构参考commodity"`
	TotalRefundFee uint   `json:"total_refund_fee" gorm:"column:total_refund_fee;type:int(11) unsigned;default:0;comment:退款总金额"`
	ContactInfo    string `json:"contact_info" gorm:"column:contact_info;type:varchar(255);comment:联系人信息"`
	DeliveryTime   uint   `json:"delivery_time" gorm:"column:delivery_time;type:int(11) unsigned;default:0;comment:发货时间"`
	DeliveryType   int32  `json:"delivery_type" gorm:"column:delivery_type;type:tinyint(4);default:0;comment:订单类型 1自提 2发货"`
	SalesOrderNo   string `json:"sales_order_no" gorm:"column:sales_order_no;type:varchar(255);index:member_id_sales_order_no_index;comment:销售单号发货时候填写的"`
	LogisticsNo    string `json:"logistics_no" gorm:"column:logistics_no;type:varchar(255);comment:物流单号"`
	LogisticsCom   string `json:"logistics_com" gorm:"column:logistics_com;type:varchar(255);comment:物流编码"`
	LogisticsName  string `json:"logistics_name" gorm:"column:logistics_name;type:varchar(255);comment:物流公司"`
	LogisticsTrack string `json:"logistics_track" gorm:"column:logistics_track;type:text;comment:物流轨迹"`
	SignAt         int    `json:"sign_at" gorm:"column:sign_at;type:int(11);default:0;comment:物流签收时间"`
	Tickets        string `json:"tickets" gorm:"column:tickets;type:varchar(255);comment:小票"`
}

// TableName returns the table name of the ExternalPayOrder model
func (ExternalPayOrderModel) TableName() string {
	return "external_pay_order"
}
