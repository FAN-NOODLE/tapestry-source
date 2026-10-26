import { request } from '../utils/request'

const guideTimeApi = data => request({
  url: 'reserve/user-time-list',
  method: 'POST'
}, data)

const orderInfoApi = data => request({
  url: 'reserve/order-info',
  method: 'POST'
}, data)

const cancelOrderApi = data => request({
  url: 'reserve/cancel-order',
  method: 'POST'
}, data)

const makeScoreApi = data => request({
  url: 'reserve/make-score',
  method: 'POST'
}, data)


const getConfigMessageApi = data => request({
  url: 'reserve/subscribe-message/get-config',
  method: 'POST'
}, data)

// 确认取消 都调这个字段，通过参数区别
const commitOrderApi = data => request({
  url: 'reserve/edit-order',
  method: 'POST'
}, data)

const scoreInfoApi = data => request({
  url: 'reserve/score-info',
  method: 'POST'
}, data)


export {
  guideTimeApi,
  orderInfoApi,
  cancelOrderApi,
  makeScoreApi,
  commitOrderApi,
  getConfigMessageApi,
  scoreInfoApi
}
