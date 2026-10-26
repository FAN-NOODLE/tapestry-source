import { request } from '../utils/request'

/** 订单列表 */
export const getOrderList = (data) =>
  request(
    {
      url: 'user-external-pay/order-list',
      method: 'POST'
    },
    data
  )

/** 退款列表 */
export const getRefundOrderList = (data) =>
  request(
    {
      url: 'user-external-pay/refund-order-list',
      method: 'POST'
    },
    data
  )

/** 物流信息 */
export const getOrderLogisticsTrack = (data) =>
  request(
    {
      url: 'user-external-pay/order-logistics_track',
      method: 'POST'
    },
    data
  )
