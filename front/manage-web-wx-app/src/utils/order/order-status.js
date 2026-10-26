// -----------------------------------------
// 订单状态
// -----------------------------------------

// 未支付
export const NO_PAY = 20
// 已取消
export const CANCELED = 21
// 待支付
export const WAIT_PAY = 23
// 支付中
export const PAYING = 25
// 待确认
export const WAIT_CONFIRM = 30
// 待发货
export const WAIT_DELIVER = 40
// 待收货
export const WAIT_RECEIVED = 60
// 已签收
export const HAD_RECEIVED = 80
// 部分退款
export const REFUND_PART = 83
// 已完成
export const COMPLETE = 100
// 待退款
export const PENDING_REFUND = [41, 81]
// 退款中
export const REFUNDING = 200
// 退款失败
export const REFUND_FAILED = 201
// 全部退款成功
export const REFUND_ALL_SUCCESS = 205
// 售后中
export const AFTER_SALE_STATUS = [41, 61, 81, 82, 84]
// 全部退款
export const ALL_REFUND_COMPLETE = 400


// 订单状态：20:未支付;21:已取消;23:待支付;25:支付中;
// 30:待确认(已支付);40:待发货;41:售后中(待发货,仅退款);
// 60:待收货;61:售后中;80:已签收;81:售后中(已签收,退货退款);
// 83:部分退款;100:已完成;200:退款中;201:退款失败;400:全部退款

// 20 待支付，21 已取消，30 待确认, 40 待发货, 41 售后中(待发货)，
// 60 待收货，80 已签收，81 售后中(已发货)，83 部分退款，
// 100 已完成，200退款中，201 退款失败，205 全部退款成功
