import commonRequest from '../utils/common-request'

/** 订单列表 */
const getOrderList = commonRequest({
  url: 'external-pay/order-list',
  method: 'POST'
})

/** 退款列表 */
const getRefundOrderList = commonRequest({
  url: '/external-pay/refund-order-list',
  method: 'POST'
})

/** 订单详情 */
const getOrderDetail = commonRequest({
  url: '/external-pay/order-detail',
  method: 'POST'
})

/** 退款单订单详情 */
const getRefoundOrderDetail = commonRequest({
  url: '/external-pay/refund-order-detail',
  method: 'POST'
})

/** 快递公司 */
const getDeliveryCompany = commonRequest({
  url: '/external-pay/logistics-company',
  method: 'POST'
})

/** 发货 */
const handleToDelivery = commonRequest({
  url: '/external-pay/deliver-order',
  method: 'POST'
})

/** 上传文件 cos-token */
const cosUploadToken = commonRequest({
  url: '/partner-member/upload',
  method: 'POST'
})

const getOrderLogisticsTrack = commonRequest({
  url: '/external-pay/order-logistics_track',
  method: 'POST'
})

/** 获取线上订单配置 */
const getOnlineOrderPage = commonRequest({
  url: 'store/d1m-order-h5',
  method: 'POST'
})

/** 是否展示会员查询 */
const customWhiteListApi = commonRequest({
  url: 'customer/white-list',
  method: 'POST'
})

export {
  getOrderList,
  getRefundOrderList,
  getOrderDetail,
  getRefoundOrderDetail,
  getDeliveryCompany,
  handleToDelivery,
  cosUploadToken,
  getOrderLogisticsTrack,
  getOnlineOrderPage,
  customWhiteListApi
}
