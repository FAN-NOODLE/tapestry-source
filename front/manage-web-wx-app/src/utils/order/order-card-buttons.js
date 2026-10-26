/**
 * 提示：服务订单-待接订单：只展示接单按钮，数组的其他项都是（服务订单-我的订单）按钮逻辑判断
 */
const buttons = [
  [{ value: 1, text: '确认发货' }],
  [{ value: 10, text: '确认提货' }],
  [{ value: 20, text: '驳回' }, { value: 21, text: '通过审核' }],
  [{ value: 30, text: '转退款审核' }],
  [{ value: 40, text: '驳回' }, { value: 41, text: '退款' }],
  [{ value: 50, text: '转退款审核' }],
  [{ value: 60, text: '转退款审核' }],
  [{ value: 70, text: '接单' }],
  [{ value: 61, text: '添加包裹' }]
]

/**
 * 普通订单
 * 20 待支付，21 已取消，30 待确认, 40 待发货, 41 售后中(待发货)，
 * 60 待收货，80 已签收，81 售后中(已发货)，83 部分退款，
 * 100 已完成，200退款中，201 退款失败，205 全部退款成功
 *
 * 售后
 * 21:待审核;22:待open审核;23:待用户上传物流;24:待open收货;
 * 25:人工中;26:同意退款;27:已发货;28:已取消;29:已拒绝
 *
 * */

/**
 * 初始化按钮
 * order 订单信息
 * showBtn 是否有退款权限
 */
const initButtons = (order, showBtn) => {
  let temp_buttons = []
  // 服务订单-待接订单，只展示接单按钮
  if (order.service_order_type === 1) {
    temp_buttons = buttons[7]
  } else if (!order.is_refund_order) { // 是否是售后订单 is_refund_order 1:是 0:否
    // 快递订单处理
    // 是否是快递单
    if (order.status === 40 && order.delivery_type === 1) {
      // 快递
      temp_buttons = buttons[0]
    } else if (order.status === 60 && order.delivery_type === 1) {
      // 待收货
      temp_buttons = buttons[8]
    } else if (order.status === 40 && order.delivery_type === 2) {
      // 自提
      temp_buttons = buttons[1]
    }
  } else {
    // 售后单处理
    switch (order.refund_status) {
      case 21:
        temp_buttons = buttons[2]
        break
      case 22:
        // 是否需要寄回 refund_apply_type 1: 无需, 2: 需要
        if (order.refund_apply_type === 2) {
          // 不需要寄回
          temp_buttons = buttons[5]
        } else {
          // 需要寄回
          temp_buttons = buttons[3]
        }
        break
      case 24:
        temp_buttons = buttons[6]
        break
      case 25:
        if (showBtn) temp_buttons = buttons[4]
        break
      default:
        break
    }
  }
  return temp_buttons
}
export default initButtons
