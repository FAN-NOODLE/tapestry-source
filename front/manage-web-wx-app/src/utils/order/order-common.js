/**
 * 订单用公用方法
 * */
import {
  NO_PAY,
  PAYING,
  AFTER_SALE_STATUS,
  CANCELED,
  COMPLETE,
  HAD_RECEIVED,
  REFUND_FAILED,
  REFUND_PART,
  REFUNDING,
  WAIT_CONFIRM,
  WAIT_DELIVER,
  WAIT_PAY,
  WAIT_RECEIVED,
  ALL_REFUND_COMPLETE,
  REFUND_ALL_SUCCESS
} from './order-status'

import {
  _WAIT_CHECK,
  _WAIT_CHECK_A,
  _WAIT_CHECK_B,
  _WAIT_CHECK_C,
  _WAIT_REFUND,
  _REFUND_SUCCESS
} from './after-order-status'

/**
 * 订单状态解析
 * */
export const parseOrderStatus = function (status) {
  const ORDER_STATUS = [
    { status: WAIT_PAY, text: '待支付', color: '#ECAB49' },
    { status: PAYING, text: '支付中', color: '#41C56B' },
    { status: NO_PAY, text: '未支付', color: '#808080' },
    { status: WAIT_DELIVER, text: '待发货', color: '#ECAB49' },
    { status: CANCELED, text: '已取消', color: '#808080' },
    { status: WAIT_CONFIRM, text: '待确认', color: '#ECAB49' },
    { status: REFUNDING, text: '退款中', color: '#41C56B' },
    { status: REFUND_PART, text: '部分退款', color: '#41C56B' },
    { status: REFUND_FAILED, text: '退款失败', color: '#EE6352' },
    { status: WAIT_RECEIVED, text: '待收货', color: '#ECAB49' },
    { status: HAD_RECEIVED, text: '已签收', color: '#808080' },
    { status: COMPLETE, text: '已完成', color: '#808080' },
    { status: AFTER_SALE_STATUS, text: '售后中', color: '#41C56B' },
    { status: ALL_REFUND_COMPLETE, text: '全部退款', color: '#41C56B' }, // 全部退款的颜色
    { status: REFUND_ALL_SUCCESS, text: '全部退款成功', color: '#808080' }
    // { status: COMPLETE, text: '已完成', color: '#808080' }// 待售后 待退款
  ]
  let zhText = '--'
  let color = ''
  ORDER_STATUS.forEach(order => {
    if (order.status === parseInt(status)) {
      zhText = order.text
      color = order.color
    } else if (Array.isArray(order.status)) {
      if (order.status.indexOf(parseInt(status)) >= 0) {
        zhText = order.text
        color = order.color
      }
    }
  })
  return { zhText, color }
}
/**
 * 售后状态解析
 * */
export const parseAfterOrderStatus = function (refundStatus) {
  let zhText = '--'
  const AFTER_ORDER_STATUS = [
    { status: _WAIT_CHECK, text: '待审核' },
    { status: _WAIT_CHECK_A, text: '待审核待ERP审核' },
    { status: _WAIT_CHECK_B, text: '待审核待用户寄回商品' },
    { status: _WAIT_CHECK_C, text: '待审核用户已寄出商品' },
    { status: _WAIT_REFUND, text: '待退款' },
    { status: _REFUND_SUCCESS, text: '退款成功' }
  ]

  AFTER_ORDER_STATUS.forEach(item => {
    if (parseInt(refundStatus) === item.status) {
      zhText = item.text
    }
  })
  return zhText
}
