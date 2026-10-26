import { request } from '../utils/customer-request'

/** 手机号授权 */
export const getPhoneApi = data => request({
  url: 'activity/reserve/get-phone',
  method: 'POST'
}, data)

/** 获取活动预约时间段 */
export const getTimes = data => request({
  url: 'activity/reserve/time-lists',
  method: 'POST'
}, data)

/** 获取活动信息 */
export const getActivityInfoApi = data => request({
  url: 'activity/reserve/info',
  method: 'POST'
}, data)

/** 创建预约 */
export const createReserveApi = data => request({
  url: 'activity/reserve/create',
  method: 'POST'
}, data)

/** 取消预约 */
export const cancelReserveApi = data => request({
  url: 'activity/reserve/cancel',
  method: 'POST'
}, data)

/** 编辑预约 */
export const editReserveApi = data => request({
  url: 'activity/reserve/edit',
  method: 'POST'
}, data)

