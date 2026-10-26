import { request } from '../utils/request'

/** 我的-当前用户信息 */
export const getMyCenter = data => request({
  url: 'invite/my-center',
  method: 'POST'
}, data)

/** 我的-SA信息 */
export const getMySA = data => request({
  url: 'invite/my-sa',
  method: 'POST'
}, data)
