import { request } from '../utils/request'
/**
 * 登录
 * @param {*} data
 * @returns
 */
export const customerLoginApi = data => request({
  url: 'reserve/login',
  method: 'POST'
}, data)
/**
 * 登录
 */
export const getWxlogin = data => request({
  url: 'common/wx-login',
  method: 'POST'
}, data)
