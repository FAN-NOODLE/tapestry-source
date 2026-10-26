import { request } from '../utils/request'

// 获取导购当前导购时间
const guideTimeApi = data => request({
  url: 'reserve/time-list',
  method: 'POST'
}, data)
// 创建预约单reserve/create-order
const createOrderApi = data => request({
  url: 'reserve/create-order',
  method: 'POST'
}, data)
// 获取用户信息
const userInfoApi = data => request({
  url: 'reserve/user-simple-info',
  method: 'POST'
}, data)

// 获取预约列表
const orderListApi = data => request({
  url: 'reserve/user-order-lists',
  method: 'POST'
}, data)

const orderDetailApi = data => request({
  url: 'reserve/order-info',
  method: 'POST'
}, data)

const orderAssignApi = data => request({
  url: 'reserve/order-assign',
  method: 'POST'
}, data)

// 预约单编辑
const orderEditApi = data => request({
  url: 'reserve/edit-order',
  method: 'POST'
}, data)

const orderFinishApi = data => request({
  url: 'reserve/order-finish',
  method: 'POST'
}, data)

const writeOffApi = data => request({
  url: 'reserve/order-writeOff',
  method: 'POST'
}, data)

const getConfigMessageApi = data => request({
  url: 'reserve/subscribe-message/get-config',
  method: 'POST'
}, data)

// 消息订阅
const orderNotifyApi = data => request({
  url: 'reserve/subscribe-message/subscribe',
  method: 'POST'
}, data)

export {
  guideTimeApi,
  createOrderApi,
  userInfoApi,
  orderListApi,
  orderDetailApi,
  orderAssignApi,
  orderEditApi,
  orderFinishApi,
  writeOffApi,
  getConfigMessageApi,
  orderNotifyApi
}
