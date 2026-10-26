import commonRequest from '../utils/common-request'

// 获取导购当前导购时间
const guideTimeApi = commonRequest({
  url: 'reserve/time-list',
  method: 'POST'
})
// 创建预约单reserve/create-order
const createOrderApi = commonRequest({
  url: 'reserve/create-order',
  method: 'POST'
})
// 获取用户信息
const userInfoApi = commonRequest({
  url: 'reserve/user-simple-info',
  method: 'POST'
})
// 获取门店信息
const storeInfoApi = commonRequest({
  url: 'reserve/member-info',
  method: 'POST'
})

// 获取服务列表
const serviceListApi = commonRequest({
  url: 'reserve/service-list',
  method: 'POST'
})

// 邀约记录列表
const orderListApi = commonRequest({
  url: 'reserve/order-list',
  method: 'POST'
})

const orderDetailApi = commonRequest({
  url: 'reserve/order-detail',
  method: 'POST'
})
// 分享次数reserve/
const orderShareNumApi = commonRequest({
  url: 'reserve/order-share',
  method: 'POST'
})
// 订单编辑
const orderEditApi = commonRequest({
  url: 'reserve/order-edit',
  method: 'POST'
})
// 取消订单
const orderCancelApi = commonRequest({
  url: 'reserve/order-cancel',
  method: 'POST'
})
// 确认顾客到店
const orderConfirmStore = commonRequest({
  url: 'reserve/confirm-store',
  method: 'POST'
})
// (店长)获取导购列表
const getStaffLists = commonRequest({
  url: 'member/staff-lists',
  method: 'POST'
})
// 获取门店列表
const getStoreLists = commonRequest({
  url: 'member/store-lists',
  method: 'POST'
})
// 是否到店购买 & 选择顾客类型
const changeStatus = commonRequest({
  url: 'reserve/change-status',
  method: 'POST'
})
// 发送评价记录日志
const commentShareApi = commonRequest({
  url: 'reserve/comment-share',
  method: 'POST'
})

// 发送评价记录日志
const verifyShareApi = commonRequest({
  url: 'reserve/verify-share',
  method: 'POST'
})
// 顾客画像邀约列表接口
const inviteListApi = commonRequest({
  url: 'reserve/customer-reserve-lists',
  method: 'POST'
})

export {
  inviteListApi,
  guideTimeApi,
  createOrderApi,
  userInfoApi,
  storeInfoApi,
  serviceListApi,
  orderListApi,
  orderDetailApi,
  orderShareNumApi,
  orderEditApi,
  orderCancelApi,
  orderConfirmStore,
  getStaffLists,
  getStoreLists,
  changeStatus,
  commentShareApi,
  verifyShareApi
}
